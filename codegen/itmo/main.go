package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// ProgramInfo represents a program from the list page
type ProgramInfo struct {
	ID       int    `json:"id"`        // ID from URL (e.g. 2190)
	Code     string `json:"code"`      // Program code (e.g. "01.03.02")
	Name     string `json:"name"`      // Program name without quotes
	FullName string `json:"full_name"` // Full program name with code
	URL      string `json:"url"`       // Full URL to the program
	Capacity int    `json:"capacity"`  // КЦП capacity
}

// GeneratedSources represents the output for code generation
type GeneratedSources struct {
	Programs []ProgramInfo `json:"programs"`
	Count    int           `json:"count"`
}

var (
	// Regex to extract program code and name (handle both « and " quotes)
	listCodeRegex = regexp.MustCompile(`^(\d{2}\.\d{2}\.\d{2})\s*(?:и\s*\d{2}\.\d{2}\.\d{2})?\s*[«"]([^»"]+)[»"]`)
	// Regex to extract capacity number (handle cyrillic КЦП)
	listCapacityRegex = regexp.MustCompile(`(?:КЦП|ККП):\s*(\d+)`)
	// Regex to extract ID from URL
	urlIDRegex = regexp.MustCompile(`/(\d+)$`)
)

func generateSources(inputFile string) (*GeneratedSources, error) {
	content, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", inputFile, err)
	}

	doc, err := html.Parse(strings.NewReader(string(content)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	var programs []ProgramInfo
	seen := make(map[int]bool) // Track seen program IDs to avoid duplicates

	// Find all program cards using the HTML traversal
	findProgramCards(doc, &programs, seen)

	return &GeneratedSources{
		Programs: programs,
		Count:    len(programs),
	}, nil
}

// findProgramCards traverses the HTML tree to find program card links
func findProgramCards(n *html.Node, programs *[]ProgramInfo, seen map[int]bool) {
	if n.Type == html.ElementNode && n.Data == "a" {
		// Check if this is a program card
		if hasClass(n, "DirectionsList_card__5AVa5") {
			program := parseProgramCard(n)
			if program.ID > 0 {
				if !seen[program.ID] {
					seen[program.ID] = true
					*programs = append(*programs, program)
				} else {
					// Update existing program with better information
					for i := range *programs {
						if (*programs)[i].ID == program.ID {
							// Update if we have better capacity info
							if program.Capacity > 0 && (*programs)[i].Capacity == 0 {
								(*programs)[i].Capacity = program.Capacity
							}
							// Update if we have better name info
							if program.Name != "" && ((*programs)[i].Name == "" || len(program.Name) > len((*programs)[i].Name)) {
								(*programs)[i].Name = program.Name
								(*programs)[i].Code = program.Code
								(*programs)[i].FullName = program.FullName
							}
							break
						}
					}
				}
			}
		}
	}

	// Recursively traverse children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findProgramCards(c, programs, seen)
	}
}

// hasClass checks if an HTML node has a specific CSS class
func hasClass(n *html.Node, className string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" && strings.Contains(attr.Val, className) {
			return true
		}
	}
	return false
}

// parseProgramCard extracts program information from a card link node
func parseProgramCard(cardNode *html.Node) ProgramInfo {
	program := ProgramInfo{}

	// Extract href attribute
	for _, attr := range cardNode.Attr {
		if attr.Key == "href" {
			program.URL = attr.Val
			// Extract ID from URL
			if matches := urlIDRegex.FindStringSubmatch(attr.Val); len(matches) >= 2 {
				if id, err := strconv.Atoi(matches[1]); err == nil {
					program.ID = id
				}
			}
			break
		}
	}

	// Extract program text and capacity from child elements
	var programText, capacityText string
	extractCardContent(cardNode, &programText, &capacityText)

	// Parse program code and name
	if matches := listCodeRegex.FindStringSubmatch(programText); len(matches) >= 3 {
		program.Code = matches[1]
		program.Name = matches[2]
	} else {
		// Fallback: try to extract quotes manually
		// Handle both « » and " " quotes
		var start, end int = -1, -1

		for i, char := range programText {
			if char == '«' || char == '"' {
				start = i
				break
			}
		}

		if start >= 0 {
			for i := start + 1; i < len(programText); i++ {
				char := rune(programText[i])
				if char == '»' || char == '"' {
					end = i
					break
				}
			}
		}

		if start >= 0 && end > start {
			// Extract name between quotes
			name := programText[start+1 : end]
			// Handle UTF-8 properly
			if programText[start] == '«' {
				// For cyrillic quotes, we need to handle UTF-8
				runes := []rune(programText)
				for i, r := range runes {
					if r == '«' {
						for j := i + 1; j < len(runes); j++ {
							if runes[j] == '»' {
								program.Name = string(runes[i+1 : j])
								break
							}
						}
						break
					}
				}
			} else {
				program.Name = name
			}

			// Extract code from beginning
			codePart := strings.TrimSpace(programText[:start])
			codeRegex := regexp.MustCompile(`(\d{2}\.\d{2}\.\d{2})`)
			if matches := codeRegex.FindStringSubmatch(codePart); len(matches) >= 2 {
				program.Code = matches[1]
			}
		}
	}

	// If we couldn't parse properly, use the full text as name
	if program.Name == "" {
		program.Name = programText
	}

	program.FullName = programText

	// Parse capacity
	if matches := listCapacityRegex.FindStringSubmatch(capacityText); len(matches) >= 2 {
		if capacity, err := strconv.Atoi(matches[1]); err == nil {
			program.Capacity = capacity
		}
	}

	// Build full URL
	if program.URL != "" && !strings.HasPrefix(program.URL, "http") {
		program.URL = "https://abit.itmo.ru" + program.URL
	}

	return program
}

