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
	"time"

	"github.com/gofiber/fiber/v3"
)

// ApplicationResponse enriches Application with additional flags.
type ApplicationResponse struct {
	ID                int       `json:"id"`
	StudentID         string    `json:"student_id"`
	Priority          int       `json:"priority"`
	CompetitionType   int       `json:"competition_type"`
	RatingPlace       int       `json:"rating_place"`
	Score             int       `json:"score"`
	RunID             int       `json:"run_id"`
	UpdatedAt         time.Time `json:"updated_at"`
	HeadingID         int       `json:"heading_id"`
	OriginalSubmitted bool      `json:"original_submitted"`
	OriginalQuit      bool      `json:"original_quit"`
	PassingNow        bool      `json:"passing_now"`
}

// GetApplications retrieves a list of applications, with optional filtering.
func GetApplications(client *ent.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		limit, _ := strconv.Atoi(c.Query("limit", "1000"))
		offset, _ := strconv.Atoi(c.Query("offset", "0"))
		studentID := c.Query("studentID")
		varsityCode := c.Query("varsityCode")
		headingID, _ := strconv.Atoi(c.Query("headingId", "0"))
		runParam := c.Query("run", "latest")

		ctx := context.Background()

		// Resolve the run ID from the parameter
		runResolution, err := ResolveRunFromIteration(ctx, client, runParam)
		if err != nil {
			log.Printf("error resolving run from parameter '%s': %v", runParam, err)
			return fiber.NewError(fiber.StatusBadRequest, "invalid run parameter")
		}

		q := client.Application.Query().Where(application.RunIDEQ(runResolution.RunID))

		if studentID != "" {
			q = q.Where(application.StudentID(studentID))
		}

		if varsityCode != "" {
			q = q.Where(application.HasHeadingWith(heading.HasVarsityWith(varsity.CodeEQ(varsityCode))))
		}

		if headingID > 0 {
			q = q.Where(application.HasHeadingWith(heading.ID(headingID)))
		}

		applications, err := q.WithHeading(func(hq *ent.HeadingQuery) { hq.WithVarsity() }).WithRun().Order(application.ByRatingPlace()).Limit(limit).Offset(offset).All(ctx)
		if err != nil {
			log.Printf("error getting applications: %v", err)
			return fiber.ErrInternalServerError
		}

		// Build helper maps ---------------------------------

		studentSet := make(map[string]struct{})
		headingSet := make(map[int]struct{})

		for _, app := range applications {
			studentSet[app.StudentID] = struct{}{}
			headingSet[app.Edges.Heading.ID] = struct{}{}
		}

		// Get latest run's calculations for comparison (to determine passingNow)
		latestRunID, err := getLatestRunID(ctx, client)
		if err != nil {
			log.Printf("error getting latest run ID: %v", err)
			return fiber.ErrInternalServerError
		}

		// Fetch calculations at latest run for involved students & headings
		var studentIDs []string
		for id := range studentSet {
			studentIDs = append(studentIDs, id)
		}
		var headingIDs []int
		for id := range headingSet {
			headingIDs = append(headingIDs, id)
		}

		calcMap := make(map[string]struct{}) // key studentID|headingID
		if latestRunID > 0 && len(studentIDs) > 0 && len(headingIDs) > 0 {
			calcs, _ := client.Calculation.Query().Where(
				calculation.RunIDEQ(latestRunID),
				calculation.StudentIDIn(studentIDs...),
				calculation.HasHeadingWith(heading.IDIn(headingIDs...)),
			).WithHeading().All(ctx)
			for _, ccalc := range calcs {
				key := ccalc.StudentID + "|" + strconv.Itoa(ccalc.Edges.Heading.ID)
				calcMap[key] = struct{}{}
			}
		}

		// Fetch originalSubmitted-true apps to help compute originalQuit
		// We need to check applications from the same run for originalQuit logic
		origApps, _ := client.Application.Query().Where(
			application.OriginalSubmitted(true),
			application.RunIDEQ(runResolution.RunID),
			application.StudentIDIn(studentIDs...),
		).WithHeading(func(hq *ent.HeadingQuery) { hq.WithVarsity() }).All(ctx)

		// Build map studentID -> set of varsity IDs where originalSubmitted true (for same run)
		origMap := make(map[string]map[int]struct{})
		for _, oa := range origApps {
			key := oa.StudentID
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

			// passingNow - check if student is in calculations for latest run
			if _, ok := calcMap[app.StudentID+"|"+strconv.Itoa(app.Edges.Heading.ID)]; ok {
				passingNow = true
			}

			// originalQuit (only if originalSubmitted false)
			if !app.OriginalSubmitted {
				key := app.StudentID
				if vars, ok := origMap[key]; ok {
					if _, inSameVarsity := vars[app.Edges.Heading.Edges.Varsity.ID]; !inSameVarsity && len(vars) > 0 {
						originalQuit = true
					}
				}
			}

			runID := 0
			if app.Edges.Run != nil {
				runID = app.Edges.Run.ID
			}

			resp[i] = ApplicationResponse{
				ID:                app.ID,
				StudentID:         app.StudentID,
				Priority:          app.Priority,
				CompetitionType:   int(app.CompetitionType),
				RatingPlace:       app.RatingPlace,
				Score:             app.Score,
				RunID:             runID,
				UpdatedAt:         app.UpdatedAt,
				HeadingID:         app.Edges.Heading.ID,
				OriginalSubmitted: app.OriginalSubmitted,
				OriginalQuit:      originalQuit,
				PassingNow:        passingNow,
			}
		}

		return c.JSON(resp)
	}
}
