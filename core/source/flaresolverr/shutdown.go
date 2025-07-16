package flaresolverr

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	shutdownOnce sync.Once
	shutdownChan = make(chan struct{})
)

// InitGracefulShutdown sets up graceful shutdown handling for FlareSolverr sessions
// This should be called once during application startup
func InitGracefulShutdown() {
	shutdownOnce.Do(func() {
		go handleShutdownSignals()
	})
}

// handleShutdownSignals handles graceful shutdown of FlareSolverr sessions
func handleShutdownSignals() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutdown signal received, cleaning up FlareSolverr sessions...")

	// Use the iteration-level cleanup for application shutdown
	if err := StopForIteration(); err != nil {
		log.Printf("Failed to cleanup FlareSolverr sessions during shutdown: %v", err)
	} else {
		log.Println("FlareSolverr sessions cleaned up successfully")
	}

	// Signal that shutdown is complete
	close(shutdownChan)
}

// WaitForShutdown blocks until graceful shutdown is complete
// This can be used by the main application to wait for cleanup
func WaitForShutdown() {
	<-shutdownChan
}

// IsShuttingDown returns true if the application is in the process of shutting down
func IsShuttingDown() bool {
	select {
	case <-shutdownChan:
		return true
	default:
		return false
	}
}