// extractCardContent extracts text content from the card's div elements
func extractCardContent(cardNode *html.Node, programText, capacityText *string) {
	var extractText func(*html.Node) string
	extractText = func(n *html.Node) string {
		if n.Type == html.TextNode {
			return n.Data
		}
		var text string
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			text += extractText(c)
		}
		return text
	}

	// Look for specific elements within the card
	var findContent func(*html.Node)
	findContent = func(n *html.Node) {
		if n.Type == html.ElementNode {
			text := strings.TrimSpace(extractText(n))
			if text != "" {
				if n.Data == "p" && *programText == "" {
					// First <p> tag contains program name
					*programText = text
				} else if strings.Contains(text, "КЦП:") || strings.Contains(text, "ККП:") {
					// Any element with capacity info
					*capacityText = text
				} else if n.Data == "div" && *programText == "" && !strings.Contains(text, "КЦП:") && !strings.Contains(text, "ККП:") {
					// Fallback: first div without capacity info
					*programText = text
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findContent(c)
		}
	}

	findContent(cardNode)
}

func generateSourceCode(sources *GeneratedSources) string {
	var builder strings.Builder

	builder.WriteString("// Generated ITMO HTTPHeadingSource definitions\n")
	builder.WriteString("// Copy these into /core/registry/itmo/itmo.go as a slice:\n\n")

	for i, program := range sources.Programs {
		builder.WriteString(fmt.Sprintf("\t\t&itmo.HTTPHeadingSource{\n"))
		builder.WriteString(fmt.Sprintf("\t\t\tURL: \"%s\",\n", program.URL))
		builder.WriteString(fmt.Sprintf("\t\t\tPrettyName: \"%s\",\n", escapeString(program.Name)))
		// Capacities will be parsed from individual program pages
		builder.WriteString(fmt.Sprintf("\t\t\tCapacities: core.Capacities{Regular: %d}, // fallback total КЦП\n", program.Capacity))
		builder.WriteString("\t\t}")

		// Add comma except for last item
		if i < len(sources.Programs)-1 {
			builder.WriteString(",")
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

func escapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}

func main() {
	generateMain()
}

func generateMain() {
	inputFile := "/home/yegor/Prestart/analabit/sample_data/itmo/itmo_lists.html"

	fmt.Println("Generating ITMO sources from program list...")

	sources, err := generateSources(inputFile)
	if err != nil {
		log.Fatalf("Error generating sources: %v", err)
	}

	fmt.Printf("Found %d programs\n", sources.Count)

	// Generate the Go source code for manual copy-paste
	code := generateSourceCode(sources)

	// Print to stdout for copy-paste
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("Copy the following slice into /core/registry/itmo/itmo.go:")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println(code)
	fmt.Println(strings.Repeat("=", 80))

	// Also write to file for reference
	outputFile := "generated_sources.txt"
	err = ioutil.WriteFile(outputFile, []byte(code), 0644)
	if err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	fmt.Printf("\nAlso saved to %s for reference\n", outputFile)

	// Print some examples
	fmt.Println("\nFirst few programs:")
	for i, program := range sources.Programs {
		if i >= 5 {
			break
		}
		fmt.Printf("  %s (%s) - ID: %d, Capacity: %d\n",
			program.Name, program.Code, program.ID, program.Capacity)
	}
}
