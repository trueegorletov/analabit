package main

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

func setupTestDB(t *testing.T) *ent.Client {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	return client
}

func createTestData(t *testing.T, client *ent.Client) {
	ctx := context.Background()

	// Create a test varsity
	varsity, err := client.Varsity.Create().
		SetCode("TEST").
		SetName("Test University").
		Save(ctx)
	require.NoError(t, err)

	// Create a test heading
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

	// Create dummy runs first to satisfy foreign key constraints
	run1, err := client.Run.Create().
		SetTriggeredAt(time.Now()).
		SetPayloadMeta(map[string]any{"dummy": "run1"}).
		Save(ctx)
	require.NoError(t, err)

	run2, err := client.Run.Create().
		SetTriggeredAt(time.Now()).
		SetPayloadMeta(map[string]any{"dummy": "run2"}).
		Save(ctx)
	require.NoError(t, err)

	// Create applications with different iterations
	baseTime := time.Now().Add(-24 * time.Hour)

	// Iteration 1 data
	_, err = client.Application.Create().
		SetStudentID("student1").
		SetPriority(1).
		SetCompetitionType(1).
		SetRatingPlace(10).
		SetScore(100).
		SetIteration(1).
		SetRunID(run1.ID). // Will be updated by migration
		SetUpdatedAt(baseTime).
		SetHeading(heading).
		Save(ctx)
	require.NoError(t, err)

	_, err = client.Calculation.Create().
		SetStudentID("student1").
		SetAdmittedPlace(1).
		SetIteration(1).
		SetRunID(run1.ID). // Will be updated by migration
		SetUpdatedAt(baseTime.Add(1 * time.Hour)).
		SetHeading(heading).
		Save(ctx)
	require.NoError(t, err)

	_, err = client.DrainedResult.Create().
		SetDrainedPercent(0).
		SetAvgPassingScore(100).
		SetMinPassingScore(100).
		SetMaxPassingScore(100).
		SetMedPassingScore(100).
		SetAvgLastAdmittedRatingPlace(1).
		SetMinLastAdmittedRatingPlace(1).
		SetMaxLastAdmittedRatingPlace(1).
		SetMedLastAdmittedRatingPlace(1).
		SetIteration(1).
		SetRunID(run1.ID). // Will be updated by migration
		SetHeading(heading).
		Save(ctx)
	require.NoError(t, err)

	// Iteration 2 data
	_, err = client.Application.Create().
		SetStudentID("student2").
		SetPriority(1).
		SetCompetitionType(1).
		SetRatingPlace(5).
		SetScore(120).
		SetIteration(2).
		SetRunID(run2.ID). // Will be updated by migration
		SetUpdatedAt(baseTime.Add(2 * time.Hour)).
		SetHeading(heading).
		Save(ctx)
	require.NoError(t, err)

	_, err = client.Calculation.Create().
		SetStudentID("student2").
		SetAdmittedPlace(1).
		SetIteration(2).
		SetRunID(run2.ID). // Will be updated by migration
		SetUpdatedAt(baseTime.Add(3 * time.Hour)).
		SetHeading(heading).
		Save(ctx)
	require.NoError(t, err)
}

func TestGetDistinctIterations(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	createTestData(t, client)
	ctx := context.Background()

	iterations, err := getDistinctIterations(client, ctx)
	require.NoError(t, err)

	expected := []int{1, 2}
	assert.Equal(t, expected, iterations)
}

func TestGetEarliestUpdatedAt(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	createTestData(t, client)
	ctx := context.Background()

	// Test with transaction
	tx, err := client.Tx(ctx)
	require.NoError(t, err)
	defer tx.Rollback()

	// For iteration 1, the earliest should be from application (baseTime)
	earliest, err := getEarliestUpdatedAt(tx, ctx, 1)
	require.NoError(t, err)

	// Should be close to baseTime (allowing for small time differences)
	baseTime := time.Now().Add(-24 * time.Hour)
	timeDiff := earliest.Sub(baseTime)
	assert.True(t, timeDiff < time.Minute, "Expected earliest time to be close to baseTime")
}

