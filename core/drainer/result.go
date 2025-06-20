package drainer

import "analabit/core"

type DrainedResult struct {
	Heading *core.Heading

	MinPassingScore int
	MaxPassingScore int
	AvgPassingScore int
	MedPassingScore int

	MinLastAdmittedRatingPlace int
	MaxLastAdmittedRatingPlace int
	AvgLastAdmittedRatingPlace int
	MedLastAdmittedRatingPlace int

	DrainedPercent int
}
