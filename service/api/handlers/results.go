package handlers

import (
	"analabit/core/ent"
	"analabit/core/ent/application"
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
	HeadingID               int    `json:"heading_id"`
	HeadingCode             string `json:"heading_code"`
	PassingScore            int    `json:"passing_score"`
	LastAdmittedRatingPlace int    `json:"last_admitted_rating_place"`
	Iteration               int    `json:"iteration"`
}

// DrainedResultDTO represents aggregated drained statistics per heading.
type DrainedResultDTO struct {
	HeadingID                  int    `json:"heading_id"`
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
	Steps   map[int][]int                `json:"steps"`
	Primary map[int]CalculationResultDTO `json:"primary,omitempty"`
	Drained map[int][]DrainedResultDTO   `json:"drained,omitempty"`
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
					requestStepInt, err := strconv.Atoi(part)
					if err != nil {
						return fiber.NewError(fiber.StatusBadRequest, "invalid drained step value")
					}
					if requestStepInt > 0 { // ignore non-positive
						requestedSteps = append(requestedSteps, requestStepInt)
					}
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
			if r.DrainedPercent > 0 { // exclude illegal non-positive
				tempSet[hid][r.DrainedPercent] = struct{}{}
			}
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
			// First attempt: get from DrainedResult table with drained_percent == 0
			iter, err := resolveIteration(iterationParam == "latest", getLatestDrainedIteration)
			if err != nil {
				log.Printf("error resolving drained iteration for primary: %v", err)
				return fiber.ErrInternalServerError
			}

			drQuery := client.DrainedResult.Query().Where(
				drainedresult.IterationEQ(iter),
				drainedresult.DrainedPercentEQ(0),
			).WithHeading()

			if len(headingIDs) > 0 {
				drQuery = drQuery.Where(drainedresult.HasHeadingWith(heading.IDIn(headingIDs...)))
			} else if varsityCode != "" {
				drQuery = drQuery.Where(drainedresult.HasHeadingWith(heading.HasVarsityWith(varsity.CodeEQ(varsityCode))))
			}

			zeroDrained, err := drQuery.All(ctx)
			if err != nil {
				log.Printf("error fetching 0%% drained results: %v", err)
				return fiber.ErrInternalServerError
			}

			primaryMap := make(map[int]CalculationResultDTO)
			coveredHeading := make(map[int]struct{})
			for _, dr := range zeroDrained {
				if dr.Edges.Heading == nil {
					continue
				}
				hid := dr.Edges.Heading.ID
				dto := CalculationResultDTO{
					HeadingID:               hid,
					HeadingCode:             dr.Edges.Heading.Code,
					PassingScore:            dr.AvgPassingScore,
					LastAdmittedRatingPlace: dr.AvgLastAdmittedRatingPlace,
					Iteration:               dr.Iteration,
				}
				primaryMap[hid] = dto
				coveredHeading[hid] = struct{}{}
			}

			// Fallback for headings missing (unlikely) – compute from Calculation + Application
			if len(headingIDs) == 0 && varsityCode == "" {
				// gather all headings in calculations fallback if requested none; else use requested set
			}

			// Find headings that need fallback (requested but not covered)
			var fallbackHeadings []int
			if len(headingIDs) > 0 {
				for _, hid := range headingIDs {
					if _, ok := coveredHeading[hid]; !ok {
						fallbackHeadings = append(fallbackHeadings, hid)
					}
				}
			} else {
				// need fallback for none because we won't know; we'll skip
			}

			if len(fallbackHeadings) > 0 {
				iterCalc, err := resolveIteration(iterationParam == "latest", getLatestCalculationIteration)
				if err != nil {
					log.Printf("error resolving calculation iteration for fallback: %v", err)
					return fiber.ErrInternalServerError
				}

				calFallbackQuery := client.Calculation.Query().Where(calculation.IterationEQ(iterCalc))
				calFallbackQuery = calFallbackQuery.Where(calculation.HasHeadingWith(heading.IDIn(fallbackHeadings...))).WithHeading()

				calRows, err := calFallbackQuery.All(ctx)
				if err != nil {
					log.Printf("error fetching fallback calculations: %v", err)
					return fiber.ErrInternalServerError
				}

				// group by heading -> max admitted place row
				calcMax := make(map[int]*ent.Calculation)
				for _, calc := range calRows {
					if calc.Edges.Heading == nil {
						continue
					}
					hid := calc.Edges.Heading.ID
					if prev, ok := calcMax[hid]; !ok || calc.AdmittedPlace > prev.AdmittedPlace {
						calcMax[hid] = calc
					}
				}

				for hid, calc := range calcMax {
					// query application to get score (maybe diff iterations)
					appRow, err := client.Application.Query().Where(
						application.StudentIDEQ(calc.StudentID),
						application.HasHeadingWith(heading.IDEQ(hid)),
					).Order(ent.Desc(application.FieldIteration)).First(ctx)
					passingScore := 0
					if err == nil {
						passingScore = appRow.Score
					} else {
						log.Printf("warning: application row not found for fallback StudentID %s heading %d: %v", calc.StudentID, hid, err)
					}

					dto := CalculationResultDTO{
						HeadingID:               hid,
						HeadingCode:             calc.Edges.Heading.Code,
						PassingScore:            passingScore,
						LastAdmittedRatingPlace: calc.AdmittedPlace,
						Iteration:               calc.Iteration,
					}
					primaryMap[hid] = dto
				}
			}

			resp.Primary = primaryMap
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

			drainedMap := make(map[int][]DrainedResultDTO)
			for _, dr := range drained {
				if dr.Edges.Heading == nil {
					continue
				}
				if dr.DrainedPercent <= 0 { // skip non-positive percent in drained response
					continue
				}
				hid := dr.Edges.Heading.ID
				dto := DrainedResultDTO{
					HeadingID:                  hid,
					HeadingCode:                dr.Edges.Heading.Code,
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
				drainedMap[hid] = append(drainedMap[hid], dto)
			}
			resp.Drained = drainedMap
		}

		return c.JSON(resp)
	}
}
