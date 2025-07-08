package cmd

import (
	"analabit/cli/config"
	"analabit/cli/corestate"
	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"fmt"
	"os"
	"sort"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var headingCmd = &cobra.Command{
	Use:   "heading [varsity code or index] [heading index]",
	Short: "Prints detailed info about a specific heading and its admission results",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if !corestate.CrawlingDone {
			fmt.Println("Crawling not yet complete. Please wait.")
			return
		}
		corestate.ResultsMutex.RLock()
		defer corestate.ResultsMutex.RUnlock()

		varsityIdentifier := args[0]
		headingIdxArg := args[1]

		var targetVarsity *source.Varsity
		vIndex, err := strconv.Atoi(varsityIdentifier)
		if err == nil {
			if vIndex < 0 || vIndex >= len(corestate.LoadedVarsities) {
				fmt.Printf("Error: Invalid varsity index %d. Use 'varsities' command.\n", vIndex)
				return
			}
			targetVarsity = corestate.LoadedVarsities[vIndex]
		} else {
			found := false
			for _, v := range corestate.LoadedVarsities {
				if v.Code == varsityIdentifier {
					targetVarsity = v
					found = true
					break
				}
			}
			if !found {
				fmt.Printf("Error: Varsity with code '%s' not found. Use 'varsities' command.\n", varsityIdentifier)
				return
			}
		}

		if targetVarsity == nil || targetVarsity.VarsityCalculator == nil {
			fmt.Println("Error: Target varsity or its calculator is not properly loaded.")
			return
		}

		// Get headings from VarsityCalculator, sort them by pretty name to match 'headings' command output
		var sortedHeadings []*core.Heading
		// Iterate over the HeadingsCache map directly
		for _, h := range targetVarsity.VarsityCalculator.Headings() {
			sortedHeadings = append(sortedHeadings, h)
		}
		sort.Slice(sortedHeadings, func(i, j int) bool {
			return sortedHeadings[i].PrettyName() < sortedHeadings[j].PrettyName()
		})

		hIndex, err := strconv.Atoi(headingIdxArg)
		if err != nil || hIndex < 0 || hIndex >= len(sortedHeadings) {
			fmt.Printf("Error: Invalid heading index '%s'. Use 'headings %s' command to see available indexes.\n", headingIdxArg, varsityIdentifier)
			return
		}
		targetHeading := sortedHeadings[hIndex]

		fmt.Printf("Varsity: %s (%s)\n", targetVarsity.Name, targetVarsity.Code)
		fmt.Printf("Heading: %s (%s)\n", targetHeading.PrettyName(), targetHeading.Code())
		fmt.Printf("Total Capacity: %d\n\n", targetHeading.TotalCapacity())

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Drained\tPassing Score\tLast Admitted Rating Place")
		fmt.Fprintln(w, "-------\t-------------\t--------------------------")

		// Primary Results (DrainedPercent = 0)
		primaryResultsForVarsity, okPrimary := corestate.PrimaryResults[targetVarsity.Code]
		if !okPrimary {
			fmt.Println("Primary calculation results not found for this varsity.")
			w.Flush() // Flush before early return
			return
		}

		foundPrimaryResultForHeading := false
		for _, res := range primaryResultsForVarsity {
			if res.Heading.Code() == targetHeading.Code() {
				passingScore, psErr := res.PassingScore()
				lastAdmittedPlace, larpErr := res.LastAdmittedRatingPlace()
				psStr := "N/A"
				larpStr := "N/A"
				if psErr == nil {
					psStr = strconv.Itoa(passingScore)
				}
				if larpErr == nil {
					larpStr = "#" + strconv.Itoa(lastAdmittedPlace)
				}
				fmt.Fprintf(w, "--\t%s\t%s\n", psStr, larpStr)
				foundPrimaryResultForHeading = true
				break
			}
		}
		if !foundPrimaryResultForHeading {
			fmt.Fprintf(w, "--\tN/A\tN/A\n")
		}

		// Drained Results
		if !corestate.SimulationsDone && (corestate.TotalSimulations == 0 || corestate.CompletedSimulations.Load() < corestate.TotalSimulations) {
			fmt.Fprintln(w, "Simulation in progress...\t\t")
		} else {
			drainedResultsForVarsity, okDrained := corestate.DrainedResults[targetVarsity.Code]
			if !okDrained {
				fmt.Println("Drained simulation results not found for this varsity.")
			} else {
				// Sort stages for consistent output
				stages := make([]int, 0, len(config.AppConfig.DrainSim.Stages))
				for _, stage := range config.AppConfig.DrainSim.Stages { // Use configured stages
					stages = append(stages, stage)
				}
				sort.Ints(stages)

				for _, stagePercent := range stages {

					resultsForStage, okStage := drainedResultsForVarsity[stagePercent]
					if !okStage {
						fmt.Fprintf(w, "%d%%\t(No data)\t(No data)\n", stagePercent)
						continue
					}

					foundDrainedResultForHeading := false
					for _, res := range resultsForStage {

						if res.Heading.Code() == targetHeading.Code() {
							psStr := fmt.Sprintf("%d (Est.)", res.AvgPassingScore)
							larpStr := fmt.Sprintf("#%d (Est.)", res.AvgLastAdmittedRatingPlace)
							fmt.Fprintf(w, "%d%%\t%s\t%s\n", res.DrainedPercent, psStr, larpStr)
							foundDrainedResultForHeading = true
							break
						}
					}
					if !foundDrainedResultForHeading {
						fmt.Fprintf(w, "%d%%\tN/A (Est.)\tN/A (Est.)\n", stagePercent)
					}
				}
			}
		}
		w.Flush()
	},
}
