// Package mipt provides parsing functionality for MIPT (Moscow Institute of Physics and Technology) admission lists.
// It implements the HeadingSource interface to extract application data from MIPT's admission list HTML files.
package mipt

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"golang.org/x/net/html"
)

// Common patterns and mappings for MIPT parsing
var (
	// Competition type mapping for MIPT categories
	competitionTypeMap = map[string]core.Competition{
		"БВИ":       core.CompetitionBVI,
		"Целевая":   core.CompetitionTargetQuota,
		"Особая":    core.CompetitionSpecialQuota,
		"Отдельная": core.CompetitionDedicatedQuota,
		"Общий":     core.CompetitionRegular,
		"Основная":  core.CompetitionRegular,
	}

	// Regex patterns for data extraction from MIPT HTML
	// Based on actual structure: <td>315</td><td>4102004</td><td>0</td><td>0</td><td>262</td>
	positionRegex    = regexp.MustCompile(`^\d+`)
	studentIDRegex   = regexp.MustCompile(`^\d{7}$`)
	scoreRegex       = regexp.MustCompile(`^\d{1,3}$`)
	priorityRegex    = regexp.MustCompile(`^\d+$`)
	originalDocRegex = regexp.MustCompile(`^(Да|Нет)$`)
)

// extractHeadingFromTitle extracts the program name from the HTML title
func extractHeadingFromTitle(doc *html.Node) string {
	var title string
	var findTitle func(*html.Node)
	findTitle = func(n *html.Node) {
		if title != "" {
			return
		}
		if n.Type == html.ElementNode && n.Data == "title" {
			title = getTextContent(n)
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findTitle(c)
		}
	}
	findTitle(doc)

	if title != "" {
		// MIPT titles typically contain the program name
		// Remove common prefixes and suffixes to get clean program name
		title = strings.TrimSpace(title)

		// Remove common MIPT prefixes
		prefixes := []string{"МФТИ", "Московский физико-технический институт", "Конкурсные списки"}
		for _, prefix := range prefixes {
			if strings.HasPrefix(title, prefix) {
				title = strings.TrimSpace(strings.TrimPrefix(title, prefix))
				title = strings.TrimLeft(title, "- ")
			}
		}

		return title
	}

	return ""
}

// tableFormat stores the column indices for different data fields.
type tableFormat struct {
	studentIDColumn   int
	priorityColumn    int
	totalScoreColumn  int
	bviColumn         int
	consentColumn     int
}

// detectTableFormat returns a hardcoded table format for the new MIPT layout.
func detectTableFormat(headerRow *html.Node) *tableFormat {
	return &tableFormat{
		studentIDColumn:   2,  // "Уникальный код"
		priorityColumn:    1,  // Priority is the 2nd column
		totalScoreColumn:  5,  // "Сумма баллов"
		bviColumn:         16, // "Без вступительных испытаний"
		consentColumn:     11, // "Согласие на зачисление"
	}
}

// parseApplicantFromTableRow extracts application data from a MIPT table row.
func parseApplicantFromTableRow(row *html.Node, defaultCompetitionType core.Competition, format *tableFormat) (*source.ApplicationData, error) {
	cells := extractTableCells(row)
	if len(cells) < 9 { // Need at least 9 columns for basic data
		return nil, fmt.Errorf("insufficient columns in row: got %d, expected at least 9", len(cells))
	}

	// Column 0: Position (e.g., "315")
	positionText := strings.TrimSpace(getTextContent(cells[0]))
	if !positionRegex.MatchString(positionText) {
		return nil, fmt.Errorf("invalid position format: %s", positionText)
	}
	position, err := strconv.Atoi(positionText)
	if err != nil {
		return nil, fmt.Errorf("invalid position number: %v", err)
	}

	// Column with Student ID
	if format.studentIDColumn >= len(cells) {
		return nil, fmt.Errorf("student ID column not found")
	}
	studentIDText := strings.TrimSpace(getTextContent(cells[format.studentIDColumn]))
	if !studentIDRegex.MatchString(studentIDText) {
		return nil, fmt.Errorf("invalid student ID format: %s", studentIDText)
	}

	// Column with Total Score
	var scoresSum int
	if format.totalScoreColumn < len(cells) {
		totalScoreText := strings.TrimSpace(getTextContent(cells[format.totalScoreColumn]))
		if scoreRegex.MatchString(totalScoreText) {
			scoresSum, _ = strconv.Atoi(totalScoreText)
		}
	}

	// Column with Priority
	priority := 1 // Default priority
	if format.priorityColumn < len(cells) {
		priorityText := strings.TrimSpace(getTextContent(cells[format.priorityColumn]))
		if priorityRegex.MatchString(priorityText) && priorityText != "0" {
			if p, err := strconv.Atoi(priorityText); err == nil && p > 0 {
				priority = p
			}
		}
	}

	// Determine competition type: Check BVI column
	competitionType := defaultCompetitionType
	if defaultCompetitionType == core.CompetitionRegular && format.bviColumn < len(cells) {
		bviColumnText := strings.TrimSpace(getTextContent(cells[format.bviColumn]))
		if strings.Contains(bviColumnText, "✓") || strings.Contains(bviColumnText, "Диплом") {
			competitionType = core.CompetitionBVI
		}
	}

	// Column with "Согласие на зачисление" (consent for enrollment)
	originalSubmitted := false
	if format.consentColumn < len(cells) {
		consentText := strings.TrimSpace(getTextContent(cells[format.consentColumn]))
		originalSubmitted = strings.Contains(consentText, "✓")
	}

	return &source.ApplicationData{
		StudentID:         studentIDText,
		RatingPlace:       position,
		Priority:          priority,
		CompetitionType:   competitionType,
		ScoresSum:         scoresSum,
		OriginalSubmitted: originalSubmitted,
	}, nil
}

