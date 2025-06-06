package main

import (
	"analabit/core/source/hse"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

var permHeadingToSource = make(map[string]hse.HttpHeadingSource)

// printPermSourcesListFunc prints a ready to paste func `permSourcesList` declaration
// see defs/hse.go for permSourcesList() definition
func printPermSourcesListFunc() {
	capacitiesPath := "resources/hse_perm_places.html"
	linksPath := "resources/hse_perm_links.html"

	capacitiesMap, capOriginalNames, err := parseCapacitiesHTML(capacitiesPath)
	if err != nil {
		log.Fatalf("Error parsing Perm capacities: %v", err)
	}

	linksMap, linkOriginalNames, err := parseLinksHTML(linksPath)
	if err != nil {
		log.Fatalf("Error parsing Perm links: %v", err)
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
	sort.Strings(headingNamesForSorting) // Sort by normalized name for consistent output

	// --- Revised Initial Linking (Exact + Containment) ---
	usedNormalizedLinkNamesInInitialPass := make(map[string]bool) // Tracks normalized link names used in this phase

	for _, normalizedCapName := range headingNamesForSorting { // Iterate over sorted normalized capacity names
		caps := capacitiesMap[normalizedCapName]
		var currentLinks ProgramLinks // MODIFIED: Corrected type to local ProgramLinks, not hse.ProgramLinks

		// 1. Attempt Exact Match
		// Ensure that the key exists in linkOriginalNames for logging, though linksMap is the source of truth for links.
		exactMatchOriginalLinkName := linkOriginalNames[normalizedCapName]
		if linkData, ok := linksMap[normalizedCapName]; ok {
			currentLinks = linkData
			usedNormalizedLinkNamesInInitialPass[normalizedCapName] = true // Mark as used by exact match
			if capOriginalNames[normalizedCapName] != "" && exactMatchOriginalLinkName != "" {
				// Optional: log exact match if detailed logging is desired.
				// log.Printf("Info: Perm capacity '%s' (orig: '%s') exactly matched link '%s' (orig: '%s').",
				// normalizedCapName, capOriginalNames[normalizedCapName], normalizedCapName, exactMatchOriginalLinkName)
			}
		}

		sourceEntry := hse.HttpHeadingSource{
			RCListURL:         validateUrl(currentLinks.RCURL),  // Corrected variable
			TQListURL:         validateUrl(currentLinks.TQURL),  // Corrected variable
			DQListURL:         validateUrl(currentLinks.DQURL),  // Corrected variable
			SQListURL:         validateUrl(currentLinks.SQURL),  // Corrected variable
			BListURL:          validateUrl(currentLinks.BVIURL), // Corrected variable
			HeadingCapacities: caps,
		}
		permHeadingToSource[normalizedCapName] = sourceEntry // Corrected map key
	}

	// --- Enhanced Fuzzy Matching Logic (adapted from msk.go) ---
	var capacitiesOriginalNamesWithoutExactLinks []string
	for normCapName := range capacitiesMap {
		entry := permHeadingToSource[normCapName]
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
		// Check if this normalized link name (from links file) has a corresponding capacity entry
		// AND that capacity entry actually uses these links (not links from a fuzzy match for example)
		if entry, exists := permHeadingToSource[normLinkName]; exists {
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
		for _, sourceEntry := range permHeadingToSource {
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
	processedCapOriginalNames := make(map[string]bool) // Tracks original capacity names that received links via fuzzy matching
	usedOriginalLinkNames := make(map[string]bool)     // Tracks original link names that were used in fuzzy matching

	if len(capacitiesOriginalNamesWithoutExactLinks) > 0 && len(linksOriginalNamesWithoutExactCaps) > 0 {
		availableLinksForFuzzy := make([]string, 0, len(linksOriginalNamesWithoutExactCaps))
		for _, name := range linksOriginalNamesWithoutExactCaps {
			availableLinksForFuzzy = append(availableLinksForFuzzy, name)
		}

		for _, originalCapName := range capacitiesOriginalNamesWithoutExactLinks {
			if len(availableLinksForFuzzy) == 0 {
				break // No more links to match against
			}

			matches := fuzzy.RankFindFold(originalCapName, availableLinksForFuzzy)
			if len(matches) > 0 {
				bestMatch := matches[0]
				originalLinkNameMatched := bestMatch.Target

				normalizedCapName := normalizeHeadingName(originalCapName)
				normalizedLinkNameMatched := normalizeHeadingName(originalLinkNameMatched)

				if linksData, ok := linksMap[normalizedLinkNameMatched]; ok {
					entryToUpdate, entryExists := permHeadingToSource[normalizedCapName]
					if !entryExists {
						log.Printf("Error: Perm Capacity name '%s' (normalized: '%s') not found in permHeadingToSource for fuzzy update. This should not happen.", originalCapName, normalizedCapName)
						continue
					}

					// Only update if the current links are empty (to avoid overwriting containment matches from nn/spb style if that was mixed in)
					// For Perm, following msk.go, this check might be redundant if no prior fuzzy logic was applied.
					// However, it's safer to keep it.
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
						permHeadingToSource[normalizedCapName] = entryToUpdate

						fuzzyMatchedPairsLog = append(fuzzyMatchedPairsLog, fmt.Sprintf("(Cap: '%s' ~ Link: '%s' [Dist: %d])", originalCapName, originalLinkNameMatched, bestMatch.Distance))
						processedCapOriginalNames[originalCapName] = true
						usedOriginalLinkNames[originalLinkNameMatched] = true

						// Remove the matched link name from availableLinksForFuzzy
						tempLinks := []string{}
						for _, ln := range availableLinksForFuzzy {
							if ln != originalLinkNameMatched {
								tempLinks = append(tempLinks, ln)
							}
						}
						availableLinksForFuzzy = tempLinks
					}
				} else {
					log.Printf("Warning: Perm Normalized link name '%s' (from original '%s') not found in linksMap during fuzzy matching. This is unexpected.", normalizedLinkNameMatched, originalLinkNameMatched)
				}
			}
		}
	}

	if len(fuzzyMatchedPairsLog) > 0 {
		log.Printf("Info: Applied fuzzy matches for Perm: [ %s ]", strings.Join(fuzzyMatchedPairsLog, ", "))
	}

	var finalCapacitiesWithoutLinks []string
	for _, originalCapName := range capacitiesOriginalNamesWithoutExactLinks {
		if !processedCapOriginalNames[originalCapName] {
			// Re-check if links are still empty, as some other logic might have populated them (though unlikely with msk.go structure)
			normalizedCapName := normalizeHeadingName(originalCapName)
			entry, exists := permHeadingToSource[normalizedCapName]
			if !exists || (entry.RCListURL == "" && entry.TQListURL == "" && entry.DQListURL == "" && entry.SQListURL == "" && entry.BListURL == "") {
				finalCapacitiesWithoutLinks = append(finalCapacitiesWithoutLinks, originalCapName)
			}
		}
	}
	if len(finalCapacitiesWithoutLinks) > 0 {
		sort.Strings(finalCapacitiesWithoutLinks)
		log.Printf("Info: Perm Programs in CAPACITIES list but still no LINKS (after fuzzy attempts): [ %s ]", strings.Join(finalCapacitiesWithoutLinks, ", "))
	}

	var finalLinksWithoutCapacities []string
	for _, originalLinkName := range linksOriginalNamesWithoutExactCaps {
		if !usedOriginalLinkNames[originalLinkName] {
			finalLinksWithoutCapacities = append(finalLinksWithoutCapacities, originalLinkName)
		}
	}
	if len(finalLinksWithoutCapacities) > 0 {
		sort.Strings(finalLinksWithoutCapacities)
		log.Printf("Info: Perm Programs in LINKS list but not matched to CAPACITIES (after fuzzy attempts): [ %s ]", strings.Join(finalLinksWithoutCapacities, ", "))
	}

	var sb strings.Builder
	sb.WriteString("// permSourcesList returns a list of HeadingSource for HSE Perm.\n") // MODIFIED: Joined to single line with escaped newline
	sb.WriteString("// Generated by codegen/hse/perm.go\n")
	sb.WriteString("func permSourcesList() []source.HeadingSource {\n") // MODIFIED: Joined to single line with escaped newline
	sb.WriteString("\treturn []source.HeadingSource{\n")                // MODIFIED: Joined to single line with escaped newline

	var mainOutput strings.Builder
	var zeroCapacityOutput strings.Builder
	var missingURLsOutput strings.Builder

	// Use headingNamesForSorting which is derived from capacitiesMap for consistent order
	for _, normalizedName := range headingNamesForSorting {
		sourceEntry, ok := permHeadingToSource[normalizedName]
		if !ok {
			log.Printf("Error: Perm Normalized name '%s' from capacitiesMap not found in permHeadingToSource during final print. Skipping.", normalizedName)
			continue
		}

		originalName, nameOk := capOriginalNames[normalizedName]
		if !nameOk {
			// Fallback: try to find original name from linksMap if it was a link-only entry (less likely with current logic)
			originalName, nameOk = linkOriginalNames[normalizedName]
			if !nameOk {
				originalName = normalizedName // Absolute fallback
				log.Printf("Warning: Perm Original name not found for normalized name '%s' during print. Using normalized name.", normalizedName)
			}
		}

		entryComment := fmt.Sprintf("\t\t// %s\n", strings.ReplaceAll(originalName, "`", "'")) // MODIFIED: Joined to single line with escaped newline
		entryCode := fmt.Sprintf("\t\t&%s,\n", printRvalueSource(&sourceEntry))                // MODIFIED: Joined to single line with escaped newline

		isZeroCapacity := sourceEntry.HeadingCapacities.Regular == 0 &&
			sourceEntry.HeadingCapacities.TargetQuota == 0 &&
			sourceEntry.HeadingCapacities.DedicatedQuota == 0 &&
			sourceEntry.HeadingCapacities.SpecialQuota == 0 &&
			// Also check if KCP itself was zero, as Regular could be 0 if quotas sum up to KCP
			(capacitiesMap[normalizedName].Regular+capacitiesMap[normalizedName].TargetQuota+capacitiesMap[normalizedName].DedicatedQuota+capacitiesMap[normalizedName].SpecialQuota) == 0

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
		sb.WriteString("\n\t\t// TODO The following Perm headings do not have capacities determined (or KCP is zero):\n\n") // MODIFIED: Joined to single line with escaped newline
		sb.WriteString(zeroCapacityOutput.String())
	}

	if missingURLsOutput.Len() > 0 {
		sb.WriteString("\n\t\t// TODO The following Perm headings do not have list URLs determined:\n\n") // MODIFIED: Joined to single line with escaped newline
		sb.WriteString(missingURLsOutput.String())
	}

	sb.WriteString("\t}\n}\n") // MODIFIED: Joined to single line with escaped newline
	fmt.Print(sb.String())
}
