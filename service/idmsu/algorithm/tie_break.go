package algorithm

// Candidate represents a possible match with confidence score.
// Used in tie-break phase when multiple candidate canonical IDs exist for an internal applicant (or vice-versa).
type Candidate struct {
	CanonicalID string
	Confidence  float64
}

// ResolveConflicts takes a map from internalID to candidate list and selects a single canonical ID per internal
// while ensuring each canonical ID is assigned at most once.
// Selection rules:
//  1. Highest confidence wins (ties resolved by internal lexical order stable pass).
//  2. If a canonical ID is already claimed with equal confidence, keep the first claimant.
//
// Returns final mapping internal->canonical and slice of unresolved internal IDs.
func ResolveConflicts(candidates map[string][]Candidate) (map[string]string, []string) {
	final := make(map[string]string)
	claimed := make(map[string]struct{})
	unresolved := []string{}

	// Iterate deterministically by internalID lexical order for reproducibility.
	keys := make([]string, 0, len(candidates))
	for k := range candidates {
		keys = append(keys, k)
	}
	// simple insertion sort for small slices
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j-1] > keys[j]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}

	for _, internalID := range keys {
		best := Candidate{Confidence: -1}
		for _, cand := range candidates[internalID] {
			if cand.Confidence > best.Confidence {
				best = cand
			}
		}
		if best.Confidence < 0 {
			unresolved = append(unresolved, internalID)
			continue
		}
		if _, taken := claimed[best.CanonicalID]; taken {
			unresolved = append(unresolved, internalID)
			continue
		}
		final[internalID] = best.CanonicalID
		claimed[best.CanonicalID] = struct{}{}
	}
	return final, unresolved
}
