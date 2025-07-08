package hse

import (
	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"testing"

	"github.com/xuri/excelize/v2"
)

// TestParseSampleFile tests parsing the sample HSE XLSX file
func TestParseSampleFile(t *testing.T) {
	// Open the sample XLSX file
	f, err := excelize.OpenFile("../../../sample_data/hse/BD_moscow_AMI_O.xlsx")
	if err != nil {
		t.Fatalf("Failed to open sample XLSX file: %v", err)
	}
	defer f.Close()

	// Test extracting pretty name
	prettyName, err := extractPrettyNameFromXLSX(f, "test_file")
	if err != nil {
		t.Fatalf("Failed to extract pretty name: %v", err)
	}

	if prettyName == "" {
		t.Fatal("Pretty name is empty")
	}

	t.Logf("Extracted pretty name: %s", prettyName)

	// Test parsing applications
	mockReceiver := &MockDataReceiver{}
	err = parseApplicationsFromXLSX(f, "TEST_CODE", mockReceiver, "test_file")
	if err != nil {
		t.Fatalf("Failed to parse applications: %v", err)
	}

	if len(mockReceiver.Applications) == 0 {
		t.Fatal("No applications were parsed")
	}

	t.Logf("Parsed %d applications", len(mockReceiver.Applications))

	// Test some specific applications
	app := mockReceiver.Applications[0]
	if app.HeadingCode != "TEST_CODE" {
		t.Errorf("Expected heading code 'TEST_CODE', got '%s'", app.HeadingCode)
	}

	if app.StudentID == "" {
		t.Error("Student ID is empty")
	}

	if app.RatingPlace <= 0 {
		t.Errorf("Rating place should be positive, got %d", app.RatingPlace)
	}

	t.Logf("First application: StudentID=%s, RatingPlace=%d, CompetitionType=%v, ScoresSum=%d",
		app.StudentID, app.RatingPlace, app.CompetitionType, app.ScoresSum)

	// Check that we have different competition types
	competitionTypes := make(map[core.Competition]int)
	for _, app := range mockReceiver.Applications {
		competitionTypes[app.CompetitionType]++
	}

	t.Logf("Competition type distribution:")
	for compType, count := range competitionTypes {
		t.Logf("  %v: %d", compType, count)
	}

	if len(competitionTypes) < 2 {
		t.Error("Expected multiple competition types in the sample data")
	}
}

// MockDataReceiver implements source.DataReceiver for testing
type MockDataReceiver struct {
	Headings     []*source.HeadingData
	Applications []*source.ApplicationData
}

func (m *MockDataReceiver) PutHeadingData(heading *source.HeadingData) {
	m.Headings = append(m.Headings, heading)
}

func (m *MockDataReceiver) PutApplicationData(application *source.ApplicationData) {
	m.Applications = append(m.Applications, application)
}
