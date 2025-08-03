package spbstu

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/source/flaresolverr"
	"github.com/trueegorletov/analabit/core/utils"
)

// HTTPHeadingSource loads SPbSTU heading data from multiple JSON list IDs.
type HTTPHeadingSource struct {
	PrettyName           string
	RegularListID        int
	TargetQuotaListIDs   []int
	DedicatedQuotaListID int
	SpecialQuotaListID   int
	Capacities           core.Capacities
}

// SpbstuCapacityResponse represents the capacity response from SPbSTU API.
type SpbstuCapacityResponse struct {
	Places int `json:"places"`
}

// getSpbstuHeaders returns the headers to be used for SPbSTU requests
func getSpbstuHeaders() map[string]string {
	return map[string]string{
		"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Accept-Language":    "ru-RU,ru;q=0.9,en;q=0.8",
		"Accept-Encoding":    "gzip, deflate, br",
		"Cache-Control":      "no-cache",
		"Pragma":             "no-cache",
		"Sec-Ch-Ua":          `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`,
		"Sec-Ch-Ua-Mobile":   "?0",
		"Sec-Ch-Ua-Platform": `"Windows"`,
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "same-origin",
		"Referer":            "https://my.spbstu.ru/",
	}
}

// getSpbstuRequestDelay returns the configured delay for SPbSTU requests
func getSpbstuRequestDelay() time.Duration {
	if envDelay := os.Getenv("SPBSTU_REQUEST_DELAY_SECONDS"); envDelay != "" {
		if seconds, err := strconv.Atoi(envDelay); err == nil && seconds > 0 {
			return time.Duration(seconds) * time.Second
		}
	}
	return 100 * time.Millisecond // Default delay
}

// fetchSpbstuListByID fetches and decodes the SPbSTU list from a list ID using FlareSolverr
func fetchSpbstuListByID(listID int, competitionFilter int) ([]SpbstuApplicationEntry, error) {
	if listID == -1 {
		return nil, nil
	}
	url := fmt.Sprintf("https://my.spbstu.ru/home/get-abit-list?filter_1=2&filter_2=%d&filter_3=%d&education_level=bachelor", competitionFilter, listID)

	// Create context with timeout to prevent deadlocks
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Acquire semaphore for rate limiting
	release, err := source.AcquireHTTPSemaphores(ctx, "spbstu")
	if err != nil {
		return nil, fmt.Errorf("failed to acquire semaphores for %s: %w", url, err)
	}
	defer release()

	// Make HTTP request through FlareSolverr with session management
	resp, err := flaresolverr.SafeGetWithDomain(url, getSpbstuHeaders())
	if err != nil {
		if flaresolverr.IsFlareSolverrError(err) {
			return nil, fmt.Errorf("FlareSolverr unavailable for %s: %w", url, err)
		}
		return nil, fmt.Errorf("failed to download %s via FlareSolverr: %w", url, err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to download %s (status code %d)", url, resp.StatusCode)
	}

	// Extract JSON from HTML wrapper if present
	body := resp.Body
	if start := strings.Index(body, "<pre>"); start != -1 {
		if end := strings.LastIndex(body, "</pre>"); end != -1 && end > start {
			body = body[start+len("<pre>") : end]
		}
	}
	// Fallback: extract JSON by locating the first '{' and last '}'
	if !strings.HasPrefix(strings.TrimSpace(body), "{") && !strings.HasPrefix(strings.TrimSpace(body), "[") {
		if start := strings.Index(body, "{"); start != -1 {
			if end := strings.LastIndex(body, "}"); end != -1 && end > start {
				body = body[start : end+1]
			}
		} else if start := strings.Index(body, "["); start != -1 {
			if end := strings.LastIndex(body, "]"); end != -1 && end > start {
				body = body[start : end+1]
			}
		}
	}

	// Parse JSON response using existing decode function
	entries, err := decodeSpbstuList(strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response for %s: %w", url, err)
	}

	// Apply delay after each request
	time.Sleep(getSpbstuRequestDelay())

	return entries, nil
}

// fetchSpbstuCapacity fetches capacity data for a specific list ID using FlareSolverr
func fetchSpbstuCapacity(listID int) (int, error) {
	if listID == -1 {
		return 0, nil
	}

	url := "https://my.spbstu.ru/home/get-direction-info"
	postData := map[string]interface{}{
		"condition":       "1",
		"id_3":            strconv.Itoa(listID),
		"education_level": "bachelor",
	}

	// Create context with timeout to prevent deadlocks
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Acquire semaphore for rate limiting
	release, err := source.AcquireHTTPSemaphores(ctx, "spbstu")
	if err != nil {
		return 0, fmt.Errorf("failed to acquire semaphores for %s: %w", url, err)
	}
	defer release()

	// Set headers for POST request
	headers := getSpbstuHeaders()
	headers["Content-Type"] = "application/json"

	// Make HTTP POST request through FlareSolverr
	resp, err := flaresolverr.SafePostWithData(url, postData, headers)
	if err != nil {
		if flaresolverr.IsFlareSolverrError(err) {
			return 0, fmt.Errorf("FlareSolverr unavailable for capacity request %s: %w", url, err)
		}
		return 0, fmt.Errorf("failed to fetch capacity from %s via FlareSolverr: %w", url, err)
	}

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("failed to fetch capacity from %s (status code %d)", url, resp.StatusCode)
	}

	// Extract JSON from HTML wrapper if present
	body := resp.Body
	if start := strings.Index(body, "<pre>"); start != -1 {
		if end := strings.LastIndex(body, "</pre>"); end != -1 && end > start {
			body = body[start+len("<pre>") : end]
		}
	}
	// Fallback: extract JSON by locating the first '[' and last ']'
	if !strings.HasPrefix(strings.TrimSpace(body), "[") {
		if start := strings.Index(body, "["); start != -1 {
			if end := strings.LastIndex(body, "]"); end != -1 && end > start {
				body = body[start : end+1]
			}
		}
	}

	var responses []SpbstuCapacityResponse
	if err := json.Unmarshal([]byte(body), &responses); err != nil {
		return 0, fmt.Errorf("failed to decode capacity response: %w", err)
	}

	if len(responses) == 0 {
		return 0, fmt.Errorf("empty capacity response")
	}

	// Apply delay after each request
	time.Sleep(getSpbstuRequestDelay())

	return responses[0].Places, nil
}

