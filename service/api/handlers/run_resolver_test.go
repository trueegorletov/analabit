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

func TestResolveRunFromIteration_BackwardCompatibility(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// Create a test varsity and heading
	varsity, err := client.Varsity.Create().
		SetCode("TEST").
		SetName("Test University").
		Save(ctx)
	require.NoError(t, err)

	heading, err := client.Heading.Create().
		SetCode("TEST-H1").
		SetName("Test Heading").
		SetRegularCapacity(100).
		SetTargetQuotaCapacity(10).
		SetDedicatedQuotaCapacity(10).
		SetSpecialQuotaCapacity(5).
		SetVarsity(varsity).
		Save(ctx)
	require.NoError(t, err)

	runs := createTestRuns(t, client, 3)

	// Create test applications with different iteration values
	_, err = client.Application.Create().
		SetStudentID("student1").
		SetPriority(1).
		SetCompetitionType(1).
		SetRatingPlace(10).
		SetScore(100).
		SetIteration(100).
		SetRunID(runs[0].ID).
		SetHeading(heading).
		Save(ctx)
	require.NoError(t, err)

	// Test finding run by iteration
	resolution, err := ResolveRunFromIteration(ctx, client, "100")
	require.NoError(t, err)
	assert.Equal(t, runs[0].ID, resolution.RunID)
	assert.False(t, resolution.IsLatest)
	assert.False(t, resolution.IsRelative)
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
		{"nonexistent iteration", "999"},
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
