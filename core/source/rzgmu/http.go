// Package rzgmu provides support for loading RZGMU (Рязанский государственный медицинский университет) admission data.
// RZGMU admission lists are provided in PDF format and require Tabula for data extraction.
package rzgmu

import (
	"analabit/core"
	"analabit/core/source"
	"analabit/core/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sync/semaphore"
)

// Global semaphore to limit concurrent HTTP requests
var httpRequestSemaphore *semaphore.Weighted

func init() {
	// Default to 2 concurrent requests for RZGMU, but allow override via environment variable
	maxConcurrentRequests := int64(2)
	if envVal := os.Getenv("RZGMU_HTTP_MAX_CONCURRENT_REQUESTS"); envVal != "" {
		if parsed, err := strconv.ParseInt(envVal, 10, 64); err == nil && parsed > 0 {
			maxConcurrentRequests = parsed
		} else {
			log.Printf("Warning: Invalid RZGMU_HTTP_MAX_CONCURRENT_REQUESTS value '%s', using default %d", envVal, maxConcurrentRequests)
		}
	}

	httpRequestSemaphore = semaphore.NewWeighted(maxConcurrentRequests)
	log.Printf("Initialized RZGMU HTTP request semaphore with limit: %d concurrent requests", maxConcurrentRequests)
}

// HTTPHeadingSource defines how to load RZGMU heading data from a PDF file URL.
// RZGMU provides admission lists in PDF format that need to be processed with Tabula.
type HTTPHeadingSource struct {
	URL         string           // URL to the PDF file containing the admission list
	ProgramName string           // Name of the educational program (if empty, will use default)
	Capacities  *core.Capacities // Capacities for this heading. If nil, will be extracted from PDF
}

// LoadTo loads data from HTTP source, downloading PDF, extracting table data with Tabula,
// and sending HeadingData and ApplicationData to the provided receiver.
func (s *HTTPHeadingSource) LoadTo(receiver source.DataReceiver) error {
	if s.URL == "" {
		return fmt.Errorf("URL is required for RZGMU HttpHeadingSource")
	}

	log.Printf("Processing RZGMU admission list from: %s", s.URL)

	// Acquire a semaphore slot, respecting context cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := httpRequestSemaphore.Acquire(ctx, 1); err != nil {
		return fmt.Errorf("failed to acquire semaphore for RZGMU list from %s: %w", s.URL, err)
	}
	defer httpRequestSemaphore.Release(1)

	// Download PDF to temporary file
	tempPDFPath, err := downloadPDFToTemp(ctx, s.URL)
	if err != nil {
		return fmt.Errorf("failed to download RZGMU PDF from %s: %w", s.URL, err)
	}
	defer func() {
		if removeErr := os.Remove(tempPDFPath); removeErr != nil {
			log.Printf("Warning: failed to remove temporary PDF file %s: %v", tempPDFPath, removeErr)
		}
	}()

	// Extract CSV data using Tabula
	csvData, err := extractCSVWithTabula(ctx, tempPDFPath)
	if err != nil {
		return fmt.Errorf("failed to extract CSV from RZGMU PDF %s: %w", s.URL, err)
	}

	// Parse the CSV data
	programs, err := parseRZGMUTabulaData(csvData)
	if err != nil {
		return fmt.Errorf("failed to parse RZGMU data from %s: %w", s.URL, err)
	}

	// Process each program found in the data (usually just one per PDF)
	for _, program := range programs {
		// Use specified program name if provided, otherwise use extracted name
		programName := s.ProgramName
		if programName == "" {
			programName = program.PrettyName
		}

		capacities := s.Capacities
		if capacities == nil {
			// Extract capacities from the program data
			capacities = &program.ExtractedCapacities
		}

		headingCode := utils.GenerateHeadingCode(programName)

		// Send HeadingData to the receiver
		receiver.PutHeadingData(&source.HeadingData{
			Code:       headingCode,
			Capacities: *capacities,
			PrettyName: programName,
		})

		log.Printf("Sent RZGMU heading: %s (Code: %s, Caps: %v)", programName, headingCode, *capacities)

		// Send ApplicationData for each application in this program
		for _, app := range program.Applications {
			app.HeadingCode = headingCode
			receiver.PutApplicationData(app)
		}

		log.Printf("Sent %d applications for RZGMU heading %s", len(program.Applications), programName)
	}

	log.Printf("Successfully processed RZGMU heading(s) from %s", s.URL)
	return nil
}

