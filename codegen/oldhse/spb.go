package main

import (
	"analabit/core/source/oldhse"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

var spbHeadingToSource = make(map[string]oldhse.HttpHeadingSource)

// printSpbSourcesListFunc prints a ready to paste func `spbSourcesList` declaration
// see defs/hse.go for spbSourcesList() definition
func printSpbSourcesListFunc() {
	capacitiesPath := "resources/hse_spb_places.html"
	linksPath := "resources/hse_spb_links.html"

	capacitiesMap, capOriginalNames, err := parseCapacitiesHTML(capacitiesPath)
	if err != nil {
		log.Fatalf("Error parsing SPB capacities: %v", err)
	}

	linksMap, linkOriginalNames, err := parseLinksHTML(linksPath)
	if err != nil {
		log.Fatalf("Error parsing SPB links: %v", err)
	}

	validateUrl := func(u string) string {
		if u == "" || u == "." || u == "/" || u == "#" || strings.TrimSpace(u) == "&nbsp;" {
			return ""
		}
		if strings.HasPrefix(u, "/storage/public_report_2024") {
			return "https://enrol.hse.ru" + u
		}
		// Handle cases where the URL might already be absolute but from a different domain (e.g. http://enrol.hse.ru)
		if strings.HasPrefix(u, "http://enrol.hse.ru/storage/public_report_2024") {
			return strings.Replace(u, "http://enrol.hse.ru", "https://enrol.hse.ru", 1)
		}
		return u
	}

	var headingNamesForSorting []string
	for normalizedName := range capacitiesMap {
		headingNamesForSorting = append(headingNamesForSorting, normalizedName)
	}
	sort.Strings(headingNamesForSorting)

	for _, normalizedName := range headingNamesForSorting {
		caps := capacitiesMap[normalizedName]
		links, linksOk := linksMap[normalizedName]
		if !linksOk {
			// Attempt to find a link by checking if the capacity name contains the link name
			foundFuzzyLink := false
			for linkNormName, linkData := range linksMap {
				if strings.Contains(normalizedName, linkNormName) {
					links = linkData
					linksOk = true
					log.Printf("Info: SPB capacity name '%s' (original: '%s') did not have exact link match, but fuzzy matched to link name '%s' (original: '%s') by containment.", normalizedName, capOriginalNames[normalizedName], linkNormName, linkOriginalNames[linkNormName])
					foundFuzzyLink = true
					break
				}
			}
			if !foundFuzzyLink {
				links = ProgramLinks{} // Default to empty links
			}
		}

		sourceEntry := oldhse.HttpHeadingSource{
			RCListURL:  validateUrl(links.RCURL),
			TQListURL:  validateUrl(links.TQURL),
			DQListURL:  validateUrl(links.DQURL),
			SQListURL:  validateUrl(links.SQURL),
			BListURL:   validateUrl(links.BVIURL),
			Capacities: caps,
		}
		spbHeadingToSource[normalizedName] = sourceEntry
	}

	var capacitiesOriginalNamesWithoutExactLinks []string
	for normCapName := range capacitiesMap {
		entry := spbHeadingToSource[normCapName]
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
	for normLinkName, origLinkName := range linkOriginalNames {
		var isDirectlyMatched bool
		if _, exists := spbHeadingToSource[normLinkName]; exists {
			entry := spbHeadingToSource[normLinkName]
			linksFromThisOrigName := linksMap[normLinkName]
			if entry.RCListURL == validateUrl(linksFromThisOrigName.RCURL) &&
				entry.TQListURL == validateUrl(linksFromThisOrigName.TQURL) &&
				entry.DQListURL == validateUrl(linksFromThisOrigName.DQURL) &&
				entry.SQListURL == validateUrl(linksFromThisOrigName.SQURL) &&
				entry.BListURL == validateUrl(linksFromThisOrigName.BVIURL) {
				isDirectlyMatched = true
			}
		}

		isUsedByAnyCapacity := false
		linksFromThisOrigName := linksMap[normLinkName]
		for _, sourceEntry := range spbHeadingToSource {
			if (linksFromThisOrigName.RCURL != "" && sourceEntry.RCListURL != "" && sourceEntry.RCListURL == validateUrl(linksFromThisOrigName.RCURL)) ||
				(linksFromThisOrigName.TQURL != "" && sourceEntry.TQListURL != "" && sourceEntry.TQListURL == validateUrl(linksFromThisOrigName.TQURL)) ||
				(linksFromThisOrigName.DQURL != "" && sourceEntry.DQListURL != "" && sourceEntry.DQListURL == validateUrl(linksFromThisOrigName.DQURL)) ||
				(linksFromThisOrigName.SQURL != "" && sourceEntry.SQListURL != "" && sourceEntry.SQListURL == validateUrl(linksFromThisOrigName.SQURL)) ||
				(linksFromThisOrigName.BVIURL != "" && sourceEntry.BListURL != "" && sourceEntry.BListURL == validateUrl(linksFromThisOrigName.BVIURL)) {
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
					entryToUpdate, entryExists := spbHeadingToSource[normalizedCapName]
					if !entryExists {
						log.Printf("Error: SPB Capacity name '%s' (normalized: '%s') not found in spbHeadingToSource for fuzzy update. This should not happen.", originalCapName, normalizedCapName)
						continue
					}

					canUpdate := entryToUpdate.RCListURL == "" &&
						entryToUpdate.TQListURL == "" &&
						entryToUpdate.DQListURL == "" &&
						entryToUpdate.SQListURL == "" &&
						entryToUpdate.BListURL == ""

					if canUpdate {
						entryToUpdate.RCListURL = validateUrl(linksData.RCURL)
						entryToUpdate.TQListURL = validateUrl(linksData.TQURL)
						entryToUpdate.DQListURL = validateUrl(linksData.DQURL)
						entryToUpdate.SQListURL = validateUrl(linksData.SQURL)
						entryToUpdate.BListURL = validateUrl(linksData.BVIURL)
						spbHeadingToSource[normalizedCapName] = entryToUpdate

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
					}
				} else {
					log.Printf("Warning: SPB Normalized link name '%s' (from original '%s') not found in linksMap during fuzzy matching. This is unexpected.", normalizedLinkNameMatched, originalLinkNameMatched)
				}
			}
		}
	}

	if len(fuzzyMatchedPairsLog) > 0 {
		log.Printf("Info: Applied fuzzy matches for SPB: [ %s ]", strings.Join(fuzzyMatchedPairsLog, ", "))
	}

	var finalCapacitiesWithoutLinks []string
	for _, originalCapName := range capacitiesOriginalNamesWithoutExactLinks {
		if !processedCapOriginalNames[originalCapName] {
			normalizedCapName := normalizeHeadingName(originalCapName)
			entry, exists := spbHeadingToSource[normalizedCapName]
			if !exists || (entry.RCListURL == "" && entry.TQListURL == "" && entry.DQListURL == "" && entry.SQListURL == "" && entry.BListURL == "") {
				finalCapacitiesWithoutLinks = append(finalCapacitiesWithoutLinks, originalCapName)
			}
		}
	}
	if len(finalCapacitiesWithoutLinks) > 0 {
		sort.Strings(finalCapacitiesWithoutLinks)
		log.Printf("Info: SPB Programs in CAPACITIES list but still no LINKS (after fuzzy attempts): [ %s ]", strings.Join(finalCapacitiesWithoutLinks, ", "))
	}

	var finalLinksWithoutCapacities []string
	for _, originalLinkName := range linksOriginalNamesWithoutExactCaps {
		if !usedOriginalLinkNames[originalLinkName] {
			finalLinksWithoutCapacities = append(finalLinksWithoutCapacities, originalLinkName)
		}
	}
	if len(finalLinksWithoutCapacities) > 0 {
		sort.Strings(finalLinksWithoutCapacities)
		log.Printf("Info: SPB Programs in LINKS list but not matched to CAPACITIES (after fuzzy attempts): [ %s ]", strings.Join(finalLinksWithoutCapacities, ", "))
	}

	var sb strings.Builder
	sb.WriteString("// spbSourcesList returns a list of HeadingSource for HSE Saint Petersburg.\n")
	sb.WriteString("// Generated by tools/hse_generate_lists.go\n")
	sb.WriteString("func spbSourcesList() []source.HeadingSource {\n")
	sb.WriteString("\treturn []source.HeadingSource{\n")

	var mainOutput strings.Builder
	var zeroCapacityOutput strings.Builder
	var missingURLsOutput strings.Builder

	for _, normalizedName := range headingNamesForSorting {
		sourceEntry, ok := spbHeadingToSource[normalizedName]
		if !ok {
			log.Printf("Error: SPB Normalized name '%s' from capacitiesMap not found in spbHeadingToSource during final print. Skipping.", normalizedName)
			continue
		}

		originalName, nameOk := capOriginalNames[normalizedName]
		if !nameOk {
			originalName, nameOk = linkOriginalNames[normalizedName]
			if !nameOk {
				originalName = normalizedName
				log.Printf("Warning: SPB Original name not found for normalized name '%s' during print. Using normalized name.", normalizedName)
			}
		}

		entryComment := fmt.Sprintf("\t\t// %s\n", strings.ReplaceAll(originalName, "`", "'"))
		entryCode := fmt.Sprintf("\t\t&%s,\n", printRvalueSource(&sourceEntry))

		isZeroCapacity := sourceEntry.Capacities.Regular == 0 &&
			sourceEntry.Capacities.TargetQuota == 0 &&
			sourceEntry.Capacities.DedicatedQuota == 0 &&
			sourceEntry.Capacities.SpecialQuota == 0

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
		sb.WriteString("\n\t\t// TODO The following SPB headings do not have capacities determined:\n\n")
		sb.WriteString(zeroCapacityOutput.String())
	}

	if missingURLsOutput.Len() > 0 {
		sb.WriteString("\n\t\t// TODO The following SPB headings do not have list URLs determined:\n\n")
		sb.WriteString(missingURLsOutput.String())
	}

	sb.WriteString("\t}\n}\n")
	fmt.Print(sb.String())
}
