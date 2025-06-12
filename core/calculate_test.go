package core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

// Helper to get admitted student IDs for a specific heading code from results.
func getAdmittedStudentIDs(results []CalculationResult, headingCode string) []string {
	for _, res := range results {
		if res.Heading.Code() == headingCode { // Use exported field Heading
			ids := make([]string, len(res.Admitted))
			for i, s := range res.Admitted { // Use exported field Admitted
				ids[i] = s.id // Use exported field ID
			}
			return ids
		}
	}
	return []string{} // Return empty slice if heading not found or no admitted students
}

// Helper to create a student ID string.
func sid(i int) string {
	return fmt.Sprintf("student%d", i)
}

// TestCalculateAdmissions_BasicQuotas tests basic quota filling.
func TestCalculateAdmissions_BasicQuotas(t *testing.T) {
	v := NewVarsityCalculator("TEST_VARSITY", "")
	v.AddHeading("H1", Capacities{Regular: 0, TargetQuota: 1, DedicatedQuota: 1, SpecialQuota: 1}, "Heading 1")

	// Target Quota
	v.AddApplication("H1", sid(1), 10, 1, CompetitionTargetQuota, 0) // Should be admitted
	v.AddApplication("H1", sid(2), 5, 1, CompetitionTargetQuota, 0)  // Should be admitted, displacing student1
	v.AddApplication("H1", sid(3), 15, 1, CompetitionTargetQuota, 0) // Should not be admitted

	// Dedicated Quota
	v.AddApplication("H1", sid(4), 10, 1, CompetitionDedicatedQuota, 0) // Should be admitted
	v.AddApplication("H1", sid(5), 5, 1, CompetitionDedicatedQuota, 0)  // Should be admitted, displacing student4

	// Special Quota
	v.AddApplication("H1", sid(6), 10, 1, CompetitionSpecialQuota, 0) // Should be admitted

	results := v.CalculateAdmissions()
	admittedH1 := getAdmittedStudentIDs(results, "H1")

	// Expected: student2 (Target), student5 (Dedicated), student6 (Special)
	// Order within quotas is by rating, but here we check presence.
	// The final list is sorted by outscores, so Quota types then rating.
	// Since all are different quota types, their original rating place determines order among them if they had same quota type.
	// Here, they are distinct quotas, so the order is Target, Dedicated, Special (due to enum values) then rating.
	// However, getAdmittedStudentIDs doesn't guarantee order based on internal admission logic,
	// so we check for presence and count.
	assert.Len(t, admittedH1, 3, "H1 should have 3 admitted students")
	assert.Contains(t, admittedH1, sid(2)) // Best Target
	assert.Contains(t, admittedH1, sid(5)) // Best Dedicated
	assert.Contains(t, admittedH1, sid(6)) // Only Special
}

// TestCalculateAdmissions_QuotaFailureNoFallback tests that students failing a quota don't fall back.
func TestCalculateAdmissions_QuotaFailureNoFallback(t *testing.T) {
	v := NewVarsityCalculator("TEST_VARSITY", "")
	v.AddHeading("H1", Capacities{Regular: 1, TargetQuota: 1}, "Heading 1")

	v.AddApplication("H1", sid(1), 10, 1, CompetitionTargetQuota, 0) // Admitted to Target
	v.AddApplication("H1", sid(2), 5, 1, CompetitionTargetQuota, 0)  // Admitted to Target (displaces s1)
	v.AddApplication("H1", sid(3), 15, 1, CompetitionTargetQuota, 0) // Fails Target, should not get Regular
	v.AddApplication("H1", sid(4), 20, 1, CompetitionRegular, 0)     // Admitted to Regular

	results := v.CalculateAdmissions()
	admittedH1 := getAdmittedStudentIDs(results, "H1")

	assert.Len(t, admittedH1, 2, "H1 should have 2 admitted students")
	assert.Contains(t, admittedH1, sid(2)) // student2 from TargetQuota
	assert.Contains(t, admittedH1, sid(4)) // student4 from Regular
	assert.NotContains(t, admittedH1, sid(1))
	assert.NotContains(t, admittedH1, sid(3))
}

