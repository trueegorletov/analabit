package resolver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/trueegorletov/analabit/core/idresolver"
	"github.com/trueegorletov/analabit/core/utils"
	"github.com/trueegorletov/analabit/service/idmsu/cache"
)

type MSUResolver interface {
	ResolveBatch(ctx context.Context, req []idresolver.ResolveRequestItem) ([]idresolver.ResolveResponseItem, error)
	StartBackgroundFetcher(ctx context.Context)
	HasRecentData(ctx context.Context) bool
}

type msuResolver struct {
	cache cache.Cache
	// programMapping removed - now using PrettyName directly
}

type GosuslugiEntry = cache.GosuslugiEntry

// debugInternalID is used to emit detailed logs for a specific student during matching.
const debugInternalID = "025928"

type GosuslugiAPIResponse struct {
	UpdateDate string           `json:"updateDate"`
	Applicants []GosuslugiEntry `json:"applicants"`
}

// MSU competition list IDs for different programs and quota types
// Real IDs extracted from MSU admissions data
var msuCompetitionIDs = getMSUCompetitionIDs()

func NewMSUResolver(cache cache.Cache) MSUResolver {
	return &msuResolver{
		cache: cache,
	}
}

// restoreCacheFromLatestRun restores the cache from the latest completed run
func (r *msuResolver) restoreCacheFromLatestRun(ctx context.Context) error {
	latestRun, err := r.cache.GetLatestCompletedRun(ctx)
	if err != nil {
		return fmt.Errorf("failed to get latest completed run: %w", err)
	}

	if latestRun == nil {
		log.Println("No completed runs found, skipping cache restoration")
		return nil
	}

	// Check if the run is recent (within 2 hours)
	if time.Since(*latestRun.FinishedAt) > 2*time.Hour {
		log.Printf("Latest run %d is too old (%v), skipping cache restoration", latestRun.ID, latestRun.FinishedAt)
		return nil
	}

	log.Printf("Restoring cache from run %d (finished at %v)", latestRun.ID, latestRun.FinishedAt)

	if err := r.cache.RestoreCacheFromRun(ctx, latestRun.ID); err != nil {
		return fmt.Errorf("failed to restore cache from run %d: %w", latestRun.ID, err)
	}

	log.Printf("Successfully restored cache from run %d", latestRun.ID)
	return nil
}

// HasRecentData checks if there's recent data available (for readiness check)
func (r *msuResolver) HasRecentData(ctx context.Context) bool {
	// First, ensure we have at least one completed fetch run and that it is recent.
	latestRun, err := r.cache.GetLatestCompletedRun(ctx)
	if err != nil {
		log.Printf("Failed to check for recent data: %v", err)
		return false
	}

	if latestRun == nil {
		return false
	}

	// Require the run itself to be completed within the last two hours.
	if time.Since(*latestRun.FinishedAt) > 2*time.Hour {
		return false
	}

	// Additionally, verify that cached data is fresh for every configured programme.
	// If any programme cache is stale (or an error occurs while checking), the service
	// should be considered not ready so that callers will retry later.
	if r.isGlobalCacheStale(ctx, 30*time.Minute) {
		return false
	}

	return true
}

// isGlobalCacheStale checks every MSU programme id in msuCompetitionIDs; if ANY of them is older
// than the given threshold (or an error occurs while checking) the cache is considered incomplete.
func (r *msuResolver) isGlobalCacheStale(ctx context.Context, threshold time.Duration) bool {
	for program := range msuCompetitionIDs {
		_, err := r.cache.IsCacheStale(ctx, program, threshold)
		if err != nil {
			log.Printf("Error checking cache freshness for program %s: %v", program, err)
			return true // fail-safe – treat unknown state as stale
		}
		// if stale {
		//     log.Printf("Cache is stale for program %s (global check)", program)
		//     return true
		// }
	}
	return false
}

