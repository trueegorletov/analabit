package rzgmu

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/utils"
)

// HTTPHeadingSource loads RZGMU admission data from HTML pages
type HTTPHeadingSource struct {
	ProgramName string
	Capacities  core.Capacities // Required: pre-defined capacities for each competition type
}

// LoadTo downloads and parses all four RZGMU HTML pages, filtering for the specified program
func (h *HTTPHeadingSource) LoadTo(receiver source.DataReceiver) error {
	client := &http.Client{}
	
	// Generate consistent heading code based on program name
	headingCode := utils.GenerateHeadingCode(h.ProgramName)
	
	// Use pre-defined capacities directly from the struct
	totalCapacities := h.Capacities
	
	// Emit single heading data with combined capacities
	headingData := &source.HeadingData{
		Code:       headingCode,
		Capacities: totalCapacities,
		PrettyName: h.ProgramName,
	}
	receiver.PutHeadingData(headingData)
	
	// Parse application data from all pages
	for _, pageInfo := range pageURLs {
		resp, err := client.Get(pageInfo.URL)
		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", pageInfo.URL, err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("HTTP %d for %s", resp.StatusCode, pageInfo.URL)
		}
		
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to parse HTML from %s: %w", pageInfo.URL, err)
		}
		
		if err := h.parsePage(doc, pageInfo.Competition, headingCode, receiver); err != nil {
			return fmt.Errorf("failed to parse page %s: %w", pageInfo.URL, err)
		}
	}
	
	return nil
}



// parsePage extracts application data for the target program from a single HTML page
func (h *HTTPHeadingSource) parsePage(doc *goquery.Document, competition core.Competition, headingCode string, receiver source.DataReceiver) error {
	var foundTargetProgram bool
	
	// Look for program headings and tables
	doc.Find("h3, h4, h5, table").Each(func(i int, s *goquery.Selection) {
		tagName := goquery.NodeName(s)
		
		if strings.HasPrefix(tagName, "h") {
			// This is a heading - check if it's our target program
			headingText := strings.TrimSpace(s.Text())
			
			if matchesProgram(headingText, h.ProgramName) {
				foundTargetProgram = true
			} else {
				// Different program - reset state
				foundTargetProgram = false
			}
		} else if tagName == "table" && foundTargetProgram {
			// This is a table following our target program heading
			h.parseTable(s, headingCode, competition, receiver)
		}
	})
	
	return nil
}

// parseTable extracts application data from a table
func (h *HTTPHeadingSource) parseTable(table *goquery.Selection, headingCode string, competition core.Competition, receiver source.DataReceiver) {
	table.Find("tr").Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() < 7 { // Need at least 7 columns based on the sample data
			return
		}
		
		// Extract data from table cells
		// Based on sample: Position | StudentID | ScoresSum | Subjects | OriginalSubmitted | SpecialRight | Priority | ...
		positionText := strings.TrimSpace(cells.Eq(0).Text())
		studentIDText := strings.TrimSpace(cells.Eq(1).Text())
		scoreText := strings.TrimSpace(cells.Eq(2).Text())
		originalText := strings.TrimSpace(cells.Eq(4).Text())
		priorityText := strings.TrimSpace(cells.Eq(6).Text())
		
		// Skip header rows or empty rows
		if positionText == "" || studentIDText == "" {
			return
		}
		
		// Parse position
		position, err := strconv.Atoi(positionText)
		if err != nil {
			return // Skip non-numeric positions (likely header)
		}
		
		// Parse student ID
		studentID, err := strconv.Atoi(studentIDText)
		if err != nil {
			return // Skip invalid student IDs
		}
		
		// Parse scores sum (extract number from bold text like "<b>144</b>")
		scoreHTML, _ := cells.Eq(2).Html()
		var scoresSum int
		var actualCompetitionType = competition
		
		if isBVI(scoreHTML) {
			scoresSum = 0 // BVI candidates don't have numeric scores
			// Override competition type to BVI only if original competition is Regular
			if competition == core.CompetitionRegular {
				actualCompetitionType = core.CompetitionBVI
			}
		} else {
			// Extract numeric score from bold text
			if strings.Contains(scoreHTML, "<b>") {
				scoreText = strings.TrimSpace(cells.Eq(2).Text())
			}
			if scoresSum, err = strconv.Atoi(scoreText); err != nil {
				return // Skip rows with invalid scores
			}
		}
		
		// Parse original submitted (0 or positive integer)
		originalSubmitted := 0
		if originalText != "" {
			if val, err := strconv.Atoi(originalText); err == nil {
				originalSubmitted = val
			}
		}
		
		// Parse priority (0 or positive integer)
		priority := 0
		if priorityText != "" {
			if val, err := strconv.Atoi(priorityText); err == nil {
				priority = val
			}
		}
		
		// Create application data
		appData := &source.ApplicationData{
			HeadingCode:       headingCode,
			StudentID:         strconv.Itoa(studentID),
			ScoresSum:         scoresSum,
			RatingPlace:       position,
			Priority:          priority,
			CompetitionType:   actualCompetitionType,
			OriginalSubmitted: originalSubmitted > 0,
		}
		
		receiver.PutApplicationData(appData)
	})
}