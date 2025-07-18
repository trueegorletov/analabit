package mipt

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trueegorletov/analabit/core"
	"golang.org/x/net/html"
)

func TestParseNewTableFormat(t *testing.T) {
	// Load the new format HTML file
	htmlContent, err := os.ReadFile("/home/yegor/Documents/Prest/analabit/sample_data/mipt/mipt-NEW-TABLE-FORMAT.html")
	assert.NoError(t, err)

	// Parse the HTML
	doc, err := html.Parse(strings.NewReader(string(htmlContent)))
	assert.NoError(t, err)

	// Parse the applications
	applications, err := parseApplicationsFromHTML(doc, core.CompetitionRegular)
	assert.NoError(t, err)

	// Validate the parsed data
	assert.Len(t, applications, 34, "Expected 34 applications to be parsed")

	// Print debug info for the first 5 applications
	for i := 0; i < 5 && i < len(applications); i++ {
		app := applications[i]
		t.Logf("App %d: ID=%s, Priority=%d, Original=%v, Competition=%s", i+1, app.StudentID, app.Priority, app.OriginalSubmitted, app.CompetitionType)
	}
}