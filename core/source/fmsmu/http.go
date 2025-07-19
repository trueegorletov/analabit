package fmsmu

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/utils"
	"golang.org/x/net/html"
)

// HTTPHeadingSource defines how to load FMSMU heading data from HTML URLs.
// FMSMU provides admission lists in HTML format with paginated individual list pages.
type HTTPHeadingSource struct {
	PrettyName           string          `json:"pretty_name"`
	RegularListID        string          `json:"regular_list_id"`
	SpecialQuotaListID   string          `json:"special_quota_list_id"`
	DedicatedQuotaListID string          `json:"dedicated_quota_list_id"`
	TargetQuotaListIDs   []string        `json:"target_quota_list_ids"`
	Capacities           core.Capacities `json:"capacities"`
}

// LoadTo loads data from HTTP source, downloading HTML pages and sending HeadingData and ApplicationData to the provided receiver.
func (s *HTTPHeadingSource) LoadTo(receiver source.DataReceiver) error {
	// Acquire a semaphore slot, respecting context cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	release, err := source.AcquireHTTPSemaphores(ctx, "fmsmu")
	if err != nil {
		return fmt.Errorf("failed to acquire semaphores for FMSMU: %w", err)
	}
	defer release()

	log.Printf("Processing FMSMU admission data for program: %s", s.PrettyName)

	headingCode := utils.GenerateHeadingCode(s.PrettyName)

	// Send HeadingData to the receiver
	receiver.PutHeadingData(&source.HeadingData{
		Code:       headingCode,
		Capacities: s.Capacities,
		PrettyName: s.PrettyName,
	})

	log.Printf("Sent FMSMU heading: %s (Code: %s, Caps: %v)", s.PrettyName, headingCode, s.Capacities)

	// Define list configurations
	listConfigs := []struct {
		ListID      string
		Competition core.Competition
		ListName    string
	}{
		{s.RegularListID, core.CompetitionRegular, "Regular List"},
		{s.SpecialQuotaListID, core.CompetitionSpecialQuota, "Special Quota List"},
		{s.DedicatedQuotaListID, core.CompetitionDedicatedQuota, "Dedicated Quota List"},
	}

	// Process individual lists
	totalApplicants := 0
	for _, config := range listConfigs {
		if config.ListID == "" {
			continue
		}

		applications, err := s.fetchListByID(ctx, config.ListID, config.Competition)
		if err != nil {
			log.Printf("Error fetching %s (ID: %s): %v", config.ListName, config.ListID, err)
			continue
		}

		if len(applications) == 0 {
			log.Printf("No applications found in %s (ID: %s)", config.ListName, config.ListID)
			continue
		}

		// Set heading code and send applications
		for _, app := range applications {
			app.HeadingCode = headingCode
			receiver.PutApplicationData(app)
			totalApplicants++
		}
	}

	// Handle multiple target quota list IDs
	for i, listID := range s.TargetQuotaListIDs {
		if listID == "" {
			continue
		}

		applications, err := s.fetchListByID(ctx, listID, core.CompetitionTargetQuota)
		if err != nil {
			log.Printf("Error fetching Target Quota List %d (ID: %s): %v", i+1, listID, err)
			continue
		}

		if len(applications) == 0 {
			log.Printf("No applications found in Target Quota List %d (ID: %s)", i+1, listID)
			continue
		}

		// Set heading code and send applications
		for _, app := range applications {
			app.HeadingCode = headingCode
			receiver.PutApplicationData(app)
			totalApplicants++
		}
	}

	log.Printf("Sent %d applications for FMSMU heading %s", totalApplicants, s.PrettyName)
	return nil
}

// fetchListByID fetches and parses all pages of a FMSMU list by ID
func (s *HTTPHeadingSource) fetchListByID(ctx context.Context, listID string, competitionType core.Competition) ([]*source.ApplicationData, error) {
	var allApplications []*source.ApplicationData
	page := 1
	maxPages := 300

	for {
		// Construct URL for this page
		url := fmt.Sprintf("https://priem.sechenov.ru/local/components/firstbit/competition.list/templates/.default/applications.php?COMPETITIVE_GROUP_ID=%s&appPage_%s=page-%d&ADMISSION_LISTS=N&CONTRACT_IS_PAID=N&lang=ru&search=", listID, listID, page)

		applications, doc, err := s.fetchListPage(ctx, url, competitionType)
		if err != nil {
			// If a single page fails, we might want to continue or stop.
			// For now, let's log and continue to the next page if we have a max page count.
			log.Printf("Error fetching page %d for list %s: %v. Skipping page.", page, listID, err)
			if page < maxPages {
				page++
				continue
			}
		}

		if page == 1 {
			maxPages = s.getMaxPageFromPagination(doc)
			if maxPages == 0 {
				log.Printf("Warning: Could not determine total page count for list %s on the first page. Falling back to a limit of 300 pages.", listID)
				maxPages = 300 // Fallback
			}
		}

		allApplications = append(allApplications, applications...)

		// Stop if there are no more applications on the page, as some pages might be empty before the end.
		if len(applications) == 0 && page > 1 {
			log.Printf("No applications found on page %d for list %s. Assuming end of list.", page, listID)
			break
		}

		page++
		if page > maxPages {
			break
		}
	}

	return allApplications, nil
}

