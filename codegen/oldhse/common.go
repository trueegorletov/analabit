package main

import (
	"analabit/core"
	"analabit/core/source/oldhse"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// printRvalueSource Prints HttpHeadingSource as it can be used by generated code to assign to
// a HttpHeadingSource typed variable.
func printRvalueSource(s *oldhse.HttpHeadingSource) string {
	return fmt.Sprintf(`oldhse.HttpHeadingSource{
		RCListURL:         "%s",
		TQListURL:         "%s",
		DQListURL:         "%s",
		SQListURL:         "%s",
		BListURL:          "%s",
		Capacities: analabit.%s,
	}`,
		s.RCListURL,
		s.TQListURL,
		s.DQListURL,
		s.SQListURL,
		s.BListURL,
		s.Capacities.PrintRvalue(), // Use the new PrintRvalue method
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

	isPermFileGlobal := strings.Contains(filePath, "_perm_places.html")

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			var cellsContent []string
			var programNameFromCell string    // Extracted program name
			var firstCellThTextForPerm string // Content of the first <th>, if it's Perm

			isRowSuitableForData := true
			tempCellIndex := 0

			// Iterate over children of <tr> (which are <th> or <td>)
			for cChildOfTr := n.FirstChild; cChildOfTr != nil; cChildOfTr = cChildOfTr.NextSibling {
				if cChildOfTr.Type == html.ElementNode {
					// Rule 1: Skip row if any cell has colspan or rowspan
					for _, attr := range cChildOfTr.Attr {
						if attr.Key == "colspan" || attr.Key == "rowspan" {
							isRowSuitableForData = false
							break
						}
					}
					if !isRowSuitableForData {
						break
					}

					// Rule 2: Collect content from <td> or (if Perm) <th>
					if cChildOfTr.Data == "td" || (isPermFileGlobal && cChildOfTr.Data == "th") {
						cellText := getTextContent(cChildOfTr)
						cellsContent = append(cellsContent, cellText)

						if isPermFileGlobal {
							if cChildOfTr.Data == "th" && tempCellIndex == 0 {
								firstCellThTextForPerm = cellText
							}
							// For Perm, program name is in the first <td>, which is the second cell overall.
							if cChildOfTr.Data == "td" && tempCellIndex == 1 {
								programNameFromCell = cellText
							}
						} else { // Not Perm
							// Program name is in the first <td>
							if cChildOfTr.Data == "td" && tempCellIndex == 0 {
								programNameFromCell = cellText
							}
						}
						tempCellIndex++
					}
				}
			}

			if !isRowSuitableForData { // Row had colspan/rowspan, skip processing its content further as a data row
				// Continue recursion for its children normally via the loop at the end of f
			} else {
				programNameFromCell = strings.TrimSpace(programNameFromCell)

				// Rule 3: Skip known header text rows
				isKnownHeaderTextRow := false
				if programNameFromCell == "" && !isPermFileGlobal { // For non-Perm, if first cell (program name) is empty, likely not data. For Perm, programNameFromCell is 2nd cell.
					isKnownHeaderTextRow = true
				} else if programNameFromCell == "" && isPermFileGlobal && firstCellThTextForPerm == "" { // Both key cells empty for Perm
					isKnownHeaderTextRow = true
				}

				if !isKnownHeaderTextRow { // Only proceed if not already marked as header
					lowerProgramName := strings.ToLower(programNameFromCell)
					headerKeywords := []string{
						"специальность/образовательная программа", "направление подготовки / образовательная программа",
						"наименование образовательной программы",
						"образовательная программа бакалавриата", // Perm's <td> can be this for a header
					}
					for _, kw := range headerKeywords {
						if lowerProgramName == kw { // Exact match for these specific headers
							isKnownHeaderTextRow = true
							break
						}
					}
					// Check for "особой квоты" etc. only if it's Perm and the program name cell contains it
					if isPermFileGlobal && (strings.Contains(lowerProgramName, "особой квоты") || strings.Contains(lowerProgramName, "отдельной квоты") || strings.Contains(lowerProgramName, "целевой квоты")) {
						isKnownHeaderTextRow = true
					}

					if strings.HasPrefix(lowerProgramName, "итого") || strings.HasPrefix(lowerProgramName, "всего") {
						isKnownHeaderTextRow = true
					}
					if isPermFileGlobal && strings.ToLower(firstCellThTextForPerm) == "всего" { // Perm summary row by <th>
						isKnownHeaderTextRow = true
					}
				}

				if !isKnownHeaderTextRow {
					// Rule 4: Check expected cell count and extract data
					expectedCols := 0
					kcpColIdx, specColIdx, dedColIdx, targColIdx := 0, 0, 0, 0

					if isPermFileGlobal {
						expectedCols = 8 // 1 th + 7 td. We collected all of them in cellsContent.
						// cellsContent[0] is th text (firstCellThTextForPerm)
						// cellsContent[1] is program name td text (programNameFromCell)
						// KCP from cellsContent[2], Special from [3], Dedicated from [4], Target from [5]
						kcpColIdx, specColIdx, dedColIdx, targColIdx = 2, 3, 4, 5
					} else { // MSK, SPB, NN
						expectedCols = 7 // all td
						// cellsContent[0] is program name td text (programNameFromCell)
						// KCP from cellsContent[1]
						kcpColIdx = 1
						isNN := strings.Contains(filePath, "_nn_places.html")
						// isSpb := strings.Contains(filePath, "_spb_places.html") // Not needed for quota order if NN is handled

						if isNN {
							specColIdx, dedColIdx, targColIdx = 2, 3, 4 // Special, Dedicated, Target
						} else { // MSK, SPB (Perm is handled above)
							specColIdx, targColIdx, dedColIdx = 2, 3, 4 // Special, Target, Dedicated
						}
					}

					if len(cellsContent) == expectedCols {
						// All checks passed, this is a data row
						kcp := extractIntFromCell(cellsContent[kcpColIdx])
						specialQuota := extractIntFromCell(cellsContent[specColIdx])
						dedicatedQuota := extractIntFromCell(cellsContent[dedColIdx])
						targetQuota := extractIntFromCell(cellsContent[targColIdx])

						generalCapacity := kcp - specialQuota - targetQuota - dedicatedQuota
						if generalCapacity < 0 {
							generalCapacity = 0
						}

						normalizedName := normalizeHeadingName(programNameFromCell)
						if _, exists := capacitiesMap[normalizedName]; exists {
							log.Printf("Warning: Duplicate normalized program name in capacities file '%s': '%s' (from original '%s')", filePath, normalizedName, programNameFromCell)
						}
						capacitiesMap[normalizedName] = core.Capacities{
							Regular:        generalCapacity,
							TargetQuota:    targetQuota,
							DedicatedQuota: dedicatedQuota,
							SpecialQuota:   specialQuota,
						}
						originalNamesMap[normalizedName] = programNameFromCell
					}
				}
			}
		} // End of "if n.Data == tr"

		// Standard recursion for all nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	// Helper to extract cell texts from a given node (e.g., thead or tr) for specific cell tag (th or td)
	extractHeaders := func(parentNode *html.Node, cellTagName string) []string {
		var texts []string
		// Iterate direct children of parentNode (which is expected to be <thead> or <tr>)
		for c := parentNode.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == cellTagName {
				texts = append(texts, getTextContent(c))
			} else if cellTagName == "th" && c.Type == html.ElementNode && c.Data == "tr" { // Handle cases where <th> are inside <tr> inside <thead>
				for thCell := c.FirstChild; thCell != nil; thCell = thCell.NextSibling {
					if thCell.Type == html.ElementNode && thCell.Data == "th" {
						texts = append(texts, getTextContent(thCell))
					}
				}
			}
		}
		return texts
	}

	var findMainTable func(*html.Node, string) *html.Node
	findMainTable = func(n *html.Node, currentFilePath string) *html.Node {
		if n.Type == html.ElementNode && n.Data == "table" {
			isPermFile := strings.Contains(currentFilePath, "_perm_places.html")

			var tableThead *html.Node
			for child := n.FirstChild; child != nil; child = child.NextSibling {
				if child.Type == html.ElementNode && child.Data == "thead" {
					tableThead = child
					break
				}
			}

			if tableThead != nil {
				if isPermFile {
					// For Perm, data rows are children of thead.
					// The parsing function `f` will iterate its <tr> children and skip non-data rows.
					log.Printf("Info: Using thead as data row container for Perm file %s.", currentFilePath)
					return tableThead // Return thead itself
				}

				// Original logic for MSK/NN (they also have thead but data is in tbody)
				headerTexts := extractHeaders(tableThead, "th")
				if len(headerTexts) >= 3 &&
					(strings.Contains(headerTexts[0], "Направление подготовки") || strings.Contains(headerTexts[0], "Наименование образовательной программы") || strings.Contains(headerTexts[0], "Образовательная программа")) &&
					(strings.Contains(headerTexts[1], "КЦП") || strings.Contains(headerTexts[1], "Бюджетные места")) &&
					strings.Contains(headerTexts[2], "особое право") {
					for tbodyCandidate := tableThead.NextSibling; tbodyCandidate != nil; tbodyCandidate = tbodyCandidate.NextSibling {
						if tbodyCandidate.Type == html.ElementNode && tbodyCandidate.Data == "tbody" {
							log.Printf("Info: Found table body for %s using thead-based detection.", currentFilePath)
							return tbodyCandidate
						}
					}
					for child := n.FirstChild; child != nil; child = child.NextSibling {
						if child.Type == html.ElementNode && child.Data == "tbody" {
							log.Printf("Info: Found table body for %s using thead-based detection (fallback table child tbody).", currentFilePath)
							return child
						}
					}
				}
			}

			// Attempt 2: SPB style (tbody -> tr[0] -> td) or general fallback if no thead or thead didn't lead to tbody
			var tableTbody *html.Node
			for child := n.FirstChild; child != nil; child = child.NextSibling {
				if child.Type == html.ElementNode && child.Data == "tbody" {
					tableTbody = child
					break
				}
			}

			if tableTbody != nil {
				var firstTr *html.Node
				for trCandidate := tableTbody.FirstChild; trCandidate != nil; trCandidate = trCandidate.NextSibling {
					if trCandidate.Type == html.ElementNode && trCandidate.Data == "tr" {
						firstTr = trCandidate
						break
					}
				}

				if firstTr != nil {
					headerTexts := extractHeaders(firstTr, "td")
					if len(headerTexts) >= 3 &&
						strings.Contains(headerTexts[0], "Специальность/образовательная программа") &&
						(strings.Contains(headerTexts[1], "Распределение КЦП") || strings.Contains(headerTexts[1], "КЦП")) &&
						strings.Contains(headerTexts[2], "особое право") {
						log.Printf("Info: Found table body for %s using tbody/tr[0]/td-based detection (SPB style).", currentFilePath)
						return tableTbody
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if foundTbody := findMainTable(c, currentFilePath); foundTbody != nil {
				return foundTbody
			}
		}
		return nil
	}

	tableBody := findMainTable(doc, filePath)
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
