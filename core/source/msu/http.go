package msu

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/utils"
)

// HTTPHeadingSource implements source.HeadingSource for Moscow State University (MSU).
//
// One instance of HTTPHeadingSource corresponds to exactly one educational program
// (heading) at MSU.  Each program can have multiple competitive lists
// (общий конкурс, БВИ, особая квота, отдельная квота, целевая квота and
// детализированные целевые квоты).  The fields of this struct hold the
// information required to locate and parse these lists.
//
// A typical MSU rating page contains several anchors of the form
// `<h4 id="<anchor>">…</h4>` followed by a table with applicants.  The
// anchor names referenced here correspond to those fragment identifiers
// (without the leading '#').  For example, on
// https://cpk.msu.ru/rating/dep_01 the general list for МАТЕМАТИКА has
// anchor "01_01_1_01", its БВИ list has anchor "01_01_1_01_bvi", the
// special quota list has anchor "01_01_1_02", etc.  See the codegen
// utility for automatically populating these fields.
type HTTPHeadingSource struct {
	// PrettyName is the human‑readable name of the educational program (ОП).
	PrettyName string
	// FacultyURL is the absolute URL of the faculty page containing all
	// competitive lists for this program.  Example:
	// https://cpk.msu.ru/rating/dep_01
	FacultyURL string
	// RegularAnchor holds the anchor id for the general (regular) list.
	// If empty, no general list will be parsed.
	RegularAnchor string
	// BVIAnchor holds the anchor id for the list of applicants who have
	// the right to be admitted without entrance exams (БВИ).  If empty,
	// no BVI list will be parsed.
	BVIAnchor string
	// TargetQuotaAnchors holds anchor ids for all target quota lists,
	// including detalisised target quotas.  Applicants from all such
	// lists are parsed with CompetitionType set to CompetitionTargetQuota.
	TargetQuotaAnchors []string
	// DedicatedQuotaAnchors holds anchor ids for all dedicated quota lists.
	DedicatedQuotaAnchors []string
	// SpecialQuotaAnchors holds anchor ids for all special quota lists.
	SpecialQuotaAnchors []string
	// Capacities defines the number of places available on this program for
	// each competition type.  When LoadTo is called the value of this
	// field is copied to the produced HeadingData.
	Capacities core.Capacities
}

