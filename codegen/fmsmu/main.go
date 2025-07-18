//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// Local type definitions for codegen
type SpoilerData struct {
	ID              string
	HeadingName     string
	Capacity        int
	CompetitionType string
	Organization    string
}

type HeadingGroup struct {
	Spoilers    []SpoilerData
	Capacities  Capacities
	CleanedName string
}

type Capacities struct {
	Total          int
	Regular        int
	TargetQuota    int
	SpecialQuota   int
	DedicatedQuota int
}

type HTTPHeadingSource struct {
	PrettyName           string
	RegularListID        string
	SpecialQuotaListID   string
	DedicatedQuotaListID string
	TargetQuotaListIDs   []string
	Capacities           Capacities
}

// Competition type constants
const (
	CompetitionRegular        = "regular"
	CompetitionTargetQuota    = "target"
	CompetitionSpecialQuota   = "special"
	CompetitionDedicatedQuota = "dedicated"
)

func main() {
	// Fetch all registry pages
	allSpoilers, err := fetchAllRegistryPages()
	if err != nil {
		log.Fatalf("Error fetching registry pages: %v", err)
	}

	// Generate source definitions
	sources := generateSources(allSpoilers)

	// Output as Go code
	outputGoCode(sources)
}

// fetchAllRegistryPages fetches spoilers from all registry pages
func fetchAllRegistryPages() ([]SpoilerData, error) {
	var allSpoilers []SpoilerData
	page := 1

	for {
		spoilers, err := fetchRegistryPage(page)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch page %d: %w", page, err)
		}

		if len(spoilers) == 0 {
			break
		}

		allSpoilers = append(allSpoilers, spoilers...)

		page++
		// Safety check to prevent infinite loops
		if page > 100 {
			break
		}
	}

	return allSpoilers, nil
}

// generateSources converts spoilers to HTTPHeadingSource definitions
func generateSources(spoilers []SpoilerData) []HTTPHeadingSource {
	// Group spoilers by heading name first
	groups := groupByHeading(spoilers)

	var sources []HTTPHeadingSource

	// Create one HTTPHeadingSource per group
	for _, group := range groups {
		httpSource := HTTPHeadingSource{
			PrettyName: group.CleanedName,
			Capacities: group.Capacities,
		}

		// Set list IDs based on spoilers in the group
		for _, spoiler := range group.Spoilers {
			switch spoiler.CompetitionType {
			case CompetitionRegular:
				httpSource.RegularListID = spoiler.ID
			case CompetitionSpecialQuota:
				httpSource.SpecialQuotaListID = spoiler.ID
			case CompetitionDedicatedQuota:
				httpSource.DedicatedQuotaListID = spoiler.ID
			case CompetitionTargetQuota:
				httpSource.TargetQuotaListIDs = append(httpSource.TargetQuotaListIDs, spoiler.ID)
			}
		}

		sources = append(sources, httpSource)
	}

	return sources
}

// outputGoCode outputs the sources as Go code
func outputGoCode(sources []HTTPHeadingSource) {
	for _, httpSrc := range sources {
		fmt.Printf("&fmsmu.HTTPHeadingSource{\n")
		fmt.Printf("\tPrettyName: %s,\n", formatGoString(httpSrc.PrettyName))

		if httpSrc.RegularListID != "" {
			fmt.Printf("\tRegularListID: %s,\n", formatGoString(httpSrc.RegularListID))
		}
		if httpSrc.SpecialQuotaListID != "" {
			fmt.Printf("\tSpecialQuotaListID: %s,\n", formatGoString(httpSrc.SpecialQuotaListID))
		}
		if httpSrc.DedicatedQuotaListID != "" {
			fmt.Printf("\tDedicatedQuotaListID: %s,\n", formatGoString(httpSrc.DedicatedQuotaListID))
		}
		if len(httpSrc.TargetQuotaListIDs) > 0 {
			fmt.Printf("\tTargetQuotaListIDs: []string{\n")
			for _, id := range httpSrc.TargetQuotaListIDs {
				fmt.Printf("\t\t%s,\n", formatGoString(id))
			}
			fmt.Printf("\t},\n")
		}

		fmt.Printf("\tCapacities: source.Capacities{\n")
		if httpSrc.Capacities.Regular > 0 {
			fmt.Printf("\t\tRegular: %d,\n", httpSrc.Capacities.Regular)
		}
		if httpSrc.Capacities.TargetQuota > 0 {
			fmt.Printf("\t\tTargetQuota: %d,\n", httpSrc.Capacities.TargetQuota)
		}
		if httpSrc.Capacities.SpecialQuota > 0 {
			fmt.Printf("\t\tSpecialQuota: %d,\n", httpSrc.Capacities.SpecialQuota)
		}
		if httpSrc.Capacities.DedicatedQuota > 0 {
			fmt.Printf("\t\tDedicatedQuota: %d,\n", httpSrc.Capacities.DedicatedQuota)
		}
		fmt.Printf("\t},\n")
		fmt.Printf("},\n")
		fmt.Println()
	}
}

// fetchRegistryPage fetches spoilers from a single registry page
func fetchRegistryPage(pageNum int) ([]SpoilerData, error) {
	url := fmt.Sprintf("https://priem.sechenov.ru/submitted-applicants/?search=&search_terms=&type-of-financing%%5B0%%5D=57&page=page-%d", pageNum)

	// Create request with proper headers
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for page %d: %w", pageNum, err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page %d: %w", pageNum, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d for page %d", resp.StatusCode, pageNum)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML for page %d: %w", pageNum, err)
	}

	spoilers := extractSpoilers(doc)
	return spoilers, nil
}

