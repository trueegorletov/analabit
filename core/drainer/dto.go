package drainer

import "analabit/core"

// NewDrainedResultDTOs converts a slice of DrainedResult to a slice of DrainedResultDTO.
func NewDrainedResultDTOs(results []DrainedResult) []core.DrainedResultDTO {
	dtos := make([]core.DrainedResultDTO, 0, len(results))
	for _, result := range results {
		dtos = append(dtos, core.DrainedResultDTO{
			HeadingCode:                result.Heading.FullCode(),
			DrainedPercent:             result.DrainedPercent,
			AvgPassingScore:            result.AvgPassingScore,
			MinPassingScore:            result.MinPassingScore,
			MaxPassingScore:            result.MaxPassingScore,
			MedPassingScore:            result.MedPassingScore,
			AvgLastAdmittedRatingPlace: result.AvgLastAdmittedRatingPlace,
			MinLastAdmittedRatingPlace: result.MinLastAdmittedRatingPlace,
			MaxLastAdmittedRatingPlace: result.MaxLastAdmittedRatingPlace,
			MedLastAdmittedRatingPlace: result.MedLastAdmittedRatingPlace,
		})
	}
	return dtos
}
