package handlers

import (
	"context"
	"log"
	"time"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/ent"
	"github.com/trueegorletov/analabit/core/ent/application"
	"github.com/trueegorletov/analabit/core/utils"

	"github.com/gofiber/fiber/v3"
)

// StudentApplicationResponse represents an application with prettified student ID
type StudentApplicationResponse struct {
	ID                int              `json:"id"`
	StudentID         string           `json:"student_id"`
	Priority          int              `json:"priority"`
	CompetitionType   core.Competition `json:"competition_type"`
	RatingPlace       int              `json:"rating_place"`
	Score             int              `json:"score"`
	RunID             int              `json:"run_id"`
	UpdatedAt         time.Time        `json:"updated_at"`
	OriginalSubmitted bool             `json:"original_submitted"`
	Heading           *ent.Heading     `json:"heading"`
}

// GetStudentByID retrieves information about a student's applications.
func GetStudentByID(client *ent.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		studentIDRaw := c.Params("id")

		// Validate and prepare student ID
		studentID, err := utils.PrepareStudentID(studentIDRaw)
		if err != nil {
			log.Printf("invalid student ID parameter '%s': %v", studentIDRaw, err)
			return fiber.NewError(fiber.StatusBadRequest, "invalid student ID parameter")
		}

		applications, err := client.Application.
			Query().
			Where(application.StudentID(studentID)).
			WithHeading(func(q *ent.HeadingQuery) {
				q.WithVarsity()
			}).
			All(context.Background())

		if err != nil {
			log.Printf("error getting student applications: %v", err)
			return fiber.ErrInternalServerError
		}

		if len(applications) == 0 {
			return fiber.NewError(fiber.StatusNotFound, "Student not found")
		}

		// Transform applications to prettify student IDs
		response := make([]StudentApplicationResponse, len(applications))
		for i, app := range applications {
			response[i] = StudentApplicationResponse{
				ID:                app.ID,
				StudentID:         utils.PrettifyStudentID(app.StudentID),
				Priority:          app.Priority,
				CompetitionType:   app.CompetitionType,
				RatingPlace:       app.RatingPlace,
				Score:             app.Score,
				RunID:             app.RunID,
				UpdatedAt:         app.UpdatedAt,
				OriginalSubmitted: app.OriginalSubmitted,
				Heading:           app.Edges.Heading,
			}
		}

		return c.JSON(response)
	}
}