// TestCalculateAdmissions_UnfilledQuotasAddToGeneral tests that unfilled quota spots go to Regular.
func TestCalculateAdmissions_UnfilledQuotasAddToGeneral(t *testing.T) {
	v := NewVarsityCalculator("TEST_VARSITY", "")
	// Target:1, Dedicated:1, Special:1, Regular:1. Total initial: 4
	// Effective Regular should be 1 (base) + 2 (unfilled TQ, DQ) = 3
	v.AddHeading("H1", Capacities{Regular: 1, TargetQuota: 2, DedicatedQuota: 2, SpecialQuota: 1}, "Heading 1")

	// Special Quota (1 spot)
	v.AddApplication("H1", sid(1), 10, 1, CompetitionSpecialQuota, 0) // Admitted to Special

	// Target Quota (2 spots) - 0 applicants, 2 unfilled
	// Dedicated Quota (2 spots) - 0 applicants, 2 unfilled

	// Regular Competition (Base 1 + 2 TQ_unfilled + 2 DQ_unfilled = 5 effective spots)
	v.AddApplication("H1", sid(2), 100, 1, CompetitionBVI, 0)    // Admitted BVI
	v.AddApplication("H1", sid(3), 50, 1, CompetitionBVI, 0)     // Admitted BVI (better rating)
	v.AddApplication("H1", sid(4), 10, 1, CompetitionRegular, 0) // Admitted Regular
	v.AddApplication("H1", sid(5), 5, 1, CompetitionRegular, 0)  // Admitted Regular (better rating)
	v.AddApplication("H1", sid(6), 1, 1, CompetitionRegular, 0)  // Admitted Regular (best rating)
	v.AddApplication("H1", sid(7), 20, 1, CompetitionRegular, 0) // Not admitted (Regular full)

	results := v.CalculateAdmissions()
	admittedH1 := getAdmittedStudentIDs(results, "H1")

	// Expected: s1 (SQ), s3 (BVI), s2 (BVI), s6 (Reg), s5 (Reg), s4 (Reg)
	// Total 1 (SQ) + 2 (BVI) + 3 (Regular from effective general) = 6
	// Capacities: Regular:1, TQ:2, DQ:2, SQ:1.
	// SQ: s1 (1/1 filled)
	// TQ: 0/2 filled (2 unfilled)
	// DQ: 0/2 filled (2 unfilled)
	// Regular: Base 1 + 2 (from TQ) + 2 (from DQ) = 5 spots
	// BVI: s3 (50), s2 (100) -> both admitted
	// Regular: s6 (1), s5 (5), s4 (10) -> all 3 admitted
	// Total admitted: 1 (SQ) + 2 (BVI) + 3 (Regular) = 6
	assert.Len(t, admittedH1, 6, "H1 should have 6 admitted students")
	assert.Contains(t, admittedH1, sid(1)) // Special
	assert.Contains(t, admittedH1, sid(2)) // BVI
	assert.Contains(t, admittedH1, sid(3)) // BVI
	assert.Contains(t, admittedH1, sid(4)) // Regular
	assert.Contains(t, admittedH1, sid(5)) // Regular
	assert.Contains(t, admittedH1, sid(6)) // Regular
	assert.NotContains(t, admittedH1, sid(7))
}

// TestCalculateAdmissions_BVIvsRegular tests BVI priority over Regular in Regular competition.
func TestCalculateAdmissions_BVIvsRegular(t *testing.T) {
	v := NewVarsityCalculator("TEST_VARSITY", "")
	v.AddHeading("H1", Capacities{Regular: 2}, "Heading 1") // 2 Regular spots

	v.AddApplication("H1", sid(1), 200, 1, CompetitionRegular, 0) // Regular
	v.AddApplication("H1", sid(2), 100, 1, CompetitionBVI, 0)     // BVI (better rating than s1, also BVI)
	v.AddApplication("H1", sid(3), 10, 1, CompetitionRegular, 0)  // Regular (better rating than s1)
	v.AddApplication("H1", sid(4), 50, 1, CompetitionBVI, 0)      // BVI (better rating than s2)

	results := v.CalculateAdmissions()
	admittedH1 := getAdmittedStudentIDs(results, "H1")

	// Expected: s4 (BVI, 50), s2 (BVI, 100). s3 and s1 (Regulars) should not be admitted.
	assert.Len(t, admittedH1, 2, "H1 should have 2 admitted students")
	assert.Contains(t, admittedH1, sid(4)) // BVI, rating 50
	assert.Contains(t, admittedH1, sid(2)) // BVI, rating 100
	assert.NotContains(t, admittedH1, sid(1))
	assert.NotContains(t, admittedH1, sid(3))
}

