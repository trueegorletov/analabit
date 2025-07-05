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

func main() {
	primaryRel := "codegen/spbsu/headings_metadata/primary.json"
	targetQuotaRel := "codegen/spbsu/headings_metadata/target_quota.json"

	primaryPath, err := findFileUpwards(primaryRel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find primary.json: %v\n", err)
		os.Exit(1)
	}
	targetQuotaPath, err := findFileUpwards(targetQuotaRel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find target_quota.json: %v\n", err)
		os.Exit(1)
	}

	primaryEntries, err := readEntriesFromFile(primaryPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read primary.json: %v\n", err)
		os.Exit(1)
	}
	targetQuotaEntries, err := readEntriesFromFile(targetQuotaPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read target_quota.json: %v\n", err)
		os.Exit(1)
	}

	headings := make(map[string]*HeadingInfo)
	anomalies := []string{}

	// First, process primary.json entries (all except TargetQuota)
	for _, entry := range primaryEntries {
		parts := strings.Split(entry.CompetitiveName, ";")
		if len(parts) == 0 {
			anomalies = append(anomalies, fmt.Sprintf("No semicolon in competitive_name: %q (id=%d)", entry.CompetitiveName, entry.ID))
			continue
		}
		heading := strings.TrimSpace(parts[len(parts)-1])
		if heading == "" {
			anomalies = append(anomalies, fmt.Sprintf("Empty heading for id=%d", entry.ID))
			continue
		}
		if _, ok := headings[heading]; !ok {
			headings[heading] = &HeadingInfo{}
		}
		typeAdded := false
		idcap := ListIDCapacity{ID: entry.ID, Capacity: entry.CompetitiveGroupPlaces}
		if entry.NewName != nil && strings.Contains(entry.NewName.CompetitiveGroupName, "Особая квота") {
			headings[heading].SpecialQuota = append(headings[heading].SpecialQuota, idcap)
			typeAdded = true
		}
		if strings.Contains(entry.CompetitiveName, "Отдельная квота") {
			headings[heading].DedicatedQuota = append(headings[heading].DedicatedQuota, idcap)
			typeAdded = true
		}
		// Never assign TargetQuota here
		if !typeAdded {
			headings[heading].RegularBVI = append(headings[heading].RegularBVI, idcap)
		}
	}

	// Now process target_quota.json entries (always TargetQuota)
	for _, entry := range targetQuotaEntries {
		parts := strings.Split(entry.CompetitiveName, ";")
		if len(parts) == 0 {
			anomalies = append(anomalies, fmt.Sprintf("No semicolon in competitive_name: %q (id=%d)", entry.CompetitiveName, entry.ID))
			continue
		}
		heading := strings.TrimSpace(parts[len(parts)-1])
		if heading == "" {
			anomalies = append(anomalies, fmt.Sprintf("Empty heading for id=%d", entry.ID))
			continue
		}
		if _, ok := headings[heading]; !ok {
			headings[heading] = &HeadingInfo{}
		}
		idcap := ListIDCapacity{ID: entry.ID, Capacity: entry.CompetitiveGroupPlaces}
		headings[heading].TargetQuota = append(headings[heading].TargetQuota, idcap)
	}

	// Output the result as Go struct literals for sourcesList
	problematic := []string{}
	for heading, info := range headings {
		// Find first list ID for each type, or -1 if none
		regularID := -1
		if len(info.RegularBVI) > 0 {
			regularID = info.RegularBVI[0].ID
		}
		targetID := -1
		if len(info.TargetQuota) > 0 {
			targetID = info.TargetQuota[0].ID
		}
		dedicatedID := -1
		if len(info.DedicatedQuota) > 0 {
			dedicatedID = info.DedicatedQuota[0].ID
		}
		specialID := -1
		if len(info.SpecialQuota) > 0 {
			specialID = info.SpecialQuota[0].ID
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
		fmt.Fprintf(&out, "\tTargetQuotaListID: %d,\n", targetID)
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

	// Log anomalies
	if len(anomalies) > 0 {
		fmt.Fprintln(os.Stderr, "\nAnomalies detected:")
		for _, a := range anomalies {
			fmt.Fprintln(os.Stderr, a)
		}
	}
}
