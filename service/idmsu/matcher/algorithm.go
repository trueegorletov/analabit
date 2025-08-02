package matcher

import (
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"github.com/trueegorletov/analabit/core/idresolver"
)

// MatchResult represents the result of matching an internal ID to a canonical ID
type MatchResult struct {
	InternalID      string
	CanonicalID     string
	Confidence      float64
	MatchType       string // "BVI", "Strong", "Weak", "Fallback"
	ProgramName     string
	CompetitionType string
}

// GosuslugiApplicant represents an applicant from the Gosuslugi API
type GosuslugiApplicant struct {
	Rating                   int         `json:"rating"`
	Priority                 int         `json:"priority"`
	Consent                  string      `json:"consent"`
	ConsentDate              string      `json:"consentDate"`
	SumMark                  float64     `json:"sumMark"`
	WithoutTests             bool        `json:"withoutTests"`
	Result1                  float64     `json:"result1"`
	Result2                  float64     `json:"result2"`
	Result3                  float64     `json:"result3"`
	Result4                  float64     `json:"result4"`
	Result5                  float64     `json:"result5"`
	Result6                  float64     `json:"result6"`
	Result7                  float64     `json:"result7"`
	Result8                  float64     `json:"result8"`
	AchievementsMark         float64     `json:"achievementsMark"`
	StatusID                 int         `json:"statusId"`
	StatusName               string      `json:"statusName"`
	IDApplication            int         `json:"idApplication"`
	CompetitionSelectionDate string      `json:"competitionSelectionDate"`
	OfferPlace               *int        `json:"offerPlace"`
	OfferOrganization        interface{} `json:"offerOrganization"`
	TargetAchievementsMark   float64     `json:"targetAchievementsMark"`
	PaidContract             bool        `json:"paidContract"`
}

// MatchingEngine implements the sophisticated MSU ID matching algorithm
type MatchingEngine struct {
	// Global map to store matches across all competition types to avoid conflicts
	globalMatches map[string]MatchResult
	// Track used canonical IDs to enforce mapping rules
	usedCanonicalIDs map[string][]string // canonical ID -> list of internal IDs that map to it
}

// NewMatchingEngine creates a new matching engine
func NewMatchingEngine() *MatchingEngine {
	return &MatchingEngine{
		globalMatches:    make(map[string]MatchResult),
		usedCanonicalIDs: make(map[string][]string),
	}
}

// ExecuteMatching performs the complete matching algorithm as specified
func (me *MatchingEngine) ExecuteMatching(
	requests []idresolver.ResolveRequestItem,
	gosuslugiData map[string]map[string][]GosuslugiApplicant,
) []MatchResult {
	slog.Info("Starting sophisticated MSU ID matching algorithm",
		"requestCount", len(requests),
		"programsInCache", len(gosuslugiData))

	// Group requests by program name
	programGroups := me.groupRequestsByProgram(requests)

	// Process each program separately
	for programName, programRequests := range programGroups {
		slog.Info("Processing program", "name", programName, "requestCount", len(programRequests))

		programData, exists := gosuslugiData[programName]
		if !exists {
			slog.Warn("No Gosuslugi data found for program", "name", programName)
			me.handleMissingProgramData(programRequests, programName)
			continue
		}

		// Process competition types in specific order: BVI first, then others
		me.processBVI(programRequests, programData, programName)
		me.processOtherCompetitionTypes(programRequests, programData, programName)
	}

	// Convert global matches to slice and handle remaining unmatched
	results := me.finalizeResults(requests)

	slog.Info("MSU ID matching completed",
		"totalMatches", len(results),
		"highConfidenceMatches", me.countByConfidence(results, 0.8))

	return results
}

// groupRequestsByProgram groups requests by program name
func (me *MatchingEngine) groupRequestsByProgram(requests []idresolver.ResolveRequestItem) map[string][]idresolver.ResolveRequestItem {
	groups := make(map[string][]idresolver.ResolveRequestItem)

	for _, req := range requests {
		for _, app := range req.Apps {
			programName := app.PrettyName
			if _, exists := groups[programName]; !exists {
				groups[programName] = make([]idresolver.ResolveRequestItem, 0)
			}

			// Create a new request item for this specific program
			programReq := idresolver.ResolveRequestItem{
				InternalID: req.InternalID,
				Apps:       []idresolver.MSUAppDetails{app},
			}
			groups[programName] = append(groups[programName], programReq)
		}
	}

	return groups
}

