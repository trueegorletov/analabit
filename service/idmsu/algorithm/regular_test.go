package algorithm

import "testing"

func TestMatchRegular_AllMatched(t *testing.T) {
	internal := []RegularApplicant{
		{InternalID: "I1", ScoreSum: 300, Priority: 1, RatingPlace: 1, Competition: CompetitionRegular},
		{InternalID: "I2", ScoreSum: 290, Priority: 2, RatingPlace: 2, Competition: CompetitionRegular},
	}
	external := []RegularApplicant{
		{CanonicalID: "C1", ScoreSum: 300, Priority: 1, RatingPlace: 10, Competition: CompetitionRegular},
		{CanonicalID: "C2", ScoreSum: 290, Priority: 2, RatingPlace: 20, Competition: CompetitionRegular},
	}

	mapping, ok := MatchRegular(internal, external)
	if !ok {
		t.Fatalf("expected all matched")
	}
	if mapping["I1"] != "C1" || mapping["I2"] != "C2" {
		t.Fatalf("unexpected mapping: %+v", mapping)
	}
}

func TestMatchRegular_PartialMatch(t *testing.T) {
	internal := []RegularApplicant{{InternalID: "I1", ScoreSum: 300, Priority: 1, RatingPlace: 1, Competition: CompetitionRegular}}
	external := []RegularApplicant{{CanonicalID: "C2", ScoreSum: 290, Priority: 2, RatingPlace: 5, Competition: CompetitionRegular}}

	_, ok := MatchRegular(internal, external)
	if ok {
		t.Fatalf("expected partial match failure")
	}
}
