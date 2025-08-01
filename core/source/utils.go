package source

import (
	"context"
	"log/slog"
	"net"
	"strings"
	"time"
)

// AcquireHTTPSemaphores acquires from per-varsity (if exists) and global semaphores, returning a release function.
func AcquireHTTPSemaphores(ctx context.Context, varsityCode string) (func(), error) {
	var releases []func()

	// Acquire per-varsity if exists
	if sem, ok := VarsitySemaphores[varsityCode]; ok {
		acquireCtx := ctx
		if varsityCode == "spbsu" {
			acquireCtx = context.Background()
		}
		if err := sem.Acquire(acquireCtx, 1); err != nil {
			slog.Error("Failed to acquire varsity semaphore", "varsity", varsityCode, "error", err)
			return nil, err
		}
		releases = append(releases, func() { sem.Release(1) })
	}

	// Acquire global
	globalAcquireCtx := ctx
	if varsityCode == "spbsu" {
		globalAcquireCtx = context.Background()
	}
	if err := GlobalHTTPSemaphore.Acquire(globalAcquireCtx, 1); err != nil {
		slog.Error("Failed to acquire global semaphore", "error", err)
		// Release any acquired per-varsity
		for _, rel := range releases {
			rel()
		}
		return nil, err
	}
	releases = append(releases, func() { GlobalHTTPSemaphore.Release(1) })

	// Return composed release (reverse order: global first, then per-varsity)
	return func() {
		for i := len(releases) - 1; i >= 0; i-- {
			releases[i]()
		}
	}, nil
}

// Retry executes the operation with exponential backoff retries.
func Retry(operation func() error, maxAttempts int, backoff func(attempt int) time.Duration) error {
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := operation()
		if err == nil {
			return nil
		}
		if attempt < maxAttempts {
			sleepDuration := backoff(attempt)
			// Enhanced logging with error type classification
			errorType := "unknown"
			if isTimeoutError(err) {
				errorType = "timeout"
			} else if isNetworkError(err) {
				errorType = "network"
			}
			slog.Warn("Retry attempt failed, retrying after backoff", 
				"attempt", attempt, 
				"max_attempts", maxAttempts,
				"error", err, 
				"error_type", errorType,
				"backoff", sleepDuration)
			time.Sleep(sleepDuration)
		} else {
			slog.Error("Operation failed after max retries", "attempts", maxAttempts, "error", err)
			return err
		}
	}
	return nil // Unreachable, but for completeness
}

// isTimeoutError checks if the error is a timeout-related error
func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	errorStr := err.Error()
	return strings.Contains(errorStr, "timeout") || 
		   strings.Contains(errorStr, "i/o timeout") ||
		   strings.Contains(errorStr, "context deadline exceeded")
}

// isNetworkError checks if the error is a network-related error
func isNetworkError(err error) bool {
	if err == nil {
		return false
	}
	// Check for net.Error interface
	if netErr, ok := err.(net.Error); ok {
		return netErr.Timeout() || netErr.Temporary()
	}
	errorStr := err.Error()
	return strings.Contains(errorStr, "connection refused") ||
		   strings.Contains(errorStr, "no such host") ||
		   strings.Contains(errorStr, "network unreachable")
}
