package handlers

import (
	"analabit/core/ent"
	"analabit/core/ent/application"
	"analabit/core/ent/calculation"
	"analabit/core/ent/drainedresult"
	"analabit/core/ent/run"
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
// - positive "n>0" → backward compatibility: find run with iteration=n
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
	} else if offset > 0 {
		// Backward compatibility: treat as old iteration value
		runID, err := getRunByIteration(ctx, client, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to get run for iteration %d: %w", offset, err)
		}
		return &RunResolution{
			RunID:      runID,
			IsLatest:   false,
			IsRelative: false,
			Offset:     0,
		}, nil
	}

	// offset == 0 case (already handled above, but just in case)
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

// getLatestRunID returns the ID of the most recent run
func getLatestRunID(ctx context.Context, client *ent.Client) (int, error) {
	latestRun, err := client.Run.Query().
		Order(ent.Desc(run.FieldID)).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return 0, fmt.Errorf("no runs found")
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

	// Get runs ordered by ID descending, skip by absolute offset
	skip := -offset // Convert negative offset to positive skip value
	runs, err := client.Run.Query().
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

// getRunByIteration finds a run that contains records with the specified iteration value
// This is for backward compatibility with old iteration-based queries
func getRunByIteration(ctx context.Context, client *ent.Client, iteration int) (int, error) {
	// First try to find from application table
	app, err := client.Application.Query().
		Where(application.IterationEQ(iteration)).
		WithRun().
		First(ctx)
	if err == nil && app.Edges.Run != nil {
		return app.Edges.Run.ID, nil
	}

	// Try calculation table
	calc, err := client.Calculation.Query().
		Where(calculation.IterationEQ(iteration)).
		WithRun().
		First(ctx)
	if err == nil && calc.Edges.Run != nil {
		return calc.Edges.Run.ID, nil
	}

	// Try drained result table
	dr, err := client.DrainedResult.Query().
		Where(drainedresult.IterationEQ(iteration)).
		WithRun().
		First(ctx)
	if err == nil && dr.Edges.Run != nil {
		return dr.Edges.Run.ID, nil
	}

	return 0, fmt.Errorf("no run found with iteration %d", iteration)
}
