package core

import (
	"reflect"
	"testing"
)

// Helper to get admitted student IDs for a specific heading code from results.
func getAdmittedStudentIDs(results []CalculationResult, headingCode string) []string {
	var ids []string
	for _, res := range results {
		if res.heading.code == headingCode {
			for _, s := range res.admitted {
				ids = append(ids, s.id)
			}
			return ids
		}
	}
	return ids // Should not happen if headingCode exists
}

func TestVarsityCalculator_AddHeading(t *testing.T) {
	vc := NewVarsityCalculator("TestUni")
	vc.AddHeading("CS101", 1, "Computer Science")

	if len(vc.headings) != 1 {
		t.Errorf("Expected 1 heading, got %d", len(vc.headings))
	}
	if _, ok := vc.headings["CS101"]; !ok {
		t.Errorf("Heading CS101 not found")
	}
	if vc.headings["CS101"].capacity != 1 {
		t.Errorf("Expected capacity 1 for CS101, got %d", vc.headings["CS101"].capacity)
	}
	if vc.headings["CS101"].prettyName != "Computer Science" {
		t.Errorf("Expected prettyName 'Computer Science', got '%s'", vc.headings["CS101"].prettyName)
	}
}

func TestVarsityCalculator_AddApplication(t *testing.T) {
	vc := NewVarsityCalculator("TestUni")
	vc.AddHeading("CS101", 1, "Computer Science")
	vc.AddApplication("CS101", "1", 100, 1, CompetitionRegular)

	if len(vc.students) != 1 {
		t.Errorf("Expected 1 student, got %d", len(vc.students))
	}
	studentS1, ok := vc.students["1"]
	if !ok {
		t.Errorf("Student %s not found", prepareStudentID("1"))
	}
	if len(studentS1.applications) != 1 {
		t.Errorf("Expected 1 application for %s, got %d", prepareStudentID("1"), len(studentS1.applications))
	}
	app := studentS1.applications[0]
	if app.heading.code != "CS101" {
		t.Errorf("Expected application for CS101, got %s", app.heading.code)
	}
	if app.ratingPlace != 100 {
		t.Errorf("Expected ratingPlace 100, got %d", app.ratingPlace)
	}
	if app.priority != 1 {
		t.Errorf("Expected priority 1, got %d", app.priority)
	}
}

func TestVarsityCalculator_SetQuit(t *testing.T) {
	vc := NewVarsityCalculator("TestUni")
	vc.AddHeading("CS101", 1, "Computer Science")
	vc.AddApplication("CS101", "1", 100, 1, CompetitionRegular)

	if vc.students["1"].quit {
		t.Errorf("Student %s should not be quit initially", prepareStudentID("1"))
	}
	vc.SetQuit("1")
	if !vc.students["1"].quit {
		t.Errorf("Student %s should be quit after SetQuit", prepareStudentID("1"))
	}

	// Test SetQuit for non-existent student (should not panic)
	vc.SetQuit("S_NonExistent")
}

func TestCalculateAdmissions_BasicAdmission(t *testing.T) {
	vc := NewVarsityCalculator("TestUni")
	vc.AddHeading("CS101", 2, "Computer Science")
	vc.AddHeading("MA101", 1, "Mathematics")

	vc.AddApplication("CS101", "1", 10, 1, CompetitionRegular)
	vc.AddApplication("MA101", "2", 5, 1, CompetitionRegular)
	vc.AddApplication("CS101", "3", 12, 1, CompetitionRegular)

	results := vc.CalculateAdmissions()

	csAdmitted := getAdmittedStudentIDs(results, "CS101")
	maAdmitted := getAdmittedStudentIDs(results, "MA101")

	expectedCS := []string{prepareStudentID("1"), prepareStudentID("3")}

	if !reflect.DeepEqual(csAdmitted, expectedCS) {
		t.Errorf("CS101: Expected %v, got %v", expectedCS, csAdmitted)
	}

	expectedMA := []string{prepareStudentID("2")}
	if !reflect.DeepEqual(maAdmitted, expectedMA) {
		t.Errorf("MA101: Expected %v, got %v", expectedMA, maAdmitted)
	}
}

func TestCalculateAdmissions_CapacityLimit(t *testing.T) {
	vc := NewVarsityCalculator("TestUni")
	vc.AddHeading("CS101", 1, "Software Engineering") // Capacity 1

	vc.AddApplication("CS101", "1", 10, 1, CompetitionRegular)
	vc.AddApplication("CS101", "2", 5, 1, CompetitionRegular)
	vc.AddApplication("CS101", "3", 12, 1, CompetitionRegular)

	results := vc.CalculateAdmissions()
	csAdmitted := getAdmittedStudentIDs(results, "CS101")

	expectedCS := []string{prepareStudentID("2")}
	if !reflect.DeepEqual(csAdmitted, expectedCS) {
		t.Errorf("CS101: Expected %v, got %v", expectedCS, csAdmitted)
	}
}

