// Package source provides HTTP timeout coordination and configuration for analabit.
package source

import (
	"context"
	"log/slog"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

// TimeoutConfig holds timeout configuration for a varsity
type TimeoutConfig struct {
	Enabled        bool          // whether timeouts are enabled for this varsity
	BatchSize      int           // number of requests per batch before main timeout
	MicroTimeout   time.Duration // timeout between individual requests
	MainTimeout    time.Duration // timeout between batches
	MicroRatioLow  float64       // randomization range low for micro timeouts
	MicroRatioHigh float64       // randomization range high for micro timeouts
	MainRatioLow   float64       // randomization range low for main timeouts
	MainRatioHigh  float64       // randomization range high for main timeouts
}

// timeoutState tracks request counting and batch state for a varsity
type timeoutState struct {
	mu           sync.Mutex // per-varsity synchronization
	requestCount int        // current number of requests in the current batch
	lastRequest  time.Time  // timestamp of the last request
}

// Global timeout configuration and state
var (
	VarsityTimeouts      map[string]TimeoutConfig
	varsityTimeoutStates sync.Map // map[string]*timeoutState with lock-free access
)

// Default timeout configurations for supported varsities
var defaultTimeoutConfigs = map[string]TimeoutConfig{
	"fmsmu": {
		Enabled:        true,
		BatchSize:      60,
		MicroTimeout:   40 * time.Millisecond,
		MainTimeout:    4500 * time.Millisecond,
		MicroRatioLow:  0.8,
		MicroRatioHigh: 1.2,
		MainRatioLow:   0.75,
		MainRatioHigh:  1.5,
	},
	"mipt": {
		Enabled:        true,
		BatchSize:      8,
		MicroTimeout:   50 * time.Millisecond,
		MainTimeout:    1500 * time.Millisecond,
		MicroRatioLow:  0.2,
		MicroRatioHigh: 1.8,
		MainRatioLow:   0.8,
		MainRatioHigh:  1.2,
	},
	"mephi": {
		Enabled:        true,
		BatchSize:      3,
		MicroTimeout:   200 * time.Millisecond,
		MainTimeout:    3311 * time.Millisecond,
		MicroRatioLow:  0.2,
		MicroRatioHigh: 1.8,
		MainRatioLow:   0.5,
		MainRatioHigh:  1.2,
	},
	// SPbSU uses both timeout coordination and dedicated rate limiting:
	// - Timeout system provides baseline request spacing for politeness
	// - Rate limiter handles 429 errors with exponential backoff (20s, 30s, 60s)
	// - Rate limiter provides additional protection beyond these timeout values
	"spbsu": {
		Enabled:        true,
		BatchSize:      8,                       // Larger batches reduce main timeout frequency, decreasing log noise
		MicroTimeout:   100 * time.Millisecond,  // Reduced since rate limiter handles 429 errors
		MainTimeout:    3333 * time.Millisecond, // Rate limiter provides additional protection
		MicroRatioLow:  0.8,
		MicroRatioHigh: 1.2,
		MainRatioLow:   0.8,
		MainRatioHigh:  1.2,
	},
}

func init() {
	loadTimeoutConfigs()
}

// loadTimeoutConfigs loads timeout configurations from environment variables
func loadTimeoutConfigs() {
	VarsityTimeouts = make(map[string]TimeoutConfig)

	for varsity, defaultConfig := range defaultTimeoutConfigs {
		config := defaultConfig

		// Load enabled flag
		if envVal := os.Getenv(varsity + "_HTTP_TIMEOUT_ENABLED"); envVal != "" {
			if enabled, err := strconv.ParseBool(envVal); err == nil {
				config.Enabled = enabled
			} else {
				slog.Warn("Invalid timeout enabled flag, using default", "varsity", varsity, "env", envVal, "default", config.Enabled)
			}
		}

		// Load batch size
		if envVal := os.Getenv(varsity + "_HTTP_TIMEOUT_BATCH_SIZE"); envVal != "" {
			if batchSize, err := strconv.Atoi(envVal); err == nil && batchSize > 0 {
				config.BatchSize = batchSize
			} else {
				slog.Warn("Invalid timeout batch size, using default", "varsity", varsity, "env", envVal, "default", config.BatchSize)
			}
		}

		// Load micro timeout (in milliseconds)
		if envVal := os.Getenv(varsity + "_HTTP_TIMEOUT_MICRO_MS"); envVal != "" {
			if microMs, err := strconv.Atoi(envVal); err == nil && microMs > 0 {
				config.MicroTimeout = time.Duration(microMs) * time.Millisecond
			} else {
				slog.Warn("Invalid micro timeout, using default", "varsity", varsity, "env", envVal, "default", config.MicroTimeout)
			}
		}

		// Load main timeout (in milliseconds)
		if envVal := os.Getenv(varsity + "_HTTP_TIMEOUT_MAIN_MS"); envVal != "" {
			if mainMs, err := strconv.Atoi(envVal); err == nil && mainMs > 0 {
				config.MainTimeout = time.Duration(mainMs) * time.Millisecond
			} else {
				slog.Warn("Invalid main timeout, using default", "varsity", varsity, "env", envVal, "default", config.MainTimeout)
			}
		}

		// Load micro ratio low
		if envVal := os.Getenv(varsity + "_HTTP_TIMEOUT_MICRO_RATIO_LOW"); envVal != "" {
			if ratio, err := strconv.ParseFloat(envVal, 64); err == nil && ratio > 0 {
				config.MicroRatioLow = ratio
			} else {
				slog.Warn("Invalid micro ratio low, using default", "varsity", varsity, "env", envVal, "default", config.MicroRatioLow)
			}
		}

		// Load micro ratio high
		if envVal := os.Getenv(varsity + "_HTTP_TIMEOUT_MICRO_RATIO_HIGH"); envVal != "" {
			if ratio, err := strconv.ParseFloat(envVal, 64); err == nil && ratio > 0 {
				config.MicroRatioHigh = ratio
			} else {
				slog.Warn("Invalid micro ratio high, using default", "varsity", varsity, "env", envVal, "default", config.MicroRatioHigh)
			}
		}

		// Load main ratio low
		if envVal := os.Getenv(varsity + "_HTTP_TIMEOUT_MAIN_RATIO_LOW"); envVal != "" {
			if ratio, err := strconv.ParseFloat(envVal, 64); err == nil && ratio > 0 {
				config.MainRatioLow = ratio
			} else {
				slog.Warn("Invalid main ratio low, using default", "varsity", varsity, "env", envVal, "default", config.MainRatioLow)
			}
		}

		// Load main ratio high
		if envVal := os.Getenv(varsity + "_HTTP_TIMEOUT_MAIN_RATIO_HIGH"); envVal != "" {
			if ratio, err := strconv.ParseFloat(envVal, 64); err == nil && ratio > 0 {
				config.MainRatioHigh = ratio
			} else {
				slog.Warn("Invalid main ratio high, using default", "varsity", varsity, "env", envVal, "default", config.MainRatioHigh)
			}
		}

		VarsityTimeouts[varsity] = config

		if config.Enabled {
			slog.Info("Loaded timeout configuration",
				"varsity", varsity,
				"enabled", config.Enabled,
				"batch_size", config.BatchSize,
				"micro_timeout", config.MicroTimeout,
				"main_timeout", config.MainTimeout,
				"micro_ratio_range", []float64{config.MicroRatioLow, config.MicroRatioHigh},
				"main_ratio_range", []float64{config.MainRatioLow, config.MainRatioHigh})
		} else {
			slog.Debug("Timeout configuration disabled", "varsity", varsity)
		}
	}
}

// WaitBeforeHTTPRequest implements batch-based timeout coordination for HTTP requests
func WaitBeforeHTTPRequest(varsityCode string, ctx context.Context) error {
	// Check if timeout configuration exists for this varsity
	config, exists := VarsityTimeouts[varsityCode]
	if !exists || !config.Enabled {
		// No timeout configuration or disabled - proceed immediately
		return nil
	}

	// Record synchronization start time for time accounting
	syncStartTime := time.Now()

	// Get or create state for this varsity using sync.Map
	stateInterface, _ := varsityTimeoutStates.LoadOrStore(varsityCode, &timeoutState{
		requestCount: 0,
		lastRequest:  time.Time{},
	})
	state := stateInterface.(*timeoutState)

	// Acquire per-varsity mutex
	state.mu.Lock()
	defer state.mu.Unlock()

	// Calculate synchronization overhead
	syncEndTime := time.Now()
	syncOverhead := syncEndTime.Sub(syncStartTime)

	now := syncEndTime

	// Determine if we need a timeout and which type
	var timeoutDuration time.Duration
	var timeoutType string

	if state.requestCount >= config.BatchSize {
		// We've reached the batch limit - apply main timeout and reset counter
		timeoutDuration = calculateRandomizedTimeout(config.MainTimeout, config.MainRatioLow, config.MainRatioHigh)
		timeoutType = "main"
		state.requestCount = 0
		slog.Debug("Applying main timeout after batch completion",
			"varsity", varsityCode,
			"batch_size", config.BatchSize,
			"timeout", timeoutDuration)
	} else if !state.lastRequest.IsZero() {
		// Not the first request and not at batch limit - apply micro timeout
		timeoutDuration = calculateRandomizedTimeout(config.MicroTimeout, config.MicroRatioLow, config.MicroRatioHigh)
		timeoutType = "micro"
		slog.Debug("Applying micro timeout between requests",
			"varsity", varsityCode,
			"request_count", state.requestCount+1,
			"timeout", timeoutDuration)
	}

	// Update state
	state.requestCount++
	state.lastRequest = now

	// Apply timeout if needed, accounting for synchronization overhead
	if timeoutDuration > 0 {
		// Subtract synchronization overhead from intended timeout duration
		adjustedTimeout := timeoutDuration - syncOverhead

		// Skip sleep if synchronization time already exceeded intended timeout
		if adjustedTimeout <= 0 {
			slog.Debug("Skipping sleep due to synchronization overhead",
				"varsity", varsityCode,
				"timeout_type", timeoutType,
				"intended_duration", timeoutDuration,
				"sync_overhead", syncOverhead,
				"request_count", state.requestCount)
			return nil
		}

		slog.Debug("Waiting before HTTP request",
			"varsity", varsityCode,
			"timeout_type", timeoutType,
			"intended_duration", timeoutDuration,
			"adjusted_duration", adjustedTimeout,
			"sync_overhead", syncOverhead,
			"request_count", state.requestCount)

		// Context-aware sleep with cancellation support
		select {
		case <-time.After(adjustedTimeout):
			// Timeout completed normally
			return nil
		case <-ctx.Done():
			// Context was cancelled
			slog.Debug("HTTP timeout cancelled by context",
				"varsity", varsityCode,
				"timeout_type", timeoutType,
				"error", ctx.Err())
			return ctx.Err()
		}
	}

	return nil
}

// calculateRandomizedTimeout applies randomization to a base timeout duration
func calculateRandomizedTimeout(baseDuration time.Duration, ratioLow, ratioHigh float64) time.Duration {
	if ratioLow >= ratioHigh {
		// Invalid range - return base duration
		return baseDuration
	}

	// Generate random ratio within the specified range
	ratio := ratioLow + rand.Float64()*(ratioHigh-ratioLow)

	// Apply ratio to base duration
	randomizedDuration := time.Duration(float64(baseDuration) * ratio)

	return randomizedDuration
}

// resetBatchCounter resets the batch counter for a varsity (useful for testing or manual reset)
func resetBatchCounter(varsityCode string) {
	if stateInterface, ok := varsityTimeoutStates.Load(varsityCode); ok {
		state := stateInterface.(*timeoutState)
		state.mu.Lock()
		defer state.mu.Unlock()

		state.requestCount = 0
		slog.Debug("Reset batch counter", "varsity", varsityCode)
	}
}

// GetTimeoutState returns the current timeout state for a varsity (useful for monitoring/debugging)
func GetTimeoutState(varsityCode string) (requestCount int, lastRequest time.Time, exists bool) {
	if stateInterface, ok := varsityTimeoutStates.Load(varsityCode); ok {
		state := stateInterface.(*timeoutState)
		state.mu.Lock()
		defer state.mu.Unlock()

		return state.requestCount, state.lastRequest, true
	}
	return 0, time.Time{}, false
}
