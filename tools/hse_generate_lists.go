package main

import (
	"analabit/core"
	"analabit/source/hse"
	"analabit/utils"
	"fmt"
	"github.com/lithammer/fuzzysearch/fuzzy" // Added for fuzzy matching
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var mskHeadingToSource map[string]hse.HttpHeadingSource = make(map[string]hse.HttpHeadingSource)

// printRvalueSource Prints HttpHeadingSource as it can be used by generated code to assign to
// a HttpHeadingSource typed variable.
func printRvalueSource(s *hse.HttpHeadingSource) string {
	return fmt.Sprintf(`hse.HttpHeadingSource{
		RCListURL:         utils.MustParseURL("%s"),
		TQListURL:         utils.MustParseURL("%s"),
		DQListURL:         utils.MustParseURL("%s"),
		SQListURL:         utils.MustParseURL("%s"),
		BListURL:          utils.MustParseURL("%s"),
		HeadingCapacities: %s,
	}`,
		s.RCListURL.String(),
		s.TQListURL.String(),
		s.DQListURL.String(),
		s.SQListURL.String(),
		s.BListURL.String(),
		s.HeadingCapacities.PrintRvalue(), // Use the new PrintRvalue method
	)
}

// ProgramLinks stores the URLs for different competition types for a program.
type ProgramLinks struct {
	BVIURL string
	SQURL  string
	DQURL  string
	TQURL  string
	RCURL  string
}

func normalizeHeadingName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "&nbsp;", " ")
	name = strings.Join(strings.Fields(name), " ") // Normalize spaces
	// Specific normalizations observed from data
	name = strings.ReplaceAll(name, "программа двух дипломов ниу вшэ и университета кёнхи \"экономика и политика азии\"", "программа двух дипломов ниу вшэ и университета кёнхи «экономика и политика в азии»")
	name = strings.ReplaceAll(name, "международная программа «международные отношения и глобальные исследования»/ international program «international relations and global studies»", "международная программа «международные отношения и глобальные исследования»")
	return name
}

func getTextContent(n *html.Node) string {
	if n == nil {
		return ""
	}
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}
	if n.Type == html.ElementNode && n.Data == "a" && n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
		// Prefer direct text of <a> if it's simple, otherwise full content
		return strings.TrimSpace(n.FirstChild.Data)
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

