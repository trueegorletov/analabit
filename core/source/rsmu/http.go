// Package rsmu provides support for loading RSMU (Russian National Research Medical University) admission data.
// RSMU admission lists are provided in JSON format with target quota lists and individual applicant data.
package rsmu

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/utils"
)

// HTTPHeadingSource defines how to load RSMU heading data from JSON URLs.
// RSMU provides admission lists in JSON format with target quota lists and individual applicant data.
type HTTPHeadingSource struct {
	ProgramName            string
	TargetQuotaListURLs    []string
	RegularListURL         string
	SpecialQuotaListURL    string
	DedicatedQuotaListURL  string
}

// LoadTo loads data from HTTP source, downloading JSON files and sending HeadingData and ApplicationData to the provided receiver.
func (s *HTTPHeadingSource) LoadTo(receiver source.DataReceiver) error {
	// Acquire a semaphore slot, respecting context cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	release, err := source.AcquireHTTPSemaphores(ctx, "rsmu")
	if err != nil {
		return fmt.Errorf("failed to acquire semaphores for RSMU: %w", err)
	}
	defer release()

	log.Printf("Processing RSMU admission data for program: %s", s.ProgramName)

	// Collect all individual lists and calculate total capacities
	var allLists []*IndividualList
	capacities := &core.Capacities{}

	// Process target quota lists
	for _, url := range s.TargetQuotaListURLs {
		list, err := s.downloadIndividualList(ctx, url)
		if err != nil {
			log.Printf("Warning: failed to download target quota list from %s: %v", url, err)
			continue
		}
		allLists = append(allLists, list)
		capacities.TargetQuota += list.Plan
	}

	// Process regular list
	if s.RegularListURL != "" {
		list, err := s.downloadIndividualList(ctx, s.RegularListURL)
		if err != nil {
			log.Printf("Warning: failed to download regular list from %s: %v", s.RegularListURL, err)
		} else {
			allLists = append(allLists, list)
			capacities.Regular = list.Plan
		}
	}

	// Process special quota list
	if s.SpecialQuotaListURL != "" {
		list, err := s.downloadIndividualList(ctx, s.SpecialQuotaListURL)
		if err != nil {
			log.Printf("Warning: failed to download special quota list from %s: %v", s.SpecialQuotaListURL, err)
		} else {
			allLists = append(allLists, list)
			capacities.SpecialQuota = list.Plan
		}
	}

	// Process dedicated quota list
	if s.DedicatedQuotaListURL != "" {
		list, err := s.downloadIndividualList(ctx, s.DedicatedQuotaListURL)
		if err != nil {
			log.Printf("Warning: failed to download dedicated quota list from %s: %v", s.DedicatedQuotaListURL, err)
		} else {
			allLists = append(allLists, list)
			capacities.DedicatedQuota = list.Plan
		}
	}

	if len(allLists) == 0 {
		return fmt.Errorf("no valid lists found for RSMU program: %s", s.ProgramName)
	}

	headingCode := utils.GenerateHeadingCode(s.ProgramName)

	// Send HeadingData to the receiver
	receiver.PutHeadingData(&source.HeadingData{
		Code:       headingCode,
		Capacities: *capacities,
		PrettyName: s.ProgramName,
	})

	log.Printf("Sent RSMU heading: %s (Code: %s, Caps: %v)", s.ProgramName, headingCode, *capacities)

	// Send ApplicationData for each applicant from all lists
	totalApplicants := 0
	for _, list := range allLists {
		for _, applicant := range list.Applicants {
			competitionType := mapCompetitionType(list.Type)
			// Override to BVI if noExam is true
			if applicant.NoExam {
				competitionType = core.CompetitionBVI
			}

			appData := &source.ApplicationData{
				HeadingCode:       headingCode,
				StudentID:         applicant.Title,
				ScoresSum:         applicant.Total,
				RatingPlace:       applicant.Order,
				Priority:          applicant.Priority,
				CompetitionType:   competitionType,
				OriginalSubmitted: applicant.Original,
			}
			receiver.PutApplicationData(appData)
			totalApplicants++
		}
	}

	log.Printf("Sent %d applications for RSMU heading %s", totalApplicants, s.ProgramName)
	return nil
}



// downloadIndividualList downloads and parses an individual list JSON
func (s *HTTPHeadingSource) downloadIndividualList(ctx context.Context, url string) (*IndividualList, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download individual list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download individual list (status code %d)", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var individualList IndividualList
	if err := json.Unmarshal(body, &individualList); err != nil {
		return nil, fmt.Errorf("failed to parse individual list JSON: %w", err)
	}

	return &individualList, nil
}