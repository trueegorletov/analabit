package oldrzgmu

import (
	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

// TestParseRZGMUTextData tests parsing of RZGMU text data
func TestParseRZGMUTextData(t *testing.T) {
	// Sample text data based on the actual RZGMU format
	textData := `Направление подготовки: Лечебное дело
№ Код Балл ВИ ИД ПП Приоритет Согласие
ЛЕЧЕБНОЕ ДЕЛО (БЮДЖЕТ) (Специалитет)
Мест: 19
Конкурсная группа: ОСНОВНЫЕ МЕСТА
1. 3867113 - БВИ 10 1 Согласие
2. 4235271 - БВИ 10 1 Согласие
3. 4237110 - БВИ 10 1 Согласие
4. 3652269 302
Биология - 95
Химия - 100
Русский язык - 97
10 Пр.право 2 Нет
5. 3630713 302
Биология - 95
Химия - 100
Русский язык - 97
10 3 Согласие`

	programs, err := parseRZGMUTextData(textData)
	if err != nil {
		t.Fatalf("Failed to parse RZGMU text data: %v", err)
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

	// Check applications
	if len(program.Applications) != 5 {
		t.Errorf("Expected 5 applications, got %d", len(program.Applications))
	}

	// Check first application (BVI)
	app1 := program.Applications[0]
	if app1.StudentID != "3867113" {
		t.Errorf("Expected student ID '3867113', got '%s'", app1.StudentID)
	}
	if app1.RatingPlace != 1 {
		t.Errorf("Expected rating place 1, got %d", app1.RatingPlace)
	}
	if app1.CompetitionType != core.CompetitionBVI {
		t.Errorf("Expected BVI competition type, got %v", app1.CompetitionType)
	}
	if !app1.OriginalSubmitted {
		t.Error("Expected original submitted to be true")
	}

	// Check fourth application (Regular with score)
	app4 := program.Applications[3]
	if app4.StudentID != "3652269" {
		t.Errorf("Expected student ID '3652269', got '%s'", app4.StudentID)
	}
	if app4.Priority != 2 {
		t.Errorf("Expected priority 2, got %d", app4.Priority)
	}
	if app4.CompetitionType != core.CompetitionRegular {
		t.Errorf("Expected Regular competition type, got %v", app4.CompetitionType)
	}
	if app4.OriginalSubmitted {
		t.Error("Expected original submitted to be false")
	}
}

// TestExtractTextFromPDFDebug tests text extraction from PDF and compares with sample
func TestExtractTextFromPDFDebug(t *testing.T) {
	// Path to the sample PDF file
	samplePDFPath := filepath.Join("..", "..", "..", "sample_data", "rzgmu", "rzgmu_l_b.pdf")

	// Extract text using rsc.io/pdf
	extractedText, err := extractTextFromPDFFile(samplePDFPath)
	if err != nil {
		t.Fatalf("Failed to extract text from PDF: %v", err)
	}

	// Read the sample text file for comparison
	sampleTextPath := filepath.Join("..", "..", "..", "sample_data", "rzgmu", "rzgmu_l_b-text-converted.txt")
	sampleTextBytes, err := ioutil.ReadFile(sampleTextPath)
	if err != nil {
		t.Fatalf("Failed to read sample text file: %v", err)
	}
	sampleText := string(sampleTextBytes)

	// Debug: Print both texts (first 500 chars) for comparison
	t.Logf("Extracted text (first 500 chars):\\n%s", truncateString(extractedText, 500))
	t.Logf("Sample text (first 500 chars):\\n%s", truncateString(sampleText, 500))

	// Check that both contain the key elements
	requiredElements := []string{
		"Направление подготовки:",
		"Лечебное дело",
		"ЛЕЧЕБНОЕ ДЕЛО (БЮДЖЕТ)",
		"Мест:",
		"3867113",
		"БВИ",
		"Биология",
		"Химия",
		"Русский язык",
	}

	for _, element := range requiredElements {
		if !strings.Contains(extractedText, element) {
			t.Errorf("Extracted text missing required element: '%s'", element)
		}
		if !strings.Contains(sampleText, element) {
			t.Errorf("Sample text missing required element: '%s'", element)
		}
	}

	// Test parsing both texts
	extractedPrograms, err := parseRZGMUTextData(extractedText)
	if err != nil {
		t.Errorf("Failed to parse extracted text: %v", err)
	} else {
		t.Logf("Successfully parsed %d programs from extracted text with %d applications",
			len(extractedPrograms), len(extractedPrograms[0].Applications))
	}

	samplePrograms, err := parseRZGMUTextData(sampleText)
	if err != nil {
		t.Errorf("Failed to parse sample text: %v", err)
	} else {
		t.Logf("Successfully parsed %d programs from sample text with %d applications",
			len(samplePrograms), len(samplePrograms[0].Applications))
	}

	// Compare application counts
	if len(extractedPrograms) > 0 && len(samplePrograms) > 0 {
		extractedCount := len(extractedPrograms[0].Applications)
		sampleCount := len(samplePrograms[0].Applications)

		if extractedCount != sampleCount {
			t.Logf("Warning: Different application counts - extracted: %d, sample: %d",
				extractedCount, sampleCount)
		}
	}
}

// TestPriorityParsingBug tests the specific bug where priority is incorrectly parsed
// for student 3950875 who should have priority 2 in Педиатрия but gets parsed as priority 1
func TestPriorityParsingBug(t *testing.T) {
	// Test case from rzgmu_p_b-text-converted.txt where student 3950875 should have priority 2
	// Line: "88.   3950875   275     Химия - 93          5         2               Нет"
	testData := `Направление подготовки: Педиатрия
№ Код Балл ВИ ИД ПП Приоритет Согласие
ПЕДИАТРИЯ (БЮДЖЕТ) (Специалитет)
Мест: 5
Конкурсная группа: ОСНОВНЫЕ МЕСТА
88. 3950875 275
Биология - 88
Химия - 93
5 2 Нет`

	programs, err := parseRZGMUTextData(testData)
	if err != nil {
		t.Fatalf("Failed to parse RZGMU text data: %v", err)
	}

	if len(programs) != 1 {
		t.Fatalf("Expected 1 program, got %d", len(programs))
	}

	program := programs[0]
	if len(program.Applications) != 1 {
		t.Fatalf("Expected 1 application, got %d", len(program.Applications))
	}

	app := program.Applications[0]
	if app.StudentID != "3950875" {
		t.Errorf("Expected student ID '3950875', got '%s'", app.StudentID)
	}

	// This is the bug: should be priority 2, but currently parses as priority 1
	if app.Priority != 2 {
		t.Errorf("BUG REPRODUCED: Expected priority 2, got %d", app.Priority)
	}

	if app.ScoresSum != 275 {
		t.Errorf("Expected score 275, got %d", app.ScoresSum)
	}

	if app.OriginalSubmitted {
		t.Error("Expected original submitted to be false (Нет)")
	}
}

// TestPriorityParsingVariousCases tests various priority parsing scenarios
func TestPriorityParsingVariousCases(t *testing.T) {
	testCases := []struct {
		name             string
		line             string
		expectedPriority int
		expectedConsent  bool
	}{
		{
			name:             "БВИ with priority 1 and consent",
			line:             "1. 3867113 - БВИ 10 1 Согласие",
			expectedPriority: 1,
			expectedConsent:  true,
		},
		{
			name:             "Regular with Пр.право and priority 2",
			line:             "4. 3652269 302 10 Пр.право 2 Нет",
			expectedPriority: 2,
			expectedConsent:  false,
		},
		{
			name:             "Regular with priority 3 and consent",
			line:             "5. 3630713 302 10 3 Согласие",
			expectedPriority: 3,
			expectedConsent:  true,
		},
		{
			name:             "Student 3950875 case",
			line:             "88. 3950875 275 5 2 Нет",
			expectedPriority: 2,
			expectedConsent:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app := parseApplicationStart(tc.line)
			if app == nil {
				t.Fatalf("Failed to parse application from line: %s", tc.line)
			}

			if app.Priority != tc.expectedPriority {
				t.Errorf("Expected priority %d, got %d for line: %s", tc.expectedPriority, app.Priority, tc.line)
			}

			if app.OriginalSubmitted != tc.expectedConsent {
				t.Errorf("Expected consent %v, got %v for line: %s", tc.expectedConsent, app.OriginalSubmitted, tc.line)
			}
		})
	}
}

// TestRealDataPriorityParsing tests priority parsing using actual sample data files
func TestRealDataPriorityParsing(t *testing.T) {
	// Test with the actual rzgmu_p_b-text-converted.txt file
	sampleTextPath := filepath.Join("..", "..", "..", "sample_data", "rzgmu", "rzgmu_p_b-text-converted.txt")
	sampleTextBytes, err := ioutil.ReadFile(sampleTextPath)
	if err != nil {
		t.Skipf("Skipping test: sample file not found: %v", err)
	}
	sampleText := string(sampleTextBytes)

	programs, err := parseRZGMUTextData(sampleText)
	if err != nil {
		t.Fatalf("Failed to parse RZGMU text data: %v", err)
	}

	if len(programs) != 1 {
		t.Fatalf("Expected 1 program, got %d", len(programs))
	}

	program := programs[0]

	// Find student 3950875 in the applications
	var student3950875 *source.ApplicationData
	for _, app := range program.Applications {
		if app.StudentID == "3950875" {
			student3950875 = app
			break
		}
	}

	if student3950875 == nil {
		t.Fatal("Student 3950875 not found in parsed applications")
	}

	// Verify correct priority parsing - should be 2 according to source data
	if student3950875.Priority != 2 {
		t.Errorf("Expected priority 2 for student 3950875 in Педиатрия, got %d", student3950875.Priority)
	}

	// Verify score
	if student3950875.ScoresSum != 275 {
		t.Errorf("Expected score 275 for student 3950875, got %d", student3950875.ScoresSum)
	}

	// Verify consent (should be false for "Нет")
	if student3950875.OriginalSubmitted {
		t.Error("Expected original submitted to be false for student 3950875")
	}

	t.Logf("SUCCESS: Student 3950875 correctly parsed with priority=%d, score=%d, consent=%v",
		student3950875.Priority, student3950875.ScoresSum, student3950875.OriginalSubmitted)
}

// TestActualSingleLineFormat tests parsing with the actual single-line format from the sample files
func TestActualSingleLineFormat(t *testing.T) {
	// Test case from rzgmu_p_b-text-converted.txt - should have priority 2
	testDataPediatrics := `Направление подготовки: Педиатрия
№ Код Балл ВИ ИД ПП Приоритет Согласие
ПЕДИАТРИЯ (БЮДЖЕТ) (Специалитет)
Мест: 5
Конкурсная группа: ОСНОВНЫЕ МЕСТА
88. 3950875 275     Химия - 93          5         2               Нет`

	programs, err := parseRZGMUTextData(testDataPediatrics)
	if err != nil {
		t.Fatalf("Failed to parse Pediatrics data: %v", err)
	}

	if len(programs) != 1 {
		t.Fatalf("Expected 1 program, got %d", len(programs))
	}

	if len(programs[0].Applications) != 1 {
		t.Fatalf("Expected 1 application, got %d", len(programs[0].Applications))
	}

	app := programs[0].Applications[0]
	if app.StudentID != "3950875" {
		t.Errorf("Expected student ID '3950875', got '%s'", app.StudentID)
	}

	if app.Priority != 2 {
		t.Errorf("PEDIATRICS BUG: Expected priority 2, got %d", app.Priority)
	}

	if app.ScoresSum != 275 {
		t.Errorf("Expected score 275, got %d", app.ScoresSum)
	}

	if app.OriginalSubmitted {
		t.Error("Expected original submitted to be false (Нет)")
	}

	// Test case from rzgmu_l_b-text-converted.txt - should have priority 1
	testDataMedicine := `Направление подготовки: Лечебное дело
№ Код Балл ВИ ИД ПП Приоритет Согласие
ЛЕЧЕБНОЕ ДЕЛО (БЮДЖЕТ) (Специалитет)
Мест: 19
Конкурсная группа: ОСНОВНЫЕ МЕСТА
93. 3950875 278     Химия - 93          8         1               Нет`

	programs2, err := parseRZGMUTextData(testDataMedicine)
	if err != nil {
		t.Fatalf("Failed to parse Medicine data: %v", err)
	}

	if len(programs2) != 1 {
		t.Fatalf("Expected 1 program, got %d", len(programs2))
	}

	if len(programs2[0].Applications) != 1 {
		t.Fatalf("Expected 1 application, got %d", len(programs2[0].Applications))
	}

	app2 := programs2[0].Applications[0]
	if app2.StudentID != "3950875" {
		t.Errorf("Expected student ID '3950875', got '%s'", app2.StudentID)
	}

	if app2.Priority != 1 {
		t.Errorf("MEDICINE BUG: Expected priority 1, got %d", app2.Priority)
	}

	if app2.ScoresSum != 278 {
		t.Errorf("Expected score 278, got %d", app2.ScoresSum)
	}

	if app2.OriginalSubmitted {
		t.Error("Expected original submitted to be false (Нет)")
	}
}

// Helper function to truncate string for display
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
