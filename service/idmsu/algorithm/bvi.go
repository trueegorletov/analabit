package algorithm

import "sort"

// InternalApplicant represents an MSU (internal) applicant appearing on the university list.
// Only the fields required for BVI matching are present.
type InternalApplicant struct {
	InternalID  string // Raw MSU internal ID (e.g. "028478")
	RatingPlace int    // Position on the MSU rating list (1 = best)
}

// ExternalApplicant represents a Gosuslugi applicant row for the same program/list.
// CanonicalID is the final nationwide ID (idApplication).
type ExternalApplicant struct {
	CanonicalID string // 13-digit canonical ID from Gosuslugi
	RatingPlace int    // Position on Gosuslugi list (1 = best)
}

// MatchBVI aligns BVI applicants by rating position order.
//
// The logic is straightforward:
//  1. Sort both slices by RatingPlace ascending.
//  2. Pair elements by index up to min(len(internal), len(external)).
//  3. Return mapping internalID -> canonicalID.
//  4. The boolean return indicates whether every internal applicant found a match.
//     (false means some fallback handling will be required by caller).
func MatchBVI(internal []InternalApplicant, external []ExternalApplicant) (map[string]string, bool) {
	result := make(map[string]string, len(internal))

	// Defensive copy to avoid mutating caller slices.
	in := make([]InternalApplicant, len(internal))
	copy(in, internal)
	ex := make([]ExternalApplicant, len(external))
	copy(ex, external)

	sort.Slice(in, func(i, j int) bool { return in[i].RatingPlace < in[j].RatingPlace })
	sort.Slice(ex, func(i, j int) bool { return ex[i].RatingPlace < ex[j].RatingPlace })

	paired := 0
	for i := 0; i < len(in) && i < len(ex); i++ {
		result[in[i].InternalID] = ex[i].CanonicalID
		paired++
	}
	return result, paired == len(in)
}
