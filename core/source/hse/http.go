package hse

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/utils"

	"github.com/xuri/excelize/v2"
)

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

	release, err := source.AcquireHTTPSemaphores(ctx, "hse")
	if err != nil {
		return fmt.Errorf("failed to acquire semaphores for HSE list from %s: %w", s.URL, err)
	}
	defer release()

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
