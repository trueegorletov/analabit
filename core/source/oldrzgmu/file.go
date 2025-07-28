// Package rzgmu provides support for loading RZGMU (Рязанский государственный медицинский университет) admission data.
// RZGMU admission lists are provided in PDF format and are parsed using rsc.io/pdf for text extraction.
package oldrzgmu

import (
	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/utils"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"rsc.io/pdf"
)

// FileHeadingSource defines how to load RZGMU heading data from a local PDF file.
// RZGMU provides admission lists in PDF format that are parsed using rsc.io/pdf.
type FileHeadingSource struct {
	Path        string           // Path to the PDF file containing the admission list
	ProgramName string           // Name of the educational program (if empty, will use default)
	Capacities  *core.Capacities // Capacities for this heading. If nil, will be extracted from PDF
}

// LoadTo loads data from file source, extracting text with rsc.io/pdf,
// and sending HeadingData and ApplicationData to the provided receiver.
func (s *FileHeadingSource) LoadTo(receiver source.DataReceiver) error {
	if s.Path == "" {
		return fmt.Errorf("Path is required for RZGMU FileHeadingSource")
	}

	log.Printf("Processing RZGMU admission list from file: %s", s.Path)

	// Extract text data using rsc.io/pdf
	textData, err := extractTextFromPDFFile(s.Path)
	if err != nil {
		return fmt.Errorf("failed to extract text from RZGMU PDF %s: %w", s.Path, err)
	}

	// Parse the text data
	programs, err := parseRZGMUTextData(textData)
	if err != nil {
		return fmt.Errorf("failed to parse RZGMU data from %s: %w", s.Path, err)
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
			// Use extracted capacities from the program data
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

	log.Printf("Successfully processed RZGMU heading(s) from %s", s.Path)
	return nil
}

// extractTextFromPDFFile uses pdftotext (if available) or rsc.io/pdf to extract text from a PDF file
func extractTextFromPDFFile(pdfPath string) (string, error) {
	// Try pdftotext first with layout preservation (better handling of font encoding and structure)
	if text, err := extractTextWithPDFToText(pdfPath); err == nil {
		log.Printf("Successfully extracted text using pdftotext from %s", pdfPath)
		return text, nil
	} else {
		log.Printf("pdftotext failed (%v), falling back to rsc.io/pdf for %s", err, pdfPath)
	}

	// Fallback to rsc.io/pdf
	return extractTextWithRscPDFFile(pdfPath)
}

// extractTextWithPDFToText uses the pdftotext command-line utility to extract text with layout preservation
func extractTextWithPDFToText(pdfPath string) (string, error) {
	// Check if pdftotext is available
	if _, err := exec.LookPath("pdftotext"); err != nil {
		return "", fmt.Errorf("pdftotext not available: %w", err)
	}

	// Create temporary output file
	tempFile, err := os.CreateTemp("", "rzgmu_text_*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Run pdftotext with layout preservation and UTF-8 encoding
	cmd := exec.Command("pdftotext", "-layout", "-enc", "UTF-8", pdfPath, tempFile.Name())
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("pdftotext command failed: %w", err)
	}

	// Read the extracted text
	textBytes, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read extracted text: %w", err)
	}

	return string(textBytes), nil
}

// extractTextWithRscPDFFile uses rsc.io/pdf to extract text from a PDF file (fallback method)
func extractTextWithRscPDFFile(pdfPath string) (string, error) {
	f, err := pdf.Open(pdfPath)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF with rsc.io/pdf: %w", err)
	}

	var buf bytes.Buffer
	numPages := f.NumPage()

	for i := 1; i <= numPages; i++ {
		page := f.Page(i)
		content := page.Content()

		// Group text elements by position to reconstruct lines
		var currentLine bytes.Buffer
		var lastY float64 = -1

		for _, t := range content.Text {
			// If Y position changed significantly, we're on a new line
			if lastY != -1 && (lastY-t.Y) > 5 {
				if currentLine.Len() > 0 {
					buf.WriteString(strings.TrimSpace(currentLine.String()))
					buf.WriteString("\n")
					currentLine.Reset()
				}
			}

			// Add text to current line
			currentLine.WriteString(t.S)
			if t.S != "" && t.S != " " {
				currentLine.WriteString(" ")
			}

			lastY = t.Y
		}

		// Add the last line
		if currentLine.Len() > 0 {
			buf.WriteString(strings.TrimSpace(currentLine.String()))
			buf.WriteString("\n")
		}
	}

	result := buf.String()

	// Clean up extra spaces and empty lines
	lines := strings.Split(result, "\n")
	var cleanLines []string
	for _, line := range lines {
		cleaned := strings.TrimSpace(line)
		if cleaned != "" {
			cleanLines = append(cleanLines, cleaned)
		}
	}

	return strings.Join(cleanLines, "\n"), nil
}
