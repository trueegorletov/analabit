package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/trueegorletov/analabit/core/ent"
	"github.com/trueegorletov/analabit/core/ent/application"
	"github.com/trueegorletov/analabit/core/ent/heading"
	"github.com/trueegorletov/analabit/core/ent/varsity"
	"github.com/trueegorletov/analabit/core/utils"

	"encoding/base64"
	"strings"

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

// GetApplications retrieves a list of applications with cursor-based pagination.
func GetApplications(client *ent.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := context.Background()

		// Parse query parameters
		studentIDRaw := c.Query("studentID")
		varsityCode := c.Query("varsityCode")
		headingID, _ := strconv.Atoi(c.Query("headingId", "0"))
		runParam := c.Query("run", "latest")
		first, err := strconv.Atoi(c.Query("first", "100"))
		if err != nil || first <= 0 {
			first = 100
		}
		after := c.Query("after")

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

		// Base query
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

		// Order: by rating place ASC (assuming lower is better), then by ID ASC for stability
		q = q.Order(ent.Asc(application.FieldRatingPlace), ent.Asc(application.FieldID))

		// Handle cursor
		limit := first + 1 // fetch one extra to check hasNextPage
		if after != "" {
			decoded, err := base64.StdEncoding.DecodeString(after)
			if err != nil {
				return fiber.NewError(fiber.StatusBadRequest, "invalid cursor")
			}
			parts := strings.Split(string(decoded), ":")
			if len(parts) != 2 {
				return fiber.NewError(fiber.StatusBadRequest, "invalid cursor format")
			}
			lastRatingPlace, err := strconv.Atoi(parts[0])
			if err != nil {
				return fiber.NewError(fiber.StatusBadRequest, "invalid cursor")
			}
			lastID, err := strconv.Atoi(parts[1])
			if err != nil {
				return fiber.NewError(fiber.StatusBadRequest, "invalid cursor")
			}
			// Predicate: (rating_place > last) OR (rating_place == last AND id > lastID)
			q = q.Where(
				func(s *sql.Selector) {
					s.Where(sql.Or(
						sql.GT(s.C(application.FieldRatingPlace), lastRatingPlace),
						sql.And(
							sql.EQ(s.C(application.FieldRatingPlace), lastRatingPlace),
							sql.GT(s.C(application.FieldID), lastID),
						),
					))
				},
			)
		}

		// Fetch applications
		applications, err := q.WithHeading(func(hq *ent.HeadingQuery) { hq.WithVarsity() }).WithRun().Limit(limit).All(ctx)
		if err != nil {
			log.Printf("error getting applications: %v", err)
			return fiber.ErrInternalServerError
		}

		// Determine hasNextPage
		hasNextPage := len(applications) > first
		if hasNextPage {
			applications = applications[:first]
		}

		// Fetch precomputed flags from materialized view
		appIDs := make([]int, len(applications))
		for i, app := range applications {
			appIDs[i] = app.ID
		}
		flagsQuery := fmt.Sprintf("SELECT application_id, passing_now, passing_to_more_priority, another_varsities_count, original_submitted FROM application_flags WHERE application_id IN (%s)", utils.IntsToSQLIn(appIDs))
		rows, err := client.QueryContext(ctx, flagsQuery)
		if err != nil {
			log.Printf("error querying materialized view: %v", err)
			return fiber.ErrInternalServerError
		}
		defer rows.Close()

		flagsMap := make(map[int]struct {
			PassingNow, PassingToMorePriority bool
			AnotherVarsitiesCount             int
			OriginalSubmitted                 bool
		})
		for rows.Next() {
			var appID int
			var passingNow, passingToMorePriority, originalSubmitted bool
			var anotherVarsitiesCount int
			if err := rows.Scan(&appID, &passingNow, &passingToMorePriority, &anotherVarsitiesCount, &originalSubmitted); err != nil {
				log.Printf("error scanning flags: %v", err)
				return fiber.ErrInternalServerError
			}
			flagsMap[appID] = struct {
				PassingNow, PassingToMorePriority bool
				AnotherVarsitiesCount             int
				OriginalSubmitted                 bool
			}{passingNow, passingToMorePriority, anotherVarsitiesCount, originalSubmitted}
		}

		// Build edges
		edges := make([]ApplicationEdge, len(applications))
		for i, app := range applications {
			flags := flagsMap[app.ID]
			node := ApplicationResponse{
				ID:                    app.ID,
				StudentID:             app.StudentID,
				Priority:              app.Priority,
				CompetitionType:       app.CompetitionType.String(),
				RatingPlace:           app.RatingPlace,
				Score:                 app.Score,
				RunID:                 app.RunID,
				UpdatedAt:             app.UpdatedAt,
				HeadingID:             app.Edges.Heading.ID,
				OriginalSubmitted:     flags.OriginalSubmitted,
				OriginalQuit:          false, // TODO: Compute or precompute if needed
				PassingNow:            flags.PassingNow,
				PassingToMorePriority: flags.PassingToMorePriority,
				AnotherVarsitiesCount: flags.AnotherVarsitiesCount,
			}
			cursorStr := fmt.Sprintf("%d:%d", app.RatingPlace, app.ID)
			edges[i] = ApplicationEdge{
				Node:   node,
				Cursor: base64.StdEncoding.EncodeToString([]byte(cursorStr)),
			}
		}

		// PageInfo
		var endCursor string
		if len(edges) > 0 {
			endCursor = edges[len(edges)-1].Cursor
		}

		// Total count (approximate or exact)
		totalCount, err := q.Clone().Count(ctx)
		if err != nil {
			log.Printf("error counting applications: %v", err)
			totalCount = 0
		}

		connection := ApplicationsConnection{
			Edges:      edges,
			PageInfo:   PageInfo{HasNextPage: hasNextPage, EndCursor: endCursor},
			TotalCount: totalCount,
		}

		return c.JSON(connection)
	}
}

// Cursor-based pagination structures
type PageInfo struct {
	HasNextPage bool   `json:"hasNextPage"`
	EndCursor   string `json:"endCursor"`
}

type ApplicationEdge struct {
	Node   ApplicationResponse `json:"node"`
	Cursor string              `json:"cursor"`
}

type ApplicationsConnection struct {
	Edges      []ApplicationEdge `json:"edges"`
	PageInfo   PageInfo          `json:"pageInfo"`
	TotalCount int               `json:"totalCount"`
}
