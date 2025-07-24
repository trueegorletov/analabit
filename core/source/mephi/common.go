package mephi

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"golang.org/x/net/html"
)

// ParseCapacitiesRegistry parses the MEPhI capacities registry HTML to extract heading names and their KCP values.
// It handles rowspan attributes to map one KCP value to multiple headings.
func ParseCapacitiesRegistry(doc *html.Node) (map[string]int, error) {
	capacities := make(map[string]int)

	// Find the table
	table := findNode(doc, "table", "id", "myTable")
	if table == nil {
		return nil, fmt.Errorf("capacities table not found")
	}

	// Find tbody
	tbody := findNode(table, "tbody", "", "")
	if tbody == nil {
		return nil, fmt.Errorf("tbody not found in capacities table")
	}

	// Parse rows
	rows := findAllNodes(tbody, "tr", "", "")
	var pendingKCP int
	var pendingRowspan int

	for _, row := range rows {
		// Skip header row
		if hasClass(row, "header") {
			continue
		}

		cells := findAllNodes(row, "td", "", "")
		if len(cells) < 4 {
			continue
		}

		// Extract heading name from second column
		headingName := strings.TrimSpace(getTextContent(cells[1]))
		if headingName == "" {
			continue
		}

		// Check if there's a new KCP value in the 4th column ("Оч")
		kcpText := strings.TrimSpace(getTextContent(cells[3]))
		if kcpText != "" && kcpText != "&nbsp;" {
			// Parse KCP value
			if kcp, err := strconv.Atoi(kcpText); err == nil {
				pendingKCP = kcp
				// Check for rowspan
				if rowspanAttr := getAttr(cells[3], "rowspan"); rowspanAttr != "" {
					if rowspan, err := strconv.Atoi(rowspanAttr); err == nil {
						pendingRowspan = rowspan - 1 // -1 because current row is included
					}
				}
			}
		}

		// Assign KCP to current heading
		if pendingKCP > 0 {
			capacities[headingName] = pendingKCP
			if pendingRowspan > 0 {
				pendingRowspan--
			} else {
				pendingKCP = 0 // Reset if no more rowspan
			}
		}
	}

	return capacities, nil
}

// ParseListLinksRegistry parses the MEPhI list links registry HTML to extract heading names and competition URLs.
func ParseListLinksRegistry(doc *html.Node) (map[string]map[core.Competition][]string, error) {
	result := make(map[string]map[core.Competition][]string)

	// Find all tables
	tables := findAllNodes(doc, "table", "class", "w100")
	if len(tables) == 0 {
		return nil, fmt.Errorf("no tables found in links registry")
	}

	for _, table := range tables {
		tbody := findNode(table, "tbody", "", "")
		if tbody == nil {
			continue
		}

		rows := findAllNodes(tbody, "tr", "", "")
		for _, row := range rows {
			cells := findAllNodes(row, "td", "", "")
			if len(cells) < 3 {
				continue
			}

			// Extract competition description from first cell
			competitionDesc := strings.TrimSpace(getTextContent(cells[0]))
			if competitionDesc == "" || strings.Contains(competitionDesc, "Конкурсная единица") {
				continue
			}

			// Parse heading name and competition type from description
			headingName, competitionType := parseCompetitionDescription(competitionDesc)
			if headingName == "" {
				continue
			}

			// Extract URLs from second and third cells
			urls := extractURLsFromCells(cells[1:3])
			if len(urls) == 0 {
				continue
			}

			// Initialize map for heading if not exists
			if result[headingName] == nil {
				result[headingName] = make(map[core.Competition][]string)
			}

			// Add URLs to the appropriate competition type
			result[headingName][competitionType] = append(result[headingName][competitionType], urls...)
		}
	}

	return result, nil
}