func TestCalculateAdmissions_PriorityHandling(t *testing.T) {
	vc := NewVarsityCalculator("TestUni")
	vc.AddHeading("CS101", 1, "Data Science")
	vc.AddHeading("MA101", 1, "Pure Mathematics")

	vc.AddApplication("CS101", "1", 10, 1, CompetitionRegular)
	vc.AddApplication("MA101", "1", 5, 2, CompetitionRegular)

	vc.AddApplication("MA101", "2", 8, 1, CompetitionRegular)

	results := vc.CalculateAdmissions()
	csAdmitted := getAdmittedStudentIDs(results, "CS101")
	maAdmitted := getAdmittedStudentIDs(results, "MA101")

	expectedCS := []string{prepareStudentID("1")}
	if !reflect.DeepEqual(csAdmitted, expectedCS) {
		t.Errorf("CS101: Expected %v, got %v for student %s by priority", expectedCS, csAdmitted, prepareStudentID("1"))
	}

	expectedMA := []string{prepareStudentID("2")}
	if !reflect.DeepEqual(maAdmitted, expectedMA) {
		t.Errorf("MA101: Expected %v, got %v", expectedMA, maAdmitted)
	}
}

func TestCalculateAdmissions_Displacement(t *testing.T) {
	vc := NewVarsityCalculator("TestUni")
	vc.AddHeading("CS101", 1, "Cybersecurity") // Capacity 1

	vc.AddApplication("CS101", "1", 10, 1, CompetitionRegular)
	vc.AddApplication("CS101", "2", 5, 1, CompetitionRegular)

	_ = vc.student("1")
	_ = vc.student("2")

	results := vc.CalculateAdmissions()
	csAdmitted := getAdmittedStudentIDs(results, "CS101")

	expectedCS := []string{prepareStudentID("2")}
	if !reflect.DeepEqual(csAdmitted, expectedCS) {
		t.Errorf("CS101: Expected %v to be admitted (student %s displaces %s), got %v", expectedCS, prepareStudentID("2"), prepareStudentID("1"), csAdmitted)
	}
}

func TestCalculateAdmissions_StudentQuit(t *testing.T) {
	vc := NewVarsityCalculator("TestUni")
	vc.AddHeading("CS101", 1, "Information Technology")
	vc.AddHeading("MA101", 1, "Theoretical Physics")

	vc.AddApplication("CS101", "1", 10, 1, CompetitionRegular)
	vc.AddApplication("MA101", "2", 5, 1, CompetitionRegular)
	vc.AddApplication("CS101", "3", 12, 1, CompetitionRegular)

	vc.SetQuit("1")

	results := vc.CalculateAdmissions()
	csAdmitted := getAdmittedStudentIDs(results, "CS101")
	maAdmitted := getAdmittedStudentIDs(results, "MA101")

	expectedCS := []string{prepareStudentID("3")}
	if !reflect.DeepEqual(csAdmitted, expectedCS) {
		t.Errorf("CS101: Expected %v, got %v", expectedCS, csAdmitted)
	}
	expectedMA := []string{prepareStudentID("2")}
	if !reflect.DeepEqual(maAdmitted, expectedMA) {
		t.Errorf("MA101: Expected %v, got %v", expectedMA, maAdmitted)
	}
}