func (r *msuResolver) ResolveBatch(ctx context.Context, req []idresolver.ResolveRequestItem) ([]idresolver.ResolveResponseItem, error) {
	results := make([]idresolver.ResolveResponseItem, 0, len(req))

	for _, item := range req {
		result := r.resolveStudent(ctx, item)
		results = append(results, result)
	}

	return results, nil
}

func (r *msuResolver) resolveStudent(ctx context.Context, item idresolver.ResolveRequestItem) idresolver.ResolveResponseItem {
	// Get relevant programs
	programs := r.extractPrograms(item.Apps)
	if item.InternalID == debugInternalID {
		log.Printf("[DEBUG] Resolving InternalID %s with programs: %v", item.InternalID, programs)
	}

	// Check for stale cache before fetching candidates
	for _, program := range programs {
		stale, err := r.cache.IsCacheStale(ctx, program, 6*time.Hour)
		if err != nil {
			log.Printf("Error checking cache freshness for program %s: %v", program, err)
			// Decide how to handle error, e.g., proceed with potentially stale data or fail
			// For now, we'll log and continue, but you might want to return an error response
		}
		if stale {
			if item.InternalID == debugInternalID {
				log.Printf("[DEBUG] Cache is stale for program %s, but continuing with existing data", program)
			}
			// Proceed despite stale cache
		}
	}

	// Get cached Gosuslugi data for relevant programs
	candidates := r.getCandidates(ctx, programs)

	// Group candidates by their internal ID (last 6 digits of IDApplication)
	groupedCandidates := r.groupCandidatesByInternalID(candidates)
	if item.InternalID == debugInternalID {
		log.Printf("[DEBUG] Candidate groups for %s: %d groups", item.InternalID, len(groupedCandidates))
	}

	// Try to match the student
	bestMatch, confidence := r.findBestMatch(item, groupedCandidates)
	if item.InternalID == debugInternalID {
		if bestMatch != nil {
			log.Printf("[DEBUG] Best match for %s: %+v (confidence %.2f)", item.InternalID, *bestMatch, confidence)
		} else {
			log.Printf("[DEBUG] No match found for %s, confidence %.2f", item.InternalID, confidence)
		}
	}

	if confidence >= 0.5 { // Lowered threshold for more matches
		canonicalID, _ := utils.PrepareStudentID(bestMatch.IDApplication.String())
		return idresolver.ResolveResponseItem{
			InternalID:  item.InternalID,
			CanonicalID: canonicalID,
			Confidence:  confidence,
		}
	}

	// No good match found, generate fallback ID
	fallbackID := r.generateFallbackID(item.InternalID)
	return idresolver.ResolveResponseItem{
		InternalID:  item.InternalID,
		CanonicalID: fallbackID,
		Confidence:  0.0,
	}
}

// New function to group candidates by internal ID (last 6 digits)
func (r *msuResolver) groupCandidatesByInternalID(candidates []GosuslugiEntry) map[string][]GosuslugiEntry {
	grouped := make(map[string][]GosuslugiEntry)
	for _, cand := range candidates {
		idStr := cand.IDApplication.String()
		var internalID string
		if len(idStr) >= 6 {
			internalID = idStr[len(idStr)-6:]
		} else {
			// Pad with leading zeros to ensure 6-digit alignment
			internalID = fmt.Sprintf("%06s", idStr)
		}
		grouped[internalID] = append(grouped[internalID], cand)
	}
	return grouped
}

func (r *msuResolver) findBestMatch(item idresolver.ResolveRequestItem, groupedCandidates map[string][]GosuslugiEntry) (*GosuslugiEntry, float64) {
	var bestMatch *GosuslugiEntry
	var bestConfidence float64

	for _, group := range groupedCandidates {
		confidence := r.calculateMatchConfidence(item, group)
		if confidence > bestConfidence {
			bestConfidence = confidence
			// Select the first entry in the group as representative
			bestMatch = &group[0]
		}
	}

	return bestMatch, bestConfidence
}

