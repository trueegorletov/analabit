package itmo

import (
	"io"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/html"
)

func TestDebugCapacityExtraction(t *testing.T) {
	url := "https://abit.itmo.ru/rating/bachelor/budget/2190"

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		t.Fatalf("Failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	bodySize := len(body)
	if bodySize > 1000 {
		bodySize = 1000
	}
	t.Logf("First %d chars of body: %s", bodySize, string(body[:bodySize]))

	// Parse HTML
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	// Extract raw text
	text := getTextContent(doc)
	textSize := len(text)
	if textSize > 1000 {
		textSize = 1000
	}
	t.Logf("First %d chars of extracted text: %s", textSize, text[:textSize])

	// Look for capacity-related text
	capacityIndex := strings.Index(text, "Количество мест")
	if capacityIndex >= 0 {
		start := capacityIndex
		end := capacityIndex + 200
		if end > len(text) {
			end = len(text)
		}
		t.Logf("Capacity section: %s", text[start:end])
	} else {
		t.Logf("Capacity text not found in extracted text")
	}

	// Test the regex directly
	matches := capacityRegex.FindStringSubmatch(text)
	t.Logf("Regex matches: %v", matches)

	// Also test some variations
	capacityRegexDebug := regexp.MustCompile(`Количество мест:.*?(\d+).*?\(.*?(\d+).*?ЦК.*?(\d+).*?ОcК.*?(\d+).*?ОтК.*?\)`)
	debugMatches := capacityRegexDebug.FindStringSubmatch(text)
	t.Logf("Debug regex matches: %v", debugMatches)

	// Extract capacities
	capacities := extractCapacities(doc)
	t.Logf("Extracted capacities: %+v", capacities)

	total := capacities.Regular + capacities.TargetQuota + capacities.SpecialQuota + capacities.DedicatedQuota
	t.Logf("Total capacity: %d", total)
}