// LoadTo implements source.HeadingSource for HTTPHeadingSource.
func (s *HTTPHeadingSource) LoadTo(receiver source.DataReceiver) error {
	if s.PrettyName == "" {
		return fmt.Errorf("PrettyName is required for SPbSTU HTTPHeadingSource")
	}

	headingCode := utils.GenerateHeadingCode(s.PrettyName)

	// Fetch capacities at runtime if not provided
	capacities := s.Capacities
	if capacities.Regular == 0 && s.RegularListID != -1 {
		regular, err := fetchSpbstuCapacity(s.RegularListID)
		if err != nil {
			log.Printf("Error fetching regular capacity for %s: %v", s.PrettyName, err)
		} else {
			capacities.Regular = regular
		}
	}

	if capacities.DedicatedQuota == 0 && s.DedicatedQuotaListID != -1 {
		dedicated, err := fetchSpbstuCapacity(s.DedicatedQuotaListID)
		if err != nil {
			log.Printf("Error fetching dedicated quota capacity for %s: %v", s.PrettyName, err)
		} else {
			capacities.DedicatedQuota = dedicated
		}
	}

	if capacities.SpecialQuota == 0 && s.SpecialQuotaListID != -1 {
		special, err := fetchSpbstuCapacity(s.SpecialQuotaListID)
		if err != nil {
			log.Printf("Error fetching special quota capacity for %s: %v", s.PrettyName, err)
		} else {
			capacities.SpecialQuota = special
		}
	}

	// Sum TargetQuota capacities when multiple lists exist
	if capacities.TargetQuota == 0 && len(s.TargetQuotaListIDs) > 0 {
		totalTarget := 0
		for _, listID := range s.TargetQuotaListIDs {
			if listID == -1 {
				continue
			}
			target, err := fetchSpbstuCapacity(listID)
			if err != nil {
				log.Printf("Error fetching target quota capacity for list %d in %s: %v", listID, s.PrettyName, err)
				continue
			}
			totalTarget += target
		}
		capacities.TargetQuota = totalTarget
	}

	receiver.PutHeadingData(&source.HeadingData{
		Code:       headingCode,
		Capacities: capacities,
		PrettyName: s.PrettyName,
	})

	// Define list types with their competition filters
	listDefs := []struct {
		ListID            int
		Competition       core.Competition
		CompetitionFilter int
		ListName          string
	}{
		{s.RegularListID, core.CompetitionRegular, 1, "Regular&BVI List"},
		{s.DedicatedQuotaListID, core.CompetitionDedicatedQuota, 4, "Dedicated Quota List"},
		{s.SpecialQuotaListID, core.CompetitionSpecialQuota, 3, "Special Quota List"},
	}

	for _, def := range listDefs {
		if def.ListID == -1 {
			continue
		}
		entries, err := fetchSpbstuListByID(def.ListID, def.CompetitionFilter)
		if err != nil {
			log.Printf("Error fetching %s (%d): %v", def.ListName, def.ListID, err)
			continue
		}
		if entries == nil {
			log.Printf("No entries found in %s (%d)", def.ListName, def.ListID)
			continue
		}
		parseAndLoadApplications(entries, def.Competition, headingCode, receiver)
	}

	// Handle multiple target quota list IDs
	for i, listID := range s.TargetQuotaListIDs {
		if listID == -1 {
			continue
		}
		entries, err := fetchSpbstuListByID(listID, 5) // Competition filter 5 for TargetQuota
		if err != nil {
			log.Printf("Error fetching Target Quota List %d (%d): %v", i+1, listID, err)
			continue
		}
		if entries == nil {
			log.Printf("No entries found in Target Quota List %d (%d)", i+1, listID)
			continue
		}
		parseAndLoadApplications(entries, core.CompetitionTargetQuota, headingCode, receiver)
	}

	return nil
}