// processBVI handles BVI matching with strict positional mapping
func (me *MatchingEngine) processBVI(
	requests []idresolver.ResolveRequestItem,
	programData map[string][]GosuslugiApplicant,
	programName string,
) {
	// Extract BVI applications from MSU requests
	bviRequests := me.extractBVIRequests(requests)
	if len(bviRequests) == 0 {
		return
	}

	// Extract BVI candidates from Gosuslugi data (withoutTests == true in RegularBVI)
	bviCandidates := me.extractBVICandidates(programData)
	if len(bviCandidates) == 0 {
		slog.Warn("No BVI candidates found in Gosuslugi data", "program", programName)
		return
	}

	// Sort both by rating place
	sort.Slice(bviRequests, func(i, j int) bool {
		return bviRequests[i].Apps[0].RatingPlace < bviRequests[j].Apps[0].RatingPlace
	})
	sort.Slice(bviCandidates, func(i, j int) bool {
		return bviCandidates[i].Rating < bviCandidates[j].Rating
	})

	// Positional mapping
	matches := min(len(bviRequests), len(bviCandidates))
	for i := 0; i < matches; i++ {
		internalID := bviRequests[i].InternalID
		candidate := bviCandidates[i]
		canonicalID := fmt.Sprintf("%d", candidate.IDApplication)

		match := MatchResult{
			InternalID:      internalID,
			CanonicalID:     canonicalID,
			Confidence:      1.0, // BVI matches have highest confidence
			MatchType:       "BVI",
			ProgramName:     programName,
			CompetitionType: "BVI",
		}

		me.addGlobalMatch(internalID, match)
		slog.Debug("BVI positional match",
			"internalID", internalID,
			"canonicalID", canonicalID,
			"msuRating", bviRequests[i].Apps[0].RatingPlace,
			"gosuslugiRating", candidate.Rating)
	}

	// Handle remaining unmatched BVI requests with fallback
	for i := matches; i < len(bviRequests); i++ {
		internalID := bviRequests[i].InternalID
		match := MatchResult{
			InternalID:      internalID,
			CanonicalID:     me.generateFallbackID(internalID),
			Confidence:      0.0,
			MatchType:       "Fallback",
			ProgramName:     programName,
			CompetitionType: "BVI",
		}
		me.addGlobalMatch(internalID, match)
	}

	slog.Info("BVI matching completed",
		"program", programName,
		"msuRequests", len(bviRequests),
		"gosuslugiCandidates", len(bviCandidates),
		"successfulMatches", matches)
}

// processOtherCompetitionTypes handles Regular, SpecialQuota, TargetQuota, DedicatedQuota
func (me *MatchingEngine) processOtherCompetitionTypes(
	requests []idresolver.ResolveRequestItem,
	programData map[string][]GosuslugiApplicant,
	programName string,
) {
	competitionTypes := []string{"Regular", "SpecialQuota", "TargetQuota", "DedicatedQuota"}

	for _, competitionType := range competitionTypes {
		me.processCompetitionType(requests, programData, programName, competitionType)
	}
}

// processCompetitionType processes a specific competition type
func (me *MatchingEngine) processCompetitionType(
	requests []idresolver.ResolveRequestItem,
	programData map[string][]GosuslugiApplicant,
	programName string,
	competitionType string,
) {
	// Extract requests for this competition type
	typeRequests := me.extractRequestsByCompetitionType(requests, competitionType)
	if len(typeRequests) == 0 {
		return
	}

	// Extract candidates for this competition type
	candidates := me.extractCandidatesByCompetitionType(programData, competitionType)
	if len(candidates) == 0 {
		slog.Warn("No candidates found",
			"program", programName,
			"competitionType", competitionType)
		me.handleMissingCandidates(typeRequests, programName, competitionType)
		return
	}

	// Sort both by rating place
	sort.Slice(typeRequests, func(i, j int) bool {
		return typeRequests[i].Apps[0].RatingPlace < typeRequests[j].Apps[0].RatingPlace
	})
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Rating < candidates[j].Rating
	})

	// Two-pass matching: Strong matches first, then weak matches
	me.performTwoPassMatching(typeRequests, candidates, programName, competitionType)
}

