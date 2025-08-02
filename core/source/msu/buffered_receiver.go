package msu

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"sync"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/idresolver"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/utils"
)

// MSUReceiverFactory creates MSU buffered receivers
type MSUReceiverFactory struct {
	resolver idresolver.StudentIDResolver
}

// NewMSUReceiverFactory creates a new factory with the given resolver
func NewMSUReceiverFactory(resolver idresolver.StudentIDResolver) *MSUReceiverFactory {
	return &MSUReceiverFactory{resolver: resolver}
}

// CreateMSUReceiver implements source.MSUReceiverFactory
func (f *MSUReceiverFactory) CreateMSUReceiver(downstream source.DataReceiver) source.MSUReceiver {
	return NewMSUBufferedReceiver(downstream, f.resolver)
}

// RegisterFactory registers the MSU receiver factory globally
func RegisterFactory(resolver idresolver.StudentIDResolver) {
	factory := NewMSUReceiverFactory(resolver)
	source.SetMSUReceiverFactory(factory)
}

// MSUBufferedReceiver buffers MSU application data by internal ID and resolves canonical IDs
type MSUBufferedReceiver struct {
	downstream source.DataReceiver
	resolver   idresolver.StudentIDResolver
	buffer     map[string][]*MSUAppData
	headings   []*source.HeadingData
	appCount   int          // Track total applications received
	mutex      sync.RWMutex // Protects buffer, appCount, and headings from concurrent access
}

// MSUAppData holds the raw MSU application data before ID resolution
type MSUAppData struct {
	prettyName        string
	scoreSum          int
	ratingPlace       int
	priority          int
	competitionType   string
	originalSubmitted bool
	dviScore          int
	egeScores         []int
	rawStudentID      string // Original ID from MSU site
}

// NewMSUBufferedReceiver creates a new MSU buffered receiver
func NewMSUBufferedReceiver(downstream source.DataReceiver, resolver idresolver.StudentIDResolver) *MSUBufferedReceiver {
	return &MSUBufferedReceiver{
		downstream: downstream,
		resolver:   resolver,
		buffer:     make(map[string][]*MSUAppData),
		headings:   make([]*source.HeadingData, 0),
	}
}

// PutHeadingData forwards heading data directly to downstream
func (r *MSUBufferedReceiver) PutHeadingData(heading *source.HeadingData) {
	r.mutex.Lock()
	r.headings = append(r.headings, heading)
	r.mutex.Unlock()
	r.downstream.PutHeadingData(heading)
}

// PutApplicationData buffers application data by extracted internal ID
func (r *MSUBufferedReceiver) PutApplicationData(application *source.ApplicationData) {
	// Extract internal ID according to MSU rules from the plan
	internalID := r.extractInternalID(application.StudentID, application.CompetitionType.String())

	// Create MSU app data
	appData := &MSUAppData{
		prettyName:        application.HeadingName, // Use HeadingName for MSU-specific pretty name
		scoreSum:          application.ScoresSum,
		ratingPlace:       application.RatingPlace,
		priority:          application.Priority,
		competitionType:   application.CompetitionType.String(),
		originalSubmitted: application.OriginalSubmitted,
		dviScore:          application.DVIScore,
		egeScores:         application.EGEScores,
		rawStudentID:      application.StudentID,
	}

	// Lock to protect buffer from concurrent writes
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Buffer by internal ID
	r.buffer[internalID] = append(r.buffer[internalID], appData)
	r.appCount++

	// Log every 100 applications to track progress
	if r.appCount%100 == 0 {
		slog.Info("MSU applications buffered",
			"totalReceived", r.appCount,
			"uniqueStudents", len(r.buffer))
	}
}

// extractInternalID extracts internal ID according to MSU rules from the plan
func (r *MSUBufferedReceiver) extractInternalID(rawID string, competitionType string) string {
	// Remove all non-numeric characters
	numericOnly := regexp.MustCompile(`[^0-9]`).ReplaceAllString(rawID, "")

	// Apply MSU ID extraction rules from the plan
	if strings.Contains(competitionType, "DedicatedQuota") {
		// For DedicatedQuota: Use full value, pad with leading zeros to exactly 6 digits if shorter
		if len(numericOnly) < 6 {
			return fmt.Sprintf("%06s", numericOnly)
		}
		return numericOnly
	} else {
		// For most lists: Extract last 6 digits
		if len(numericOnly) >= 6 {
			return numericOnly[len(numericOnly)-6:]
		}
		// If less than 6 digits, pad with leading zeros
		return fmt.Sprintf("%06s", numericOnly)
	}
}

