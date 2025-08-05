package resolver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/trueegorletov/analabit/core/idresolver"
	"github.com/trueegorletov/analabit/service/idmsu/cache"
	"github.com/trueegorletov/analabit/service/idmsu/matcher"
)

// MSUResolver implements idresolver.StudentIDResolver interface for MSU applicants.
// It maintains a background goroutine to fetch and update data from Gosuslugi API.
type MSUResolver struct {
	cache               *cache.LayeredCache
	dbStore             *cache.DatabaseStore
	fetcherStarted      bool
	lastSuccessfulFetch time.Time
	mu                  sync.RWMutex
	fetchInProgress     bool
}

// NewMSUResolver creates a new MSUResolver with the provided cache and database store.
func NewMSUResolver(cache *cache.LayeredCache, dbStore *cache.DatabaseStore) *MSUResolver {
	resolver := &MSUResolver{
		cache:   cache,
		dbStore: dbStore,
	}

	// Restore last fetch time from database
	var lastFetch time.Time
	if lf, err := dbStore.GetLastSuccessfulFetch(); err == nil && !lf.IsZero() {
		lastFetch = lf
		resolver.lastSuccessfulFetch = lastFetch
		slog.Info("Restored last fetch time from database", "lastFetch", lastFetch)
	}

	// Repair logic for existing cache
	if !lastFetch.IsZero() && time.Since(lastFetch) > 16*time.Hour {
	    repaired, err := dbStore.RepairLastFetch()
	    if err != nil {
	        slog.Error("Failed to repair last fetch", "error", err)
	    } else if repaired {
	        resolver.lastSuccessfulFetch = time.Now()
	        slog.Info("Repaired last successful fetch")
	    }
	} else if lastFetch.IsZero() {
	    created, err := dbStore.CreateRepairRun()
	    if err != nil {
	        slog.Error("Failed to create repair run", "error", err)
	    } else if created {
	        resolver.lastSuccessfulFetch = time.Now()
	        slog.Info("Created repair fetch run")
	    }
	}

	return resolver
}

// ResolveBatch implements idresolver.StudentIDResolver.ResolveBatch.
// It performs sophisticated matching using the new algorithm and persistent caching.
func (r *MSUResolver) ResolveBatch(ctx context.Context, req []idresolver.ResolveRequestItem) ([]idresolver.ResolveResponseItem, error) {
	slog.Info("Starting batch resolution", "requestCount", len(req))

	// Check if data is stale and trigger fetch if needed
	if stale, err := r.dbStore.IsDataStale(time.Hour); err == nil && stale {
		slog.Info("Data is stale, checking for fresh data fetch")
		r.triggerFetchIfNeeded(ctx)
	}

	// First, try to get cached results
	cachedResults := make(map[string]idresolver.ResolveResponseItem)
	var uncachedRequests []idresolver.ResolveRequestItem

	for _, item := range req {
		if canonicalID, confidence, found := r.dbStore.GetMatchCache(item.InternalID); found {
			cachedResults[item.InternalID] = idresolver.ResolveResponseItem{
				InternalID:  item.InternalID,
				CanonicalID: canonicalID,
				Confidence:  confidence,
			}
		} else {
			uncachedRequests = append(uncachedRequests, item)
		}
	}

	slog.Info("Cache lookup results",
		"cachedResults", len(cachedResults),
		"uncachedRequests", len(uncachedRequests))

	// If we have uncached requests, perform sophisticated matching
	var newResults []matcher.MatchResult
	if len(uncachedRequests) > 0 {
		// Get all cached Gosuslugi data
		gosuslugiData := r.getAllCachedGosuslugiData()

		if len(gosuslugiData) == 0 {
			slog.Warn("No Gosuslugi data available, using fallback for all uncached requests")
			newResults = r.createFallbackResults(uncachedRequests)
		} else {
			// Execute sophisticated matching algorithm
			engine := matcher.NewMatchingEngine()
			newResults = engine.ExecuteMatching(uncachedRequests, gosuslugiData)
		}

		// Cache the new results
		for _, result := range newResults {
			err := r.dbStore.SetMatchCache(
				result.InternalID,
				result.CanonicalID,
				result.Confidence,
				result.ProgramName,
				result.CompetitionType,
			)
			if err != nil {
				slog.Error("Failed to cache match result", "error", err, "internalID", result.InternalID)
			}
		}
	}

	// Combine cached and new results
	resp := make([]idresolver.ResolveResponseItem, len(req))
	for i, item := range req {
		if cached, exists := cachedResults[item.InternalID]; exists {
			resp[i] = cached
		} else {
			// Find in new results
			found := false
			for _, result := range newResults {
				if result.InternalID == item.InternalID {
					resp[i] = idresolver.ResolveResponseItem{
						InternalID:  result.InternalID,
						CanonicalID: result.CanonicalID,
						Confidence:  result.Confidence,
					}
					found = true
					break
				}
			}
			if !found {
				// Fallback
				resp[i] = idresolver.ResolveResponseItem{
					InternalID:  item.InternalID,
					CanonicalID: r.generateFallbackID(item.InternalID),
					Confidence:  0.0,
				}
			}
		}
	}

	slog.Info("Batch resolution completed",
		"totalRequests", len(req),
		"cachedResults", len(cachedResults),
		"newResults", len(newResults))

	return resp, nil
}

