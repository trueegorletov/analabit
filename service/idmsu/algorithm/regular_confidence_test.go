package algorithm

import "testing"

func TestScoreRegular(t *testing.T) {
	cases := []struct {
		name string
		in   RegularApplicant
		ex   RegularApplicant
		want float64
	}{
		{"perfect", RegularApplicant{ScoreSum: 300, Priority: 1}, RegularApplicant{ScoreSum: 300, Priority: 1}, 1.0},
		{"same score diff priority small", RegularApplicant{ScoreSum: 290, Priority: 1}, RegularApplicant{ScoreSum: 290, Priority: 3}, 0.8},
		{"score diff small same priority", RegularApplicant{ScoreSum: 295, Priority: 2}, RegularApplicant{ScoreSum: 292, Priority: 2}, 0.7},
		{"score diff medium", RegularApplicant{ScoreSum: 280, Priority: 2}, RegularApplicant{ScoreSum: 285, Priority: 3}, 0.5},
		{"mismatch", RegularApplicant{ScoreSum: 270, Priority: 1}, RegularApplicant{ScoreSum: 250, Priority: 5}, 0.0},
	}

	for _, tc := range cases {
		got := ScoreRegular(tc.in, tc.ex)
		if got != tc.want {
			t.Errorf("%s: expected %.1f got %.1f", tc.name, tc.want, got)
		}
	}
}