func TestCalculateAdmissions_ComplexScenario(t *testing.T) {
	vc := NewVarsityCalculator("TestUni")
	vc.AddHeading("H1", 1, "Applied Mathematics and Informatics")
	vc.AddHeading("H2", 2, "Modern Programming")
	vc.AddHeading("H3", 1, "System Architecture")

	vc.AddApplication("H1", "104", 10, 1, CompetitionRegular)
	vc.AddApplication("H2", "104", 5, 2, CompetitionRegular)

	vc.AddApplication("H2", "102", 8, 1, CompetitionRegular)
	vc.AddApplication("H3", "102", 10, 2, CompetitionRegular)

	vc.AddApplication("H1", "103", 9, 1, CompetitionRegular)
	vc.AddApplication("H2", "103", 12, 2, CompetitionRegular)

	vc.AddApplication("H2", "101", 7, 1, CompetitionRegular)

	vc.AddApplication("H3", "105", 6, 1, CompetitionRegular)
	vc.AddApplication("H2", "105", 20, 2, CompetitionRegular)

	vc.AddApplication("H1", "106", 1, 1, CompetitionRegular)
	vc.SetQuit("106")

	results := vc.CalculateAdmissions()

	h1Admitted := getAdmittedStudentIDs(results, "H1")
	h2Admitted := getAdmittedStudentIDs(results, "H2")
	h3Admitted := getAdmittedStudentIDs(results, "H3")

	expectedH1 := []string{prepareStudentID("103")}
	if !reflect.DeepEqual(h1Admitted, expectedH1) {
		t.Errorf("H1: Expected %v, got %v", expectedH1, h1Admitted)
	}

	expectedH2 := []string{prepareStudentID("104"), prepareStudentID("101")}

	if !reflect.DeepEqual(h2Admitted, expectedH2) {
		t.Errorf("H2: Expected %v, got %v", expectedH2, h2Admitted)
	}

	expectedH3 := []string{prepareStudentID("105")}
	if !reflect.DeepEqual(h3Admitted, expectedH3) {
		t.Errorf("H3: Expected %v, got %v", expectedH3, h3Admitted)
	}
}

func TestCalculateAdmissions_NoApplicants(t *testing.T) {
	vc := NewVarsityCalculator("TestUni")
	vc.AddHeading("CS101", 2, "Network Engineering")
	vc.AddHeading("MA101", 1, "Statistics")

	results := vc.CalculateAdmissions()

	if len(results) != 2 {
		t.Fatalf("Expected 2 results (one for each heading), got %d", len(results))
	}

	csAdmitted := getAdmittedStudentIDs(results, "CS101")
	maAdmitted := getAdmittedStudentIDs(results, "MA101")

	if len(csAdmitted) != 0 {
		t.Errorf("CS101: Expected 0 admitted, got %v", csAdmitted)
	}
	if len(maAdmitted) != 0 {
		t.Errorf("MA101: Expected 0 admitted, got %v", maAdmitted)
	}
}

func TestCalculateAdmissions_AllQuit(t *testing.T) {
	vc := NewVarsityCalculator("TestUni")
	vc.AddHeading("CS101", 1, "Artificial Intelligence")

	vc.AddApplication("CS101", "1", 10, 1, CompetitionRegular)
	vc.AddApplication("CS101", "2", 5, 1, CompetitionRegular)
	vc.SetQuit("1")
	vc.SetQuit("2")

	results := vc.CalculateAdmissions()
	csAdmitted := getAdmittedStudentIDs(results, "CS101")

	if len(csAdmitted) != 0 {
		t.Errorf("CS101: Expected 0 admitted as all quit, got %v", csAdmitted)
	}
}

func TestCalculateAdmissions_StudentGetsSecondPriorityAfterFirstChoiceFull(t *testing.T) {
	vc := NewVarsityCalculator("TestUni")
	vc.AddHeading("H1", 1, "Biotechnology")
	vc.AddHeading("H2", 1, "Chemical Engineering")

	vc.AddApplication("H1", "801", 1, 1, CompetitionRegular)

	vc.AddApplication("H1", "802", 10, 1, CompetitionRegular)
	vc.AddApplication("H2", "802", 5, 2, CompetitionRegular)

	results := vc.CalculateAdmissions()
	h1Admitted := getAdmittedStudentIDs(results, "H1")
	h2Admitted := getAdmittedStudentIDs(results, "H2")

	expectedH1 := []string{prepareStudentID("801")}
	if !reflect.DeepEqual(h1Admitted, expectedH1) {
		t.Errorf("H1: Expected %v, got %v", expectedH1, h1Admitted)
	}

	expectedH2 := []string{prepareStudentID("802")}
	if !reflect.DeepEqual(h2Admitted, expectedH2) {
		t.Errorf("H2: Expected %v, got %v", expectedH2, h2Admitted)
	}
}

func TestCalculateAdmissions_DisplacementCascades(t *testing.T) {
	vc := NewVarsityCalculator("TestUni")
	vc.AddHeading("H1", 1, "Mechanical Engineering")
	vc.AddHeading("H2", 1, "Aerospace Engineering")

	vc.AddApplication("H2", "901", 10, 1, CompetitionRegular)

	vc.AddApplication("H1", "902", 10, 1, CompetitionRegular)
	vc.AddApplication("H2", "902", 5, 2, CompetitionRegular)

	vc.AddApplication("H1", "903", 1, 1, CompetitionRegular)

	results := vc.CalculateAdmissions()
	h1Admitted := getAdmittedStudentIDs(results, "H1")
	h2Admitted := getAdmittedStudentIDs(results, "H2")

	expectedH1 := []string{prepareStudentID("903")}
	if !reflect.DeepEqual(h1Admitted, expectedH1) {
		t.Errorf("H1: Expected %v, got %v", expectedH1, h1Admitted)
	}

	expectedH2 := []string{prepareStudentID("902")}
	if !reflect.DeepEqual(h2Admitted, expectedH2) {
		t.Errorf("H2: Expected %v, got %v", expectedH2, h2Admitted)
	}
}

