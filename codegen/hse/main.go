package main

import (
	"analabit/core"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// ProgramInfo stores information about an educational program
type ProgramInfo struct {
	Name           string
	NormalizedName string
	URL            string
	Capacities     core.Capacities
}

// Campus configuration
type CampusConfig struct {
	Name       string
	LinksPath  string
	PlacesPath string
	URLPrefix  string
}

var campusConfigs = map[string]CampusConfig{
	"msk": {
		Name:       "Москва",
		LinksPath:  "lists_and_caps/lists.html",
		PlacesPath: "lists_and_caps/places_msk.html",
		URLPrefix:  "moscow",
	},
	"nn": {
		Name:       "Нижний Новгород",
		LinksPath:  "lists_and_caps/lists.html",
		PlacesPath: "lists_and_caps/places_nn.html",
		URLPrefix:  "nn",
	},
	"perm": {
		Name:       "Пермь",
		LinksPath:  "lists_and_caps/lists.html",
		PlacesPath: "lists_and_caps/places_perm.html",
		URLPrefix:  "perm",
	},
	"spb": {
		Name:       "Санкт-Петербург",
		LinksPath:  "lists_and_caps/lists.html",
		PlacesPath: "lists_and_caps/places_spb.html",
		URLPrefix:  "spb",
	},
}

func main() {
	var campus string
	flag.StringVar(&campus, "c", "", "Campus to process: msk, nn, perm, spb")
	flag.StringVar(&campus, "campus", "", "Campus to process: msk, nn, perm, spb")
	flag.Parse()

	if campus == "" {
		log.Fatal("Campus must be specified with -c or --campus flag")
	}

	config, exists := campusConfigs[campus]
	if !exists {
		log.Fatalf("Invalid campus: %s. Valid options: msk, nn, perm, spb", campus)
	}

	// Parse links from lists.html
	programs, err := parseLinksForCampus(config.LinksPath, config.Name)
	if err != nil {
		log.Fatalf("Error parsing links: %v", err)
	}

	// Parse capacities from places file
	capacities, err := parseCapacitiesHTML(config.PlacesPath)
	if err != nil {
		log.Fatalf("Error parsing capacities: %v", err)
	}

	// Match programs with capacities
	matchedPrograms, filteredPrograms, anomalies := matchProgramsWithCapacities(programs, capacities)

	// Output results
	outputResults(matchedPrograms, filteredPrograms, anomalies)
}

func normalizeHeadingName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "&nbsp;", " ")
	name = strings.Join(strings.Fields(name), " ") // Normalize spaces

	// Remove common suffixes/prefixes that might cause mismatches
	name = strings.TrimSpace(name)

	// Handle common variations
	name = strings.ReplaceAll(name, "«", "\"")
	name = strings.ReplaceAll(name, "»", "\"")
	name = strings.ReplaceAll(name, `"`, "\"")
	name = strings.ReplaceAll(name, `"`, "\"")

	// Remove parenthetical clarifications that cause mismatches
	re := regexp.MustCompile(`\([^)]*\)`)
	name = re.ReplaceAllString(name, "")
	name = strings.Join(strings.Fields(name), " ") // Re-normalize after removal

	// Handle specific problematic cases
	name = strings.ReplaceAll(name, "ниу вшэ", "")
	name = strings.ReplaceAll(name, "нию вшэ", "")
	name = strings.ReplaceAll(name, "центра педагогического мастерства", "")
	name = strings.ReplaceAll(name, ": правовое регулирование бизнеса", "")
	name = strings.ReplaceAll(name, ": цифровой юрист", "")

	// Remove specific form indicators
	name = strings.ReplaceAll(name, "очная форма обучения", "")
	name = strings.ReplaceAll(name, "очно-заочная форма обучения", "")
	name = strings.ReplaceAll(name, "(онлайн)", "")

	// Clean up multiple spaces and trim
	name = strings.Join(strings.Fields(name), " ")
	name = strings.TrimSpace(name)

	return name
}

func getTextContent(n *html.Node) string {
	if n == nil {
		return ""
	}
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}

	var sb strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		textContent := getTextContent(c)
		if textContent != "" {
			if sb.Len() > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(textContent)
		}
	}
	return strings.TrimSpace(sb.String())
}

var numRegexp = regexp.MustCompile(`^\s*(\d+)`)

func extractIntFromCell(cellContent string) int {
	if strings.TrimSpace(cellContent) == "-" {
		return 0
	}
	matches := numRegexp.FindStringSubmatch(cellContent)
	if len(matches) > 1 {
		val, err := strconv.Atoi(matches[1])
		if err == nil {
			return val
		}
	}
	return 0 // Default to 0 if not parsable or "-"
}

