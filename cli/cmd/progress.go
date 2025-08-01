package cmd

import (
	"github.com/trueegorletov/analabit/cli/corestate"
	"fmt"

	"github.com/spf13/cobra"
)

var progressCmd = &cobra.Command{
	Use:   "progress",
	Short: "Prints the current progress of background drain simulations",
	Run: func(cmd *cobra.Command, args []string) {
		if !corestate.CrawlingDone {
			fmt.Println("Crawling has not started or completed yet.")
			return
		}

		if corestate.SimulationsDone {
			fmt.Println("All drain simulations have finished.")
			completed := corestate.CompletedSimulations.Load()
			total := corestate.TotalSimulations
			if total > 0 {
				fmt.Printf("Completed %d out of %d simulations.\n", completed, total)
			} else {
				fmt.Println("No simulations were scheduled (e.g., no varsities or stages defined).")
			}
			return
		}

		completed := corestate.CompletedSimulations.Load()
		total := corestate.TotalSimulations

		if total == 0 { // Should be caught by SimulationsDone if truly 0, but as a safeguard
			fmt.Println("Drain simulations are not scheduled (total is 0).")
			return
		}

		percentage := 0.0
		if total > 0 { // Avoid division by zero if total somehow becomes 0 mid-process
			percentage = (float64(completed) / float64(total)) * 100
		}
		fmt.Printf("Drain simulation progress: %.2f%% (%d/%d simulations completed)\n", percentage, completed, total)
	},
}
