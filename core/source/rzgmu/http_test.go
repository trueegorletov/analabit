package rzgmu

import (
	"analabit/core"
	"analabit/core/source"
	"testing"
)

// TestParseRZGMUData tests parsing of RZGMU CSV data
func TestParseRZGMUData(t *testing.T) {
	// Sample CSV data based on the actual RZGMU format
	csvData := `Направление подготовки: Лечебное дело,,,,,,,
№,Код,Балл,ВИ,ИД,ПП,Приоритет,Согласие
ЛЕЧЕБНОЕ ДЕЛО (БЮДЖЕТ) (Специалитет) Мест: 19,,,,,,,
Конкурсная группа: ОСНОВНЫЕ МЕСТА,,,,,,,
1,3867113,-,БВИ,10,,1,Согласие
2,4235271,-,БВИ,10,,1,Согласие
3,4237110,-,БВИ,10,,1,Согласие
4,3652269,302,Биология - 95 Химия - 100 Русский язык - 97,10,Пр.право,2,Нет
5,3630713,302,Биология - 95 Химия - 100 Русский язык - 97,10,,3,Согласие`

	programs, err := parseRZGMUData(csvData)
	if err != nil {
		t.Fatalf("Failed to parse RZGMU data: %v", err)
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
		t.Fatalf("Expected 5 applications, got %d", len(program.Applications))
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
		t.Errorf("Expected original submitted to be true")
	}

	// Check fourth application (Regular with score)
	app4 := program.Applications[3]
	if app4.StudentID != "3652269" {
		t.Errorf("Expected student ID '3652269', got '%s'", app4.StudentID)
	}
	if app4.ScoresSum != 312 { // 302 + 10 bonus points
		t.Errorf("Expected total score 312 (302 + 10), got %d", app4.ScoresSum)
	}
	if app4.Priority != 2 {
		t.Errorf("Expected priority 2, got %d", app4.Priority)
	}
	if app4.CompetitionType != core.CompetitionRegular {
		t.Errorf("Expected Regular competition type, got %v", app4.CompetitionType)
	}
	if app4.OriginalSubmitted {
		t.Errorf("Expected original submitted to be false")
	}
}

// TestExtractProgramName tests program name extraction
func TestExtractProgramName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Направление подготовки: Лечебное дело", "Лечебное дело"},
		{"Направление подготовки: Педиатрия", "Педиатрия"},
		{"invalid line", ""},
		{"Направление подготовки:", ""},
	}

	for _, test := range tests {
		result := extractProgramName(test.input)
		if result != test.expected {
			t.Errorf("extractProgramName(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

// TestExtractCapacity tests capacity extraction
func TestExtractCapacity(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"ЛЕЧЕБНОЕ ДЕЛО (БЮДЖЕТ) (Специалитет) Мест: 19", 19},
		{"Some text Мест: 100", 100},
		{"Мест: 5", 5},
		{"no capacity info", 0},
		{"Мест:", 0},
	}

	for _, test := range tests {
		result := extractCapacity(test.input)
		if result != test.expected {
			t.Errorf("extractCapacity(%q) = %d, expected %d", test.input, result, test.expected)
		}
	}
}

// TestIsApplicationDataRow tests application data row detection
func TestIsApplicationDataRow(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"1,3867113,-,БВИ,10,,1,Согласие", true},
		{"4,3652269,302,exam_info,10,Пр.право,2,Нет", true},
		{"№,Код,Балл,ВИ,ИД,ПП,Приоритет,Согласие", false}, // header
		{"Конкурсная группа: ОСНОВНЫЕ МЕСТА", false},
		{"", false},
		{"1,2,3", false}, // too few columns
	}

	for _, test := range tests {
		result := isApplicationDataRow(test.input)
		if result != test.expected {
			t.Errorf("isApplicationDataRow(%q) = %v, expected %v", test.input, result, test.expected)
		}
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

// TestHTTPHeadingSourceLoadTo tests the main LoadTo method with mock data
func TestHTTPHeadingSourceLoadTo(t *testing.T) {
	// This test would require mocking HTTP and Tabula calls
	// For now, we'll test the parsing logic directly
	// A more complete test would require integration testing with actual PDF files
}
