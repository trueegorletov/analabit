package fmsmu

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/trueegorletov/analabit/core"
	"golang.org/x/net/html"
)

func TestDebugParsing(t *testing.T) {
	fmt.Println("=== FMSMU Debug Test - Using Production Parsing Logic ===")
	fmt.Println()

	// Create HTTPHeadingSource instance to use production methods
	source := &HTTPHeadingSource{
		PrettyName:           "Лечебное дело",
		RegularListID:        "4261",
		SpecialQuotaListID:   "5291",
		DedicatedQuotaListID: "5292",
		Capacities: core.Capacities{
			Regular:        274,
			SpecialQuota:   110,
			DedicatedQuota: 110,
		},
	}

	// Test 1: Parse sample HTML file using production logic
	fmt.Println("=== TEST 1: Sample HTML File Parsing ===")
	testSampleHTML(t, source)

	// Test 2: Optionally test with real HTTP data (commented out by default to avoid network calls in tests)
	fmt.Println("\n=== TEST 2: Real HTTP Data ===")
	testRealHTTP(t, source)

	fmt.Println("\nDebug test completed successfully.")
}

func testSampleHTML(t *testing.T, source *HTTPHeadingSource) {
	// Load the sample HTML file
	htmlContent, err := os.ReadFile("/home/yegor/Documents/Prest/analabit/sample_data/fmsmu/sample_table_element.html")
	if err != nil {
		t.Fatalf("Failed to load sample HTML file: %v", err)
	}

	// Parse HTML using golang.org/x/net/html (same as production)
	doc, err := html.Parse(strings.NewReader(string(htmlContent)))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	// Use production parsing method
	applications, err := source.parseApplicationsFromHTML(doc, core.CompetitionRegular)
	if err != nil {
		t.Fatalf("Failed to parse applications: %v", err)
	}

	fmt.Printf("Successfully parsed %d applications from sample HTML\n", len(applications))

	// Analyze each application with detailed debugging
	for i, app := range applications {
		fmt.Printf("\n--- Application %d ---\n", i+1)
		fmt.Printf("Student ID: %s\n", app.StudentID)
		fmt.Printf("Rating Place: %d\n", app.RatingPlace)
		fmt.Printf("Scores Sum: %d\n", app.ScoresSum)
		fmt.Printf("Priority: %d\n", app.Priority)
		fmt.Printf("Competition Type: %s\n", app.CompetitionType)
		fmt.Printf("Original Submitted: %t\n", app.OriginalSubmitted)

		// Validate score range (should be reasonable for FMSMU applications)
		if app.ScoresSum < 200 || app.ScoresSum > 320 {
			fmt.Printf("⚠️  Warning: Score %d outside typical range [200, 320]\n", app.ScoresSum)
		}

		// Validate that scores are not equal to student IDs (the bug we're fixing)
		if app.StudentID == fmt.Sprintf("%d", app.ScoresSum) {
			t.Errorf("BUG DETECTED: Score %d equals Student ID %s - this indicates cell mapping bug!", app.ScoresSum, app.StudentID)
		}

		// Validate student ID format (should be numeric)
		if app.StudentID == "" {
			t.Errorf("Empty student ID for application %d", i+1)
		}

		// Validate rating place (should be positive)
		if app.RatingPlace <= 0 {
			t.Errorf("Invalid rating place %d for application %d", app.RatingPlace, i+1)
		}
	}

	// Test the difference between flawed and correct cell extraction
	fmt.Println("\n=== CELL EXTRACTION COMPARISON ===")
	testCellExtractionDifference(t, source, doc)

	// Assertions
	if len(applications) == 0 {
		t.Error("No applications were parsed from sample data")
	}

	// Validate that all scores are in reasonable range
	for i, app := range applications {
		if app.ScoresSum < 0 || app.ScoresSum > 400 {
			t.Errorf("Application %d: Score %d outside valid range [0, 400]", i+1, app.ScoresSum)
		}
	}

	fmt.Printf("\n✅ Sample HTML test completed successfully with %d applications\n", len(applications))
}

