package rzgmu

import (
	"analabit/core"
	"analabit/core/source"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// parseRZGMUTabulaData parses the CSV data extracted directly from Tabula
// Unlike manual conversion, Tabula output has a different structure:
// - No clear program name header (assumed to be one program per PDF)
// - Capacity info is mangled across multiple cells
// - Application data rows start with numbers followed by dots
// - Multiple pages with similar table structures
// - Competition type changes are indicated by rating place resets and header rows
func parseRZGMUTabulaData(csvData string) ([]ProgramData, error) {
	lines := strings.Split(csvData, "\n")

	// Since each PDF is for one educational program, we'll collect all applications
	// and create a single program with a default name (can be overridden by FileHeadingSource)
	programName := "Лечебное дело" // Default for RZGMU medical programs
	var capacity int
	var applications []*source.ApplicationData

	// Track competition type detection
	currentCompetitionType := core.CompetitionRegular
	lastRatingPlace := 0
	pendingCompetitionTypeChange := false

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Skip header lines
		if isHeaderLine(line) {
			continue
		}

		// Check for competition type headers (ЦЕЛЕВАЯ, ОСОБАЯ, ОТДЕЛЬНАЯ)
		if competitionType := detectCompetitionTypeHeader(line); competitionType != "" {
			switch competitionType {
			case "ЦЕЛЕВАЯ":
				currentCompetitionType = core.CompetitionTargetQuota
				log.Printf("Detected Target Quota section at line %d: %s", i+1, line)
			case "ОСОБАЯ":
				currentCompetitionType = core.CompetitionSpecialQuota
				log.Printf("Detected Special Quota section at line %d: %s", i+1, line)
			case "ОТДЕЛЬНАЯ":
				currentCompetitionType = core.CompetitionDedicatedQuota
				log.Printf("Detected Dedicated Quota section at line %d: %s", i+1, line)
			}
			pendingCompetitionTypeChange = true // Set flag that we found a header
			continue
		}

		// Extract capacity if available in the line
		if extractedCapacity := extractCapacityFromTabulaLine(line); extractedCapacity > 0 {
			capacity = extractedCapacity
			continue
		}

		// Skip competition group info and target quota company info
		if isMetadataLine(line) {
			continue
		}

		// Try to parse as application data
		if isApplicationDataRowTabula(line) {
			app, err := parseApplicationRowTabula(line)
			if err != nil {
				log.Printf("Warning: Failed to parse Tabula application row at line %d: %v (line: %s)", i+1, err, line)
				continue
			}
			if app != nil {
				// Check for rating place reset (indicates new competition type)
				if lastRatingPlace > 0 && app.RatingPlace < lastRatingPlace && app.RatingPlace <= 5 {
					// Rating place reset detected - there should have been a header
					if !pendingCompetitionTypeChange {
						log.Printf("Warning: Rating place reset detected (from %d to %d) at line %d but no competition type header was found",
							lastRatingPlace, app.RatingPlace, i+1)
						log.Printf("Warning: Previous line context around reset: line %d", i)

						// Try to infer competition type from context or keep current
						// For now, we'll assume we missed a header and continue with current type
						// In future iterations, we could implement smarter inference
					}
					pendingCompetitionTypeChange = false
				}

				// Apply current competition type to the application
				app.CompetitionType = currentCompetitionType
				applications = append(applications, app)
				lastRatingPlace = app.RatingPlace
			}
		}
	}

	if len(applications) == 0 {
		return nil, fmt.Errorf("no applications found in RZGMU Tabula data")
	}

	// If no capacity was found, use a reasonable default or estimate
	if capacity == 0 {
		// Try to estimate capacity from the data structure or use a default
		capacity = estimateCapacityFromApplications(applications)
	}

	// Calculate capacities based on detected competition types
	regularCap, targetCap, specialCap, dedicatedCap := calculateCapacitiesByType(applications, capacity)

	program := ProgramData{
		PrettyName: programName,
		ExtractedCapacities: core.Capacities{
			Regular:        regularCap,
			TargetQuota:    targetCap,
			DedicatedQuota: dedicatedCap,
			SpecialQuota:   specialCap,
		},
		Applications: applications,
	}

	return []ProgramData{program}, nil
}

// isHeaderLine checks if a line is a table header
func isHeaderLine(line string) bool {
	upperLine := strings.ToUpper(line)
	return strings.Contains(upperLine, "NO,КОД,БАЛЛ") ||
		strings.Contains(upperLine, "№,КОД,БАЛЛ") ||
		strings.Contains(upperLine, "ВИ,ИД,ПП,ПРИОРИТЕТ,СОГЛАСИЕ")
}