// performTwoPassMatching implements the two-pass matching algorithm
func (me *MatchingEngine) performTwoPassMatching(
	requests []idresolver.ResolveRequestItem,
	candidates []GosuslugiApplicant,
	programName string,
	competitionType string,
) {
	usedCandidates := make(map[int]bool) // Track used candidate indices
	strongMatches := make(map[string]MatchResult)

	// Pass A: Strong matching
	for _, req := range requests {
		internalID := req.InternalID

		// Skip if already matched with higher confidence
		if existing, exists := me.globalMatches[internalID]; exists && existing.Confidence >= 0.8 {
			continue
		}

		bestCandidateIdx := -1
		bestScore := 0.0

		for i, candidate := range candidates {
			if usedCandidates[i] {
				continue
			}

			// Check if this candidate can be used (mapping rules)
			if !me.canUseCanonicalID(fmt.Sprintf("%d", candidate.IDApplication), competitionType) {
				continue
			}

			score := me.calculateStrongMatchScore(req.Apps[0], candidate, requests)
			if score >= 0.8 && score > bestScore {
				bestScore = score
				bestCandidateIdx = i
			}
		}

		if bestCandidateIdx >= 0 {
			candidate := candidates[bestCandidateIdx]
			canonicalID := fmt.Sprintf("%d", candidate.IDApplication)

			match := MatchResult{
				InternalID:      internalID,
				CanonicalID:     canonicalID,
				Confidence:      bestScore,
				MatchType:       "Strong",
				ProgramName:     programName,
				CompetitionType: competitionType,
			}

			strongMatches[internalID] = match
			usedCandidates[bestCandidateIdx] = true
			me.addGlobalMatch(internalID, match)

			slog.Debug("Strong match found",
				"internalID", internalID,
				"canonicalID", canonicalID,
				"confidence", bestScore,
				"competitionType", competitionType)
		}
	}

	// Pass B: Weak matching for remaining unmatched requests
	me.performWeakMatching(requests, candidates, usedCandidates, programName, competitionType)
}

// performWeakMatching implements the weak matching pass
func (me *MatchingEngine) performWeakMatching(
	requests []idresolver.ResolveRequestItem,
	candidates []GosuslugiApplicant,
	usedCandidates map[int]bool,
	programName string,
	competitionType string,
) {
	// Collect unmatched requests
	var unmatchedRequests []idresolver.ResolveRequestItem
	for _, req := range requests {
		if _, exists := me.globalMatches[req.InternalID]; !exists {
			unmatchedRequests = append(unmatchedRequests, req)
		}
	}

	if len(unmatchedRequests) == 0 {
		return
	}

	// Sort unmatched by rating place (ascending)
	sort.Slice(unmatchedRequests, func(i, j int) bool {
		return unmatchedRequests[i].Apps[0].RatingPlace < unmatchedRequests[j].Apps[0].RatingPlace
	})

	// Collect available candidates (not used and respect mapping rules)
	var availableCandidates []struct {
		candidate GosuslugiApplicant
		index     int
	}

	for i, candidate := range candidates {
		if usedCandidates[i] {
			continue
		}
		if !me.canUseCanonicalID(fmt.Sprintf("%d", candidate.IDApplication), competitionType) {
			continue
		}
		availableCandidates = append(availableCandidates, struct {
			candidate GosuslugiApplicant
			index     int
		}{candidate, i})
	}

	// Sort available candidates by rating (preserve order)
	sort.Slice(availableCandidates, func(i, j int) bool {
		return availableCandidates[i].candidate.Rating < availableCandidates[j].candidate.Rating
	})

	// Weak matching: maintain rating order
	matches := min(len(unmatchedRequests), len(availableCandidates))
	for i := 0; i < matches; i++ {
		req := unmatchedRequests[i]
		candidateInfo := availableCandidates[i]
		candidate := candidateInfo.candidate

		// Apply weak matching heuristics for tie-breaking
		confidence := me.calculateWeakMatchConfidence(req.Apps[0], candidate)
		canonicalID := fmt.Sprintf("%d", candidate.IDApplication)

		match := MatchResult{
			InternalID:      req.InternalID,
			CanonicalID:     canonicalID,
			Confidence:      confidence,
			MatchType:       "Weak",
			ProgramName:     programName,
			CompetitionType: competitionType,
		}

		usedCandidates[candidateInfo.index] = true
		me.addGlobalMatch(req.InternalID, match)

		slog.Debug("Weak match found",
			"internalID", req.InternalID,
			"canonicalID", canonicalID,
			"confidence", confidence,
			"competitionType", competitionType)
	}

	// Fallback for remaining unmatched
	for i := matches; i < len(unmatchedRequests); i++ {
		req := unmatchedRequests[i]
		match := MatchResult{
			InternalID:      req.InternalID,
			CanonicalID:     me.generateFallbackID(req.InternalID),
			Confidence:      0.0,
			MatchType:       "Fallback",
			ProgramName:     programName,
			CompetitionType: competitionType,
		}
		me.addGlobalMatch(req.InternalID, match)
	}
}

