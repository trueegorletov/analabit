package cmd

import (
	"analabit/cli/config"
	"analabit/cli/corestate"
	"analabit/core"
	"analabit/core/ent"
	"analabit/core/upload"
	"context"
	"fmt"
	"log"

	"analabit/core/drainer"

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
			// The upload.Primary function expects the *core.VarsityCalculator that contains the student applications.
			// The `results` are the CalculationResult from this calculator.
			if err := upload.Primary(ctx, client, targetVarsityCalculator, results); err != nil {
				log.Printf("Error uploading primary results for varsity %s: %v", varsityCode, err)
			} else {
				fmt.Printf("Successfully uploaded primary results for %s.\n", varsityCode)
			}
		}
		corestate.ResultsMutex.RUnlock()
		fmt.Println("Primary results upload finished.")

		// Upload Drained Results
		fmt.Println("Uploading drained simulation results...")
		corestate.ResultsMutex.RLock()
		for varsityCode, stageMap := range corestate.DrainedResults {
			for stage, results := range stageMap {
				if len(results) == 0 {
					continue // Skip if no results for this stage
				}

				// Convert drainer.DrainedResult to core.DrainedResultDTO
				drainedDTOs := drainer.NewDrainedResultDTOs(results)

				// Call the refactored upload.DrainedResults function with DTOs
				if err := upload.DrainedResults(ctx, client, drainedDTOs); err != nil {
					log.Printf("Error uploading drained results for varsity %s, stage %d%%: %v", varsityCode, stage, err)
				} else {
					fmt.Printf("Successfully uploaded drained results for %s, stage %d%%.\n", varsityCode, stage)
				}
			}
		}
		corestate.ResultsMutex.RUnlock()
		fmt.Println("Drained simulation results upload finished.")

		fmt.Println("Upload process complete.")
	},
}
