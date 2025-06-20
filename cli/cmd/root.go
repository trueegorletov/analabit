package cmd

import (
	"analabit/cli/config"
	"analabit/cli/corestate"
	"analabit/cli/shell"
	"analabit/core/drainer"
	"analabit/core/registry"
	"analabit/core/source"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
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
	allDefs := registry.AllDefinitions
	var filteredDefs []source.VarsityDefinition

	log.Println("performCrawling: Starting...") // New log

	varsitiesToUse := make(map[string]bool)
	if len(config.AppConfig.Varsities.List) == 1 && config.AppConfig.Varsities.List[0] == "all" {
		for _, def := range allDefs {
			varsitiesToUse[def.Code] = true
		}
	} else {
		for _, code := range config.AppConfig.Varsities.List {
			varsitiesToUse[code] = true
		}
	}
	for _, code := range config.AppConfig.Varsities.Excluded {
		delete(varsitiesToUse, code)
	}

	for _, def := range allDefs {
		if varsitiesToUse[def.Code] {
			filteredDefs = append(filteredDefs, def)
		}
	}

	if len(filteredDefs) == 0 {
		log.Println("No varsities selected after filtering. Skipping crawling.")
		corestate.LoadedVarsities = []*source.Varsity{}
		log.Println("performCrawling: Finished (no varsities selected).") // New log
		return nil
	}

	cacheDir := config.AppConfig.Cache.Directory
	ttlSeconds := int64(config.AppConfig.Cache.TTLMinutes * 60)
	var validCacheFile string
	var latestTimestamp int64 = -1

	log.Printf("performCrawling: Checking cache directory '%s' with TTL %d minutes.", cacheDir, config.AppConfig.Cache.TTLMinutes) // New log
	if _, err := os.Stat(cacheDir); !os.IsNotExist(err) {
		log.Println("performCrawling: Walking cache directory...") // New log
		err := filepath.WalkDir(cacheDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				log.Printf("Error walking cache directory at %s: %v", path, err)
				return err // Stop walking this path if error occurs
			}
			if !d.IsDir() && strings.HasSuffix(d.Name(), ".gob") { // Ensure it's .gob, not .dob
				nameWithoutExt := strings.TrimSuffix(d.Name(), ".gob")
				ts, err := strconv.ParseInt(nameWithoutExt, 10, 64)
				if err == nil {
					if time.Now().Unix()-ts < ttlSeconds && ts > latestTimestamp {
						latestTimestamp = ts
						validCacheFile = path
					}
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("performCrawling: Error during cache directory walk: %v", err) // Log error from WalkDir itself
		}
		log.Println("performCrawling: Finished walking cache directory.") // New log
	}

	var loadedVarsities []*source.Varsity
	log.Println("performCrawling: About to load varsity data...") // New log

	if validCacheFile != "" {
		log.Printf("Attempting to use cache file: %s\n", validCacheFile)
		file, err := os.Open(validCacheFile)
		if err != nil {
			log.Printf("Failed to open cache file %s: %v. Falling back to full crawl.\n", validCacheFile, err)
			log.Println("performCrawling: Attempting to load from definitions (cache open failed)...") // New log
			loadedVarsities = source.LoadFromDefinitions(filteredDefs)
			log.Println("performCrawling: Finished loading from definitions (cache open failed).") // New log
		} else {
			defer file.Close()
			caches, err := source.DeserializeList(file)
			if err != nil {
				log.Printf("Failed to deserialize cache file %s: %v. Falling back to full crawl.\n", validCacheFile, err)
				log.Println("performCrawling: Attempting to load from definitions (cache deserialize failed)...") // New log
				loadedVarsities = source.LoadFromDefinitions(filteredDefs)
				log.Println("performCrawling: Finished loading from definitions (cache deserialize failed).") // New log
			} else {
				log.Println("Successfully deserialized cache. Loading with caches.")
				log.Println("performCrawling: Attempting to load with caches...") // New log
				loadedVarsities = source.LoadWithCaches(filteredDefs, caches)
				log.Println("performCrawling: Finished loading with caches.") // New log
			}
		}
	} else {
		log.Println("No valid cache file found or cache is disabled/empty. Performing full crawl.")
		log.Println("performCrawling: Attempting to load from definitions (no cache)...") // New log
		loadedVarsities = source.LoadFromDefinitions(filteredDefs)
		log.Println("performCrawling: Finished loading from definitions (no cache).") // New log
	}

	log.Printf("performCrawling: Data loaded. Number of varsities: %d. About to save cache...", len(loadedVarsities)) // New log
	if len(loadedVarsities) > 0 {
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			return fmt.Errorf("failed to create cache directory %s: %w", cacheDir, err)
		}
		newCacheFilename := filepath.Join(cacheDir, fmt.Sprintf("%d.gob", time.Now().Unix()))
		file, err := os.Create(newCacheFilename)
		if err != nil {
			log.Printf("Failed to create new cache file %s: %v. Proceeding without saving cache this time.\n", newCacheFilename, err)
		} else {
			defer file.Close()
			var cachesToSave []*source.VarsityDataCache
			for _, v := range loadedVarsities {
				if v.VarsityDataCache != nil {
					cachesToSave = append(cachesToSave, v.VarsityDataCache)
				}
			}
			if err := source.SerializeList(cachesToSave, file); err != nil {
				log.Printf("Failed to serialize data to cache file %s: %v. Proceeding without saving cache this time.\n", newCacheFilename, err)
			} else {
				log.Printf("Saved new cache to %s\n", newCacheFilename)
			}
		}
	}

	sort.Slice(loadedVarsities, func(i, j int) bool {
		return loadedVarsities[i].Name < loadedVarsities[j].Name
	})
	corestate.LoadedVarsities = loadedVarsities
	log.Println("performCrawling: Finished successfully.") // New log
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
