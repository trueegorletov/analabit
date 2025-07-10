package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type MiptProgram struct {
	Name       string
	Capacities MiptCapacities
}

type MiptCapacities struct {
	Total          int
	Regular        int
	BVI            int
	TargetQuota    int
	DedicatedQuota int
	SpecialQuota   int
}

type MiptURLs struct {
	RegularBVI      string
	SpecialQuota    string
	DedicatedQuota  string
	TargetQuotaURLs []string
	Contract        string
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <capacities-html-file> <registry-html-file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s ../../sample_data/mipt/mipt-CapacitiesByHeadingName-capacities-page-UNCOMPRESSED.html ../../sample_data/mipt/mipt-lists-registry_UNCOMPRESSED.html\n", os.Args[0])
		os.Exit(1)
	}

	capacitiesFile := os.Args[1]
	registryFile := os.Args[2]

	// Read the capacities HTML file
	capacitiesContent, err := ioutil.ReadFile(capacitiesFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading capacities file: %v\n", err)
		os.Exit(1)
	}

	// Read the registry HTML file
	registryContent, err := ioutil.ReadFile(registryFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading registry file: %v\n", err)
		os.Exit(1)
	}

	// Parse educational programs from capacities file
	programs, err := parsePrograms(string(capacitiesContent))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing programs: %v\n", err)
		os.Exit(1)
	}

	// Parse URLs from the registry file using table structure
	urlMap, err := parseRegistryURLsTableStructured(string(registryContent))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing registry URLs: %v\n", err)
		os.Exit(1)
	}

	// Generate HTTPHeadingSource entries
	generateHTTPHeadingSources(programs, urlMap)
}

func parsePrograms(content string) ([]MiptProgram, error) {
	fmt.Println("// Parsing educational programs from capacities file...")

	programCapacities := make(map[string]MiptCapacities)

	// Split content into table rows
	rows := strings.Split(content, "<tr")
	fmt.Printf("// Found %d table row sections\n", len(rows))

	for _, row := range rows {
		// Skip header and empty rows
		if len(row) < 100 {
			continue
		}

		// Look for educational program names in the row
		// Pattern for program names in table cells
		programPattern := regexp.MustCompile(`<p class="MsoNormal"[^>]*><span[^>]*><font face="Arial">([А-Яа-я\s,\-\(\)«»]+)<o:p></o:p></font></span></p>`)
		programMatches := programPattern.FindAllStringSubmatch(row, -1)

		if len(programMatches) == 0 {
			continue
		}

		for _, programMatch := range programMatches {
			programName := strings.TrimSpace(programMatch[1])

			// Filter out direction codes, headers, organization names, and irrelevant text
			if strings.Contains(programName, ".") || // Skip codes like "03.03.01"
				len(programName) < 15 || // Skip very short strings
				strings.Contains(strings.ToLower(programName), "направление") ||
				strings.Contains(strings.ToLower(programName), "уровень") ||
				strings.Contains(strings.ToLower(programName), "подготовка") ||
				strings.Contains(strings.ToLower(programName), "программа") ||
				strings.Contains(strings.ToLower(programName), "специальность") ||
				strings.Contains(strings.ToLower(programName), "профиль") ||
				strings.HasPrefix(programName, "0") || // Skip numeric codes
				strings.Contains(programName, "ФГОС") ||
				strings.Contains(programName, "г.") || // Skip years like "2023 г."
				strings.Contains(strings.ToLower(programName), "акционерное") ||
				strings.Contains(strings.ToLower(programName), "общество") ||
				strings.Contains(strings.ToLower(programName), "федеральное") ||
				strings.Contains(strings.ToLower(programName), "государственное") ||
				strings.Contains(strings.ToLower(programName), "автономное") ||
				strings.Contains(strings.ToLower(programName), "учреждение") ||
				strings.Contains(strings.ToLower(programName), "предприятие") ||
				strings.Contains(strings.ToLower(programName), "институт") ||
				strings.Contains(strings.ToLower(programName), "завод") ||
				strings.Contains(strings.ToLower(programName), "центр") ||
				strings.Contains(strings.ToLower(programName), "организации") ||
				strings.Contains(strings.ToLower(programName), "иные") ||
				strings.Contains(programName, "АО") ||
				strings.Contains(programName, "ПАО") ||
				strings.Contains(programName, "ФГУП") ||
				strings.Contains(programName, "ООО") {
				continue
			}

			// Extract capacity numbers from this row
			// Look for centered numbers in table cells
			numberPattern := regexp.MustCompile(`<p class="MsoNormal" align="center"[^>]*><span[^>]*><font[^>]*>(\d+)<`)
			numberMatches := numberPattern.FindAllStringSubmatch(row, -1)

			if len(numberMatches) >= 4 {
				totalStr := numberMatches[0][1]
				regularStr := numberMatches[1][1]
				specialStr := numberMatches[2][1]
				dedicatedStr := numberMatches[3][1]

				total, _ := strconv.Atoi(totalStr)
				regular, _ := strconv.Atoi(regularStr)
				special, _ := strconv.Atoi(specialStr)
				dedicated, _ := strconv.Atoi(dedicatedStr)

				// Validate reasonable capacity numbers
				if total > 0 && total < 1000 && regular >= 0 && special >= 0 && dedicated >= 0 {
					programCapacities[programName] = MiptCapacities{
						Total:          total,
						Regular:        regular,
						SpecialQuota:   special,
						DedicatedQuota: dedicated,
						BVI:            0, // BVI will be calculated or set separately
						TargetQuota:    0, // Target quota will be calculated
					}

					fmt.Printf("// Found program: %s - Total: %d, Regular: %d, Special: %d, Dedicated: %d\n",
						programName, total, regular, special, dedicated)
				}
			}
		}
	}

	// Convert map to slice
	var programs []MiptProgram
	for name, capacities := range programCapacities {
		programs = append(programs, MiptProgram{
			Name:       name,
			Capacities: capacities,
		})
	}

	fmt.Printf("// Extracted capacities for %d programs\n", len(programs))
	return programs, nil
}

