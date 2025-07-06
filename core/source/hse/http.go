package hse

import (
	"analabit/core"
	"analabit/core/source"
	"analabit/core/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/xuri/excelize/v2"
	"golang.org/x/sync/semaphore"
)

// Global semaphore to limit concurrent HTTP requests
var httpRequestSemaphore *semaphore.Weighted

func init() {
	// Default to 3 concurrent requests for HSE, but allow override via environment variable
	maxConcurrentRequests := int64(3)
	if envVal := os.Getenv("HSE_HTTP_MAX_CONCURRENT_REQUESTS"); envVal != "" {
		if parsed, err := strconv.ParseInt(envVal, 10, 64); err == nil && parsed > 0 {
			maxConcurrentRequests = parsed
		} else {
			log.Printf("Warning: Invalid HSE_HTTP_MAX_CONCURRENT_REQUESTS value '%s', using default %d", envVal, maxConcurrentRequests)
		}
	}

	httpRequestSemaphore = semaphore.NewWeighted(maxConcurrentRequests)
	log.Printf("Initialized HSE HTTP request semaphore with limit: %d concurrent requests", maxConcurrentRequests)
}

// HTTPHeadingSource defines how to load HSE heading data from a single XLSX file URL.
// The new HSE format contains all applications for one program in a single file.
type HTTPHeadingSource struct {
	URL        string // URL to the XLSX file containing all applications for this heading
	Capacities core.Capacities
}

// LoadTo loads data from HTTP source, sending HeadingData and ApplicationData to the provided receiver.
func (s *HTTPHeadingSource) LoadTo(receiver source.DataReceiver) error {
	if s.URL == "" {
		return fmt.Errorf("URL is required for HSE HttpHeadingSource")
	}

	log.Printf("Downloading HSE admission list from: %s", s.URL)

	// Acquire a semaphore slot, respecting context cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := httpRequestSemaphore.Acquire(ctx, 1); err != nil {
		return fmt.Errorf("failed to acquire semaphore for HSE list from %s: %w", s.URL, err)
	}
	defer httpRequestSemaphore.Release(1)

	resp, err := http.Get(s.URL)
	if err != nil {
		return fmt.Errorf("failed to download HSE list from %s: %w", s.URL, err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("Error closing response body for HSE list from %s: %v", s.URL, closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download HSE list from %s (status code %d)", s.URL, resp.StatusCode)
	}

	// Open the XLSX file from the response body
	f, err := excelize.OpenReader(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to open Excel content for HSE list from %s: %w", s.URL, err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			log.Printf("Error closing Excel file for HSE list from %s: %v", s.URL, closeErr)
		}
	}()

	// Extract heading name from the XLSX
	prettyName, err := extractPrettyNameFromXLSX(f, s.URL)
	if err != nil {
		return fmt.Errorf("failed to extract pretty name from HSE list at %s: %w", s.URL, err)
	}

	headingCode := utils.GenerateHeadingCode(prettyName)

	// Send HeadingData to the receiver
	receiver.PutHeadingData(&source.HeadingData{
		Code:       headingCode,
		Capacities: s.Capacities,
		PrettyName: prettyName,
	})

	log.Printf("Sent HSE heading: %s (Code: %s, Caps: %v)", prettyName, headingCode, s.Capacities)

	// Parse applications from the XLSX and send to receiver
	err = parseApplicationsFromXLSX(f, headingCode, receiver, s.URL)
	if err != nil {
		return fmt.Errorf("failed to parse applications from HSE list at %s: %w", s.URL, err)
	}

	log.Printf("Successfully processed HSE heading from %s", s.URL)
	return nil
}
