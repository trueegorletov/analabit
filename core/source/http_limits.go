// Package source provides data source implementations and HTTP request limiting for analabit.
package source

import (
	"log/slog"
	"os"
	"strconv"

	"golang.org/x/sync/semaphore"
)

// GlobalHTTPSemaphore is the global semaphore for limiting concurrent HTTP requests.
var GlobalHTTPSemaphore *semaphore.Weighted

// VarsitySemaphores maps varsity codes to their per-varsity semaphores.
var VarsitySemaphores map[string]*semaphore.Weighted

var defaultLimits = map[string]int64{
	"hse":    6,
	"itmo":   10,
	"mipt":   1,
	"mirea":  4, // Increased limit for session-based requests (was 3 for sessionless)
	"oldhse": 1,
	"spbstu": 6,
	"spbsu":  6,
	"rzgmu":  3,
	"fmsmu":  16,
	"rsmu":   6,
}

// Session-based limits for sources that use FlareSolverr sessions
var sessionBasedLimits = map[string]int64{
	"mirea": 4, // Higher limit for session-based requests
}

// Fallback limits for sessionless requests
var sessionlessLimits = map[string]int64{
	"mirea": 2, // Lower limit for sessionless fallback requests
}

var envVars = map[string]string{
	"hse":    "HSE_HTTP_MAX_CONCURRENT",
	"itmo":   "ITMO_HTTP_MAX_CONCURRENT",
	"mipt":   "MIPT_HTTP_MAX_CONCURRENT",
	"mirea":  "MIREA_HTTP_MAX_CONCURRENT",
	"oldhse": "OLDHSE_HTTP_MAX_CONCURRENT",
	"spbstu": "SPBSTU_HTTP_MAX_CONCURRENT",
	"spbsu":  "SPBSU_HTTP_MAX_CONCURRENT",
	"rzgmu":  "RZGMU_HTTP_MAX_CONCURRENT",
	"fmsmu":  "FMSMU_HTTP_MAX_CONCURRENT",
	"rsmu":   "RSMU_HTTP_MAX_CONCURRENT",
}

// Environment variables for FlareSolverr session management
var sessionEnvVars = map[string]string{
	"pool_size":                "FLARESOLVERR_SESSION_POOL_SIZE",
	"idle_timeout":             "FLARESOLVERR_SESSION_IDLE_TIMEOUT_MINUTES",
	"max_requests_per_session": "FLARESOLVERR_SESSION_MAX_REQUESTS",
	"health_check_interval":    "FLARESOLVERR_SESSION_HEALTH_CHECK_INTERVAL_MINUTES",
}

const fallbackGlobalLimit = 48

func init() {
	// Load global limit
	globalEnv := os.Getenv("GLOBAL_HTTP_MAX_CONCURRENT")
	globalLimit := int64(fallbackGlobalLimit)
	if globalEnv != "" {
		if parsed, err := strconv.ParseInt(globalEnv, 10, 64); err == nil && parsed > 0 {
			globalLimit = parsed
		} else {
			slog.Warn("Invalid GLOBAL_HTTP_MAX_CONCURRENT, using default", "default", globalLimit)
		}
	}
	GlobalHTTPSemaphore = semaphore.NewWeighted(globalLimit)
	slog.Info("Loaded global HTTP limit", "limit", globalLimit)

	// Load per-varsity limits
	VarsitySemaphores = make(map[string]*semaphore.Weighted)
	for code, envName := range envVars {
		limit := defaultLimits[code]
		envVal := os.Getenv(envName)
		if envVal != "" {
			if parsed, err := strconv.ParseInt(envVal, 10, 64); err == nil && parsed > 0 {
				limit = parsed
			} else {
				slog.Warn("Invalid env for varsity, using default", "varsity", code, "env", envName, "default", limit)
			}
		}
		if limit > 0 {
			VarsitySemaphores[code] = semaphore.NewWeighted(limit)
			// Special logging for MIREA to explain the low limit
			if code == "mirea" {
				slog.Info("Loaded varsity HTTP limit (low due to FlareSolverr browser instances)", "varsity", code, "limit", limit)
			} else {
				slog.Info("Loaded varsity HTTP limit", "varsity", code, "limit", limit)
			}
		} else {
			slog.Info("No per-varsity limit for", "varsity", code)
		}
	}
}