// fetchListPage fetches and parses a single page of a FMSMU list
func (s *HTTPHeadingSource) fetchListPage(ctx context.Context, url string, competitionType core.Competition) ([]*source.ApplicationData, *html.Node, error) {
	var lastErr error
	const maxRetries = 3

	for i := 0; i <= maxRetries; i++ {
		if i > 0 {
			log.Printf("Retrying request for %s (attempt %d/%d) after error: %v", url, i, maxRetries, lastErr)
			time.Sleep(500 * time.Millisecond)
		}

		// Apply timeout coordination before making HTTP request
		if err := source.WaitBeforeHTTPRequest("fmsmu", ctx); err != nil {
			// This is likely a context cancellation or global rate limit issue, not worth retrying.
			return nil, nil, fmt.Errorf("timeout coordination failed: %w", err)
		}

		timeout := 30 * time.Second
		if i > 0 {
			timeout = 60 * time.Second
		}

		client := &http.Client{
			Timeout: timeout,
		}

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create request: %w", err) // non-recoverable
		}

		resp, err := client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("failed to download list page: %w", err)
			continue // retry
		}

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("failed to download list page (status code %d)", resp.StatusCode)
			resp.Body.Close()
			continue
		}

		htmlContent, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("failed to read HTML content: %w", err)
			continue // retry, could be a network issue during read
		}

		doc, err := html.Parse(strings.NewReader(string(htmlContent)))
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse HTML: %w", err)
		}

		// Parse applications from the table
		applications, err := s.parseApplicationsFromHTML(doc, competitionType)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse applications: %w", err)
		}

		return applications, doc, nil
	}

	return nil, nil, lastErr
}

// parseApplicationsFromHTML extracts applications from FMSMU HTML document
func (s *HTTPHeadingSource) parseApplicationsFromHTML(doc *html.Node, defaultCompetitionType core.Competition) ([]*source.ApplicationData, error) {
	// Find the table with class "table-competition-lists"
	table := findElementByClass(doc, "table-competition-lists")
	if table == nil {
		return nil, fmt.Errorf("table-competition-lists not found")
	}

	// Find all table rows with data-app attribute
	rows := s.findDataRows(table)
	if len(rows) == 0 {
		return []*source.ApplicationData{}, nil // Empty page
	}

	var applications []*source.ApplicationData
	var errors []string

	// Process each row
	for i, row := range rows {
		app, err := s.parseApplicationFromTableRow(row, defaultCompetitionType)
		if err != nil {
			errors = append(errors, fmt.Sprintf("row %d: %v", i+1, err))
			continue
		}

		if app != nil {
			applications = append(applications, app)
		}
	}

	if len(errors) > 0 {
		log.Printf("FMSMU HTTP parser encountered %d errors:", len(errors))
		for _, err := range errors {
			log.Printf("  - %s", err)
		}
	}

	return applications, nil
}

// findDataRows finds all table rows with data-app attribute
func (s *HTTPHeadingSource) findDataRows(table *html.Node) []*html.Node {
	var rows []*html.Node
	var find func(*html.Node)
	find = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			if getAttr(n, "data-app") != "" {
				rows = append(rows, n)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			find(c)
		}
	}
	find(table)
	return rows
}

