package hse

import (
	"analabit/core"
	"analabit/source"
	"analabit/utils"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/xuri/excelize/v2"
)

// HttpHeadingSource defines how to load HSE heading data from URLs.
// It assumes RCListURL is the primary source for the heading name.
type HttpHeadingSource struct {
	RCListURL url.URL // For "Основные конкурсные места" (CompetitionRegular)
	TQListURL url.URL // For "Целевая квота" (CompetitionTargetQuota)
	DQListURL url.URL // For "Отдельная квота" (CompetitionDedicatedQuota)
	SQListURL url.URL // For "Особая квота" (CompetitionSpecialQuota)
	BListURL  url.URL // For "Без вступительных испытаний" (CompetitionBVI)
	Capacity  int
}

// openHttpExcelFile downloads and opens an Excel file from a URL.
// Returns (nil, nil) if urlStr is empty or invalid, to allow skipping.
// Returns (nil, error) for actual download/open errors.
func openHttpExcelFile(urlStr string, listName string) (*excelize.File, error) {
	if urlStr == "" || urlStr == "." || urlStr == "/" { // Check for effectively empty URLs
		log.Printf("Skipping %s: URL ('%s') is empty or invalid.", listName, urlStr)
		return nil, nil // Indicate skippable
	}

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
func (s *HttpHeadingSource) LoadTo(headings chan<- source.HeadingData, applications chan<- source.ApplicationData) error {
	defer close(headings)
	defer close(applications)

	if s.RCListURL.Path == "" || s.RCListURL.Path == "." || s.RCListURL.Path == "/" {
		return fmt.Errorf("RCListURL is mandatory and was not provided or is invalid in HttpHeadingSource")
	}
	rcListURLString := s.RCListURL.String()

	log.Printf("Attempting to extract heading name from RCListURL: %s", rcListURLString)
	rcListFile, err := openHttpExcelFile(rcListURLString, "RC List (for name extraction)")
	if err != nil {
		return fmt.Errorf("failed to open primary RCList file from %s: %w", rcListURLString, err)
	}
	if rcListFile == nil { // Should not happen due to the check above, but as a safeguard
		return fmt.Errorf("primary RCList file from %s could not be opened (was skipped by openHttpExcelFile)", rcListURLString)
	}
	defer func() {
		if err := rcListFile.Close(); err != nil {
			log.Printf("Error closing primary RCList Excel file (opened from %s): %v", rcListURLString, err)
		}
	}()

	prettyName, err := extractPrettyNameFromXLS(rcListFile, rcListURLString)
	if err != nil {
		return fmt.Errorf("failed to extract pretty name using RCList %s: %w", rcListURLString, err)
	}

	headingCode := utils.GenerateHeadingCode(prettyName)
	// Send HeadingData to the channel
	headings <- source.HeadingData{
		Code:       headingCode,
		Capacity:   s.Capacity,
		PrettyName: prettyName,
	}
	log.Printf("Sent heading: %s (Code: %s, Capacity: %d) using name from %s", prettyName, headingCode, s.Capacity, rcListURLString)

	definitions := []listDefinition{
		{Source: s.BListURL.String(), CompetitionType: core.CompetitionBVI, ListName: "BVI List"},
		{Source: s.SQListURL.String(), CompetitionType: core.CompetitionSpecialQuota, ListName: "Special Quota List"},
		{Source: s.DQListURL.String(), CompetitionType: core.CompetitionDedicatedQuota, ListName: "Dedicated Quota List"},
		{Source: s.TQListURL.String(), CompetitionType: core.CompetitionTargetQuota, ListName: "Target Quota List"},
		{Source: rcListURLString, CompetitionType: core.CompetitionRegular, ListName: "Common List"},
	}

	// Pass the applications channel to processApplicationsFromLists
	return processApplicationsFromLists(applications, headingCode, prettyName, definitions, rcListFile, rcListURLString, openHttpExcelFile)
}