// extractCapacityFromTabulaLine tries to extract capacity from a Tabula line
func extractCapacityFromTabulaLine(line string) int {
	// In Tabula output, capacity might appear as "Мест: 19" in various positions
	if strings.Contains(line, "Мест:") {
		return extractCapacity(line)
	}
	return 0
}

// isMetadataLine checks if a line contains metadata (competition group, target quota company info, etc.)
func isMetadataLine(line string) bool {
	lowerLine := strings.ToLower(line)
	return strings.Contains(lowerLine, "конкурсная группа") ||
		strings.Contains(lowerLine, "основные места") ||
		strings.Contains(lowerLine, "целевое обучение") ||
		strings.Contains(lowerLine, "особые права") ||
		strings.Contains(lowerLine, "без вступительных испытаний") ||
		strings.Contains(lowerLine, "льготы") ||
		// Target quota company patterns
		strings.Contains(lowerLine, "министерство") ||
		strings.Contains(lowerLine, "департамент") ||
		strings.Contains(lowerLine, "управление") ||
		strings.Contains(lowerLine, "комитет") ||
		strings.Contains(lowerLine, "администрация") ||
		(strings.Contains(lowerLine, "область") && !strings.Contains(line, ",")) ||
		(strings.Contains(lowerLine, "район") && !strings.Contains(line, ","))
}

// isApplicationDataRowTabula checks if a line contains application data (Tabula format)
func isApplicationDataRowTabula(line string) bool {
	parts := strings.Split(line, ",")
	if len(parts) < 6 { // Need at least 6 columns for basic data
		return false
	}

	// Check if first field is a number with optional dot (e.g., "1.", "2.", "10")
	firstField := strings.TrimSpace(parts[0])
	firstField = strings.TrimSuffix(firstField, ".") // Remove trailing dot

	if _, err := strconv.Atoi(firstField); err != nil {
		return false
	}

	// Additional validation: second field should look like a student ID
	secondField := strings.TrimSpace(parts[1])
	if len(secondField) < 4 || len(secondField) > 10 {
		return false
	}

	// Check if it's all digits (student ID should be numeric)
	if _, err := strconv.Atoi(secondField); err != nil {
		return false
	}

	return true
}

