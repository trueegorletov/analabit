package flaresolverr

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
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

// handleShutdownSignals listens for termination signals and triggers graceful shutdown
func handleShutdownSignals() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutdown signal received, cleaning up FlareSolverr sessions...")

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown session manager if it exists
	if sessionManager != nil {
		if err := sessionManager.Shutdown(ctx); err != nil {
			log.Printf("Error during session manager shutdown: %v", err)
		} else {
			log.Println("FlareSolverr sessions cleaned up successfully")
		}
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