package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"
)

// Competition represents a single competition within a heading
type Competition struct {
	CompTypeID int      `json:"comp_type_id"`
	CompType   string   `json:"comp_type"`
	Plan       int      `json:"plan"`
	AppCount   int      `json:"app_count"`
	CompIDs    []string `json:"comp_ids"`
}

// RegistryEntry represents a single heading entry in the MIREA registry JSON
type RegistryEntry struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Plan         int           `json:"plan"`
	AppCount     int           `json:"app_count"`
	Competitions []Competition `json:"competitions"`
	// Skip foreign_competitions entirely
}

// HeadingInfo stores information about a heading and its associated lists
type HeadingInfo struct {
	Name                  string
	RegularListIDs        []string
	BVIListIDs            []string
	TargetQuotaListIDs    []string
	DedicatedQuotaListIDs []string
	SpecialQuotaListIDs   []string
}

// extractHeadingName extracts clean heading names by cutting at the first slash
func extractHeadingName(title string) string {
	// Cut at first slash as specified in the format reference
	if slashIndex := strings.LastIndex(title, "/"); slashIndex != -1 {
		return strings.TrimSpace(title[:slashIndex])
	}
	return strings.TrimSpace(title)
}

// mapCompetitionType maps comp_type_id to competition type name
func mapCompetitionType(compTypeID int) string {
	switch compTypeID {
	case 4:
		return "Regular"
	case 1:
		return "BVI"
	case 2:
		return "SpecialQuota"
	case 3:
		return "TargetQuota"
	case 7:
		return "DedicatedQuota"
	case 6:
		return "NonBudgetary" // Will be ignored
	default:
		return "Unknown"
	}
}

// shouldIgnoreCompetition checks if a competition should be ignored
func shouldIgnoreCompetition(compTypeID int) bool {
	// Ignore non-budgetary competitions (comp_type_id 6)
	return compTypeID == 6
}

// readRegistryFile reads and parses the MIREA registry JSON file
func readRegistryFile(filename string) ([]RegistryEntry, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	var registry []RegistryEntry
	if err := json.Unmarshal(data, &registry); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from %s: %w", filename, err)
	}

	return registry, nil
}

