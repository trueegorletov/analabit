// Package fmsmu provides support for loading FMSMU (First Moscow State Medical University) admission data.
// FMSMU admission lists are provided in HTML format with paginated registry and individual list pages.
package fmsmu

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/trueegorletov/analabit/core"
	"golang.org/x/net/html"
)

// SpoilerData represents parsed data from a registry spoiler element
type SpoilerData struct {
	ID              string
	HeadingName     string
	Capacity        int
	CompetitionType core.Competition
	Organization    string // For target quota lists
}

// HeadingGroup represents grouped spoilers by heading name
type HeadingGroup struct {
	HeadingName           string
	RegularListID         string
	SpecialQuotaListID    string
	DedicatedQuotaListID  string
	TargetQuotaListIDs    []string
	Spoilers              []SpoilerData
	Capacities            core.Capacities
}

// fetchRegistryPage fetches a single registry page
func fetchRegistryPage(ctx context.Context, pageURL string) (*html.Node, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", pageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch registry page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch registry page (status code %d)", resp.StatusCode)
	}

	htmlContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTML content: %w", err)
	}

	doc, err := html.Parse(strings.NewReader(string(htmlContent)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	return doc, nil
}

// extractSpoilers extracts spoiler data from a registry page HTML document
func extractSpoilers(doc *html.Node) ([]SpoilerData, error) {
	var spoilers []SpoilerData

	// Find the spoiler-box div
	spoilerBox := findElementByClass(doc, "spoiler-box")
	if spoilerBox == nil {
		return nil, fmt.Errorf("spoiler-box not found")
	}

	// Find all spoiler divs within the spoiler-box
	var findSpoilers func(*html.Node)
	findSpoilers = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			if hasClass(n, "spoiler") {
				if spoiler := parseSpoiler(n); spoiler != nil {
					spoilers = append(spoilers, *spoiler)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findSpoilers(c)
		}
	}
	findSpoilers(spoilerBox)

	return spoilers, nil
}

// parseSpoiler parses a single spoiler div element
func parseSpoiler(spoilerDiv *html.Node) *SpoilerData {
	// Extract data-group-id
	id := getAttr(spoilerDiv, "data-group-id")
	if id == "" {
		return nil
	}

	// Find spoiler__course-name span
	courseNameSpan := findElementByClass(spoilerDiv, "spoiler__course-name")
	if courseNameSpan == nil {
		return nil
	}

	// Extract heading name (strip code prefix)
	courseText := getTextContent(courseNameSpan)
	headingName := extractHeadingName(courseText)
	if headingName == "" {
		return nil
	}

	// Extract capacity from admission-plan span
	admissionPlanSpan := findElementByClass(courseNameSpan, "admission-plan")
	if admissionPlanSpan == nil {
		return nil
	}
	capacity := extractCapacity(getTextContent(admissionPlanSpan))

	// Determine competition type and organization
	competitionType, organization := determineCompetitionType(spoilerDiv)

	return &SpoilerData{
		ID:              id,
		HeadingName:     headingName,
		Capacity:        capacity,
		CompetitionType: competitionType,
		Organization:    organization,
	}
}

// extractHeadingName extracts heading name from course text, removing code prefix
func extractHeadingName(courseText string) string {
	// Use regex to match pattern like "37.05.01 Клиническая психология"
	re := regexp.MustCompile(`^\d+\.\d+\.\d+\s+(.+?)\s*Количество мест:`)
	matches := re.FindStringSubmatch(courseText)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

// extractCapacity extracts capacity number from admission plan text
func extractCapacity(admissionText string) int {
	// Extract number after "Количество мест:"
	re := regexp.MustCompile(`Количество мест:\s*(\d+)`)
	matches := re.FindStringSubmatch(admissionText)
	if len(matches) > 1 {
		if capacity, err := strconv.Atoi(matches[1]); err == nil {
			return capacity
		}
	}
	return 0
}

// determineCompetitionType determines competition type from form-education spans
func determineCompetitionType(spoilerDiv *html.Node) (core.Competition, string) {
	// Find form-education-box
	formEducationBox := findElementByClass(spoilerDiv, "form-education-box")
	if formEducationBox == nil {
		return core.CompetitionRegular, ""
	}

	// Find all form-education spans with accent class
	accentSpans := findElementsByClass(formEducationBox, "form-education--accent")
	var organization string

	for _, span := range accentSpans {
		text := strings.TrimSpace(getTextContent(span))
		switch text {
		case "Особая квота":
			return core.CompetitionSpecialQuota, ""
		case "Отдельная квота":
			return core.CompetitionDedicatedQuota, ""
		case "Целевой прием":
			// Continue to find organization
			competitionType := core.CompetitionTargetQuota
			// Look for organization name in subsequent spans
			for _, orgSpan := range accentSpans {
				orgText := strings.TrimSpace(getTextContent(orgSpan))
				if orgText != "Целевой прием" && orgText != "Бюджетное" && orgText != "Очная" {
					organization = orgText
					break
				}
			}
			return competitionType, organization
		}
	}

	return core.CompetitionRegular, ""
}

// groupByHeading groups spoilers by heading name
func groupByHeading(spoilers []SpoilerData) map[string]*HeadingGroup {
	groups := make(map[string]*HeadingGroup)

	for _, spoiler := range spoilers {
		group, exists := groups[spoiler.HeadingName]
		if !exists {
			group = &HeadingGroup{
				HeadingName: spoiler.HeadingName,
				Capacities:  core.Capacities{},
			}
			groups[spoiler.HeadingName] = group
		}

		switch spoiler.CompetitionType {
		case core.CompetitionRegular:
			group.RegularListID = spoiler.ID
			group.Capacities.Regular = spoiler.Capacity
		case core.CompetitionSpecialQuota:
			group.SpecialQuotaListID = spoiler.ID
			group.Capacities.SpecialQuota = spoiler.Capacity
		case core.CompetitionDedicatedQuota:
			group.DedicatedQuotaListID = spoiler.ID
			group.Capacities.DedicatedQuota = spoiler.Capacity
		case core.CompetitionTargetQuota:
			group.TargetQuotaListIDs = append(group.TargetQuotaListIDs, spoiler.ID)
			group.Capacities.TargetQuota += spoiler.Capacity
		}
	}

	return groups
}

// GroupByHeading groups spoilers by heading name and computes total capacities
func GroupByHeading(spoilers []SpoilerData) map[string]HeadingGroup {
	groups := make(map[string]HeadingGroup)

	for _, spoiler := range spoilers {
		group, exists := groups[spoiler.HeadingName]
		if !exists {
			group = HeadingGroup{
				HeadingName: spoiler.HeadingName,
				Spoilers:    []SpoilerData{},
				Capacities:  core.Capacities{},
			}
		}

		// Add spoiler to group
		group.Spoilers = append(group.Spoilers, spoiler)

		// Add capacity based on competition type
		switch spoiler.CompetitionType {
		case core.CompetitionRegular:
			group.Capacities.Regular += spoiler.Capacity
		case core.CompetitionTargetQuota:
			group.Capacities.TargetQuota += spoiler.Capacity
		case core.CompetitionSpecialQuota:
			group.Capacities.SpecialQuota += spoiler.Capacity
		case core.CompetitionDedicatedQuota:
			group.Capacities.DedicatedQuota += spoiler.Capacity
		}

		groups[spoiler.HeadingName] = group
	}

	return groups
}

// Helper functions for HTML parsing

func findElementByClass(node *html.Node, className string) *html.Node {
	if node.Type == html.ElementNode && hasClass(node, className) {
		return node
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if result := findElementByClass(c, className); result != nil {
			return result
		}
	}
	return nil
}

func findElementsByClass(node *html.Node, className string) []*html.Node {
	var results []*html.Node
	var find func(*html.Node)
	find = func(n *html.Node) {
		if n.Type == html.ElementNode && hasClass(n, className) {
			results = append(results, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			find(c)
		}
	}
	find(node)
	return results
}

func hasClass(node *html.Node, className string) bool {
	for _, attr := range node.Attr {
		if attr.Key == "class" {
			classes := strings.Fields(attr.Val)
			for _, class := range classes {
				if class == className {
					return true
				}
			}
		}
	}
	return false
}

func getAttr(node *html.Node, attrName string) string {
	for _, attr := range node.Attr {
		if attr.Key == attrName {
			return attr.Val
		}
	}
	return ""
}

func getTextContent(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}
	var text strings.Builder
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		text.WriteString(getTextContent(c))
	}
	return text.String()
}