package handlers

import (
	"analabit/core/ent"
	"analabit/core/ent/enttest"
	"context"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestClient(t *testing.T) *ent.Client {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	return client
}

func createTestRuns(t *testing.T, client *ent.Client, count int) []*ent.Run {
	runs := make([]*ent.Run, count)
	for i := 0; i < count; i++ {
		run, err := client.Run.Create().
			SetTriggeredAt(time.Now().Add(time.Duration(i) * time.Hour)).
			SetPayloadMeta(map[string]any{"test_run": i}).
			SetFinished(true). // Mark as finished for testing
			Save(context.Background())
		require.NoError(t, err)
		runs[i] = run
	}
	return runs
}

func TestResolveRunFromIteration_Basic(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// Create a test run
	run, err := client.Run.Create().
		SetTriggeredAt(time.Now()).
		SetPayloadMeta(map[string]any{"test": "data"}).
		SetFinished(true). // Mark as finished for testing
		Save(ctx)
	require.NoError(t, err)

	// Test latest resolution
	resolution, err := ResolveRunFromIteration(ctx, client, "latest")
	require.NoError(t, err)
	assert.Equal(t, run.ID, resolution.RunID)
	assert.True(t, resolution.IsLatest)
	assert.True(t, resolution.IsRelative)
	assert.Equal(t, 0, resolution.Offset)
}

func TestResolveRunFromIteration_Latest(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	runs := createTestRuns(t, client, 3)
	ctx := context.Background()

	tests := []struct {
		name     string
		param    string
		expected int
	}{
		{"empty string", "", runs[2].ID},
		{"latest", "latest", runs[2].ID},
		{"zero", "0", runs[2].ID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolution, err := ResolveRunFromIteration(ctx, client, tt.param)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, resolution.RunID)
			assert.True(t, resolution.IsLatest)
			assert.True(t, resolution.IsRelative)
			assert.Equal(t, 0, resolution.Offset)
		})
	}
}

func TestResolveRunFromIteration_RelativeOffset(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	runs := createTestRuns(t, client, 5)
	ctx := context.Background()

	tests := []struct {
		name     string
		param    string
		expected int
		offset   int
	}{
		{"previous run", "-1", runs[3].ID, -1},
		{"two runs back", "-2", runs[2].ID, -2},
		{"three runs back", "-3", runs[1].ID, -3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolution, err := ResolveRunFromIteration(ctx, client, tt.param)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, resolution.RunID)
			assert.False(t, resolution.IsLatest)
			assert.True(t, resolution.IsRelative)
			assert.Equal(t, tt.offset, resolution.Offset)
		})
	}
}

func TestResolveRunFromIteration_ErrorCases(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	tests := []struct {
		name  string
		param string
	}{
		{"invalid string", "invalid"},
		{"positive iteration not supported", "999"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ResolveRunFromIteration(ctx, client, tt.param)
			assert.Error(t, err)
		})
	}
}

func TestGetLatestRunID(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	runs := createTestRuns(t, client, 3)

	latestID, err := getLatestRunID(ctx, client)
	require.NoError(t, err)
	assert.Equal(t, runs[2].ID, latestID)
}

func TestGetLatestRunID_OnlyFinishedRuns(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// Create finished runs
	finishedRuns := createTestRuns(t, client, 2)

	// Create an unfinished run (newer than finished ones)
	_, err := client.Run.Create().
		SetTriggeredAt(time.Now().Add(24 * time.Hour)).
		SetPayloadMeta(map[string]any{"test": "unfinished"}).
		SetFinished(false). // This should be ignored
		Save(ctx)
	require.NoError(t, err)

	// Should return the latest finished run, not the unfinished one
	latestID, err := getLatestRunID(ctx, client)
	require.NoError(t, err)
	assert.Equal(t, finishedRuns[1].ID, latestID)
}