func TestCalculateAdmissions_EmptyApplicationsForStudent(t *testing.T) {
	vc := NewVarsityCalculator("TestUni")
	vc.AddHeading("H1", 1, "Civil Engineering")

	_ = vc.student("1001")

	vc.AddApplication("H1", "1002", 1, 1, CompetitionRegular)

	results := vc.CalculateAdmissions()
	h1Admitted := getAdmittedStudentIDs(results, "H1")

	expectedH1 := []string{prepareStudentID("1002")}
	if !reflect.DeepEqual(h1Admitted, expectedH1) {
		t.Errorf("H1: Expected %v, got %v", expectedH1, h1Admitted)
	}
	if len(vc.students["1001"].applications) != 0 {
		t.Errorf("Student %s should have 0 applications, got %d", prepareStudentID("1001"), len(vc.students["1001"].applications))
	}
}

// --- New Quota Test Cases ---

func TestCalculateAdmissions_QuotaStudentPrioritizedOverNonQuota(t *testing.T) {
	vc := NewVarsityCalculator("TestUni")
	vc.AddHeading("H1", 1, "Quota Studies")

	vc.AddApplication("H1", "201", 1, 1, CompetitionRegular)
	vc.AddApplication("H1", "202", 10, 1, CompetitionTargetQuota)

	results := vc.CalculateAdmissions()
	h1Admitted := getAdmittedStudentIDs(results, "H1")

	expectedH1 := []string{prepareStudentID("202")}
	if !reflect.DeepEqual(h1Admitted, expectedH1) {
		t.Errorf("H1: Expected quota student %v to be admitted, got %v", expectedH1, h1Admitted)
	}
}

func TestCalculateAdmissions_QuotaStudentsComparedByQuotaRating(t *testing.T) {
	vc := NewVarsityCalculator("TestUni")
	vc.AddHeading("H1", 1, "Advanced Quota Program")

	vc.AddApplication("H1", "301", 10, 1, CompetitionTargetQuota)
	vc.AddApplication("H1", "302", 5, 1, CompetitionTargetQuota)

	results := vc.CalculateAdmissions()
	h1Admitted := getAdmittedStudentIDs(results, "H1")

	expectedH1 := []string{prepareStudentID("302")}
	if !reflect.DeepEqual(h1Admitted, expectedH1) {
		t.Errorf("H1: Expected quota student %v with better quota rating to be admitted, got %v", expectedH1, h1Admitted)
	}
}

func TestCalculateAdmissions_QuotaStudentDisplacesNonQuota(t *testing.T) {
	vc := NewVarsityCalculator("TestUni")
	vc.AddHeading("H1", 1, "Competitive Program")

	vc.AddApplication("H1", "401", 5, 1, CompetitionRegular)

	vc.AddApplication("H1", "402", 8, 1, CompetitionTargetQuota)

	_ = vc.student("401")
	_ = vc.student("402")

	results := vc.CalculateAdmissions()
	h1Admitted := getAdmittedStudentIDs(results, "H1")

	expectedH1 := []string{prepareStudentID("402")}
	if !reflect.DeepEqual(h1Admitted, expectedH1) {
		t.Errorf("H1: Expected quota student %v to displace non-quota, got %v", expectedH1, h1Admitted)
	}
}

func TestCalculateAdmissions_HeadingFullWithQuotaNonQuotaCannotEnter(t *testing.T) {
	vc := NewVarsityCalculator("TestUni")
	vc.AddHeading("H1", 1, "Exclusive Program")

	vc.AddApplication("H1", "501", 1, 1, CompetitionTargetQuota)
	vc.AddApplication("H1", "502", 1, 1, CompetitionRegular)

	results := vc.CalculateAdmissions()
	h1Admitted := getAdmittedStudentIDs(results, "H1")

	expectedH1 := []string{prepareStudentID("501")}
	if !reflect.DeepEqual(h1Admitted, expectedH1) {
		t.Errorf("H1: Expected %v (quota) to hold the spot, got %v", expectedH1, h1Admitted)
	}
}

