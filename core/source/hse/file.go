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

// FileHeadingSource defines how to load HSE heading data from a local XLSX file.
// The new HSE format contains all applications for one program in a single file.
type FileHeadingSource struct {
	FilePath   string // Path to the XLSX file containing all applications for this heading
	Capacities core.Capacities
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

// LoadTo loads data from local file system sources, sending HeadingData and ApplicationData to the provided receiver.
func (s *FileHeadingSource) LoadTo(receiver source.DataReceiver) error {
	if s.FilePath == "" {
		return fmt.Errorf("FilePath is required for HSE FileHeadingSource")
	}

	log.Printf("Opening HSE admission list from: %s", s.FilePath)

	f, err := openLocalExcelFile(s.FilePath, "HSE XLSX file")
	if err != nil {
		return fmt.Errorf("failed to open HSE list from %s: %w", s.FilePath, err)
	}
	if f == nil {
		return fmt.Errorf("HSE file from %s could not be opened", s.FilePath)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			log.Printf("Error closing HSE Excel file from %s: %v", s.FilePath, closeErr)
		}
	}()

	// Extract heading name from the XLSX
	prettyName, err := extractPrettyNameFromXLSX(f, s.FilePath)
	if err != nil {
		return fmt.Errorf("failed to extract pretty name from HSE list at %s: %w", s.FilePath, err)
	}

	headingCode := utils.GenerateHeadingCode(prettyName)

	// Send HeadingData to the receiver
	receiver.PutHeadingData(&source.HeadingData{
		Code:       headingCode,
		Capacities: s.Capacities,
		PrettyName: prettyName,
	})

	log.Printf("Sent HSE heading: %s (Code: %s, Caps: %v)", prettyName, headingCode, s.Capacities)

	// Parse applications from the XLSX and send to receiver
	err = parseApplicationsFromXLSX(f, headingCode, receiver, s.FilePath)
	if err != nil {
		return fmt.Errorf("failed to parse applications from HSE list at %s: %w", s.FilePath, err)
	}

	log.Printf("Successfully processed HSE heading from %s", s.FilePath)
	return nil
}
