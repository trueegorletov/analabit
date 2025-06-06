package main

import (
	"analabit/core/source/hse"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy" // Added for fuzzy matching
)

var mskHeadingToSource = make(map[string]hse.HttpHeadingSource)

// printMskSourcesList prints a ready to paste func `mskSourcesList` declaration
// see defs/hse.go for mskSourcesList() definition
func printMskSourcesListFunc() {
	capacitiesPath := "resources/hse_msk_places.html"
	linksPath := "resources/hse_msk_links.html"

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
			RCListURL:         validateUrl(links.RCURL),  // Changed to MustParseURL
			TQListURL:         validateUrl(links.TQURL),  // Changed to MustParseURL
			DQListURL:         validateUrl(links.DQURL),  // Changed to MustParseURL
			SQListURL:         validateUrl(links.SQURL),  // Changed to MustParseURL
			BListURL:          validateUrl(links.BVIURL), // Changed to MustParseURL
			HeadingCapacities: caps,
		}
		mskHeadingToSource[normalizedName] = sourceEntry
	}

	// --- Enhanced Fuzzy Matching Logic ---
	var capacitiesOriginalNamesWithoutExactLinks []string
	for normCapName := range capacitiesMap {
		entry := mskHeadingToSource[normCapName]
		isMissingLinks := entry.RCListURL == "" &&
			entry.TQListURL == "" &&
			entry.DQListURL == "" &&
			entry.SQListURL == "" &&
			entry.BListURL == ""
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
			if (linksFromThisOrigName.RCURL != "" && sourceEntry.RCListURL != "" && sourceEntry.RCListURL == linksFromThisOrigName.RCURL) || // Changed to MustParseURL
				(linksFromThisOrigName.TQURL != "" && sourceEntry.TQListURL != "" && sourceEntry.TQListURL == linksFromThisOrigName.TQURL) || // Changed to MustParseURL
				(linksFromThisOrigName.DQURL != "" && sourceEntry.DQListURL != "" && sourceEntry.DQListURL == linksFromThisOrigName.DQURL) || // Changed to MustParseURL
				(linksFromThisOrigName.SQURL != "" && sourceEntry.SQListURL != "" && sourceEntry.SQListURL == linksFromThisOrigName.SQURL) || // Changed to MustParseURL
				(linksFromThisOrigName.BVIURL != "" && sourceEntry.BListURL != "" && sourceEntry.BListURL == linksFromThisOrigName.BVIURL) { // Changed to MustParseURL
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

					entryToUpdate.RCListURL = validateUrl(linksData.RCURL) // Changed to MustParseURL
					entryToUpdate.TQListURL = validateUrl(linksData.TQURL) // Changed to MustParseURL
					entryToUpdate.DQListURL = validateUrl(linksData.DQURL) // Changed to MustParseURL
					entryToUpdate.SQListURL = validateUrl(linksData.SQURL) // Changed to MustParseURL
					entryToUpdate.BListURL = validateUrl(linksData.BVIURL) // Changed to MustParseURL
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
	sb.WriteString("// mskSourcesList returns a list of HeadingSource for HSE Moscow.\n")
	sb.WriteString("// Generated by codegen/hse/msk.go\n")
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

		isZeroCapacity := sourceEntry.HeadingCapacities.Regular == 0 &&
			sourceEntry.HeadingCapacities.TargetQuota == 0 &&
			sourceEntry.HeadingCapacities.DedicatedQuota == 0 &&
			sourceEntry.HeadingCapacities.SpecialQuota == 0

		isMissingAllURLs := sourceEntry.RCListURL == "" &&
			sourceEntry.TQListURL == "" &&
			sourceEntry.DQListURL == "" &&
			sourceEntry.SQListURL == "" &&
			sourceEntry.BListURL == ""

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
