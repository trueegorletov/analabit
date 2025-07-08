package rzgmu

import (
	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// TextProgramData represents a single program's admission data with improved capacity tracking
type TextProgramData struct {
	PrettyName          string
	ExtractedCapacities core.Capacities
	Applications        []*source.ApplicationData
}

// TextElement represents a text element with its position
type TextElement struct {
	Text string
	X, Y float64
}

// detectCompetitionTypeHeader detects competition type from header lines
func detectCompetitionTypeHeader(line string) string {
	upperLine := strings.ToUpper(line)

	// Check for Regular/BVI section headers
	if strings.Contains(upperLine, "ЛЕЧЕБНОЕ ДЕЛО (БЮДЖЕТ)") ||
		(strings.Contains(upperLine, "ЛЕЧЕБНОЕ ДЕЛО") && strings.Contains(upperLine, "(СПЕЦИАЛИТЕТ)") &&
			!strings.Contains(upperLine, "КВОТА")) {
		return "REGULAR"
	}

	// Check for quota type keywords
	if strings.Contains(upperLine, "ЦЕЛЕВАЯ") && strings.Contains(upperLine, "КВОТА") {
		return "ЦЕЛЕВАЯ"
	}

	if strings.Contains(upperLine, "ОСОБАЯ") && strings.Contains(upperLine, "КВОТА") {
		return "ОСОБАЯ"
	}

	if strings.Contains(upperLine, "ОТДЕЛЬНАЯ") && strings.Contains(upperLine, "КВОТА") {
		return "ОТДЕЛЬНАЯ"
	}

	return ""
}

// findCapacityInNextLines looks for capacity information in the next few lines
func findCapacityInNextLines(lines []string, startIndex, maxLookahead int) int {
	for i := startIndex; i < len(lines) && i < startIndex+maxLookahead; i++ {
		line := strings.TrimSpace(lines[i])
		if capacity := extractCapacityFromLine(line); capacity > 0 {
			return capacity
		}
		// Stop looking if we hit another competition type header or application start
		if detectCompetitionTypeHeader(line) != "" || isApplicationStartLine(line) {
			break
		}
	}
	return 0
}

// isApplicationStartLine checks if a line starts an application entry
func isApplicationStartLine(line string) bool {
	// Look for patterns like "1. 3867113" or "1.     3867113"
	// But exclude column header lines that contain symbols like "№"
	if strings.Contains(line, "№") || strings.Contains(line, "Код") || strings.Contains(line, "Приоритет") {
		return false
	}
	re := regexp.MustCompile(`^\s*\d+\.\s+\d{7}`)
	return re.MatchString(line)
}

// parseRZGMUTextData parses the clean text extracted from PDF using rsc.io/pdf
// This handles the structured text format that is much cleaner than Tabula CSV output
func parseRZGMUTextData(textData string) ([]TextProgramData, error) {
	lines := strings.Split(textData, "\n")

	programName := "Лечебное дело" // Default for RZGMU medical programs
	var applications []*source.ApplicationData

	// Capacity tracking with summation instead of reset
	var totalCapacities core.Capacities

	// Track competition type detection
	currentCompetitionType := core.CompetitionRegular
	lastRatingPlace := 0
	justChangedCompetitionType := false // Track when we've recently changed competition types
	lastProposalLineIndex := -1         // Track when we last saw a "Предложение РвР" line

	// Application parsing state
	var currentApp *source.ApplicationData
	var examBuffer []string

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Extract program name from header
		if strings.HasPrefix(line, "Направление подготовки:") {
			extractedName := extractProgramNameFromText(line)
			if extractedName != "" {
				programName = extractedName
				log.Printf("Extracted program name: %s", programName)
			}
			continue
		}

		// Skip column headers
		if strings.Contains(line, "№ Код Балл ВИ ИД ПП Приоритет Согласие") {
			continue
		}

		// Check for competition type headers with lookahead for capacity
		if competitionType := detectCompetitionTypeHeader(line); competitionType != "" {
			// Look ahead for capacity in the next few lines
			capacity := findCapacityInNextLines(lines, i, 5)

			switch competitionType {
			case "REGULAR":
				currentCompetitionType = core.CompetitionRegular
				totalCapacities.Regular += capacity
				log.Printf("Detected Regular section with capacity %d (total: %d)", capacity, totalCapacities.Regular)
			case "ЦЕЛЕВАЯ":
				currentCompetitionType = core.CompetitionTargetQuota
				totalCapacities.TargetQuota += capacity
				log.Printf("Detected Target Quota section with capacity %d (total: %d)", capacity, totalCapacities.TargetQuota)
			case "ОСОБАЯ":
				currentCompetitionType = core.CompetitionSpecialQuota
				totalCapacities.SpecialQuota += capacity
				log.Printf("Detected Special Quota section with capacity %d (total: %d)", capacity, totalCapacities.SpecialQuota)
			case "ОТДЕЛЬНАЯ":
				currentCompetitionType = core.CompetitionDedicatedQuota
				totalCapacities.DedicatedQuota += capacity
				log.Printf("Detected Dedicated Quota section with capacity %d (total: %d)", capacity, totalCapacities.DedicatedQuota)
			}
			justChangedCompetitionType = true // Mark that we've changed competition types
			continue
		}

		// Check for standalone capacity lines that might have been missed
		if strings.HasPrefix(line, "Мест:") && !strings.Contains(line, "№") {
			capacity := extractCapacityFromLine(line)
			if capacity > 0 {
				// Add to current competition type
				switch currentCompetitionType {
				case core.CompetitionRegular:
					totalCapacities.Regular += capacity
					log.Printf("Added standalone capacity %d to Regular (total: %d)", capacity, totalCapacities.Regular)
				case core.CompetitionTargetQuota:
					totalCapacities.TargetQuota += capacity
					log.Printf("Added standalone capacity %d to Target Quota (total: %d)", capacity, totalCapacities.TargetQuota)
				case core.CompetitionSpecialQuota:
					totalCapacities.SpecialQuota += capacity
					log.Printf("Added standalone capacity %d to Special Quota (total: %d)", capacity, totalCapacities.SpecialQuota)
				case core.CompetitionDedicatedQuota:
					totalCapacities.DedicatedQuota += capacity
					log.Printf("Added standalone capacity %d to Dedicated Quota (total: %d)", capacity, totalCapacities.DedicatedQuota)
				}
			}
			continue
		}

		// Check for "Предложение РвР" lines which legitimately reset rankings
		if isProposalLine(line) {
			lastProposalLineIndex = i
			log.Printf("Detected administrative proposal line: %s", line)
			continue
		}

		// Skip metadata lines
		if isMetadataLineText(line) {
			continue
		}

		// Check if this line starts a new application (numbered entry)
		if appData := parseApplicationStart(line); appData != nil {
			// Save previous application if we were building one
			if currentApp != nil {
				// Apply current competition type only if not already set (e.g., for BVI)
				if currentApp.CompetitionType == 0 {
					currentApp.CompetitionType = currentCompetitionType
				}

				applications = append(applications, currentApp)

				// Check for rating place reset (indicates missing header)
				// Skip this check if we just changed competition types or recently saw a proposal
				recentProposal := lastProposalLineIndex >= 0 && (i-lastProposalLineIndex) <= 15
				if !justChangedCompetitionType && !recentProposal && lastRatingPlace > 0 && appData.RatingPlace < lastRatingPlace && appData.RatingPlace <= 5 {
					log.Printf("Warning: Rating place reset detected (from %d to %d) at line %d - possible missing competition type header",
						lastRatingPlace, appData.RatingPlace, i+1)
				}
			}

			// Start new application and update tracking variables
			currentApp = appData
			examBuffer = nil
			lastRatingPlace = appData.RatingPlace // Always update lastRatingPlace

			// Reset competition type flag only after first application in new section
			if justChangedCompetitionType {
				justChangedCompetitionType = false
			}
			continue
		}

		// Check if this line contains exam subject information
		if currentApp != nil && isExamSubjectLine(line) {
			examBuffer = append(examBuffer, line)
			continue
		}

		// If we have additional fields for the current application, parse them
		if currentApp != nil && len(examBuffer) > 0 {
			// We've finished collecting exam data, now parse the remaining fields
			if remainingFields := parseRemainingFields(line); remainingFields != nil {
				// Only update priority if it's still the default (1) and we found a better value
				if currentApp.Priority == 1 && remainingFields.Priority != 1 {
					currentApp.Priority = remainingFields.Priority
				}

				// Only update consent status if it was actually found in this line
				if remainingFields.ConsentFound {
					currentApp.OriginalSubmitted = remainingFields.OriginalSubmitted
				}

				// Apply current competition type only if not already set (e.g., for BVI)
				if currentApp.CompetitionType == 0 {
					currentApp.CompetitionType = currentCompetitionType
				}

				applications = append(applications, currentApp)
				lastRatingPlace = currentApp.RatingPlace
				currentApp = nil
				examBuffer = nil
			}
		}
	}

	// Don't forget the last application
	if currentApp != nil {
		// Apply current competition type only if not already set (e.g., for BVI)
		if currentApp.CompetitionType == 0 {
			currentApp.CompetitionType = currentCompetitionType
		}
		applications = append(applications, currentApp)
	}

	if len(applications) == 0 {
		return nil, fmt.Errorf("no applications found in RZGMU text data")
	}

	log.Printf("Parsed %d applications with capacities: Regular=%d, Target=%d, Special=%d, Dedicated=%d",
		len(applications), totalCapacities.Regular, totalCapacities.TargetQuota, totalCapacities.SpecialQuota, totalCapacities.DedicatedQuota)

	program := TextProgramData{
		PrettyName:          programName,
		ExtractedCapacities: totalCapacities,
		Applications:        applications,
	}

	return []TextProgramData{program}, nil
}

