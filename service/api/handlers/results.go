package handlers

import (
	"analabit/core/ent"
	"analabit/core/ent/calculation"
	"analabit/core/ent/drainedresult"
	"analabit/core/ent/heading"
	"analabit/core/ent/varsity"
	"context"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
)

// CalculationResultDTO represents one admitted student record.
type CalculationResultDTO struct {
	HeadingCode   string `json:"heading_code"`
	StudentID     string `json:"student_id"`
	AdmittedPlace int    `json:"admitted_place"`
	Iteration     int    `json:"iteration"`
}

// DrainedResultDTO represents aggregated drained statistics per heading.
type DrainedResultDTO struct {
	HeadingCode                string `json:"heading_code"`
	DrainedPercent             int    `json:"drained_percent"`
	AvgPassingScore            int    `json:"avg_passing_score"`
	MinPassingScore            int    `json:"min_passing_score"`
	MaxPassingScore            int    `json:"max_passing_score"`
	MedPassingScore            int    `json:"med_passing_score"`
	AvgLastAdmittedRatingPlace int    `json:"avg_last_admitted_rating_place"`
	MinLastAdmittedRatingPlace int    `json:"min_last_admitted_rating_place"`
	MaxLastAdmittedRatingPlace int    `json:"max_last_admitted_rating_place"`
	MedLastAdmittedRatingPlace int    `json:"med_last_admitted_rating_place"`
	Iteration                  int    `json:"iteration"`
}

// ResultsResponse aggregates requested result kinds.
type ResultsResponse struct {
	Steps   map[int][]int          `json:"steps"`
	Primary []CalculationResultDTO `json:"primary,omitempty"`
	Drained []DrainedResultDTO     `json:"drained,omitempty"`
}

