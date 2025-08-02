package algorithm

import "testing"

func TestResolveConflicts(t *testing.T) {
	cases := []struct {
		name      string
		input     map[string][]Candidate
		wantMap   map[string]string
		wantUnres []string
	}{
		{
			name: "simple unique",
			input: map[string][]Candidate{
				"A": []Candidate{{CanonicalID: "X", Confidence: 1.0}},
				"B": []Candidate{{CanonicalID: "Y", Confidence: 0.9}},
			},
			wantMap:   map[string]string{"A": "X", "B": "Y"},
			wantUnres: nil,
		},
		{
			name: "duplicate canonical keep first",
			input: map[string][]Candidate{
				"A": []Candidate{{CanonicalID: "Z", Confidence: 1.0}},
				"B": []Candidate{{CanonicalID: "Z", Confidence: 1.0}},
			},
			wantMap:   map[string]string{"A": "Z"},
			wantUnres: []string{"B"},
		},
		{
			name: "different confidences choose higher",
			input: map[string][]Candidate{
				"A": []Candidate{{CanonicalID: "X", Confidence: 0.8}, {CanonicalID: "Y", Confidence: 0.9}},
				"B": []Candidate{{CanonicalID: "Y", Confidence: 0.9}},
			},
			wantMap:   map[string]string{"A": "Y"},
			wantUnres: []string{"B"},
		},
	}

	equalMap := func(a, b map[string]string) bool {
		if len(a) != len(b) {
			return false
		}
		for k, v := range a {
			if b[k] != v {
				return false
			}
		}
		return true
	}

	for _, tc := range cases {
		gotMap, gotUnres := ResolveConflicts(tc.input)
		if !equalMap(gotMap, tc.wantMap) {
			t.Errorf("%s: mapping mismatch got %v want %v", tc.name, gotMap, tc.wantMap)
		}
		if len(gotUnres) != len(tc.wantUnres) {
			t.Errorf("%s: unresolved len mismatch got %v want %v", tc.name, gotUnres, tc.wantUnres)
		}
	}
}