func parseCapacitiesHTML(filePath string) (map[string]core.Capacities, map[string]string, error) {
	capacitiesMap := make(map[string]core.Capacities)
	originalNamesMap := make(map[string]string) // Stores normalizedName -> originalName

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read capacities file %s: %w", filePath, err)
	}

	doc, err := html.Parse(strings.NewReader(string(content)))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse HTML from %s: %w", filePath, err)
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			var cells []string
			var firstCellOriginalName string
			isDataRow := true
			tdCount := 0

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "td" {
					tdCount++
					// Check for colspan, common in header/category rows
					for _, attr := range c.Attr {
						if attr.Key == "colspan" {
							isDataRow = false
							break
						}
					}
					if !isDataRow {
						break
					}
					cellText := getTextContent(c)
					if tdCount == 1 {
						firstCellOriginalName = cellText // Capture original name from the first cell
					}
					cells = append(cells, cellText)
				}
			}

			if isDataRow && len(cells) == 7 {
				// This is likely a program row
				programNameOriginal := strings.TrimSpace(cells[0])
				if programNameOriginal == "" || strings.HasPrefix(strings.ToLower(programNameOriginal), "итого") || strings.HasPrefix(strings.ToLower(programNameOriginal), "всего:") {
					return // Skip summary rows or empty names
				}
				// Check if second cell looks like a number or "-", not a sub-heading
				if strings.TrimSpace(cells[1]) == "" && !strings.Contains(programNameOriginal, "НАПРАВЛЕНИЕ ПОДГОТОВКИ") && !strings.Contains(programNameOriginal, "очная форма обучения") {
					// Potentially a program with no budget places, check if it's a known pattern
					if strings.TrimSpace(cells[2]) == "-" && strings.TrimSpace(cells[3]) == "-" && strings.TrimSpace(cells[4]) == "-" {
						// All quotas are '-', likely a fully paid program listed among budget ones.
					} else {
						// If it's not a clear non-data row, but second cell is empty, log it.
						// log.Printf("Skipping row with empty KCP but not a clear header: %s", programNameOriginal)
						// return
					}
				}

				kcp := extractIntFromCell(cells[1])
				specialQuota := extractIntFromCell(cells[2])
				targetQuota := extractIntFromCell(cells[3])
				dedicatedQuota := extractIntFromCell(cells[4])
				// paidPlaces := extractIntFromCell(cells[5]) // Not used for core.Capacities

				if programNameOriginal == "Математика" && kcp == 60 { // ensure we are on the right row
					// This is a specific check if needed, but general logic should handle it.
				}

				// General capacity calculation
				generalCapacity := kcp - specialQuota - targetQuota - dedicatedQuota
				if generalCapacity < 0 {
					generalCapacity = 0 // Cannot be negative
				}

				if kcp == 0 && specialQuota == 0 && targetQuota == 0 && dedicatedQuota == 0 && generalCapacity == 0 {
					//log.Printf("Info: Program '%s' from capacities file has zero KCP and zero quotas. Added with zero capacities (may still have BVI or other links).", programNameOriginal)
				}

				normalizedName := normalizeHeadingName(programNameOriginal)
				if _, exists := capacitiesMap[normalizedName]; exists {
					log.Printf("Warning: Duplicate normalized program name in capacities file: '%s' (from original '%s')", normalizedName, programNameOriginal)
				}
				capacitiesMap[normalizedName] = core.Capacities{
					General:        generalCapacity,
					TargetQuota:    targetQuota,
					DedicatedQuota: dedicatedQuota,
					SpecialQuota:   specialQuota,
				}
				originalNamesMap[normalizedName] = firstCellOriginalName
			} else if len(cells) > 0 && (strings.Contains(cells[0], "НАПРАВЛЕНИЕ ПОДГОТОВКИ") || strings.Contains(cells[0], "очная форма обучения") || strings.Contains(cells[0], "очно-заочная форма обучения")) {
				// This is a header/category row, skip.
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	// More robust table finding: look for a table with a thead containing specific headers
	var findMainTable func(*html.Node) *html.Node
	findMainTable = func(n *html.Node) *html.Node {
		if n.Type == html.ElementNode && n.Data == "table" {
			// Optional: log table attributes
			var headerTexts []string
			var visitThead func(*html.Node)
			visitThead = func(node *html.Node) {
				if node.Type == html.ElementNode && node.Data == "th" {
					headerTexts = append(headerTexts, getTextContent(node))
				}
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					visitThead(c)
				}
			}

			for cTheadCandidate := n.FirstChild; cTheadCandidate != nil; cTheadCandidate = cTheadCandidate.NextSibling {
				if cTheadCandidate.Type == html.ElementNode && cTheadCandidate.Data == "thead" {
					visitThead(cTheadCandidate)

					// Check if we found enough characteristic headers
					if len(headerTexts) >= 5 &&
						strings.Contains(headerTexts[0], "Направление подготовки") &&
						strings.Contains(headerTexts[1], "Бюджетные места") &&
						strings.Contains(headerTexts[2], "особое право") { // Corrected typo here
						// This is likely the main table, find its tbody
						for tbodyCandidate := cTheadCandidate.NextSibling; tbodyCandidate != nil; tbodyCandidate = tbodyCandidate.NextSibling {
							if tbodyCandidate.Type == html.ElementNode && tbodyCandidate.Data == "tbody" {
								return tbodyCandidate
							}
						}
					}
					break // only check first thead found in this table
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if foundTbody := findMainTable(c); foundTbody != nil {
				return foundTbody
			}
		}
		return nil
	}

	tableBody := findMainTable(doc)
	if tableBody == nil {
		return nil, nil, fmt.Errorf("could not find the main data table body in %s", filePath)
	}
	f(tableBody)

	return capacitiesMap, originalNamesMap, nil
}

func parseLinksHTML(filePath string) (map[string]ProgramLinks, map[string]string, error) {
	linksMap := make(map[string]ProgramLinks)
	originalNamesMap := make(map[string]string) // normalizedName -> originalName

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read links file %s: %w", filePath, err)
	}

	doc, err := html.Parse(strings.NewReader(string(content)))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse HTML from %s: %w", filePath, err)
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			var cells []*html.Node
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "td" {
					cells = append(cells, c)
				}
			}

			if len(cells) == 7 { // Expect 7 columns for data rows
				programNameOriginal := getTextContent(cells[0])
				if programNameOriginal == "" || programNameOriginal == "Образовательная программа" { // Skip header or empty
					return
				}

				var links ProgramLinks
				extractLink := func(cell *html.Node) string {
					if cell.FirstChild != nil && cell.FirstChild.Type == html.ElementNode && cell.FirstChild.Data == "a" {
						for _, attr := range cell.FirstChild.Attr {
							if attr.Key == "href" {
								return strings.TrimSpace(attr.Val)
							}
						}
					} else if cell.FirstChild != nil && cell.FirstChild.NextSibling != nil && cell.FirstChild.NextSibling.Type == html.ElementNode && cell.FirstChild.NextSibling.Data == "a" {
						// Sometimes there's a text node (like a line break) then the <a>
						for _, attr := range cell.FirstChild.NextSibling.Attr {
							if attr.Key == "href" {
								return strings.TrimSpace(attr.Val)
							}
						}
					}
					return ""
				}

				links.BVIURL = extractLink(cells[1])
				links.SQURL = extractLink(cells[2])
				links.DQURL = extractLink(cells[3])
				links.TQURL = extractLink(cells[4])
				links.RCURL = extractLink(cells[5])
				// links.PaidURL = extractLink(cells[6]) // Not used

				normalizedName := normalizeHeadingName(programNameOriginal)
				if _, exists := linksMap[normalizedName]; exists {
					log.Printf("Warning: Duplicate normalized program name in links file: '%s' (from original '%s')", normalizedName, programNameOriginal)
				}
				linksMap[normalizedName] = links
				originalNamesMap[normalizedName] = programNameOriginal
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	var findTableBody func(*html.Node) *html.Node
	findTableBody = func(n *html.Node) *html.Node {
		if n.Type == html.ElementNode && n.Data == "tbody" {
			// Basic check: first row, first cell (th or td) contains "Образовательная программа"
			// This is for the header row. The actual data rows follow.
			headerRow := n.FirstChild
			for headerRow != nil && headerRow.Type != html.ElementNode { // Find first element node (tr)
				headerRow = headerRow.NextSibling
			}
			if headerRow != nil && headerRow.Data == "tr" {
				firstHeaderCell := headerRow.FirstChild
				for firstHeaderCell != nil && firstHeaderCell.Type != html.ElementNode { // Find first element node (th/td)
					firstHeaderCell = firstHeaderCell.NextSibling
				}
				if firstHeaderCell != nil && (firstHeaderCell.Data == "th" || firstHeaderCell.Data == "td") {
					if strings.Contains(getTextContent(firstHeaderCell), "Образовательная программа") {
						return n // This tbody contains the links table
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if found := findTableBody(c); found != nil {
				return found
			}
		}
		return nil
	}

	tableBody := findTableBody(doc)
	if tableBody == nil {
		return nil, nil, fmt.Errorf("could not find the main data table body in %s", filePath)
	}
	f(tableBody)

	return linksMap, originalNamesMap, nil
}

// printMskSourcesList prints a ready to paste func `mskSourcesList` declaration
// see defs/hse.go for mskSourcesList() definition
func printMskSourcesListFunc() {
	capacitiesPath := "tools/resources/hse_msk_places.html"
	linksPath := "tools/resources/hse_msk_links.html"

	capacitiesMap, capOriginalNames, err := parseCapacitiesHTML(capacitiesPath)
	if err != nil {
		log.Fatalf("Error parsing capacities: %v", err)
	}

	linksMap, linkOriginalNames, err := parseLinksHTML(linksPath)
	if err != nil {
		log.Fatalf("Error parsing links: %v", err)
	}

	// Define validateUrl once here so it can be used by both loops
	validateUrl := func(u string) string {
		if u == "" || u == "." || u == "/" || u == "#" || strings.TrimSpace(u) == "&nbsp;" {
			return ""
		}
		if strings.HasPrefix(u, "/storage/public_report_2024") {
			return "https://enrol.hse.ru" + u
		}
		return u
	}

	var headingNamesForSorting []string
	for normalizedName := range capacitiesMap {
		headingNamesForSorting = append(headingNamesForSorting, normalizedName)
	}
	sort.Strings(headingNamesForSorting) // Sort by normalized name for consistent output

	// Populate mskHeadingToSource based on capacitiesMap and linksMap (initial exact matches)
	for _, normalizedName := range headingNamesForSorting {
		caps := capacitiesMap[normalizedName]
		links, linksOk := linksMap[normalizedName]
		if !linksOk {
			links = ProgramLinks{} // Default to empty links if no exact match
		}

		sourceEntry := hse.HttpHeadingSource{
			RCListURL:         utils.MustParseURL(validateUrl(links.RCURL)),
			TQListURL:         utils.MustParseURL(validateUrl(links.TQURL)),
			DQListURL:         utils.MustParseURL(validateUrl(links.DQURL)),
			SQListURL:         utils.MustParseURL(validateUrl(links.SQURL)),
			BListURL:          utils.MustParseURL(validateUrl(links.BVIURL)),
			HeadingCapacities: caps,
		}
		mskHeadingToSource[normalizedName] = sourceEntry
	}

	// --- Enhanced Fuzzy Matching Logic ---
	var capacitiesOriginalNamesWithoutExactLinks []string
	for normCapName := range capacitiesMap {
		entry := mskHeadingToSource[normCapName]
		isMissingLinks := (entry.RCListURL == nil || entry.RCListURL.String() == "") &&
			(entry.TQListURL == nil || entry.TQListURL.String() == "") &&
			(entry.DQListURL == nil || entry.DQListURL.String() == "") &&
			(entry.SQListURL == nil || entry.SQListURL.String() == "") &&
			(entry.BListURL == nil || entry.BListURL.String() == "")
		if isMissingLinks {
			capacitiesOriginalNamesWithoutExactLinks = append(capacitiesOriginalNamesWithoutExactLinks, capOriginalNames[normCapName])
		}
	}
	sort.Strings(capacitiesOriginalNamesWithoutExactLinks)

	var linksOriginalNamesWithoutExactCaps []string
	// normalizedCapNamesSet was previously defined but not used for this logic directly.
	// The logic for linksOriginalNamesWithoutExactCaps should find link names that don't have a direct capacity match
	// AND whose links are not already used by any capacity entry.
	for normLinkName, origLinkName := range linkOriginalNames {
		// Check if this normalized link name (from links file) has a corresponding capacity entry by the same normalized name
		var isDirectlyMatched bool
		if _, exists := mskHeadingToSource[normLinkName]; exists {
			isDirectlyMatched = true
		}

		// More importantly, check if the links from this originalLinkName are used by ANY capacity entry
		isUsedByAnyCapacity := false
		linksFromThisOrigName := linksMap[normLinkName]
		for _, sourceEntry := range mskHeadingToSource {
			if (linksFromThisOrigName.RCURL != "" && sourceEntry.RCListURL != nil && sourceEntry.RCListURL.String() == utils.MustParseURL(validateUrl(linksFromThisOrigName.RCURL)).String()) ||
				(linksFromThisOrigName.TQURL != "" && sourceEntry.TQListURL != nil && sourceEntry.TQListURL.String() == utils.MustParseURL(validateUrl(linksFromThisOrigName.TQURL)).String()) ||
				(linksFromThisOrigName.DQURL != "" && sourceEntry.DQListURL != nil && sourceEntry.DQListURL.String() == utils.MustParseURL(validateUrl(linksFromThisOrigName.DQURL)).String()) ||
				(linksFromThisOrigName.SQURL != "" && sourceEntry.SQListURL != nil && sourceEntry.SQListURL.String() == utils.MustParseURL(validateUrl(linksFromThisOrigName.SQURL)).String()) ||
				(linksFromThisOrigName.BVIURL != "" && sourceEntry.BListURL != nil && sourceEntry.BListURL.String() == utils.MustParseURL(validateUrl(linksFromThisOrigName.BVIURL)).String()) {
				isUsedByAnyCapacity = true
				break
			}
		}

		if !isDirectlyMatched && !isUsedByAnyCapacity {
			alreadyAdded := false
			for _, addedName := range linksOriginalNamesWithoutExactCaps {
				if addedName == origLinkName {
					alreadyAdded = true
					break
				}
			}
			if !alreadyAdded {
				linksOriginalNamesWithoutExactCaps = append(linksOriginalNamesWithoutExactCaps, origLinkName)
			}
		}
	}
	sort.Strings(linksOriginalNamesWithoutExactCaps)

	fuzzyMatchedPairsLog := []string{}
	processedCapOriginalNames := make(map[string]bool)
	usedOriginalLinkNames := make(map[string]bool)

	if len(capacitiesOriginalNamesWithoutExactLinks) > 0 && len(linksOriginalNamesWithoutExactCaps) > 0 {
		availableLinksForFuzzy := make([]string, 0, len(linksOriginalNamesWithoutExactCaps))
		for _, name := range linksOriginalNamesWithoutExactCaps {
			availableLinksForFuzzy = append(availableLinksForFuzzy, name)
		}

		for _, originalCapName := range capacitiesOriginalNamesWithoutExactLinks {
			if len(availableLinksForFuzzy) == 0 {
				break
			}

			matches := fuzzy.RankFindFold(originalCapName, availableLinksForFuzzy)
			if len(matches) > 0 {
				bestMatch := matches[0]
				originalLinkNameMatched := bestMatch.Target

				normalizedCapName := normalizeHeadingName(originalCapName)
				normalizedLinkNameMatched := normalizeHeadingName(originalLinkNameMatched)

				if linksData, ok := linksMap[normalizedLinkNameMatched]; ok {
					entryToUpdate, entryExists := mskHeadingToSource[normalizedCapName]
					if !entryExists {
						log.Printf("Error: Capacity name '%s' (normalized: '%s') not found in mskHeadingToSource for fuzzy update. This should not happen.", originalCapName, normalizedCapName)
						continue
					}

					entryToUpdate.RCListURL = utils.MustParseURL(validateUrl(linksData.RCURL))
					entryToUpdate.TQListURL = utils.MustParseURL(validateUrl(linksData.TQURL))
					entryToUpdate.DQListURL = utils.MustParseURL(validateUrl(linksData.DQURL))
					entryToUpdate.SQListURL = utils.MustParseURL(validateUrl(linksData.SQURL))
					entryToUpdate.BListURL = utils.MustParseURL(validateUrl(linksData.BVIURL))
					mskHeadingToSource[normalizedCapName] = entryToUpdate

					fuzzyMatchedPairsLog = append(fuzzyMatchedPairsLog, fmt.Sprintf("(Cap: '%s' ~ Link: '%s' [Dist: %d])", originalCapName, originalLinkNameMatched, bestMatch.Distance))
					processedCapOriginalNames[originalCapName] = true
					usedOriginalLinkNames[originalLinkNameMatched] = true

					tempLinks := []string{}
					for _, ln := range availableLinksForFuzzy {
						if ln != originalLinkNameMatched {
							tempLinks = append(tempLinks, ln)
						}
					}
					availableLinksForFuzzy = tempLinks
				} else {
					log.Printf("Warning: Normalized link name '%s' (from original '%s') not found in linksMap during fuzzy matching. This is unexpected.", normalizedLinkNameMatched, originalLinkNameMatched)
				}
			}
		}
	}

	if len(fuzzyMatchedPairsLog) > 0 {
		log.Printf("Info: Applied fuzzy matches for MSK: [ %s ]", strings.Join(fuzzyMatchedPairsLog, ", "))
	}

	var finalCapacitiesWithoutLinks []string
	for _, originalCapName := range capacitiesOriginalNamesWithoutExactLinks {
		if !processedCapOriginalNames[originalCapName] {
			finalCapacitiesWithoutLinks = append(finalCapacitiesWithoutLinks, originalCapName)
		}
	}
	if len(finalCapacitiesWithoutLinks) > 0 {
		sort.Strings(finalCapacitiesWithoutLinks)
		log.Printf("Info: MSK Programs in CAPACITIES list but still no LINKS (after fuzzy attempts): [ %s ]", strings.Join(finalCapacitiesWithoutLinks, ", "))
	}

	var finalLinksWithoutCapacities []string
	for _, originalLinkName := range linksOriginalNamesWithoutExactCaps {
		if !usedOriginalLinkNames[originalLinkName] {
			finalLinksWithoutCapacities = append(finalLinksWithoutCapacities, originalLinkName)
		}
	}
	if len(finalLinksWithoutCapacities) > 0 {
		sort.Strings(finalLinksWithoutCapacities)
		log.Printf("Info: MSK Programs in LINKS list but not matched to CAPACITIES (after fuzzy attempts): [ %s ]", strings.Join(finalLinksWithoutCapacities, ", "))
	}

	var sb strings.Builder
	sb.WriteString("package defs\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"analabit/core\"\n")
	sb.WriteString("\t\"analabit/source\"\n")
	sb.WriteString("\t\"analabit/source/hse\"\n")
	sb.WriteString("\t\"analabit/utils\"\n")
	sb.WriteString(")\n\n")
	sb.WriteString("// mskSourcesList returns a list of HeadingSource for HSE Moscow.\n")
	sb.WriteString("// Generated by tools/hse_generate_lists.go\n")
	sb.WriteString("func mskSourcesList() []source.HeadingSource {\n")
	sb.WriteString("\treturn []source.HeadingSource{\n")

	var mainOutput strings.Builder
	var zeroCapacityOutput strings.Builder
	var missingURLsOutput strings.Builder

	for _, normalizedName := range headingNamesForSorting {
		sourceEntry, ok := mskHeadingToSource[normalizedName]
		if !ok {
			log.Printf("Error: Normalized name '%s' from capacitiesMap not found in mskHeadingToSource during final print. Skipping.", normalizedName)
			continue
		}

		originalName, nameOk := capOriginalNames[normalizedName]
		if !nameOk {
			originalName = normalizedName
			log.Printf("Warning: Original name not found for normalized name '%s' during print. Using normalized name.", normalizedName)
		}

		entryComment := fmt.Sprintf("\t\t// %s\n", strings.ReplaceAll(originalName, "`", "'"))
		entryCode := fmt.Sprintf("\t\t&%s,\n", printRvalueSource(&sourceEntry))

		isZeroCapacity := sourceEntry.HeadingCapacities.General == 0 &&
			sourceEntry.HeadingCapacities.TargetQuota == 0 &&
			sourceEntry.HeadingCapacities.DedicatedQuota == 0 &&
			sourceEntry.HeadingCapacities.SpecialQuota == 0

		isMissingAllURLs := (sourceEntry.RCListURL == nil || sourceEntry.RCListURL.String() == "") &&
			(sourceEntry.TQListURL == nil || sourceEntry.TQListURL.String() == "") &&
			(sourceEntry.DQListURL == nil || sourceEntry.DQListURL.String() == "") &&
			(sourceEntry.SQListURL == nil || sourceEntry.SQListURL.String() == "") &&
			(sourceEntry.BListURL == nil || sourceEntry.BListURL.String() == "")

		if isMissingAllURLs {
			missingURLsOutput.WriteString(entryComment)
			missingURLsOutput.WriteString(entryCode)
		} else if isZeroCapacity {
			zeroCapacityOutput.WriteString(entryComment)
			zeroCapacityOutput.WriteString(entryCode)
		} else {
			mainOutput.WriteString(entryComment)
			mainOutput.WriteString(entryCode)
		}
	}

	sb.WriteString(mainOutput.String())

	if zeroCapacityOutput.Len() > 0 {
		sb.WriteString("\n\t\t// TODO The following headings do not have capacities determined:\n\n")
		sb.WriteString(zeroCapacityOutput.String())
	}

	if missingURLsOutput.Len() > 0 {
		sb.WriteString("\n\t\t// TODO The following headings do not have list URLs determined:\n\n")
		sb.WriteString(missingURLsOutput.String())
	}

	sb.WriteString("\t}\n}\n")
	fmt.Print(sb.String())
}
