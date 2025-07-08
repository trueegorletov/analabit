package cmd

import (
	"analabit/cli/config"
	"analabit/cli/corestate"
	"analabit/core"
	"analabit/core/source"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var studentCmd = &cobra.Command{
	Use:   "student [student ID]",
	Short: "Prints admission info for a specific student across varsities and headings",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !corestate.CrawlingDone {
			fmt.Println("Crawling not yet complete. Please wait.")
			return
		}
		corestate.ResultsMutex.RLock()
		defer corestate.ResultsMutex.RUnlock()

		studentID := args[0]

		// --- Primary Results Table ---
		fmt.Printf("Admission status for student ID: %s\n\n", studentID)
		fmt.Println("Primary Results (Drained: --):")

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
		var tableHeader []string
		var primaryRow []string
		var studentFoundInPrimary bool

		// Collect all varsities where student is admitted in primary results
		// and the corresponding headings. Keep varsity order same as LoadedVarsities.
		varsityToPrimaryHeading := make(map[string]string)               // varsityCode -> headingPrettyName
		studentPrimaryApplications := make(map[string]*core.Application) // varsityCode -> student's winning Application

		for _, varsity := range corestate.LoadedVarsities { // Iterate in fixed order
			primaryResultsForVarsity, ok := corestate.PrimaryResults[varsity.Code]
			if !ok {
				continue
			}
			studentObj := varsity.VarsityCalculator.GetStudent(studentID)
			if studentObj == nil || studentObj.Quit() { // Check if student exists in this varsity and hasn't quit
				continue
			}

			for _, calcResult := range primaryResultsForVarsity {
				for _, admittedStudent := range calcResult.Admitted {
					if admittedStudent.ID() == studentID {
						studentFoundInPrimary = true
						varsityToPrimaryHeading[varsity.Code] = calcResult.Heading.PrettyName()
						// Find the application for this admission
						for _, app := range studentObj.Applications() {
							if app.Heading().Code() == calcResult.Heading.Code() {
								studentPrimaryApplications[varsity.Code] = app
								break
							}
						}
						break // Found student in this calcResult
					}
				}
				if _, exists := varsityToPrimaryHeading[varsity.Code]; exists {
					break // Student found in this varsity, move to next varsity
				}
			}
		}

		if !studentFoundInPrimary {
			fmt.Println("Student not found in any primary admission results.")
		} else {
			tableHeader = append(tableHeader, "Drained")
			primaryRow = append(primaryRow, "--")
			for _, varsity := range corestate.LoadedVarsities { // Keep consistent order
				if headingName, ok := varsityToPrimaryHeading[varsity.Code]; ok {
					tableHeader = append(tableHeader, varsity.Name)
					primaryRow = append(primaryRow, headingName)
				}
			}
			fmt.Fprintln(w, strings.Join(tableHeader, "\t"))
			fmt.Fprintln(w, strings.Join(primaryRow, "\t"))
		}
		w.Flush() // Flush primary results table

		// --- Drained Results Table ---
		fmt.Println("\nDrained Simulation Results:")
		if !corestate.SimulationsDone && (corestate.TotalSimulations == 0 || corestate.CompletedSimulations.Load() < corestate.TotalSimulations) {
			fmt.Println("Simulation in progress...")
			return
		}
		if !studentFoundInPrimary && !corestate.SimulationsDone { // If not in primary and sims not done, no point proceeding for drained
			fmt.Println("Student not in primary results, and simulations might still be running or produced no data.")
			return
		}

		wDrained := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
		if len(tableHeader) > 0 { // Print header only if there was one from primary
			fmt.Fprintln(wDrained, strings.Join(tableHeader, "\t"))
		} else { // Construct a header if student wasn't in primary but might be in drained
			tableHeader = append(tableHeader, "Drained")
			// Add all varsity names to header for potential drained results
			// This is a bit speculative but ensures columns exist if data appears.
			tempVarsityNames := make(map[string]string) // code -> name
			for _, v := range corestate.LoadedVarsities {
				tempVarsityNames[v.Code] = v.Name
			}
			// Sort varsity codes to ensure consistent header order if built here
			var codes []string
			for code := range tempVarsityNames {
				codes = append(codes, code)
			}
			sort.Strings(codes)
			for _, code := range codes {
				tableHeader = append(tableHeader, tempVarsityNames[code])
			}
			fmt.Fprintln(wDrained, strings.Join(tableHeader, "\t"))
		}

		stages := make([]int, 0, len(config.AppConfig.DrainSim.Stages))
		for _, stage := range config.AppConfig.DrainSim.Stages {
			stages = append(stages, stage)
		}
		sort.Ints(stages)

		for _, stagePercent := range stages {
			drainedRow := make([]string, len(tableHeader))
			drainedRow[0] = fmt.Sprintf("%d%%", stagePercent)
			foundForStage := false

			for i, th := range tableHeader {
				if i == 0 {
					continue
				} // Skip "Drained" column title
				// Find varsity by name (th)
				var currentVarsity *source.Varsity
				for _, v := range corestate.LoadedVarsities {
					if v.Name == th {
						currentVarsity = v
						break
					}
				}
				if currentVarsity == nil {
					drainedRow[i] = "--"
					continue
				}

				studentObj := currentVarsity.VarsityCalculator.GetStudent(studentID)
				if studentObj == nil || studentObj.Quit() { // Student might have quit in this specific simulation clone
					drainedRow[i] = "--"
					continue
				}

				drainedResultsForVarsity, okVarsity := corestate.DrainedResults[currentVarsity.Code]
				if !okVarsity {
					drainedRow[i] = "--"
					continue
				}
				resultsForStage, okStage := drainedResultsForVarsity[stagePercent]
				if !okStage {
					drainedRow[i] = "--"
					continue
				}

				// Find student's application and check against DrainedResult.AvgLastAdmittedRatingPlace
				admittedToHeadingName := "--"
				studentApplications := studentObj.Applications() // ApplicationsCache are sorted by priority

				for _, app := range studentApplications {
					appRatingPlace := app.RatingPlace()
					for _, drainedRes := range resultsForStage {
						if drainedRes.Heading.Code() == app.Heading().Code() {
							if appRatingPlace <= drainedRes.AvgLastAdmittedRatingPlace {
								admittedToHeadingName = drainedRes.Heading.PrettyName()
								foundForStage = true
								break
							}
						}
					}
					if admittedToHeadingName != "--" { // Found admission for this student in this varsity for this stage
						break
					}
				}
				drainedRow[i] = admittedToHeadingName
			}
			if foundForStage || studentFoundInPrimary { // Print row if student was in primary or found in this drained stage
				fmt.Fprintln(wDrained, strings.Join(drainedRow, "\t"))
			}
		}
		wDrained.Flush()
	},
}