// StartBackgroundFetcher starts the background data fetching process.
// It performs an immediate fetch and then sets up periodic fetching every 2 hours.
func (r *MSUResolver) StartBackgroundFetcher(ctx context.Context) {
	r.mu.Lock()
	if r.fetcherStarted {
		r.mu.Unlock()
		return
	}
	r.fetcherStarted = true
	r.mu.Unlock()

	slog.Info("Starting background fetcher...")

	go func() {
		// Perform immediate fetch only if data is stale (older than 2 hours)
		if time.Since(r.lastSuccessfulFetch) > 1*time.Hour {
			slog.Info("Data is stale, performing initial fetch")
			if err := r.fetchGosuslugiDataWithRun(ctx); err != nil {
				slog.Error("Initial fetch failed", "error", err)
			}
		} else {
			slog.Info("Data is fresh, skipping initial fetch",
				"lastFetch", r.lastSuccessfulFetch,
				"age", time.Since(r.lastSuccessfulFetch))
		}

		// Set up periodic fetching every 2 hours
		ticker := time.NewTicker(45 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				slog.Info("Background fetcher stopped")
				return
			case <-ticker.C:
				slog.Info("Performing periodic fetch")
				if err := r.fetchGosuslugiDataWithRun(ctx); err != nil {
					slog.Error("Periodic fetch failed", "error", err)
				}
			}
		}
	}()
}

// GosuslugiAPIResponse represents the response from Gosuslugi API
type GosuslugiAPIResponse struct {
	UpdateDate string               `json:"updateDate"`
	Applicants []GosuslugiApplicant `json:"applicants"`
}

// GosuslugiApplicant represents an applicant from the Gosuslugi API
type GosuslugiApplicant struct {
	Rating                   int     `json:"rating"`
	Priority                 int     `json:"priority"`
	Consent                  string  `json:"consent"`
	ConsentDate              string  `json:"consentDate"`
	SumMark                  float64 `json:"sumMark"`
	WithoutTests             bool    `json:"withoutTests"`
	Result1                  float64 `json:"result1"`
	Result2                  float64 `json:"result2"`
	Result3                  float64 `json:"result3"`
	Result4                  float64 `json:"result4"`
	Result5                  float64 `json:"result5"`
	Result6                  float64 `json:"result6"`
	Result7                  float64 `json:"result7"`
	Result8                  float64 `json:"result8"`
	AchievementsMark         float64 `json:"achievementsMark"`
	StatusID                 int     `json:"statusId"`
	StatusName               string  `json:"statusName"`
	IDApplication            int     `json:"idApplication"`
	CompetitionSelectionDate string  `json:"competitionSelectionDate"`

	OfferPlace             *int        `json:"offerPlace"`
	OfferOrganization      interface{} `json:"offerOrganization"`
	TargetAchievementsMark float64     `json:"targetAchievementsMark"`
	PaidContract           bool        `json:"paidContract"`
}