// ParseApplicationList parses the MEPhI application list HTML to extract application data.
func ParseApplicationList(doc *html.Node, competitionType core.Competition) ([]*source.ApplicationData, error) {
	var applications []*source.ApplicationData

	// Find the rating table
	table := findNode(doc, "table", "id", "ratingTable")
	if table == nil {
		return nil, fmt.Errorf("rating table not found")
	}

	tbody := findNode(table, "tbody", "", "")
	if tbody == nil {
		return nil, fmt.Errorf("tbody not found in rating table")
	}

	rows := findAllNodes(tbody, "tr", "", "")
	for _, row := range rows {
		// Skip header rows
		if hasClass(row, "throw") {
			continue
		}

		cells := findAllNodes(row, "td", "", "")
		if len(cells) < 10 {
			continue
		}

		// Extract student ID from third cell (№ дела)
		studentIDText := strings.TrimSpace(getTextContent(cells[2]))
		if studentIDText == "" {
			continue
		}

		// Extract scores sum from seventh cell (Сумма баллов)
		scoresSumText := strings.TrimSpace(getTextContent(cells[6]))
		var scoresSum *int
		if scoresSumText != "" && scoresSumText != "-" {
			if sum, err := strconv.Atoi(scoresSumText); err == nil {
				scoresSum = &sum
			}
		}

		// Extract document status from eighth cell (Документы)
		documentStatus := strings.TrimSpace(getTextContent(cells[7]))
		originalSubmitted := strings.Contains(documentStatus, "Согласие подано")

		// Extract priority from tenth cell (Приоритет)
		priorityText := strings.TrimSpace(getTextContent(cells[9]))
		var priority *int
		if priorityText != "" {
			if p, err := strconv.Atoi(priorityText); err == nil {
				priority = &p
			}
		}

		// Determine actual competition type based on scores
		actualCompetitionType := competitionType
		if scoresSum == nil && strings.Contains(getTextContent(cells[6]), "-") {
			// No numeric score, likely BVI (olympiad winner)
			actualCompetitionType = core.CompetitionBVI
		} else {
				// Default to the provided competition type
				actualCompetitionType = competitionType
		}

		// Handle nil pointers
		scoreValue := 0
		if scoresSum != nil {
			scoreValue = *scoresSum
		}
		priorityValue := 1
		if priority != nil {
			priorityValue = *priority
		}

		application := &source.ApplicationData{
			StudentID:         studentIDText,
			ScoresSum:         scoreValue,
			Priority:          priorityValue,
			OriginalSubmitted: originalSubmitted,
			CompetitionType:   actualCompetitionType,
		}
		applications = append(applications, application)
	}

	return applications, nil
}

// parseCompetitionDescription extracts heading name and competition type from MEPhI competition description.
func parseCompetitionDescription(desc string) (string, core.Competition) {
	// Remove extra whitespace
	desc = regexp.MustCompile(`\s+`).ReplaceAllString(desc, " ")
	desc = strings.TrimSpace(desc)

	// Extract heading name (everything before the first comma)
	parts := strings.Split(desc, ",")
	if len(parts) == 0 {
		return "", core.CompetitionRegular
	}

	headingName := strings.TrimSpace(parts[0])

	// Determine competition type based on description
	descLower := strings.ToLower(desc)
	if strings.Contains(descLower, "целевой прием") {
		return headingName, core.CompetitionTargetQuota
	} else if strings.Contains(descLower, "отдельная квота") {
		return headingName, core.CompetitionDedicatedQuota
	} else if strings.Contains(descLower, "особое право") {
		return headingName, core.CompetitionSpecialQuota
	}

	return headingName, core.CompetitionRegular // Default to regular competition
}

// extractURLsFromCells extracts relative URLs from table cells and converts them to absolute URLs.
func extractURLsFromCells(cells []*html.Node) []string {
	var urls []string
	baseURL := "https://org.mephi.ru"

	for _, cell := range cells {
		links := findAllNodes(cell, "a", "", "")
		for _, link := range links {
			href := getAttr(link, "href")
			if href != "" {
				// Convert relative URL to absolute
				if strings.HasPrefix(href, "/") {
					urls = append(urls, baseURL+href)
				} else {
					urls = append(urls, href)
				}
			}
		}
	}

	return urls
}

// Helper functions for HTML parsing

func findNode(root *html.Node, tag, attrName, attrValue string) *html.Node {
	if root.Type == html.ElementNode && root.Data == tag {
		if attrName == "" || getAttr(root, attrName) == attrValue {
			return root
		}
	}

	for child := root.FirstChild; child != nil; child = child.NextSibling {
		if result := findNode(child, tag, attrName, attrValue); result != nil {
			return result
		}
	}

	return nil
}

func findAllNodes(root *html.Node, tag, attrName, attrValue string) []*html.Node {
	var nodes []*html.Node

	if root.Type == html.ElementNode && root.Data == tag {
		if attrName == "" || getAttr(root, attrName) == attrValue {
			nodes = append(nodes, root)
		}
	}

	for child := root.FirstChild; child != nil; child = child.NextSibling {
		nodes = append(nodes, findAllNodes(child, tag, attrName, attrValue)...)
	}

	return nodes
}

func getAttr(node *html.Node, attrName string) string {
	for _, attr := range node.Attr {
		if attr.Key == attrName {
			return attr.Val
		}
	}
	return ""
}

func hasClass(node *html.Node, className string) bool {
	classAttr := getAttr(node, "class")
	classes := strings.Fields(classAttr)
	for _, class := range classes {
		if class == className {
			return true
		}
	}
	return false
}

func getTextContent(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}

	var text strings.Builder
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		text.WriteString(getTextContent(child))
	}

	return text.String()
}