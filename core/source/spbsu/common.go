package spbsu

import (
	"analabit/core"
	"analabit/core/source"
	"encoding/json"
	"fmt"
	"io"
)

// SpbsuApplicationEntry represents a single application entry in the SPbSU JSON list.
type SpbsuApplicationEntry struct {
	UserCode           string      `json:"user_code"`
	ScoreOverall       int         `json:"score_overall"`
	OrderNumber        int         `json:"order_number"`
	PriorityNumber     int         `json:"priority_number"`
	OriginalDocument   bool        `json:"original_document"`
	TargetOrganization interface{} `json:"target_organization"`
	WithoutTrials      bool        `json:"without_trials"`
}

// SpbsuListResponse represents the top-level JSON structure.
type SpbsuListResponse struct {
	List []SpbsuApplicationEntry `json:"list"`
}

// FetchListFunc abstracts fetching and decoding a list from a source (e.g., URL or file).
type FetchListFunc func(source string) ([]SpbsuApplicationEntry, error)

// parseAndLoadApplications parses entries and emits ApplicationData to the receiver.
func parseAndLoadApplications(entries []SpbsuApplicationEntry, competitionType core.Competition, headingCode string, receiver source.DataReceiver) {
	for _, entry := range entries {
		competition := competitionType

		if entry.WithoutTrials && competitionType == core.CompetitionRegular {
			competition = core.CompetitionBVI
		}

		receiver.PutApplicationData(&source.ApplicationData{
			HeadingCode:       headingCode,
			StudentID:         entry.UserCode,
			ScoresSum:         entry.ScoreOverall,
			RatingPlace:       entry.OrderNumber,
			Priority:          entry.PriorityNumber,
			CompetitionType:   competition,
			OriginalSubmitted: entry.OriginalDocument,
		})
	}
}

// decodeSpbsuList decodes the JSON from an io.Reader into a slice of entries.
func decodeSpbsuList(r io.Reader) ([]SpbsuApplicationEntry, error) {
	var resp SpbsuListResponse
	dec := json.NewDecoder(r)
	if err := dec.Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode SPbSU list JSON: %w", err)
	}
	return resp.List, nil
}