func testCellExtractionDifference(t *testing.T, source *HTTPHeadingSource, doc *html.Node) {
	// Find the table
	table := findElementByClass(doc, "table-competition-lists")
	if table == nil {
		t.Fatal("table-competition-lists not found")
	}

	// Find data rows
	rows := source.findDataRows(table)
	if len(rows) == 0 {
		t.Fatal("No data rows found")
	}

	fmt.Printf("Found %d data rows for cell extraction comparison\n", len(rows))

	// Analyze first row in detail
	if len(rows) > 0 {
		row := rows[0]
		fmt.Println("\n--- Cell Extraction Analysis (First Row) ---")

		// Method 1: extractTableCells (flawed - includes nested cells)
		allCells := source.extractTableCells(row)
		fmt.Printf("extractTableCells (flawed): %d cells\n", len(allCells))
		for i, cell := range allCells {
			if i < 10 { // Show first 10 cells
				cellText := strings.TrimSpace(getTextContent(cell))
				fmt.Printf("  All[%d]: '%s'\n", i, cellText)
			}
		}

		// Method 2: extractOuterTableCells (correct - only direct td children)
		outerCells := source.extractOuterTableCells(row)
		fmt.Printf("\nextractOuterTableCells (correct): %d cells\n", len(outerCells))
		for i, cell := range outerCells {
			if i < 10 { // Show first 10 cells
				cellText := strings.TrimSpace(getTextContent(cell))
				fmt.Printf("  Outer[%d]: '%s'\n", i, cellText)
			}
		}

		// Show the difference
		fmt.Printf("\nDifference: %d extra cells from nested tables\n", len(allCells)-len(outerCells))

		// Demonstrate how this affects parsing
		if len(outerCells) >= 3 {
			// Extract from first cell (contains nested table)
			firstCell := outerCells[0]
			innerTable := source.findInnerTable(firstCell)
			if innerTable != nil {
				innerCells := source.extractTableCells(innerTable)
				if len(innerCells) >= 2 {
					ratingPlace := strings.TrimSpace(getTextContent(innerCells[0]))
					studentID := strings.TrimSpace(getTextContent(innerCells[1]))
					fmt.Printf("\nFrom nested table: Rating=%s, StudentID=%s\n", ratingPlace, studentID)
				}
			}

			// BVI basis (2nd outer cell)
			bviBasis := strings.TrimSpace(getTextContent(outerCells[1]))
			fmt.Printf("BVI Basis (outer[1]): '%s'\n", bviBasis)

			// Scores sum (3rd outer cell)
			scoresText := strings.TrimSpace(getTextContent(outerCells[2]))
			fmt.Printf("Scores Sum (outer[2]): '%s'\n", scoresText)
		}
	}
}

func testRealHTTP(t *testing.T, source *HTTPHeadingSource) {
	fmt.Println("Testing with real HTTP data from FMSMU...")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test URL for page 25 of "Лечебное дело" Regular list (ID: 4261)
	url := "https://priem.sechenov.ru/local/components/firstbit/competition.list/templates/.default/applications.php?COMPETITIVE_GROUP_ID=4261&appPage_4261=page-15&ADMISSION_LISTS=N&CONTRACT_IS_PAID=N&lang=ru&search="

	fmt.Printf("Fetching: %s\n", url)

	applications, doc, err := source.fetchListPage(ctx, url, core.CompetitionRegular)
	if err != nil {
		// Don't fail the test for network issues, just log
		fmt.Printf("⚠️  Could not fetch real HTTP data: %v\n", err)
		fmt.Println("This is expected in CI/CD environments without internet access.")
		return
	}

	fmt.Printf("Successfully fetched and parsed %d applications from real HTTP data\n", len(applications))

	if len(applications) > 0 {
		fmt.Println("\n--- Real HTTP Data Analysis ---")

		// Analyze score distribution
		var scoresInRange, scoresOutOfRange int
		minScore, maxScore := 400, 0

		for i, app := range applications {
			if i < 5 { // Show first 5 applications
				fmt.Printf("App %d: ID=%s, Place=%d, Score=%d, Priority=%d\n",
					i+1, app.StudentID, app.RatingPlace, app.ScoresSum, app.Priority)
			}

			// Track score statistics
			if app.ScoresSum >= 280 && app.ScoresSum <= 310 {
				scoresInRange++
			} else {
				scoresOutOfRange++
			}

			if app.ScoresSum > 0 {
				if app.ScoresSum < minScore {
					minScore = app.ScoresSum
				}
				if app.ScoresSum > maxScore {
					maxScore = app.ScoresSum
				}
			}
		}

		fmt.Printf("\nScore Analysis:\n")
		fmt.Printf("  Total applications: %d\n", len(applications))
		fmt.Printf("  Scores in range [280-310]: %d\n", scoresInRange)
		fmt.Printf("  Scores outside range: %d\n", scoresOutOfRange)
		fmt.Printf("  Score range: %d - %d\n", minScore, maxScore)

		// Validate that most scores are in expected range for competitive programs
		expectedRatio := 0.7 // At least 70% should be in reasonable range
		actualRatio := float64(scoresInRange) / float64(len(applications))
		if actualRatio < expectedRatio {
			t.Errorf("Too many scores outside expected range: %.2f%% (expected at least %.2f%%)",
				actualRatio*100, expectedRatio*100)
		}

		// Test pagination parsing
		maxPages := source.getMaxPageFromPagination(doc)
		fmt.Printf("Detected max pages: %d\n", maxPages)
		if maxPages == 0 {
			fmt.Println("⚠️  Could not detect pagination - might be single page or parsing issue")
		}
	}

	fmt.Println("✅ Real HTTP test completed successfully")
}
