package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/trueegorletov/analabit/core"
	"golang.org/x/net/html"
)

type ApplicationData struct {
	StudentID         string
	OriginalSubmitted bool
}

func testOriginalDetection() {
	// Test URL for MIPT Прикладная математика и информатика
	testURL := "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJpa2xhZG5heWEgbWF0ZW1hdGlrYSBpIGluZm9ybWF0aWthX0J5dWR6aGV0X05hIG9ic2hjaGlraCBvc25vdmFuaXlha2guaHRtbA=="

	fmt.Printf("Testing MIPT original document detection fix with URL:\n%s\n\n", testURL)

	// Fetch the HTML content
	resp, err := http.Get(testURL)
	if err != nil {
		log.Fatalf("Failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	// Parse HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	// Find all table rows
	rows := findTableRows(doc)
	fmt.Printf("Found %d table rows\n\n", len(rows))

	// Test the first 50 student entries
	fmt.Println("Testing first 50 students:")
	fmt.Println("Student ID\tOriginal Submitted")
	fmt.Println("----------\t------------------")

	count := 0
	for _, row := range rows {
		if count >= 50 {
			break
		}

		// Try to parse the row as applicant data
		appData, err := parseApplicantFromTableRow(row, core.CompetitionRegular)
		if err != nil {
			// Skip rows that can't be parsed (headers, etc.)
			continue
		}

		// Log the student ID and original submission status
		fmt.Printf("%s\t\t%t\n", appData.StudentID, appData.OriginalSubmitted)
		count++
	}

	fmt.Printf("\nTested %d student records\n", count)
}

// Helper functions from common.go that we need for the test

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

func parseApplicantFromTableRow(row *html.Node, defaultCompetitionType core.Competition) (*ApplicationData, error) {
	cells := extractTableCells(row)
	if len(cells) < 9 {
		return nil, fmt.Errorf("insufficient columns in row: got %d, expected at least 9", len(cells))
	}

	// Column 1: Student ID (e.g., "4102004")
	studentIDText := strings.TrimSpace(getTextContent(cells[1]))
	if len(studentIDText) != 7 || !isNumeric(studentIDText) {
		return nil, fmt.Errorf("invalid student ID format: %s", studentIDText)
	}

	// Column 9: "Согласие на зачисление" (consent for enrollment) - look for checkmark "✓"
	originalSubmitted := false
	if len(cells) > 9 {
		consentText := strings.TrimSpace(getTextContent(cells[9]))
		originalSubmitted = strings.Contains(consentText, "✓")
	}

	return &ApplicationData{
		StudentID:         studentIDText,
		OriginalSubmitted: originalSubmitted,
	}, nil
}

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

func isNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
