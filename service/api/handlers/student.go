package handlers

import (
	"github.com/trueegorletov/analabit/core/ent"
	"github.com/trueegorletov/analabit/core/ent/application"
	"context"
	"log"

	"github.com/gofiber/fiber/v3"
)

// GetStudentByID retrieves information about a student's applications.
func GetStudentByID(client *ent.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		studentID := c.Params("id")

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

		return c.JSON(applications)
	}
}
