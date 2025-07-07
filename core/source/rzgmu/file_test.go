package rzgmu

import (
	"analabit/core"
	"context"
	"path/filepath"
	"strings"
	"testing"
)

// TestFileHeadingSourceLoadTo tests the FileHeadingSource with a real PDF file
func TestFileHeadingSourceLoadTo(t *testing.T) {
	// Path to the sample PDF file
	samplePDFPath := filepath.Join("..", "..", "..", "sample_data", "rzgmu", "rmgu_l_b.pdf")

	// Create a file source with specified program name
	fileSource := &FileHeadingSource{
		FilePath:    samplePDFPath,
		ProgramName: "Лечебное дело",
		Capacities:  nil, // Let it extract from the PDF
	}

	// Create a mock receiver to capture the data
	receiver := &MockDataReceiver{}

	// Load the data
	err := fileSource.LoadTo(receiver)
	if err != nil {
		t.Fatalf("Failed to load data from PDF file: %v", err)
	}

	// Verify we got at least one heading
	if len(receiver.Headings) == 0 {
		t.Fatal("Expected at least one heading, got none")
	}

	// Check the first heading
	heading := receiver.Headings[0]
	expectedName := "Лечебное дело"
	if heading.PrettyName != expectedName {
		t.Errorf("Expected heading name '%s', got '%s'", expectedName, heading.PrettyName)
	}

	// Check that heading code was generated
	if heading.Code == "" {
		t.Error("Expected non-empty heading code")
	}

	// Check capacities were extracted
	if heading.Capacities.Regular == 0 {
		t.Error("Expected non-zero regular capacity")
	}

	// Verify we got applications
	if len(receiver.Applications) == 0 {
		t.Fatal("Expected at least one application, got none")
	}

	// Check that applications have the correct heading code
	for i, app := range receiver.Applications {
		if app.HeadingCode != heading.Code {
			t.Errorf("Application %d has wrong heading code: expected '%s', got '%s'",
				i, heading.Code, app.HeadingCode)
		}
		if app.StudentID == "" {
			t.Errorf("Application %d has empty student ID", i)
		}
		if app.RatingPlace <= 0 {
			t.Errorf("Application %d has invalid rating place: %d", i, app.RatingPlace)
		}
	}

	t.Logf("Successfully processed %d headings and %d applications from PDF",
		len(receiver.Headings), len(receiver.Applications))
}

// TestFileHeadingSourceWithPresetCapacities tests FileHeadingSource with predefined capacities
func TestFileHeadingSourceWithPresetCapacities(t *testing.T) {
	// Path to the sample PDF file
	samplePDFPath := filepath.Join("..", "..", "..", "sample_data", "rzgmu", "rmgu_l_b.pdf")

	// Create custom capacities
	customCapacities := &core.Capacities{
		Regular:        25,
		TargetQuota:    5,
		DedicatedQuota: 2,
		SpecialQuota:   1,
	}

	// Create a file source with preset capacities
	fileSource := &FileHeadingSource{
		FilePath:    samplePDFPath,
		ProgramName: "Лечебное дело",
		Capacities:  customCapacities,
	}

	// Create a mock receiver to capture the data
	receiver := &MockDataReceiver{}

	// Load the data
	err := fileSource.LoadTo(receiver)
	if err != nil {
		t.Fatalf("Failed to load data from PDF file: %v", err)
	}

	// Verify we got at least one heading
	if len(receiver.Headings) == 0 {
		t.Fatal("Expected at least one heading, got none")
	}

	// Check that the custom capacities were used
	heading := receiver.Headings[0]
	if heading.Capacities.Regular != customCapacities.Regular {
		t.Errorf("Expected regular capacity %d, got %d",
			customCapacities.Regular, heading.Capacities.Regular)
	}
	if heading.Capacities.TargetQuota != customCapacities.TargetQuota {
		t.Errorf("Expected target quota %d, got %d",
			customCapacities.TargetQuota, heading.Capacities.TargetQuota)
	}
}

// TestFileHeadingSourceMissingFile tests error handling for missing files
func TestFileHeadingSourceMissingFile(t *testing.T) {
	fileSource := &FileHeadingSource{
		FilePath: "/nonexistent/path/to/file.pdf",
	}

	receiver := &MockDataReceiver{}
	err := fileSource.LoadTo(receiver)

	if err == nil {
		t.Fatal("Expected error for missing file, got nil")
	}

	// Should contain information about the missing file
	if !contains(err.Error(), "does not exist") {
		t.Errorf("Expected error message to mention missing file, got: %v", err)
	}
}

