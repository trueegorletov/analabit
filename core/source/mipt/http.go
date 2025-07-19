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

	// Apply timeout coordination before HTTP request
	if err := source.WaitBeforeHTTPRequest("mipt", ctx); err != nil {
		return nil, fmt.Errorf("timeout coordination failed for %s: %w", listURL, err)
	}

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
	// Find the header row to detect the table format
	headerRow, tableNode := findTableHeaderRow(doc)
	if headerRow == nil {
		return nil, fmt.Errorf("table header not found in MIPT HTML")
	}
	format := detectTableFormat(headerRow)

	// Find all table rows
	rows := findTableRows(tableNode)
	if len(rows) <= 1 { // Should have at least a header and one data row
		return nil, fmt.Errorf("no data rows found in MIPT HTML")
	}

	var applications []*source.ApplicationData
	var errors []string

	// Process each row, skipping the header row (index 0)
	for i, row := range rows[1:] {
		app, err := parseApplicantFromTableRow(row, defaultCompetitionType, format)
		if err != nil {
			errors = append(errors, fmt.Sprintf("row %d: %v", i+2, err)) // i+2 because we skip the header
			continue
		}
		applications = append(applications, app)
	}

	if len(errors) > 0 {
		log.Printf("MIPT HTTP parser encountered %d errors:", len(errors))
		for _, err := range errors {
			log.Printf("  - %s", err)
		}
	}

	return applications, nil
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
