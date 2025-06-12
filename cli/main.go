package main

import (
	"analabit/cli/cmd"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

var logFile *os.File // Global variable to keep the log file open

func main() {
	// Load configuration early to set up logging
	// We need to parse flags to get the config path *before* Cobra does it fully,
	// or make a preliminary pass, or rely on default config path for initial log setup.
	// For simplicity, we'll try to load config with a default path first for logging setup.
	// Cobra's PersistentPreRunE in root.go will load it again properly with flag parsing.

	// Attempt to load config to get logging.File path
	// This is a bit of a chicken-and-egg, as flags are parsed by Cobra later.
	// A more robust solution might involve a two-pass flag parsing or a dedicated log config flag.
	prelimCfgPath := "./cli/config.toml" // Default or common path
	if _, err := os.Stat(prelimCfgPath); os.IsNotExist(err) {
		prelimCfgPath = "config.toml" // Try root path
	}

	// Load config just for the logging path initially
	tempCfg := struct {
		Logging struct {
			File string `mapstructure:"file"`
		}
	}{}

	// Use a temporary viper instance for this initial load
	// to avoid interfering with the main viper instance in config.LoadConfig
	v := viper.New()
	v.SetConfigFile(prelimCfgPath)
	if err := v.ReadInConfig(); err == nil {
		v.Unmarshal(&tempCfg)
	}

	if tempCfg.Logging.File != "" {
		var err error
		logFile, err = os.OpenFile(tempCfg.Logging.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			log.Printf("Error opening global log file '%s': %v. Logging to stderr.", tempCfg.Logging.File, err)
		} else {
			log.SetOutput(logFile)
			log.Printf("Global logging redirected to file: %s", tempCfg.Logging.File)
		}
	} else {
		log.Println("Global logging configured to: stderr (default or no log file specified)")
	}

	// Ensure log file is closed on exit
	defer func() {
		if logFile != nil {
			log.Println("Closing global log file.")
			err := logFile.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error closing global log file: %v\n", err)
			}
		}
	}()

	log.Println("CLI main: Starting application...")
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'\n", err)
		log.Printf("CLI main: Error during cmd.Execute(): %v", err)
		os.Exit(1)
	}
	log.Println("CLI main: Application finished successfully.")
}
