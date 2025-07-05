package core

import (
	"log/slog"
	"sort"
)

// applicationsNormalizer takes a set of applications for a heading and normalizes them.
// It ensures that each student is represented only by their best application and that
// rating places are re-numbered consecutively according to official precedence rules.
type applicationsNormalizer struct {
	apps []*Application
}

// newApplicationsNormalizer creates a normalizer for a given slice of applications.
func newApplicationsNormalizer(apps []*Application) *applicationsNormalizer {
	return &applicationsNormalizer{apps: apps}
}

// normalize performs the full normalization process.
func (n *applicationsNormalizer) normalize() []*Application {
	if len(n.apps) == 0 {
		return n.apps
	}

	// Step 1: Group applications by student ID to handle multiple applications from a single student.
	groupedApps := make(map[string][]*Application)
	for _, app := range n.apps {
		groupedApps[app.StudentID()] = append(groupedApps[app.StudentID()], app)
	}

	// Step 2: For each student, select only the application with the highest competition precedence.
	var uniqueApps []*Application
	for studentID, studentApps := range groupedApps {
		if len(studentApps) > 1 {
			slog.Warn("Multiple applications for same student and heading, retaining best", "studentID", studentID)
		}

		bestApp := studentApps[0]
		for _, app := range studentApps {
			if competitionPrecedence(app.CompetitionType()) > competitionPrecedence(bestApp.CompetitionType()) {
				bestApp = app
			}
		}
		uniqueApps = append(uniqueApps, bestApp)
	}

	// Step 3: Sort the unique applications based on the defined multi-level precedence.
	sort.SliceStable(uniqueApps, func(i, j int) bool {
		appI := uniqueApps[i]
		appJ := uniqueApps[j]

		// Primary sort key: Competition Type Precedence (descending)
		precI := competitionPrecedence(appI.CompetitionType())
		precJ := competitionPrecedence(appJ.CompetitionType())
		if precI != precJ {
			return precI > precJ
		}

		// Secondary sort key (for non-BVI): Rating Place (ascending)
		if appI.CompetitionType() != CompetitionBVI {
			if appI.RatingPlace() != appJ.RatingPlace() {
				return appI.RatingPlace() < appJ.RatingPlace()
			}
		}

		// Tertiary sort key / Tie-breaker: Original Score (ascending)
		return appI.Score() < appJ.Score()
	})

	// Step 4: Re-number the rating places consecutively.
	for i, app := range uniqueApps {
		app.ratingPlace = i + 1
	}

	return uniqueApps
}

// competitionPrecedence returns the precedence value for competition types.
// Higher values have higher precedence (will be sorted first).
// Precedence order: BVI > TargetQuota > DedicatedQuota > SpecialQuota > Regular
func competitionPrecedence(ct Competition) int {
	switch ct {
	case CompetitionBVI:
		return 5
	case CompetitionTargetQuota:
		return 4
	case CompetitionDedicatedQuota:
		return 3
	case CompetitionSpecialQuota:
		return 2
	case CompetitionRegular:
		return 1
	default:
		return 0 // Unknown competition types get lowest precedence
	}
}