// TestCalculateAdmissions_StudentPriorities tests student choosing higher priority application.
func TestCalculateAdmissions_StudentPriorities(t *testing.T) {
	v := NewVarsityCalculator("TEST_VARSITY", "")
	v.AddHeading("H1", Capacities{Regular: 1}, "Heading 1")
	v.AddHeading("H2", Capacities{Regular: 1}, "Heading 2")

	// Student1: Prefers H1 (priority 1) over H2 (priority 2)
	v.AddApplication("H1", sid(1), 10, 1, CompetitionRegular, 0)
	v.AddApplication("H2", sid(1), 5, 2, CompetitionRegular, 0) // Better rating, but lower priority

	// Student2: Can only go to H2
	v.AddApplication("H2", sid(2), 20, 1, CompetitionRegular, 0)

	results := v.CalculateAdmissions()
	admittedH1 := getAdmittedStudentIDs(results, "H1")
	admittedH2 := getAdmittedStudentIDs(results, "H2")

	assert.Len(t, admittedH1, 1, "H1 should have 1 student")
	assert.Contains(t, admittedH1, sid(1), "Student1 should be in H1")

	assert.Len(t, admittedH2, 1, "H2 should have 1 student")
	assert.Contains(t, admittedH2, sid(2), "Student2 should be in H2")
}

// TestCalculateAdmissions_DisplacementAndReconsideration tests displacement and subsequent reconsideration.
func TestCalculateAdmissions_DisplacementAndReconsideration(t *testing.T) {
	v := NewVarsityCalculator("TEST_VARSITY", "")
	v.AddHeading("H1", Capacities{Regular: 1}, "Heading 1")
	v.AddHeading("H2", Capacities{Regular: 1}, "Heading 2")

	// Student1: H1 (Prio 1), H2 (Prio 2)
	v.AddApplication("H1", sid(1), 100, 1, CompetitionRegular, 0)
	v.AddApplication("H2", sid(1), 10, 2, CompetitionRegular, 0)

	// Student2: H1 (Prio 1) - better rating than Student1 for H1
	v.AddApplication("H1", sid(2), 50, 1, CompetitionRegular, 0)

	// Student3: H2 (Prio 1) - will initially take H2
	v.AddApplication("H2", sid(3), 20, 1, CompetitionRegular, 0)

	// Initial state (conceptual):
	// H1: Student1 (100) - provisional
	// H2: Student3 (20) - provisional
	// Student1 best: H1 (100)

	// Process Student2:
	// Student2 (50) applies to H1, outscores Student1 (100).
	// H1: Student2 (50)
	// Student1 is displaced from H1. Student1 needs reconsideration.

	// Reconsider Student1:
	// Student1's next best is H2 (Prio 2, rating 10).
	// Student1 (10) applies to H2, outscores Student3 (20).
	// H2: Student1 (10)
	// Student3 is displaced from H2. Student3 has no other applications.

	results := v.CalculateAdmissions()
	admittedH1 := getAdmittedStudentIDs(results, "H1")
	admittedH2 := getAdmittedStudentIDs(results, "H2")

	assert.Len(t, admittedH1, 1, "H1 should have 1 student")
	assert.Contains(t, admittedH1, sid(2), "Student2 should be in H1")

	assert.Len(t, admittedH2, 1, "H2 should have 1 student")
	assert.Contains(t, admittedH2, sid(1), "Student1 should be in H2 after displacement")

	assert.NotContains(t, admittedH1, sid(1))
	assert.NotContains(t, admittedH2, sid(3))
}