// extractTableCells extracts all <td> and <th> elements from a table row.
func extractTableCells(row *html.Node) []*html.Node {
	var cells []*html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "td" || n.Data == "th") {
			cells = append(cells, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(row)
	return cells
}

// findTableRows finds all <tr> elements within a given table node.
func findTableRows(table *html.Node) []*html.Node {
	var rows []*html.Node
	var find func(*html.Node)
	find = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			rows = append(rows, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			find(c)
		}
	}
	find(table)
	return rows
}

// getTextContent recursively extracts all text content from an HTML node
func getTextContent(n *html.Node) string {
	if n == nil {
		return ""
	}
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}

	var sb strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		textContent := getTextContent(c)
		if textContent != "" {
			if sb.Len() > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(textContent)
		}
	}
	return strings.TrimSpace(sb.String())
}

// findTableRows finds all table rows (tr elements) in the document
// findTableHeaderRow finds the header row (tr with th elements) in the document.
func findTableHeaderRow(doc *html.Node) (*html.Node, *html.Node) {
	var headerRow *html.Node
	var tableNode *html.Node
	var findHeader func(*html.Node)
	findHeader = func(n *html.Node) {
		if headerRow != nil {
			return
		}
		if n.Type == html.ElementNode && n.Data == "tr" {
			// Check if this row contains any <th> elements
			hasTh := false
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "th" {
					hasTh = true
					break
				}
			}
			if hasTh {
				headerRow = n
				// Now find the parent table
				p := n.Parent
				for p != nil {
					if p.Type == html.ElementNode && p.Data == "table" {
						tableNode = p
						break
					}
					p = p.Parent
				}
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findHeader(c)
		}
	}
	findHeader(doc)
	return headerRow, tableNode
}



// determineCompetitionType analyzes the document or context to determine the competition type
func determineCompetitionType(doc *html.Node, filename string) core.Competition {
	// Check filename for competition type hints
	filenameUpper := strings.ToUpper(filename)

	if strings.Contains(filenameUpper, "BVI") || strings.Contains(filenameUpper, "БВИ") {
		return core.CompetitionBVI
	}
	if strings.Contains(filenameUpper, "TARGET") || strings.Contains(filenameUpper, "ЦЕЛЕВАЯ") {
		return core.CompetitionTargetQuota
	}
	if strings.Contains(filenameUpper, "SPECIAL") || strings.Contains(filenameUpper, "ОСОБАЯ") {
		return core.CompetitionSpecialQuota
	}
	if strings.Contains(filenameUpper, "DEDICATED") || strings.Contains(filenameUpper, "ОТДЕЛЬНАЯ") {
		return core.CompetitionDedicatedQuota
	}

	// Check document content for competition type indicators
	docText := strings.ToUpper(getTextContent(doc))
	if strings.Contains(docText, "БВИ") {
		return core.CompetitionBVI
	}
	if strings.Contains(docText, "ЦЕЛЕВАЯ") {
		return core.CompetitionTargetQuota
	}
	if strings.Contains(docText, "ОСОБАЯ") {
		return core.CompetitionSpecialQuota
	}
	if strings.Contains(docText, "ОТДЕЛЬНАЯ") {
		return core.CompetitionDedicatedQuota
	}

	// Default to regular competition
	return core.CompetitionRegular
}
