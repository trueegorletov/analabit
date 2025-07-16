package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RegistryEntry represents a single entry in the registry JSON
type RegistryEntry struct {
	Title string `json:"title"`
	File  string `json:"file"`
}



// ProgramData holds categorized URLs for a program
type ProgramData struct {
	TargetQuotaURLs    []string
	RegularURL         string
	SpecialQuotaURL    string
	DedicatedQuotaURL  string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <registry-json-file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s ../../sample_data/rsmu/p2_560.json\n", os.Args[0])
		os.Exit(1)
	}

	registryFile := os.Args[1]

	// Read the registry JSON file
	registryContent, err := os.ReadFile(registryFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading registry file: %v\n", err)
		os.Exit(1)
	}

	// Parse registry
	var registryEntries []RegistryEntry
	if err := json.Unmarshal(registryContent, &registryEntries); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing registry JSON: %v\n", err)
		os.Exit(1)
	}

	// Get the base directory for individual list files
	baseDir := filepath.Dir(registryFile)

	// Process registry entries and categorize by program
	programs := processRegistryEntries(registryEntries, baseDir)

	// Generate HTTPHeadingSource entries
	generateHTTPHeadingSources(programs, registryFile)
}

func processRegistryEntries(registryEntries []RegistryEntry, baseDir string) map[string]*ProgramData {
	programs := make(map[string]*ProgramData)

	for _, entry := range registryEntries {
		// Skip contract entries as they're not budget-related
		if strings.Contains(strings.ToLower(entry.Title), "контракт") {
			continue
		}

		// Extract and clean program name
		programName := extractProgramName(entry)
		cleanProgramName := cleanProgramName(programName)

		// Initialize program data if not exists
		if programs[cleanProgramName] == nil {
			programs[cleanProgramName] = &ProgramData{}
		}

		// Convert file path to URL
		url := convertToURL(entry.File)

		// Categorize by competition type
		competitionType := detectCompetitionType(entry.Title)
		switch competitionType {
		case "target":
			programs[cleanProgramName].TargetQuotaURLs = append(programs[cleanProgramName].TargetQuotaURLs, url)
		case "regular":
			programs[cleanProgramName].RegularURL = url
		case "special":
			programs[cleanProgramName].SpecialQuotaURL = url
		case "dedicated":
			programs[cleanProgramName].DedicatedQuotaURL = url
		}
	}

	return programs
}

func generateHTTPHeadingSources(programs map[string]*ProgramData, registryFile string) {
	fmt.Println("// Generated RSMU HTTPHeadingSource entries")
	fmt.Println("// Based on registry:", registryFile)
	fmt.Println()

	for programName, data := range programs {
		// Generate HTTPHeadingSource entry
		fmt.Printf("&rsmu.HTTPHeadingSource{\n")
		fmt.Printf("\tProgramName: \"%s\",\n", programName)

		// Add TargetQuotaListURLs if any
		if len(data.TargetQuotaURLs) > 0 {
			fmt.Printf("\tTargetQuotaListURLs: []string{\n")
			for _, url := range data.TargetQuotaURLs {
				fmt.Printf("\t\t\"%s\",\n", url)
			}
			fmt.Printf("\t},\n")
		}

		// Add other URLs if they exist
		if data.RegularURL != "" {
			fmt.Printf("\tRegularListURL: \"%s\",\n", data.RegularURL)
		}
		if data.SpecialQuotaURL != "" {
			fmt.Printf("\tSpecialQuotaListURL: \"%s\",\n", data.SpecialQuotaURL)
		}
		if data.DedicatedQuotaURL != "" {
			fmt.Printf("\tDedicatedQuotaListURL: \"%s\",\n", data.DedicatedQuotaURL)
		}

		fmt.Printf("},\n")
		fmt.Println()
	}

	fmt.Printf("// Total RSMU headings generated: %d\n", len(programs))
}

func extractProgramName(entry RegistryEntry) string {
	return entry.Title
}

func cleanProgramName(programName string) string {
	cleaned := strings.TrimSpace(programName)
	
	// Competition type markers to split on
	competitionMarkers := []string{
		" Целевая квота",
		" Общий конкурс",
		" Особая квота",
		" Отдельная квота",
		" Контракт",
	}
	
	// Find the first competition marker and extract everything before it
	programPart := cleaned
	for _, marker := range competitionMarkers {
		if idx := strings.Index(cleaned, marker); idx != -1 {
			programPart = strings.TrimSpace(cleaned[:idx])
			break
		}
	}
	
	// Look for the FIRST parentheses pair in the program part
	start := strings.Index(programPart, "(")
	if start != -1 {
		// Find the matching closing parenthesis
		end := strings.Index(programPart[start:], ")")
		if end != -1 {
			// Extract text within the first parentheses pair
			return strings.TrimSpace(programPart[start+1 : start+end])
		}
	}
	
	// If no parentheses found, return the cleaned program part
	return strings.TrimSpace(programPart)
}

func detectCompetitionType(title string) string {
	titleLower := strings.ToLower(title)
	
	if strings.Contains(titleLower, "целевая квота") {
		return "target"
	}
	if strings.Contains(titleLower, "общий конкурс") {
		return "regular"
	}
	if strings.Contains(titleLower, "особая квота") {
		return "special"
	}
	if strings.Contains(titleLower, "отдельная квота") {
		return "dedicated"
	}
	// Note: contract competitions are excluded
	
	return "regular" // Default to regular
}

func convertToURL(filePath string) string {
	// Convert local file path to a URL format
	// This is a placeholder - in real usage, this would be the actual URL
	// For now, we'll use a placeholder URL structure
	return fmt.Sprintf("https://submitted.rsmu.ru/data/%s", filePath)
}