// extractSpoilers finds and parses spoiler elements from HTML
func extractSpoilers(doc *html.Node) []SpoilerData {
	var spoilers []SpoilerData

	// Find spoiler-box div
	spoilerBox := findElementByClass(doc, "spoiler-box")
	if spoilerBox == nil {
		return spoilers
	}

	// Find all spoiler divs within spoiler-box
	var findSpoilers func(*html.Node)
	findSpoilers = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "div" {
			class := getAttr(node, "class")
			if strings.Contains(class, "spoiler") && strings.Contains(class, "spoiler--great-spoiler") {
				if spoiler := parseSpoiler(node); spoiler != nil {
					spoilers = append(spoilers, *spoiler)
				}
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			findSpoilers(c)
		}
	}

	findSpoilers(spoilerBox)
	return spoilers
}

// parseSpoiler extracts data from a single spoiler element
func parseSpoiler(spoilerNode *html.Node) *SpoilerData {
	// Extract data-group-id
	id := getAttr(spoilerNode, "data-group-id")
	if id == "" {
		return nil
	}

	// Extract heading name
	headingName := extractHeadingName(spoilerNode)
	if headingName == "" {
		return nil
	}

	// Extract capacity
	capacity := extractCapacity(spoilerNode)

	// Determine competition type
	competitionType, _ := determineCompetitionType(spoilerNode)

	return &SpoilerData{
		ID:              id,
		HeadingName:     headingName,
		Capacity:        capacity,
		CompetitionType: competitionType,
		Organization:    "FMSMU",
	}
}

// extractHeadingName extracts the heading name from spoiler
func extractHeadingName(spoilerNode *html.Node) string {
	// Try to find the course name element first
	courseNameNode := findElementByClass(spoilerNode, "spoiler__course-name")
	if courseNameNode == nil {
		// Fallback: extract from the entire spoiler content
		courseNameNode = spoilerNode
	}

	text := getTextContent(courseNameNode)

	// First clean whitespace to normalize the text
	re := regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)

	// Remove state-code prefix using regex (XX.XX.XX)
	re = regexp.MustCompile(`^\d{2}\.\d{2}\.\d{2}\s+`)
	text = re.ReplaceAllString(text, "")

	// Remove capacity suffix ("Количество мест: N" and everything after)
	re = regexp.MustCompile(`\s*Количество мест:.*$`)
	text = re.ReplaceAllString(text, "")

	return strings.TrimSpace(text)
}

// extractCapacity extracts capacity from spoiler
func extractCapacity(spoilerNode *html.Node) int {
	text := getTextContent(spoilerNode)

	// Look for "Количество мест: N"
	re := regexp.MustCompile(`Количество мест:\s*(\d+)`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		if capacity, err := strconv.Atoi(matches[1]); err == nil {
			return capacity
		}
	}

	return 0
}

// determineCompetitionType determines competition type from spoiler content
func determineCompetitionType(spoilerNode *html.Node) (string, string) {
	text := strings.ToLower(getTextContent(spoilerNode))

	if strings.Contains(text, "целевой") {
		return CompetitionTargetQuota, "FMSMU"
	}
	if strings.Contains(text, "особая") {
		return CompetitionSpecialQuota, "FMSMU"
	}
	if strings.Contains(text, "отдельная") {
		return CompetitionDedicatedQuota, "FMSMU"
	}

	return CompetitionRegular, "FMSMU"
}

// groupByHeading groups spoilers by heading name and computes capacities
func groupByHeading(spoilers []SpoilerData) map[string]HeadingGroup {
	groups := make(map[string]HeadingGroup)

	for _, spoiler := range spoilers {
		// Use the cleaned heading name as the key for grouping
		cleanedName := spoiler.HeadingName
		group := groups[cleanedName]
		group.Spoilers = append(group.Spoilers, spoiler)

		// Store the cleaned name in the group for later use
		group.CleanedName = cleanedName

		// Add to appropriate capacity
		switch spoiler.CompetitionType {
		case CompetitionRegular:
			group.Capacities.Regular += spoiler.Capacity
		case CompetitionTargetQuota:
			group.Capacities.TargetQuota += spoiler.Capacity
		case CompetitionSpecialQuota:
			group.Capacities.SpecialQuota += spoiler.Capacity
		case CompetitionDedicatedQuota:
			group.Capacities.DedicatedQuota += spoiler.Capacity
		}

		group.Capacities.Total = group.Capacities.Regular + group.Capacities.TargetQuota + group.Capacities.SpecialQuota + group.Capacities.DedicatedQuota

		groups[cleanedName] = group
	}

	return groups
}

// Helper functions for HTML parsing
func findElementByClass(node *html.Node, className string) *html.Node {
	if node.Type == html.ElementNode {
		for _, attr := range node.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, className) {
				return node
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if result := findElementByClass(c, className); result != nil {
			return result
		}
	}

	return nil
}

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
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

// formatGoString formats a string for Go code output
func formatGoString(s string) string {
	// Use JSON marshaling to properly escape the string
	b, _ := json.Marshal(s)
	return string(b)
}
