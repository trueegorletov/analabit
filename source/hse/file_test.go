package hse_test

import (
	"analabit/source"
	"analabit/source/hse"
	"reflect"
	"sync"
	"testing"
)

func TestFileHeadingSource_LoadTo_ConsistentData(t *testing.T) {
	s := hse.FileHeadingSource{
		RCListPath: "../../sample_data/hse/rc.xlsx",
		TQListPath: "../../sample_data/hse/tq.xlsx",
		DQListPath: "../../sample_data/hse/dq.xlsx",
		SQListPath: "../../sample_data/hse/sq.xlsx",
		BListPath:  "../../sample_data/hse/bvi.xlsx",
		Capacity:   10, // Arbitrary capacity, as it's part of the struct
	}

	collectData := func(src *hse.FileHeadingSource) ([]source.HeadingData, []source.ApplicationData, error) {
		headingsChan := make(chan source.HeadingData, 100) // Buffered for robustness
		applicationsChan := make(chan source.ApplicationData, 100) // Buffered for robustness

		var wg sync.WaitGroup
		var receivedHeadings []source.HeadingData
		var receivedApplications []source.ApplicationData

		wg.Add(1)
		go func() {
			defer wg.Done()
			for h := range headingsChan {
				receivedHeadings = append(receivedHeadings, h)
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			for a := range applicationsChan {
				receivedApplications = append(receivedApplications, a)
			}
		}()

		err := src.LoadTo(headingsChan, applicationsChan)
		wg.Wait() // Ensure all data is collected after LoadTo returns and channels are closed

		return receivedHeadings, receivedApplications, err
	}

	// First run
	headings1, applications1, err1 := collectData(&s)
	if err1 != nil {
		t.Fatalf("First call to LoadTo failed: %v", err1)
	}

	// Second run
	// Re-initialize FileHeadingSource if it had mutable state, but this one doesn't.
	// Or, ensure LoadTo is idempotent if called on the same instance.
	// The current FileHeadingSource is safe to reuse.
	headings2, applications2, err2 := collectData(&s)
	if err2 != nil {
		t.Fatalf("Second call to LoadTo failed: %v", err2)
	}

	// Compare results
	if !reflect.DeepEqual(headings1, headings2) {
		t.Errorf("Headings data is not consistent between calls.\nRun 1: %v\nRun 2: %v", headings1, headings2)
	}

	if !reflect.DeepEqual(applications1, applications2) {
		t.Errorf("Applications data is not consistent between calls.\nRun 1: %v\nRun 2: %v", applications1, applications2)
	}

	// For the current LoadTo implementation, we expect empty slices.
	// This part of the test will need adjustment if LoadTo starts producing data.
	if len(headings1) != 0 {
		t.Errorf("Expected no heading data with current LoadTo, got %d items", len(headings1))
	}
	if len(applications1) != 0 {
		t.Errorf("Expected no application data with current LoadTo, got %d items", len(applications1))
	}
}

