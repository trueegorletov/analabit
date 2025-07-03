package cmd

import (
	"analabit/cli/config"
	"analabit/cli/corestate"
	"analabit/cli/shell"
	"analabit/core/drainer"
	"analabit/core/registry"
	"analabit/core/source"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "analabit-cli",
	Short: "Analabit CLI provides tools for varsity data analysis and simulation.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		log.Println("rootCmd PersistentPreRunE: Entered")
		corestate.InitMutex.Lock()
		log.Println("rootCmd PersistentPreRunE: InitMutex locked") // New log
		defer corestate.InitMutex.Unlock()

		if corestate.Initialized { // Directly access the exported field
			log.Println("rootCmd PersistentPreRunE: Already initialized, returning previous error (if any).") // New log
			return corestate.InitError
		}

		log.Println("rootCmd PersistentPreRunE: Calling corestate.InitializeState()") // New log
		corestate.InitializeState()

		cfgPath, _ := cmd.Flags().GetString("config")
		log.Printf("rootCmd PersistentPreRunE: Loading config from: %s", cfgPath) // New log
		if err := config.LoadConfig(cfgPath); err != nil {
			corestate.InitError = fmt.Errorf("failed to load configuration: %w", err)
			log.Printf("rootCmd PersistentPreRunE: Error loading config: %v", corestate.InitError) // New log
			return corestate.InitError
		}
		log.Println("Configuration loaded.")

		if err := performCrawling(); err != nil {
			corestate.InitError = fmt.Errorf("failed during crawling: %w", err)
			log.Printf("rootCmd PersistentPreRunE: Error during crawling: %v", corestate.InitError) // New log
			return corestate.InitError
		}
		corestate.CrawlingDone = true
		log.Printf("Crawling finished, varsities loaded: %d\n", len(corestate.LoadedVarsities))

		performPrimaryCalculations()
		log.Println("Primary calculations finished.")

		go startDrainSimulations()
		log.Println("Drain simulations started in background.")

		corestate.Initialized = true                                       // Corrected: Directly set the exported field
		log.Println("rootCmd PersistentPreRunE: Initialization complete.") // New log
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !corestate.Initialized { // Use the exported field directly
			if corestate.InitError != nil {
				fmt.Fprintf(os.Stderr, "Initialization failed: %s\n", corestate.InitError)
				return
			}
			fmt.Fprintln(os.Stderr, "CLI is not initialized. This should not happen.")
			return
		}
		shell.Run(cmd) // Pass the rootCmd (which is 'cmd' in this context) to shell.Run
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	log.Println("rootCmd init(): Starting") // New log
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is ./cli/config.toml, ./config.toml)")
	rootCmd.AddCommand(varsitiesCmd)
	rootCmd.AddCommand(headingsCmd)
	rootCmd.AddCommand(headingCmd)
	rootCmd.AddCommand(studentCmd)
	rootCmd.AddCommand(progressCmd)
	rootCmd.AddCommand(uploadCmd)
	log.Println("rootCmd init(): Finished") // New log
}

func performCrawling() error {
	params := registry.CrawlOptions{
		VarsitiesList:    config.AppConfig.Varsities.List,
		VarsitiesExclude: config.AppConfig.Varsities.Excluded,
		CacheDir:         config.AppConfig.Cache.Directory,
		CacheTTLMinutes:  config.AppConfig.Cache.TTLMinutes,
		DrainStages:      config.AppConfig.DrainSim.Stages,
		DrainIterations:  config.AppConfig.DrainSim.Iterations,
	}
	result, err := registry.CrawlWithOptions(registry.AllDefinitions, params)
	if err != nil {
		return err
	}
	corestate.LoadedVarsities = result.LoadedVarsities
	return nil
}

