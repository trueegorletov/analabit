package source

import (
	"context"
	"log/slog"
	"time"
)

// AcquireHTTPSemaphores acquires from per-varsity (if exists) and global semaphores, returning a release function.
func AcquireHTTPSemaphores(ctx context.Context, varsityCode string) (func(), error) {
	var releases []func()

	// Acquire per-varsity if exists
	if sem, ok := VarsitySemaphores[varsityCode]; ok {
		if err := sem.Acquire(ctx, 1); err != nil {
			slog.Error("Failed to acquire varsity semaphore", "varsity", varsityCode, "error", err)
			return nil, err
		}
		releases = append(releases, func() { sem.Release(1) })
	}

	// Acquire global
	if err := GlobalHTTPSemaphore.Acquire(ctx, 1); err != nil {
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
			slog.Warn("Retry attempt failed, retrying after backoff", "attempt", attempt, "error", err, "backoff", sleepDuration)
			time.Sleep(sleepDuration)
		} else {
			slog.Error("Operation failed after max retries", "attempts", maxAttempts, "error", err)
			return err
		}
	}
	return nil // Unreachable, but for completeness
}