// downloadPDFToTemp downloads a PDF from the given URL to a temporary file
func downloadPDFToTemp(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download PDF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download PDF (status code %d)", resp.StatusCode)
	}

	// Create temporary file
	tempFile, err := os.CreateTemp("", "rzgmu_*.pdf")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer tempFile.Close()

	// Copy response body to temporary file
	_, err = tempFile.ReadFrom(resp.Body)
	if err != nil {
		os.Remove(tempFile.Name())
		return "", fmt.Errorf("failed to write PDF to temporary file: %w", err)
	}

	return tempFile.Name(), nil
}

// extractCSVWithTabula uses Tabula CLI to extract CSV data from a PDF file
func extractCSVWithTabula(ctx context.Context, pdfPath string) (string, error) {
	// Find the Tabula JAR file - try multiple possible paths
	var tabulaJarPath string
	possiblePaths := []string{
		filepath.Join("tools", "tabula.jar"),                   // From project root
		filepath.Join("..", "..", "..", "tools", "tabula.jar"), // From test directory
		filepath.Join("../../../../tools", "tabula.jar"),       // From deeper test directory
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			tabulaJarPath = path
			break
		}
	}

	if tabulaJarPath == "" {
		return "", fmt.Errorf("tabula JAR not found at any of the expected paths: %v. Run 'make tools' to download it", possiblePaths)
	}

	// Prepare Tabula command
	// Using --format=CSV, --guess to automatically detect table areas, and --pages=all to parse all pages
	cmd := exec.CommandContext(ctx, "java", "-jar", tabulaJarPath,
		"--format=CSV",
		"--guess",
		"--pages=all",
		"--silent",
		pdfPath)

	// Execute the command and capture output
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("tabula extraction failed: %s (stderr: %s)", err, string(exitErr.Stderr))
		}
		return "", fmt.Errorf("failed to run Tabula: %w", err)
	}

	return string(output), nil
}

// ProgramData represents a single program's admission data
type ProgramData struct {
	PrettyName          string
	ExtractedCapacities core.Capacities
	Applications        []*source.ApplicationData
}

// parseRZGMUData parses the CSV data extracted from RZGMU PDF
func parseRZGMUData(csvData string) ([]ProgramData, error) {
	lines := strings.Split(csvData, "\n")
	var programs []ProgramData
	var currentProgram *ProgramData

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check if this line contains program information
		if strings.HasPrefix(line, "Направление подготовки:") {
			// If we already have a current program, save it before starting a new one
			if currentProgram != nil && len(currentProgram.Applications) > 0 {
				programs = append(programs, *currentProgram)
			}

			// Start new program
			prettyName := extractProgramName(line)
			if prettyName == "" {
				log.Printf("Warning: Could not extract program name from line %d: %s", i+1, line)
				continue
			}

			currentProgram = &ProgramData{
				PrettyName:   prettyName,
				Applications: make([]*source.ApplicationData, 0),
			}
			continue
		}

		// Check if this line contains capacity information
		if strings.Contains(line, "Мест:") && currentProgram != nil {
			capacity := extractCapacity(line)
			if capacity > 0 {
				// For RZGMU, we assume all places are Regular competition
				// This can be refined if we find quota information in the data
				currentProgram.ExtractedCapacities = core.Capacities{
					Regular:        capacity,
					TargetQuota:    0,
					DedicatedQuota: 0,
					SpecialQuota:   0,
				}
			}
			continue
		}

		// Check if this line contains table header
		if strings.Contains(line, "№,Код,Балл") {
			// This indicates we're starting the applications table
			continue
		}

		// Check if this line contains competition group info
		if strings.HasPrefix(line, "Конкурсная группа:") {
			// We can use this to determine competition type if needed
			continue
		}

		// Parse application data rows
		if currentProgram != nil && isApplicationDataRow(line) {
			app, err := parseApplicationRow(line)
			if err != nil {
				log.Printf("Warning: Failed to parse application row at line %d: %v (line: %s)", i+1, err, line)
				continue
			}
			if app != nil {
				currentProgram.Applications = append(currentProgram.Applications, app)
			}
		}
	}

	// Add the last program if it exists
	if currentProgram != nil && len(currentProgram.Applications) > 0 {
		programs = append(programs, *currentProgram)
	}

	if len(programs) == 0 {
		return nil, fmt.Errorf("no programs found in RZGMU data")
	}

	return programs, nil
}