// TestCalculateAdmissions_FullScenarioWithQuotasAndGeneral tests a more complex scenario.
func TestCalculateAdmissions_FullScenarioWithQuotasAndGeneral(t *testing.T) {
	v := NewVarsityCalculator("TEST_VARSITY", "")
	// H1: TQ=1, DQ=1, SQ=1, Gen=1. Total initial capacity = 4.
	v.AddHeading("H1", Capacities{TargetQuota: 1, DedicatedQuota: 1, SpecialQuota: 1, Regular: 1}, "H1")
	// H2: Gen=2
	v.AddHeading("H2", Capacities{Regular: 2}, "H2")

	// Student ApplicationsCache:
	// S1: H1 (TQ, 10, P1), H2 (Reg, 20, P2)
	v.AddApplication("H1", sid(1), 10, 1, CompetitionTargetQuota, 0)
	v.AddApplication("H2", sid(1), 20, 2, CompetitionRegular, 0)

	// S2: H1 (TQ, 5, P1) -> better TQ for H1
	v.AddApplication("H1", sid(2), 5, 1, CompetitionTargetQuota, 0)

	// S3: H1 (DQ, 100, P1)
	v.AddApplication("H1", sid(3), 100, 1, CompetitionDedicatedQuota, 0)

	// S4: H1 (SQ, 100, P1)
	v.AddApplication("H1", sid(4), 100, 1, CompetitionSpecialQuota, 0)

	// S5: H1 (BVI, 100, P1), H2 (BVI, 5, P2)
	v.AddApplication("H1", sid(5), 100, 1, CompetitionBVI, 0)
	v.AddApplication("H2", sid(5), 5, 2, CompetitionBVI, 0)

	// S6: H1 (Reg, 10, P1)
	v.AddApplication("H1", sid(6), 10, 1, CompetitionRegular, 0) // Will compete for general H1

	// S7: H2 (Reg, 10, P1)
	v.AddApplication("H2", sid(7), 10, 1, CompetitionRegular, 0)

	// S8: H2 (BVI, 1, P1) -> best BVI for H2
	v.AddApplication("H2", sid(8), 1, 1, CompetitionBVI, 0)

	// Expected H1 admissions:
	// TQ: S2 (5) - S1 (10) is outbid for TQ, S1 will try H2.
	// DQ: S3 (100)
	// SQ: S4 (100)
	// Regular (1 spot, no unfilled quotas): S5 (BVI, 100) beats S6 (Regular, 10)
	// H1: S2, S3, S4, S5. (4 students)

	// Expected H2 admissions (2 spots):
	// S1 (Reg, 20, P2 from H1 TQ failure)
	// S5 (BVI, 5, P2 from H1 Regular failure - but S5 got into H1 Regular, so this H2 app is moot for S5)
	// S7 (Reg, 10, P1)
	// S8 (BVI, 1, P1)
	//
	// Iteration:
	// Student S2 gets H1 TQ.
	// Student S3 gets H1 DQ.
	// Student S4 gets H1 SQ.
	// Student S5 applies H1 BVI (P1). H1 Regular has 1 spot. S5 gets it. studentBestPlacement[S5] = H1/BVI.
	// Student S1 applies H1 TQ (P1), rating 10. S2 (rating 5) already has it. S1 fails H1 TQ.
	// Student S6 applies H1 Reg (P1). H1 Regular is full (S5). S6 fails H1 Regular.
	//
	// Student S1 now tries H2 (Reg, 20, P2). H2 has 2 spots. S1 gets one. studentBestPlacement[S1] = H2/Reg.
	// Student S7 applies H2 (Reg, 10, P1). H2 has 1 spot left. S7 gets it. studentBestPlacement[S7] = H2/Reg.
	// Student S8 applies H2 (BVI, 1, P1). H2 is full (S1, S7). S8 (BVI) outscores S1 (Reg) and S7 (Reg).
	//   S8 outscores S1 (rating 20). S1 is displaced. studentBestPlacement[S8] = H2/BVI.
	//   S8 outscores S7 (rating 10). S7 is displaced. (S8 takes one spot, one more to fill or re-evaluate)
	//   Let's trace carefully: H2 spots: 2.
	//   S1 (Reg 20 P2), S7 (Reg 10 P1), S8 (BVI 1 P1).
	//   S8 (BVI 1 P1) is best. Takes a spot.
	//   S7 (Reg 10 P1) is next best of remaining. Takes a spot.
	//   S1 (Reg 20 P2) is out.
	//   So H2: S8, S7.
	//
	// S1 was displaced from H2. S1 has no more applications.
	// S5 is happy in H1 (BVI). Its H2 app (BVI 5 P2) is lower priority.
	//
	// Final H1: S2 (TQ), S3 (DQ), S4 (SQ), S5 (BVI).
	// Final H2: S8 (BVI), S7 (Reg).

	results := v.CalculateAdmissions()
	admittedH1 := getAdmittedStudentIDs(results, "H1")
	admittedH2 := getAdmittedStudentIDs(results, "H2")

	// Sort for consistent comparison if needed, though Contains checks are primary.
	sort.Strings(admittedH1)
	sort.Strings(admittedH2)

	expectedH1 := []string{sid(2), sid(3), sid(4), sid(5)}
	sort.Strings(expectedH1)
	assert.Equal(t, expectedH1, admittedH1, "H1 admitted students mismatch")

	expectedH2 := []string{sid(7), sid(8)}
	sort.Strings(expectedH2)
	assert.Equal(t, expectedH2, admittedH2, "H2 admitted students mismatch")

	// Check that S1 (who lost H1 TQ and then H2 Regular) is not admitted anywhere.
	assert.NotContains(t, admittedH1, sid(1))
	assert.NotContains(t, admittedH2, sid(1))
	// Check S6 (lost H1 Regular)
	assert.NotContains(t, admittedH1, sid(6))
}