// Finalize resolves all buffered applications and forwards them to downstream
func (r *MSUBufferedReceiver) Finalize(ctx context.Context) error {
	// Lock for reading buffer
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if len(r.buffer) == 0 {
		slog.Info("MSU buffered receiver finalize: no applications to process")
		return nil
	}

	// Log buffer statistics
	totalApps := 0
	for _, apps := range r.buffer {
		totalApps += len(apps)
	}
	slog.Info("MSU buffered receiver starting finalization",
		"uniqueStudents", len(r.buffer),
		"totalApplications", totalApps)

	// Build resolve request
	request := make([]idresolver.ResolveRequestItem, 0, len(r.buffer))
	for internalID, apps := range r.buffer {
		resolveApps := make([]idresolver.MSUAppDetails, len(apps))
		for i, app := range apps {
			resolveApps[i] = idresolver.MSUAppDetails{
				PrettyName:  app.prettyName,
				ScoreSum:    app.scoreSum,
				RatingPlace: app.ratingPlace,
				Priority:    app.priority,
				DVIScore:    app.dviScore,
				EGEScores:   app.egeScores,
			}
		}
		request = append(request, idresolver.ResolveRequestItem{
			InternalID: internalID,
			Apps:       resolveApps,
		})
	}

	// Call resolver
	slog.Info("Requesting ID resolution for MSU students", "requestCount", len(request))
	response, err := r.resolver.ResolveBatch(ctx, request)
	if err != nil {
		slog.Error("Failed to resolve MSU student IDs", "error", err, "requestCount", len(request))
		// Fallback: use hashed internal IDs
		response = r.createFallbackResponse(request)
		slog.Warn("Using fallback ID resolution for all MSU students", "fallbackCount", len(response))
	} else {
		slog.Info("Received ID resolution response", "responseCount", len(response))
	}

	// Create response map for quick lookup
	resolutionMap := make(map[string]idresolver.ResolveResponseItem)
	for _, resp := range response {
		resolutionMap[resp.InternalID] = resp
	}

	// Forward resolved applications to downstream
	for internalID, apps := range r.buffer {
		resolution, found := resolutionMap[internalID]
		if !found {
			// Create fallback resolution
			resolution = idresolver.ResolveResponseItem{
				InternalID:  internalID,
				CanonicalID: r.createFallbackID(internalID),
				Confidence:  0.0,
			}
		}

		// Log resolution result
		if resolution.Confidence < 0.7 {
			// slog.Warn("Low confidence MSU ID resolution",
			// 	"internalID", internalID,
			// 	"canonicalID", resolution.CanonicalID,
			// 	"confidence", resolution.Confidence,
			// 	"appCount", len(apps))
		} else {
			slog.Debug("MSU ID resolved",
				"internalID", internalID,
				"canonicalID", resolution.CanonicalID,
				"confidence", resolution.Confidence)
		}

		// Forward each application with resolved canonical ID
		for _, app := range apps {
			applicationData := &source.ApplicationData{
				HeadingCode:       utils.GenerateHeadingCode(app.prettyName),
				StudentID:         resolution.CanonicalID,
				ScoresSum:         app.scoreSum,
				RatingPlace:       app.ratingPlace,
				Priority:          app.priority,
				CompetitionType:   parseCompetitionType(app.competitionType),
				OriginalSubmitted: app.originalSubmitted,
				DVIScore:          app.dviScore,
				EGEScores:         app.egeScores,
				HeadingName:       app.prettyName,
			}
			r.downstream.PutApplicationData(applicationData)
		}
	}

	// Log processing summary
	highConfidenceCount := 0
	prettyHighConfidenceCount := 0
	mediumConfidenceCount := 0
	lowConfidenceCount := 0
	fallbackCount := 0
	processedApps := 0

	for _, resolution := range resolutionMap {
		switch {
		case resolution.Confidence >= 0.8:
			highConfidenceCount++
		case resolution.Confidence >= 0.6:
			prettyHighConfidenceCount++
		case resolution.Confidence >= 0.4:
			mediumConfidenceCount++
		case resolution.Confidence > 0.0:
			lowConfidenceCount++
		default:
			fallbackCount++
		}
	}

	for _, apps := range r.buffer {
		processedApps += len(apps)
	}

	slog.Info("MSU ID resolution statistics",
		"highConfidence", highConfidenceCount,
		"prettyHighConfidence", prettyHighConfidenceCount,
		"mediumConfidence", mediumConfidenceCount,
		"lowConfidence", lowConfidenceCount,
		"fallbackIDs", fallbackCount)

	slog.Info("MSU buffered receiver finalization complete",
		"processedStudents", len(r.buffer),
		"processedApplications", processedApps)

	return nil
}

// createFallbackResponse creates fallback responses when resolver fails
func (r *MSUBufferedReceiver) createFallbackResponse(request []idresolver.ResolveRequestItem) []idresolver.ResolveResponseItem {
	response := make([]idresolver.ResolveResponseItem, len(request))
	for i, req := range request {
		response[i] = idresolver.ResolveResponseItem{
			InternalID:  req.InternalID,
			CanonicalID: r.createFallbackID(req.InternalID),
			Confidence:  0.0,
		}
	}
	return response
}

// createFallbackID creates a fallback canonical ID from internal ID
func (r *MSUBufferedReceiver) createFallbackID(internalID string) string {
	// Create fallback ID: MSU- prefix + padded internal ID
	fallbackID := fmt.Sprintf("MSU-%s", internalID)
	// Pad to 13 characters as required by utils.PrepareStudentID
	for len(fallbackID) < 13 {
		fallbackID = "0" + fallbackID
	}
	// If still too long, truncate from the beginning but keep MSU prefix recognizable
	if len(fallbackID) > 13 {
		// Take last 13 characters, which should preserve most of the internal ID
		fallbackID = fallbackID[len(fallbackID)-13:]
	}
	return fallbackID
}

// parseCompetitionType converts string back to Competition enum
func parseCompetitionType(competitionStr string) core.Competition {
	switch competitionStr {
	case "Regular":
		return core.CompetitionRegular
	case "BVI":
		return core.CompetitionBVI
	case "TargetQuota":
		return core.CompetitionTargetQuota
	case "DedicatedQuota":
		return core.CompetitionDedicatedQuota
	case "SpecialQuota":
		return core.CompetitionSpecialQuota
	default:
		slog.Warn("Unknown competition type, defaulting to Regular", "competitionStr", competitionStr)
		return core.CompetitionRegular
	}
}