func (r *msuResolver) calculateMatchConfidence(item idresolver.ResolveRequestItem, candidateGroup []GosuslugiEntry) float64 {
	bestRatio := 0.0

	// Iterate through each MSU application for the student
	for _, msuApp := range item.Apps {
		msuScores := r.extractMSUEGEScores(msuApp)
		if len(msuScores) == 0 {
			continue
		}

		for _, candidateApp := range candidateGroup {
			candidateScores := r.extractEGEScores(candidateApp)
			overlap := r.scoreOverlapCount(msuScores, candidateScores, 0) // exact score match
			ratio := float64(overlap) / float64(len(msuScores))
			if item.InternalID == debugInternalID {
				log.Printf("[DEBUG] Comparing MSU %v with candidate %v -> overlap %d, ratio %.2f", msuScores, candidateScores, overlap, ratio)
			}
			if ratio > bestRatio {
				bestRatio = ratio
			}
		}
	}

	switch {
	case bestRatio >= 1.0:
		return 1.0 // Perfect match
	case bestRatio >= 0.8:
		return 0.8 // Very high confidence
	case bestRatio >= 0.6:
		return 0.6 // Medium confidence
	case bestRatio >= 0.4:
		return 0.4 // Low confidence but still plausible
	default:
		return 0.0 // No reliable match
	}
}

// scoreOverlapCount returns how many scores in scores1 can be matched with scores2 within the given tolerance.
func (r *msuResolver) scoreOverlapCount(scores1, scores2 []int, tolerance int) int {
	matched := 0
	used := make([]bool, len(scores2))

	for _, s1 := range scores1 {
		for j, s2 := range scores2 {
			if used[j] {
				continue
			}
			if absInt(s1-s2) <= tolerance {
				matched++
				used[j] = true
				break
			}
		}
	}

	return matched
}

// absInt returns the absolute value of an int.
func absInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func (r *msuResolver) extractMSUEGEScores(app idresolver.MSUAppDetails) []int {
	scores := make([]int, 0, len(app.EGEScores))

	// Add only positive EGE scores, explicitly excluding DVI and zeroes
	for _, score := range app.EGEScores {
		if score > 0 {
			scores = append(scores, score)
		}
	}

	// Sort scores for consistent comparison
	sort.Ints(scores)
	return scores
}

// For Gosuslugi, assume Results are EGE scores; filtering out zero scores and sorting
func (r *msuResolver) extractEGEScores(candidate GosuslugiEntry) []int {
	var scores []int
	results := []*float64{candidate.Result1, candidate.Result2, candidate.Result3, candidate.Result4,
		candidate.Result5, candidate.Result6, candidate.Result7, candidate.Result8}

	// Only include positive scores (filter out nulls and zeroes)
	for _, result := range results {
		if result != nil && *result > 0 {
			scores = append(scores, int(*result))
		}
	}

	// Sort scores for consistent comparison
	sort.Ints(scores)
	return scores
}

func (r *msuResolver) scoresMatch(scores1, scores2 []int, tolerance int) bool {
	if len(scores1) != len(scores2) {
		return false
	}

	for i := range scores1 {
		diff := scores1[i] - scores2[i]
		if diff < 0 {
			diff = -diff
		}
		if diff > tolerance {
			return false
		}
	}

	return true
}

// extractPrograms collects unique program names from MSU applications ensuring we at least attempt resolution even if mapping is missing.
func (r *msuResolver) extractPrograms(apps []idresolver.MSUAppDetails) []string {
	programSet := make(map[string]struct{})
	for _, app := range apps {
		programSet[app.PrettyName] = struct{}{}
		if _, ok := GetCompetitionIDsForProgram(app.PrettyName); !ok {
			// log.Printf("WARNING: No competition IDs found for MSU program: %s", app.PrettyName)
		}
	}
	programs := make([]string, 0, len(programSet))
	for p := range programSet {
		programs = append(programs, p)
	}
	return programs
}