// TestFileHeadingSourceEmptyFilePath tests error handling for empty file path
func TestFileHeadingSourceEmptyFilePath(t *testing.T) {
	fileSource := &FileHeadingSource{
		FilePath: "",
	}

	receiver := &MockDataReceiver{}
	err := fileSource.LoadTo(receiver)

	if err == nil {
		t.Fatal("Expected error for empty file path, got nil")
	}

	expectedMsg := "FilePath is required"
	if !contains(err.Error(), expectedMsg) {
		t.Errorf("Expected error message to contain '%s', got: %v", expectedMsg, err)
	}
}

// TestParseRZGMUDataWithRealSample tests parsing with real converted CSV data
func TestParseRZGMUDataWithRealSample(t *testing.T) {
	// Use a subset of the real CSV data from the sample file
	csvData := `Направление подготовки: Лечебное дело,,,,,,,
№,Код,Балл,ВИ,ИД,ПП,Приоритет,Согласие
ЛЕЧЕБНОЕ ДЕЛО (БЮДЖЕТ) (Специалитет) Мест: 19,,,,,,,
Конкурсная группа: ОСНОВНЫЕ МЕСТА,,,,,,,
1,3867113,-,БВИ,10,,1,Согласие
2,4235271,-,БВИ,10,,1,Согласие
3,4237110,-,БВИ,10,,1,Согласие
4,3652269,302,Биология,10,Пр.право,2,Нет
5,3630713,302,Биология,10,,3,Согласие`

	programs, err := parseRZGMUData(csvData)
	if err != nil {
		t.Fatalf("Failed to parse real RZGMU data: %v", err)
	}

	if len(programs) != 1 {
		t.Fatalf("Expected 1 program, got %d", len(programs))
	}

	program := programs[0]

	// Check program name
	expectedName := "Лечебное дело"
	if program.PrettyName != expectedName {
		t.Errorf("Expected program name '%s', got '%s'", expectedName, program.PrettyName)
	}

	// Check extracted capacities
	if program.ExtractedCapacities.Regular != 19 {
		t.Errorf("Expected regular capacity 19, got %d", program.ExtractedCapacities.Regular)
	}

	// Check applications count
	if len(program.Applications) != 5 {
		t.Fatalf("Expected 5 applications, got %d", len(program.Applications))
	}

	// Check BVI applications (first 3)
	for i := 0; i < 3; i++ {
		app := program.Applications[i]
		if app.CompetitionType != core.CompetitionBVI {
			t.Errorf("Application %d should be BVI, got %v", i+1, app.CompetitionType)
		}
		if app.ScoresSum != 10 { // Only bonus points for BVI
			t.Errorf("Application %d expected score 10, got %d", i+1, app.ScoresSum)
		}
		if !app.OriginalSubmitted {
			t.Errorf("Application %d should have original submitted", i+1)
		}
	}

	// Check regular applications (last 2)
	app4 := program.Applications[3]
	if app4.CompetitionType != core.CompetitionRegular {
		t.Errorf("Application 4 should be Regular, got %v", app4.CompetitionType)
	}
	if app4.ScoresSum != 312 { // 302 + 10 bonus
		t.Errorf("Application 4 expected score 312, got %d", app4.ScoresSum)
	}
	if app4.Priority != 2 {
		t.Errorf("Application 4 expected priority 2, got %d", app4.Priority)
	}
	if app4.OriginalSubmitted {
		t.Errorf("Application 4 should not have original submitted")
	}

	app5 := program.Applications[4]
	if app5.CompetitionType != core.CompetitionRegular {
		t.Errorf("Application 5 should be Regular, got %v", app5.CompetitionType)
	}
	if app5.Priority != 3 {
		t.Errorf("Application 5 expected priority 3, got %d", app5.Priority)
	}
	if !app5.OriginalSubmitted {
		t.Errorf("Application 5 should have original submitted")
	}
}

