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
	"mipt":   6,
	"mirea":  3, // Low limit due to FlareSolverr browser instances being memory-intensive
	"oldhse": 1,
	"spbstu": 6,
	"spbsu":  6,
	"rzgmu":  3,
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
}

const fallbackGlobalLimit = 32

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
