package mipt

import (
	"testing"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
)

// TestReceiver implements source.DataReceiver for testing
type TestReceiver struct {
	HeadingData     []*source.HeadingData
	ApplicationData []*source.ApplicationData
}

func (tr *TestReceiver) PutHeadingData(heading *source.HeadingData) {
	tr.HeadingData = append(tr.HeadingData, heading)
}

func (tr *TestReceiver) PutApplicationData(application *source.ApplicationData) {
	tr.ApplicationData = append(tr.ApplicationData, application)
}

func TestMIPTFileParser(t *testing.T) {
	// Test with the sample MIPT file
	filePath := "../../../sample_data/mipt/mipt-RegularBVI-sample-list-UNCOMPRESSED.html"

	source := &FileHeadingSource{
		FilePath: filePath,
		Capacities: core.Capacities{
			Regular:        30,
			TargetQuota:    3,
			SpecialQuota:   3,
			DedicatedQuota: 3,
		},
	}

	receiver := &TestReceiver{}

	err := source.LoadTo(receiver)
	if err != nil {
		t.Fatalf("Failed to parse MIPT file: %v", err)
	}

	// Check that we got heading data
	if len(receiver.HeadingData) != 1 {
		t.Errorf("Expected 1 heading, got %d", len(receiver.HeadingData))
	}

	heading := receiver.HeadingData[0]
	if heading.Code == "" {
		t.Error("Heading code should not be empty")
	}
	if heading.PrettyName == "" {
		t.Error("Heading pretty name should not be empty")
	}

	// Check that we got application data
	if len(receiver.ApplicationData) == 0 {
		t.Error("Expected at least one application, got 0")
	}

	t.Logf("Parsed heading: %s (code: %s)", heading.PrettyName, heading.Code)
	t.Logf("Parsed %d applications", len(receiver.ApplicationData))

	// Validate some application data
	for i, app := range receiver.ApplicationData {
		if i >= 250 { // Only check first 5 for brevity
			break
		}

		if app.StudentID == "" {
			t.Errorf("Application %d: Student ID should not be empty", i+1)
		}
		if app.RatingPlace <= 0 {
			t.Errorf("Application %d: Rating place should be positive, got %d", i+1, app.RatingPlace)
		}
		if app.HeadingCode != heading.Code {
			t.Errorf("Application %d: Heading code mismatch", i+1)
		}

		t.Logf("App %d: ID=%s, Place=%d, Score=%d, Priority=%d, Original=%t, Competition=%s",
			i+1, app.StudentID, app.RatingPlace, app.ScoresSum, app.Priority, app.OriginalSubmitted, app.CompetitionType)
	}
}