// TestCalculateAdmissions_QuitStudent tests that a quit student is not admitted.
func TestCalculateAdmissions_QuitStudent(t *testing.T) {
	v := NewVarsityCalculator("TEST_VARSITY", "")
	v.AddHeading("H1", Capacities{Regular: 1}, "H1")

	v.AddApplication("H1", sid(1), 10, 1, CompetitionRegular, 0)
	v.SetQuit(sid(1)) // Student1 quits

	v.AddApplication("H1", sid(2), 20, 1, CompetitionRegular, 0) // Student2 should get the spot

	results := v.CalculateAdmissions()
	admittedH1 := getAdmittedStudentIDs(results, "H1")

	assert.Len(t, admittedH1, 1, "H1 should have 1 student")
	assert.Contains(t, admittedH1, sid(2))
	assert.NotContains(t, admittedH1, sid(1))
}

// TestCalculateAdmissions_OrderOfResults tests that results are sorted by heading code.
func TestCalculateAdmissions_OrderOfResults(t *testing.T) {
	v := NewVarsityCalculator("TEST_VARSITY", "")
	v.AddHeading("H2", Capacities{Regular: 1}, "Heading 2")
	v.AddHeading("H1", Capacities{Regular: 1}, "Heading 1") // Added out of order

	v.AddApplication("H1", sid(1), 10, 1, CompetitionRegular, 0)
	v.AddApplication("H2", sid(2), 10, 1, CompetitionRegular, 0)

	results := v.CalculateAdmissions()

	if len(results) == 2 {
		assert.Equal(t, "H1", results[0].Heading.Code(), "First result should be H1")  // Use exported field Heading
		assert.Equal(t, "H2", results[1].Heading.Code(), "Second result should be H2") // Use exported field Heading
	} else {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}
}

// TestCalculateAdmissions_SameRatingDifferentCompetitionTypeInGeneral
// Tests that BVI is preferred over Regular even if Regular has "better" rating number.
func TestCalculateAdmissions_SameRatingDifferentCompetitionTypeInGeneral(t *testing.T) {
	v := NewVarsityCalculator("TEST_VARSITY", "")
	v.AddHeading("H1", Capacities{Regular: 1}, "H1") // Only one spot

	v.AddApplication("H1", sid(1), 1, 1, CompetitionRegular, 0) // Best possible rating, but Regular
	v.AddApplication("H1", sid(2), 999, 1, CompetitionBVI, 0)   // Worst possible rating, but BVI

	results := v.CalculateAdmissions()
	admittedH1 := getAdmittedStudentIDs(results, "H1")

	assert.Len(t, admittedH1, 1)
	assert.Contains(t, admittedH1, sid(2), "BVI student should be admitted over Regular student")
	assert.NotContains(t, admittedH1, sid(1))
}