func performPrimaryCalculations() {
	log.Println("performPrimaryCalculations: Starting...")
	corestate.ResultsMutex.Lock()
	log.Println("performPrimaryCalculations: ResultsMutex locked") // New log
	defer func() {
		log.Println("performPrimaryCalculations: Unlocking ResultsMutex...") // New log
		corestate.ResultsMutex.Unlock()
	}()

	for i, v := range corestate.LoadedVarsities {
		log.Printf("performPrimaryCalculations: Processing varsity %d/%d: %s (%s) - Before Clone", i+1, len(corestate.LoadedVarsities), v.Name, v.Code) // New log
		clonedVarsity := v.Clone()
		log.Printf("performPrimaryCalculations: Processing varsity %s (%s) - After Clone, Before CalculateAdmissions", v.Name, v.Code) // New log
		results := clonedVarsity.VarsityCalculator.CalculateAdmissions()
		log.Printf("performPrimaryCalculations: Processing varsity %s (%s) - After CalculateAdmissions, Before storing results", v.Name, v.Code) // New log
		corestate.PrimaryResults[v.Code] = results
		log.Printf("performPrimaryCalculations: Processing varsity %s (%s) - Results stored", v.Name, v.Code) // New log
		//panic("for test")
	}
	log.Println("performPrimaryCalculations: Finished.")
}

func startDrainSimulations() {
	log.Println("startDrainSimulations: Starting...") // New log
	startTime := time.Now()                           // Record start time
	numVarsities := len(corestate.LoadedVarsities)
	numStages := len(config.AppConfig.DrainSim.Stages)
	if numVarsities == 0 || numStages == 0 {
		log.Println("No varsities or drain stages configured. Skipping drain simulations.")
		corestate.SimulationsDone = true
		log.Println("startDrainSimulations: Finished (no simulations needed).") // New log
		return
	}
	corestate.TotalSimulations = int32(numVarsities * numStages)
	corestate.CompletedSimulations.Store(0)

	var wg sync.WaitGroup
	// Add the total number of expected simulation tasks (stage goroutines) upfront
	// to prevent a race condition where wg.Wait() might be called before wg.Add().
	if corestate.TotalSimulations > 0 {
		wg.Add(int(corestate.TotalSimulations))
	}

	for _, loopVarsity := range corestate.LoadedVarsities { // Renamed varsity to loopVarsity for clarity
		go func(currentVarsity *source.Varsity) { // Launch a goroutine for each varsity
			corestate.ResultsMutex.Lock() // Lock before accessing/modifying DrainedResults for this specific varsity
			if _, ok := corestate.DrainedResults[currentVarsity.Code]; !ok {
				corestate.DrainedResults[currentVarsity.Code] = make(map[int][]drainer.DrainedResult)
			}
			corestate.ResultsMutex.Unlock()

			for _, stage := range config.AppConfig.DrainSim.Stages { // Iterate over stages for the current varsity
				go func(v *source.Varsity, s int) { // Launch a goroutine for each simulation task (varsity-stage pair)
					defer wg.Done() // Decrement WaitGroup counter when this simulation task is done

					// Drainer.New takes the prototype; its Run method clones it internally.
					drainerInstance := drainer.New(v, s)
					drainedResultSlice := drainerInstance.Run(config.AppConfig.DrainSim.Iterations)

					corestate.ResultsMutex.Lock()
					corestate.DrainedResults[v.Code][s] = drainedResultSlice // Store results
					corestate.ResultsMutex.Unlock()

					corestate.CompletedSimulations.Add(1)
					log.Printf("Simulation finished for %s, stage %d%% (%d/%d total)\\\\n", v.Name, s, corestate.CompletedSimulations.Load(), corestate.TotalSimulations)
				}(currentVarsity, stage) // Pass the currentVarsity (from the outer goroutine's scope) and stage
			}
		}(loopVarsity) // Pass the loopVarsity to the goroutine to capture its current value
	}

	log.Println("startDrainSimulations: Initiated parallel simulation processing per varsity. Waiting for all individual simulations to complete...") // Updated log message
	if corestate.TotalSimulations > 0 {                                                                                                               // Only wait if there were simulations to run
		wg.Wait()
	}
	corestate.SimulationsDone = true
	elapsedTime := time.Since(startTime)                                        // Calculate elapsed time
	fmt.Printf("All drain simulations finished in %s.\n", elapsedTime)          // Print to stdout
	log.Printf("All drain simulations finished. Elapsed time: %s", elapsedTime) // Log with elapsed time
}

// Individual command files (varsities.go, headings.go, etc.) will be created separately.
