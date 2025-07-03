package handlers

import (
	"analabit/core/ent"
	"analabit/core/ent/application"
	"analabit/core/ent/calculation"
	"analabit/core/ent/heading"
	"analabit/core/ent/varsity"
	"context"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// ApplicationResponse enriches Application with additional flags.
type ApplicationResponse struct {
	*ent.Application  `json:"-"`
	OriginalSubmitted bool `json:"original_submitted"`
	OriginalQuit      bool `json:"original_quit"`
	PassingNow        bool `json:"passing_now"`
}

// GetApplications retrieves a list of applications, with optional filtering.
func GetApplications(client *ent.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		limit, _ := strconv.Atoi(c.Query("limit", "100"))
		offset, _ := strconv.Atoi(c.Query("offset", "0"))
		studentID := c.Query("studentID")
		varsityCode := c.Query("varsityCode")
		headingID, _ := strconv.Atoi(c.Query("headingId", "0"))

		q := client.Application.Query()

		if studentID != "" {
			q = q.Where(application.StudentID(studentID))
		}

		if varsityCode != "" {
			q = q.Where(application.HasHeadingWith(heading.HasVarsityWith(varsity.CodeEQ(varsityCode))))
		}

		if headingID > 0 {
			q = q.Where(application.HasHeadingWith(heading.ID(headingID)))
		}

		ctx := context.Background()

		applications, err := q.WithHeading(func(hq *ent.HeadingQuery) { hq.WithVarsity() }).Limit(limit).Offset(offset).All(ctx)
		if err != nil {
			log.Printf("error getting applications: %v", err)
			return fiber.ErrInternalServerError
		}

		// Build helper maps ---------------------------------

		studentSet := make(map[string]struct{})
		headingSet := make(map[int]struct{})
		iterationsSet := make(map[int]struct{})

		for _, app := range applications {
			studentSet[app.StudentID] = struct{}{}
			headingSet[app.Edges.Heading.ID] = struct{}{}
			iterationsSet[app.Iteration] = struct{}{}
		}

		// Get last calc iteration
		var v []struct {
			Max int `json:"max"`
		}
		if err := client.Calculation.Query().Aggregate(ent.Max(calculation.FieldIteration)).Scan(ctx, &v); err != nil {
			log.Printf("error getting calc iter: %v", err)
			return fiber.ErrInternalServerError
		}
		lastCalcIter := 0
		if len(v) > 0 {
			lastCalcIter = v[0].Max
		}

		// Fetch calculations at last iteration for involved students & headings
		var studentIDs []string
		for id := range studentSet {
			studentIDs = append(studentIDs, id)
		}
		var headingIDs []int
		for id := range headingSet {
			headingIDs = append(headingIDs, id)
		}

		calcMap := make(map[string]struct{}) // key studentID|headingID
		if lastCalcIter > 0 {
			calcs, _ := client.Calculation.Query().Where(
				calculation.IterationEQ(lastCalcIter),
				calculation.StudentIDIn(studentIDs...),
				calculation.HasHeadingWith(heading.IDIn(headingIDs...)),
			).WithHeading().All(ctx)
			for _, ccalc := range calcs {
				key := ccalc.StudentID + "|" + strconv.Itoa(ccalc.Edges.Heading.ID)
				calcMap[key] = struct{}{}
			}
		}

		// Fetch originalSubmitted-true apps to help compute originalQuit
		var iterations []int
		for it := range iterationsSet {
			iterations = append(iterations, it)
		}
		origApps, _ := client.Application.Query().Where(
			application.OriginalSubmitted(true),
			application.StudentIDIn(studentIDs...),
			application.IterationIn(iterations...),
		).WithHeading(func(hq *ent.HeadingQuery) { hq.WithVarsity() }).All(ctx)

		// Build map studentID|iteration -> set of varsity IDs where originalSubmitted true
		origMap := make(map[string]map[int]struct{})
		for _, oa := range origApps {
			key := oa.StudentID + "|" + strconv.Itoa(oa.Iteration)
			if _, ok := origMap[key]; !ok {
				origMap[key] = make(map[int]struct{})
			}
			if oa.Edges.Heading != nil && oa.Edges.Heading.Edges.Varsity != nil {
				origMap[key][oa.Edges.Heading.Edges.Varsity.ID] = struct{}{}
			}
		}

		// Compose response ----------------------------------
		resp := make([]ApplicationResponse, len(applications))

		for i, app := range applications {
			var passingNow bool
			var originalQuit bool

			// passingNow
			if _, ok := calcMap[app.StudentID+"|"+strconv.Itoa(app.Edges.Heading.ID)]; ok {
				passingNow = true
			}

			// originalQuit (only if originalSubmitted false)
			if !app.OriginalSubmitted {
				key := app.StudentID + "|" + strconv.Itoa(app.Iteration)
				if vars, ok := origMap[key]; ok {
					if _, inSameVarsity := vars[app.Edges.Heading.Edges.Varsity.ID]; !inSameVarsity && len(vars) > 0 {
						originalQuit = true
					}
				}
			}

			resp[i] = ApplicationResponse{
				Application:       app,
				OriginalSubmitted: app.OriginalSubmitted,
				OriginalQuit:      originalQuit,
				PassingNow:        passingNow,
			}
		}

		return c.JSON(resp)
	}
}