func parseRegistryURLsTableStructured(content string) (map[string]MiptURLs, error) {
	fmt.Println("// Parsing URLs from registry file using table structure...")

	urlMap := make(map[string]MiptURLs)

	// Split by table rows to process each row individually
	rows := strings.Split(content, "<tr")
	fmt.Printf("// Processing %d table rows\n", len(rows))

	for rowIndex, row := range rows {
		// Skip very short rows
		if len(row) < 200 {
			continue
		}

		// Extract the first table cell to check for program name
		// Look for the first <td> element in the row
		tdStart := strings.Index(row, "<td")
		if tdStart == -1 {
			continue
		}

		// Find the end of the opening <td> tag
		tdContentStart := strings.Index(row[tdStart:], ">")
		if tdContentStart == -1 {
			continue
		}
		tdContentStart += tdStart + 1

		// Find the closing </td> tag
		tdEnd := strings.Index(row[tdContentStart:], "</td>")
		if tdEnd == -1 {
			continue
		}

		firstCellContent := row[tdContentStart : tdContentStart+tdEnd]

		// Look for educational program names in the first cell
		// Remove HTML tags from the cell content
		cleanContent := regexp.MustCompile(`<[^>]*>`).ReplaceAllString(firstCellContent, "")
		cleanContent = strings.TrimSpace(cleanContent)

		// Handle HTML entities
		cleanContent = strings.ReplaceAll(cleanContent, "&nbsp;", " ")
		cleanContent = strings.ReplaceAll(cleanContent, "&amp;", "&")
		cleanContent = strings.ReplaceAll(cleanContent, "&lt;", "<")
		cleanContent = strings.ReplaceAll(cleanContent, "&gt;", ">")
		cleanContent = strings.ReplaceAll(cleanContent, "&quot;", "\"")

		// Normalize whitespace
		cleanContent = regexp.MustCompile(`\s+`).ReplaceAllString(cleanContent, " ")
		cleanContent = strings.TrimSpace(cleanContent)

		// Check if this looks like an educational program name
		if !isEducationalProgram(cleanContent) {
			continue
		}

		programName := cleanContent
		fmt.Printf("// Found program in row %d: '%s'\n", rowIndex, programName)

		// Extract all URLs from this same row
		urls := MiptURLs{}
		urlPattern := regexp.MustCompile(`href="([^"]+)"`)
		urlMatches := urlPattern.FindAllStringSubmatch(row, -1)

		for _, urlMatch := range urlMatches {
			if len(urlMatch) < 2 {
				continue
			}

			url := urlMatch[1]

			// Only process MIPT application URLs
			if !strings.Contains(url, "applications_v2/") {
				continue
			}

			// Decode base64 URLs to classify them
			pathStart := strings.Index(url, "applications_v2/") + len("applications_v2/")
			if pathStart < len(url) {
				encodedPath := url[pathStart:]
				if decodedBytes, err := base64.StdEncoding.DecodeString(encodedPath); err == nil {
					decodedPath := string(decodedBytes)

					// Classify URL by type based on decoded path
					if strings.Contains(decodedPath, "_Byudzhet_Na") || strings.Contains(decodedPath, "_Na obshchikh") {
						urls.RegularBVI = url
						fmt.Printf("//   + RegularBVI URL\n")
					} else if strings.Contains(decodedPath, "_Imeyushchie osoboe") || strings.Contains(decodedPath, "Osoboe pravo") {
						urls.SpecialQuota = url
						fmt.Printf("//   + SpecialQuota URL\n")
					} else if strings.Contains(decodedPath, "_Otdelnaya kvota") {
						urls.DedicatedQuota = url
						fmt.Printf("//   + DedicatedQuota URL\n")
					} else if strings.Contains(decodedPath, "_Tselevoe") {
						urls.TargetQuotaURLs = append(urls.TargetQuotaURLs, url)
						fmt.Printf("//   + TargetQuota URL\n")
					}
				}
			}
		}

		// Only include programs that have at least a budget URL
		if urls.RegularBVI != "" {
			urlMap[programName] = urls
			fmt.Printf("// Stored program with %d total URLs\n", 1+len(urls.TargetQuotaURLs)+boolToInt(urls.SpecialQuota != "")+boolToInt(urls.DedicatedQuota != ""))
		}
	}

	fmt.Printf("// Extracted URLs for %d programs with budget places\n", len(urlMap))
	return urlMap, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Helper function to identify if a string is likely an educational program name
func isEducationalProgram(candidate string) bool {
	candidate = strings.TrimSpace(candidate)

	// Must be long enough
	if len(candidate) < 10 {
		return false
	}

	// Must contain Cyrillic characters
	if !regexp.MustCompile(`[А-Яа-я]`).MatchString(candidate) {
		return false
	}

	// Skip obvious non-program text
	skipPatterns := []string{
		"программы", "реализуемые", "русском", "языке", "английском",
		"иностранных", "граждан", "платные", "места", "контракт",
		"бюджет", "квота", "целевая", "особая", "отдельная",
		"конкурсная", "группа", "основа", "направление",
		"уровень", "подготовка", "специальность", "профиль",
		"физтех", "школа", "высшая", "им.", "имени",
		"организации", "иные", "граждан", "только",
		"фгос", "г.", "года", "курс",
	}

	candidateLower := strings.ToLower(candidate)
	for _, skip := range skipPatterns {
		if strings.Contains(candidateLower, skip) {
			return false
		}
	}

	// Skip if it looks like a direction code
	if regexp.MustCompile(`^\d+\.\d+\.\d+`).MatchString(candidate) {
		return false
	}

	// Skip organization names
	orgPatterns := []string{
		"акционерное", "общество", "федеральное", "государственное",
		"автономное", "учреждение", "предприятие", "институт",
		"завод", "центр", "пао", "ао", "фгуп", "ооо", "нии",
		"цнии", "нпо", "нпп", "мцст", "циам", "цаги",
	}

	for _, org := range orgPatterns {
		if strings.Contains(candidateLower, org) {
			return false
		}
	}

	// Look for typical educational program keywords
	programKeywords := []string{
		"математика", "физика", "информатика", "техника", "технология",
		"инженерия", "наука", "биология", "химия", "электроника",
		"программирование", "системный", "анализ", "управление",
		"биотехнология", "биофизика", "геокосмические", "авиационные",
		"радиотехника", "компьютерные", "природоподобные", "плазменные",
		"ядерные", "фундаментальная", "прикладная", "техническая",
		"перспективных", "наноэлектроника", "моделирование", "теория",
		"естественные", "проектирование", "разработка", "комплексных",
		"бизнес", "приложений",
	}

	hasKeyword := false
	for _, keyword := range programKeywords {
		if strings.Contains(candidateLower, keyword) {
			hasKeyword = true
			break
		}
	}

	return hasKeyword
}

func generateHTTPHeadingSources(programs []MiptProgram, urlMap map[string]MiptURLs) {
	fmt.Println("// Generated MIPT HTTPHeadingSource entries")
	fmt.Println("// Educational program names used as keys (not direction codes)")
	fmt.Println("// URLs extracted from MIPT registry file - REAL URLs, not hallucinated")
	fmt.Println("// Capacities extracted from MIPT capacities file")
	fmt.Println("// TABLE STRUCTURE PARSING - Program names matched to URLs from same table row")
	fmt.Println("// EXACT MATCHES ONLY - No fuzzy matching")
	fmt.Println("")

	// Debug: Print all available program names from registry
	fmt.Println("// Available programs in registry:")
	for programName := range urlMap {
		fmt.Printf("//   '%s'\n", programName)
	}
	fmt.Println("")

	// Debug: Print all program names from capacities
	fmt.Println("// Available programs in capacities:")
	for _, program := range programs {
		fmt.Printf("//   '%s'\n", program.Name)
	}
	fmt.Println("")

	exactMatched := 0
	unmatched := []string{}

	for _, program := range programs {
		urls, found := urlMap[program.Name]

		if found {
			exactMatched++
			fmt.Printf("// EXACT MATCH: '%s'\n", program.Name)
			generateSingleHTTPHeadingSource(program, urls)
		} else {
			unmatched = append(unmatched, program.Name)
		}
	}

	fmt.Printf("\n// EXACT MATCHES SUMMARY: %d/%d programs matched, %d unmatched\n", exactMatched, len(programs), len(unmatched))

	// Detailed anomaly reporting
	if len(unmatched) > 0 {
		fmt.Println("\n// ANOMALIES - Programs from capacities file with NO EXACT MATCH in registry:")
		for _, name := range unmatched {
			fmt.Printf("//   UNMATCHED CAPACITY: '%s'\n", name)
		}
	}

	// Report URLs without matching capacities (exact match only)
	registryOnly := []string{}
	for registryName := range urlMap {
		foundInCapacities := false
		for _, program := range programs {
			if program.Name == registryName {
				foundInCapacities = true
				break
			}
		}
		if !foundInCapacities {
			registryOnly = append(registryOnly, registryName)
		}
	}

	if len(registryOnly) > 0 {
		fmt.Println("\n// ANOMALIES - Programs in registry with NO EXACT MATCH in capacities:")
		for _, name := range registryOnly {
			fmt.Printf("//   UNMATCHED REGISTRY: '%s'\n", name)
		}
	}

	// Additional detailed comparison for better debugging
	fmt.Println("\n// DETAILED COMPARISON FOR DEBUGGING:")
	fmt.Println("// Programs that might be similar but don't match exactly:")
	for _, program := range programs {
		if _, found := urlMap[program.Name]; !found {
			fmt.Printf("//   CAPACITY: '%s'\n", program.Name)
			// Look for potential similar names in registry
			for registryName := range urlMap {
				if strings.Contains(strings.ToLower(program.Name), strings.ToLower(registryName)) ||
					strings.Contains(strings.ToLower(registryName), strings.ToLower(program.Name)) {
					fmt.Printf("//     -> Possible match in registry: '%s'\n", registryName)
				}
			}
		}
	}
}

func generateSingleHTTPHeadingSource(program MiptProgram, urls MiptURLs) {
	fmt.Printf("&mipt.HTTPHeadingSource{\n")
	fmt.Printf("\tPrettyName: \"%s\",\n", program.Name)

	if urls.RegularBVI != "" {
		fmt.Printf("\tRegularBVIListURL: \"%s\",\n", urls.RegularBVI)
	}

	if len(urls.TargetQuotaURLs) > 0 {
		fmt.Printf("\tTargetQuotaListURLs: []string{\n")
		for _, targetURL := range urls.TargetQuotaURLs {
			fmt.Printf("\t\t\"%s\",\n", targetURL)
		}
		fmt.Printf("\t},\n")
	}

	if urls.DedicatedQuota != "" {
		fmt.Printf("\tDedicatedQuotaListURL: \"%s\",\n", urls.DedicatedQuota)
	}

	if urls.SpecialQuota != "" {
		fmt.Printf("\tSpecialQuotaListURL: \"%s\",\n", urls.SpecialQuota)
	}

	fmt.Printf("\tCapacities: core.Capacities{\n")
	fmt.Printf("\t\tRegular: %d,\n", program.Capacities.Regular)
	fmt.Printf("\t\tTargetQuota: %d,\n", program.Capacities.TargetQuota)
	fmt.Printf("\t\tDedicatedQuota: %d,\n", program.Capacities.DedicatedQuota)
	fmt.Printf("\t\tSpecialQuota: %d,\n", program.Capacities.SpecialQuota)
	fmt.Printf("\t},\n")
	fmt.Printf("},\n\n")
}
