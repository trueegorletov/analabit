package itmo

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"golang.org/x/net/html"
)

// mockDataReceiver is a test implementation of DataReceiver
type mockDataReceiver struct {
	headings     []*source.HeadingData
	applications []*source.ApplicationData
}

func (m *mockDataReceiver) PutHeadingData(heading *source.HeadingData) {
	m.headings = append(m.headings, heading)
}

func (m *mockDataReceiver) PutApplicationData(application *source.ApplicationData) {
	m.applications = append(m.applications, application)
}

func TestParseApplicant(t *testing.T) {
	// Test HTML content for a sample application
	htmlContent := `
	<div class="RatingPage_table__item__qMY0F">
		<div class="RatingPage_table__name__rhwwJ">
			<div>
				<p class="RatingPage_table__position__uYWvi">163 № 4274517</p>
			</div>
		</div>
		<div class="RatingPage_table__block__5Sf4O">
			<div class="RatingPage_table__info__quwhV">
				<div class="RatingPage_table__infoLeft__Y_9cA">
					<p>Приоритет: 1</p>
					<p>Математика: 100</p>
					<p>Информатика: 100</p>
					<p>Русский язык: 100</p>
				</div>
			</div>
			<div class="RatingPage_table__info__quwhV">
				<div class="RatingPage_table__infoLeft__Y_9cA">
					<p>ИД: 10</p>
					<p>Балл ВИ+ИД: 310</p>
				</div>
				<div>
					<p>Есть согласие: нет</p>
				</div>
			</div>
		</div>
	</div>`

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		t.Fatalf("Failed to parse test HTML: %v", err)
	}

	// Find the application node
	var appNode *html.Node
	var findAppNode func(*html.Node)
	findAppNode = func(n *html.Node) {
		if appNode != nil {
			return
		}
		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				if attr.Key == "class" && strings.Contains(attr.Val, "table__item__qMY0F") {
					appNode = n
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findAppNode(c)
		}
	}
	findAppNode(doc)

	if appNode == nil {
		t.Fatal("Could not find application node in test HTML")
	}

	// Parse the application
	app, err := parseApplicant(appNode, core.CompetitionRegular)
	if err != nil {
		t.Fatalf("Failed to parse applicant: %v", err)
	}

	// Verify the parsed data
	if app.StudentID != "4274517" {
		t.Errorf("Expected StudentID '4274517', got '%s'", app.StudentID)
	}
	if app.RatingPlace != 163 {
		t.Errorf("Expected RatingPlace 163, got %d", app.RatingPlace)
	}
	if app.Priority != 1 {
		t.Errorf("Expected Priority 1, got %d", app.Priority)
	}
	if app.CompetitionType != core.CompetitionRegular {
		t.Errorf("Expected CompetitionType Regular, got %v", app.CompetitionType)
	}
	if app.ScoresSum != 310 {
		t.Errorf("Expected ScoresSum 310, got %d", app.ScoresSum)
	}
	if app.OriginalSubmitted != false {
		t.Errorf("Expected OriginalSubmitted false, got %t", app.OriginalSubmitted)
	}
}

func TestExtractCapacities(t *testing.T) {
	htmlContent := `
	<div>
		<p>Количество мест: 170 (17 ЦК, 17 ОcК, 17 ОтК)</p>
	</div>`

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		t.Fatalf("Failed to parse test HTML: %v", err)
	}

	capacities := extractCapacities(doc)

	expectedRegular := 170 - 17 - 17 - 17 // 119
	if capacities.Regular != expectedRegular {
		t.Errorf("Expected Regular capacity %d, got %d", expectedRegular, capacities.Regular)
	}
	if capacities.TargetQuota != 17 {
		t.Errorf("Expected TargetQuota 17, got %d", capacities.TargetQuota)
	}
	if capacities.SpecialQuota != 17 {
		t.Errorf("Expected SpecialQuota 17, got %d", capacities.SpecialQuota)
	}
	if capacities.DedicatedQuota != 17 {
		t.Errorf("Expected DedicatedQuota 17, got %d", capacities.DedicatedQuota)
	}
}

func TestExtractPrettyName(t *testing.T) {
	htmlContent := `
	<html>
		<body>
			<h2>01.03.02 «Прикладная математика и информатика»</h2>
		</body>
	</html>`

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		t.Fatalf("Failed to parse test HTML: %v", err)
	}

	name := extractPrettyName(doc)
	expected := "Прикладная математика и информатика"
	if name != expected {
		t.Errorf("Expected pretty name '%s', got '%s'", expected, name)
	}
}

func TestHTTPHeadingSourceWithSampleFile(t *testing.T) {
	// Try to find the sample file
	samplePath := filepath.Join("..", "..", "..", "sample_data", "itmo", "itmo_applications_sample.html")

	// Check if the sample file exists
	if _, err := os.Stat(samplePath); os.IsNotExist(err) {
		t.Skip("Sample file not found, skipping integration test")
	}

	// Read the sample file
	content, err := os.ReadFile(samplePath)
	if err != nil {
		t.Fatalf("Failed to read sample file: %v", err)
	}

	// Parse the HTML
	doc, err := html.Parse(strings.NewReader(string(content)))
	if err != nil {
		t.Fatalf("Failed to parse sample HTML: %v", err)
	}

	// Test capacity extraction
	capacities := extractCapacities(doc)
	if capacities.Regular == 0 && capacities.TargetQuota == 0 &&
		capacities.SpecialQuota == 0 && capacities.DedicatedQuota == 0 {
		t.Error("Failed to extract any capacities from sample file")
	}

	// Test pretty name extraction
	prettyName := extractPrettyName(doc)
	if prettyName == "" {
		t.Error("Failed to extract pretty name from sample file")
	}

	// Test application parsing
	receiver := &mockDataReceiver{}
	source := &HTTPHeadingSource{
		PrettyName: prettyName,
		Capacities: capacities,
	}

	err = source.parseApplicationsByCategory(doc, receiver, "test-code")
	if err != nil {
		t.Fatalf("Failed to parse applications: %v", err)
	}

	if len(receiver.applications) == 0 {
		t.Error("No applications were parsed from sample file")
	}

	// Verify we have applications in different categories
	categoryCount := make(map[core.Competition]int)
	for _, app := range receiver.applications {
		categoryCount[app.CompetitionType]++
	}

	if len(categoryCount) == 0 {
		t.Error("No competition categories found")
	}

	t.Logf("Parsed %d applications across %d categories", len(receiver.applications), len(categoryCount))
	for comp, count := range categoryCount {
		t.Logf("  %s: %d applications", comp.String(), count)
	}

	// Test that we have reasonable data
	if len(receiver.applications) < 100 {
		t.Errorf("Expected at least 100 applications, got %d", len(receiver.applications))
	}
}
