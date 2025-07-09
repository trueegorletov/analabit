// Package itmo provides parsing functionality for ITMO University admission lists.
// It implements the HeadingSource interface to extract application data from
// ITMO's competitive lists at https://abit.itmo.ru/ratings/bachelor.
package itmo

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"golang.org/x/net/html"
)

// Common patterns and mappings for ITMO parsing
var (
	// Competition type mapping for ITMO section headings
	competitionTypeMap = map[string]core.Competition{
		"Без вступительных испытаний": core.CompetitionBVI,
		"Целевая квота":               core.CompetitionTargetQuota,
		"Особая квота":                core.CompetitionSpecialQuota,
		"Отдельная квота":             core.CompetitionDedicatedQuota,
		"Общий конкурс":               core.CompetitionRegular,
	}

	// Regex patterns for data extraction
	positionAndIDRegex = regexp.MustCompile(`(\d+)\s*№\s*(\d+)`)
	capacityRegex      = regexp.MustCompile(`Количество мест:\s*(\d+)\s*\(\s*(\d+)\s*ЦК\s*,\s*(\d+)\s*ОcК\s*,\s*(\d+)\s*ОтК\s*\)`)
	priorityRegex      = regexp.MustCompile(`Приоритет:\s*(\d+)`)
	individualRegex    = regexp.MustCompile(`ИД:\s*(\d+)`)
	totalScoreRegex    = regexp.MustCompile(`Балл ВИ\+ИД:\s*(\d+)`)
	mathRegex          = regexp.MustCompile(`Математика:\s*(\d+)`)
	csRegex            = regexp.MustCompile(`Информатика:\s*(\d+)`)
	russianRegex       = regexp.MustCompile(`Русский язык:\s*(\d+)`)
)

// extractCapacities parses the capacity information from ITMO HTML
// It looks for text like "Количество мест: 170 (17 ЦК, 17 ОcК, 17 ОтК)"
// If parsing fails, uses fallback capacities with 10% for each quota type
func extractCapacities(doc *html.Node) core.Capacities {
	text := getTextContent(doc)
	matches := capacityRegex.FindStringSubmatch(text)

	if len(matches) >= 5 {
		total, _ := strconv.Atoi(matches[1])
		target, _ := strconv.Atoi(matches[2])    // ЦК
		special, _ := strconv.Atoi(matches[3])   // ОcК
		dedicated, _ := strconv.Atoi(matches[4]) // ОтК

		regular := total - target - special - dedicated
		if regular < 0 {
			regular = 0
		}

		return core.Capacities{
			Regular:        regular,
			TargetQuota:    target,
			SpecialQuota:   special,
			DedicatedQuota: dedicated,
		}
	}

	// Parsing failed - return zero capacities to trigger fallback in LoadTo
	return core.Capacities{}
}

// CalculateFallbackCapacities computes quota distribution when parsing fails
// Gives 10% each to target/special/dedicated quotas, 70% to regular
func CalculateFallbackCapacities(totalKCP int) core.Capacities {
	if totalKCP <= 0 {
		return core.Capacities{}
	}

	// Each quota gets 10% (rounded up)
	targetQuota := (totalKCP + 9) / 10    // 10% rounded up
	specialQuota := (totalKCP + 9) / 10   // 10% rounded up
	dedicatedQuota := (totalKCP + 9) / 10 // 10% rounded up

	// Regular gets the remainder (approximately 70%)
	regular := totalKCP - targetQuota - specialQuota - dedicatedQuota
	if regular < 0 {
		regular = 0
	}

	return core.Capacities{
		Regular:        regular,
		TargetQuota:    targetQuota,
		SpecialQuota:   specialQuota,
		DedicatedQuota: dedicatedQuota,
	}
}

// parseApplicant extracts application data from an ITMO application HTML node
func parseApplicant(node *html.Node, competitionType core.Competition) (*source.ApplicationData, error) {
	text := getTextContent(node)

	// Extract position and application ID
	matches := positionAndIDRegex.FindStringSubmatch(text)
	if len(matches) < 3 {
		return nil, fmt.Errorf("could not extract position and ID from application")
	}

	position, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("invalid position number: %v", err)
	}

	studentID := matches[2]

	// Extract priority (default to 1 if not found, which is common for ITMO)
	priority := 1
	if priorityMatches := priorityRegex.FindStringSubmatch(text); len(priorityMatches) > 1 {
		priority, _ = strconv.Atoi(priorityMatches[1])
	}

	// Extract scores - for ITMO we need to calculate total score
	var scoresSum int

	// Try to get total score directly first
	if totalMatches := totalScoreRegex.FindStringSubmatch(text); len(totalMatches) > 1 {
		scoresSum, _ = strconv.Atoi(totalMatches[1])
	} else {
		// If no total score, try to calculate from individual subjects
		var mathScore, csScore, russianScore, individualScore int

		if mathMatches := mathRegex.FindStringSubmatch(text); len(mathMatches) > 1 {
			mathScore, _ = strconv.Atoi(mathMatches[1])
		}
		if csMatches := csRegex.FindStringSubmatch(text); len(csMatches) > 1 {
			csScore, _ = strconv.Atoi(csMatches[1])
		}
		if russianMatches := russianRegex.FindStringSubmatch(text); len(russianMatches) > 1 {
			russianScore, _ = strconv.Atoi(russianMatches[1])
		}
		if indMatches := individualRegex.FindStringSubmatch(text); len(indMatches) > 1 {
			individualScore, _ = strconv.Atoi(indMatches[1])
		}

		scoresSum = mathScore + csScore + russianScore + individualScore
	}

	// Extract consent status
	originalSubmitted := strings.Contains(strings.ToLower(text), "согласие: да")

	return &source.ApplicationData{
		StudentID:         studentID,
		RatingPlace:       position,
		Priority:          priority,
		CompetitionType:   competitionType,
		ScoresSum:         scoresSum,
		OriginalSubmitted: originalSubmitted,
	}, nil
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

// extractPrettyName extracts the program name from HTML title or heading
func extractPrettyName(doc *html.Node) string {
	// Look for the program title in h2 element
	var title string
	var findTitle func(*html.Node)
	findTitle = func(n *html.Node) {
		if title != "" {
			return
		}
		if n.Type == html.ElementNode && n.Data == "h2" {
			text := getTextContent(n)
			if text != "" && strings.Contains(text, "«") && strings.Contains(text, "»") {
				title = text
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findTitle(c)
		}
	}
	findTitle(doc)

	if title != "" {
		// Extract just the program name in quotes
		start := strings.Index(title, "«")
		end := strings.Index(title, "»")
		if start != -1 && end != -1 && end > start {
			// Extract the text between the quotes, handling UTF-8 properly
			return title[start+len("«") : end]
		}
		return title
	}

	return ""
}
