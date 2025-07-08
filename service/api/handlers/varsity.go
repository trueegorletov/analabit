package handlers

import (
	"github.com/trueegorletov/analabit/core/ent"
	"context"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func GetVarsities(client *ent.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		limit, _ := strconv.Atoi(c.Query("limit", "100"))
		offset, _ := strconv.Atoi(c.Query("offset", "0"))

		varsities, err := client.Varsity.
			Query().
			Limit(limit).
			Offset(offset).
			All(context.Background())

		if err != nil {
			log.Printf("error getting varsities: %v", err)
			return fiber.ErrInternalServerError
		}

		return c.JSON(varsities)
	}
}
