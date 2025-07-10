package drainer

import (
	"fmt"
	"log/slog"
	"sort"
	"sync"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
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
		psValues         []int
		larpValues       []int
		psSum            int
		larpSum          int
	}

	codeToResult := make(map[string]headingResults)

	resultsChan := make(chan []core.CalculationResult, iterations)

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
				results.psValues = make([]int, 0, iterations)
				results.larpValues = make([]int, 0, iterations)
			}

			passingScore, err := result.PassingScore()

			if err != nil {
				slog.Debug("Passing score unavailable", "error", err, "heading", code)
				continue
			}

			lastAdmittedRatingPlace, err := result.LastAdmittedRatingPlace()

			if err != nil {
				slog.Debug("Last admitted rating place unavailable", "error", err, "heading", code)
				continue
			}

			results.psValues = append(results.psValues, passingScore)
			results.larpValues = append(results.larpValues, lastAdmittedRatingPlace)
			results.psSum += passingScore
			results.larpSum += lastAdmittedRatingPlace

			codeToResult[code] = results
		}
	}

	drained := make([]DrainedResult, 0, len(codeToResult))

	for _, results := range codeToResult {
		if len(results.psValues) == 0 || len(results.larpValues) == 0 {
			fmt.Println("No results for heading", "heading", results.prototypeHeading.PrettyName())
			continue
		}

		sort.Ints(results.psValues)
		sort.Ints(results.larpValues)

		drained = append(drained, DrainedResult{
			Heading:                    results.prototypeHeading,
			MinPassingScore:            results.psValues[0],
			MaxPassingScore:            results.psValues[len(results.psValues)-1],
			AvgPassingScore:            results.psSum / len(results.psValues),
			MedPassingScore:            median(results.psValues),
			MinLastAdmittedRatingPlace: results.larpValues[0],
			MaxLastAdmittedRatingPlace: results.larpValues[len(results.larpValues)-1],
			AvgLastAdmittedRatingPlace: results.larpSum / len(results.larpValues),
			MedLastAdmittedRatingPlace: median(results.larpValues),
			DrainedPercent:             d.drainPercent,
		})
	}

	return drained
}

func median(data []int) int {
	n := len(data)
	if n == 0 {
		return 0
	}
	// data is expected to be sorted
	if n%2 != 0 {
		return data[n/2]
	}
	return (data[n/2-1] + data[n/2]) / 2
}
