package handlers

import (
	"context"
	"log"
	"log/slog"
	"strconv"
	"time"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/ent"
	"github.com/trueegorletov/analabit/core/ent/application"
	"github.com/trueegorletov/analabit/core/ent/calculation"
	"github.com/trueegorletov/analabit/core/ent/heading"
	"github.com/trueegorletov/analabit/core/ent/varsity"
	"github.com/trueegorletov/analabit/core/utils"

	"github.com/gofiber/fiber/v3"
)

// ApplicationResponse enriches Application with additional flags.
type ApplicationResponse struct {
	ID                    int       `json:"id"`
	StudentID             string    `json:"student_id"`
	Priority              int       `json:"priority"`
	CompetitionType       string    `json:"competition_type"`
	RatingPlace           int       `json:"rating_place"`
	Score                 int       `json:"score"`
	RunID                 int       `json:"run_id"`
	UpdatedAt             time.Time `json:"updated_at"`
	HeadingID             int       `json:"heading_id"`
	OriginalSubmitted     bool      `json:"original_submitted"`
	OriginalQuit          bool      `json:"original_quit"`
	PassingNow            bool      `json:"passing_now"`
	PassingToMorePriority bool      `json:"passing_to_more_priority"`
	AnotherVarsitiesCount int       `json:"another_varsities_count"`
}