// LoadTo fetches the faculty page, parses all configured lists and sends
// HeadingData and ApplicationData to the provided receiver.  The caller is
// responsible for closing the receiver's channels after this method
// returns.  If any list cannot be downloaded or parsed the corresponding
// applications are silently skipped, allowing the remaining lists to
// continue.  A non‑nil error is returned only if the faculty page itself
// cannot be retrieved or parsed.
func (hs *HTTPHeadingSource) LoadTo(receiver source.DataReceiver) error {
	// Acquire semaphore slots for rate limiting
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Minute) // Increased timeout
	defer cancel()

	release, err := source.AcquireHTTPSemaphores(ctx, "msu")
	if err != nil {
		return fmt.Errorf("msu: failed to acquire semaphores: %w", err)
	}
	defer release()

	// prepare and send heading data first
	headingCode := utils.GenerateHeadingCode(hs.PrettyName)
	hd := &source.HeadingData{
		Code:       headingCode,
		Capacities: hs.Capacities,
		PrettyName: hs.PrettyName,
	}
	receiver.PutHeadingData(hd)

	// Create HTTP client with longer timeout for network resilience
	client := &http.Client{
		Timeout: 7 * time.Minute, // Individual request timeout
	}

	// fetch the faculty page with context and enhanced error logging
	req, err := http.NewRequestWithContext(ctx, "GET", hs.FacultyURL, nil)
	if err != nil {
		return fmt.Errorf("msu: failed to create request for %s: %w", hs.FacultyURL, err)
	}

	// Add diagnostic logging for network attempts
	slog.Info("MSU: Attempting to fetch faculty page", "url", hs.FacultyURL, "program", hs.PrettyName)
	resp, err := client.Do(req)
	if err != nil {
		// Enhanced error logging with network diagnostics
		slog.Error("MSU: Network error fetching faculty page", "url", hs.FacultyURL, "program", hs.PrettyName, "error", err)
		return fmt.Errorf("msu: failed to fetch faculty page %s: %w", hs.FacultyURL, err)
	}
	defer resp.Body.Close()
	slog.Info("MSU: Successfully fetched faculty page", "url", hs.FacultyURL, "program", hs.PrettyName, "status", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("msu: unexpected status %d while fetching %s", resp.StatusCode, hs.FacultyURL)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("msu: failed to parse faculty page %s: %w", hs.FacultyURL, err)
	}

	// build list descriptors
	type listDesc struct {
		anchor      string
		competition core.Competition
	}
	lists := make([]listDesc, 0)
	// order matters here: BVI is processed before regular so that the
	// competition types reflect the intended semantics of the lists.
	if hs.BVIAnchor != "" {
		lists = append(lists, listDesc{anchor: hs.BVIAnchor, competition: core.CompetitionBVI})
	}
	if hs.RegularAnchor != "" {
		lists = append(lists, listDesc{anchor: hs.RegularAnchor, competition: core.CompetitionRegular})
	}
	for _, a := range hs.TargetQuotaAnchors {
		if strings.TrimSpace(a) != "" {
			lists = append(lists, listDesc{anchor: a, competition: core.CompetitionTargetQuota})
		}
	}
	for _, a := range hs.DedicatedQuotaAnchors {
		if strings.TrimSpace(a) != "" {
			lists = append(lists, listDesc{anchor: a, competition: core.CompetitionDedicatedQuota})
		}
	}
	for _, a := range hs.SpecialQuotaAnchors {
		if strings.TrimSpace(a) != "" {
			lists = append(lists, listDesc{anchor: a, competition: core.CompetitionSpecialQuota})
		}
	}

	// iterate over each configured list and parse its rows
	for _, ld := range lists {
		// find the <h4 id="anchor"> element
		h4sel := doc.Find(fmt.Sprintf("h4#%s", ld.anchor))
		if h4sel.Length() == 0 {
			// Enhanced debugging: log missing anchors and show available anchors
			slog.Warn("MSU: Anchor not found on page", "url", hs.FacultyURL, "program", hs.PrettyName, "missingAnchor", ld.anchor, "competition", ld.competition.String())

			// Debug: show all available h4 anchors on the page
			availableAnchors := make([]string, 0)
			doc.Find("h4[id]").Each(func(i int, s *goquery.Selection) {
				if id, exists := s.Attr("id"); exists {
					availableAnchors = append(availableAnchors, id)
				}
			})
			slog.Info("MSU: Available anchors on page", "url", hs.FacultyURL, "anchors", availableAnchors)
			continue
		}
		// the table with the data is inside the div that follows the heading
		table := h4sel.Next().Find("table").First()
		if table.Length() == 0 {
			slog.Warn("MSU: No table found after anchor", "url", hs.FacultyURL, "program", hs.PrettyName, "anchor", ld.anchor, "competition", ld.competition.String())
			continue
		}
		slog.Info("MSU: Found table for anchor", "url", hs.FacultyURL, "program", hs.PrettyName, "anchor", ld.anchor, "competition", ld.competition.String())
		// parse all rows within the table body
		rowCount := 0
		table.Find("tbody tr").Each(func(i int, tr *goquery.Selection) {
			rowCount++
			cells := tr.Find("td")
			// we expect at least an ID cell; skip rows with too few columns
			if cells.Length() < 2 {
				slog.Debug("MSU: Skipping row with insufficient columns", "anchor", ld.anchor, "rowIndex", i, "cellCount", cells.Length())
				return
			}
			// rating place
			ratingPlace := 0
			if cells.Length() > 0 {
				rpStr := strings.TrimSpace(cells.Eq(0).Text())
				if rp, err := strconv.Atoi(rpStr); err == nil {
					ratingPlace = rp
				}
			}
			// student ID
			studentID := strings.TrimSpace(cells.Eq(1).Text())
			if studentID == "" {
				// no ID means a malformed row
				return
			}
			// original / consent submitted (3rd column contains "Да" or "Нет")
			originalSubmitted := false
			if cells.Length() > 2 {
				val := strings.TrimSpace(cells.Eq(2).Text())
				if val == "Да" {
					originalSubmitted = true
				}
			}

			// priority
			priority := 0
			if cells.Length() > 3 {
				pStr := strings.TrimSpace(cells.Eq(3).Text())
				if p, err := strconv.Atoi(pStr); err == nil {
					priority = p
				}
			}

			// scores sum
			scoresSum := 0
			if ld.competition != core.CompetitionBVI && cells.Length() > 7 {
				ssStr := strings.TrimSpace(cells.Eq(7).Text())
				if ss, err := strconv.Atoi(ssStr); err == nil {
					scoresSum = ss
				}
			}

			// dvi score
			dviScore := 0
			if ld.competition != core.CompetitionBVI && cells.Length() > 9 {
				dviStr := strings.TrimSpace(cells.Eq(9).Text())
				if dvi, err := strconv.Atoi(dviStr); err == nil {
					dviScore = dvi
				}
			}

			// ege scores
			var egeScores []int
			if ld.competition != core.CompetitionBVI && cells.Length() > 10 {
				for i := 10; i < cells.Length(); i++ {
					cellText := strings.TrimSpace(cells.Eq(i).Text())
					// Stop parsing scores when a non-numeric cell is encountered (e.g., "Нет")
					if score, err := strconv.Atoi(cellText); err == nil {
						egeScores = append(egeScores, score)
					} else {
						break
					}
				}
			}

			// create and send application data
			ad := &source.ApplicationData{
				HeadingCode:       headingCode,
				StudentID:         studentID,
				ScoresSum:         scoresSum,
				RatingPlace:       ratingPlace,
				Priority:          priority,
				CompetitionType:   ld.competition,
				OriginalSubmitted: originalSubmitted,
				DVIScore:          dviScore,
				EGEScores:         egeScores,
				HeadingName:       hs.PrettyName,
			}
			receiver.PutApplicationData(ad)
		})
		slog.Info("MSU: Processed rows for anchor", "anchor", ld.anchor, "competition", ld.competition.String(), "rowCount", rowCount)
	}

	return nil
}