func TestMigrateIteration(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	createTestData(t, client)
	ctx := context.Background()

	// Migrate iteration 1
	err := migrateIteration(client, ctx, 1)
	require.NoError(t, err)

	// Verify that a new run was created (in addition to the 2 dummy runs)
	runs, err := client.Run.Query().All(ctx)
	require.NoError(t, err)
	assert.Len(t, runs, 3) // 2 dummy runs + 1 new migration run

	// Find the migration run (should have migration metadata)
	var run *ent.Run
	for _, r := range runs {
		if meta, ok := r.PayloadMeta["migrated_from_iteration"]; ok && meta == float64(1) {
			run = r
			break
		}
	}
	require.NotNil(t, run, "Migration run should exist")
	assert.NotNil(t, run.PayloadMeta)
	if meta, ok := run.PayloadMeta["migrated_from_iteration"]; ok {
		assert.Equal(t, float64(1), meta) // JSON numbers are float64
	}

	// Verify that all records for iteration 1 now reference the new run
	apps, err := client.Application.Query().All(ctx)
	require.NoError(t, err)

	for _, app := range apps {
		if app.Iteration == 1 {
			assert.Equal(t, run.ID, app.RunID)
		}
	}

	calcs, err := client.Calculation.Query().All(ctx)
	require.NoError(t, err)

	for _, calc := range calcs {
		if calc.Iteration == 1 {
			assert.Equal(t, run.ID, calc.RunID)
		}
	}

	drained, err := client.DrainedResult.Query().All(ctx)
	require.NoError(t, err)

	for _, dr := range drained {
		if dr.Iteration == 1 {
			assert.Equal(t, run.ID, dr.RunID)
		}
	}
}

func TestFullMigration(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	createTestData(t, client)
	ctx := context.Background()

	// Get all iterations
	iterations, err := getDistinctIterations(client, ctx)
	require.NoError(t, err)

	// Migrate all iterations
	for _, iteration := range iterations {
		err := migrateIteration(client, ctx, iteration)
		require.NoError(t, err)
	}

	// Verify that we have the correct number of runs
	runs, err := client.Run.Query().All(ctx)
	require.NoError(t, err)
	assert.Len(t, runs, 4) // Should have 2 dummy runs + 2 migration runs for iterations 1 and 2

	// Verify that all records have been updated with proper run IDs
	apps, err := client.Application.Query().All(ctx)
	require.NoError(t, err)

	for _, app := range apps {
		assert.NotEqual(t, app.Iteration, app.RunID, "RunID should not equal iteration after migration")
		assert.True(t, app.RunID > 0, "RunID should be set")
	}

	calcs, err := client.Calculation.Query().All(ctx)
	require.NoError(t, err)

	for _, calc := range calcs {
		assert.NotEqual(t, calc.Iteration, calc.RunID, "RunID should not equal iteration after migration")
		assert.True(t, calc.RunID > 0, "RunID should be set")
	}

	drained, err := client.DrainedResult.Query().All(ctx)
	require.NoError(t, err)

	for _, dr := range drained {
		assert.NotEqual(t, dr.Iteration, dr.RunID, "RunID should not equal iteration after migration")
		assert.True(t, dr.RunID > 0, "RunID should be set")
	}

	// Verify that records with the same iteration have the same run ID
	appsByIteration := make(map[int][]int) // iteration -> run IDs
	for _, app := range apps {
		appsByIteration[app.Iteration] = append(appsByIteration[app.Iteration], app.RunID)
	}

	for iteration, runIDs := range appsByIteration {
		// All records for the same iteration should have the same run ID
		if len(runIDs) > 1 {
			firstRunID := runIDs[0]
			for _, runID := range runIDs[1:] {
				assert.Equal(t, firstRunID, runID, "All records for iteration %d should have the same run ID", iteration)
			}
		}
	}
}

func TestEmptyDatabase(t *testing.T) {
	client := setupTestDB(t)
	defer client.Close()

	ctx := context.Background()

	// Test with empty database
	iterations, err := getDistinctIterations(client, ctx)
	require.NoError(t, err)
	assert.Empty(t, iterations)
}