// TestCalculateAdmissions_QuotaStudentsDoNotCompeteForGeneralIfQuotaFull
// Ensures that if a student applies for a quota and it's full (even if they are outranked),
// they don't automatically compete for general spots for THAT SAME application.
func TestCalculateAdmissions_QuotaStudentsDoNotCompeteForGeneralIfQuotaFull(t *testing.T) {
	v := NewVarsityCalculator("TEST_VARSITY", "")
	v.AddHeading("H1", Capacities{TargetQuota: 1, Regular: 1}, "H1")

	// S1 fills the TargetQuota
	v.AddApplication("H1", sid(1), 10, 1, CompetitionTargetQuota, 0)
	// S2 applies for TargetQuota but S1 is better or already there. S2 should NOT get general for this app.
	v.AddApplication("H1", sid(2), 20, 1, CompetitionTargetQuota, 0)
	// S3 applies for Regular and should get it.
	v.AddApplication("H1", sid(3), 1, 1, CompetitionRegular, 0)

	results := v.CalculateAdmissions()
	admittedH1 := getAdmittedStudentIDs(results, "H1")
	sort.Strings(admittedH1)

	expected := []string{sid(1), sid(3)}
	sort.Strings(expected)

	assert.Equal(t, expected, admittedH1)
	assert.NotContains(t, admittedH1, sid(2))
}