func parseLinksForCampus(filePath, campusName string) ([]ProgramInfo, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", filePath, err)
	}

	doc, err := html.Parse(strings.NewReader(string(content)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	var programs []ProgramInfo

	// Find the section for the specific campus
	var campusSection *html.Node
	var findCampusSection func(*html.Node)
	findCampusSection = func(n *html.Node) {
		if campusSection != nil {
			return
		}

		if n.Type == html.ElementNode && n.Data == "h2" {
			text := getTextContent(n)
			if strings.TrimSpace(text) == campusName {
				campusSection = n.Parent
				return
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findCampusSection(c)
		}
	}
	findCampusSection(doc)

	if campusSection == nil {
		return nil, fmt.Errorf("campus section not found for %s", campusName)
	}

	// Find the table within this campus section
	var findTable func(*html.Node)
	findTable = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" {
			parseTable(n, &programs)
			return
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findTable(c)
		}
	}
	findTable(campusSection)

	return programs, nil
}

func parseTable(table *html.Node, programs *[]ProgramInfo) {
	var parseTableRows func(*html.Node)
	parseTableRows = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			cells := []*html.Node{}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "td" {
					cells = append(cells, c)
				}
			}

			if len(cells) == 2 {
				nameCell := cells[0]
				linkCell := cells[1]

				name := getTextContent(nameCell)
				if name != "" && name != "Наименование конкурса" {
					// Find the link URL
					var url string
					var findLink func(*html.Node)
					findLink = func(n *html.Node) {
						if n.Type == html.ElementNode && n.Data == "a" {
							for _, attr := range n.Attr {
								if attr.Key == "href" {
									url = attr.Val
									return
								}
							}
						}
						for c := n.FirstChild; c != nil; c = c.NextSibling {
							findLink(c)
						}
					}
					findLink(linkCell)

					if url != "" {
						*programs = append(*programs, ProgramInfo{
							Name:           name,
							NormalizedName: normalizeHeadingName(name),
							URL:            url,
						})
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseTableRows(c)
		}
	}
	parseTableRows(table)
}

