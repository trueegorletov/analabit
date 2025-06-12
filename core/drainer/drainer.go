package drainer

import (
	"analabit/core"
	"analabit/core/source"
	"fmt"
	"log/slog"
	"sync"
)

type Drainer struct {
	prototype    *source.Varsity
	drainPercent int
}

func New(prototype *source.Varsity, drainPercent int) *Drainer {
	if drainPercent < 0 || drainPercent > 100 {
		panic("drain percent must be between 0 and 100")
	}

	return &Drainer{
		prototype:    prototype,
		drainPercent: drainPercent,
	}
}

const maxComputeGoroutines = 128

func (d *Drainer) Run(iterations int) []DrainedResult {
	type headingResults struct {
		prototypeHeading *core.Heading
		psSum, psN       int
		larpSum, larpN   int
	}

	codeToResult := make(map[string]headingResults)

	resultsChan := make(chan []core.CalculationResult, maxComputeGoroutines*2)

	sema := make(chan struct{}, maxComputeGoroutines)

	wg := sync.WaitGroup{}
	for it := 0; it < iterations; it++ {
		wg.Add(1)
		go func() {
			sema <- struct{}{}
			defer func() {
				wg.Done()
				<-sema
			}()

			vc := d.prototype.Clone().VarsityCalculator
			vc.SimulateOriginalsDrain(d.drainPercent)
			resultsChan <- vc.CalculateAdmissions()
		}()
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for step := range resultsChan {
		for _, result := range step {
			code := result.Heading.Code()

			results, ok := codeToResult[code]

			if !ok {
				results.prototypeHeading = d.prototype.GetHeading(code)
			}

			passingScore, err := result.PassingScore()

			if err != nil {
				slog.Error("Failed to get passing score", "error", err, "heading", code)
				continue
			}

			lastAdmittedRatingPlace, err := result.LastAdmittedRatingPlace()

			if err != nil {
				slog.Error("Failed to get last admitted rating place", "error", err, "heading", code)
				continue
			}

			results.psSum += passingScore
			results.psN++
			results.larpSum += lastAdmittedRatingPlace
			results.larpN++

			codeToResult[code] = results
		}
	}

	drained := make([]DrainedResult, 0, len(codeToResult))

	for _, results := range codeToResult {
		if results.psN == 0 || results.larpN == 0 {
			fmt.Println("No results for heading", "heading", results.prototypeHeading.PrettyName())
			continue
		}

		avgPassingScore := results.psSum / results.psN
		avgLastAdmittedRatingPlace := results.larpSum / results.larpN
		drained = append(drained, DrainedResult{
			Heading:                    results.prototypeHeading,
			AvgPassingScore:            avgPassingScore,
			AvgLastAdmittedRatingPlace: avgLastAdmittedRatingPlace,
			DrainedPercent:             d.drainPercent,
		})
	}

	return drained
}