// extractProgramName extracts the program name from a line like "Направление подготовки: Лечебное дело"
func extractProgramName(line string) string {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) < 2 {
		return ""
	}
	// Remove any trailing commas and trim whitespace
	name := strings.TrimSpace(parts[1])
	name = strings.TrimRight(name, ",")
	return name
}

// extractCapacity extracts the capacity number from a line containing "Мест: X"
func extractCapacity(line string) int {
	// Look for pattern "Мест: X" where X is a number
	if idx := strings.Index(line, "Мест:"); idx >= 0 {
		remaining := line[idx+5:] // Skip "Мест:"
		// Find the first number in the remaining string
		var numStr strings.Builder
		for _, r := range remaining {
			if r >= '0' && r <= '9' {
				numStr.WriteRune(r)
			} else if numStr.Len() > 0 {
				break // Stop at first non-digit after we've started collecting digits
			}
		}
		if numStr.Len() > 0 {
			if num, err := strconv.Atoi(numStr.String()); err == nil {
				return num
			}
		}
	}
	return 0
}

// isApplicationDataRow checks if a line contains application data
func isApplicationDataRow(line string) bool {
	// A data row should start with a number (rating place) and contain commas
	parts := strings.Split(line, ",")
	if len(parts) < 7 { // Need at least 7 columns: №,Код,Балл,ВИ,ИД,ПП,Приоритет,Согласие
		return false
	}

	// Check if first field is a number (rating place)
	firstField := strings.TrimSpace(parts[0])
	if _, err := strconv.Atoi(firstField); err != nil {
		return false
	}

	return true
}

// parseApplicationRow parses a single application data row
func parseApplicationRow(line string) (*source.ApplicationData, error) {
	parts := strings.Split(line, ",")
	if len(parts) < 8 {
		return nil, fmt.Errorf("insufficient columns in application row (expected at least 8, got %d)", len(parts))
	}

	// Parse rating place (№)
	ratingPlaceStr := strings.TrimSpace(parts[0])
	ratingPlace, err := strconv.Atoi(ratingPlaceStr)
	if err != nil {
		return nil, fmt.Errorf("invalid rating place '%s': %w", ratingPlaceStr, err)
	}

	// Parse student ID (Код)
	studentID := strings.TrimSpace(parts[1])
	if studentID == "" {
		return nil, fmt.Errorf("empty student ID")
	}

	// Parse total score (Балл)
	scoreStr := strings.TrimSpace(parts[2])
	var totalScore int
	if scoreStr != "" && scoreStr != "-" {
		totalScore, err = strconv.Atoi(scoreStr)
		if err != nil {
			log.Printf("Warning: invalid total score '%s' for student %s, using 0", scoreStr, studentID)
			totalScore = 0
		}
	}

	// Parse individual achievements (ИД) - column 4
	var bonusPoints int
	if len(parts) > 4 {
		bonusStr := strings.TrimSpace(parts[4])
		if bonusStr != "" {
			bonusPoints, err = strconv.Atoi(bonusStr)
			if err != nil {
				bonusPoints = 0
			}
		}
	}

	// Add bonus points to total score
	totalScore += bonusPoints

	// Parse priority (Приоритет) - column 6
	priority := 1 // Default priority
	if len(parts) > 6 {
		priorityStr := strings.TrimSpace(parts[6])
		if priorityStr != "" {
			priority, err = strconv.Atoi(priorityStr)
			if err != nil {
				priority = 1
			}
		}
	}

	// Parse consent status (Согласие) - column 7
	var originalSubmitted bool
	if len(parts) > 7 {
		consentStr := strings.TrimSpace(parts[7])
		originalSubmitted = (consentStr == "Согласие")
	}

	// Determine competition type
	competitionType := core.CompetitionRegular
	if len(parts) > 3 {
		examInfo := strings.TrimSpace(parts[3]) // ВИ column
		if examInfo == "БВИ" {
			competitionType = core.CompetitionBVI
		}
	}

	return &source.ApplicationData{
		HeadingCode:       "", // Will be set by caller
		StudentID:         studentID,
		ScoresSum:         totalScore,
		RatingPlace:       ratingPlace,
		Priority:          priority,
		CompetitionType:   competitionType,
		OriginalSubmitted: originalSubmitted,
	}, nil
}