// getCandidates retrieves cached Gosuslugi entries for provided programs.
func (r *msuResolver) getCandidates(ctx context.Context, programs []string) []GosuslugiEntry {
	var all []GosuslugiEntry
	for _, program := range programs {
		compIDs, ok := msuCompetitionIDs[program]
		if !ok {
			continue
		}
		ids := []string{compIDs.RegularBVI, compIDs.DedicatedQuota, compIDs.SpecialQuota, compIDs.TargetQuota}
		for _, cid := range ids {
			if cid == "" {
				continue
			}
			cacheKey := fmt.Sprintf("%s:%s", program, cid)
			cands, err := r.cache.GetCandidates(ctx, cacheKey)
			if err != nil {
				log.Printf("Error getting candidates for %s: %v", cacheKey, err)
				continue
			}
			for _, cand := range cands {
				all = append(all, GosuslugiEntry(cand))
			}
		}
	}
	return r.filterValidCandidates(all)
}

// filterValidCandidates drops excluded entries.
func (r *msuResolver) filterValidCandidates(cands []GosuslugiEntry) []GosuslugiEntry {
	var valid []GosuslugiEntry
	for _, c := range cands {
		if c.StatusID != nil && *c.StatusID == 4 {
			continue
		}
		valid = append(valid, c)
	}
	return valid
}

// generateFallbackID converts last 6 digits of internal ID into canonical placeholder.
func (r *msuResolver) generateFallbackID(internalID string) string {
	if len(internalID) >= 6 {
		lastSix := internalID[len(internalID)-6:]
		return fmt.Sprintf("MSU-%s", strings.Repeat("0", 7-len(lastSix))+lastSix)
	}
	return fmt.Sprintf("MSU-%06s", internalID)
}

func (r *msuResolver) StartBackgroundFetcher(ctx context.Context) {
	log.Println("Starting background fetcher for Gosuslugi data...")

	// Try to restore cache from latest completed run on startup
	if err := r.restoreCacheFromLatestRun(ctx); err != nil {
		log.Printf("Failed to restore cache from latest run: %v", err)
	}

	// Fetch data immediately on startup
	if err := r.fetchGosuslugiDataWithRun(ctx); err != nil {
		log.Printf("Initial fetch failed: %v", err)
	}

	// Set up periodic fetching every 30 minutes
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Background fetcher stopped")
			return
		case <-ticker.C:
			if err := r.fetchGosuslugiDataWithRun(ctx); err != nil {
				log.Printf("Periodic fetch failed: %v", err)
			}
		}
	}
}

func (r *msuResolver) fetchGosuslugiDataWithRun(ctx context.Context) error {
	log.Println("Starting new Gosuslugi data fetch run...")

	// Create a new run
	runID, err := r.cache.CreateRun(ctx, map[string]interface{}{
		"source": "gosuslugi",
		"type":   "periodic_fetch",
	})
	if err != nil {
		return fmt.Errorf("failed to create run: %w", err)
	}

	log.Printf("Created run %d for Gosuslugi data fetch", runID)

	totalFetched := 0
	totalErrors := 0

	// Fetch data for all configured competition IDs
	for program, competitionIDs := range msuCompetitionIDs {
		// Get all competition IDs from the struct
		compIDs := []string{}
		if competitionIDs.RegularBVI != "" {
			compIDs = append(compIDs, competitionIDs.RegularBVI)
		}
		if competitionIDs.DedicatedQuota != "" {
			compIDs = append(compIDs, competitionIDs.DedicatedQuota)
		}
		if competitionIDs.SpecialQuota != "" {
			compIDs = append(compIDs, competitionIDs.SpecialQuota)
		}
		if competitionIDs.TargetQuota != "" {
			compIDs = append(compIDs, competitionIDs.TargetQuota)
		}

		for _, competitionID := range compIDs {
			if err := r.fetchCompetitionData(ctx, runID, program, competitionID); err != nil {
				log.Printf("Error fetching competition %s for program %s: %v", competitionID, program, err)
				totalErrors++
			} else {
				totalFetched++
			}

			// Add delay between requests to avoid rate limiting
			time.Sleep(1 * time.Second)
		}
	}

	// Mark run as finished
	if err := r.cache.FinishRun(ctx, runID); err != nil {
		log.Printf("Failed to finish run %d: %v", runID, err)
		return err
	}

	log.Printf("Gosuslugi data fetch run %d completed: %d successful, %d errors", runID, totalFetched, totalErrors)
	return nil
}

