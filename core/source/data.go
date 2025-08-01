package source

import "github.com/trueegorletov/analabit/core"

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
	// MSU-specific fields for ID resolution
	DVIScore          int   // DVI (additional entrance exam) score, 0 if not applicable
	EGEScores         []int // Individual EGE scores, empty if not available
	HeadingName       string // MSU-specific field for pretty name
}

// HeadingData mirrors the data needed for core.VarsityCalculator.AddHeading.
type HeadingData struct {
	Code       string          // Unique code for the heading
	Capacities core.Capacities // Number of places available in this heading
	PrettyName string          // Name of the heading
}
