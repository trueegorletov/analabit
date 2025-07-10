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
	positionRegex    = regexp.MustCompile(`^\d+$`)
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

// parseApplicantFromTableRow extracts application data from a MIPT table row
// Based on structure: Position|StudentID|Priority1|Priority2|TotalScore|ExamScore|IndividualScore|Status|OriginalDoc|...|BVI_Column|...
func parseApplicantFromTableRow(row *html.Node, defaultCompetitionType core.Competition) (*source.ApplicationData, error) {
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

	// Column 1: Student ID (e.g., "4102004")
	studentIDText := strings.TrimSpace(getTextContent(cells[1]))
	if !studentIDRegex.MatchString(studentIDText) {
		return nil, fmt.Errorf("invalid student ID format: %s", studentIDText)
	}

	// Column 4: Total Score (e.g., "262")
	totalScoreText := strings.TrimSpace(getTextContent(cells[4]))
	var scoresSum int
	if scoreRegex.MatchString(totalScoreText) {
		scoresSum, _ = strconv.Atoi(totalScoreText)
	}

	// Column 2 or 3: Priority (usually column 2, fallback to 3)
	priority := 1 // Default priority
	for _, col := range []int{2, 3} {
		if col < len(cells) {
			priorityText := strings.TrimSpace(getTextContent(cells[col]))
			if priorityRegex.MatchString(priorityText) && priorityText != "0" {
				if p, err := strconv.Atoi(priorityText); err == nil && p > 0 {
					priority = p
					break
				}
			}
		}
	}

	// Determine competition type: Check BVI column (around column 13 based on sample data)
	competitionType := defaultCompetitionType
	if defaultCompetitionType == core.CompetitionRegular && len(cells) > 13 {
		bviColumnText := strings.TrimSpace(getTextContent(cells[13]))
		// Look for checkmark or "БВИ" indicators in the BVI column
		if strings.Contains(bviColumnText, "✓") || strings.Contains(bviColumnText, "Диплом") {
			competitionType = core.CompetitionBVI
		} else {
			competitionType = core.CompetitionRegular
		}
	}

	// Column 10: "Согласие на зачисление" (consent for enrollment) - look for checkmark "✓"
	originalSubmitted := false
	if len(cells) > 10 {
		consentText := strings.TrimSpace(getTextContent(cells[10]))
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

// extractTableCells extracts all td elements from a table row
func extractTableCells(row *html.Node) []*html.Node {
	var cells []*html.Node
	var findCells func(*html.Node)
	findCells = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "td" {
			cells = append(cells, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findCells(c)
		}
	}
	findCells(row)
	return cells
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
func findTableRows(doc *html.Node) []*html.Node {
	var rows []*html.Node
	var findRows func(*html.Node)
	findRows = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			rows = append(rows, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findRows(c)
		}
	}
	findRows(doc)
	return rows
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