// GetApplications retrieves a list of applications, with optional filtering.
func GetApplications(client *ent.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		limit, _ := strconv.Atoi(c.Query("limit", "1000"))
		offset, _ := strconv.Atoi(c.Query("offset", "0"))
		studentIDRaw := c.Query("studentID")
		varsityCode := c.Query("varsityCode")
		headingID, _ := strconv.Atoi(c.Query("headingId", "0"))
		runParam := c.Query("run", "latest")

		ctx := context.Background()

		// Validate and prepare student ID if provided
		var studentID string
		if studentIDRaw != "" {
			preparedID, err := utils.PrepareStudentID(studentIDRaw)
			if err != nil {
				log.Printf("invalid student ID parameter '%s': %v", studentIDRaw, err)
				return fiber.NewError(fiber.StatusBadRequest, "invalid student ID parameter")
			}
			studentID = preparedID
		}

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

		// Collect student and heading IDs for batch queries
		studentIDsSet := make(map[string]struct{}, len(applications))
		studentIDs := make([]string, 0, len(applications))
		for _, app := range applications {
			if _, exists := studentIDsSet[app.StudentID]; !exists {
				studentIDs = append(studentIDs, app.StudentID)
				studentIDsSet[app.StudentID] = struct{}{}
			}
		}

		// Fetch all applications for the students in this run to calculate related fields
		allStudentApplications, err := client.Application.Query().
			Where(
				application.RunIDEQ(runResolution.RunID),
				application.StudentIDIn(studentIDs...),
			).
			WithHeading(func(hq *ent.HeadingQuery) {
				hq.WithVarsity()
			}).
			All(ctx)
		if err != nil {
			log.Printf("error getting all student applications: %v", err)
			return fiber.ErrInternalServerError
		}

		// Group applications by student ID for quick lookup
		appsByStudent := make(map[string][]*ent.Application)
		for _, app := range allStudentApplications {
			appsByStudent[app.StudentID] = append(appsByStudent[app.StudentID], app)
		}

		// Collect all heading IDs from all applications of the students on this page
		allHeadingIDs := make([]int, 0)
		for _, studentApps := range appsByStudent {
			for _, app := range studentApps {
				if app.Edges.Heading == nil {
					slog.Error("application has no heading, skipping", "appID", app.ID)
					continue // Skip applications without a heading
				}

				allHeadingIDs = append(allHeadingIDs, app.Edges.Heading.ID)
			}
		}

		// Fetch passing calculation results for all relevant headings
		passingCalculations, err := client.Calculation.Query().
			Where(
				calculation.RunIDEQ(runResolution.RunID),
				calculation.HasHeadingWith(heading.IDIn(allHeadingIDs...)),
			).
			WithHeading(func(hq *ent.HeadingQuery) { hq.WithVarsity() }).
			All(ctx)
		if err != nil {
			log.Printf("error getting passing calculations: %v", err)
			return fiber.ErrInternalServerError
		}

		type studentPassingInfo struct {
			appHeadingID int
			appPriority  int
		}

		saferAppHeadingID := func(app *ent.Application) int {
			if app.Edges.Heading == nil {
				slog.Error("application has no heading, returning 0 for safer handling", "appID", app.ID)
				return 0 // No heading, return 0 to avoid panic
			}
			return app.Edges.Heading.ID
		}
		saferHeadingVarsityID := func(heading *ent.Heading) int {
			if heading.Edges.Varsity == nil {
				slog.Error("heading has no varsity, returning 0 for safer handling", "headingID", heading.ID)
				return 0 // No varsity, return 0 to avoid panic
			}
			return heading.Edges.Varsity.ID
		}
		saferAppVarsityID := func(app *ent.Application) int {
			if app.Edges.Heading == nil {
				slog.Error("application has no heading, returning 0 for safer handling", "appID", app.ID)
				return 0 // No heading or varsity, return 0 to avoid panic
			}

			if app.Edges.Heading.Edges.Varsity == nil {
				slog.Error("application heading has no varsity, returning 0 for safer handling", "appID", app.ID, "headingID", app.Edges.Heading.ID)
				return 0 // No varsity, return 0 to avoid panic
			}

			return app.Edges.Heading.Edges.Varsity.ID
		}

		// Create a map for quick lookup of passing students
		passingStudents := make(map[int]map[string]struct{})

		studentsPassingInfosByVarsity := make(map[string]map[int]studentPassingInfo)

		for _, calc := range passingCalculations {
			calcHeading := calc.Edges.Heading
			if calcHeading == nil {
				slog.Error("calculation has no heading, skipping", "calcID", calc.ID)
				continue
			}
			if _, ok := passingStudents[calcHeading.ID]; !ok {
				passingStudents[calcHeading.ID] = make(map[string]struct{})
			}
			passingStudents[calcHeading.ID][calc.StudentID] = struct{}{}

			varsityID := saferHeadingVarsityID(calcHeading)

			if _, ok := studentsPassingInfosByVarsity[calc.StudentID]; !ok {
				studentsPassingInfosByVarsity[calc.StudentID] = make(map[int]studentPassingInfo)
			}

			if _, ok := studentsPassingInfosByVarsity[calc.StudentID][varsityID]; !ok {
				studentsPassingInfosByVarsity[calc.StudentID][varsityID] = studentPassingInfo{
					appHeadingID: calcHeading.ID,
					appPriority:  -1,
				}
			}
		}

		origVarsityIDByStudent := make(map[string]int)
		varsitiesCountByStudent := make(map[string]int)

		for id, apps := range appsByStudent {
			varsities := make(map[int]struct{})

			for _, app := range apps {
				appVarsityID := saferAppVarsityID(app)

				varsities[appVarsityID] = struct{}{}

				if _, ok := origVarsityIDByStudent[id]; !ok {
					if app.OriginalSubmitted {
						origVarsityIDByStudent[id] = appVarsityID
					}
				}

				if passingInfo, ok := studentsPassingInfosByVarsity[id][appVarsityID]; ok {
					if passingInfo.appHeadingID == saferAppHeadingID(app) {
						passingInfo.appPriority = app.Priority
						studentsPassingInfosByVarsity[id][appVarsityID] = passingInfo
					}
				}
			}

			varsitiesCountByStudent[id] = len(varsities)
		}

		// Build the response
		response := make([]ApplicationResponse, len(applications))
		for i, app := range applications {
			originalSubmitted, originalQuit := false, false

			appVarsityID := saferAppVarsityID(app)
			appHeadingID := saferAppHeadingID(app)

			studentOriginalVarsityID, originalDetermined := origVarsityIDByStudent[app.StudentID]

			if originalDetermined {
				if appVarsityID == studentOriginalVarsityID {
					originalSubmitted = true
					originalQuit = false
				} else {
					originalSubmitted = false
					originalQuit = true
				}
			}

			// Calculate passing_now and passing_to_more_priority
			passingNow := false
			passingToMorePriority := false

			if passingInfo, ok := studentsPassingInfosByVarsity[app.StudentID][appVarsityID]; ok {
				if passingInfo.appPriority == app.Priority || passingInfo.appHeadingID == appHeadingID {
					passingNow = true
					passingToMorePriority = false
				} else if passingInfo.appPriority < app.Priority {
					passingNow = false
					passingToMorePriority = true
				}
			}

			anotherVarsitiesCount := varsitiesCountByStudent[app.StudentID] - 1

			response[i] = ApplicationResponse{
				ID:                    app.ID,
				StudentID:             utils.PrettifyStudentID(app.StudentID),
				Priority:              app.Priority,
				CompetitionType:       core.Competition(app.CompetitionType).String(),
				RatingPlace:           app.RatingPlace,
				Score:                 app.Score,
				RunID:                 app.RunID,
				UpdatedAt:             app.UpdatedAt,
				HeadingID:             appHeadingID,
				OriginalSubmitted:     originalSubmitted,
				OriginalQuit:          originalQuit,
				PassingNow:            passingNow,
				PassingToMorePriority: passingToMorePriority,
				AnotherVarsitiesCount: anotherVarsitiesCount,
			}
		}

		return c.JSON(response)
	}
}
