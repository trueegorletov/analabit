package hse

import (
	"analabit/core"
	"analabit/core/source"
	"analabit/core/utils"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

// extractPrettyNameFromXLS extracts the heading name from cell A3 of the given XLSX file.
// sourceHint is for logging (URL or file path).
func extractPrettyNameFromXLS(f *excelize.File, sourceHint string) (string, error) {
	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		return "", fmt.Errorf("no sheets found in Excel file (source: %s)", sourceHint)
	}
	sheetName := sheetList[0]

	cellValue, err := f.GetCellValue(sheetName, "A3")
	if err != nil {
		return "", fmt.Errorf("failed to get cell A3 from Excel file (source: %s, sheet: %s): %w", sourceHint, sheetName, err)
	}

	re := regexp.MustCompile(`"([^"]*)"`)
	matches := re.FindStringSubmatch(cellValue)
	if len(matches) >= 2 {
		return matches[1], nil
	}

	parts := strings.SplitN(cellValue, "Образовательная программа ", 2)
	if len(parts) == 2 {
		name := strings.TrimSpace(parts[1])
		if name != "" {
			return name, nil
		}
	}
	return "", fmt.Errorf("could not extract pretty name from '%s' in Excel file (source: %s, sheet: %s, cell A3)", cellValue, sourceHint, sheetName)
}

// parseAndLoadApplicationsInternal parses an Excel file and sends application data to the applications channel.
// sourceHint is for logging (URL or file path).
func parseAndLoadApplicationsInternal(
	f *excelize.File, competitionType core.Competition, headingCode string, receiver source.DataReceiver, sourceHint string) error {
	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		return fmt.Errorf("no sheets found in Excel file (hint: %s)", sourceHint)
	}
	sheetName := sheetList[0]

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("could not get rows from Excel file (hint: %s, sheet %s): %w", sourceHint, sheetName, err)
	}

	startRowIndex := 12 // Data starts from row 12
	if len(rows) <= startRowIndex {
		log.Printf("No data rows found in %s (sheet %s) starting from row 12 (expected at least %d rows, got %d).", sourceHint, sheetName, startRowIndex+1, len(rows))
		return nil // No data to process is not an error itself
	}

	for i := startRowIndex; i < len(rows); i++ {
		row := rows[i]

		var ratingPlaceStr string

		if len(row) == 0 {
			continue
		} else if ratingPlaceStr = strings.TrimSpace(row[0]); ratingPlaceStr == "" {
			continue
		}

		var ratingPlace int

		ratingPlace, err = strconv.Atoi(ratingPlaceStr)
		if err != nil {
			log.Printf("Warning: Invalid rating place '%s' in row %d of %s (sheet %s): %v. Skipping this row.", ratingPlaceStr, i+1, sourceHint, sheetName, err)
			continue // Skip this row if conversion fails
		}

		var studentID string
		if len(row) > 1 {
			studentID = strings.TrimSpace(row[1])
		}
		if studentID == "" {
			continue
		}

		preparedID, errPrepare := utils.PrepareStudentID(studentID) // Changed to exported function
		if errPrepare != nil {
			log.Printf("Warning: Invalid student ID '%s' in row %d of %s (sheet %s): %v. Skipping this row.", studentID, i+1, sourceHint, sheetName, errPrepare)
			continue
		}
		studentID = preparedID

		var originalSubmitted bool // This variable was removed in previous steps as unused
		if len(row) > 4 {
			originalStr := strings.TrimSpace(strings.ToLower(row[4]))
			if originalStr == "да" {
				originalSubmitted = true // This variable was removed in previous steps as unused
			} else if originalStr == "нет" {
				originalSubmitted = false
			} else if originalStr != "" {
				originalSubmitted = false
			}
		}

		var scoresSum int
		if len(row) > 10 {
			scoresStr := strings.TrimSpace(row[10])
			if scoresStr != "" {
				scores, errConv := strconv.Atoi(scoresStr)
				if errConv != nil {
					scoresSum = 0
				} else {
					scoresSum = scores
				}
			}
		}

		var priority int
		if len(row) > 12 {
			priorityStr := strings.TrimSpace(row[12])
			if priorityStr != "" {
				prio, errConv := strconv.Atoi(priorityStr)
				if errConv != nil {
					priority = 0
				} else {
					priority = prio
				}
			}
		}

		receiver.PutApplicationData(&source.ApplicationData{
			HeadingCode:       headingCode,
			StudentID:         studentID,
			RatingPlace:       ratingPlace,
			ScoresSum:         scoresSum,
			Priority:          priority,
			CompetitionType:   competitionType,
			OriginalSubmitted: originalSubmitted,
		})
	}
	return nil
}

// listDefinition holds information about a specific competition list.
type listDefinition struct {
	Source          string // URL string for HTTP or file path for File sources.
	CompetitionType core.Competition
	ListName        string
}

// processApplicationsFromLists iterates through defined lists, opens/downloads them,
// and sends application data via the applications channel.
func processApplicationsFromLists(
	receiver source.DataReceiver,
	headingCode string,
	prettyName string,
	definitions []listDefinition,
	primaryFile *excelize.File,
	primaryFileSourceHint string,
	openFileFunc func(source string, listName string) (*excelize.File, error),
) error {
	for _, listDef := range definitions {
		if listDef.Source == "" { // Already checked by callers before adding to definitions
			log.Printf("Skipping %s: source is empty (should have been filtered before).", listDef.ListName)
			continue
		}

		log.Printf("Processing %s from %s for heading %s", listDef.ListName, listDef.Source, prettyName)

		var currentFile *excelize.File
		var err error
		closeCurrentFile := false

		if listDef.Source == primaryFileSourceHint {
			currentFile = primaryFile
			log.Printf("Reusing already opened Excel file for %s from %s", listDef.ListName, listDef.Source)
		} else {
			currentFile, err = openFileFunc(listDef.Source, listDef.ListName)
			if err != nil { // An actual error occurred during opening/downloading
				log.Printf("Error opening/processing source for %s from %s: %v. Continuing.", listDef.ListName, listDef.Source, err)
				continue
			}
			if currentFile == nil { // openFileFunc decided to skip this list (e.g., invalid URL/path but not a hard error)
				// Log message should have been printed by openFileFunc
				continue
			}
			closeCurrentFile = true
		}

		// Pass the applications channel to parseAndLoadApplicationsInternal
		err = parseAndLoadApplicationsInternal(currentFile, listDef.CompetitionType, headingCode, receiver, listDef.Source)
		if err != nil {
			log.Printf("Warning: Failed to load applications from %s (%s): %v. Continuing.",
				listDef.ListName, listDef.Source, err)
		}

		if closeCurrentFile {
			if closeErr := currentFile.Close(); closeErr != nil {
				log.Printf("Error closing Excel file for %s (from %s): %v", listDef.ListName, listDef.Source, closeErr)
			}
		}
	}
	return nil
}
