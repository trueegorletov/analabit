package rzgmu

import (
	"analabit/core"
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
	samplePDFPath := filepath.Join("..", "..", "..", "sample_data", "rzgmu", "rmgu_l_b.pdf")

	// Extract text using rsc.io/pdf
	extractedText, err := extractTextFromPDFFile(samplePDFPath)
	if err != nil {
		t.Fatalf("Failed to extract text from PDF: %v", err)
	}

	// Read the sample text file for comparison
	sampleTextPath := filepath.Join("..", "..", "..", "sample_data", "rzgmu", "rmgiu_l_b-text-converted.txt")
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

// TestPDFTextExtractionBasic tests basic PDF text extraction functionality
func TestPDFTextExtractionBasic(t *testing.T) {
	// Path to the sample PDF file
	samplePDFPath := filepath.Join("..", "..", "..", "sample_data", "rzgmu", "rmgu_l_b.pdf")

	// Extract text using rsc.io/pdf
	text, err := extractTextFromPDFFile(samplePDFPath)
	if err != nil {
		t.Fatalf("Failed to extract text from PDF: %v", err)
	}

	if len(text) == 0 {
		t.Fatal("Extracted text is empty")
	}

	// Should contain Russian text
	if !strings.Contains(text, "Лечебное дело") {
		t.Error("Extracted text should contain 'Лечебное дело'")
	}

	// Should contain student IDs
	if !strings.Contains(text, "3867113") {
		t.Error("Extracted text should contain student ID '3867113'")
	}

	t.Logf("Successfully extracted %d characters of text from PDF", len(text))
}

// Helper function to truncate string for display
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
