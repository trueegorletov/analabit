package msu

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
)

type testReceiver struct {
	applications []*source.ApplicationData
	headings     []*source.HeadingData
}

func (r *testReceiver) PutApplicationData(application *source.ApplicationData) {
	r.applications = append(r.applications, application)
}

func (r *testReceiver) PutHeadingData(heading *source.HeadingData) {
	r.headings = append(r.headings, heading)
}

func TestHTTPHeadingSource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html, err := os.ReadFile("testdata/dep_01.html")
		if err != nil {
			t.Fatalf("failed to read testdata: %v", err)
		}
		w.Write(html)
	}))
	defer server.Close()

	// Choose a source from core/registry/msu/msu.go for testing
	var s source.HeadingSource = &HTTPHeadingSource{
		FacultyURL:            server.URL,
		PrettyName:            "Математика",
		RegularAnchor:         "01_01_1_01",
		BVIAnchor:             "01_01_1_01_bvi",
		TargetQuotaAnchors:    []string{"01_01_1_03"},
		DedicatedQuotaAnchors: []string{"01_01_1_04_1", "01_01_1_04_2"},
		SpecialQuotaAnchors:   []string{"01_01_1_02"},
		Capacities:            core.Capacities{Regular: 189, TargetQuota: 11, DedicatedQuota: 25, SpecialQuota: 25},
	}

	receiver := &testReceiver{}
	err := s.LoadTo(receiver)
	if err != nil {
		t.Fatalf("error loading data: %v", err)
	}

	printedCount := 0
	for _, app := range receiver.applications {
		if app.CompetitionType == core.CompetitionBVI {
			continue
		}

		if printedCount >= 5 {
			break
		}

		fmt.Printf("--- Application %d ---\n", printedCount+1)
		fmt.Printf("HeadingCode: %s\n", app.HeadingCode)
		fmt.Printf("StudentID: %s\n", app.StudentID)
		fmt.Printf("ScoresSum: %d\n", app.ScoresSum)
		fmt.Printf("RatingPlace: %d\n", app.RatingPlace)
		fmt.Printf("Priority: %d\n", app.Priority)
		fmt.Printf("CompetitionType: %v\n", app.CompetitionType)
		fmt.Printf("OriginalSubmitted: %v\n", app.OriginalSubmitted)
		fmt.Printf("DVI Score: %d\n", app.DVIScore)
		fmt.Printf("EGE Scores: %v\n", app.EGEScores)
		printedCount++
	}
}