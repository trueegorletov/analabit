package core

import (
	"testing"
)

func TestNewUploadPayloadFromCalculator(t *testing.T) {
	// Create a simple VarsityCalculator with test data
	vc := NewVarsityCalculator("test_varsity", "Test Varsity")

	// Add a heading
	capacities := Capacities{Regular: 10, TargetQuota: 2, DedicatedQuota: 1, SpecialQuota: 1}
	vc.AddHeading("TEST01", capacities, "Test Heading")

	// Add students and applications
	vc.AddApplication("TEST01", "student1", 1, 1, CompetitionRegular, 85)
	vc.AddApplication("TEST01", "student2", 2, 1, CompetitionTargetQuota, 90)
	vc.SetOriginalSubmitted("student1")

	// Calculate admissions to get results
	results := vc.CalculateAdmissions()

	// Convert to UploadPayload
	payload := NewUploadPayloadFromCalculator(vc, results, nil)

	// Verify basic structure
	if payload.VarsityCode != "test_varsity" {
		t.Errorf("Expected varsity code 'test_varsity', got '%s'", payload.VarsityCode)
	}

	if payload.VarsityName != "Test Varsity" {
		t.Errorf("Expected varsity name 'Test Varsity', got '%s'", payload.VarsityName)
	}

	// Verify students
	if len(payload.Students) != 2 {
		t.Errorf("Expected 2 students, got %d", len(payload.Students))
	}

	// Find student1 and verify OriginalSubmitted flag
	var student1Found bool
	for _, student := range payload.Students {
		if student.ID == "student1" {
			student1Found = true
			if !student.OriginalSubmitted {
				t.Error("Expected student1 to have OriginalSubmitted=true")
			}
		}
	}
	if !student1Found {
		t.Error("student1 not found in payload")
	}

	// Verify applications
	if len(payload.Applications) != 2 {
		t.Errorf("Expected 2 applications, got %d", len(payload.Applications))
	}

	// Verify at least one application has correct data
	found := false
	for _, app := range payload.Applications {
		if app.StudentID == "student1" && app.HeadingCode == "test_varsity:TEST01" {
			found = true
			if app.RatingPlace != 1 {
				t.Errorf("Expected rating place 1, got %d", app.RatingPlace)
			}
			if app.CompetitionType != CompetitionRegular {
				t.Errorf("Expected competition type Regular, got %v", app.CompetitionType)
			}
		}
	}
	if !found {
		t.Error("Expected application for student1 not found")
	}

	// Verify calculations
	if len(payload.Calculations) != len(results) {
		t.Errorf("Expected %d calculations, got %d", len(results), len(payload.Calculations))
	}
}

func TestStudentDTO(t *testing.T) {
	dto := StudentDTO{
		ID:                "test123",
		OriginalSubmitted: true,
	}

	if dto.ID != "test123" {
		t.Errorf("Expected ID 'test123', got '%s'", dto.ID)
	}

	if !dto.OriginalSubmitted {
		t.Error("Expected OriginalSubmitted to be true")
	}
}

func TestApplicationDTO(t *testing.T) {
	dto := ApplicationDTO{
		StudentID:       "student1",
		HeadingCode:     "varsity:heading",
		Priority:        1,
		CompetitionType: CompetitionRegular,
		RatingPlace:     5,
		Score:           85,
	}

	if dto.StudentID != "student1" {
		t.Errorf("Expected StudentID 'student1', got '%s'", dto.StudentID)
	}

	if dto.HeadingCode != "varsity:heading" {
		t.Errorf("Expected HeadingCode 'varsity:heading', got '%s'", dto.HeadingCode)
	}

	if dto.CompetitionType != CompetitionRegular {
		t.Errorf("Expected CompetitionType Regular, got %v", dto.CompetitionType)
	}
}
