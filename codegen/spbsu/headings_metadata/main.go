package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type NewName struct {
	ID                   int    `json:"id"`
	CompetitiveGroupID   int    `json:"competitive_group_id"`
	CompetitiveGroupName string `json:"competitive_group_name"`
	CreatedAt            string `json:"created_at"`
	UpdatedAt            string `json:"updated_at"`
}

// ListEntry represents a LEGACY format entry
type ListEntry struct {
	ID                     int      `json:"id"`
	CompetitiveName        string   `json:"competitive_name"`
	FormPaymentsID         int      `json:"form_payments_id"`
	FormEducationID        int      `json:"form_education_id"`
	ScenarioID             int      `json:"scenario_id"`
	ExternalID             string   `json:"external_id"`
	CompetitiveGroupPlaces int      `json:"competitive_group_places"`
	NewName                *NewName `json:"new_name"`
}

// NewFormatEntry represents a NEW format entry
type NewFormatEntry struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Code       int     `json:"code"`
	Visible    bool    `json:"visible"`
	ExternalID string  `json:"external_id"`
	CustomName *string `json:"custom_name"`
}

// NewFormatResponse represents the NEW format API response
type NewFormatResponse struct {
	Data []NewFormatEntry `json:"data"`
}

type ListIDCapacity struct {
	ID       int `json:"id"`
	Capacity int `json:"capacity"`
}

type HeadingInfo struct {
	RegularBVI     []ListIDCapacity `json:"regular_bvi,omitempty"`
	SpecialQuota   []ListIDCapacity `json:"special_quota,omitempty"`
	DedicatedQuota []ListIDCapacity `json:"dedicated_quota,omitempty"`
	TargetQuota    []ListIDCapacity `json:"target_quota,omitempty"`
}

// findFileUpwards searches for relPath by walking up from the current directory
func findFileUpwards(relPath string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		candidate := filepath.Join(cwd, relPath)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
		parent := filepath.Dir(cwd)
		if parent == cwd {
			break // reached root
		}
		cwd = parent
	}
	return "", fmt.Errorf("file %s not found in any parent directory", relPath)
}

