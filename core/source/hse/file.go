package hse

import (
	"analabit/core"
	"analabit/core/source"
	"analabit/core/utils"
	"fmt"
	"log"
	"os"

	"github.com/xuri/excelize/v2"
)

// FileHeadingSource defines how to load HSE heading data from local file paths.
// It assumes RCListPath is the primary source for the heading name.
type FileHeadingSource struct {
	RCListPath        string // For "Основные конкурсные места" (CompetitionRegular)
	TQListPath        string // For "Целевая квота" (CompetitionTargetQuota)
	DQListPath        string // For "Отдельная квота" (CompetitionDedicatedQuota)
	SQListPath        string // For "Особая квота" (CompetitionSpecialQuota)
	BListPath         string // For "Без вступительных испытаний" (CompetitionBVI)
	HeadingCapacities core.Capacities
}

// openLocalExcelFile opens an Excel file from the local filesystem.
// Returns (nil, nil) if filePath is empty, to allow skipping.
// Returns (nil, error) for actual open errors.
func openLocalExcelFile(filePath string, listName string) (*excelize.File, error) {
	if filePath == "" {
		log.Printf("Skipping %s: file path is empty.", listName)
		return nil, nil // Indicate skippable
	}

	// Basic check if file exists, though excelize.OpenFile will also check
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file for %s does not exist at path %s: %w", listName, filePath, err)
	}

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file for %s from %s: %w", listName, filePath, err)
	}
	return f, nil
}

// LoadTo loads data from local file system sources, sending HeadingData and ApplicationData to the provided channels.
func (s *FileHeadingSource) LoadTo(receiver source.DataReceiver) error {
	if s.RCListPath == "" {
		return fmt.Errorf("RCListPath is mandatory and was not provided in HseFileHeadingSource")
	}

	log.Printf("Attempting to extract heading name from RCListPath: %s", s.RCListPath)
	rcListFile, err := openLocalExcelFile(s.RCListPath, "RC List (for name extraction)")
	if err != nil {
		return fmt.Errorf("failed to open primary RCList file from %s: %w", s.RCListPath, err)
	}
	if rcListFile == nil { // Should not happen due to the check above
		return fmt.Errorf("primary RCList file from %s could not be opened (was skipped by openLocalExcelFile)", s.RCListPath)
	}
	defer func() {
		if err := rcListFile.Close(); err != nil {
			log.Printf("Error closing primary RCList Excel file (opened from %s): %v", s.RCListPath, err)
		}
	}()

	prettyName, err := extractPrettyNameFromXLS(rcListFile, s.RCListPath)
	if err != nil {
		return fmt.Errorf("failed to extract pretty name using RCList %s: %w", s.RCListPath, err)
	}

	headingCode := utils.GenerateHeadingCode(prettyName)
	// Send HeadingData to the channel
	receiver.PutHeadingData(&source.HeadingData{
		Code:       headingCode,
		Capacities: s.HeadingCapacities,
		PrettyName: prettyName,
	})
	log.Printf("Sent heading: %s (Code: %s, HeadingCapacities: %d) using name from %s", prettyName, headingCode, s.HeadingCapacities, s.RCListPath)

	definitions := []listDefinition{
		{Source: s.BListPath, CompetitionType: core.CompetitionBVI, ListName: "BVI List"},
		{Source: s.SQListPath, CompetitionType: core.CompetitionSpecialQuota, ListName: "Special Quota List"},
		{Source: s.DQListPath, CompetitionType: core.CompetitionDedicatedQuota, ListName: "Dedicated Quota List"},
		{Source: s.TQListPath, CompetitionType: core.CompetitionTargetQuota, ListName: "Target Quota List"},
		{Source: s.RCListPath, CompetitionType: core.CompetitionRegular, ListName: "Common List"},
	}

	// Pass the applications channel to processApplicationsFromLists
	return processApplicationsFromLists(receiver, headingCode, prettyName, definitions, rcListFile, s.RCListPath, openLocalExcelFile)
}
