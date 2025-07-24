package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source/mephi"
)

const baseURL = "https://org.mephi.ru"

type MephiHeading struct {
	Name                   string
	Capacities             core.Capacities
	RegularURLs            []string
	TargetQuotaURLs        []string
	DedicatedQuotaURLs     []string
	SpecialQuotaURLs       []string
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <capacities-html-file> <links-html-file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s ../../sample_data/mephi/mephi_heading_names_and_capacities_registry_body.html ../../sample_data/mephi/mephi_links_to_lists_registry_body.html\n", os.Args[0])
		os.Exit(1)
	}

	capacitiesFile := os.Args[1]
	linksFile := os.Args[2]

	// Parse capacities from the first file
	capacities, err := parseCapacitiesFromFile(capacitiesFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing capacities: %v\n", err)
		os.Exit(1)
	}

	// Parse links from the second file
	links, err := parseLinksFromFile(linksFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing links: %v\n", err)
		os.Exit(1)
	}

	// Combine data and generate HTTPHeadingSource structs
	headings := combineDataAndGenerate(capacities, links)

	// Output Go source code
	generateGoSourceCode(headings)
}

func parseCapacitiesFromFile(filename string) (map[string]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open capacities file: %w", err)
	}
	defer file.Close()

	doc, err := html.Parse(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	return mephi.ParseCapacitiesRegistry(doc)
}

func parseLinksFromFile(filename string) (map[string]map[core.Competition][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open links file: %w", err)
	}
	defer file.Close()

	doc, err := html.Parse(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	return mephi.ParseListLinksRegistry(doc)
}

func combineDataAndGenerate(capacities map[string]int, links map[string]map[core.Competition][]string) []MephiHeading {
	var headings []MephiHeading

	// Create a map to track which headings we've processed
	processed := make(map[string]bool)

	// Process all headings that have either capacities or links
	allHeadings := make(map[string]bool)
	for name := range capacities {
		allHeadings[name] = true
	}
	for name := range links {
		allHeadings[name] = true
	}

	for headingName := range allHeadings {
		if processed[headingName] {
			continue
		}

		// Get total capacity for this heading
		totalCapacity := capacities[headingName]
		if totalCapacity == 0 {
			fmt.Printf("// Warning: No capacity found for heading: %s\n", headingName)
			continue
		}

		// Calculate capacities using 10% rule
		caps := calculateCapacities(totalCapacity)

		// Get URLs for this heading
		headingLinks := links[headingName]
		if headingLinks == nil {
			headingLinks = make(map[core.Competition][]string)
		}

		// Convert relative URLs to absolute URLs
		heading := MephiHeading{
			Name:                   headingName,
			Capacities:             caps,
			RegularURLs:            makeAbsoluteURLs(headingLinks[core.CompetitionRegular]),
			TargetQuotaURLs:        makeAbsoluteURLs(headingLinks[core.CompetitionTargetQuota]),
			DedicatedQuotaURLs:     makeAbsoluteURLs(headingLinks[core.CompetitionDedicatedQuota]),
			SpecialQuotaURLs:       makeAbsoluteURLs(headingLinks[core.CompetitionSpecialQuota]),
		}

		headings = append(headings, heading)
		processed[headingName] = true

		fmt.Printf("// Generated heading: %s (Total: %d, Regular: %d, Target: %d, Dedicated: %d, Special: %d)\n",
			headingName, totalCapacity, caps.Regular, caps.TargetQuota, caps.DedicatedQuota, caps.SpecialQuota)
	}

	return headings
}

func calculateCapacities(total int) core.Capacities {
	// Use 10% rule for quotas as specified in the implementation plan
	targetQuota := total / 10
	dedicatedQuota := total / 10
	specialQuota := total / 10
	regular := total - targetQuota - dedicatedQuota - specialQuota

	// Ensure regular is not negative
	if regular < 0 {
		regular = 0
	}

	return core.Capacities{
		Regular:        regular,
		TargetQuota:    targetQuota,
		DedicatedQuota: dedicatedQuota,
		SpecialQuota:   specialQuota,
	}
}

func makeAbsoluteURLs(urls []string) []string {
	var absoluteURLs []string
	for _, url := range urls {
		if url != "" {
			// Check if URL is already absolute
			if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
				absoluteURLs = append(absoluteURLs, url)
			} else if strings.HasPrefix(url, "/") {
				// Relative URL starting with /
				absoluteURLs = append(absoluteURLs, baseURL+url)
			} else {
				// Assume it's already absolute or handle as needed
				absoluteURLs = append(absoluteURLs, url)
			}
		}
	}
	return absoluteURLs
}

func generateGoSourceCode(headings []MephiHeading) {
	fmt.Println("// Generated MEPhI HTTPHeadingSource structs")
	fmt.Println("// Copy and paste the following into core/registry/mephi/mephi.go")
	fmt.Println()

	for i, heading := range headings {
		if i > 0 {
			fmt.Println(",")
		}

		fmt.Printf("\t\t&mephi.HTTPHeadingSource{\n")
		fmt.Printf("\t\t\tHeadingName: %q,\n", heading.Name)
		fmt.Printf("\t\t\tCapacities: core.Capacities{\n")
		fmt.Printf("\t\t\t\tRegular:        %d,\n", heading.Capacities.Regular)
		fmt.Printf("\t\t\t\tTargetQuota:    %d,\n", heading.Capacities.TargetQuota)
		fmt.Printf("\t\t\t\tDedicatedQuota: %d,\n", heading.Capacities.DedicatedQuota)
		fmt.Printf("\t\t\t\tSpecialQuota:   %d,\n", heading.Capacities.SpecialQuota)
		fmt.Printf("\t\t\t},\n")

		if len(heading.RegularURLs) > 0 {
			fmt.Printf("\t\t\tRegularURLs: []string{\n")
			for _, url := range heading.RegularURLs {
				fmt.Printf("\t\t\t\t%q,\n", url)
			}
			fmt.Printf("\t\t\t},\n")
		}

		if len(heading.TargetQuotaURLs) > 0 {
			fmt.Printf("\t\t\tTargetQuotaURLs: []string{\n")
			for _, url := range heading.TargetQuotaURLs {
				fmt.Printf("\t\t\t\t%q,\n", url)
			}
			fmt.Printf("\t\t\t},\n")
		}

		if len(heading.DedicatedQuotaURLs) > 0 {
			fmt.Printf("\t\t\tDedicatedQuotaURLs: []string{\n")
			for _, url := range heading.DedicatedQuotaURLs {
				fmt.Printf("\t\t\t\t%q,\n", url)
			}
			fmt.Printf("\t\t\t},\n")
		}

		if len(heading.SpecialQuotaURLs) > 0 {
			fmt.Printf("\t\t\tSpecialQuotaURLs: []string{\n")
			for _, url := range heading.SpecialQuotaURLs {
				fmt.Printf("\t\t\t\t%q,\n", url)
			}
			fmt.Printf("\t\t\t},\n")
		}

		fmt.Printf("\t\t}")
	}

	fmt.Println()
	fmt.Printf("\n// Total headings generated: %d\n", len(headings))
}