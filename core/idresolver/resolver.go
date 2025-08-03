package idresolver

import "context"

// MSUAppDetails contains the details of a single MSU application needed for ID resolution
type MSUAppDetails struct {
	PrettyName        string // Program name (e.g., "Математика", "Физика")
	ScoreSum          int
	RatingPlace       int
	Priority          int
	DVIScore          int
	EGEScores         []int
	AchievementsScore int
}

// ResolveRequestItem represents a single internal ID with all its applications
type ResolveRequestItem struct {
	InternalID string
	Apps       []MSUAppDetails
}

// ResolveResponseItem represents the resolution result for a single internal ID
type ResolveResponseItem struct {
	InternalID  string
	CanonicalID string
	Confidence  float64 // 0-1 score, where 1.0 is exact match
}

// StudentIDResolver is the interface for resolving MSU internal student IDs to canonical IDs
type StudentIDResolver interface {
	ResolveBatch(ctx context.Context, req []ResolveRequestItem) ([]ResolveResponseItem, error)
}
