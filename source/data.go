package source

import "analabit/core"

// ApplicationData mirrors the data needed for core.VarsityCalculator.AddApplication
// and adds additional fields which could be useful for further processing.
type ApplicationData struct {
	HeadingCode       string
	StudentID         string
	ScoresSum         int
	RatingPlace       int
	Priority          int
	CompetitionType   core.Competition
	OriginalSubmitted bool
}

// HeadingData mirrors the data needed for core.VarsityCalculator.AddHeading.
type HeadingData struct {
	Code       string          // Unique code for the heading
	Capacities core.Capacities // Number of places available in this heading
	PrettyName string          // Name of the heading
}
