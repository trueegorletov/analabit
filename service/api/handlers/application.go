package handlers

import (
	"analabit/core"
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

		// Collect student and heading IDs for batch queries
		studentIDs := make([]string, 0, len(applications))
		headingIDs := make([]int, 0, len(applications))
		for _, app := range applications {
			studentIDs = append(studentIDs, app.StudentID)
			headingIDs = append(headingIDs, app.Edges.Heading.ID)
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
				allHeadingIDs = append(allHeadingIDs, app.Edges.Heading.ID)
			}
		}

		// Fetch passing calculation results for all relevant headings
		passingCalculations, err := client.Calculation.Query().
			Where(
				calculation.RunIDEQ(runResolution.RunID),
				calculation.HasHeadingWith(heading.IDIn(allHeadingIDs...)),
			).
			All(ctx)
		if err != nil {
			log.Printf("error getting passing calculations: %v", err)
			return fiber.ErrInternalServerError
		}

		// Create a map for quick lookup of passing students
		passingStudents := make(map[int]map[string]struct{})
		for _, calc := range passingCalculations {
			calcHeading, err := calc.QueryHeading().Only(ctx)
			if err != nil {
				log.Printf("error getting heading for calculation: %v", err)
				return fiber.ErrInternalServerError
			}
			if _, ok := passingStudents[calcHeading.ID]; !ok {
				passingStudents[calcHeading.ID] = make(map[string]struct{})
			}
			passingStudents[calcHeading.ID][calc.StudentID] = struct{}{}
		}

		// Build the response
		response := make([]ApplicationResponse, len(applications))
		for i, app := range applications {
			studentApps := appsByStudent[app.StudentID]

			// Calculate original_quit
			originalQuit := false
			if !app.OriginalSubmitted {
				for _, otherApp := range studentApps {
					if otherApp.ID != app.ID && otherApp.OriginalSubmitted {
						originalQuit = true
						break
					}
				}
			}

			// Calculate passing_now and passing_to_more_priority
			passingNow := false
			passingToMorePriority := false
			if headingPassing, ok := passingStudents[app.Edges.Heading.ID]; ok {
				if _, ok := headingPassing[app.StudentID]; ok {
					passingNow = true
				}
			}

			if !passingNow {
				for _, otherApp := range studentApps {
					if otherApp.Edges.Heading.Edges.Varsity.ID == app.Edges.Heading.Edges.Varsity.ID && otherApp.Priority < app.Priority {
						otherAppHeading, err := otherApp.QueryHeading().Only(ctx)
						if err != nil {
							log.Printf("error getting heading for otherApp: %v", err)
							return fiber.ErrInternalServerError
						}
						if headingPassing, ok := passingStudents[otherAppHeading.ID]; ok {
							if _, ok := headingPassing[app.StudentID]; ok {
								passingToMorePriority = true
								break
							}
						}
					}
				}
			}

			// Calculate another_varsities_count
			varsitySet := make(map[int]struct{})
			for _, studentApp := range studentApps {
				varsitySet[studentApp.Edges.Heading.Edges.Varsity.ID] = struct{}{}
			}
			delete(varsitySet, app.Edges.Heading.Edges.Varsity.ID)
			anotherVarsitiesCount := len(varsitySet)

			response[i] = ApplicationResponse{
				ID:                    app.ID,
				StudentID:             app.StudentID,
				Priority:              app.Priority,
				CompetitionType:       core.Competition(app.CompetitionType).String(),
				RatingPlace:           app.RatingPlace,
				Score:                 app.Score,
				RunID:                 app.RunID,
				UpdatedAt:             app.UpdatedAt,
				HeadingID:             app.Edges.Heading.ID,
				OriginalSubmitted:     app.OriginalSubmitted,
				OriginalQuit:          originalQuit,
				PassingNow:            passingNow,
				PassingToMorePriority: passingToMorePriority,
				AnotherVarsitiesCount: anotherVarsitiesCount,
			}
		}

		return c.JSON(response)
	}
}