// extractProgramNameFromText extracts program name from header line
func extractProgramNameFromText(line string) string {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) < 2 {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

// detectCompetitionTypeAndCapacity detects competition type headers and extracts capacity
func detectCompetitionTypeAndCapacity(line string) (string, int) {
	upperLine := strings.ToUpper(line)

	// Check for Regular/BVI section headers
	if strings.Contains(upperLine, "ЛЕЧЕБНОЕ ДЕЛО (БЮДЖЕТ)") ||
		(strings.Contains(upperLine, "ЛЕЧЕБНОЕ ДЕЛО") && strings.Contains(upperLine, "(СПЕЦИАЛИТЕТ)") &&
			!strings.Contains(upperLine, "КВОТА")) {
		// Look for capacity in next lines or same line
		if capacity := extractCapacityFromLine(line); capacity > 0 {
			return "REGULAR", capacity
		}
		return "REGULAR", 0
	}

	// Check for quota type keywords
	if strings.Contains(upperLine, "ЦЕЛЕВАЯ") && strings.Contains(upperLine, "КВОТА") {
		capacity := extractCapacityFromLine(line)
		return "ЦЕЛЕВАЯ", capacity
	}

	if strings.Contains(upperLine, "ОСОБАЯ") && strings.Contains(upperLine, "КВОТА") {
		capacity := extractCapacityFromLine(line)
		return "ОСОБАЯ", capacity
	}

	if strings.Contains(upperLine, "ОТДЕЛЬНАЯ") && strings.Contains(upperLine, "КВОТА") {
		capacity := extractCapacityFromLine(line)
		return "ОТДЕЛЬНАЯ", capacity
	}

	// Check for standalone capacity lines
	if strings.HasPrefix(line, "Мест:") {
		capacity := extractCapacityFromLine(line)
		if capacity > 0 {
			// This is a capacity line following a header - determine type based on previous context
			// For now, we'll treat standalone capacity as regular
			return "REGULAR", capacity
		}
	}

	return "", 0
}

// extractCapacityFromLine extracts capacity number from lines containing "Мест: X"
func extractCapacityFromLine(line string) int {
	re := regexp.MustCompile(`Мест:\s*(\d+)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) > 1 {
		if capacity, err := strconv.Atoi(matches[1]); err == nil {
			return capacity
		}
	}
	return 0
}

// isMetadataLineText checks if a line contains metadata
func isMetadataLineText(line string) bool {
	lowerLine := strings.ToLower(line)
	return strings.Contains(lowerLine, "конкурсная группа") ||
		strings.Contains(lowerLine, "основные места") ||
		strings.Contains(lowerLine, "целевое обучение") ||
		strings.Contains(lowerLine, "особые права") ||
		strings.Contains(lowerLine, "без вступительных испытаний") ||
		strings.Contains(lowerLine, "льготы")
}

// parseApplicationStart parses the start of an application entry
func parseApplicationStart(line string) *source.ApplicationData {
	// Match pattern like "1. 3867113 - БВИ 10 1 Согласие" or "1. 3785711 - БПВИ 8 Пр.право 8 Нет"
	re := regexp.MustCompile(`^(\d+)\.\s+(\d+)\s+(-|\d+)\s*(.*)`)
	matches := re.FindStringSubmatch(line)

	if len(matches) < 4 {
		return nil
	}

	ratingPlace, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}

	studentID := matches[2]
	scoreStr := matches[3]
	remaining := strings.TrimSpace(matches[4])

	// Parse remaining fields to detect БВИ first
	fields := strings.Fields(remaining)
	isBVI, isBPVI := false, false
	for _, field := range fields {
		if field == "БВИ" {
			isBVI = true
		} else if field == "БПВИ" {
			isBPVI = true
		}
	}

	// Handle score parsing - only use the third column ("Балл")
	var totalScore int
	if scoreStr == "-" {
		// Missing score - set defaults based on competition type
		if isBVI {
			totalScore = 310 // Default for БВИ
			log.Printf("Missing score for БВИ student %s, setting to 310", studentID)
		} else if isBPVI {
			totalScore = 282 // Default for БПВИ
			log.Printf("Missing score for БПВИ student %s, setting to 310", studentID)
		} else {
			totalScore = 0 // Default for others
			log.Printf("Missing score for student %s, setting to 0", studentID)
		}
	} else {
		totalScore, err = strconv.Atoi(scoreStr)
		if err != nil {
			log.Printf("Invalid score '%s' for student %s, setting to 0", scoreStr, studentID)
			totalScore = 0
		}
	}

	app := &source.ApplicationData{
		StudentID:   studentID,
		RatingPlace: ratingPlace,
		ScoresSum:   totalScore, // Only use the third column value
		Priority:    1,          // Default
	}

	// First pass: find the consent field position ("Согласие" or "Нет")
	consentIndex := -1
	for j, field := range fields {
		if field == "Согласие" || field == "Нет" {
			consentIndex = j
			break
		}
	}

	// Parse other fields (БВИ, consent, priority) but NOT add to ScoresSum
	for i, field := range fields {
		if field == "БВИ" {
			app.CompetitionType = core.CompetitionBVI
		} else if field == "Согласие" {
			app.OriginalSubmitted = true
		} else if field == "Нет" {
			app.OriginalSubmitted = false
		} else if field == "Пр.право" {
			// Preferential right indicator
		} else if num, err := strconv.Atoi(field); err == nil {
			// If this numeric field is immediately before consent, it's the priority
			if consentIndex >= 0 && i == consentIndex-1 {
				app.Priority = num
			}
		}
	}

	return app
}

// isExamSubjectLine checks if a line contains exam subject information
func isExamSubjectLine(line string) bool {
	lowerLine := strings.ToLower(line)
	return strings.Contains(lowerLine, "биология") ||
		strings.Contains(lowerLine, "химия") ||
		strings.Contains(lowerLine, "русский язык") ||
		strings.Contains(lowerLine, "математика") ||
		strings.Contains(lowerLine, "физика")
}

// parseExamScores parses exam scores from collected subject lines
func parseExamScores(examLines []string) int {
	totalScore := 0
	re := regexp.MustCompile(`(\d+)`)

	for _, line := range examLines {
		matches := re.FindAllString(line, -1)
		for _, match := range matches {
			if score, err := strconv.Atoi(match); err == nil {
				totalScore += score
			}
		}
	}

	return totalScore
}

// RemainingFields holds the remaining application fields
type RemainingFields struct {
	BonusPoints       int
	Priority          int
	OriginalSubmitted bool
	ConsentFound      bool // Whether consent field was actually found in this line
}

// parseRemainingFields parses remaining fields from a line
func parseRemainingFields(line string) *RemainingFields {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return nil
	}

	// Skip header lines that contain column names
	if isHeaderLine(line) {
		return nil
	}

	result := &RemainingFields{
		Priority: 1, // Default
	}

	// First pass: find the consent field position ("Согласие" or "Нет")
	consentIndex := -1
	for j, field := range fields {
		if field == "Согласие" || field == "Нет" {
			consentIndex = j
			break
		}
	}

	// Parse fields using the same logic as parseApplicationStart
	for i, field := range fields {
		if field == "Согласие" {
			result.OriginalSubmitted = true
			result.ConsentFound = true
		} else if field == "Нет" {
			result.OriginalSubmitted = false
			result.ConsentFound = true
		} else if field == "Пр.право" {
			// Preferential right indicator
		} else if num, err := strconv.Atoi(field); err == nil {
			// If this numeric field is immediately before consent, it's the priority
			if consentIndex >= 0 && i == consentIndex-1 {
				result.Priority = num
			} else if num > 20 {
				// Likely bonus points (higher values)
				result.BonusPoints = num
			}
			// Don't assign random numeric values as priority unless they're in the right position
		}
	}

	return result
}

// DebugTextExtraction helps debug PDF text extraction by comparing with sample data
func DebugTextExtraction(pdfPath, sampleTextPath string) error {
	log.Printf("=== DEBUG: Comparing PDF extraction with sample text ===")

	// Extract text from PDF using rsc.io/pdf
	extractedText, err := extractTextFromPDFFile(pdfPath)
	if err != nil {
		return fmt.Errorf("failed to extract text from PDF: %w", err)
	}

	// Read sample text file if provided
	var sampleText string
	if sampleTextPath != "" {
		data, err := os.ReadFile(sampleTextPath)
		if err != nil {
			return fmt.Errorf("failed to read sample text file: %w", err)
		}
		sampleText = string(data)
	}

	// Show first few lines of each
	extractedLines := strings.Split(extractedText, "\n")
	log.Printf("=== EXTRACTED TEXT (first 20 lines) ===")
	for i, line := range extractedLines {
		if i >= 20 {
			break
		}
		log.Printf("%3d: %s", i+1, line)
	}

	if sampleText != "" {
		sampleLines := strings.Split(sampleText, "\n")
		log.Printf("=== SAMPLE TEXT (first 20 lines) ===")
		for i, line := range sampleLines {
			if i >= 20 {
				break
			}
			log.Printf("%3d: %s", i+1, line)
		}
	}

	// Try parsing both
	log.Printf("=== PARSING EXTRACTED TEXT ===")
	extractedPrograms, err := parseRZGMUTextData(extractedText)
	if err != nil {
		log.Printf("ERROR parsing extracted text: %v", err)
	} else {
		for _, prog := range extractedPrograms {
			log.Printf("Extracted - Program: %s, Apps: %d, Capacities: R=%d T=%d S=%d D=%d",
				prog.PrettyName, len(prog.Applications),
				prog.ExtractedCapacities.Regular, prog.ExtractedCapacities.TargetQuota,
				prog.ExtractedCapacities.SpecialQuota, prog.ExtractedCapacities.DedicatedQuota)
		}
	}

	if sampleText != "" {
		log.Printf("=== PARSING SAMPLE TEXT ===")
		samplePrograms, err := parseRZGMUTextData(sampleText)
		if err != nil {
			log.Printf("ERROR parsing sample text: %v", err)
		} else {
			for _, prog := range samplePrograms {
				log.Printf("Sample - Program: %s, Apps: %d, Capacities: R=%d T=%d S=%d D=%d",
					prog.PrettyName, len(prog.Applications),
					prog.ExtractedCapacities.Regular, prog.ExtractedCapacities.TargetQuota,
					prog.ExtractedCapacities.SpecialQuota, prog.ExtractedCapacities.DedicatedQuota)
			}
		}
	}

	return nil
}

// isProposalLine checks if a line contains "Предложение РвР" which legitimately resets rankings
func isProposalLine(line string) bool {
	return strings.Contains(line, "Предложение РвР")
}

// isHeaderLine checks if a line contains column headers like "№ Код Балл ВИ ИД ПП Приоритет Согласие"
func isHeaderLine(line string) bool {
	lowerLine := strings.ToLower(line)

	// Check for common column header combinations
	// A header line typically contains multiple of these column names
	headerKeywords := []string{"№", "код", "балл", "приоритет"}
	matchCount := 0

	for _, keyword := range headerKeywords {
		if strings.Contains(lowerLine, keyword) {
			matchCount++
		}
	}

	// If it contains 3 or more header keywords, it's likely a header line
	return matchCount >= 2
}
