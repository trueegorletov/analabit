package mipt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/utils"
	"golang.org/x/net/html"
)

// FileHeadingSource implements source.HeadingSource for MIPT HTML files
type FileHeadingSource struct {
	FilePath   string          `json:"file_path"`
	PrettyName string          `json:"pretty_name"` // Fallback if parsing fails
	Capacities core.Capacities `json:"capacities"`  // Manual capacities if known
}

// LoadTo parses MIPT application data from a local HTML file
func (f *FileHeadingSource) LoadTo(receiver source.DataReceiver) error {
	if f.FilePath == "" {
		return fmt.Errorf("FilePath is required for MIPT FileHeadingSource")
	}

	// Check if file exists
	if _, err := os.Stat(f.FilePath); os.IsNotExist(err) {
		return fmt.Errorf("MIPT file does not exist at path %s: %w", f.FilePath, err)
	}

	// Read the HTML file
	htmlContent, err := os.ReadFile(f.FilePath)
	if err != nil {
		return fmt.Errorf("failed to read MIPT HTML file %s: %v", f.FilePath, err)
	}

	// Parse the HTML
	doc, err := html.Parse(strings.NewReader(string(htmlContent)))
	if err != nil {
		return fmt.Errorf("failed to parse MIPT HTML: %v", err)
	}

	// Extract program name
	prettyName := f.PrettyName
	if prettyName == "" {
		prettyName = extractHeadingFromTitle(doc)
	}
	if prettyName == "" {
		// Use filename as fallback
		filename := filepath.Base(f.FilePath)
		prettyName = strings.TrimSuffix(filename, filepath.Ext(filename))
	}

	// Determine competition type from filename or content
	competitionType := determineCompetitionType(doc, filepath.Base(f.FilePath))

	// Generate heading code
	headingCode := utils.GenerateHeadingCode(prettyName)

	// Send heading data
	headingData := &source.HeadingData{
		Code:       headingCode,
		PrettyName: prettyName,
		Capacities: f.Capacities, // Use provided capacities or zero values
	}
	receiver.PutHeadingData(headingData)

	// Parse applications from the table
	err = f.parseApplicationsFromTable(doc, receiver, headingCode, competitionType)
	if err != nil {
		return fmt.Errorf("failed to parse MIPT applications: %v", err)
	}

	return nil
}

// parseApplicationsFromTable extracts applications from the MIPT HTML table
func (f *FileHeadingSource) parseApplicationsFromTable(doc *html.Node, receiver source.DataReceiver, headingCode string, competitionType core.Competition) error {
	// Find all table rows
	rows := findTableRows(doc)

	fmt.Printf("DEBUG: Found %d total table rows\n", len(rows))

	if len(rows) == 0 {
		return fmt.Errorf("no table rows found in MIPT HTML file")
	}

	var applicationsProcessed int
	var errors []string

	// Process each row
	for i, row := range rows {
		// Skip header rows (usually the first row or rows with th elements)
		if f.isHeaderRow(row) {
			//fmt.Printf("DEBUG: Row %d is header row, skipping\n", i+1)
			continue
		}

		// Skip empty rows
		cells := extractTableCells(row)
		if len(cells) == 0 {
			//fmt.Printf("DEBUG: Row %d has no cells, skipping\n", i+1)
			continue
		}

		// Check if this looks like a data row (first cell should be a position number)
		firstCellText := strings.TrimSpace(getTextContent(cells[0]))
		//fmt.Printf("DEBUG: Row %d first cell: '%s' (cells: %d)\n", i+1, firstCellText, len(cells))

		// Skip header rows and invalid rows - MIPT data has position numbers like "1", "2", "315" etc.
		if firstCellText == "" || !positionRegex.MatchString(firstCellText) {
			//fmt.Printf("DEBUG: Row %d first cell doesn't match position format, skipping\n", i+1)
			continue
		}

		// Parse the applicant data
		app, err := parseApplicantFromTableRow(row, competitionType)
		if err != nil {
			errors = append(errors, fmt.Sprintf("row %d: %v", i+1, err))
			continue
		}

		// Set the heading code for this application
		app.HeadingCode = headingCode

		// Send the application data
		receiver.PutApplicationData(app)
		applicationsProcessed++
	}

	// Log summary
	fmt.Printf("MIPT parser processed %d applications from %s\n", applicationsProcessed, f.FilePath)

	if len(errors) > 0 {
		fmt.Printf("MIPT parser encountered %d errors:\n", len(errors))
		for _, err := range errors {
			fmt.Printf("  - %s\n", err)
		}
	}

	if applicationsProcessed == 0 {
		return fmt.Errorf("no valid applications found in MIPT HTML file")
	}

	return nil
}

// isHeaderRow determines if a table row is a header row
func (f *FileHeadingSource) isHeaderRow(row *html.Node) bool {
	// Check if row contains th elements (table header cells)
	var hasThElements bool
	var checkForTh func(*html.Node)
	checkForTh = func(n *html.Node) {
		if hasThElements {
			return
		}
		if n.Type == html.ElementNode && n.Data == "th" {
			hasThElements = true
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			checkForTh(c)
		}
	}
	checkForTh(row)

	if hasThElements {
		return true
	}

	// Check if first cell looks like a header (contains column names)
	cells := extractTableCells(row)
	if len(cells) > 0 {
		firstCellText := strings.ToLower(strings.TrimSpace(getTextContent(cells[0])))
		headerKeywords := []string{"№", "номер", "место", "позиция", "фио", "имя"}
		for _, keyword := range headerKeywords {
			if strings.Contains(firstCellText, keyword) {
				return true
			}
		}
	}

	return false
}
