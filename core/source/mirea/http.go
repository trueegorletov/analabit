package mirea

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/utils"
	"resty.dev/v3"
)

// Configuration constants for MIREA API
const (
	MireaDefaultDelay = 2 * time.Second // Default delay after each request
)

// Global Resty client with cookie persistence and browser-like behavior
var mireaClient *resty.Client

// initMireaClient initializes the global Resty client with realistic browser settings
func initMireaClient() error {
	if mireaClient != nil {
		return nil // Already initialized
	}

	// Create Resty client with realistic browser settings
	mireaClient = resty.New().
		SetTimeout(30*time.Second).
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36").
		SetHeader("Accept", "application/json, text/plain, */*").
		SetHeader("Accept-Language", "ru-RU,ru;q=0.9,en;q=0.8").
		SetHeader("Accept-Encoding", "gzip, deflate, br").
		SetHeader("Cache-Control", "no-cache").
		SetHeader("Pragma", "no-cache").
		SetHeader("Sec-Ch-Ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`).
		SetHeader("Sec-Ch-Ua-Mobile", "?0").
		SetHeader("Sec-Ch-Ua-Platform", `"Windows"`).
		SetHeader("Sec-Fetch-Dest", "empty").
		SetHeader("Sec-Fetch-Mode", "cors").
		SetHeader("Sec-Fetch-Site", "same-origin").
		SetHeader("Referer", "https://priem.mirea.ru/")

	return nil
}

// getMireaRequestDelay returns the configured delay for MIREA requests
func getMireaRequestDelay() time.Duration {
	if envDelay := os.Getenv("MIREA_REQUEST_DELAY_SECONDS"); envDelay != "" {
		if seconds, err := strconv.Atoi(envDelay); err == nil && seconds > 0 {
			return time.Duration(seconds) * time.Second
		}
	}
	return MireaDefaultDelay
}

// HTTPHeadingSource loads MIREA heading data from JSON list IDs.
type HTTPHeadingSource struct {
	RegularListIDs        []string
	BVIListIDs            []string
	TargetQuotaListIDs    []string
	DedicatedQuotaListIDs []string
	SpecialQuotaListIDs   []string
}

// fetchMireaListByID fetches and decodes a single MIREA list using go-resty
func fetchMireaListByID(listID string) (*MireaListResponse, error) {
	if listID == "" {
		return nil, nil
	}

	// Initialize client if needed
	if err := initMireaClient(); err != nil {
		return nil, fmt.Errorf("failed to initialize MIREA client: %w", err)
	}

	// Build API URL
	apiURL := "https://priem.mirea.ru/competitions_api/entrants?competitions[]=" + listID

	// Create context with timeout to prevent deadlocks
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Acquire semaphore for rate limiting
	release, err := source.AcquireHTTPSemaphores(ctx, "mirea")
	if err != nil {
		return nil, fmt.Errorf("failed to acquire semaphores for list %s: %w", listID, err)
	}
	defer release()

	// Make HTTP request with go-resty
	var mireaResp MireaListResponse
	resp, err := mireaClient.R().
		SetResult(&mireaResp).
		Get(apiURL)

	if err != nil {
		return nil, fmt.Errorf("failed to download list %s: %w", listID, err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to download list %s (status code %d)", listID, resp.StatusCode())
	}

	// Apply delay after each request
	time.Sleep(getMireaRequestDelay())

	return &mireaResp, nil
}

// LoadTo implements source.HeadingSource for HTTPHeadingSource.
func (s *HTTPHeadingSource) LoadTo(receiver source.DataReceiver) error {
	log.Printf("MIREA: Starting data load")

	// Extract metadata from the first available list
	var prettyName string
	var regularCapacity, targetQuotaCapacity, dedicatedQuotaCapacity, specialQuotaCapacity int

	// Try to get metadata from the first available list
	allListIDGroups := [][]string{
		s.RegularListIDs,
		s.BVIListIDs,
		s.TargetQuotaListIDs,
		s.DedicatedQuotaListIDs,
		s.SpecialQuotaListIDs,
	}

	for _, listGroup := range allListIDGroups {
		for _, listID := range listGroup {
			if listID != "" {
				resp, err := fetchMireaListByID(listID)
				if err != nil {
					continue // Try next list
				}
				if resp != nil && len(resp.Data) > 0 {
					prettyName = extractHeadingName(resp.Data[0].Title)
					break
				}
			}
		}
		if prettyName != "" {
			break
		}
	}

	if prettyName == "" {
		return fmt.Errorf("could not extract heading name from any available list")
	}

	// Calculate capacities
	// Regular/BVI: use the same capacity value
	if len(s.RegularListIDs) > 0 && s.RegularListIDs[0] != "" {
		resp, err := fetchMireaListByID(s.RegularListIDs[0])
		if err == nil && resp != nil && len(resp.Data) > 0 {
			regularCapacity = resp.Data[0].Plan
		}
	}

	// TargetQuota: sum capacities from all lists
	for _, listID := range s.TargetQuotaListIDs {
		if listID != "" {
			resp, err := fetchMireaListByID(listID)
			if err == nil && resp != nil && len(resp.Data) > 0 {
				targetQuotaCapacity += resp.Data[0].Plan
			}
		}
	}

	// DedicatedQuota: sum capacities
	for _, listID := range s.DedicatedQuotaListIDs {
		if listID != "" {
			resp, err := fetchMireaListByID(listID)
			if err == nil && resp != nil && len(resp.Data) > 0 {
				dedicatedQuotaCapacity += resp.Data[0].Plan
			}
		}
	}

	// SpecialQuota: sum capacities
	for _, listID := range s.SpecialQuotaListIDs {
		if listID != "" {
			resp, err := fetchMireaListByID(listID)
			if err == nil && resp != nil && len(resp.Data) > 0 {
				specialQuotaCapacity += resp.Data[0].Plan
			}
		}
	}

	headingCode := utils.GenerateHeadingCode(prettyName)

	// Send heading data
	receiver.PutHeadingData(&source.HeadingData{
		Code: headingCode,
		Capacities: core.Capacities{
			Regular:        regularCapacity,
			TargetQuota:    targetQuotaCapacity,
			DedicatedQuota: dedicatedQuotaCapacity,
			SpecialQuota:   specialQuotaCapacity,
		},
		PrettyName: prettyName,
	})

	// Process all competition lists
	// Regular Lists
	for _, listID := range s.RegularListIDs {
		if listID == "" {
			continue
		}
		resp, err := fetchMireaListByID(listID)
		if err != nil {
			continue // Skip failed requests
		}
		if resp == nil || len(resp.Data) == 0 {
			continue
		}
		parseAndLoadApplications(resp.Data[0].Entrants, core.CompetitionRegular, headingCode, receiver)
	}

	// BVI Lists
	for _, listID := range s.BVIListIDs {
		if listID == "" {
			continue
		}
		resp, err := fetchMireaListByID(listID)
		if err != nil {
			continue // Skip failed requests
		}
		if resp == nil || len(resp.Data) == 0 {
			continue
		}
		parseAndLoadApplications(resp.Data[0].Entrants, core.CompetitionBVI, headingCode, receiver)
	}

	// TargetQuota Lists
	for _, listID := range s.TargetQuotaListIDs {
		if listID == "" {
			continue
		}
		resp, err := fetchMireaListByID(listID)
		if err != nil {
			continue // Skip failed requests
		}
		if resp == nil || len(resp.Data) == 0 {
			continue
		}
		parseAndLoadApplications(resp.Data[0].Entrants, core.CompetitionTargetQuota, headingCode, receiver)
	}

	// DedicatedQuota Lists
	for _, listID := range s.DedicatedQuotaListIDs {
		if listID == "" {
			continue
		}
		resp, err := fetchMireaListByID(listID)
		if err != nil {
			continue // Skip failed requests
		}
		if resp == nil || len(resp.Data) == 0 {
			continue
		}
		parseAndLoadApplications(resp.Data[0].Entrants, core.CompetitionDedicatedQuota, headingCode, receiver)
	}

	// SpecialQuota Lists
	for _, listID := range s.SpecialQuotaListIDs {
		if listID == "" {
			continue
		}
		resp, err := fetchMireaListByID(listID)
		if err != nil {
			continue // Skip failed requests
		}
		if resp == nil || len(resp.Data) == 0 {
			continue
		}
		parseAndLoadApplications(resp.Data[0].Entrants, core.CompetitionSpecialQuota, headingCode, receiver)
	}

	log.Printf("MIREA: Data load completed for %s", prettyName)
	return nil
}