// parseApplicationFromTableRow parses a single table row into ApplicationData
func (s *HTTPHeadingSource) parseApplicationFromTableRow(row *html.Node, defaultCompetitionType core.Competition) (*source.ApplicationData, error) {
	// Extract only outer-level td elements to avoid nested table cell offset
	cells := s.extractOuterTableCells(row)
	if len(cells) < 11 {
		return nil, fmt.Errorf("insufficient outer table cells: expected at least 11, got %d", len(cells))
	}

	// Check status column (last column) - skip if "Отозвано поступающим"
	statusText := strings.TrimSpace(getTextContent(cells[len(cells)-1]))
	if statusText == "Отозвано поступающим" {
		return nil, nil // Skip this row
	}

	// Extract data from the first cell (contains nested table with № and УИД)
	firstCell := cells[0]
	innerTable := s.findInnerTable(firstCell)
	if innerTable == nil {
		return nil, fmt.Errorf("inner table not found in first cell")
	}

	innerCells := s.extractTableCells(innerTable)
	if len(innerCells) < 2 {
		return nil, fmt.Errorf("insufficient inner table cells: expected 2, got %d", len(innerCells))
	}

	// Extract rating place (№)
	ratingPlaceText := strings.TrimSpace(getTextContent(innerCells[0]))
	ratingPlace, err := strconv.Atoi(ratingPlaceText)
	if err != nil {
		return nil, fmt.Errorf("failed to parse rating place '%s': %w", ratingPlaceText, err)
	}

	// Extract student ID (УИД)
	studentID := strings.TrimSpace(getTextContent(innerCells[1]))
	if studentID == "" {
		return nil, fmt.Errorf("empty student ID")
	}

	// Extract BVI basis (2nd column in outer cells)
	bviBasis := strings.TrimSpace(getTextContent(cells[1]))

	// Extract scores sum (3rd column in outer cells) with robust error handling
	scoresSumText := strings.TrimSpace(getTextContent(cells[2]))
	var scoresSum int
	if scoresSumText == "—" || scoresSumText == "-" || scoresSumText == "" {
		scoresSum = 285 // fallback value for non-numeric scores
		log.Printf("Using fallback score 285 for non-numeric value: '%s' (Student ID: %s)", scoresSumText, studentID)
	} else {
		var err error
		scoresSum, err = strconv.Atoi(scoresSumText)
		if err != nil {
			scoresSum = 285 // fallback for any parsing error
			log.Printf("Using fallback score 285 for parsing error: '%s' (Student ID: %s, Error: %v)", scoresSumText, studentID, err)
		}
	}

	// Validate score range and log warnings for unexpected values
	if scoresSum < 0 || scoresSum > 320 {
		log.Printf("Warning: Score %d outside expected range [0, 320] for Student ID: %s", scoresSum, studentID)
	}

	// Extract original submitted (11th column in outer cells)
	originalSubmittedText := strings.TrimSpace(getTextContent(cells[10]))
	originalSubmitted := originalSubmittedText == "Да"

	// Extract priority (10th column in outer cells)
	priorityText := strings.TrimSpace(getTextContent(cells[9]))
	priority := 1 // Default priority
	if priorityText != "" {
		if p, err := strconv.Atoi(priorityText); err == nil {
			priority = p
		}
	}

	// Determine competition type
	competitionType := defaultCompetitionType
	// Override to BVI only if bviBasis is not "—" AND it's a regular competition
	if bviBasis != "—" && defaultCompetitionType == core.CompetitionRegular {
		competitionType = core.CompetitionBVI
	}

	return &source.ApplicationData{
		StudentID:         studentID,
		ScoresSum:         scoresSum,
		RatingPlace:       ratingPlace,
		Priority:          priority,
		CompetitionType:   competitionType,
		OriginalSubmitted: originalSubmitted,
	}, nil
}

// extractTableCells extracts all td elements from a table row or table
func (s *HTTPHeadingSource) extractTableCells(node *html.Node) []*html.Node {
	var cells []*html.Node
	var find func(*html.Node)
	find = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "td" {
			cells = append(cells, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			find(c)
		}
	}
	find(node)
	return cells
}

// extractOuterTableCells extracts only direct td children of a table row, excluding nested table cells
func (s *HTTPHeadingSource) extractOuterTableCells(row *html.Node) []*html.Node {
	var cells []*html.Node
	for c := row.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "td" {
			cells = append(cells, c)
		}
	}
	return cells
}

// findInnerTable finds the inner table within a cell
func (s *HTTPHeadingSource) findInnerTable(cell *html.Node) *html.Node {
	var find func(*html.Node) *html.Node
	find = func(n *html.Node) *html.Node {
		if n.Type == html.ElementNode && n.Data == "table" {
			return n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if result := find(c); result != nil {
				return result
			}
		}
		return nil
	}
	return find(cell)
}

// getMaxPageFromPagination extracts the maximum page number from the pagination block.
// It returns 0 if the page number cannot be determined.
func (s *HTTPHeadingSource) getMaxPageFromPagination(doc *html.Node) int {
	pagination := findElementByClass(doc, "page-navigation")
	if pagination == nil {
		// If there's no pagination block, it's likely a single page.
		return 1
	}

	maxPage := 0
	var find func(*html.Node)
	find = func(n *html.Node) {
		// The page number can be in the text of a link or list item
		if n.Type == html.ElementNode && (n.Data == "a" || n.Data == "li") {
			textContent := strings.TrimSpace(getTextContent(n))
			// Check if the text content is a number
			if pageNum, err := strconv.Atoi(textContent); err == nil {
				if pageNum > maxPage {
					maxPage = pageNum
				}
			}
		}

		// Also check for page numbers in href attributes of <a> tags
		if n.Type == html.ElementNode && n.Data == "a" {
			href := getAttr(n, "href")
			re := regexp.MustCompile(`page-(\d+)`)
			matches := re.FindStringSubmatch(href)
			if len(matches) > 1 {
				if pageNum, err := strconv.Atoi(matches[1]); err == nil {
					if pageNum > maxPage {
						maxPage = pageNum
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			find(c)
		}
	}

	find(pagination)

	// If maxPage is still 0, but we found a pagination block, there might be just one page.
	if maxPage == 0 {
		// Let's check if there are any list items. If so, it's one page.
		var hasLi bool
		var checkLi func(*html.Node)
		checkLi = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "li" {
				hasLi = true
			}
			for c := n.FirstChild; c != nil && !hasLi; c = c.NextSibling {
				checkLi(c)
			}
		}
		checkLi(pagination)
		if hasLi {
			return 1
		}
	}

	return maxPage
}
