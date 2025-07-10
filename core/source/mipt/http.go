package mipt

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/utils"
	"golang.org/x/net/html"
)

// HTTPHeadingSource loads MIPT heading data from multiple HTML list URLs.
type HTTPHeadingSource struct {
	PrettyName            string          `json:"pretty_name"`
	RegularBVIListURL     string          `json:"regular_bvi_list_url"`     // Combined Regular+BVI list
	TargetQuotaListURLs   []string        `json:"target_quota_list_urls"`   // Multiple target quota lists
	DedicatedQuotaListURL string          `json:"dedicated_quota_list_url"` // Dedicated quota list
	SpecialQuotaListURL   string          `json:"special_quota_list_url"`   // Special quota list
	Capacities            core.Capacities `json:"capacities"`
}

// fetchMiptListByURL fetches and parses the MIPT HTML list from a URL.
func fetchMiptListByURL(listURL string, competitionType core.Competition) ([]*source.ApplicationData, error) {
	if listURL == "" {
		return nil, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	release, err := source.AcquireHTTPSemaphores(ctx, "mipt")
	if err != nil {
		return nil, fmt.Errorf("failed to acquire semaphores for %s: %w", listURL, err)
	}
	defer release()

	resp, err := http.Get(listURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download %s: %w", listURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download %s (status code %d)", listURL, resp.StatusCode)
	}

	// Read the HTML content
	htmlContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTML content from %s: %w", listURL, err)
	}

	// Parse the HTML
	doc, err := html.Parse(strings.NewReader(string(htmlContent)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML from %s: %w", listURL, err)
	}

	// Extract applications from the table
	return parseApplicationsFromHTML(doc, competitionType)
}

// parseApplicationsFromHTML extracts applications from MIPT HTML document
func parseApplicationsFromHTML(doc *html.Node, defaultCompetitionType core.Competition) ([]*source.ApplicationData, error) {
	// Find all table rows
	rows := findTableRows(doc)

	if len(rows) == 0 {
		return nil, fmt.Errorf("no table rows found in MIPT HTML")
	}

	var applications []*source.ApplicationData
	var errors []string

	// Process each row
	for i, row := range rows {
		// Skip header rows
		if isHeaderRowFromHTML(row) {
			continue
		}

		// Skip empty rows
		cells := extractTableCells(row)
		if len(cells) == 0 {
			continue
		}

		// Check if this looks like a data row (first cell should be a position number)
		firstCellText := strings.TrimSpace(getTextContent(cells[0]))
		if firstCellText == "" || !positionRegex.MatchString(firstCellText) {
			continue
		}

		// Parse the applicant data
		app, err := parseApplicantFromTableRow(row, defaultCompetitionType)
		if err != nil {
			errors = append(errors, fmt.Sprintf("row %d: %v", i+1, err))
			continue
		}

		applications = append(applications, app)
	}

	log.Printf("MIPT HTTP parser processed %d applications", len(applications))

	if len(errors) > 0 {
		log.Printf("MIPT HTTP parser encountered %d errors:", len(errors))
		for _, err := range errors {
			log.Printf("  - %s", err)
		}
	}

	return applications, nil
}

// isHeaderRowFromHTML determines if a table row is a header row
func isHeaderRowFromHTML(row *html.Node) bool {
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

// LoadTo implements source.HeadingSource for HTTPHeadingSource.
func (s *HTTPHeadingSource) LoadTo(receiver source.DataReceiver) error {
	if s.PrettyName == "" {
		return fmt.Errorf("PrettyName is required for MIPT HTTPHeadingSource")
	}

	headingCode := utils.GenerateHeadingCode(s.PrettyName)
	receiver.PutHeadingData(&source.HeadingData{
		Code:       headingCode,
		Capacities: s.Capacities,
		PrettyName: s.PrettyName,
	})

	// Define list configurations
	listConfigs := []struct {
		URL         string
		Competition core.Competition
		ListName    string
	}{
		{s.RegularBVIListURL, core.CompetitionRegular, "Regular+BVI List"}, // Note: parseApplicantFromTableRow will detect BVI vs Regular per row
		{s.DedicatedQuotaListURL, core.CompetitionDedicatedQuota, "Dedicated Quota List"},
		{s.SpecialQuotaListURL, core.CompetitionSpecialQuota, "Special Quota List"},
	}

	// Process individual lists
	for _, config := range listConfigs {
		if config.URL == "" {
			continue
		}

		applications, err := fetchMiptListByURL(config.URL, config.Competition)
		if err != nil {
			log.Printf("Error fetching %s (%s): %v", config.ListName, config.URL, err)
			continue
		}

		if applications == nil {
			log.Printf("No applications found in %s (%s)", config.ListName, config.URL)
			continue
		}

		// Set heading code and send applications
		for _, app := range applications {
			app.HeadingCode = headingCode
			receiver.PutApplicationData(app)
		}
	}

	// Handle multiple target quota list URLs
	for i, listURL := range s.TargetQuotaListURLs {
		if listURL == "" {
			continue
		}

		applications, err := fetchMiptListByURL(listURL, core.CompetitionTargetQuota)
		if err != nil {
			log.Printf("Error fetching Target Quota List %d (%s): %v", i+1, listURL, err)
			continue
		}

		if applications == nil {
			log.Printf("No applications found in Target Quota List %d (%s)", i+1, listURL)
			continue
		}

		// Set heading code and send applications
		for _, app := range applications {
			app.HeadingCode = headingCode
			receiver.PutApplicationData(app)
		}
	}

	return nil
}
