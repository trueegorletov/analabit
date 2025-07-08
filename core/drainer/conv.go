package drainer

import (
	"github.com/trueegorletov/analabit/core"
	"log/slog"
)

func ConvResults(results []core.CalculationResult) []DrainedResult {
	drained := make([]DrainedResult, 0, len(results))

	for _, result := range results {
		passingScore, err := result.PassingScore()
		if err != nil {
			slog.Debug("Passing score unavailable", "error", err, "heading", result.Heading.Code())
			continue
		}

		lastAdmittedRatingPlace, err := result.LastAdmittedRatingPlace()
		if err != nil {
			slog.Debug("Last admitted rating place unavailable", "error", err, "heading", result.Heading.Code())
			continue
		}

		drained = append(drained, DrainedResult{
			Heading:                    result.Heading,
			AvgPassingScore:            passingScore,
			AvgLastAdmittedRatingPlace: lastAdmittedRatingPlace,
			DrainedPercent:             0,
		})
	}

	return drained
}
