package algorithm

// ScoreRegular returns a confidence score [0,1] for pairing an internal and external regular applicant.
// Deterministic rule set (MVP):
//  1. Perfect match: identical ScoreSum and Priority -> 1.0
//  2. Same ScoreSum, different Priority (<=2 diff) -> 0.8
//  3. ScoreSum diff <= 3 points & same Priority -> 0.7
//  4. ScoreSum diff <= 5 points -> 0.5
//  5. Otherwise -> 0.0 (not considered a match).
func ScoreRegular(in, ex RegularApplicant) float64 {
	switch {
	case in.ScoreSum == ex.ScoreSum && in.Priority == ex.Priority:
		return 1.0
	case in.ScoreSum == ex.ScoreSum && abs(in.Priority-ex.Priority) <= 2:
		return 0.8
	case abs(in.ScoreSum-ex.ScoreSum) <= 3 && in.Priority == ex.Priority:
		return 0.7
	case abs(in.ScoreSum-ex.ScoreSum) <= 5:
		return 0.5
	default:
		return 0.0
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
