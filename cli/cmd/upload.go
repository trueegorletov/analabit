package cmd

import (
	"github.com/trueegorletov/analabit/cli/config"
	"github.com/trueegorletov/analabit/cli/corestate"
	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/database"
	"github.com/trueegorletov/analabit/core/ent"
	"github.com/trueegorletov/analabit/core/upload"
	"context"
	"fmt"
	"log"

	"github.com/trueegorletov/analabit/core/drainer"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Uploads all collected primary and drained results to the database",
	Run: func(cmd *cobra.Command, args []string) {
		if !corestate.CrawlingDone {
			fmt.Println("Error: Crawling must be completed before uploading.")
			return
		}
		if !corestate.SimulationsDone {
			fmt.Println("Error: Drain simulations must be completed before uploading. Use 'progress' to check status.")
			return
		}

		fmt.Println("Starting upload process...")

		// Database Connection
		dbCfg := config.AppConfig.Upload.Database
		connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			dbCfg.Host, dbCfg.Port, dbCfg.User, dbCfg.Password, dbCfg.DBName)

		client, err := ent.Open("postgres", connStr)
		if err != nil {
			log.Fatalf("Failed opening connection to postgres: %v", err)
		}
		defer client.Close()

		ctx := context.Background()

		// Run the auto migration tool.
		if err := client.Schema.Create(ctx); err != nil {
			log.Fatalf("Failed creating schema resources: %v", err)
		}
		fmt.Println("Database schema migration/check complete.")

		// Create a new run record for this upload session
		run, err := client.Run.Create().Save(ctx)
		if err != nil {
			log.Fatalf("Failed to create run record: %v", err)
		}
		fmt.Printf("Created new run with ID: %d\n", run.ID)

		// Upload Primary Results
		fmt.Println("Uploading primary results...")
		corestate.ResultsMutex.RLock() // Ensure read lock while accessing results
		for varsityCode, results := range corestate.PrimaryResults {
			var targetVarsityCalculator *core.VarsityCalculator
			for _, v := range corestate.LoadedVarsities {
				if v.Code == varsityCode {
					targetVarsityCalculator = v.VarsityCalculator // This is the original calculator from loading phase
					break
				}
			}
			if targetVarsityCalculator == nil {
				log.Printf("Warning: VarsityCalculator not found for code %s when uploading primary results. Skipping.", varsityCode)
				continue
			}

			// Convert VarsityCalculator and results to UploadPayload
			// We need to convert drained results to the right format for this varsity
			var drainedDTOs map[int][]core.DrainedResultDTO
			corestate.ResultsMutex.RLock()
			if stageMap, exists := corestate.DrainedResults[varsityCode]; exists {
				drainedDTOs = make(map[int][]core.DrainedResultDTO)
				for stage, drainedResults := range stageMap {
					if len(drainedResults) > 0 {
						drainedDTOs[stage] = drainer.NewDrainedResultDTOs(drainedResults)
					}
				}
			}
			corestate.ResultsMutex.RUnlock()

			payload := core.NewUploadPayloadFromCalculator(targetVarsityCalculator, results, drainedDTOs)

			// Call the updated upload.Primary function with runID and payload
			if err := upload.Primary(ctx, client, run.ID, payload); err != nil {
				log.Printf("Error uploading primary results for varsity %s: %v", varsityCode, err)
			} else {
				fmt.Printf("Successfully uploaded primary results for %s.\n", varsityCode)
			}
		}
		corestate.ResultsMutex.RUnlock()
		fmt.Println("Primary results upload finished.")

		// Upload Drained Results (if not already uploaded via Primary)
		fmt.Println("Uploading remaining drained simulation results...")
		corestate.ResultsMutex.RLock()
		for varsityCode, stageMap := range corestate.DrainedResults {
			// Skip if we already uploaded drained results for this varsity via Primary
			if _, primaryExists := corestate.PrimaryResults[varsityCode]; primaryExists {
				continue
			}

			for stage, results := range stageMap {
				if len(results) == 0 {
					continue // Skip if no results for this stage
				}

				// Convert drainer.DrainedResult to core.DrainedResultDTO
				drainedDTOs := drainer.NewDrainedResultDTOs(results)

				// Call the refactored upload.DrainedResults function with runID and DTOs
				if err := upload.DrainedResults(ctx, client, run.ID, drainedDTOs); err != nil {
					log.Printf("Error uploading drained results for varsity %s, stage %d%%: %v", varsityCode, stage, err)
				} else {
					fmt.Printf("Successfully uploaded drained results for %s, stage %d%%.\n", varsityCode, stage)
				}
			}
		}
		corestate.ResultsMutex.RUnlock()
		fmt.Println("Drained simulation results upload finished.")

		// Refresh materialized views after all uploads are complete
		fmt.Println("Refreshing materialized views...")
		dbClient, err := database.NewClient(client)
		if err != nil {
			log.Printf("Warning: Failed to create database client for view refresh: %v", err)
		} else {
			if err := upload.RefreshMaterializedViews(ctx, dbClient); err != nil {
				log.Printf("Warning: Failed to refresh materialized views: %v", err)
			} else {
			fmt.Println("Materialized views refreshed successfully.")
			}
		}

		// Mark run as finished
		_, err = client.Run.UpdateOneID(run.ID).SetFinished(true).Save(ctx)
		if err != nil {
			log.Printf("Warning: Failed to mark run as finished: %v", err)
		} else {
			fmt.Printf("Run %d marked as finished.\n", run.ID)
			// Run cleanup job
			if cleanupErr := dbClient.PerformBackupAndCleanup(ctx, config.AppConfig.Cleanup.RetentionRuns, config.AppConfig.Cleanup.BackupDir); cleanupErr != nil {
				log.Printf("Warning: Cleanup job failed for run %d: %v", run.ID, cleanupErr)
			}
		}

		fmt.Printf("Upload process complete. All data uploaded under run ID: %d\n", run.ID)
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