// parseApplicationRowTabula parses a single application data row (Tabula format)
func parseApplicationRowTabula(line string) (*source.ApplicationData, error) {
	parts := strings.Split(line, ",")
	if len(parts) < 6 {
		return nil, fmt.Errorf("insufficient columns in Tabula application row (expected at least 6, got %d)", len(parts))
	}

	// Parse rating place (№) - remove trailing dot if present
	ratingPlaceStr := strings.TrimSpace(parts[0])
	ratingPlaceStr = strings.TrimSuffix(ratingPlaceStr, ".")
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

	// Determine competition type from exam info (ВИ column - usually parts[3])
	competitionType := core.CompetitionRegular
	if len(parts) > 3 {
		examInfo := strings.TrimSpace(parts[3])
		examInfoLower := strings.ToLower(examInfo)

		if examInfo == "БВИ" || strings.Contains(examInfoLower, "без вступительных") {
			competitionType = core.CompetitionBVI
		} else if strings.Contains(examInfoLower, "целевое") || strings.Contains(examInfoLower, "целевой") {
			competitionType = core.CompetitionTargetQuota
		} else if strings.Contains(examInfoLower, "особые права") || strings.Contains(examInfoLower, "льгот") {
			competitionType = core.CompetitionDedicatedQuota
		} else if strings.Contains(examInfoLower, "специаль") {
			competitionType = core.CompetitionSpecialQuota
		}
	}

	// Parse individual achievements (ИД) - usually in column 4
	var bonusPoints int
	if len(parts) > 4 {
		bonusStr := strings.TrimSpace(parts[4])
		if bonusStr != "" && bonusStr != "-" {
			bonusPoints, err = strconv.Atoi(bonusStr)
			if err != nil {
				bonusPoints = 0
			}
		}
	}

	// Add bonus points to total score
	totalScore += bonusPoints

	// Parse priority - look for it in the available columns (usually column 6)
	priority := 1 // Default priority
	if len(parts) > 6 {
		priorityStr := strings.TrimSpace(parts[6])
		if priorityStr != "" && priorityStr != "-" {
			priority, err = strconv.Atoi(priorityStr)
			if err != nil {
				priority = 1
			}
		}
	}

	// Parse consent status - usually in the last meaningful column (column 7)
	var originalSubmitted bool
	if len(parts) > 7 {
		consentStr := strings.TrimSpace(parts[7])
		originalSubmitted = (consentStr == "Согласие")
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

// estimateCapacityFromApplications tries to estimate capacity from application data
func estimateCapacityFromApplications(applications []*source.ApplicationData) int {
	// Count different competition types to get a better estimate
	regularCount := 0
	bviCount := 0
	quotaCount := 0

	for _, app := range applications {
		switch app.CompetitionType {
		case core.CompetitionRegular:
			regularCount++
		case core.CompetitionBVI:
			bviCount++
		case core.CompetitionTargetQuota, core.CompetitionDedicatedQuota, core.CompetitionSpecialQuota:
			quotaCount++
		}
	}

	// Estimate based on typical admission patterns
	// Usually regular capacity is the largest portion
	if regularCount > 0 {
		// Assume regular places are about 70-80% of total capacity
		return int(float64(regularCount) / 0.75)
	}

	// Fallback: use total application count as an upper bound estimate
	return len(applications)
}

// detectCompetitionTypeHeader detects competition type headers containing keywords
func detectCompetitionTypeHeader(line string) string {
	upperLine := strings.ToUpper(line)

	// Check for target quota keywords
	if strings.Contains(upperLine, "ЦЕЛЕВАЯ") && (strings.Contains(upperLine, "КВОТА") || strings.Contains(upperLine, "ОБУЧЕНИЕ")) {
		return "ЦЕЛЕВАЯ"
	}

	// Check for special quota keywords
	if strings.Contains(upperLine, "ОСОБАЯ") && strings.Contains(upperLine, "КВОТА") {
		return "ОСОБАЯ"
	}

	// Check for dedicated quota keywords
	if strings.Contains(upperLine, "ОТДЕЛЬНАЯ") && strings.Contains(upperLine, "КВОТА") {
		return "ОТДЕЛЬНАЯ"
	}

	// Additional patterns that might indicate quota headers
	if strings.Contains(upperLine, "УФСИН") || strings.Contains(upperLine, "МИНИСТЕРСТВО") ||
		strings.Contains(upperLine, "ДЕПАРТАМЕНТ") || strings.Contains(upperLine, "УПРАВЛЕНИЕ") {
		// These often appear in target quota headers
		if strings.Contains(upperLine, "ЦЕЛЕВАЯ") || strings.Contains(upperLine, "ОБУЧЕНИЕ") {
			return "ЦЕЛЕВАЯ"
		}
	}

	return ""
}

// calculateCapacitiesByType calculates capacities based on detected competition types
func calculateCapacitiesByType(applications []*source.ApplicationData, totalCapacity int) (regular, target, special, dedicated int) {
	// Count applications by competition type
	regularCount := 0
	targetCount := 0
	specialCount := 0
	dedicatedCount := 0
	bviCount := 0

	for _, app := range applications {
		switch app.CompetitionType {
		case core.CompetitionRegular:
			regularCount++
		case core.CompetitionTargetQuota:
			targetCount++
		case core.CompetitionSpecialQuota:
			specialCount++
		case core.CompetitionDedicatedQuota:
			dedicatedCount++
		case core.CompetitionBVI:
			bviCount++
		}
	}

	// For capacity calculation, we need to be smart about estimating
	// BVI applications typically don't count against the main capacity

	if totalCapacity > 0 {
		// If we have a known total capacity, distribute based on proportions
		nonBviApps := regularCount + targetCount + specialCount + dedicatedCount
		if nonBviApps > 0 {
			// Calculate proportions
			regular = int(float64(regularCount) / float64(nonBviApps) * float64(totalCapacity))
			target = int(float64(targetCount) / float64(nonBviApps) * float64(totalCapacity))
			special = int(float64(specialCount) / float64(nonBviApps) * float64(totalCapacity))
			dedicated = int(float64(dedicatedCount) / float64(nonBviApps) * float64(totalCapacity))

			// Adjust for rounding errors
			calculated := regular + target + special + dedicated
			if calculated < totalCapacity {
				regular += totalCapacity - calculated
			}
		} else {
			// Only BVI applications found, set regular capacity to total
			regular = totalCapacity
		}
	} else {
		// No total capacity known, use application counts as estimates
		// Add some buffer for typical admission patterns
		regular = regularCount
		target = targetCount
		special = specialCount
		dedicated = dedicatedCount

		// If we only have regular applications, estimate total capacity
		if target == 0 && special == 0 && dedicated == 0 && regular > 0 {
			// Typical pattern: regular is about 70-80% of total
			estimatedTotal := int(float64(regular) / 0.75)
			regular = estimatedTotal
		}
	}

	return regular, target, special, dedicated
}
