// Package spbstu provides source implementation for St. Petersburg State Technical University (SPbSTU).
package spbstu

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
)

// SpbstuApplicationEntry represents a single application entry in the SPbSTU JSON list.
type SpbstuApplicationEntry struct {
	Num      int    `json:"num"`
	Code     string `json:"code"`
	Base     string `json:"base"`
	Sum      int    `json:"sum"`
	Priority int    `json:"priority"`
	Approval string `json:"approval"`
}

// SpbstuListResponse represents the top-level JSON structure.
type SpbstuListResponse struct {
	Results []SpbstuApplicationEntry `json:"results"`
}

// FetchListFunc abstracts fetching and decoding a list from a source (e.g., URL or file).
type FetchListFunc func(source string) ([]SpbstuApplicationEntry, error)

// parseAndLoadApplications parses entries and emits ApplicationData to the receiver.
func parseAndLoadApplications(entries []SpbstuApplicationEntry, competitionType core.Competition, headingCode string, receiver source.DataReceiver) {
	for _, entry := range entries {
		competition := competitionType

		// BVI detection through HTML parsing in the Base field
		if strings.Contains(entry.Base, "Да") && competitionType == core.CompetitionRegular {
			competition = core.CompetitionBVI
		}

		// Map approval field: "+" means original submitted, other values mean not submitted
		originalSubmitted := entry.Approval == "+"

		receiver.PutApplicationData(&source.ApplicationData{
			HeadingCode:       headingCode,
			StudentID:         entry.Code,
			ScoresSum:         entry.Sum,
			RatingPlace:       entry.Num,
			Priority:          entry.Priority,
			CompetitionType:   competition,
			OriginalSubmitted: originalSubmitted,
		})
	}
}

// decodeSpbstuList decodes the JSON from an io.Reader into a slice of entries.
func decodeSpbstuList(r io.Reader) ([]SpbstuApplicationEntry, error) {
	var resp SpbstuListResponse
	dec := json.NewDecoder(r)
	if err := dec.Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode SPbSTU list JSON: %w", err)
	}
	return resp.Results, nil
}

// extractHeadingName extracts clean heading names by removing registry code prefixes
// and optionally target organization suffixes in parentheses for target quota lists.
func extractHeadingName(title string, isTargetQuota bool) string {
	// Remove DD.DD.DD pattern prefix
	codePattern := regexp.MustCompile(`^\d{2}\.\d{2}\.\d{2}\s*`)
	cleaned := codePattern.ReplaceAllString(title, "")

	// Only remove organization suffix for target quota lists (for debugging purposes)
	if isTargetQuota {
		suffixPattern := regexp.MustCompile(`\s*\([^)]*\)\s*$`)
		cleaned = suffixPattern.ReplaceAllString(cleaned, "")
	}

	return strings.TrimSpace(cleaned)
}
