package itmo

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/utils"
	"golang.org/x/net/html"
)

// HTTPHeadingSource implements source.HeadingSource for ITMO University
type HTTPHeadingSource struct {
	URL        string          `json:"url"`
	PrettyName string          `json:"pretty_name"` // Fallback if parsing fails
	Capacities core.Capacities `json:"capacities"`  // May be updated by parser
}

// LoadTo fetches and parses ITMO application data from the configured URL
func (h *HTTPHeadingSource) LoadTo(receiver source.DataReceiver) error {
	// Fetch the HTML content
	ctx := context.Background()
	release, err := source.AcquireHTTPSemaphores(ctx, "itmo")
	if err != nil {
		return fmt.Errorf("failed to acquire semaphores for ITMO list from %s: %v", h.URL, err)
	}
	defer release()

	resp, err := http.Get(h.URL)
	if err != nil {
		return fmt.Errorf("failed to fetch URL %s: %v", h.URL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status %d for URL %s", resp.StatusCode, h.URL)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the HTML
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Extract capacities from the page content
	parsedCapacities := extractCapacities(doc)

	// Check if parsing succeeded (total > 0) or use fallback
	totalParsed := parsedCapacities.Regular + parsedCapacities.TargetQuota +
		parsedCapacities.SpecialQuota + parsedCapacities.DedicatedQuota
	if totalParsed == 0 {
		// Parsing failed, use fallback calculation from pre-set КЦП total
		totalKCP := h.Capacities.Regular // This contains the total КЦП from list page
		if totalKCP > 0 {
			fmt.Printf("Warning: Could not parse capacity details from %s, using fallback calculation for КЦП=%d\n", h.URL, totalKCP)
			h.Capacities = CalculateFallbackCapacities(totalKCP)
		} else {
			fmt.Printf("Warning: No capacity information available for %s\n", h.URL)
		}
	} else {
		// Use parsed capacities
		h.Capacities = parsedCapacities
	}

	// Extract program name if not already set
	prettyName := h.PrettyName
	if prettyName == "" {
		prettyName = extractPrettyName(doc)
	}
	if prettyName == "" {
		return fmt.Errorf("could not extract program name from URL %s", h.URL)
	}

	// Send heading data
	headingData := &source.HeadingData{
		Code:       utils.GenerateHeadingCode(prettyName),
		PrettyName: prettyName,
		Capacities: h.Capacities,
	}
	receiver.PutHeadingData(headingData)

	// Parse and send applications by category
	err = h.parseApplicationsByCategory(doc, receiver, headingData.Code)
	if err != nil {
		return fmt.Errorf("failed to parse applications: %v", err)
	}

	return nil
}

// parseApplicationsByCategory extracts applications from each competition category
func (h *HTTPHeadingSource) parseApplicationsByCategory(doc *html.Node, receiver source.DataReceiver, headingCode string) error {
	var currentCategory string
	var extractData func(*html.Node)

	extractData = func(n *html.Node) {
		if n.Type == html.ElementNode {
			// Check if this is a category title (h5 with specific class)
			if n.Data == "h5" {
				for _, attr := range n.Attr {
					if attr.Key == "class" && strings.Contains(attr.Val, "title__zlsGy") {
						categoryText := getTextContent(n)
						if categoryText != "" {
							currentCategory = categoryText
						}
						break
					}
				}
			}

			// Check if this is an application item
			if currentCategory != "" {
				for _, attr := range n.Attr {
					if attr.Key == "class" && strings.Contains(attr.Val, "table__item__qMY0F") {
						// Parse this application
						competitionType, exists := competitionTypeMap[currentCategory]
						if !exists {
							// Skip unknown competition types
							break
						}

						app, err := parseApplicant(n, competitionType)
						if err != nil {
							// Log error but continue processing
							fmt.Printf("Warning: failed to parse application in category %s: %v\n", currentCategory, err)
							break
						}

						// Set the heading code for this application
						app.HeadingCode = headingCode

						// Send the application data
						receiver.PutApplicationData(app)
						break
					}
				}
			}
		}

		// Recursively process all child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractData(c)
		}
	}

	extractData(doc)
	return nil
}