// fetchGosuslugiDataWithRun fetches data from Gosuslugi API and stores it with run tracking
func (r *MSUResolver) fetchGosuslugiDataWithRun(ctx context.Context) error {
	// Prevent concurrent fetches
	r.mu.Lock()
	if r.fetchInProgress {
		r.mu.Unlock()
		slog.Info("Fetch already in progress, skipping")
		return nil
	}
	r.fetchInProgress = true
	r.mu.Unlock()

	defer func() {
		r.mu.Lock()
		r.fetchInProgress = false
		r.mu.Unlock()
	}()

	slog.Info("Starting Gosuslugi data fetch...")

	// Get competition IDs for all programs
	competitionIDs := getMSUCompetitionIDs()
	totalPrograms := len(competitionIDs)

	// Start a fetch run
	runID, err := r.dbStore.StartFetchRun(totalPrograms)
	if err != nil {
		return fmt.Errorf("failed to start fetch run: %w", err)
	}

	processedPrograms := 0
	var fetchErrors []string

	// Fetch data for each program
	for programName, programIDs := range competitionIDs {
		processedPrograms++
		slog.Info("Processing program",
			"progress", fmt.Sprintf("%d/%d", processedPrograms, totalPrograms),
			"name", programName)

		// Update progress
		if err := r.dbStore.UpdateFetchRun(runID, processedPrograms); err != nil {
			slog.Error("Failed to update fetch run progress", "error", err)
		}

		// Fetch data for each quota type
		for quotaType, competitionID := range map[string]string{
			"RegularBVI":     programIDs.RegularBVI,
			"DedicatedQuota": programIDs.DedicatedQuota,
			"SpecialQuota":   programIDs.SpecialQuota,
			"TargetQuota":    programIDs.TargetQuota,
		} {
			if competitionID == "" {
				continue // Skip empty competition IDs
			}

			slog.Debug("Fetching competition data",
				"quotaType", quotaType,
				"competitionID", competitionID)

			if err := r.fetchCompetitionData(ctx, competitionID, programName, quotaType); err != nil {
				errorMsg := fmt.Sprintf("Failed to fetch %s for %s (ID: %s): %v", quotaType, programName, competitionID, err)
				fetchErrors = append(fetchErrors, errorMsg)
				slog.Error("Competition data fetch failed", "error", err, "quotaType", quotaType, "competitionID", competitionID)
				continue
			}

			// Add delay between requests to avoid overwhelming the API and reduce network issues
			time.Sleep(200 * time.Millisecond)
		}

		// Add longer delay between programs to reduce network stress
		time.Sleep(1 * time.Second)
	}

	// Complete the fetch run
	status := "completed"
	errorMessage := ""
	if len(fetchErrors) > 0 {
		status = "partial"
		errorMessage = fmt.Sprintf("Some fetches failed: %v", fetchErrors)
		slog.Warn("Fetch completed with errors", "errorCount", len(fetchErrors))
	}

	if err := r.dbStore.CompleteFetchRun(runID, status, errorMessage); err != nil {
		slog.Error("Failed to complete fetch run", "error", err)
	}

	// Clear match cache since we have new data
	if err := r.dbStore.ClearMatchCache(); err != nil {
		slog.Error("Failed to clear match cache", "error", err)
	} else {
		slog.Info("Match cache cleared for fresh matching")
	}

	// Store last fetch time
	r.mu.Lock()
	r.lastSuccessfulFetch = time.Now()
	r.mu.Unlock()

	slog.Info("Gosuslugi data fetch completed",
		"totalPrograms", totalPrograms,
		"processedPrograms", processedPrograms,
		"errors", len(fetchErrors),
		"status", status)

	return nil
}

// fetchCompetitionData fetches data for a specific competition from Gosuslugi API with retry logic
func (r *MSUResolver) fetchCompetitionData(ctx context.Context, competitionID, programName, quotaType string) error {
	const maxRetries = 3
	const baseDelay = 1 * time.Second

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff with jitter
			delay := baseDelay * time.Duration(1<<attempt)
			if delay > 10*time.Second {
				delay = 10 * time.Second
			}

			slog.Debug("Retrying request after delay",
				"attempt", attempt+1,
				"maxRetries", maxRetries,
				"delay", delay,
				"competitionID", competitionID)

			// Don't let parent context cancellation interrupt retry delays
			// Use a simple sleep instead of context-aware sleep
			time.Sleep(delay)
		}

		err := r.fetchCompetitionDataWithTimeout(competitionID, programName, quotaType)
		if err == nil {
			if attempt > 0 {
				slog.Info("Request succeeded after retry",
					"attempt", attempt+1,
					"competitionID", competitionID)
			}
			return nil
		}

		lastErr = err
		slog.Warn("Request failed, will retry if attempts remain",
			"attempt", attempt+1,
			"maxRetries", maxRetries,
			"error", err,
			"competitionID", competitionID)
	}

	return fmt.Errorf("failed after %d attempts: %w", maxRetries, lastErr)
}

