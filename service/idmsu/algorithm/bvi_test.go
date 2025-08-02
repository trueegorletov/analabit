package algorithm

import "testing"

func TestMatchBVI_AllMatched(t *testing.T) {
	internal := []InternalApplicant{{"A1", 1}, {"A2", 2}, {"A3", 3}}
	external := []ExternalApplicant{{"C1", 1}, {"C2", 2}, {"C3", 3}}

	mapping, ok := MatchBVI(internal, external)
	if !ok {
		t.Fatalf("expected all applicants matched")
	}
	want := map[string]string{"A1": "C1", "A2": "C2", "A3": "C3"}
	if len(mapping) != len(want) {
		t.Fatalf("unexpected map length: got %d, want %d", len(mapping), len(want))
	}
	for k, v := range want {
		if mapping[k] != v {
			t.Errorf("mismatch for %s: got %s, want %s", k, mapping[k], v)
		}
	}
}

func TestMatchBVI_PartialMatched(t *testing.T) {
	internal := []InternalApplicant{{"A1", 1}, {"A2", 2}, {"A3", 3}}
	external := []ExternalApplicant{{"C1", 5}, {"C2", 7}} // Only two entries

	mapping, ok := MatchBVI(internal, external)
	if ok {
		t.Fatalf("expected partial match flag false")
	}
	if len(mapping) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(mapping))
	}
	if mapping["A1"] == "" || mapping["A2"] == "" {
		t.Errorf("expected A1 and A2 to be matched")
	}
	if _, found := mapping["A3"]; found {
		t.Errorf("did not expect A3 to be matched")
	}
}
