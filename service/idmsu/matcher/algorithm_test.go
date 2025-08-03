package matcher

import (
	"testing"

	"github.com/trueegorletov/analabit/core/idresolver"
)

func TestCalculateWeakMatchScore(t *testing.T) {
	engine := NewMatchingEngine()

	// Test data
	msuApp := idresolver.MSUAppDetails{
		PrettyName:        "Test Program",
		ScoreSum:          250,
		RatingPlace:       5,
		Priority:          1,
		DVIScore:          0,
		EGEScores:         []int{85, 90, 95},
		AchievementsScore: 10,
	}

	candidate := GosuslugiApplicant{
		Rating:           5,
		Priority:         1,
		SumMark:          251.0, // Close to MSU sum
		Result1:          85.0,
		Result2:          90.0,
		Result3:          95.0,
		AchievementsMark: 10.0,
		IDApplication:    12345,
	}

	score := engine.calculateWeakMatchScore(msuApp, candidate)

	// Expected score calculation:
	// Base: 0.4
	// EGE identical: +0.25
	// Score delta = 1: +0.2
	// Priority match: +0.1
	// Achievements match: +0.05
	// Total: 1.0

	expectedScore := 1.0
	if score != expectedScore {
		t.Errorf("Expected score %f, got %f", expectedScore, score)
	}
}

func TestCalculateEGEScoreSimilarity(t *testing.T) {
	engine := NewMatchingEngine()

	// Test perfect match
	candidate := GosuslugiApplicant{
		Result1: 85.0,
		Result2: 90.0,
		Result3: 95.0,
	}

	msuScores := []int{85, 90, 95}
	similarity := engine.calculateEGEScoreSimilarity(msuScores, candidate)

	if similarity != 0.25 {
		t.Errorf("Expected perfect match similarity 0.25, got %f", similarity)
	}

	// Test partial match (2 out of 3)
	candidate2 := GosuslugiApplicant{
		Result1: 85.0,
		Result2: 90.0,
		Result3: 75.0, // Different from MSU
	}

	similarity2 := engine.calculateEGEScoreSimilarity(msuScores, candidate2)

	if similarity2 != 0.1 {
		t.Errorf("Expected partial match similarity 0.1, got %f", similarity2)
	}
}

func TestCountEGEScoreMatches(t *testing.T) {
	engine := NewMatchingEngine()

	msuScores := []int{85, 90, 95}
	candidateScores := []int{85, 90, 75}

	matches := engine.countEGEScoreMatches(msuScores, candidateScores)

	if matches != 2 {
		t.Errorf("Expected 2 matches, got %d", matches)
	}
}