// TestCalculateAdmissions_ComplexPrioritiesAndDisplacement
func TestCalculateAdmissions_ComplexPrioritiesAndDisplacement(t *testing.T) {
	v := NewVarsityCalculator("TEST_VARSITY", "")
	v.AddHeading("H_HIGH_CAP", Capacities{Regular: 3}, "High Capacity Heading")
	v.AddHeading("H_LOW_CAP_PRIO1", Capacities{Regular: 1}, "Low Cap Prio 1 Heading")
	v.AddHeading("H_LOW_CAP_PRIO2", Capacities{Regular: 1}, "Low Cap Prio 2 Heading")

	// S1:
	// 1. H_LOW_CAP_PRIO1 (Rating 100, P1)
	// 2. H_HIGH_CAP (Rating 10, P2)
	v.AddApplication("H_LOW_CAP_PRIO1", sid(1), 100, 1, CompetitionRegular, 0)
	v.AddApplication("H_HIGH_CAP", sid(1), 10, 2, CompetitionRegular, 0)

	// S2:
	// 1. H_LOW_CAP_PRIO1 (Rating 50, P1) -> Will take this spot from S1
	v.AddApplication("H_LOW_CAP_PRIO1", sid(2), 50, 1, CompetitionRegular, 0)

	// S3:
	// 1. H_HIGH_CAP (Rating 20, P1)
	v.AddApplication("H_HIGH_CAP", sid(3), 20, 1, CompetitionRegular, 0)

	// S4:
	// 1. H_LOW_CAP_PRIO2 (Rating 10, P1)
	// 2. H_HIGH_CAP (Rating 5, P2)
	v.AddApplication("H_LOW_CAP_PRIO2", sid(4), 10, 1, CompetitionRegular, 0)
	v.AddApplication("H_HIGH_CAP", sid(4), 5, 2, CompetitionRegular, 0)

	// Expected outcome based on students prioritizing their highest-numbered priority applications
	// and displacements occurring based on rating. The core logic is that once a student
	// is placed via an application of a certain priority (e.g., P1), they will not
	// attempt to take a lower priority application (e.g., P2), even if the P2 application
	// has a better rating or could displace someone.
	//
	// 1. Initial P1 considerations:
	//    - S1 provisionally takes H_LOW_CAP_PRIO1 (R100).
	//    - S2 applies to H_LOW_CAP_PRIO1 (R50), displaces S1.
	//      Result: H_LOW_CAP_PRIO1 = [S2(R50)]. S1 is unplaced. studentBestPlacement[S2] = H_LOW_CAP_PRIO1 (P1)
	//    - S3 takes H_HIGH_CAP (R20, P1).
	//      Result: H_HIGH_CAP = [S3(R20)]. (Capacity 3). studentBestPlacement[S3] = H_HIGH_CAP (P1)
	//    - S4 takes H_LOW_CAP_PRIO2 (R10, P1).
	//      Result: H_LOW_CAP_PRIO2 = [S4(R10)]. studentBestPlacement[S4] = H_LOW_CAP_PRIO2 (P1)
	//
	// Current student placements after initial P1 processing:
	//    S1: Unplaced
	//    S2: H_LOW_CAP_PRIO1 (P1)
	//    S3: H_HIGH_CAP (P1)
	//    S4: H_LOW_CAP_PRIO2 (P1)
	//
	// 2. Reconsider S1 (displaced):
	//    - S1's P1 (H_LOW_CAP_PRIO1) is taken by S2.
	//    - S1 considers P2: H_HIGH_CAP (R10).
	//    - H_HIGH_CAP has S3(R20) and 2 free spots. S1(R10) takes a spot.
	//      Result: H_HIGH_CAP = [S1(R10), S3(R20)]. (Sorted by rating: S1, S3). studentBestPlacement[S1] = H_HIGH_CAP (P2)
	//
	// 3. S4's situation:
	//    - S4 is placed in H_LOW_CAP_PRIO2 (P1, R10). studentBestPlacement[S4] has priority 1.
	//    - S4's P2 application is H_HIGH_CAP (R5, priority 2).
	//    - When the algorithm processes S4's applications, after placing S4 in its P1 choice,
	//      the check `if existingPlacement.priority < currentApp.priority` (i.e., 1 < 2)
	//      will be true for S4's P2 application. The loop over S4's applications will `break`.
	//    - Thus, S4 will NOT attempt to move to H_HIGH_CAP for its P2 application,
	//      even though its rating (R5) for H_HIGH_CAP is better than S1's (R10) in H_HIGH_CAP.
	//
	// Final State:
	// H_LOW_CAP_PRIO1: S2 (sid2)
	// H_LOW_CAP_PRIO2: S4 (sid4)
	// H_HIGH_CAP: S1 (sid1), S3 (sid3)

	results := v.CalculateAdmissions()

	admittedPrio1 := getAdmittedStudentIDs(results, "H_LOW_CAP_PRIO1")
	admittedPrio2 := getAdmittedStudentIDs(results, "H_LOW_CAP_PRIO2")
	admittedHighCap := getAdmittedStudentIDs(results, "H_HIGH_CAP")

	sort.Strings(admittedHighCap) // Ensure order for comparison

	// H_LOW_CAP_PRIO1: S2 (Rating 50). S1 was displaced.
	assert.Equal(t, []string{sid(2)}, admittedPrio1, "H_LOW_CAP_PRIO1 mismatch: S2 should be admitted")

	// H_LOW_CAP_PRIO2: S4 (Rating 10, P1 for S4). S4 takes its P1.
	assert.Equal(t, []string{sid(4)}, admittedPrio2, "H_LOW_CAP_PRIO2 mismatch: S4 should be admitted via P1")

	// H_HIGH_CAP:
	// S1 (displaced from H_LOW_CAP_PRIO1) takes P2 to H_HIGH_CAP (Rating 10).
	// S3 takes P1 to H_HIGH_CAP (Rating 20).
	// Both fit. Admitted list internally sorted by rating (S1 then S3).
	// getAdmittedStudentIDs preserves this order.
	expectedHighCap := []string{sid(1), sid(3)}
	// sort.Strings(expectedHighCap) // Not strictly needed as sid(1) then sid(3) is already sorted.
	assert.Equal(t, expectedHighCap, admittedHighCap, "H_HIGH_CAP mismatch: S1 (P2) and S3 (P1) should be admitted")

	// S1 should be in H_HIGH_CAP
	assert.NotContains(t, admittedPrio1, sid(1), "S1 should not be in H_LOW_CAP_PRIO1")
	assert.NotContains(t, admittedPrio2, sid(1), "S1 should not be in H_LOW_CAP_PRIO2")
	assert.Contains(t, admittedHighCap, sid(1), "S1 should be in H_HIGH_CAP")

	// S4 should be in H_LOW_CAP_PRIO2 and NOT H_HIGH_CAP
	assert.NotContains(t, admittedPrio1, sid(4), "S4 should not be in H_LOW_CAP_PRIO1")
	assert.Contains(t, admittedPrio2, sid(4), "S4 should be in H_LOW_CAP_PRIO2")
	assert.NotContains(t, admittedHighCap, sid(4), "S4 should not have taken its P2 to H_HIGH_CAP")
}
