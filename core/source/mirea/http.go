package mirea

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/utils"
)

// HTTPHeadingSource loads MIREA heading data from multiple JSON list IDs.
type HTTPHeadingSource struct {
	RegularListIDs        []string
	BVIListIDs            []string
	TargetQuotaListIDs    []string
	DedicatedQuotaListIDs []string
	SpecialQuotaListIDs   []string
}

// fetchMireaListByID fetches and decodes the MIREA list from a list ID.
func fetchMireaListByID(listID string) (*MireaListResponse, error) {
	if listID == "" {
		return nil, nil
	}
	url := "https://priem.mirea.ru/competitions_api/entrants?competitions[]=" + listID
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	release, err := source.AcquireHTTPSemaphores(ctx, "mirea")
	if err != nil {
		return nil, fmt.Errorf("failed to acquire semaphores for %s: %w", url, err)
	}
	defer release()

	// Create request with proper headers to avoid 403 errors
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s: %w", url, err)
	}

	// Add headers to mimic a real browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en;q=0.8")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", "https://priem.mirea.ru/")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download %s (status code %d)", url, resp.StatusCode)
	}

	var mireaResp MireaListResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&mireaResp); err != nil {
		return nil, fmt.Errorf("failed to decode MIREA list JSON: %w", err)
	}

	// Add 5 second delay to work around DDoS protection
	time.Sleep(5 * time.Second)

	return &mireaResp, nil
}

// LoadTo implements source.HeadingSource for HTTPHeadingSource.
func (s *HTTPHeadingSource) LoadTo(receiver source.DataReceiver) error {
	// Find the first available list to extract metadata
	var firstResponse *MireaListResponse
	var prettyName string
	var regularCapacity, targetQuotaCapacity, dedicatedQuotaCapacity, specialQuotaCapacity int

	// Try to get metadata from the first available list
	allListIDs := [][]string{
		s.RegularListIDs,
		s.BVIListIDs,
		s.TargetQuotaListIDs,
		s.DedicatedQuotaListIDs,
		s.SpecialQuotaListIDs,
	}

	for _, listGroup := range allListIDs {
		for _, listID := range listGroup {
			if listID != "" {
				resp, err := fetchMireaListByID(listID)
				if err != nil {
					log.Printf("Error fetching metadata from list %s: %v", listID, err)
					continue
				}
				if resp != nil && len(resp.Data) > 0 {
					firstResponse = resp
					prettyName = extractHeadingName(resp.Data[0].Title)
					break
				}
			}
		}
		if firstResponse != nil {
			break
		}
	}

	if firstResponse == nil || prettyName == "" {
		return fmt.Errorf("could not extract heading name from any available list")
	}

	// Calculate capacities based on MIREA-specific logic
	// For Regular/BVI: use the same capacity value (don't sum)
	if len(s.RegularListIDs) > 0 && s.RegularListIDs[0] != "" {
		resp, err := fetchMireaListByID(s.RegularListIDs[0])
		if err == nil && resp != nil && len(resp.Data) > 0 {
			regularCapacity = resp.Data[0].Plan
		}
	}

	// For TargetQuota: sum individual list capacities
	for _, listID := range s.TargetQuotaListIDs {
		if listID != "" {
			resp, err := fetchMireaListByID(listID)
			if err == nil && resp != nil && len(resp.Data) > 0 {
				targetQuotaCapacity += resp.Data[0].Plan
			}
		}
	}

	// For other quota types: sum individual list capacities
	for _, listID := range s.DedicatedQuotaListIDs {
		if listID != "" {
			resp, err := fetchMireaListByID(listID)
			if err == nil && resp != nil && len(resp.Data) > 0 {
				dedicatedQuotaCapacity += resp.Data[0].Plan
			}
		}
	}

	for _, listID := range s.SpecialQuotaListIDs {
		if listID != "" {
			resp, err := fetchMireaListByID(listID)
			if err == nil && resp != nil && len(resp.Data) > 0 {
				specialQuotaCapacity += resp.Data[0].Plan
			}
		}
	}

	headingCode := utils.GenerateHeadingCode(prettyName)
	receiver.PutHeadingData(&source.HeadingData{
		Code: headingCode,
		Capacities: core.Capacities{
			Regular:        regularCapacity, // Same for both Regular and BVI
			TargetQuota:    targetQuotaCapacity,
			DedicatedQuota: dedicatedQuotaCapacity,
			SpecialQuota:   specialQuotaCapacity,
		},
		PrettyName: prettyName,
	})

	// Process all competition type lists
	listDefs := []struct {
		ListIDs     []string
		Competition core.Competition
		ListName    string
	}{
		{s.RegularListIDs, core.CompetitionRegular, "Regular Lists"},
		{s.BVIListIDs, core.CompetitionBVI, "BVI Lists"},
		{s.TargetQuotaListIDs, core.CompetitionTargetQuota, "Target Quota Lists"},
		{s.DedicatedQuotaListIDs, core.CompetitionDedicatedQuota, "Dedicated Quota Lists"},
		{s.SpecialQuotaListIDs, core.CompetitionSpecialQuota, "Special Quota Lists"},
	}

	for _, def := range listDefs {
		for i, listID := range def.ListIDs {
			if listID == "" {
				continue
			}
			resp, err := fetchMireaListByID(listID)
			if err != nil {
				log.Printf("Error fetching %s %d (%s): %v", def.ListName, i+1, listID, err)
				continue
			}
			if resp == nil || len(resp.Data) == 0 {
				log.Printf("No data found in %s %d (%s)", def.ListName, i+1, listID)
				continue
			}
			parseAndLoadApplications(resp.Data[0].Entrants, def.Competition, headingCode, receiver)
		}
	}

	return nil
}
