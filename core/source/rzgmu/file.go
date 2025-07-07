// Package rzgmu provides support for loading RZGMU (Рязанский государственный медицинский университет) admission data.
package rzgmu

import (
	"analabit/core"
	"analabit/core/source"
	"analabit/core/utils"
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

// FileHeadingSource defines how to load RZGMU heading data from a local PDF file.
// RZGMU provides admission lists in PDF format that need to be processed with Tabula.
type FileHeadingSource struct {
	FilePath    string           // Path to the PDF file containing the admission list
	ProgramName string           // Name of the educational program (if empty, will use default)
	Capacities  *core.Capacities // Capacities for this heading. If nil, will be extracted from PDF
}

// LoadTo loads data from local file, extracting table data with Tabula,
// and sending HeadingData and ApplicationData to the provided receiver.
func (s *FileHeadingSource) LoadTo(receiver source.DataReceiver) error {
	if s.FilePath == "" {
		return fmt.Errorf("FilePath is required for RZGMU FileHeadingSource")
	}

	// Check if file exists
	if _, err := os.Stat(s.FilePath); os.IsNotExist(err) {
		return fmt.Errorf("PDF file does not exist at path %s: %w", s.FilePath, err)
	}

	log.Printf("Processing RZGMU admission list from: %s", s.FilePath)

	// Extract CSV data using Tabula
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	csvData, err := extractCSVWithTabula(ctx, s.FilePath)
	if err != nil {
		return fmt.Errorf("failed to extract CSV from RZGMU PDF %s: %w", s.FilePath, err)
	}

	// Parse the CSV data using Tabula-specific parser
	programs, err := parseRZGMUTabulaData(csvData)
	if err != nil {
		return fmt.Errorf("failed to parse RZGMU data from %s: %w", s.FilePath, err)
	}

	// Process each program found in the data (usually just one per PDF)
	for _, program := range programs {
		// Use specified program name if provided, otherwise use extracted name
		programName := s.ProgramName
		if programName == "" {
			programName = program.PrettyName
		}

		capacities := s.Capacities
		if capacities == nil {
			// Extract capacities from the program data
			capacities = &program.ExtractedCapacities
		}

		headingCode := utils.GenerateHeadingCode(programName)

		// Send HeadingData to the receiver
		receiver.PutHeadingData(&source.HeadingData{
			Code:       headingCode,
			Capacities: *capacities,
			PrettyName: programName,
		})

		log.Printf("Sent RZGMU heading: %s (Code: %s, Caps: %v)", programName, headingCode, *capacities)

		// Send ApplicationData for each application in this program
		for _, app := range program.Applications {
			app.HeadingCode = headingCode
			receiver.PutApplicationData(app)
		}

		log.Printf("Sent %d applications for RZGMU heading %s", len(program.Applications), programName)
	}

	log.Printf("Successfully processed RZGMU heading(s) from %s", s.FilePath)
	return nil
}
