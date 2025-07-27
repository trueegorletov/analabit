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
	Name               string
	Capacities         core.Capacities
	RegularURLs        []string
	TargetQuotaURLs    []string
	DedicatedQuotaURLs []string
	SpecialQuotaURLs   []string
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

	// Output Go source code
	classifications := combineDataAndGenerate(capacities, links)
	generateGoSourceCode(classifications.Complete, classifications.Incomplete)
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

type HeadingClassification struct {
	Complete   []MephiHeading
	Incomplete []MephiHeading
}

func combineDataAndGenerate(capacities map[string]int, links map[string]map[core.Competition][]string) HeadingClassification {
	var complete []MephiHeading
	var incomplete []MephiHeading

	// Manual mapping for differently written heading names
	// Maps variations to canonical names found in capacities
	headingNameMapping := map[string]string{
		"Применение и эксплуатация автоматизированных систем специального назначения специалитет": "Применение и эксплуатация автом. систем специального назначения",
		"Системный анализ и управление": "Системный анализ  и управление",
		// Add more mappings here as needed
	}

	// Helper function to get canonical heading name
	getCanonicalName := func(name string) string {
		if canonical, exists := headingNameMapping[name]; exists {
			return canonical
		}
		return name
	}

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
		canonicalName := getCanonicalName(headingName)
		totalCapacity := capacities[headingName]
		if totalCapacity == 0 {
			totalCapacity = capacities[canonicalName]
		}
		headingLinks := links[headingName]
		if headingLinks == nil {
			headingLinks = links[canonicalName]
		}
		if headingLinks == nil {
			headingLinks = make(map[core.Competition][]string)
		}

		hasLists := len(headingLinks[core.CompetitionRegular]) > 0 || len(headingLinks[core.CompetitionTargetQuota]) > 0 || len(headingLinks[core.CompetitionDedicatedQuota]) > 0 || len(headingLinks[core.CompetitionSpecialQuota]) > 0
		hasRegular := len(headingLinks[core.CompetitionRegular]) > 0

		var caps core.Capacities
		if totalCapacity > 0 {
			caps = calculateCapacities(totalCapacity)
		} else {
			caps = core.Capacities{}
		}

		heading := MephiHeading{
			Name:               headingName,
			Capacities:         caps,
			RegularURLs:        makeAbsoluteURLs(headingLinks[core.CompetitionRegular]),
			TargetQuotaURLs:    makeAbsoluteURLs(headingLinks[core.CompetitionTargetQuota]),
			DedicatedQuotaURLs: makeAbsoluteURLs(headingLinks[core.CompetitionDedicatedQuota]),
			SpecialQuotaURLs:   makeAbsoluteURLs(headingLinks[core.CompetitionSpecialQuota]),
		}

		if totalCapacity > 0 && hasRegular {
			complete = append(complete, heading)
		} else if (totalCapacity > 0 && !hasRegular) || (totalCapacity == 0 && hasLists) {
			incomplete = append(incomplete, heading)
		}

		processed[headingName] = true
	}

	return HeadingClassification{
		Complete:   complete,
		Incomplete: incomplete,
	}

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

func generateGoSourceCode(complete []MephiHeading, incomplete []MephiHeading) {
	fmt.Println("// Generated MEPhI HTTPHeadingSource structs")
	fmt.Println("// Copy and paste the following into core/registry/mephi/mephi.go")
	fmt.Println()

	for i, heading := range complete {
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

	if len(incomplete) > 0 {
		fmt.Println("\n// NOT COMPLETE")
		for _, heading := range incomplete {
			fmt.Printf("// &mephi.HTTPHeadingSource{\n")
			fmt.Printf("// \tHeadingName: %q,\n", heading.Name)
			fmt.Printf("// \tCapacities: core.Capacities{\n")
			fmt.Printf("// \t\tRegular:        %d,\n", heading.Capacities.Regular)
			fmt.Printf("// \t\tTargetQuota:    %d,\n", heading.Capacities.TargetQuota)
			fmt.Printf("// \t\tDedicatedQuota: %d,\n", heading.Capacities.DedicatedQuota)
			fmt.Printf("// \t\tSpecialQuota:   %d,\n", heading.Capacities.SpecialQuota)
			fmt.Printf("// \t},\n")

			if len(heading.RegularURLs) > 0 {
				fmt.Printf("// \tRegularURLs: []string{\n")
				for _, url := range heading.RegularURLs {
					fmt.Printf("// \t\t%q,\n", url)
				}
				fmt.Printf("// \t},\n")
			}

			if len(heading.TargetQuotaURLs) > 0 {
				fmt.Printf("// \tTargetQuotaURLs: []string{\n")
				for _, url := range heading.TargetQuotaURLs {
					fmt.Printf("// \t\t%q,\n", url)
				}
				fmt.Printf("// \t},\n")
			}

			if len(heading.DedicatedQuotaURLs) > 0 {
				fmt.Printf("// \tDedicatedQuotaURLs: []string{\n")
				for _, url := range heading.DedicatedQuotaURLs {
					fmt.Printf("// \t\t%q,\n", url)
				}
				fmt.Printf("// \t},\n")
			}

			if len(heading.SpecialQuotaURLs) > 0 {
				fmt.Printf("// \tSpecialQuotaURLs: []string{\n")
				for _, url := range heading.SpecialQuotaURLs {
					fmt.Printf("// \t\t%q,\n", url)
				}
				fmt.Printf("// \t},\n")
			}

			fmt.Printf("// }\n")
		}
	}

	fmt.Println()
	fmt.Printf("\n// Total complete headings generated: %d", len(complete))
	fmt.Printf("\n// Total incomplete headings: %d\n", len(incomplete))

}
