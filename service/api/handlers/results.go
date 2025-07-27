package handlers

import (
	"context"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/trueegorletov/analabit/core/ent"
	"github.com/trueegorletov/analabit/core/ent/application"
	"github.com/trueegorletov/analabit/core/ent/calculation"
	"github.com/trueegorletov/analabit/core/ent/drainedresult"
	"github.com/trueegorletov/analabit/core/ent/heading"
	"github.com/trueegorletov/analabit/core/ent/varsity"

	"github.com/gofiber/fiber/v3"
)

// CalculationResultDTO represents common calculations result per heading
type CalculationResultDTO struct {
	HeadingID               int    `json:"heading_id"`
	HeadingCode             string `json:"heading_code"`
	PassingScore            int    `json:"passing_score"`
	LastAdmittedRatingPlace int    `json:"last_admitted_rating_place"`
	RunID                   int    `json:"run_id"`
	RegularsAdmitted        bool   `json:"regulars_admitted"`
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
	RunID                      int    `json:"run_id"`
	RegularsAdmitted           bool   `json:"regulars_admitted"`
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
		runParam := c.Query("run", "latest")

		// Resolve the run ID from the parameter
		runResolution, err := ResolveRunFromIteration(ctx, client, runParam)
		if err != nil {
			log.Printf("error resolving run from parameter '%s': %v", runParam, err)
			return fiber.NewError(fiber.StatusBadRequest, "invalid run parameter")
		}

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
		// Track explicit request for 100% drained
		explicit100 := false
		for _, step := range requestedSteps {
			if step == 100 {
				explicit100 = true
				break
			}
		}

		resp := ResultsResponse{}

		// ---------------------------------------------------------
		// Build steps map (available drainedPercent values) always

		stepsQuery := client.DrainedResult.Query().Where(drainedresult.RunIDEQ(runResolution.RunID)).WithHeading().WithRun()

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
			drQuery := client.DrainedResult.Query().Where(
				drainedresult.RunIDEQ(runResolution.RunID),
				drainedresult.DrainedPercentEQ(0),
			).WithHeading().WithRun()

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
				runID := 0
				if dr.Edges.Run != nil {
					runID = dr.Edges.Run.ID
				}
				dto := CalculationResultDTO{
					HeadingID:               hid,
					HeadingCode:             dr.Edges.Heading.Code,
					PassingScore:            dr.AvgPassingScore,
					LastAdmittedRatingPlace: dr.AvgLastAdmittedRatingPlace,
					RunID:                   runID,
					RegularsAdmitted:        dr.RegularsAdmitted,
				}
				primaryMap[hid] = dto
				coveredHeading[hid] = struct{}{}
			}

			// Second attempt: get from Calculation table (legacy)
			var uncoveredIDs []int
			if len(headingIDs) > 0 { // if request was filtered, only look for those missing
				for _, hid := range headingIDs {
					if _, ok := coveredHeading[hid]; !ok {
						uncoveredIDs = append(uncoveredIDs, hid)
					}
				}
			}

			if len(headingIDs) == 0 || len(uncoveredIDs) > 0 { // if no filter, try all
				calcQuery := client.Calculation.Query().Where(
					calculation.RunIDEQ(runResolution.RunID),
				).WithHeading().WithRun()
				if len(uncoveredIDs) > 0 {
					calcQuery = calcQuery.Where(calculation.HasHeadingWith(heading.IDIn(uncoveredIDs...)))
				} else if varsityCode != "" {
					calcQuery = calcQuery.Where(calculation.HasHeadingWith(heading.HasVarsityWith(varsity.CodeEQ(varsityCode))))
				}

				calcs, err := calcQuery.All(ctx)
				if err != nil {
					log.Printf("error fetching calculation results: %v", err)
					return fiber.ErrInternalServerError
				}
				for _, c := range calcs {
					if c.Edges.Heading == nil {
						continue
					}
					hid := c.Edges.Heading.ID
					if _, ok := primaryMap[hid]; !ok {
						runID := 0
						if c.Edges.Run != nil {
							runID = c.Edges.Run.ID
						}

						// Get score from application
						app, err := client.Application.Query().Where(
							application.StudentIDEQ(c.StudentID),
							application.HasHeadingWith(heading.ID(hid)),
							application.RunIDEQ(runResolution.RunID),
						).Only(ctx)

						passingScore := 0
						if err != nil {
							log.Printf("could not find application for student %s and heading %d in run %d", c.StudentID, hid, runResolution.RunID)
						} else {
							passingScore = app.Score
						}

						primaryMap[hid] = CalculationResultDTO{
							HeadingID:               hid,
							HeadingCode:             c.Edges.Heading.Code,
							PassingScore:            passingScore,
							LastAdmittedRatingPlace: c.AdmittedPlace,
							RunID:                   runID,
							RegularsAdmitted:        true, // Legacy, assume true or determine if needed
						}
					}
				}
			}
			resp.Primary = primaryMap
		}

		// DRAINED RESULTS -------------------------------------------
		if includeDrained {
			drQuery := client.DrainedResult.Query().Where(drainedresult.RunIDEQ(runResolution.RunID)).WithHeading(func(hq *ent.HeadingQuery) {
				hq.WithVarsity()
			}).WithRun()

			if len(headingIDs) > 0 {
				drQuery = drQuery.Where(drainedresult.HasHeadingWith(heading.IDIn(headingIDs...)))
			} else if varsityCode != "" {
				drQuery = drQuery.Where(drainedresult.HasHeadingWith(heading.HasVarsityWith(varsity.CodeEQ(varsityCode))))
			}

			// apply percent filter for requested stages (excluding explicit 100 fallback)
			if !drainedAll && !explicit100 && len(requestedSteps) > 0 {
				drQuery = drQuery.Where(drainedresult.DrainedPercentIn(requestedSteps...))
			}

			drResults, err := drQuery.All(ctx)
			if err != nil {
				log.Printf("error fetching drained results: %v", err)
				return fiber.ErrInternalServerError
			}

			// group drained results per heading
			drainedMap := make(map[int][]DrainedResultDTO)
			if explicit100 && !drainedAll {
				temp := make(map[int][]*ent.DrainedResult)
				for _, dr := range drResults {
					if dr.Edges.Heading == nil {
						continue
					}
					temp[dr.Edges.Heading.ID] = append(temp[dr.Edges.Heading.ID], dr)
				}
				for hid, drs := range temp {
					var chosen *ent.DrainedResult
					for _, dr := range drs {
						if dr.DrainedPercent == 100 {
							chosen = dr
							break
						}
					}
					if chosen == nil {
						max := 0
						for _, dr := range drs {
							if dr.DrainedPercent > max {
								max = dr.DrainedPercent
								chosen = dr
							}
						}
					}
					if chosen == nil {
						continue
					}
					runID := 0
					if chosen.Edges.Run != nil {
						runID = chosen.Edges.Run.ID
					}
					dto := DrainedResultDTO{
						HeadingID:                  hid,
						HeadingCode:                chosen.Edges.Heading.Code,
						DrainedPercent:             chosen.DrainedPercent,
						AvgPassingScore:            chosen.AvgPassingScore,
						MinPassingScore:            chosen.MinPassingScore,
						MaxPassingScore:            chosen.MaxPassingScore,
						MedPassingScore:            chosen.MedPassingScore,
						AvgLastAdmittedRatingPlace: chosen.AvgLastAdmittedRatingPlace,
						MinLastAdmittedRatingPlace: chosen.MinLastAdmittedRatingPlace,
						MaxLastAdmittedRatingPlace: chosen.MaxLastAdmittedRatingPlace,
						MedLastAdmittedRatingPlace: chosen.MedLastAdmittedRatingPlace,
						RunID:                      runID,
						RegularsAdmitted:           chosen.RegularsAdmitted,
					}
					drainedMap[hid] = []DrainedResultDTO{dto}
				}
			} else {
				for _, dr := range drResults {
					if dr.Edges.Heading == nil || dr.DrainedPercent <= 0 || dr.IsVirtual {
						continue
					}
					hid := dr.Edges.Heading.ID
					runID := 0
					if dr.Edges.Run != nil {
						runID = dr.Edges.Run.ID
					}
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
						RunID:                      runID,
						RegularsAdmitted:           dr.RegularsAdmitted,
					}
					drainedMap[hid] = append(drainedMap[hid], dto)
				}
			}
			resp.Drained = drainedMap
		}

		return c.JSON(resp)
	}
}
