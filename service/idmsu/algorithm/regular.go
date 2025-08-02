package algorithm

import "sort"

// CompetitionType identifies applicant category.
// Only Regular and Quota relevant for phase-2.
type CompetitionType int

const (
	CompetitionRegular CompetitionType = iota
	CompetitionQuota
)

// RegularApplicant extends InternalApplicant with score information.
// The same struct used for both internal and external lists.
type RegularApplicant struct {
	InternalID  string // empty for external side
	CanonicalID string // empty for internal side
	RatingPlace int    // list position (1 best)
	ScoreSum    int    // EGE+DVIsum for MSU, totalScore for Gosuslugi
	Priority    int    // Application priority (1 best)
	Competition CompetitionType
}

// MatchRegular performs naive 1-to-1 alignment for Regular & Quota applicants.
// Matching rule (MVP):
//   - Items sorted by RatingPlace ascending.
//   - Only pairs with identical ScoreSum & Priority are matched.
//   - Returns mapping internalID -> canonicalID with basic confidence 1.0 when exact fields match.
//   - Boolean allMatched == true when every internal applicant found candidate.
func MatchRegular(internal, external []RegularApplicant) (map[string]string, bool) {
	sort.Slice(internal, func(i, j int) bool { return internal[i].RatingPlace < internal[j].RatingPlace })
	sort.Slice(external, func(i, j int) bool { return external[i].RatingPlace < external[j].RatingPlace })

	result := map[string]string{}
	ei := 0
	for _, in := range internal {
		for ei < len(external) && (external[ei].ScoreSum != in.ScoreSum || external[ei].Priority != in.Priority || external[ei].Competition != in.Competition) {
			ei++ // skip non-matching
		}
		if ei >= len(external) {
			return result, false // unmatched remainder
		}
		result[in.InternalID] = external[ei].CanonicalID
		ei++
	}
	return result, true
}