// TestTabulaExtractionDebug helps debug what Tabula actually extracts
func TestTabulaExtractionDebug(t *testing.T) {
	// Path to the sample PDF file
	samplePDFPath := filepath.Join("..", "..", "..", "sample_data", "rzgmu", "rmgu_l_b.pdf")

	// Extract CSV data using Tabula directly
	csvData, err := extractCSVWithTabula(context.Background(), samplePDFPath)
	if err != nil {
		t.Fatalf("Failed to extract CSV with Tabula: %v", err)
	}

	t.Logf("Extracted CSV data (first 500 chars):\n%s", csvData[:min(500, len(csvData))])

	lines := strings.Split(csvData, "\n")
	t.Logf("Number of lines: %d", len(lines))
	for i, line := range lines[:min(10, len(lines))] {
		t.Logf("Line %d: %s", i+1, line)
	}
}

// TestParseRZGMUTabulaDataWithRealExtraction tests parsing with actual Tabula output
func TestParseRZGMUTabulaDataWithRealExtraction(t *testing.T) {
	// Sample data that represents what Tabula actually extracts
	csvData := `No,Код,Балл,ВИ,ИД,ПП,Приоритет,Согласие
ЛЕЧЕБНОЕ ДЕЛО (БЮДЖЕТ) (Специалитет) Мест: 19,,,,,,,
Конкурсная группа: ОСНОВНЫЕ МЕСТА,,,,,,,
1.,3867113,-,БВИ,10,,1,Согласие
2.,4235271,-,БВИ,10,,1,Согласие
3.,4237110,-,БВИ,10,,1,Согласие
4.,3652269,302,Биология,10,Пр.право,2,Нет
5.,3630713,302,Биология,10,,3,Согласие
6.,4023406,302,Биология,10,,1,Нет`

	programs, err := parseRZGMUTabulaData(csvData)
	if err != nil {
		t.Fatalf("Failed to parse Tabula RZGMU data: %v", err)
	}

	if len(programs) != 1 {
		t.Fatalf("Expected 1 program, got %d", len(programs))
	}

	program := programs[0]

	// Check program name (default)
	expectedName := "Лечебное дело"
	if program.PrettyName != expectedName {
		t.Errorf("Expected program name '%s', got '%s'", expectedName, program.PrettyName)
	}

	// Check that we got some applications
	if len(program.Applications) == 0 {
		t.Fatal("Expected some applications, got none")
	}

	t.Logf("Successfully parsed %d applications", len(program.Applications))

	// Check that BVI applications are correctly identified
	bviCount := 0
	regularCount := 0
	for _, app := range program.Applications {
		switch app.CompetitionType {
		case core.CompetitionBVI:
			bviCount++
		case core.CompetitionRegular:
			regularCount++
		}
	}

	if bviCount == 0 {
		t.Error("Expected to find some BVI applications")
	}
	if regularCount == 0 {
		t.Error("Expected to find some regular applications")
	}

	t.Logf("Found %d BVI applications and %d regular applications", bviCount, regularCount)
}

// TestCompetitionTypeDetection tests that competition types are correctly detected
func TestCompetitionTypeDetection(t *testing.T) {
	// Test BVI detection
	bviApp, err := parseApplicationRowTabula("1.,1234567,-,БВИ,10,,1,Согласие")
	if err != nil {
		t.Fatalf("Failed to parse BVI application: %v", err)
	}
	if bviApp.CompetitionType != core.CompetitionBVI {
		t.Errorf("Expected BVI competition type, got %v", bviApp.CompetitionType)
	}

	// Test regular competition detection
	regularApp, err := parseApplicationRowTabula("2.,2345678,300,Биология,5,,1,Нет")
	if err != nil {
		t.Fatalf("Failed to parse regular application: %v", err)
	}
	if regularApp.CompetitionType != core.CompetitionRegular {
		t.Errorf("Expected Regular competition type, got %v", regularApp.CompetitionType)
	}

	// Test target quota detection
	targetApp, err := parseApplicationRowTabula("3.,3456789,280,Целевое обучение,0,,1,Согласие")
	if err != nil {
		t.Fatalf("Failed to parse target quota application: %v", err)
	}
	if targetApp.CompetitionType != core.CompetitionTargetQuota {
		t.Errorf("Expected TargetQuota competition type, got %v", targetApp.CompetitionType)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && containsAt(s, substr, 0)) ||
		(len(s) > len(substr) && contains(s[1:], substr)))
}

func containsAt(s, substr string, pos int) bool {
	if pos+len(substr) > len(s) {
		return false
	}
	for i := 0; i < len(substr); i++ {
		if s[pos+i] != substr[i] {
			return false
		}
	}
	return true
}