func parseCapacitiesHTML(filePath string) (map[string]core.Capacities, error) {
	capacitiesMap := make(map[string]core.Capacities)

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", filePath, err)
	}

	doc, err := html.Parse(strings.NewReader(string(content)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Find the table with capacity information
	var table *html.Node
	var findTable func(*html.Node)
	findTable = func(n *html.Node) {
		if table != nil {
			return
		}
		if n.Type == html.ElementNode && n.Data == "table" {
			// Check if this table has the right headers
			hasCorrectHeaders := false
			var checkHeaders func(*html.Node)
			checkHeaders = func(tn *html.Node) {
				if tn.Type == html.ElementNode && (tn.Data == "th" || tn.Data == "td") {
					text := getTextContent(tn)
					// Check for various header formats across campuses
					if strings.Contains(text, "Бюджетные места") ||
						strings.Contains(text, "КЦП") ||
						strings.Contains(text, "Распределение КЦП") ||
						strings.Contains(text, "контрольных цифр приема") ||
						strings.Contains(text, "бюджетных ассигнований") {
						hasCorrectHeaders = true
					}
				}
				for c := tn.FirstChild; c != nil; c = c.NextSibling {
					checkHeaders(c)
				}
			}
			checkHeaders(n)

			if hasCorrectHeaders {
				table = n
				return
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findTable(c)
		}
	}
	findTable(doc)

	if table == nil {
		return nil, fmt.Errorf("capacity table not found")
	}

	parseCapacityTable(table, capacitiesMap)
	return capacitiesMap, nil
}

func parseCapacityTable(table *html.Node, capacitiesMap map[string]core.Capacities) {
	var parseTableRows func(*html.Node)
	parseTableRows = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			cells := []*html.Node{}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "td" {
					cells = append(cells, c)
				}
			}

			// Handle different table formats (7 cells for Moscow, potentially different for others)
			if len(cells) >= 6 {
				programName := getTextContent(cells[0])

				// Skip header rows, empty rows, and summary rows
				if programName == "" ||
					strings.Contains(programName, "Направление подготовки") ||
					strings.Contains(programName, "Специальность/образовательная программа") ||
					strings.Contains(programName, "очная форма") ||
					strings.Contains(programName, "Очная форма") ||
					strings.Contains(strings.ToLower(programName), "всего") ||
					strings.Contains(strings.ToLower(programName), "итого") {
					return
				}

				// Extract link text from program name cell or second cell for Perm format
				var linkText string
				var findLinkText func(*html.Node)
				findLinkText = func(tn *html.Node) {
					if tn.Type == html.ElementNode && tn.Data == "a" {
						linkText = getTextContent(tn)
					}
					for c := tn.FirstChild; c != nil; c = c.NextSibling {
						findLinkText(c)
					}
				}
				findLinkText(cells[0])

				// If no link in first cell, try second cell (Perm format)
				if linkText == "" && len(cells) > 1 {
					findLinkText(cells[1])
					if linkText != "" {
						programName = linkText
					}
				} else if linkText != "" {
					programName = linkText
				}

				if programName != "" {
					// Different campuses may have different column orders
					// Try to parse based on available columns
					var kcp, specialQuota, targetQuota, dedicatedQuota int

					if len(cells) >= 8 {
						// Format 3: Perm format (8 columns)
						// [Direction, Program, Total KCP, Special, Dedicated, Target, Paid, Foreign]
						kcpStr := getTextContent(cells[2])            // Total budget seats
						specialQuotaStr := getTextContent(cells[3])   // Special quota
						dedicatedQuotaStr := getTextContent(cells[4]) // Dedicated quota
						targetQuotaStr := getTextContent(cells[5])    // Target quota

						kcp = extractIntFromCell(kcpStr)
						specialQuota = extractIntFromCell(specialQuotaStr)
						dedicatedQuota = extractIntFromCell(dedicatedQuotaStr)
						targetQuota = extractIntFromCell(targetQuotaStr)
					} else if len(cells) >= 7 {
						// Format 1: Moscow format (7 columns)
						kcpStr := getTextContent(cells[1])            // Total budget seats
						specialQuotaStr := getTextContent(cells[2])   // Special quota
						targetQuotaStr := getTextContent(cells[3])    // Target quota
						dedicatedQuotaStr := getTextContent(cells[4]) // Dedicated quota

						kcp = extractIntFromCell(kcpStr)
						specialQuota = extractIntFromCell(specialQuotaStr)
						targetQuota = extractIntFromCell(targetQuotaStr)
						dedicatedQuota = extractIntFromCell(dedicatedQuotaStr)
					} else if len(cells) >= 6 {
						// Format 2: SPB format or similar (6+ columns)
						kcpStr := getTextContent(cells[1])            // Total budget seats
						specialQuotaStr := getTextContent(cells[2])   // Special quota
						dedicatedQuotaStr := getTextContent(cells[3]) // Dedicated quota
						targetQuotaStr := getTextContent(cells[4])    // Target quota

						kcp = extractIntFromCell(kcpStr)
						specialQuota = extractIntFromCell(specialQuotaStr)
						dedicatedQuota = extractIntFromCell(dedicatedQuotaStr)
						targetQuota = extractIntFromCell(targetQuotaStr)
					}

					// Calculate regular places: KCP - all quotas
					regular := kcp - specialQuota - targetQuota - dedicatedQuota
					if regular < 0 {
						regular = 0
					}

					normalizedName := normalizeHeadingName(programName)
					capacitiesMap[normalizedName] = core.Capacities{
						Regular:        regular,
						TargetQuota:    targetQuota,
						DedicatedQuota: dedicatedQuota,
						SpecialQuota:   specialQuota,
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseTableRows(c)
		}
	}
	parseTableRows(table)
}

func matchProgramsWithCapacities(programs []ProgramInfo, capacities map[string]core.Capacities) ([]ProgramInfo, []ProgramInfo, []string) {
	var matched []ProgramInfo
	var filtered []ProgramInfo
	var anomalies []string

	// Track which capacities were matched
	usedCapacities := make(map[string]bool)

	for _, program := range programs {
		if caps, exists := capacities[program.NormalizedName]; exists {
			program.Capacities = caps
			totalKCP := caps.Regular + caps.TargetQuota + caps.DedicatedQuota + caps.SpecialQuota
			if totalKCP > 0 {
				matched = append(matched, program)
			} else {
				filtered = append(filtered, program)
			}
			usedCapacities[program.NormalizedName] = true
		} else {
			anomalies = append(anomalies, fmt.Sprintf("Program found in lists but not in capacities: %s", program.Name))
		}
	}

	// Check for capacities without matching programs
	for normalizedName, _ := range capacities {
		if !usedCapacities[normalizedName] {
			anomalies = append(anomalies, fmt.Sprintf("Program found in capacities but not in lists: %s", normalizedName))
		}
	}

	return matched, filtered, anomalies
}

func outputResults(programs []ProgramInfo, filtered []ProgramInfo, anomalies []string) {
	// Sort programs by name for consistent output
	sort.Slice(programs, func(i, j int) bool {
		return programs[i].Name < programs[j].Name
	})

	for _, program := range programs {
		fmt.Printf("// %s\n", program.Name)
		fmt.Printf("&hse.HTTPHeadingSource{\n")
		fmt.Printf("    URL: \"%s\",\n", program.URL)
		fmt.Printf("    Capacities: core.Capacities{\n")
		fmt.Printf("        Regular:        %d,\n", program.Capacities.Regular)
		fmt.Printf("        TargetQuota:    %d,\n", program.Capacities.TargetQuota)
		fmt.Printf("        DedicatedQuota: %d,\n", program.Capacities.DedicatedQuota)
		fmt.Printf("        SpecialQuota:   %d,\n", program.Capacities.SpecialQuota)
		fmt.Printf("    },\n")
		fmt.Printf("},\n")
	}

	if len(filtered) > 0 {
		fmt.Printf("\n\n// ===\n")
		fmt.Printf("// FILTERED OUT\n")
		fmt.Printf("// ===\n")
		for _, program := range filtered {
			fmt.Printf("// %s (KCP = 0)\n", program.Name)
		}
	}

	if len(anomalies) > 0 {
		fmt.Printf("\n\n// ===\n")
		fmt.Printf("// FOUND ANOMALIES\n")
		fmt.Printf("// ===\n")
		for _, anomaly := range anomalies {
			fmt.Printf("// %s\n", anomaly)
		}
	}
}
