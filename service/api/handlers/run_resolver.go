package handlers

import (
	"github.com/trueegorletov/analabit/core/ent"
	"github.com/trueegorletov/analabit/core/ent/run"
	"context"
	"fmt"
	"strconv"
	"strings"
)

// RunResolution contains the resolved run information
type RunResolution struct {
	RunID      int
	IsLatest   bool
	IsRelative bool
	Offset     int // For relative queries (negative values)
}

// ResolveRunFromIteration resolves a run ID from the iteration parameter
// Supports:
// - "latest" or "0" or empty → latest run (offset 0)
// - negative "-n" → nth previous run (offset -n)
// Positive iteration values are no longer supported.
func ResolveRunFromIteration(ctx context.Context, client *ent.Client, iterationParam string) (*RunResolution, error) {
	iterationParam = strings.TrimSpace(strings.ToLower(iterationParam))

	// Handle default cases
	if iterationParam == "" || iterationParam == "latest" || iterationParam == "0" {
		runID, err := getLatestRunID(ctx, client)
		if err != nil {
			return nil, fmt.Errorf("failed to get latest run ID: %w", err)
		}
		return &RunResolution{
			RunID:      runID,
			IsLatest:   true,
			IsRelative: true,
			Offset:     0,
		}, nil
	}

	// Try to parse as integer
	offset, err := strconv.Atoi(iterationParam)
	if err != nil {
		return nil, fmt.Errorf("invalid iteration parameter: %s", iterationParam)
	}

	if offset < 0 {
		// Relative offset: get the nth previous run
		runID, err := getRunAtOffset(ctx, client, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to get run at offset %d: %w", offset, err)
		}
		return &RunResolution{
			RunID:      runID,
			IsLatest:   false,
			IsRelative: true,
			Offset:     offset,
		}, nil
	}

	return nil, fmt.Errorf("positive iteration values are no longer supported")
}

// getLatestRunID returns the ID of the most recent finished run
func getLatestRunID(ctx context.Context, client *ent.Client) (int, error) {
	latestRun, err := client.Run.Query().
		Where(run.FinishedEQ(true)).
		Order(ent.Desc(run.FieldID)).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return 0, fmt.Errorf("no finished runs found")
		}
		return 0, err
	}
	return latestRun.ID, nil
}

// getRunAtOffset returns the run ID at the specified offset from latest
// offset should be negative (e.g., -1 for previous run, -2 for the one before that)
func getRunAtOffset(ctx context.Context, client *ent.Client, offset int) (int, error) {
	if offset >= 0 {
		return 0, fmt.Errorf("offset must be negative for relative queries")
	}

	// Get finished runs ordered by ID descending, skip by absolute offset
	skip := -offset // Convert negative offset to positive skip value
	runs, err := client.Run.Query().
		Where(run.FinishedEQ(true)).
		Order(ent.Desc(run.FieldID)).
		Offset(skip).
		Limit(1).
		All(ctx)

	if err != nil {
		return 0, err
	}

	if len(runs) == 0 {
		return 0, fmt.Errorf("no run found at offset %d", offset)
	}

	return runs[0].ID, nil
}