// GetResults returns calculation and/or drained results with optional batching.
func GetResults(client *ent.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := context.Background()

		// Parse headingIds (comma separated integers)
		var headingIDs []int
		if idsParam := c.Query("headingIds"); idsParam != "" {
			for _, part := range strings.Split(idsParam, ",") {
				part = strings.TrimSpace(part)
				if part == "" {
					continue
				}
				if id, err := strconv.Atoi(part); err == nil {
					headingIDs = append(headingIDs, id)
				}
			}
		}

		varsityCode := c.Query("varsityCode")
		iterationParam := strings.ToLower(c.Query("iteration", "latest"))

		// Presence of the `primary` query parameter (any value) toggles primary results.
		includePrimary := c.Query("primary") != ""

		// `drained` parameter, if present, determines which drained results to return.
		drainedParam := c.Query("drained") // "all" | "<csv steps>" | ""
		includeDrained := drainedParam != ""

		// Determine requested drained steps
		var requestedSteps []int
		drainedAll := false
		if includeDrained {
			if strings.ToLower(drainedParam) == "all" {
				drainedAll = true
			} else {
				for _, part := range strings.Split(drainedParam, ",") {
					part = strings.TrimSpace(part)
					if part == "" {
						continue
					}
					step, err := strconv.Atoi(part)
					if err != nil {
						return fiber.NewError(fiber.StatusBadRequest, "invalid drained step value")
					}
					requestedSteps = append(requestedSteps, step)
				}
			}
		}

		resp := ResultsResponse{}

		// helper to resolve iteration when "latest" is requested
		resolveIteration := func(latest bool, aggQuery func() (int, error)) (int, error) {
			if !latest {
				it, err := strconv.Atoi(iterationParam)
				if err != nil {
					return 0, fiber.NewError(fiber.StatusBadRequest, "invalid iteration value")
				}
				return it, nil
			}
			return aggQuery()
		}

		// Fetch latest iteration numbers lazily when needed
		getLatestCalculationIteration := func() (int, error) {
			var v []struct {
				Max int `json:"max"`
			}
			if err := client.Calculation.Query().Aggregate(ent.Max(calculation.FieldIteration)).Scan(ctx, &v); err != nil {
				return 0, err
			}
			if len(v) == 0 {
				return 0, nil
			}
			return v[0].Max, nil
		}

		getLatestDrainedIteration := func() (int, error) {
			var v []struct {
				Max int `json:"max"`
			}
			if err := client.DrainedResult.Query().Aggregate(ent.Max(drainedresult.FieldIteration)).Scan(ctx, &v); err != nil {
				return 0, err
			}
			if len(v) == 0 {
				return 0, nil
			}
			return v[0].Max, nil
		}

		// ---------------------------------------------------------
		// Build steps map (available drainedPercent values) always

		iterForSteps, err := resolveIteration(iterationParam == "latest", func() (int, error) {
			var v []struct {
				Max int `json:"max"`
			}
			if err := client.DrainedResult.Query().Aggregate(ent.Max(drainedresult.FieldIteration)).Scan(ctx, &v); err != nil {
				return 0, err
			}
			if len(v) == 0 {
				return 0, nil
			}
			return v[0].Max, nil
		})
		if err != nil {
			log.Printf("error resolving iteration for steps: %v", err)
			return fiber.ErrInternalServerError
		}

		stepsQuery := client.DrainedResult.Query().Where(drainedresult.IterationEQ(iterForSteps)).WithHeading()

		if len(headingIDs) > 0 {
			stepsQuery = stepsQuery.Where(drainedresult.HasHeadingWith(heading.IDIn(headingIDs...)))
		} else if varsityCode != "" {
			stepsQuery = stepsQuery.Where(drainedresult.HasHeadingWith(heading.HasVarsityWith(varsity.CodeEQ(varsityCode))))
		}

		stepsRows, err := stepsQuery.All(ctx)
		if err != nil {
			log.Printf("error fetching steps data: %v", err)
			return fiber.ErrInternalServerError
		}

		stepsMap := make(map[int][]int)
		tempSet := make(map[int]map[int]struct{}) // headingID -> set of steps
		for _, r := range stepsRows {
			if r.Edges.Heading == nil {
				continue
			}
			hid := r.Edges.Heading.ID
			if _, ok := tempSet[hid]; !ok {
				tempSet[hid] = make(map[int]struct{})
			}
			tempSet[hid][r.DrainedPercent] = struct{}{}
		}
		for hid, set := range tempSet {
			arr := make([]int, 0, len(set))
			for step := range set {
				arr = append(arr, step)
			}
			sort.Ints(arr)
			stepsMap[hid] = arr
		}

		resp.Steps = stepsMap

		// PRIMARY RESULTS -------------------------------------------
		if includePrimary {
			iter, err := resolveIteration(iterationParam == "latest", getLatestCalculationIteration)
			if err != nil {
				log.Printf("error resolving calculation iteration: %v", err)
				return fiber.ErrInternalServerError
			}

			calQuery := client.Calculation.Query().Where(calculation.IterationEQ(iter)).WithHeading(func(hq *ent.HeadingQuery) {
				hq.WithVarsity()
			})

			if len(headingIDs) > 0 {
				calQuery = calQuery.Where(calculation.HasHeadingWith(heading.IDIn(headingIDs...)))
			} else if varsityCode != "" {
				calQuery = calQuery.Where(calculation.HasHeadingWith(heading.HasVarsityWith(varsity.CodeEQ(varsityCode))))
			}

			calculations, err := calQuery.All(ctx)
			if err != nil {
				log.Printf("error fetching calculations: %v", err)
				return fiber.ErrInternalServerError
			}

			primaryResp := make([]CalculationResultDTO, len(calculations))
			for i, calc := range calculations {
				var headingCode string
				if calc.Edges.Heading != nil {
					headingCode = calc.Edges.Heading.Code
				}
				primaryResp[i] = CalculationResultDTO{
					HeadingCode:   headingCode,
					StudentID:     calc.StudentID,
					AdmittedPlace: calc.AdmittedPlace,
					Iteration:     calc.Iteration,
				}
			}
			resp.Primary = primaryResp
		}

		// DRAINED RESULTS -------------------------------------------
		if includeDrained {
			iter, err := resolveIteration(iterationParam == "latest", getLatestDrainedIteration)
			if err != nil {
				log.Printf("error resolving drained iteration: %v", err)
				return fiber.ErrInternalServerError
			}

			drQuery := client.DrainedResult.Query().Where(drainedresult.IterationEQ(iter)).WithHeading(func(hq *ent.HeadingQuery) {
				hq.WithVarsity()
			})

			if len(headingIDs) > 0 {
				drQuery = drQuery.Where(drainedresult.HasHeadingWith(heading.IDIn(headingIDs...)))
			} else if varsityCode != "" {
				drQuery = drQuery.Where(drainedresult.HasHeadingWith(heading.HasVarsityWith(varsity.CodeEQ(varsityCode))))
			}

			if includeDrained && !drainedAll && len(requestedSteps) > 0 {
				drQuery = drQuery.Where(drainedresult.DrainedPercentIn(requestedSteps...))
			}

			drained, err := drQuery.All(ctx)
			if err != nil {
				log.Printf("error fetching drained results: %v", err)
				return fiber.ErrInternalServerError
			}

			drainedResp := make([]DrainedResultDTO, len(drained))
			for i, dr := range drained {
				var headingCode string
				if dr.Edges.Heading != nil {
					headingCode = dr.Edges.Heading.Code
				}
				drainedResp[i] = DrainedResultDTO{
					HeadingCode:                headingCode,
					DrainedPercent:             dr.DrainedPercent,
					AvgPassingScore:            dr.AvgPassingScore,
					MinPassingScore:            dr.MinPassingScore,
					MaxPassingScore:            dr.MaxPassingScore,
					MedPassingScore:            dr.MedPassingScore,
					AvgLastAdmittedRatingPlace: dr.AvgLastAdmittedRatingPlace,
					MinLastAdmittedRatingPlace: dr.MinLastAdmittedRatingPlace,
					MaxLastAdmittedRatingPlace: dr.MaxLastAdmittedRatingPlace,
					MedLastAdmittedRatingPlace: dr.MedLastAdmittedRatingPlace,
					Iteration:                  dr.Iteration,
				}
			}
			resp.Drained = drainedResp
		}

		return c.JSON(resp)
	}
}
