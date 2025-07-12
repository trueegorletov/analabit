package spbstu

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/utils"
)

// HTTPHeadingSource loads SPbSTU heading data from multiple JSON list IDs.
type HTTPHeadingSource struct {
	PrettyName           string
	RegularListID        int
	TargetQuotaListIDs   []int
	DedicatedQuotaListID int
	SpecialQuotaListID   int
	Capacities           core.Capacities
}

// SpbstuCapacityResponse represents the capacity response from SPbSTU API.
type SpbstuCapacityResponse struct {
	Places int `json:"places"`
}

// fetchSpbstuListByID fetches and decodes the SPbSTU list from a list ID.
func fetchSpbstuListByID(listID int, competitionFilter int) ([]SpbstuApplicationEntry, error) {
	if listID == -1 {
		return nil, nil
	}
	url := fmt.Sprintf("https://my.spbstu.ru/home/get-abit-list?filter_1=2&filter_2=%d&filter_3=%d&education_level=bachelor", competitionFilter, listID)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	release, err := source.AcquireHTTPSemaphores(ctx, "spbstu")
	if err != nil {
		return nil, fmt.Errorf("failed to acquire semaphores for %s: %w", url, err)
	}
	defer release()

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download %s (status code %d)", url, resp.StatusCode)
	}

	return decodeSpbstuList(resp.Body)
}

// fetchSpbstuCapacity fetches capacity data for a specific list ID.
func fetchSpbstuCapacity(listID int) (int, error) {
	if listID == -1 {
		return 0, nil
	}

	url := "https://my.spbstu.ru/home/get-direction-info"
	payload := map[string]string{
		"id_3":            strconv.Itoa(listID),
		"education_level": "bachelor",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal capacity request payload: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	release, err := source.AcquireHTTPSemaphores(ctx, "spbstu")
	if err != nil {
		return 0, fmt.Errorf("failed to acquire semaphores for %s: %w", url, err)
	}
	defer release()

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return 0, fmt.Errorf("failed to fetch capacity from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to fetch capacity from %s (status code %d)", url, resp.StatusCode)
	}

	var responses []SpbstuCapacityResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&responses); err != nil {
		return 0, fmt.Errorf("failed to decode capacity response: %w", err)
	}

	if len(responses) == 0 {
		return 0, fmt.Errorf("empty capacity response")
	}

	return responses[0].Places, nil
}

// LoadTo implements source.HeadingSource for HTTPHeadingSource.
func (s *HTTPHeadingSource) LoadTo(receiver source.DataReceiver) error {
	if s.PrettyName == "" {
		return fmt.Errorf("PrettyName is required for SPbSTU HTTPHeadingSource")
	}

	headingCode := utils.GenerateHeadingCode(s.PrettyName)

	// Fetch capacities at runtime if not provided
	capacities := s.Capacities
	if capacities.Regular == 0 && s.RegularListID != -1 {
		regular, err := fetchSpbstuCapacity(s.RegularListID)
		if err != nil {
			log.Printf("Error fetching regular capacity for %s: %v", s.PrettyName, err)
		} else {
			capacities.Regular = regular
		}
	}

	if capacities.DedicatedQuota == 0 && s.DedicatedQuotaListID != -1 {
		dedicated, err := fetchSpbstuCapacity(s.DedicatedQuotaListID)
		if err != nil {
			log.Printf("Error fetching dedicated quota capacity for %s: %v", s.PrettyName, err)
		} else {
			capacities.DedicatedQuota = dedicated
		}
	}

	if capacities.SpecialQuota == 0 && s.SpecialQuotaListID != -1 {
		special, err := fetchSpbstuCapacity(s.SpecialQuotaListID)
		if err != nil {
			log.Printf("Error fetching special quota capacity for %s: %v", s.PrettyName, err)
		} else {
			capacities.SpecialQuota = special
		}
	}

	// Sum TargetQuota capacities when multiple lists exist
	if capacities.TargetQuota == 0 && len(s.TargetQuotaListIDs) > 0 {
		totalTarget := 0
		for _, listID := range s.TargetQuotaListIDs {
			if listID == -1 {
				continue
			}
			target, err := fetchSpbstuCapacity(listID)
			if err != nil {
				log.Printf("Error fetching target quota capacity for list %d in %s: %v", listID, s.PrettyName, err)
				continue
			}
			totalTarget += target
		}
		capacities.TargetQuota = totalTarget
	}

	receiver.PutHeadingData(&source.HeadingData{
		Code:       headingCode,
		Capacities: capacities,
		PrettyName: s.PrettyName,
	})

	// Define list types with their competition filters
	listDefs := []struct {
		ListID            int
		Competition       core.Competition
		CompetitionFilter int
		ListName          string
	}{
		{s.RegularListID, core.CompetitionRegular, 1, "Regular&BVI List"},
		{s.DedicatedQuotaListID, core.CompetitionDedicatedQuota, 4, "Dedicated Quota List"},
		{s.SpecialQuotaListID, core.CompetitionSpecialQuota, 3, "Special Quota List"},
	}

	for _, def := range listDefs {
		if def.ListID == -1 {
			continue
		}
		entries, err := fetchSpbstuListByID(def.ListID, def.CompetitionFilter)
		if err != nil {
			log.Printf("Error fetching %s (%d): %v", def.ListName, def.ListID, err)
			continue
		}
		if entries == nil {
			log.Printf("No entries found in %s (%d)", def.ListName, def.ListID)
			continue
		}
		parseAndLoadApplications(entries, def.Competition, headingCode, receiver)
	}

	// Handle multiple target quota list IDs
	for i, listID := range s.TargetQuotaListIDs {
		if listID == -1 {
			continue
		}
		entries, err := fetchSpbstuListByID(listID, 5) // Competition filter 5 for TargetQuota
		if err != nil {
			log.Printf("Error fetching Target Quota List %d (%d): %v", i+1, listID, err)
			continue
		}
		if entries == nil {
			log.Printf("No entries found in Target Quota List %d (%d)", i+1, listID)
			continue
		}
		parseAndLoadApplications(entries, core.CompetitionTargetQuota, headingCode, receiver)
	}

	return nil
}
