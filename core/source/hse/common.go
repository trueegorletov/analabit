// Package hse provides support for the new format of HSE admission lists.
// Unlike the old format (oldhse), the new format groups all applications
// for one program into a single XLSX file, rather than separate files
// for each competition type.
package hse

import (
	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

// extractPrettyNameFromXLSX extracts the heading name from the first row of the XLSX file.
// The format is: "Образовательная программа "название в кавычках""
func extractPrettyNameFromXLSX(f *excelize.File, sourceHint string) (string, error) {
	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		return "", fmt.Errorf("no sheets found in Excel file (source: %s)", sourceHint)
	}
	sheetName := sheetList[0]

	// Try to get cell A1 first (where the heading name usually is)
	cellValue, err := f.GetCellValue(sheetName, "A1")
	if err != nil {
		return "", fmt.Errorf("failed to get cell A1 from Excel file (source: %s, sheet: %s): %w", sourceHint, sheetName, err)
	}

	// Extract name from quotes using regex
	re := regexp.MustCompile(`"([^"]*)"`)
	matches := re.FindStringSubmatch(cellValue)
	if len(matches) >= 2 {
		return strings.TrimSpace(matches[1]), nil
	}

	// Fallback: try to extract after "Образовательная программа "
	parts := strings.SplitN(cellValue, "Образовательная программа ", 2)
	if len(parts) == 2 {
		// Remove any quotes
		name := strings.Trim(parts[1], `"`)
		return strings.TrimSpace(name), nil
	}

	return "", fmt.Errorf("could not extract pretty name from '%s' in Excel file (source: %s, sheet: %s, cell A1)", cellValue, sourceHint, sheetName)
}

// parseApplicationsFromXLSX parses the XLSX data and sends application data to the receiver.
func parseApplicationsFromXLSX(f *excelize.File, headingCode string, receiver source.DataReceiver, sourceHint string) error {
	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		return fmt.Errorf("no sheets found in Excel file (source: %s)", sourceHint)
	}
	sheetName := sheetList[0]

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("failed to get rows from Excel file (source: %s, sheet: %s): %w", sourceHint, sheetName, err)
	}

	if len(rows) < 15 {
		log.Printf("Excel file from %s has fewer than 15 rows, no application data to process", sourceHint)
		return nil
	}

	// Find the header row (should be around row 14, 0-indexed as row 13)
	headerRowIndex := -1
	for i := 10; i < min(len(rows), 20); i++ {
		if len(rows[i]) > 1 && rows[i][0] == "№ п/п" {
			headerRowIndex = i
			break
		}
	}

	if headerRowIndex == -1 {
		return fmt.Errorf("could not find header row with '№ п/п' in Excel file from %s", sourceHint)
	}

	header := rows[headerRowIndex]

	// Find column indices
	colIndices := findColumnIndices(header)
	if colIndices.positionCol == -1 || colIndices.studentIDCol == -1 {
		return fmt.Errorf("could not find required columns in Excel file from %s", sourceHint)
	}

	// Process data rows (starting after header)
	dataStartIndex := headerRowIndex + 1
	for i := dataStartIndex; i < len(rows); i++ {
		row := rows[i]
		if len(row) <= colIndices.maxNeededCol() {
			continue // Skip rows that don't have enough columns
		}

		app, err := parseApplicationFromRow(row, colIndices, headingCode)
		if err != nil {
			log.Printf("Error parsing row %d from %s: %v", i+1, sourceHint, err)
			continue // Skip invalid rows
		}

		if app != nil {
			receiver.PutApplicationData(app)
		}
	}

	return nil
}

// columnIndices holds the indices of important columns in the CSV
type columnIndices struct {
	positionCol       int // № п/п
	studentIDCol      int // Уникальный идентификатор
	bviCol            int // Право поступления без вступительных испытаний
	specialQuotaCol   int // Поступление на места в рамках квоты для лиц, имеющих особое право
	targetQuotaCol    int // Поступление на места в рамках квоты целевого приема
	dedicatedQuotaCol int // Поступление на места в рамках отдельной квоты
	targetPriorityCol int // Приоритет целевой квоты
	otherPriorityCol  int // Приоритет иных мест
	paidPriorityCol   int // Приоритет платных мест
	totalScoreCol     int // Сумма конкурсных баллов
	consentCol        int // Наличие согласия на зачисление
}

func (c columnIndices) maxNeededCol() int {
	max := c.positionCol
	if c.studentIDCol > max {
		max = c.studentIDCol
	}
	if c.bviCol > max {
		max = c.bviCol
	}
	if c.specialQuotaCol > max {
		max = c.specialQuotaCol
	}
	if c.targetQuotaCol > max {
		max = c.targetQuotaCol
	}
	if c.dedicatedQuotaCol > max {
		max = c.dedicatedQuotaCol
	}
	if c.targetPriorityCol > max {
		max = c.targetPriorityCol
	}
	if c.otherPriorityCol > max {
		max = c.otherPriorityCol
	}
	if c.paidPriorityCol > max {
		max = c.paidPriorityCol
	}
	if c.totalScoreCol > max {
		max = c.totalScoreCol
	}
	if c.consentCol > max {
		max = c.consentCol
	}
	return max
}

