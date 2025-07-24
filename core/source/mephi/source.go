package mephi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/utils"
	"golang.org/x/net/html"
)

// HTTPHeadingSource represents a MEPhI heading source with direct URLs to competition lists.
type HTTPHeadingSource struct {
	// HeadingName is the name of the academic program/heading
	HeadingName string `json:"heading_name"`

	// Capacities contains the capacities for different competition types
	Capacities core.Capacities `json:"capacities"`

	// RegularURLs contains URLs for regular competition lists
	RegularURLs []string `json:"regular_urls,omitempty"`

	// TargetQuotaURLs contains URLs for target quota competition lists
	TargetQuotaURLs []string `json:"target_quota_urls,omitempty"`

	// DedicatedQuotaURLs contains URLs for dedicated quota competition lists
	DedicatedQuotaURLs []string `json:"dedicated_quota_urls,omitempty"`

	// SpecialQuotaURLs contains URLs for special quota competition lists
	SpecialQuotaURLs []string `json:"special_quota_urls,omitempty"`
}

// LoadTo implements the source.HeadingSource interface for MEPhI HTTP heading sources.
func (s *HTTPHeadingSource) LoadTo(receiver source.DataReceiver) error {
	// Generate heading code from the pretty name
	headingCode := utils.GenerateHeadingCode(s.HeadingName)

	// Send HeadingData to the receiver
	receiver.PutHeadingData(&source.HeadingData{
		Code:       headingCode,
		Capacities: s.Capacities,
		PrettyName: s.HeadingName,
	})

	// Load applications from all competition types
	competitionURLs := map[core.Competition][]string{
		core.CompetitionRegular:        s.RegularURLs,
		core.CompetitionTargetQuota:    s.TargetQuotaURLs,
		core.CompetitionDedicatedQuota: s.DedicatedQuotaURLs,
		core.CompetitionSpecialQuota:   s.SpecialQuotaURLs,
	}

	for competitionType, urls := range competitionURLs {
		for _, url := range urls {
			if err := s.loadApplicationsFromURL(receiver, headingCode, url, competitionType); err != nil {
				return fmt.Errorf("failed to load applications from %s: %w", url, err)
			}
		}
	}

	return nil
}

// loadApplicationsFromURL fetches and parses applications from a single URL.
func (s *HTTPHeadingSource) loadApplicationsFromURL(receiver source.DataReceiver, headingCode string, url string, competitionType core.Competition) error {
	// Fetch the HTML content
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch URL %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error %d for URL %s", resp.StatusCode, url)
	}

	// Parse HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to parse HTML from %s: %w", url, err)
	}

	// Parse applications from the HTML
	applications, err := ParseApplicationList(doc, competitionType)
	if err != nil {
		return fmt.Errorf("failed to parse applications from %s: %w", url, err)
	}

	// Send applications to the receiver
	for _, app := range applications {
		app.HeadingCode = headingCode
		receiver.PutApplicationData(app)
	}

	return nil
}

// GetName returns the heading name for this source.
func (s *HTTPHeadingSource) GetName() string {
	return s.HeadingName
}

// GetCapacities returns the capacities for this source.
func (s *HTTPHeadingSource) GetCapacities() core.Capacities {
	return s.Capacities
}

// HasCompetition checks if this source has URLs for the specified competition type.
func (s *HTTPHeadingSource) HasCompetition(competitionType core.Competition) bool {
	switch competitionType {
	case core.CompetitionRegular, core.CompetitionBVI:
		return len(s.RegularURLs) > 0
	case core.CompetitionTargetQuota:
		return len(s.TargetQuotaURLs) > 0
	case core.CompetitionDedicatedQuota:
		return len(s.DedicatedQuotaURLs) > 0
	case core.CompetitionSpecialQuota:
		return len(s.SpecialQuotaURLs) > 0
	default:
		return false
	}
}

// GetURLs returns all URLs for the specified competition type.
func (s *HTTPHeadingSource) GetURLs(competitionType core.Competition) []string {
	switch competitionType {
	case core.CompetitionRegular, core.CompetitionBVI:
		return s.RegularURLs
	case core.CompetitionTargetQuota:
		return s.TargetQuotaURLs
	case core.CompetitionDedicatedQuota:
		return s.DedicatedQuotaURLs
	case core.CompetitionSpecialQuota:
		return s.SpecialQuotaURLs
	default:
		return nil
	}
}

// String returns a string representation of the source.
func (s *HTTPHeadingSource) String() string {
	var parts []string
	parts = append(parts, fmt.Sprintf("MEPhI: %s", s.HeadingName))

	if len(s.RegularURLs) > 0 {
		parts = append(parts, fmt.Sprintf("Regular: %d URLs", len(s.RegularURLs)))
	}
	if len(s.TargetQuotaURLs) > 0 {
		parts = append(parts, fmt.Sprintf("Target: %d URLs", len(s.TargetQuotaURLs)))
	}
	if len(s.DedicatedQuotaURLs) > 0 {
		parts = append(parts, fmt.Sprintf("Dedicated: %d URLs", len(s.DedicatedQuotaURLs)))
	}
	if len(s.SpecialQuotaURLs) > 0 {
		parts = append(parts, fmt.Sprintf("Special: %d URLs", len(s.SpecialQuotaURLs)))
	}

	return strings.Join(parts, ", ")
}