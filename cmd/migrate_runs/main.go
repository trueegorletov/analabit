package main

import (
	"analabit/core/ent"
	"analabit/core/ent/application"
	"analabit/core/ent/calculation"
	"analabit/core/ent/drainedresult"
	"analabit/core/utils"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: go run main.go <postgres_connection_string>")
	}

	connStr := os.Args[1]

	client, err := ent.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Ensure schema is up-to-date
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("Failed to run schema migrations: %v", err)
	}

	log.Println("Starting run table migration...")

	// Get all distinct iterations from all tables
	iterations, err := getDistinctIterations(client, ctx)
	if err != nil {
		log.Fatalf("Failed to get distinct iterations: %v", err)
	}

	log.Printf("Found %d distinct iterations to migrate", len(iterations))

	// Migrate each iteration
	for _, iteration := range iterations {
		if err := migrateIteration(client, ctx, iteration); err != nil {
			log.Fatalf("Failed to migrate iteration %d: %v", iteration, err)
		}
		log.Printf("Successfully migrated iteration %d", iteration)
	}

	log.Println("Migration completed successfully!")
}

func getDistinctIterations(client *ent.Client, ctx context.Context) ([]int, error) {
	// Get distinct iterations from all three tables using Ent queries
	iterationSet := make(map[int]struct{})

	// Get iterations from applications
	apps, err := client.Application.Query().Select(application.FieldIteration).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query application iterations: %w", err)
	}
	for _, app := range apps {
		iterationSet[app.Iteration] = struct{}{}
	}

	// Get iterations from calculations
	calcs, err := client.Calculation.Query().Select(calculation.FieldIteration).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query calculation iterations: %w", err)
	}
	for _, calc := range calcs {
		iterationSet[calc.Iteration] = struct{}{}
	}

	// Get iterations from drained results
	drained, err := client.DrainedResult.Query().Select(drainedresult.FieldIteration).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query drained result iterations: %w", err)
	}
	for _, dr := range drained {
		iterationSet[dr.Iteration] = struct{}{}
	}

	// Convert to sorted slice
	var iterations []int
	for iteration := range iterationSet {
		iterations = append(iterations, iteration)
	}

	// Sort for deterministic order
	for i := 0; i < len(iterations)-1; i++ {
		for j := i + 1; j < len(iterations); j++ {
			if iterations[i] > iterations[j] {
				iterations[i], iterations[j] = iterations[j], iterations[i]
			}
		}
	}

	return iterations, nil
}

func migrateIteration(client *ent.Client, ctx context.Context, iteration int) error {
	return utils.WithTx(ctx, client, func(tx *ent.Tx) error {
		// Find the earliest updated_at time for this iteration across all tables
		triggeredAt, err := getEarliestUpdatedAt(tx, ctx, iteration)
		if err != nil {
			return fmt.Errorf("failed to get earliest updated_at for iteration %d: %w", iteration, err)
		}

		// Create a new run
		run, err := tx.Run.Create().
			SetTriggeredAt(triggeredAt).
			SetPayloadMeta(map[string]any{
				"migrated_from_iteration": iteration,
				"migration_timestamp":     time.Now(),
			}).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create run for iteration %d: %w", iteration, err)
		}

		// Update all tables to reference this run
		if err := updateApplications(tx, ctx, iteration, run.ID); err != nil {
			return fmt.Errorf("failed to update applications for iteration %d: %w", iteration, err)
		}

		if err := updateCalculations(tx, ctx, iteration, run.ID); err != nil {
			return fmt.Errorf("failed to update calculations for iteration %d: %w", iteration, err)
		}

		if err := updateDrainedResults(tx, ctx, iteration, run.ID); err != nil {
			return fmt.Errorf("failed to update drained results for iteration %d: %w", iteration, err)
		}

		return nil
	})
}

func getEarliestUpdatedAt(tx *ent.Tx, ctx context.Context, iteration int) (time.Time, error) {
	// Get earliest time from applications
	apps, err := tx.Application.Query().
		Where(application.IterationEQ(iteration)).
		Select(application.FieldUpdatedAt).
		All(ctx)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to query applications: %w", err)
	}

	// Get earliest time from calculations
	calcs, err := tx.Calculation.Query().
		Where(calculation.IterationEQ(iteration)).
		Select(calculation.FieldUpdatedAt).
		All(ctx)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to query calculations: %w", err)
	}

	// Find the earliest time
	var earliest time.Time
	first := true

	for _, app := range apps {
		if first || app.UpdatedAt.Before(earliest) {
			earliest = app.UpdatedAt
			first = false
		}
	}

	for _, calc := range calcs {
		if first || calc.UpdatedAt.Before(earliest) {
			earliest = calc.UpdatedAt
			first = false
		}
	}

	if first {
		// No records found, use current time
		return time.Now(), nil
	}

	return earliest, nil
}

func updateApplications(tx *ent.Tx, ctx context.Context, iteration, runID int) error {
	_, err := tx.Application.Update().
		Where(application.IterationEQ(iteration)).
		SetRunID(runID).
		Save(ctx)
	return err
}

func updateCalculations(tx *ent.Tx, ctx context.Context, iteration, runID int) error {
	_, err := tx.Calculation.Update().
		Where(calculation.IterationEQ(iteration)).
		SetRunID(runID).
		Save(ctx)
	return err
}

func updateDrainedResults(tx *ent.Tx, ctx context.Context, iteration, runID int) error {
	_, err := tx.DrainedResult.Update().
		Where(drainedresult.IterationEQ(iteration)).
		SetRunID(runID).
		Save(ctx)
	return err
}