// fetchCompetitionDataWithTimeout performs a single attempt to fetch competition data with extended timeout
func (r *MSUResolver) fetchCompetitionDataWithTimeout(competitionID, programName, quotaType string) error {
	// Create a completely independent context with longer timeout for this specific request
	// This prevents any parent context cancellation from affecting individual requests
	requestCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Construct the Gosuslugi API URL
	url := fmt.Sprintf("https://www.gosuslugi.ru/api/university/v1/public/competition/%s/applicants", competitionID)

	// Create HTTP request with the request-specific context
	req, err := http.NewRequestWithContext(requestCtx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("User-Agent", "MSU-Resolver/1.0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Connection", "close") // Prevent connection reuse issues

	// Create HTTP client with extended timeout and no connection reuse
	client := &http.Client{
		Timeout: 90 * time.Second, // Longer timeout for unreliable networks
		Transport: &http.Transport{
			DisableKeepAlives:     true, // Prevent connection reuse issues
			DisableCompression:    false,
			MaxIdleConns:          0,
			MaxIdleConnsPerHost:   0,
			IdleConnTimeout:       0,
			ResponseHeaderTimeout: 30 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Parse JSON response
	var apiResp GosuslugiAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Store data in cache with a composite key
	cacheKey := fmt.Sprintf("%s:%s:%s", programName, quotaType, competitionID)
	r.cache.Set(cacheKey, apiResp.Applicants)

	log.Printf("[idmsu] Stored %d applicants for competition %s", len(apiResp.Applicants), competitionID)
	return nil
}

// HasRecentData returns true if data has been fetched recently.
func (r *MSUResolver) HasRecentData() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Data is considered recent if it's been fetched in the last 45 minutes
	has := !r.lastSuccessfulFetch.IsZero() && time.Since(r.lastSuccessfulFetch) < 45*time.Minute

	if len(r.getAllCachedGosuslugiData()) == 0 {
		go func() {
			fetchCtx := context.Background()
			if err := r.fetchGosuslugiDataWithRun(fetchCtx); err != nil {
				slog.Error("Triggered fetch failed", "error", err)
			}
		}()

		return false
	}

	return has
}

// triggerFetchIfNeeded triggers a fetch if data is stale and no fetch is in progress
func (r *MSUResolver) triggerFetchIfNeeded(ctx context.Context) {
	r.mu.RLock()
	inProgress := r.fetchInProgress
	r.mu.RUnlock()

	if !inProgress {
		go func() {
			// Create a new independent context for the fetch operation
			// This prevents the fetch from being cancelled when the incoming request times out
			fetchCtx := context.Background()
			if err := r.fetchGosuslugiDataWithRun(fetchCtx); err != nil {
				slog.Error("Triggered fetch failed", "error", err)
			}
		}()
	}
}

// getAllCachedGosuslugiData retrieves all cached Gosuslugi data organized by program and competition type
func (r *MSUResolver) getAllCachedGosuslugiData() map[string]map[string][]matcher.GosuslugiApplicant {
	result := make(map[string]map[string][]matcher.GosuslugiApplicant)

	// Get all competition IDs
	competitionIDs := getMSUCompetitionIDs()

	for programName, programIDs := range competitionIDs {
		programData := make(map[string][]matcher.GosuslugiApplicant)

		// Try to get data for each competition type using the correct cache key format
		for _, quotaType := range []string{"RegularBVI", "DedicatedQuota", "SpecialQuota", "TargetQuota"} {
			var competitionID string
			switch quotaType {
			case "RegularBVI":
				competitionID = programIDs.RegularBVI
			case "DedicatedQuota":
				competitionID = programIDs.DedicatedQuota
			case "SpecialQuota":
				competitionID = programIDs.SpecialQuota
			case "TargetQuota":
				competitionID = programIDs.TargetQuota
			}

			if competitionID == "" {
				continue // Skip empty competition IDs
			}

			cacheKey := fmt.Sprintf("%s:%s:%s", programName, quotaType, competitionID)
			if cached, exists := r.cache.Get(cacheKey); exists {
				if applicants, ok := cached.([]GosuslugiApplicant); ok {
					// Convert to matcher format
					matcherApplicants := make([]matcher.GosuslugiApplicant, len(applicants))
					for i, app := range applicants {
						matcherApplicants[i] = matcher.GosuslugiApplicant{
							Rating:                   app.Rating,
							Priority:                 app.Priority,
							Consent:                  app.Consent,
							ConsentDate:              app.ConsentDate,
							SumMark:                  app.SumMark,
							WithoutTests:             app.WithoutTests,
							Result1:                  app.Result1,
							Result2:                  app.Result2,
							Result3:                  app.Result3,
							Result4:                  app.Result4,
							Result5:                  app.Result5,
							Result6:                  app.Result6,
							Result7:                  app.Result7,
							Result8:                  app.Result8,
							AchievementsMark:         app.AchievementsMark,
							StatusID:                 app.StatusID,
							StatusName:               app.StatusName,
							IDApplication:            app.IDApplication,
							CompetitionSelectionDate: app.CompetitionSelectionDate,
							OfferPlace:               app.OfferPlace,
							OfferOrganization:        app.OfferOrganization,
							TargetAchievementsMark:   app.TargetAchievementsMark,
							PaidContract:             app.PaidContract,
						}
					}
					programData[quotaType] = matcherApplicants
				}
			}
		}

		if len(programData) > 0 {
			result[programName] = programData
		}
	}

	slog.Debug("Retrieved cached Gosuslugi data", "programCount", len(result))
	return result
}

// createFallbackResults creates fallback results for requests when no Gosuslugi data is available
func (r *MSUResolver) createFallbackResults(requests []idresolver.ResolveRequestItem) []matcher.MatchResult {
	results := make([]matcher.MatchResult, len(requests))
	for i, req := range requests {
		results[i] = matcher.MatchResult{
			InternalID:      req.InternalID,
			CanonicalID:     r.generateFallbackID(req.InternalID),
			Confidence:      0.0,
			MatchType:       "Fallback",
			ProgramName:     "Unknown",
			CompetitionType: "Unknown",
		}
	}
	return results
}

// findBestMatch attempts to find the best matching canonical ID for the given internal ID
// by comparing the applicant's data with cached Gosuslugi entries
func (r *MSUResolver) findBestMatch(internalID string, apps []idresolver.MSUAppDetails) (string, float64) {
	if len(apps) == 0 {
		return r.generateFallbackID(internalID), 0.0
	}

	bestCanonicalID := ""
	bestConfidence := 0.0

	// Try to find matches for each application
	for _, app := range apps {
		// Get cached candidates for this program
		candidates := r.getCachedCandidates(app.PrettyName)
		if len(candidates) == 0 {
			continue
		}

		// Find the best match among candidates
		for _, candidate := range candidates {
			confidence := r.calculateMatchConfidence(app, candidate)
			if confidence > bestConfidence {
				bestConfidence = confidence
				bestCanonicalID = fmt.Sprintf("%d", candidate.IDApplication)
			}
		}
	}

	// If no good match found, return fallback
	if bestConfidence < 0.5 {
		return r.generateFallbackID(internalID), bestConfidence
	}

	return bestCanonicalID, bestConfidence
}

// getCachedCandidates retrieves cached Gosuslugi entries for a given program name
func (r *MSUResolver) getCachedCandidates(programName string) []GosuslugiApplicant {
	var candidates []GosuslugiApplicant

	// Get all cached keys and filter for this program
	if cachedData, exists := r.cache.Get(fmt.Sprintf("%s:RegularBVI", programName)); exists {
		if applicants, ok := cachedData.([]GosuslugiApplicant); ok {
			candidates = append(candidates, applicants...)
		}
	}

	if cachedData, exists := r.cache.Get(fmt.Sprintf("%s:DedicatedQuota", programName)); exists {
		if applicants, ok := cachedData.([]GosuslugiApplicant); ok {
			candidates = append(candidates, applicants...)
		}
	}

	if cachedData, exists := r.cache.Get(fmt.Sprintf("%s:SpecialQuota", programName)); exists {
		if applicants, ok := cachedData.([]GosuslugiApplicant); ok {
			candidates = append(candidates, applicants...)
		}
	}

	if cachedData, exists := r.cache.Get(fmt.Sprintf("%s:TargetQuota", programName)); exists {
		if applicants, ok := cachedData.([]GosuslugiApplicant); ok {
			candidates = append(candidates, applicants...)
		}
	}

	return candidates
}

// calculateMatchConfidence calculates how well an MSU app matches a Gosuslugi candidate
func (r *MSUResolver) calculateMatchConfidence(app idresolver.MSUAppDetails, candidate GosuslugiApplicant) float64 {
	confidence := 0.0

	// Skip excluded candidates
	if candidate.StatusID == 4 {
		return 0.0
	}

	// Score matching (most important factor)
	if app.ScoreSum > 0 {
		// Parse candidate's total score
		candidateScore := r.parseTotalScore(candidate)
		if candidateScore > 0 {
			scoreDiff := abs(app.ScoreSum - candidateScore)
			if scoreDiff <= 2 {
				confidence += 0.6 // Exact or very close score match
			} else if scoreDiff <= 5 {
				confidence += 0.4 // Close score match
			} else if scoreDiff <= 10 {
				confidence += 0.2 // Reasonable score match
			}
		}
	}

	// EGE scores matching (if available)
	if len(app.EGEScores) > 0 {
		egeMatch := r.compareEGEScores(app.EGEScores, candidate)
		confidence += egeMatch * 0.3
	}

	// Priority matching (if meaningful)
	if app.Priority > 0 {
		confidence += 0.1 // Small bonus for having priority info
	}

	return confidence
}

// parseTotalScore extracts the total score from a Gosuslugi candidate
func (r *MSUResolver) parseTotalScore(candidate GosuslugiApplicant) int {
	// Use SumMark as the total score, convert from float64 to int
	return int(candidate.SumMark)
}

// parseScoreString converts a score string to integer
func (r *MSUResolver) parseScoreString(scoreStr string) int {
	if scoreStr == "" {
		return 0
	}
	// Simple parsing - in real implementation might need more robust parsing
	var score int
	if _, err := fmt.Sscanf(scoreStr, "%d", &score); err == nil {
		return score
	}
	return 0
}

// compareEGEScores compares EGE scores between MSU app and Gosuslugi candidate
func (r *MSUResolver) compareEGEScores(msuScores []int, candidate GosuslugiApplicant) float64 {
	// Extract EGE scores from candidate
	candidateScores := r.extractEGEScores(candidate)
	if len(candidateScores) == 0 || len(msuScores) == 0 {
		return 0.0
	}

	// Sort both score arrays for comparison
	msuSorted := make([]int, len(msuScores))
	copy(msuSorted, msuScores)
	sort.Ints(msuSorted)

	candidateSorted := make([]int, len(candidateScores))
	copy(candidateSorted, candidateScores)
	sort.Ints(candidateSorted)

	// Compare sorted scores
	matches := 0
	minLen := len(msuSorted)
	if len(candidateSorted) < minLen {
		minLen = len(candidateSorted)
	}

	for i := 0; i < minLen; i++ {
		if abs(msuSorted[i]-candidateSorted[i]) <= 2 {
			matches++
		}
	}

	if minLen == 0 {
		return 0.0
	}
	return float64(matches) / float64(minLen)
}

// extractEGEScores extracts EGE scores from Gosuslugi candidate
func (r *MSUResolver) extractEGEScores(candidate GosuslugiApplicant) []int {
	var scores []int

	// Extract individual test results, convert from float64 to int
	results := []float64{candidate.Result1, candidate.Result2, candidate.Result3, candidate.Result4, candidate.Result5, candidate.Result6, candidate.Result7, candidate.Result8}

	for _, result := range results {
		if result > 0 {
			scores = append(scores, int(result))
		}
	}

	return scores
}

// generateFallbackID creates a fallback canonical ID for unmatched internal IDs
func (r *MSUResolver) generateFallbackID(internalID string) string {
	return fmt.Sprintf("MSU-%s", internalID)
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
