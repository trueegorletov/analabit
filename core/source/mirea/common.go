// Package mirea implements MIREA-specific data sources and parsers for the analabit system.
package mirea

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
)

// MireaApplicationEntry represents a single application entry in the MIREA JSON list.
type MireaApplicationEntry struct {
	Spn string `json:"spn"`           // Student ID
	Fm  int    `json:"fm"`            // Scores sum
	P   int    `json:"p"`             // Priority
	Nd  string `json:"nd,omitempty"`  // Not used in logics
	S   string `json:"s"`             // Status
	Acc int    `json:"acc,omitempty"` // Acceptance status
}

// MireaListMetadata represents metadata for a single list.
type MireaListMetadata struct {
	ID       string                  `json:"id"`       // List ID
	Plan     int                     `json:"plan"`     // Capacity
	Title    string                  `json:"title"`    // List title
	Entrants []MireaApplicationEntry `json:"entrants"` // Application entries
}

// MireaListResponse represents the top-level JSON structure.
type MireaListResponse struct {
	Data []MireaListMetadata `json:"data"` // List metadata and entrants
}

// FetchListFunc abstracts fetching and decoding a list from a source (e.g., URL or file).
type FetchListFunc func(source string) ([]MireaApplicationEntry, error)

// parseAndLoadApplications parses entries and emits ApplicationData to the receiver.
func parseAndLoadApplications(entries []MireaApplicationEntry, competitionType core.Competition, headingCode string, receiver source.DataReceiver) {
	for i, entry := range entries {

		receiver.PutApplicationData(&source.ApplicationData{
			HeadingCode:       headingCode,
			StudentID:         entry.Spn,
			ScoresSum:         entry.Fm,
			RatingPlace:       i + 1,
			Priority:          entry.P,
			CompetitionType:   competitionType,
			OriginalSubmitted: entry.Acc > 0,
		})
	}
}

// decodeMireaList decodes the JSON from an io.Reader into a slice of entries.
func decodeMireaList(r io.Reader) ([]MireaApplicationEntry, error) {
	var resp MireaListResponse
	dec := json.NewDecoder(r)
	if err := dec.Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode MIREA list JSON: %w", err)
	}

	if len(resp.Data) == 0 {
		return []MireaApplicationEntry{}, nil
	}

	return resp.Data[0].Entrants, nil
}

// extractHeadingName extracts clean heading names by cutting at the first slash character.
func extractHeadingName(title string) string {
	if slashIndex := strings.Index(title, "/"); slashIndex != -1 {
		return strings.TrimSpace(title[:slashIndex])
	}
	return strings.TrimSpace(title)
}

// validateEntry performs basic validation on a MIREA application entry.
func validateEntry(entry MireaApplicationEntry) bool {
	return entry.Spn != "" && entry.P > 0
}
