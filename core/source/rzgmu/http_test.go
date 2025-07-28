package rzgmu

import (
	"strings"
	"testing"

	"github.com/trueegorletov/analabit/core"
)

// TestHTTPHeadingSource_Basic tests basic functionality of HTTPHeadingSource
func TestHTTPHeadingSource_Basic(t *testing.T) {
	// Test source creation
	source := &HTTPHeadingSource{
		ProgramName: "Лечебное дело",
		Capacities: core.Capacities{
			Regular: 100,
		},
	}

	if source.ProgramName != "Лечебное дело" {
		t.Errorf("Expected program name 'Лечебное дело', got '%s'", source.ProgramName)
	}

	if source.Capacities.Regular != 100 {
		t.Errorf("Expected regular capacity 100, got %d", source.Capacities.Regular)
	}
}

// TestHelperFunctions tests the helper functions
func TestHelperFunctions(t *testing.T) {
	// Test normalizeProgram
	normalized := normalizeProgram("  Лечебное дело  ")
	expected := "лечебное дело"
	if normalized != expected {
		t.Errorf("Expected '%s', got '%s'", expected, normalized)
	}

	// Test matchesProgram
	headingText := "Лечебное дело (бакалавриат)"
	targetProgram := "Лечебное дело"
	normalizedHeading := normalizeProgram(headingText)
	normalizedTarget := strings.ToLower(strings.TrimSpace(targetProgram))
	t.Logf("Heading: '%s' -> '%s'", headingText, normalizedHeading)
	t.Logf("Target: '%s' -> '%s'", targetProgram, normalizedTarget)
	if !matchesProgram(headingText, targetProgram) {
		t.Errorf("Expected matchesProgram to return true for '%s' vs '%s'", normalizedHeading, normalizedTarget)
	}

	// Test isBVI
	if !isBVI("<b>БВИ</b>") {
		t.Error("Expected isBVI to return true for БВИ")
	}

	if isBVI("<b>180</b>") {
		t.Error("Expected isBVI to return false for numeric score")
	}
}