// calculateStrongMatchScore calculates the score for strong matching
func (me *MatchingEngine) calculateStrongMatchScore(
	msuApp idresolver.MSUAppDetails,
	candidate GosuslugiApplicant,
	allRequests []idresolver.ResolveRequestItem,
) float64 {
	score := 0.0

	// 1. Check priorities match (with holes) - +0.5
	if me.prioritiesMatch(msuApp, candidate, allRequests) {
		score += 0.5
	}

	// 2. Check EGE scores identical (order-agnostic) - +0.3
	if me.egeScoresMatch(msuApp.EGEScores, candidate) {
		score += 0.3
	}

	// 3. Check same priority for current program - +0.2
	if msuApp.Priority == candidate.Priority {
		score += 0.2
	}

	return score
}

// calculateWeakMatchConfidence calculates confidence for weak matches
func (me *MatchingEngine) calculateWeakMatchConfidence(
	msuApp idresolver.MSUAppDetails,
	candidate GosuslugiApplicant,
) float64 {
	confidence := 0.4 // Base weak confidence

	// Prefer candidates with identical EGE scores
	if me.egeScoresMatch(msuApp.EGEScores, candidate) {
		confidence += 0.1
	}

	// Prefer candidates with same priority for this program
	if msuApp.Priority == candidate.Priority {
		confidence += 0.1
	}

	return confidence
}

// Helper functions

func (me *MatchingEngine) extractBVIRequests(requests []idresolver.ResolveRequestItem) []idresolver.ResolveRequestItem {
	var bviRequests []idresolver.ResolveRequestItem
	for _, req := range requests {
		for _, app := range req.Apps {
			// Determine if this is a BVI application based on the application data
			// This would need to be determined from the competition type or other indicators
			// For now, we'll assume it's marked in the data somehow
			if strings.Contains(strings.ToLower(app.PrettyName), "bvi") || app.DVIScore > 0 {
				bviRequests = append(bviRequests, idresolver.ResolveRequestItem{
					InternalID: req.InternalID,
					Apps:       []idresolver.MSUAppDetails{app},
				})
				break // One BVI app per student
			}
		}
	}
	return bviRequests
}

func (me *MatchingEngine) extractBVICandidates(programData map[string][]GosuslugiApplicant) []GosuslugiApplicant {
	var bviCandidates []GosuslugiApplicant

	// BVI candidates are in RegularBVI list with WithoutTests == true
	if regularBVI, exists := programData["RegularBVI"]; exists {
		for _, candidate := range regularBVI {
			if candidate.WithoutTests {
				bviCandidates = append(bviCandidates, candidate)
			}
		}
	}

	return bviCandidates
}

func (me *MatchingEngine) extractRequestsByCompetitionType(
	requests []idresolver.ResolveRequestItem,
	competitionType string,
) []idresolver.ResolveRequestItem {
	var typeRequests []idresolver.ResolveRequestItem

	for _, req := range requests {
		for _, app := range req.Apps {
			// Match competition type - this logic would need to be refined based on actual data
			appCompetitionType := me.determineCompetitionType(app)
			if appCompetitionType == competitionType {
				typeRequests = append(typeRequests, idresolver.ResolveRequestItem{
					InternalID: req.InternalID,
					Apps:       []idresolver.MSUAppDetails{app},
				})
				break
			}
		}
	}

	return typeRequests
}

func (me *MatchingEngine) extractCandidatesByCompetitionType(
	programData map[string][]GosuslugiApplicant,
	competitionType string,
) []GosuslugiApplicant {
	var candidates []GosuslugiApplicant

	switch competitionType {
	case "Regular":
		// Regular candidates are in RegularBVI list with WithoutTests == false
		if regularBVI, exists := programData["RegularBVI"]; exists {
			for _, candidate := range regularBVI {
				if !candidate.WithoutTests {
					candidates = append(candidates, candidate)
				}
			}
		}
	case "DedicatedQuota":
		if dedicated, exists := programData["DedicatedQuota"]; exists {
			candidates = dedicated
		}
	case "SpecialQuota":
		if special, exists := programData["SpecialQuota"]; exists {
			candidates = special
		}
	case "TargetQuota":
		if target, exists := programData["TargetQuota"]; exists {
			candidates = target
		}
	}

	return candidates
}

func (me *MatchingEngine) determineCompetitionType(app idresolver.MSUAppDetails) string {
	// This would need to be implemented based on how competition type is encoded in the app data
	// For now, return a default
	return "Regular"
}

