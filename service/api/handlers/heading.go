package handlers

import (
	"analabit/core/ent"
	"analabit/core/ent/heading"
	"analabit/core/ent/varsity"
	"context"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// GetHeadings retrieves a list of headings, with optional filtering.
func GetHeadings(client *ent.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		limit, _ := strconv.Atoi(c.Query("limit", "100"))
		offset, _ := strconv.Atoi(c.Query("offset", "0"))
		varsityCode := c.Query("varsityCode")

		q := client.Heading.Query()

		if varsityCode != "" {
			q = q.Where(heading.HasVarsityWith(varsity.CodeEQ(varsityCode)))
		}

		// preload varsity to access code
		headings, err := q.WithVarsity().Limit(limit).Offset(offset).All(context.Background())
		if err != nil {
			log.Printf("error getting headings: %v", err)
			return fiber.ErrInternalServerError
		}

		// map to DTOs
		resp := make([]HeadingResponse, len(headings))
		for i, h := range headings {
			var vDTO VarsityDTO
			if v := h.Edges.Varsity; v != nil {
				vDTO = VarsityDTO{ID: v.ID, Code: v.Code, Name: v.Name}
			}
			resp[i] = HeadingResponse{
				ID:                     h.ID,
				Code:                   h.Code,
				Name:                   h.Name,
				RegularCapacity:        h.RegularCapacity,
				TargetQuotaCapacity:    h.TargetQuotaCapacity,
				DedicatedQuotaCapacity: h.DedicatedQuotaCapacity,
				SpecialQuotaCapacity:   h.SpecialQuotaCapacity,
				Varsity:                vDTO,
			}
		}

		return c.JSON(resp)
	}
}

// GetHeadingByID retrieves a single heading by its ID, along with its admission results.
func GetHeadingByID(client *ent.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid heading ID")
		}

		h, err := client.Heading.Query().Where(heading.ID(id)).WithVarsity().Only(context.Background())
		if err != nil {
			log.Printf("error getting heading: %v", err)
			return fiber.ErrNotFound
		}

		var vDTO VarsityDTO
		if v := h.Edges.Varsity; v != nil {
			vDTO = VarsityDTO{ID: v.ID, Code: v.Code, Name: v.Name}
		}

		resp := HeadingResponse{
			ID:                     h.ID,
			Code:                   h.Code,
			Name:                   h.Name,
			RegularCapacity:        h.RegularCapacity,
			TargetQuotaCapacity:    h.TargetQuotaCapacity,
			DedicatedQuotaCapacity: h.DedicatedQuotaCapacity,
			SpecialQuotaCapacity:   h.SpecialQuotaCapacity,
			Varsity:                vDTO,
		}

		return c.JSON(resp)
	}
}

// ADD: heading response DTO

type VarsityDTO struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type HeadingResponse struct {
	ID                     int        `json:"id"`
	Code                   string     `json:"code"`
	Name                   string     `json:"name"`
	RegularCapacity        int        `json:"regular_capacity"`
	TargetQuotaCapacity    int        `json:"target_quota_capacity"`
	DedicatedQuotaCapacity int        `json:"dedicated_quota_capacity"`
	SpecialQuotaCapacity   int        `json:"special_quota_capacity"`
	Varsity                VarsityDTO `json:"varsity"`
}
