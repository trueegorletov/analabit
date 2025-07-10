package main

import (
	"os"
	"path/filepath"

	"github.com/trueegorletov/analabit/core/source/mipt"
)

func main() {
	// Change to the directory where Go modules are properly configured
	dir, _ := os.Getwd()
	if filepath.Base(dir) != "analabit" {
		// If we're not in the right directory, try to change to it
		os.Chdir("/home/yegor/Prestart/analabit")
	}

	// Run the original detection test
	mipt.RunOriginalDetectionTest()
}
