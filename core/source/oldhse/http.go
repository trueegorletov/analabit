package oldhse

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/utils"

	"github.com/xuri/excelize/v2"
)

// HttpHeadingSource defines how to load HSE heading data from URLs.
// It assumes RCListURL is the primary source for the heading name.
type HttpHeadingSource struct {
	RCListURL  string // For "Основные конкурсные места" (CompetitionRegular)
	TQListURL  string // For "Целевая квота" (CompetitionTargetQuota)
	DQListURL  string // For "Отдельная квота" (CompetitionDedicatedQuota)
	SQListURL  string // For "Особая квота" (CompetitionSpecialQuota)
	BListURL   string // For "Без вступительных испытаний" (CompetitionBVI)
	Capacities core.Capacities
}

// openHttpExcelFile downloads and opens an Excel file from a URL.
// Returns (nil, nil) if urlStr is empty or invalid, to allow skipping.
// Returns (nil, error) for actual download/open errors.
func openHttpExcelFile(urlStr string, listName string) (*excelize.File, error) {
	if urlStr == "" || urlStr == "." || urlStr == "/" { // Check for effectively empty URLs
		log.Printf("Skipping %s: URL ('%s') is empty or invalid.", listName, urlStr)
		return nil, nil // Indicate skippable
	}

	// Acquire a semaphore slot, respecting context cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cancellation on function exit

	release, err := source.AcquireHTTPSemaphores(ctx, "oldhse")
	if err != nil {
		return nil, fmt.Errorf("failed to acquire semaphores for %s from %s: %w", listName, urlStr, err)
	}
	defer release()

	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to download %s from %s: %w", listName, urlStr, err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("Error closing response body for %s from %s: %v", listName, urlStr, closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download %s from %s (status code %d)", listName, urlStr, resp.StatusCode)
	}

	f, err := excelize.OpenReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel content for %s from %s: %w", listName, urlStr, err)
	}
	return f, nil
}

// LoadTo loads data from HTTP sources, sending HeadingData and ApplicationData to the provided channels.
func (s *HttpHeadingSource) LoadTo(receiver source.DataReceiver) error {
	var primaryURL string
	var primaryFileSourceHint string
	var primaryFileListName string

	if s.RCListURL != "" && s.RCListURL != "." && s.RCListURL != "/" {
		primaryURL = s.RCListURL
		primaryFileSourceHint = s.RCListURL
		primaryFileListName = "RC List (for name extraction)"
	} else if s.TQListURL != "" && s.TQListURL != "." && s.TQListURL != "/" {
		primaryURL = s.TQListURL
		primaryFileSourceHint = s.TQListURL
		primaryFileListName = "TQ List (for name extraction)"
		log.Printf("RCListURL is empty or invalid. Falling back to TQListURL for name extraction: %s", primaryURL)
	} else if s.DQListURL != "" && s.DQListURL != "." && s.DQListURL != "/" {
		primaryURL = s.DQListURL
		primaryFileSourceHint = s.DQListURL
		primaryFileListName = "DQ List (for name extraction)"
		log.Printf("RCListURL and TQListURL are empty or invalid. Falling back to DQListURL for name extraction: %s", primaryURL)
	} else if s.SQListURL != "" && s.SQListURL != "." && s.SQListURL != "/" {
		primaryURL = s.SQListURL
		primaryFileSourceHint = s.SQListURL
		primaryFileListName = "SQ List (for name extraction)"
		log.Printf("RC, TQ, and DQ ListURLs are empty or invalid. Falling back to SQListURL for name extraction: %s", primaryURL)
	} else if s.BListURL != "" && s.BListURL != "." && s.BListURL != "/" {
		primaryURL = s.BListURL
		primaryFileSourceHint = s.BListURL
		primaryFileListName = "BVI List (for name extraction)"
		log.Printf("RC, TQ, DQ, and SQ ListURLs are empty or invalid. Falling back to BListURL for name extraction: %s", primaryURL)
	} else {
		return fmt.Errorf("all provided list URLs (RC, TQ, DQ, SQ, BVI) are empty or invalid in HttpHeadingSource; cannot extract heading name")
	}

	log.Printf("Attempting to extract heading name from %s: %s", primaryFileListName, primaryURL)
	primaryFile, err := openHttpExcelFile(primaryURL, primaryFileListName)
	if err != nil {
		return fmt.Errorf("failed to open primary file (%s) from %s: %w", primaryFileListName, primaryURL, err)
	}
	if primaryFile == nil { // Should not happen due to the check above, but as a safeguard
		return fmt.Errorf("primary file (%s) from %s could not be opened (was skipped by openHttpExcelFile)", primaryFileListName, primaryURL)
	}
	defer func() {
		if err := primaryFile.Close(); err != nil {
			log.Printf("Error closing primary Excel file (%s, opened from %s): %v", primaryFileListName, primaryURL, err)
		}
	}()

	prettyName, err := extractPrettyNameFromXLS(primaryFile, primaryURL)
	if err != nil {
		return fmt.Errorf("failed to extract pretty name using %s (%s): %w", primaryFileListName, primaryURL, err)
	}

	headingCode := utils.GenerateHeadingCode(prettyName)
	// Send HeadingData to the channel
	receiver.PutHeadingData(&source.HeadingData{
		Code:       headingCode,
		Capacities: s.Capacities,
		PrettyName: prettyName,
	})
	log.Printf("Sent heading: %s (Code: %s, Caps: %d) using name from %s (%s)", prettyName, headingCode, s.Capacities, primaryFileListName, primaryURL)

	// Define rcListURLForDefinitions for clarity, ensuring it's the original s.RCListURL for the Common List definition.
	// If s.RCListURL is empty, it will be handled by processApplicationsFromLists (skipped).
	rcListURLForDefinitions := s.RCListURL
	if rcListURLForDefinitions == "" || rcListURLForDefinitions == "." || rcListURLForDefinitions == "/" {
		log.Printf("RCListURL ('%s') is empty/invalid for Common List definition; it will be skipped if not the primaryFileSourceHint", s.RCListURL)
		// processApplicationsFromLists will skip if source is empty and not the primaryFileSourceHint
	}

	definitions := []listDefinition{
		{Source: s.SQListURL, CompetitionType: core.CompetitionSpecialQuota, ListName: "Special Quota List"},
		{Source: s.DQListURL, CompetitionType: core.CompetitionDedicatedQuota, ListName: "Dedicated Quota List"},
		{Source: s.TQListURL, CompetitionType: core.CompetitionTargetQuota, ListName: "Target Quota List"},
		{Source: s.BListURL, CompetitionType: core.CompetitionBVI, ListName: "BVI List"},
		{Source: rcListURLForDefinitions, CompetitionType: core.CompetitionRegular, ListName: "Common List"},
	}

	// Pass the applications channel to processApplicationsFromLists
	// primaryFileSourceHint is the URL of the file already opened as primaryFile
	return processApplicationsFromLists(receiver, headingCode, prettyName, definitions, primaryFile, primaryFileSourceHint, openHttpExcelFile)
}

// urlPtrToString safely converts a *url.URL to a string, returning "" if nil.
func urlPtrToString(u *url.URL) string {
	if u == nil {
		return ""
	}
	return u.String()
}
