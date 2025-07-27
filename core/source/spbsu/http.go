package spbsu

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/utils"
)

type HttpHeadingSource struct {
	PrettyName           string
	RegularListID        int
	TargetQuotaListIDs   []int
	DedicatedQuotaListID int
	SpecialQuotaListID   int
	Capacities           core.Capacities
}

func makeSingleRequest(url string) (*http.Response, error) {
	var attempt int
	for {
		attempt++
		if err := source.WaitBeforeHTTPRequest("spbsu", context.Background()); err != nil {
			return nil, fmt.Errorf("timeout coordination failed: %w", err)
		}
		release, err := source.AcquireHTTPSemaphores(context.Background(), "spbsu")
		if err != nil {
			return nil, fmt.Errorf("failed to acquire semaphore: %w", err)
		}
		ctx := context.Background()
		resp, err := source.GlobalSpbsuRateLimiter.MakeRequest(ctx, url)
		release()
		if err != nil {
			slog.Warn("SPbSU request failed, will retry", "url", url, "attempt", attempt, "error", err)
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}
		if resp.StatusCode == http.StatusTooManyRequests || (resp.StatusCode >= 500 && resp.StatusCode <= 599) {
			slog.Warn("SPbSU request failed with retriable status, will retry", "url", url, "attempt", attempt, "status", resp.Status)
			resp.Body.Close()
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("received non-200 status code: %s", resp.Status)
		}
		return resp, nil
	}
}

func fetchListResiliently(listID int) ([]SpbsuApplicationEntry, error) {
	if listID == -1 {
		return nil, nil
	}
	baseURL := "https://back-enrollelists.spbu.ru/api/contest-results?per-page=300&sort=&filter[competitive_group_id]=" + strconv.Itoa(listID)
	var expectedTotal int
	for attempt := 1; ; attempt++ {
		var allEntries []SpbsuApplicationEntry
		firstPageURL := baseURL + "&page=1"
		resp, err := makeSingleRequest(firstPageURL)
		if err != nil {
			slog.Warn("Failed to fetch first page, retrying", "list_id", listID, "attempt", attempt, "error", err)
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}
		var firstPageResp SpbsuListResponse
		if err := json.NewDecoder(resp.Body).Decode(&firstPageResp); err != nil {
			resp.Body.Close()
			slog.Warn("Failed to decode first page, retrying", "list_id", listID, "attempt", attempt, "error", err)
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}
		resp.Body.Close()
		expectedTotal = firstPageResp.Meta.Total
		if expectedTotal == 0 {
			return []SpbsuApplicationEntry{}, nil
		}
		allEntries = append(allEntries, firstPageResp.List...)
		totalPages := firstPageResp.Meta.LastPage
		if totalPages == 1 {
			if len(allEntries) == expectedTotal {
				return allEntries, nil
			}
			slog.Warn("First page entry count mismatch, retrying", "list_id", listID, "expected", expectedTotal, "got", len(allEntries))
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}
		for page := 2; page <= totalPages; page++ {
			var pageAttempt int
			for {
				pageAttempt++
				pageURL := baseURL + "&page=" + strconv.Itoa(page)
				resp, err := makeSingleRequest(pageURL)
				if err != nil {
					slog.Warn("Failed to fetch page, retrying page", "list_id", listID, "page", page, "attempt", pageAttempt, "error", err)
					time.Sleep(time.Duration(pageAttempt) * time.Second)
					continue
				}
				var pageResp SpbsuListResponse
				if err := json.NewDecoder(resp.Body).Decode(&pageResp); err != nil {
					resp.Body.Close()
					slog.Warn("Failed to decode page, retrying page", "list_id", listID, "page", page, "attempt", pageAttempt, "error", err)
					time.Sleep(time.Duration(pageAttempt) * time.Second)
					continue
				}
				resp.Body.Close()
				allEntries = append(allEntries, pageResp.List...)
				break
			}
		}
		if len(allEntries) == expectedTotal {
			slog.Info("Successfully fetched all entries for list", "list_id", listID, "total", expectedTotal)
			return allEntries, nil
		}
		slog.Warn("Entry count mismatch, retrying entire list", "list_id", listID, "expected", expectedTotal, "got", len(allEntries), "attempt", attempt)
		time.Sleep(time.Duration(attempt) * 5 * time.Second)
	}
}

func (s *HttpHeadingSource) LoadTo(receiver source.DataReceiver) error {
	if s.PrettyName == "" {
		return fmt.Errorf("PrettyName is required for SPbSU HttpHeadingSource")
	}
	headingCode := utils.GenerateHeadingCode(s.PrettyName)

	// Define all lists to be fetched for this heading
	listDefs := []struct {
		ListID      int
		Competition core.Competition
	}{
		{s.RegularListID, core.CompetitionRegular},
		{s.DedicatedQuotaListID, core.CompetitionDedicatedQuota},
		{s.SpecialQuotaListID, core.CompetitionSpecialQuota},
	}

	var allListsLoaded = true

	// Fetch main lists synchronously
	for _, def := range listDefs {
		if def.ListID != -1 {
			entries, err := fetchListResiliently(def.ListID)
			if err != nil {
				slog.Error("Failed to fetch list after all retries", "list_id", def.ListID, "error", err)
				allListsLoaded = false
				continue
			}
			parseAndLoadApplications(entries, def.Competition, headingCode, receiver)
		}
	}

	// Fetch target quota lists synchronously
	for _, listID := range s.TargetQuotaListIDs {
		if listID != -1 {
			entries, err := fetchListResiliently(listID)
			if err != nil {
				slog.Error("Failed to fetch target quota list after all retries", "list_id", listID, "error", err)
				allListsLoaded = false
				continue
			}
			parseAndLoadApplications(entries, core.CompetitionTargetQuota, headingCode, receiver)
		}
	}

	if allListsLoaded {
		receiver.PutHeadingData(&source.HeadingData{
			Code:       headingCode,
			Capacities: s.Capacities,
			PrettyName: s.PrettyName,
		})
		return nil
	}
	return fmt.Errorf("not all lists were successfully loaded for heading %s", s.PrettyName)
}
