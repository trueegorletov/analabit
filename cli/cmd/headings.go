package cmd

import (
	"github.com/trueegorletov/analabit/cli/corestate"
	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"fmt"
	"os"
	"sort"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var headingsCmd = &cobra.Command{
	Use:   "headings [varsity code or index]",
	Short: "Prints all headings of the specified varsity",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !corestate.CrawlingDone {
			fmt.Println("Crawling not yet complete. Please wait.")
			return
		}
		corestate.ResultsMutex.RLock()
		defer corestate.ResultsMutex.RUnlock()

		varsityIdentifier := args[0]
		var targetVarsity *source.Varsity

		index, err := strconv.Atoi(varsityIdentifier)
		if err == nil {
			if index < 0 || index >= len(corestate.LoadedVarsities) {
				fmt.Printf("Error: Invalid varsity index %d. Use 'varsities' command to see available indexes.\n", index)
				return
			}
			targetVarsity = corestate.LoadedVarsities[index]
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
				fmt.Printf("Error: Varsity with code '%s' not found. Use 'varsities' command to see available codes.\n", varsityIdentifier)
				return
			}
		}

		if targetVarsity == nil || targetVarsity.VarsityCalculator == nil {
			fmt.Println("Error: Target varsity or its calculator is not properly loaded.")
			return
		}

		// Get headings from the VarsityCalculator, as these are the ones used in calculations
		// The VarsityCalculator stores *core.Heading in a map called HeadingsCache
		var headingsToSort []*core.Heading
		// Iterate over the HeadingsCache map directly
		for _, h := range targetVarsity.VarsityCalculator.Headings() {
			headingsToSort = append(headingsToSort, h)
		}

		if len(headingsToSort) == 0 {
			fmt.Printf("No headings found for varsity %s (%s).\n", targetVarsity.Name, targetVarsity.Code)
			return
		}

		sort.Slice(headingsToSort, func(i, j int) bool {
			return headingsToSort[i].PrettyName() < headingsToSort[j].PrettyName()
		})

		fmt.Printf("HeadingsCache for %s (%s):\n", targetVarsity.Name, targetVarsity.Code)
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Index\tCode\tPretty Name\tTotal Capacity")
		fmt.Fprintln(w, "-----\t----\t-----------\t--------------")
		for i, h := range headingsToSort {
			fmt.Fprintf(w, "%d\t%s\t%s\t%d\n", i, h.Code(), h.PrettyName(), h.TotalCapacity())
		}
		w.Flush()
	},
}
