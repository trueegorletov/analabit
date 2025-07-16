package mirea

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/source/flaresolverr"
	"github.com/trueegorletov/analabit/core/utils"
)

// Configuration constants for MIREA API
const (
	MireaDefaultDelay = 33 * time.Millisecond // Default delay after each request
)

// getMireaHeaders returns the headers to be used for MIREA requests
func getMireaHeaders() map[string]string {
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
		"Referer":            "https://priem.mirea.ru/",
	}
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

// fetchMireaListByID fetches and decodes a single MIREA list using FlareSolverr
func fetchMireaListByID(listID string) (*MireaListResponse, error) {
	if listID == "" {
		return nil, nil
	}

	// Build API URL
	apiURL := "https://priem.mirea.ru/competitions_api/entrants?competitions[]=" + listID

	// Create context with timeout to prevent deadlocks
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Acquire semaphore for rate limiting
	release, err := source.AcquireHTTPSemaphores(ctx, "mirea")
	if err != nil {
		return nil, fmt.Errorf("failed to acquire semaphores for list %s: %w", listID, err)
	}
	defer release()

	// Make HTTP request through FlareSolverr with session management
	resp, err := flaresolverr.SafeGetWithDomain(apiURL, getMireaHeaders())
	if err != nil {
		if flaresolverr.IsFlareSolverrError(err) {
			return nil, fmt.Errorf("FlareSolverr unavailable for list %s: %w", listID, err)
		}
		return nil, fmt.Errorf("failed to download list %s via FlareSolverr: %w", listID, err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to download list %s (status code %d)", listID, resp.StatusCode)
	}

	// Parse JSON response
	var mireaResp MireaListResponse
	// Extract JSON from HTML wrapper if present
	body := resp.Body
	if start := strings.Index(body, "<pre>"); start != -1 {
		if end := strings.LastIndex(body, "</pre>"); end != -1 && end > start {
			body = body[start+len("<pre>") : end]
		}
	}
	// Fallback: extract JSON by locating the first '{' and last '}'
	if !strings.HasPrefix(strings.TrimSpace(body), "{") {
		if start := strings.Index(body, "{"); start != -1 {
			if end := strings.LastIndex(body, "}"); end != -1 && end > start {
				body = body[start : end+1]
			}
		}
	}
	if err := json.Unmarshal([]byte(body), &mireaResp); err != nil {
		slog.Warn("Watafuk? First 5000 characters of response body", "body", body[:5000])
		return nil, fmt.Errorf("failed to parse JSON response for list %s: %w", listID, err)
	}

	// Apply delay after each request
	time.Sleep(getMireaRequestDelay())

	return &mireaResp, nil
}

// fetchOrGetCachedMireaList fetches a MIREA list or returns it from cache
func fetchOrGetCachedMireaList(listID string, cache map[string]*MireaListResponse) (*MireaListResponse, error) {
	if listID == "" {
		return nil, nil
	}

	// Check cache first
	if cached, exists := cache[listID]; exists {
		log.Printf("MIREA: Cache hit for list %s", listID)
		return cached, nil
	}

	// Not in cache, fetch from API
	log.Printf("MIREA: Fetching list %s from API", listID)
	resp, err := fetchMireaListByID(listID)
	if err != nil {
		slog.Warn("Watafuk? Error fetching MIREA list", "error", err)
		return nil, err
	}

	// Store in cache for future use
	cache[listID] = resp
	return resp, nil
}

// LoadTo implements source.HeadingSource for HTTPHeadingSource.
func (s *HTTPHeadingSource) LoadTo(receiver source.DataReceiver) error {
	log.Printf("MIREA: Starting data load")

	// Initialize response cache to avoid duplicate requests
	listResponseCache := make(map[string]*MireaListResponse)

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
				resp, err := fetchOrGetCachedMireaList(listID, listResponseCache)
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
		resp, err := fetchOrGetCachedMireaList(s.RegularListIDs[0], listResponseCache)
		if err == nil && resp != nil && len(resp.Data) > 0 {
			regularCapacity = resp.Data[0].Plan
		}
	}

	// TargetQuota: sum capacities from all lists
	for _, listID := range s.TargetQuotaListIDs {
		if listID != "" {
			resp, err := fetchOrGetCachedMireaList(listID, listResponseCache)
			if err == nil && resp != nil && len(resp.Data) > 0 {
				targetQuotaCapacity += resp.Data[0].Plan
			}
		}
	}

	// DedicatedQuota: sum capacities
	for _, listID := range s.DedicatedQuotaListIDs {
		if listID != "" {
			resp, err := fetchOrGetCachedMireaList(listID, listResponseCache)
			if err == nil && resp != nil && len(resp.Data) > 0 {
				dedicatedQuotaCapacity += resp.Data[0].Plan
			}
		}
	}

	// SpecialQuota: sum capacities
	for _, listID := range s.SpecialQuotaListIDs {
		if listID != "" {
			resp, err := fetchOrGetCachedMireaList(listID, listResponseCache)
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

	log.Printf("MIREA: Sent heading data for %s",
		prettyName)

	// Process all competition lists
	// Regular Lists
	for _, listID := range s.RegularListIDs {
		if listID == "" {
			continue
		}
		resp, err := fetchOrGetCachedMireaList(listID, listResponseCache)
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
		resp, err := fetchOrGetCachedMireaList(listID, listResponseCache)
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
		resp, err := fetchOrGetCachedMireaList(listID, listResponseCache)
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
		resp, err := fetchOrGetCachedMireaList(listID, listResponseCache)
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
		resp, err := fetchOrGetCachedMireaList(listID, listResponseCache)
		if err != nil {
			continue // Skip failed requests
		}
		if resp == nil || len(resp.Data) == 0 {
			continue
		}
		parseAndLoadApplications(resp.Data[0].Entrants, core.CompetitionSpecialQuota, headingCode, receiver)
	}

	// Add random delay at the end after all sources are loaded (0.05s to 0.15s)
	randomDelay := time.Duration(50+rand.Intn(100)) * time.Millisecond
	time.Sleep(randomDelay)

	log.Printf("MIREA: Data load completed for %s", prettyName)
	return nil
}
