package spbsu

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
)

// SpbsuApplicationEntry represents a single application entry in the SPbSU JSON list.
type SpbsuApplicationEntry struct {
	UserCode           string      `json:"user_code"`
	ScoreOverall       int         `json:"score_overall"`
	OrderNumber        int         `json:"order_number"`
	PriorityNumber     int         `json:"priority_number"`
	AdmissionAgreement bool        `json:"admission_agreement"`
	TargetOrganization interface{} `json:"target_organization"`
	WithoutTrials      bool        `json:"without_trials"`
	SSPVOStatus        string      `json:"sspvo_statu"`
}

// SpbsuMeta represents pagination metadata from the SPbSU API response.
type SpbsuMeta struct {
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	Total       int `json:"total"`
	PerPage     int `json:"per_page"`
}

// SpbsuListResponse represents the top-level JSON structure.
type SpbsuListResponse struct {
	List []SpbsuApplicationEntry `json:"data"`
	Meta SpbsuMeta               `json:"meta"`
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

		if strings.Contains(strings.ToLower(entry.SSPVOStatus), "отозвано") {
			continue
		}

		receiver.PutApplicationData(&source.ApplicationData{
			HeadingCode:       headingCode,
			StudentID:         entry.UserCode,
			ScoresSum:         entry.ScoreOverall,
			RatingPlace:       entry.OrderNumber,
			Priority:          entry.PriorityNumber,
			CompetitionType:   competition,
			OriginalSubmitted: entry.AdmissionAgreement,
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