func main() {
	// Read the registry file
	registryFile := "../../sample_data/mirea/lists_registry_ALL.json"
	registry, err := readRegistryFile(registryFile)
	if err != nil {
		log.Fatalf("Error reading registry: %v", err)
	}

	fmt.Printf("Loaded MIREA registry: %d entries\n", len(registry))

	// Group headings by name
	headings := make(map[string]*HeadingInfo)
	totalCompetitions := 0
	ignoredCompetitions := 0
	competitionTypeStats := make(map[string]int)

	// Process each heading entry
	for _, entry := range registry {
		headingName := extractHeadingName(entry.Title)

		if headings[headingName] == nil {
			headings[headingName] = &HeadingInfo{
				Name:                  headingName,
				RegularListIDs:        []string{},
				BVIListIDs:            []string{},
				TargetQuotaListIDs:    []string{},
				DedicatedQuotaListIDs: []string{},
				SpecialQuotaListIDs:   []string{},
			}
		}

		// Process competitions for this heading
		for _, competition := range entry.Competitions {
			totalCompetitions++

			if shouldIgnoreCompetition(competition.CompTypeID) {
				ignoredCompetitions++
				continue
			}

			compType := mapCompetitionType(competition.CompTypeID)
			competitionTypeStats[compType]++

			// Add competition IDs to appropriate slices
			switch competition.CompTypeID {
			case 4: // Regular
				headings[headingName].RegularListIDs = append(headings[headingName].RegularListIDs, competition.CompIDs...)
			case 1: // BVI
				headings[headingName].BVIListIDs = append(headings[headingName].BVIListIDs, competition.CompIDs...)
			case 2: // SpecialQuota
				headings[headingName].SpecialQuotaListIDs = append(headings[headingName].SpecialQuotaListIDs, competition.CompIDs...)
			case 3: // TargetQuota
				headings[headingName].TargetQuotaListIDs = append(headings[headingName].TargetQuotaListIDs, competition.CompIDs...)
			case 7: // DedicatedQuota
				headings[headingName].DedicatedQuotaListIDs = append(headings[headingName].DedicatedQuotaListIDs, competition.CompIDs...)
			}
		}
	}

	// Sort headings for consistent output
	var sortedNames []string
	for name := range headings {
		sortedNames = append(sortedNames, name)
	}
	sort.Strings(sortedNames)

	// Generate output
	fmt.Println("\n// Generated MIREA HeadingSource definitions")
	fmt.Println("// Paste this into core/registry/mirea/mirea.go")
	fmt.Println()

	var anomalies []string
	headingsWithRegular := 0
	headingsWithBVI := 0
	headingsWithTargetQuota := 0

	for _, name := range sortedNames {
		heading := headings[name]

		// Check for anomalies and count statistics
		if len(heading.RegularListIDs) == 0 {
			anomalies = append(anomalies, fmt.Sprintf("Missing Regular lists for: %s", name))
		} else {
			headingsWithRegular++
		}

		if len(heading.BVIListIDs) > 0 {
			headingsWithBVI++
		}

		if len(heading.TargetQuotaListIDs) > 0 {
			headingsWithTargetQuota++
		}

		// Generate HeadingSource definition
		fmt.Printf("// %s\n", name)
		fmt.Printf("&mirea.HTTPHeadingSource{\n")
		fmt.Printf("\tPrettyName: %q,\n", name)

		// Regular lists
		if len(heading.RegularListIDs) > 0 {
			fmt.Printf("\tRegularListIDs: []string{")
			for i, id := range heading.RegularListIDs {
				if i > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%q", id)
			}
			fmt.Printf("},\n")
		} else {
			fmt.Printf("\tRegularListIDs: []string{},\n")
		}

		// BVI lists
		if len(heading.BVIListIDs) > 0 {
			fmt.Printf("\tBVIListIDs: []string{")
			for i, id := range heading.BVIListIDs {
				if i > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%q", id)
			}
			fmt.Printf("},\n")
		} else {
			fmt.Printf("\tBVIListIDs: []string{},\n")
		}

		// TargetQuota lists
		if len(heading.TargetQuotaListIDs) > 0 {
			fmt.Printf("\tTargetQuotaListIDs: []string{")
			for i, id := range heading.TargetQuotaListIDs {
				if i > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%q", id)
			}
			fmt.Printf("},\n")
		} else {
			fmt.Printf("\tTargetQuotaListIDs: []string{},\n")
		}

		// DedicatedQuota lists
		if len(heading.DedicatedQuotaListIDs) > 0 {
			fmt.Printf("\tDedicatedQuotaListIDs: []string{")
			for i, id := range heading.DedicatedQuotaListIDs {
				if i > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%q", id)
			}
			fmt.Printf("},\n")
		} else {
			fmt.Printf("\tDedicatedQuotaListIDs: []string{},\n")
		}

		// SpecialQuota lists
		if len(heading.SpecialQuotaListIDs) > 0 {
			fmt.Printf("\tSpecialQuotaListIDs: []string{")
			for i, id := range heading.SpecialQuotaListIDs {
				if i > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%q", id)
			}
			fmt.Printf("},\n")
		} else {
			fmt.Printf("\tSpecialQuotaListIDs: []string{},\n")
		}

		fmt.Printf("},\n\n")
	}

	// Print statistics
	fmt.Printf("// Statistics:\n")
	fmt.Printf("// Total headings: %d\n", len(headings))
	fmt.Printf("// Headings with Regular lists: %d\n", headingsWithRegular)
	fmt.Printf("// Headings with BVI lists: %d\n", headingsWithBVI)
	fmt.Printf("// Headings with TargetQuota lists: %d\n", headingsWithTargetQuota)
	fmt.Printf("// Total competitions processed: %d\n", totalCompetitions)
	fmt.Printf("// Ignored non-budgetary competitions: %d\n", ignoredCompetitions)

	// Print competition type statistics
	fmt.Printf("// Competition type breakdown:\n")
	for compType, count := range competitionTypeStats {
		fmt.Printf("// %s: %d\n", compType, count)
	}

	// Print anomalies
	if len(anomalies) > 0 {
		fmt.Printf("\n// ANOMALIES DETECTED:\n")
		for _, anomaly := range anomalies {
			fmt.Printf("// %s\n", anomaly)
		}
	}
}
