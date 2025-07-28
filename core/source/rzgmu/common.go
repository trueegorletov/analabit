// Package rzgmu provides support for loading RZGMU (Рязанский государственный медицинский университет) admission data.
// RZGMU admission lists are provided in HTML format across four unified pages by competition type.
package rzgmu

import (
	"regexp"
	"strings"

	"github.com/trueegorletov/analabit/core"
)

// Page URLs for different competition types
var pageURLs = []struct {
	URL         string
	Competition core.Competition
}{
	{
		URL:         "https://www.rzgmu.ru/rzgmu_rang/files/ranklist_3-1.php",
		Competition: core.CompetitionRegular, // Regular + BVI
	},
	{
		URL:         "https://www.rzgmu.ru/rzgmu_rang/files/ranklist_3-2.php",
		Competition: core.CompetitionSpecialQuota,
	},
	{
		URL:         "https://www.rzgmu.ru/rzgmu_rang/files/ranklist_3-5.php",
		Competition: core.CompetitionDedicatedQuota,
	},
	{
		URL:         "https://www.rzgmu.ru/rzgmu_rang/files/ranklist_3-11.php",
		Competition: core.CompetitionTargetQuota,
	},
}

// Program names as they appear in the HTML
var programNames = []string{
	"Лечебное дело",
	"Стоматология",
	"Педиатрия",
	"Медико-профилактическое дело",
	"Фармация",
	"Клиническая психология",
}

// Regular expressions for parsing
var (
	// Regex to normalize program names (remove any parenthetical suffix)
	programNameRegex = regexp.MustCompile(`^(.+?)\s*\([^)]*\)\s*`)
)

// isBVI checks if the score cell indicates "Без вступительных испытаний" (BVI)
func isBVI(scoreText string) bool {
	return strings.Contains(scoreText, "БВИ") || strings.Contains(scoreText, "Без вступительных")
}

// normalizeProgram removes any parenthetical suffix from program name for matching and converts to lowercase
func normalizeProgram(name string) string {
	result := name
	if matches := programNameRegex.FindStringSubmatch(name); len(matches) > 1 {
		result = strings.TrimSpace(matches[1])
	} else {
		result = strings.TrimSpace(name)
	}
	return strings.ToLower(result)
}

// matchesProgram checks if the heading text matches the target program name (case-insensitive)
func matchesProgram(headingText, targetProgram string) bool {
	normalizedHeading := normalizeProgram(headingText)
	normalizedTarget := strings.ToLower(strings.TrimSpace(targetProgram))
	return normalizedHeading == normalizedTarget
}