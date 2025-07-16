// Package rsmu provides support for loading RSMU (Russian National Research Medical University) admission data.
// RSMU admission lists are provided in JSON format with target quota lists and individual applicant data.
package rsmu

import (
	"strings"

	"github.com/trueegorletov/analabit/core"
)

// TargetQuota represents a single entry in the target quota list JSON
type TargetQuota struct {
	Title          string `json:"title"`          // Title of the target quota
	File           string `json:"file"`           // Filename of the individual list JSON
	ShowDiplomaAvg bool   `json:"showDiplomaAvg"` // Whether to show diploma average
}

// Offer represents a single offer in the target quota list
type Offer struct {
	Code  string `json:"code"`  // Offer code
	Title string `json:"title"` // Offer title
	Plan  int    `json:"plan"`  // Number of places for this offer
}

// IndividualList represents the structure of an individual list JSON file
type IndividualList struct {
	Program               string           `json:"program"`      // Program name
	Type                  string           `json:"type"`         // Competition type
	Plan                  int              `json:"plan"`         // Number of places available
	Count                 int              `json:"count"`        // Total number of applicants
	Exams                 []string         `json:"exams"`        // List of required exams
	Achievements          []string         `json:"achievements"` // List of achievements
	Offers                []Offer          `json:"offers"`       // List of offers
	Applicants            []Applicant      `json:"applicants"`   // List of applicants
	SourceCompetitionType core.Competition `json:"-"`            // Set based on source URL, not serialized
}

// Applicant represents a single applicant in the individual list
type Applicant struct {
	Order                      int         `json:"order"`                      // Position in the ranking
	Title                      string      `json:"title"`                      // Student ID or identifier
	Total                      int         `json:"total"`                      // Total score
	Priority                   int         `json:"priority"`                   // Priority of application
	DiplomaAvg                 float64     `json:"diplomaAvg"`                 // Diploma average score
	Original                   bool        `json:"original"`                   // Whether original documents submitted
	Paid                       bool        `json:"paid"`                       // Whether paid education
	Approval                   bool        `json:"approval"`                   // Whether approved
	Published                  bool        `json:"published"`                  // Whether published
	Contract                   bool        `json:"contract"`                   // Whether contract signed
	Right9                     bool        `json:"right9"`                     // Special right 9
	Right10                    bool        `json:"right10"`                    // Special right 10
	NoExam                     bool        `json:"noExam"`                     // Whether exempt from exams
	Highest                    bool        `json:"highest"`                    // Whether highest priority
	AchievementScore           int         `json:"achievementScore"`           // Achievement score
	AchievementScoreGeneral    int         `json:"achievementScoreGeneral"`    // General achievement score
	Comment                    string      `json:"comment"`                    // Additional comments
	State                      string      `json:"state"`                      // Application state
	OrgCode                    interface{} `json:"orgCode"`                    // Organization code (can be null)
	Rejected                   bool        `json:"rejected"`                   // Whether rejected
	RejectionReason            string      `json:"rejectionReason"`            // Reason for rejection
	Exams                      []int       `json:"exams"`                      // Exam scores (can be null)
	Achievements               interface{} `json:"achievements"`               // Achievements (can be null)
}

// mapCompetitionType maps RSMU competition type string to core.Competition
func mapCompetitionType(competitionType string) core.Competition {
	switch strings.ToLower(strings.TrimSpace(competitionType)) {
	case "общий конкурс":
		return core.CompetitionRegular
	case "целевая квота", "целевой прием":
		return core.CompetitionTargetQuota
	case "особая квота":
		return core.CompetitionSpecialQuota
	case "отдельная квота":
		return core.CompetitionDedicatedQuota
	default:
		return core.CompetitionRegular // Default to regular competition
	}
}