func TestCalculateAdmissions_MixedQuotaAndNonQuotaMultipleSpots(t *testing.T) {
	vc := NewVarsityCalculator("TestUni")
	vc.AddHeading("H1", 3, "Interdisciplinary Studies")

	vc.AddApplication("H1", "601", 5, 1, CompetitionTargetQuota)
	vc.AddApplication("H1", "602", 2, 1, CompetitionRegular)
	vc.AddApplication("H1", "603", 2, 1, CompetitionTargetQuota)
	vc.AddApplication("H1", "604", 8, 1, CompetitionRegular)
	vc.AddApplication("H1", "605", 1, 1, CompetitionRegular)

	results := vc.CalculateAdmissions()
	h1Admitted := getAdmittedStudentIDs(results, "H1")

	expectedH1 := []string{prepareStudentID("603"), prepareStudentID("601"), prepareStudentID("605")}

	if !reflect.DeepEqual(h1Admitted, expectedH1) {
		t.Errorf("H1: Expected admitted %v, got %v", expectedH1, h1Admitted)
	}
}

func TestCalculateAdmissions_QuotaStudentLosesToBetterQuotaStudent_ThenNonQuotaGetsSpot(t *testing.T) {
	vc := NewVarsityCalculator("TestUni")
	vc.AddHeading("H1", 2, "Specialized Engineering")

	vc.AddApplication("H1", "701", 10, 1, CompetitionTargetQuota)
	vc.AddApplication("H1", "702", 3, 1, CompetitionRegular)
	vc.AddApplication("H1", "703", 1, 1, CompetitionTargetQuota)

	results := vc.CalculateAdmissions()
	h1Admitted := getAdmittedStudentIDs(results, "H1")

	expectedH1 := []string{prepareStudentID("703"), prepareStudentID("701")}

	if !reflect.DeepEqual(h1Admitted, expectedH1) {
		t.Errorf("H1: Expected admitted %v, got %v", expectedH1, h1Admitted)
	}
}

func TestCalculateAdmissions_CompetitionTypePriority(t *testing.T) {
	vc := NewVarsityCalculator("TestUni")
	vc.AddHeading("H1", 5, "Competition Type Priority Test")

	// Add applications with different competition types
	// Lower rating place is better, but competition type takes precedence
	vc.AddApplication("H1", "801", 1, 1, CompetitionRegular)         // Best rating but lowest competition type
	vc.AddApplication("H1", "802", 10, 1, CompetitionTargetQuota)    // Worse rating but better competition type
	vc.AddApplication("H1", "803", 20, 1, CompetitionDedicatedQuota) // Even worse rating but even better competition type
	vc.AddApplication("H1", "804", 30, 1, CompetitionSpecialQuota)   // Worst rating but second-best competition type
	vc.AddApplication("H1", "805", 40, 1, CompetitionBVI)            // Worst rating but best competition type

	results := vc.CalculateAdmissions()
	h1Admitted := getAdmittedStudentIDs(results, "H1")

	// Students should be admitted in order of competition type (highest to lowest)
	expectedH1 := []string{
		prepareStudentID("805"), // CompetitionBVI (highest)
		prepareStudentID("804"), // CompetitionSpecialQuota
		prepareStudentID("803"), // CompetitionDedicatedQuota
		prepareStudentID("802"), // CompetitionTargetQuota
		prepareStudentID("801"), // CompetitionRegular (lowest)
	}

	if !reflect.DeepEqual(h1Admitted, expectedH1) {
		t.Errorf("H1: Expected admitted in competition type order %v, got %v", expectedH1, h1Admitted)
	}
}

func TestCalculateAdmissions_SameCompetitionTypeDifferentRatings(t *testing.T) {
	vc := NewVarsityCalculator("TestUni")
	vc.AddHeading("H1", 3, "Same Competition Type Test")

	// Add applications with the same competition type but different rating places
	vc.AddApplication("H1", "901", 30, 1, CompetitionSpecialQuota) // Worst rating
	vc.AddApplication("H1", "902", 20, 1, CompetitionSpecialQuota) // Middle rating
	vc.AddApplication("H1", "903", 10, 1, CompetitionSpecialQuota) // Best rating

	results := vc.CalculateAdmissions()
	h1Admitted := getAdmittedStudentIDs(results, "H1")

	// Students should be admitted in order of rating place (best to worst)
	expectedH1 := []string{
		prepareStudentID("903"), // Best rating (10)
		prepareStudentID("902"), // Middle rating (20)
		prepareStudentID("901"), // Worst rating (30)
	}

	if !reflect.DeepEqual(h1Admitted, expectedH1) {
		t.Errorf("H1: Expected admitted in rating order %v, got %v", expectedH1, h1Admitted)
	}
}