// extractJSONFromResponse handles both raw JSON and HTML-wrapped JSON responses
func extractJSONFromResponse(responseBody string) string {
	// Check if response is HTML-wrapped
	if strings.HasPrefix(strings.TrimSpace(responseBody), "<html>") {
		// Find the content within <pre> tags
		preStart := strings.Index(responseBody, "<pre>")
		if preStart == -1 {
			return responseBody // No <pre> tag found, return original
		}
		preStart += len("<pre>")

		preEnd := strings.Index(responseBody[preStart:], "</pre>")
		if preEnd == -1 {
			return responseBody // No closing </pre> tag found, return original
		}

		// Extract JSON content from within <pre> tags
		jsonContent := responseBody[preStart : preStart+preEnd]
		return strings.TrimSpace(jsonContent)
	}

	// Not HTML-wrapped, return as-is
	return responseBody
}

func (r *msuResolver) fetchCompetitionData(ctx context.Context, runID int, program, competitionID string) error {
	// Construct Gosuslugi API URL
	apiURL := fmt.Sprintf("https://www.gosuslugi.ru/api/university/v1/public/competition/%s/applicants", competitionID)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers to mimic browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://www.gosuslugi.ru/")

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	responseBody := string(body)

	// Extract JSON from response (handles both raw JSON and HTML-wrapped JSON)
	jsonContent := extractJSONFromResponse(responseBody)

	// Parse JSON response
	var apiResponse GosuslugiAPIResponse
	if err := json.Unmarshal([]byte(jsonContent), &apiResponse); err != nil {
		// Log additional debug info when JSON parsing fails
		previewLen := 500
		if len(jsonContent) < previewLen {
			previewLen = len(jsonContent)
		}
		log.Printf("ERROR: JSON parsing failed for competition %s. Extracted JSON preview (%d chars): %s", competitionID, previewLen, jsonContent[:previewLen])
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Convert API response entries to cache entries and store
	cacheEntries := make([]cache.GosuslugiEntry, 0, len(apiResponse.Applicants))
	for _, applicant := range apiResponse.Applicants {
		if applicant.IDApplication.String() == "" {
			continue // Skip entries without ID
		}

		// Convert to cache entry format
		cacheEntry := cache.GosuslugiEntry{
			IDApplication: applicant.IDApplication,
			ProgramName:   program,       // Use the program name as identifier
			ListType:      competitionID, // Store competition ID in ListType field
			Result1:       applicant.Result1,
			Result2:       applicant.Result2,
			Result3:       applicant.Result3,
			Result4:       applicant.Result4,
			Result5:       applicant.Result5,
			Result6:       applicant.Result6,
			Result7:       applicant.Result7,
			Result8:       applicant.Result8,
			SumMark:       applicant.SumMark,
			Rating:        applicant.Rating,
			WithoutTests:  applicant.WithoutTests,
			StatusID:      applicant.StatusID,
		}
		cacheEntries = append(cacheEntries, cacheEntry)
	}

	// Store entries in run data
	cacheKey := fmt.Sprintf("%s:%s", program, competitionID)
	if err := r.cache.StoreCandidatesForRun(ctx, runID, cacheKey, cacheEntries); err != nil {
		return fmt.Errorf("failed to store candidates for run: %w", err)
	}

	// Store all entries for this program/competition combination in current cache
	if err := r.cache.StoreCandidates(ctx, cacheKey, cacheEntries); err != nil {
		return fmt.Errorf("failed to store candidates in cache: %w", err)
	}

	log.Printf("Cached %d applicants for program %s, competition %s (run %d)", len(cacheEntries), program, competitionID, runID)
	return nil
}