func readEntriesFromFile(filePath string) ([]ListEntry, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var entries []ListEntry
	err = json.Unmarshal(data, &entries)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func readNewFormatFromFile(filePath string) ([]NewFormatEntry, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var response NewFormatResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// isInstitutionalCode checks if a string looks like an institutional code or location
// (short, no parentheses, often all caps)
func isInstitutionalCode(s string) bool {
	s = strings.TrimSpace(s)
	// Institutional codes are typically:
	// - Short (less than 20 characters)
	// - No parentheses (educational programs always have parentheses)
	// - Simple structure
	return len(s) < 20 && !strings.Contains(s, "(") && !strings.Contains(s, "дополнительной")
}

// extractEducationalProgram extracts the educational program name from a competitive name
// by taking the last part after splitting by semicolons, unless the last part is an institutional code
// For target quota entries, it extracts the part right after "Целевой прием"
func extractEducationalProgram(competitiveName string) string {
	parts := strings.Split(competitiveName, ";")
	if len(parts) == 0 {
		return ""
	}

	// Special handling for target quota entries
	if strings.Contains(competitiveName, "Целевой прием") {
		for i, part := range parts {
			if strings.TrimSpace(part) == "Целевой прием" && i+1 < len(parts) {
				return strings.TrimSpace(parts[i+1])
			}
		}
	}

	lastPart := strings.TrimSpace(parts[len(parts)-1])

	return lastPart
}

// determineQuotaType determines the quota type based on the competitive name and new_name field
func determineQuotaType(competitiveName string, newName *NewName) string {
	if newName != nil && strings.Contains(newName.CompetitiveGroupName, "Особая квота") {
		return "special"
	}
	if strings.Contains(competitiveName, "Отдельная квота") {
		return "dedicated"
	}
	if strings.Contains(competitiveName, "Целевой прием") {
		return "target"
	}
	return "regular"
}

// determineQuotaTypeFromNewFormat determines quota type from NEW format name field
func determineQuotaTypeFromNewFormat(name string) string {
	if strings.Contains(name, "Особая квота") {
		return "special"
	}
	if strings.Contains(name, "Отдельная квота") {
		return "dedicated"
	}
	if strings.Contains(name, "Целевой прием") {
		return "target"
	}
	return "regular"
}

func main() {
	// File paths for NEW format
	primaryNewRel := "codegen/spbsu/headings_metadata/primary-NEW.json"
	targetQuotaNewRel := "codegen/spbsu/headings_metadata/target_quota-NEW.json"

	// File paths for LEGACY format
	primaryLegacyRel := "codegen/spbsu/headings_metadata/primary-LEGACY.json"
	targetQuotaLegacyRel := "codegen/spbsu/headings_metadata/target_quota-LEGACY.json"

	// Find NEW format files
	primaryNewPath, err := findFileUpwards(primaryNewRel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find primary-NEW.json: %v\n", err)
		os.Exit(1)
	}
	targetQuotaNewPath, err := findFileUpwards(targetQuotaNewRel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find target_quota-NEW.json: %v\n", err)
		os.Exit(1)
	}

	// Find LEGACY format files
	primaryLegacyPath, err := findFileUpwards(primaryLegacyRel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find primary-LEGACY.json: %v\n", err)
		os.Exit(1)
	}
	targetQuotaLegacyPath, err := findFileUpwards(targetQuotaLegacyRel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find target_quota-LEGACY.json: %v\n", err)
		os.Exit(1)
	}

	// Read NEW format files
	primaryNew, err := readNewFormatFromFile(primaryNewPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read primary-NEW.json: %v\n", err)
		os.Exit(1)
	}
	targetQuotaNew, err := readNewFormatFromFile(targetQuotaNewPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read target_quota-NEW.json: %v\n", err)
		os.Exit(1)
	}

	// Read LEGACY format files
	primaryLegacy, err := readEntriesFromFile(primaryLegacyPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read primary-LEGACY.json: %v\n", err)
		os.Exit(1)
	}
	targetQuotaLegacy, err := readEntriesFromFile(targetQuotaLegacyPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read target_quota-LEGACY.json: %v\n", err)
		os.Exit(1)
	}

	// Build lookup maps from LEGACY data
	legacyByExternalID := make(map[string]ListEntry)
	legacyByProgram := make(map[string][]ListEntry)

	// Process primary LEGACY entries
	for _, entry := range primaryLegacy {
		legacyByExternalID[entry.ExternalID] = entry
		program := extractEducationalProgram(entry.CompetitiveName)
		if program != "" {
			legacyByProgram[program] = append(legacyByProgram[program], entry)
		}
	}

	// Process target quota LEGACY entries
	for _, entry := range targetQuotaLegacy {
		legacyByExternalID[entry.ExternalID] = entry
		program := extractEducationalProgram(entry.CompetitiveName)
		if program != "" {
			legacyByProgram[program] = append(legacyByProgram[program], entry)
		}
	}

	// Track processed entries and anomalies
	processedLegacyIDs := make(map[string]bool)
	var anomalies []string

	// Build headings from NEW format files using LEGACY capacity data
	headings := make(map[string]*HeadingInfo)

	// Process primary NEW entries (all except target quota)
	for _, newEntry := range primaryNew {
		program := extractEducationalProgram(newEntry.Name)
		if program == "" {
			anomalies = append(anomalies, fmt.Sprintf("NEW: Empty educational program for ID=%d, name=%q", newEntry.ID, newEntry.Name))
			continue
		}

		// Find matching LEGACY entry
		var legacyEntry *ListEntry
		if legacy, found := legacyByExternalID[newEntry.ExternalID]; found {
			legacyEntry = &legacy
			processedLegacyIDs[newEntry.ExternalID] = true
		} else {
			// Fallback: try to match by program name
			if legacyEntries, found := legacyByProgram[program]; found {
				for _, entry := range legacyEntries {
					if determineQuotaType(entry.CompetitiveName, entry.NewName) == determineQuotaTypeFromNewFormat(newEntry.Name) {
						legacyEntry = &entry
						processedLegacyIDs[entry.ExternalID] = true
						break
					}
				}
			}
		}

		if legacyEntry == nil {
			anomalies = append(anomalies, fmt.Sprintf("NEW: No matching LEGACY entry for ID=%d, external_id=%s, program=%q", newEntry.ID, newEntry.ExternalID, program))
			continue
		}

		// Initialize heading if not exists
		if _, ok := headings[program]; !ok {
			headings[program] = &HeadingInfo{}
		}

		// Determine quota type and add to appropriate list
		quotaType := determineQuotaTypeFromNewFormat(newEntry.Name)
		idcap := ListIDCapacity{ID: newEntry.ID, Capacity: legacyEntry.CompetitiveGroupPlaces}

		switch quotaType {
		case "special":
			headings[program].SpecialQuota = append(headings[program].SpecialQuota, idcap)
		case "dedicated":
			headings[program].DedicatedQuota = append(headings[program].DedicatedQuota, idcap)
		default:
			headings[program].RegularBVI = append(headings[program].RegularBVI, idcap)
		}
	}

	// Process target quota NEW entries (always target quota)
	for _, newEntry := range targetQuotaNew {
		program := extractEducationalProgram(newEntry.Name)
		if program == "" {
			anomalies = append(anomalies, fmt.Sprintf("NEW: Empty educational program for target quota ID=%d, name=%q", newEntry.ID, newEntry.Name))
			continue
		}

		// Find matching LEGACY entry
		var legacyEntry *ListEntry
		if legacy, found := legacyByExternalID[newEntry.ExternalID]; found {
			legacyEntry = &legacy
			processedLegacyIDs[newEntry.ExternalID] = true
		} else {
			// Fallback: try to match by program name
			if legacyEntries, found := legacyByProgram[program]; found {
				for _, entry := range legacyEntries {
					if determineQuotaType(entry.CompetitiveName, entry.NewName) == "target" ||
						strings.Contains(entry.CompetitiveName, "Целевой") {
						legacyEntry = &entry
						processedLegacyIDs[entry.ExternalID] = true
						break
					}
				}
			}
		}

		if legacyEntry == nil {
			anomalies = append(anomalies, fmt.Sprintf("NEW: No matching LEGACY entry for target quota ID=%d, external_id=%s, program=%q", newEntry.ID, newEntry.ExternalID, program))
			continue
		}

		// Initialize heading if not exists
		if _, ok := headings[program]; !ok {
			headings[program] = &HeadingInfo{}
		}

		idcap := ListIDCapacity{ID: newEntry.ID, Capacity: legacyEntry.CompetitiveGroupPlaces}
		headings[program].TargetQuota = append(headings[program].TargetQuota, idcap)
	}

	// Check for unprocessed LEGACY entries
	for _, entry := range primaryLegacy {
		if !processedLegacyIDs[entry.ExternalID] {
			program := extractEducationalProgram(entry.CompetitiveName)
			anomalies = append(anomalies, fmt.Sprintf("LEGACY: Unmatched primary entry ID=%d, external_id=%s, program=%q", entry.ID, entry.ExternalID, program))
		}
	}
	for _, entry := range targetQuotaLegacy {
		if !processedLegacyIDs[entry.ExternalID] {
			program := extractEducationalProgram(entry.CompetitiveName)
			anomalies = append(anomalies, fmt.Sprintf("LEGACY: Unmatched target quota entry ID=%d, external_id=%s, program=%q", entry.ID, entry.ExternalID, program))
		}
	}

	// Output the result as Go struct literals for sourcesList
	problematic := []string{}
	for heading, info := range headings {
		// Find first list ID for each type, or -1 if none
		regularID := -1
		if len(info.RegularBVI) > 0 {
			regularID = info.RegularBVI[0].ID
		}
		dedicatedID := -1
		if len(info.DedicatedQuota) > 0 {
			dedicatedID = info.DedicatedQuota[0].ID
		}
		specialID := -1
		if len(info.SpecialQuota) > 0 {
			specialID = info.SpecialQuota[0].ID
		}

		// Collect all target quota IDs (can be multiple)
		var targetIDs []int
		for _, c := range info.TargetQuota {
			targetIDs = append(targetIDs, c.ID)
		}

		// Sum capacities for each type
		capRegular := 0
		for _, c := range info.RegularBVI {
			capRegular += c.Capacity
		}
		capTarget := 0
		for _, c := range info.TargetQuota {
			capTarget += c.Capacity
		}
		capDedicated := 0
		for _, c := range info.DedicatedQuota {
			capDedicated += c.Capacity
		}
		capSpecial := 0
		for _, c := range info.SpecialQuota {
			capSpecial += c.Capacity
		}

		var out strings.Builder
		fmt.Fprintf(&out, "// %s\n", heading)
		fmt.Fprintf(&out, "&spbsu.HttpHeadingSource{\n")
		fmt.Fprintf(&out, "\tPrettyName: \"%s\",\n", heading)
		fmt.Fprintf(&out, "\tRegularListID: %d,\n", regularID)

		// Output target quota IDs as a slice
		if len(targetIDs) == 0 {
			fmt.Fprintf(&out, "\tTargetQuotaListIDs: []int{},\n")
		} else if len(targetIDs) == 1 {
			fmt.Fprintf(&out, "\tTargetQuotaListIDs: []int{%d},\n", targetIDs[0])
		} else {
			fmt.Fprintf(&out, "\tTargetQuotaListIDs: []int{")
			for i, id := range targetIDs {
				if i > 0 {
					fmt.Fprintf(&out, ", ")
				}
				fmt.Fprintf(&out, "%d", id)
			}
			fmt.Fprintf(&out, "},\n")
		}

		fmt.Fprintf(&out, "\tDedicatedQuotaListID: %d,\n", dedicatedID)
		fmt.Fprintf(&out, "\tSpecialQuotaListID: %d,\n", specialID)
		fmt.Fprintf(&out, "\tCapacities: core.Capacities{\n")
		fmt.Fprintf(&out, "\t\tRegular: %d,\n", capRegular)
		fmt.Fprintf(&out, "\t\tTargetQuota: %d,\n", capTarget)
		fmt.Fprintf(&out, "\t\tDedicatedQuota: %d,\n", capDedicated)
		fmt.Fprintf(&out, "\t\tSpecialQuota: %d,\n", capSpecial)
		fmt.Fprintf(&out, "\t},\n")
		fmt.Fprintf(&out, "},\n\n")

		if regularID == -1 || capRegular == 0 {
			problematic = append(problematic, out.String())
		} else {
			fmt.Print(out.String())
		}
	}

	if len(problematic) > 0 {
		fmt.Println("// --- Headings with missing or zero regular capacity ---")
		for _, entry := range problematic {
			for _, line := range strings.Split(entry, "\n") {
				if line != "" {
					fmt.Printf("// %s\n", line)
				}
			}
			fmt.Println()
		}
	}

	// Output anomalies
	if len(anomalies) > 0 {
		fmt.Println("// --- ANOMALIES REPORT ---")
		for _, anomaly := range anomalies {
			fmt.Printf("// %s\n", anomaly)
		}
	}
}
