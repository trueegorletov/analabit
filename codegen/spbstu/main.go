package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"sort"
	"strings"
)

// RegistryEntry represents a single entry in the SPbSTU registry JSON.
type RegistryEntry struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// RegistryFile represents the structure of SPbSTU registry JSON files.
type RegistryFile struct {
	CodeList []RegistryEntry `json:"code_list"`
}

// HeadingInfo stores information about a heading and its associated lists.
type HeadingInfo struct {
	Name                 string
	RegularListID        int
	TargetQuotaListIDs   []int
	DedicatedQuotaListID int
	SpecialQuotaListID   int
}

// extractHeadingName extracts clean heading names by removing registry code prefixes
// and optionally target organization suffixes in parentheses for target quota lists.
func extractHeadingName(title string, isTargetQuota bool) string {
	// Remove DD.DD.DD pattern prefix
	codePattern := regexp.MustCompile(`^\d{2}\.\d{2}\.\d{2}\s*`)
	cleaned := codePattern.ReplaceAllString(title, "")

	// Only remove organization suffix for target quota lists (for debugging purposes)
	if isTargetQuota {
		suffixPattern := regexp.MustCompile(`\s*\([^)]*\)\s*$`)
		cleaned = suffixPattern.ReplaceAllString(cleaned, "")
	}

	return strings.TrimSpace(cleaned)
}

// shouldIgnoreEntry checks if an entry should be ignored (e.g., (ИНО) entries in Regular&BVI lists)
func shouldIgnoreEntry(title string, registryType string) bool {
	// Ignore (ИНО) entries in Regular&BVI lists
	if registryType == "regular" && strings.Contains(title, "(ИНО)") {
		return true
	}
	return false
}

// readRegistryFile reads and parses a SPbSTU registry JSON file.
func readRegistryFile(filename string) (*RegistryFile, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	var registry RegistryFile
	if err := json.Unmarshal(data, &registry); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from %s: %w", filename, err)
	}

	return &registry, nil
}

func main() {
	// Define the registry files to process
	registryFiles := map[string]string{
		"regular":      "../../sample_data/spbstu/full_lists_registry_REGULARBVI.json",
		"special":      "../../sample_data/spbstu/full_lists_registry_SPECIALQUOTA.json",
		"dedicated":    "../../sample_data/spbstu/full_lists_registry_DEDICATEDQUOTA.json",
		"targetquotas": "../../sample_data/spbstu/full_lists_registry_TARGETQUOTAS.json",
	}

	// Read all registry files
	registries := make(map[string]*RegistryFile)
	for name, filename := range registryFiles {
		registry, err := readRegistryFile(filename)
		if err != nil {
			log.Fatalf("Error reading %s: %v", filename, err)
		}
		registries[name] = registry
		fmt.Printf("Loaded %s: %d entries\n", name, len(registry.CodeList))
	}

	// Group headings by name
	headings := make(map[string]*HeadingInfo)

	// Process Regular/BVI lists
	for _, entry := range registries["regular"].CodeList {
		// Skip (ИНО) entries in Regular&BVI lists
		if shouldIgnoreEntry(entry.Title, "regular") {
			continue
		}

		headingName := extractHeadingName(entry.Title, false)
		if headings[headingName] == nil {
			headings[headingName] = &HeadingInfo{
				Name:                 headingName,
				RegularListID:        -1,
				DedicatedQuotaListID: -1,
				SpecialQuotaListID:   -1,
			}
		}
		headings[headingName].RegularListID = entry.ID
	}

	// Process Special Quota lists
	for _, entry := range registries["special"].CodeList {
		headingName := extractHeadingName(entry.Title, false)
		if headings[headingName] == nil {
			headings[headingName] = &HeadingInfo{
				Name:                 headingName,
				RegularListID:        -1,
				DedicatedQuotaListID: -1,
				SpecialQuotaListID:   -1,
			}
		}
		headings[headingName].SpecialQuotaListID = entry.ID
	}

	// Process Dedicated Quota lists
	for _, entry := range registries["dedicated"].CodeList {
		headingName := extractHeadingName(entry.Title, false)
		if headings[headingName] == nil {
			headings[headingName] = &HeadingInfo{
				Name:                 headingName,
				RegularListID:        -1,
				DedicatedQuotaListID: -1,
				SpecialQuotaListID:   -1,
			}
		}
		headings[headingName].DedicatedQuotaListID = entry.ID
	}

	// Process Target Quota lists
	for _, entry := range registries["targetquotas"].CodeList {
		headingName := extractHeadingName(entry.Title, true) // Remove organization suffix for target quotas
		if headings[headingName] == nil {
			headings[headingName] = &HeadingInfo{
				Name:                 headingName,
				RegularListID:        -1,
				DedicatedQuotaListID: -1,
				SpecialQuotaListID:   -1,
			}
		}
		headings[headingName].TargetQuotaListIDs = append(headings[headingName].TargetQuotaListIDs, entry.ID)
	}

	// Sort headings for consistent output
	var sortedNames []string
	for name := range headings {
		sortedNames = append(sortedNames, name)
	}
	sort.Strings(sortedNames)

	// Generate output
	fmt.Println("\n// Generated SPbSTU HeadingSource definitions")
	fmt.Println("// Paste this into core/registry/spbstu/spbstu.go")
	fmt.Println()

	var anomalies []string
	headingsWithRegular := 0

	for _, name := range sortedNames {
		heading := headings[name]

		// Check for anomalies
		if heading.RegularListID == -1 {
			anomalies = append(anomalies, fmt.Sprintf("Missing Regular&BVI list for: %s", name))
		} else {
			headingsWithRegular++
		}

		// Generate HeadingSource definition
		fmt.Printf("&spbstu.HTTPHeadingSource{\n")
		fmt.Printf("\tPrettyName:           %q,\n", heading.Name)
		fmt.Printf("\tRegularListID:        %d,\n", heading.RegularListID)

		if len(heading.TargetQuotaListIDs) > 0 {
			fmt.Printf("\tTargetQuotaListIDs:   []int{")
			for i, id := range heading.TargetQuotaListIDs {
				if i > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%d", id)
			}
			fmt.Printf("},\n")
		} else {
			fmt.Printf("\tTargetQuotaListIDs:   []int{},\n")
		}

		fmt.Printf("\tDedicatedQuotaListID: %d,\n", heading.DedicatedQuotaListID)
		fmt.Printf("\tSpecialQuotaListID:   %d,\n", heading.SpecialQuotaListID)
		fmt.Printf("},\n\n")
	}

	// Print statistics
	fmt.Printf("// Statistics:\n")
	fmt.Printf("// Total headings: %d\n", len(headings))
	fmt.Printf("// Headings with Regular&BVI lists: %d\n", headingsWithRegular)
	fmt.Printf("// Regular&BVI entries: %d\n", len(registries["regular"].CodeList))
	fmt.Printf("// Special quota entries: %d\n", len(registries["special"].CodeList))
	fmt.Printf("// Dedicated quota entries: %d\n", len(registries["dedicated"].CodeList))
	fmt.Printf("// Target quota entries: %d\n", len(registries["targetquotas"].CodeList))

	// Print anomalies
	if len(anomalies) > 0 {
		fmt.Printf("\n// ANOMALIES DETECTED:\n")
		for _, anomaly := range anomalies {
			fmt.Printf("// %s\n", anomaly)
		}
	}
}
