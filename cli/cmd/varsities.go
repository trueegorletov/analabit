package cmd

import (
	"github.com/trueegorletov/analabit/cli/corestate"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var varsitiesCmd = &cobra.Command{
	Use:   "varsities",
	Short: "Prints all loaded varsities with their codes and numerical indexes",
	Run: func(cmd *cobra.Command, args []string) {
		if !corestate.CrawlingDone {
			fmt.Println("Crawling not yet complete. Please wait.")
			return
		}
		corestate.ResultsMutex.RLock() // Use RLock for read-only access
		defer corestate.ResultsMutex.RUnlock()

		if len(corestate.LoadedVarsities) == 0 {
			fmt.Println("No varsities loaded or selected.")
			return
		}

		fmt.Println("Loaded Varsities:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Index\tCode\tPretty Name")
		fmt.Fprintln(w, "-----\t----\t-----------")
		for i, v := range corestate.LoadedVarsities {
			fmt.Fprintf(w, "%d\t%s\t%s\n", i, v.Code, v.Name)
		}
		w.Flush()
	},
}
