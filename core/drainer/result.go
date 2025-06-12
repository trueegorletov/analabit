package drainer

import "analabit/core"

type DrainedResult struct {
	Heading *core.Heading

	AvgPassingScore            int
	AvgLastAdmittedRatingPlace int

	DrainedPercent int
}