// findColumnIndices finds the indices of important columns in the header row
func findColumnIndices(header []string) columnIndices {
	indices := columnIndices{
		positionCol:       -1,
		studentIDCol:      -1,
		bviCol:            -1,
		specialQuotaCol:   -1,
		targetQuotaCol:    -1,
		dedicatedQuotaCol: -1,
		targetPriorityCol: -1,
		otherPriorityCol:  -1,
		paidPriorityCol:   -1,
		totalScoreCol:     -1,
		consentCol:        -1,
	}

	for i, col := range header {
		col = strings.TrimSpace(col)
		switch {
		case col == "№ п/п":
			indices.positionCol = i
		case col == "Уникальный идентификатор":
			indices.studentIDCol = i
		case strings.Contains(col, "Право поступления") && strings.Contains(col, "без вступительных испытаний"):
			indices.bviCol = i
		case strings.Contains(col, "Поступление на места в рамках квоты") && strings.Contains(col, "особое право"):
			indices.specialQuotaCol = i
		case strings.Contains(col, "Поступление на места в рамках квоты") && strings.Contains(col, "целевого приема"):
			indices.targetQuotaCol = i
		case strings.Contains(col, "Поступление на места") && strings.Contains(col, "отдельной квоты"):
			indices.dedicatedQuotaCol = i
		case strings.Contains(col, "Приоритет целевой квоты"):
			indices.targetPriorityCol = i
		case strings.Contains(col, "Приоритет иных мест"):
			indices.otherPriorityCol = i
		case strings.Contains(col, "Приоритет платных мест"):
			indices.paidPriorityCol = i
		case strings.Contains(col, "Сумма конкурсных баллов") && !strings.Contains(col, "целевой квоте"):
			indices.totalScoreCol = i
		case strings.Contains(col, "Наличие согласия на зачисление"):
			indices.consentCol = i
		}
	}

	return indices
}

// parseApplicationFromRow parses a single CSV row into ApplicationData
func parseApplicationFromRow(row []string, indices columnIndices, headingCode string) (*source.ApplicationData, error) {
	// Extract position (used as RatingPlace)
	positionStr := strings.TrimSpace(row[indices.positionCol])
	if positionStr == "" {
		return nil, nil // Skip empty rows
	}

	// Remove commas from position numbers (e.g., "1,000" -> "1000")
	positionStr = strings.ReplaceAll(positionStr, ",", "")

	position, err := strconv.Atoi(positionStr)
	if err != nil {
		return nil, fmt.Errorf("invalid position number: %s", positionStr)
	}

	// Extract student ID
	studentID := strings.TrimSpace(row[indices.studentIDCol])
	if studentID == "" {
		return nil, fmt.Errorf("empty student ID in row with position %d", position)
	}

	// Determine competition type based on quota columns
	competitionType := determineCompetitionType(row, indices)

	// Extract priority based on competition type
	priority := extractPriority(row, indices, competitionType)

	// Extract total score
	var totalScore int
	if indices.totalScoreCol >= 0 && indices.totalScoreCol < len(row) {
		scoreStr := strings.TrimSpace(row[indices.totalScoreCol])
		if scoreStr != "" {
			totalScore, err = strconv.Atoi(scoreStr)
			if err != nil {
				log.Printf("Warning: invalid total score '%s' for student %s, using 0", scoreStr, studentID)
				totalScore = 0
			}
		}
	}

	// Extract consent status
	var originalSubmitted bool
	if indices.consentCol >= 0 && indices.consentCol < len(row) {
		consentStr := strings.TrimSpace(row[indices.consentCol])
		originalSubmitted = (consentStr == "Да")
	}

	return &source.ApplicationData{
		HeadingCode:       headingCode,
		StudentID:         studentID,
		ScoresSum:         totalScore,
		RatingPlace:       position,
		Priority:          priority,
		CompetitionType:   competitionType,
		OriginalSubmitted: originalSubmitted,
	}, nil
}

// determineCompetitionType determines the competition type based on quota columns
func determineCompetitionType(row []string, indices columnIndices) core.Competition {
	// Check BVI first
	if indices.bviCol >= 0 && indices.bviCol < len(row) {
		bviStr := strings.TrimSpace(row[indices.bviCol])
		if bviStr == "Да" {
			return core.CompetitionBVI
		}
	}

	// Check Special Quota
	if indices.specialQuotaCol >= 0 && indices.specialQuotaCol < len(row) {
		specialStr := strings.TrimSpace(row[indices.specialQuotaCol])
		if specialStr == "Да" {
			return core.CompetitionSpecialQuota
		}
	}

	// Check Target Quota
	if indices.targetQuotaCol >= 0 && indices.targetQuotaCol < len(row) {
		targetStr := strings.TrimSpace(row[indices.targetQuotaCol])
		if targetStr == "Да" {
			return core.CompetitionTargetQuota
		}
	}

	// Check Dedicated Quota
	if indices.dedicatedQuotaCol >= 0 && indices.dedicatedQuotaCol < len(row) {
		dedicatedStr := strings.TrimSpace(row[indices.dedicatedQuotaCol])
		if dedicatedStr == "Да" {
			return core.CompetitionDedicatedQuota
		}
	}

	// Default to Regular competition
	return core.CompetitionRegular
}

// extractPriority extracts the priority based on competition type
func extractPriority(row []string, indices columnIndices, competitionType core.Competition) int {
	var priorityCol int

	switch competitionType {
	case core.CompetitionTargetQuota:
		priorityCol = indices.targetPriorityCol
	case core.CompetitionRegular, core.CompetitionBVI, core.CompetitionDedicatedQuota, core.CompetitionSpecialQuota:
		priorityCol = indices.otherPriorityCol
	default:
		priorityCol = indices.otherPriorityCol
	}

	if priorityCol >= 0 && priorityCol < len(row) {
		priorityStr := strings.TrimSpace(row[priorityCol])
		if priorityStr != "" {
			if priority, err := strconv.Atoi(priorityStr); err == nil {
				return priority
			}
		}
	}

	return 1 // Default priority
}

// min returns the minimum of two integers (Go < 1.21 compatibility)
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
