package spbsu

import (
	"analabit/core"
	"analabit/core/source"
	"analabit/core/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Global semaphore to limit concurrent HTTP requests (mirroring HSE pattern)
var httpRequestSemaphore *semaphoreWeighted

type semaphoreWeighted struct {
	ch chan struct{}
}

func newSemaphoreWeighted(n int64) *semaphoreWeighted {
	return &semaphoreWeighted{ch: make(chan struct{}, n)}
}

func (s *semaphoreWeighted) Acquire(ctx context.Context, n int64) error {
	for i := int64(0); i < n; i++ {
		select {
		case s.ch <- struct{}{}:
			// acquired
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

func (s *semaphoreWeighted) Release(n int64) {
	for i := int64(0); i < n; i++ {
		<-s.ch
	}
}

func init() {
	maxConcurrentRequests := int64(4)
	if envVal := os.Getenv("HTTP_MAX_CONCURRENT_REQUESTS"); envVal != "" {
		if parsed, err := strconv.ParseInt(envVal, 10, 64); err == nil && parsed > 0 {
			maxConcurrentRequests = parsed
		} else {
			log.Printf("Warning: Invalid HTTP_MAX_CONCURRENT_REQUESTS value '%s', using default %d", envVal, maxConcurrentRequests)
		}
	}
	httpRequestSemaphore = newSemaphoreWeighted(maxConcurrentRequests)
	log.Printf("Initialized SPbSU HTTP request semaphore with limit: %d concurrent requests", maxConcurrentRequests)
}

// HttpHeadingSource loads SPbSU heading data from up to 4 JSON list IDs.
type HttpHeadingSource struct {
	PrettyName           string
	RegularListID        int
	TargetQuotaListID    int
	DedicatedQuotaListID int
	SpecialQuotaListID   int
	Capacities           core.Capacities
}

// fetchSpbsuListByID fetches and decodes the SPbSU list from a list ID.
func fetchSpbsuListByID(listID int) ([]SpbsuApplicationEntry, error) {
	if listID == -1 {
		return nil, nil
	}
	url := "https://enrollelists.spbu.ru/lists?id=" + strconv.Itoa(listID) + "&without_control=true"
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := httpRequestSemaphore.Acquire(ctx, 1); err != nil {
		return nil, fmt.Errorf("failed to acquire semaphore for %s: %w", url, err)
	}
	defer httpRequestSemaphore.Release(1)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download %s (status code %d)", url, resp.StatusCode)
	}
	return decodeSpbsuList(resp.Body)
}

// LoadTo implements source.HeadingSource for HttpHeadingSource.
func (s *HttpHeadingSource) LoadTo(receiver source.DataReceiver) error {
	if s.PrettyName == "" {
		return fmt.Errorf("PrettyName is required for SPbSU HttpHeadingSource")
	}
	headingCode := utils.GenerateHeadingCode(s.PrettyName)
	receiver.PutHeadingData(&source.HeadingData{
		Code:       headingCode,
		Capacities: s.Capacities,
		PrettyName: s.PrettyName,
	})

	listDefs := []struct {
		ListID      int
		Competition core.Competition
		ListName    string
	}{
		{s.RegularListID, core.CompetitionRegular, "Regular List"},
		{s.TargetQuotaListID, core.CompetitionTargetQuota, "Target Quota List"},
		{s.DedicatedQuotaListID, core.CompetitionDedicatedQuota, "Dedicated Quota List"},
		{s.SpecialQuotaListID, core.CompetitionSpecialQuota, "Special Quota List"},
	}

	for _, def := range listDefs {
		if def.ListID == -1 {
			continue
		}
		entries, err := fetchSpbsuListByID(def.ListID)
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
	return nil
}