func (me *MatchingEngine) prioritiesMatch(
	msuApp idresolver.MSUAppDetails,
	candidate GosuslugiApplicant,
	allRequests []idresolver.ResolveRequestItem,
) bool {
	// This would need access to all programs/priorities for both MSU and Gosuslugi data
	// Simplified implementation for now
	return msuApp.Priority == candidate.Priority
}

func (me *MatchingEngine) egeScoresMatch(msuScores []int, candidate GosuslugiApplicant) bool {
	candidateScores := me.extractEGEScores(candidate)

	if len(msuScores) != len(candidateScores) {
		return false
	}

	// Sort both for order-agnostic comparison
	msuSorted := make([]int, len(msuScores))
	copy(msuSorted, msuScores)
	sort.Ints(msuSorted)

	candidateSorted := make([]int, len(candidateScores))
	copy(candidateSorted, candidateScores)
	sort.Ints(candidateSorted)

	for i := range msuSorted {
		if msuSorted[i] != candidateSorted[i] {
			return false
		}
	}

	return true
}

func (me *MatchingEngine) extractEGEScores(candidate GosuslugiApplicant) []int {
	var scores []int
	results := []float64{
		candidate.Result1, candidate.Result2, candidate.Result3, candidate.Result4,
		candidate.Result5, candidate.Result6, candidate.Result7, candidate.Result8,
	}

	for _, result := range results {
		if result > 0 {
			scores = append(scores, int(result))
		}
	}

	return scores
}

func (me *MatchingEngine) canUseCanonicalID(canonicalID string, competitionType string) bool {
	usedBy, exists := me.usedCanonicalIDs[canonicalID]
	if !exists {
		return true // Not used yet
	}

	// Check mapping rules:
	// - One Regular ID per canonical
	// - One quota ID per canonical (except DedicatedQuota)
	// - Multiple DedicatedQuota IDs per canonical allowed

	if competitionType == "DedicatedQuota" {
		return true // DedicatedQuota allows many-to-one
	}

	// For other types, check if canonical is already used by same type
	// This would need more sophisticated tracking of competition types per canonical ID
	return len(usedBy) == 0
}

func (me *MatchingEngine) addGlobalMatch(internalID string, match MatchResult) {
	// Check if this is a better match than existing
	if existing, exists := me.globalMatches[internalID]; exists {
		if existing.Confidence >= match.Confidence {
			return // Keep existing higher confidence match
		}
	}

	me.globalMatches[internalID] = match

	// Track canonical ID usage
	canonicalID := match.CanonicalID
	if !strings.HasPrefix(canonicalID, "MSU-") { // Not a fallback ID
		me.usedCanonicalIDs[canonicalID] = append(me.usedCanonicalIDs[canonicalID], internalID)
	}
}

func (me *MatchingEngine) generateFallbackID(internalID string) string {
	return fmt.Sprintf("MSU-%s", internalID)
}

func (me *MatchingEngine) handleMissingProgramData(requests []idresolver.ResolveRequestItem, programName string) {
	for _, req := range requests {
		match := MatchResult{
			InternalID:      req.InternalID,
			CanonicalID:     me.generateFallbackID(req.InternalID),
			Confidence:      0.0,
			MatchType:       "Fallback",
			ProgramName:     programName,
			CompetitionType: "Unknown",
		}
		me.addGlobalMatch(req.InternalID, match)
	}
}

func (me *MatchingEngine) handleMissingCandidates(requests []idresolver.ResolveRequestItem, programName, competitionType string) {
	for _, req := range requests {
		match := MatchResult{
			InternalID:      req.InternalID,
			CanonicalID:     me.generateFallbackID(req.InternalID),
			Confidence:      0.0,
			MatchType:       "Fallback",
			ProgramName:     programName,
			CompetitionType: competitionType,
		}
		me.addGlobalMatch(req.InternalID, match)
	}
}

func (me *MatchingEngine) finalizeResults(requests []idresolver.ResolveRequestItem) []MatchResult {
	results := make([]MatchResult, 0, len(requests))

	for _, req := range requests {
		if match, exists := me.globalMatches[req.InternalID]; exists {
			results = append(results, match)
		} else {
			// Create fallback for any unmatched
			fallback := MatchResult{
				InternalID:      req.InternalID,
				CanonicalID:     me.generateFallbackID(req.InternalID),
				Confidence:      0.0,
				MatchType:       "Fallback",
				ProgramName:     "Unknown",
				CompetitionType: "Unknown",
			}
			results = append(results, fallback)
		}
	}

	return results
}

func (me *MatchingEngine) countByConfidence(results []MatchResult, threshold float64) int {
	count := 0
	for _, result := range results {
		if result.Confidence >= threshold {
			count++
		}
	}
	return count
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
