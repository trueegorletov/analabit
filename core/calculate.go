package core

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

// Capacities hold the quotas and capacities of a heading for different competition types.
// If any quota is not used fully, the remaining places are available for Regular and BVI competitors.
// Max num of students who can be admitted to a heading is the sum of all quotas and the General capacity.
// Max num of students who can be admitted using each quota is limited by the corresponding field in Capacities.
// If a student fails to be admitted to a heading using their quota, they fail at all;
// quotas don't fall back to each other nor to the General capacity.
type Capacities struct {
	General        int // Guaranteed number of places for Regular & BVI competitors
	TargetQuota    int
	DedicatedQuota int
	SpecialQuota   int
}

func (c Capacities) String() string {
	return fmt.Sprintf("[S%d/D%d/T%d/G%d]",
		c.SpecialQuota, c.DedicatedQuota, c.TargetQuota, c.General,
	)
}

// PrintRvalue returns a string representation of the Capacities struct as a Go literal.
func (c Capacities) PrintRvalue() string {
	return fmt.Sprintf(`Capacities{
			General:        %d,
			TargetQuota:    %d,
			DedicatedQuota: %d,
			SpecialQuota:   %d,
		}`, c.General, c.TargetQuota, c.DedicatedQuota, c.SpecialQuota)
}

// maxInt returns the greater of two integers.
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type Competition int

// Competition represents the type of competition for a student's application.
const (
	CompetitionRegular Competition = iota
	CompetitionBVI
	CompetitionTargetQuota
	CompetitionDedicatedQuota
	CompetitionSpecialQuota
)

func (c Competition) String() string {
	switch c {
	case CompetitionRegular:
		return "Regular"
	case CompetitionBVI:
		return "BVI"
	case CompetitionTargetQuota:
		return "TargetQuota"
	case CompetitionDedicatedQuota:
		return "DedicatedQuota"
	case CompetitionSpecialQuota:
		return "SpecialQuota"
	default:
		return "UnknownCompetition"
	}
}

// Application represents a student's application to a specific heading.
type Application struct {
	// The student who made this application.
	student *Student
	// The heading to which the application is made.
	heading *Heading
	// The student's rating place for this application (lower is better).
	ratingPlace int
	// The priority of this application (e.g., 1st choice is 1; lower is higher priority).
	priority int
	// The type of competition for this application.
	competitionType Competition
}

// Student represents a student in the system.
type Student struct {
	mu sync.Mutex
	// Unique identifier of the student.
	Id string // Exported field
	// List of applications made by the student, sorted by priority (ascending, e.g., priority 1 first).
	applications []Application
	// If true, the student has withdrawn their application from this varsity and is ignored in calculations.
	quit bool
	// If true, the student has submitted their original documents for this varsity. Students who did can't be
	// thrown out when simulating percentage-drain of original certificates
	originalSubmitted bool
}

// application retrieves the student's highest priority application details for a specific heading.
// It panics if no application is found for the given heading for this student.
// WARNING: This might not be suitable for all contexts if a student has multiple applications
// to the same heading with different competition types/priorities and the specific one is needed.
func (s *Student) application(heading *Heading) Application {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, app := range s.applications {
		if app.heading == heading {
			return app
		}
	}
	panic("student has no application for the specified heading")
}

// addApplication adds a new application for the student and keeps the applications list sorted by priority.
func (s *Student) addApplication(heading *Heading, ratingPlace int, priority int, competitionType Competition) {
	s.mu.Lock()
	defer s.mu.Unlock()

	app := Application{
		student:         s, // Link student to application
		heading:         heading,
		ratingPlace:     ratingPlace,
		priority:        priority,
		competitionType: competitionType,
	}
	s.applications = append(s.applications, app)
	// Ensure applications are sorted by priority (ascending).
	sort.Slice(s.applications, func(i, j int) bool {
		return s.applications[i].priority < s.applications[j].priority
	})
}

// Heading represents a program or specialization within a varsity.
type Heading struct {
	// Unique identifier for the heading.
	code string
	// The varsity associated with this heading.
	varsity *VarsityCalculator
	// Maximum number of students that can be admitted to this heading using various quotas and general competition.
	capacities Capacities
	// A human-readable name or identifier for the heading.
	prettyName string
}

// Capacities returns the capacities of the heading, including quotas of different types.
func (h *Heading) Capacities() Capacities {
	return h.capacities
}

// PrettyName returns the human-readable name of the heading.
func (h *Heading) PrettyName() string {
	return h.prettyName
}

// Code returns the unique identifier of the heading.
func (h *Heading) Code() string {
	return h.code
}

// outscores determines if application app1 is better than application app2 for this heading.
// A lower ratingPlace is considered better.
// Applications of better competition type beat applications of any lower one.
func (h *Heading) outscores(app1 Application, app2 Application) bool {
	// First, compare competition types
	if app1.competitionType < app2.competitionType {
		return false
	}
	if app1.competitionType > app2.competitionType {
		return true
	}

	// If competition types are the same, compare by ratingPlace
	return app1.ratingPlace < app2.ratingPlace
}

// VarsityCalculator handles the admission calculation process for a university.

// CalculationResult holds the outcome of the admission process for a single heading.
type CalculationResult struct {
	// The heading for which this calculation result is generated.
	Heading *Heading // Changed to uppercase for export
	// List of students admitted to the heading, sorted according to admission criteria (e.g., quota, BVI, Regular, then rating).
	Admitted []*Student // Changed to uppercase for export
}

// HeadingAdmissionState tracks the admission status for a single heading.
type HeadingAdmissionState struct {
	heading         *Heading
	quotaAdmitted   map[Competition][]Application // Keyed by Target, Dedicated, Special; stores the applications.
	generalAdmitted []Application                 // BVI and Regular applications, sorted by BVI > Regular, then rating.
}

// NewHeadingAdmissionState creates a new state for a heading.
func NewHeadingAdmissionState(h *Heading) *HeadingAdmissionState {
	return &HeadingAdmissionState{
		heading: h,
		quotaAdmitted: map[Competition][]Application{
			CompetitionTargetQuota:    make([]Application, 0),
			CompetitionDedicatedQuota: make([]Application, 0),
			CompetitionSpecialQuota:   make([]Application, 0),
		},
		generalAdmitted: make([]Application, 0),
	}
}

// removeStudentFromOldPlacement removes a specific application from its recorded placement.
func removeStudentFromOldPlacement(
	appToRemove Application,
	admissionStates map[*Heading]*HeadingAdmissionState,
) {
	oldHeadingState, ok := admissionStates[appToRemove.heading]
	if !ok {
		// Should not happen if appToRemove was genuinely placed
		return
	}

	removed := false
	switch appToRemove.competitionType {
	case CompetitionTargetQuota, CompetitionDedicatedQuota, CompetitionSpecialQuota:
		oldAppList := oldHeadingState.quotaAdmitted[appToRemove.competitionType]
		newAppList := make([]Application, 0, len(oldAppList))
		for _, appInList := range oldAppList {
			// Compare by all key fields to ensure we remove the exact application
			if !(appInList.student == appToRemove.student &&
				appInList.heading == appToRemove.heading &&
				appInList.priority == appToRemove.priority &&
				appInList.competitionType == appToRemove.competitionType &&
				appInList.ratingPlace == appToRemove.ratingPlace) {
				newAppList = append(newAppList, appInList)
			} else {
				removed = true
			}
		}
		if removed {
			oldHeadingState.quotaAdmitted[appToRemove.competitionType] = newAppList
		}

	case CompetitionBVI, CompetitionRegular:
		oldAppList := oldHeadingState.generalAdmitted
		newAppList := make([]Application, 0, len(oldAppList))
		for _, appInList := range oldAppList {
			if !(appInList.student == appToRemove.student &&
				appInList.heading == appToRemove.heading &&
				appInList.priority == appToRemove.priority &&
				appInList.competitionType == appToRemove.competitionType &&
				appInList.ratingPlace == appToRemove.ratingPlace) {
				newAppList = append(newAppList, appInList)
			} else {
				removed = true
			}
		}
		if removed {
			oldHeadingState.generalAdmitted = newAppList
		}
	}
}

type VarsityCalculator struct {
	name     string
	students sync.Map // Stores *Student, keyed by student ID (string)
	headings sync.Map // Stores *Heading, keyed by heading code (string)
}

// NewVarsityCalculator creates a new varsity calculator.
func NewVarsityCalculator(name string) *VarsityCalculator {
	return &VarsityCalculator{
		name:     name,
		students: sync.Map{},
		headings: sync.Map{},
	}
}

// student returns an existing Student for the given ID, or creates and stores a new one if not found.
func (v *VarsityCalculator) student(id string) *Student {
	s, ok := v.students.Load(id)
	if !ok {
		newStudent := &Student{Id: id, applications: make([]Application, 0)} // Use exported Id
		s, _ = v.students.LoadOrStore(id, newStudent)
	}
	return s.(*Student)
}

// AddHeading adds a new heading to the varsity.
func (v *VarsityCalculator) AddHeading(code string, capacities Capacities, prettyName string) {
	heading := &Heading{
		code:       code,
		varsity:    v, // Link back to varsity
		capacities: capacities,
		prettyName: prettyName,
	}
	v.headings.Store(code, heading)
}

// AddApplication adds a student's application to a specific heading.
func (v *VarsityCalculator) AddApplication(headingCode string, studentID string, ratingPlace int, priority int, competitionType Competition) {
	h, ok := v.headings.Load(headingCode)
	if !ok {
		panic(fmt.Sprintf("heading with code %s not found", headingCode))
	}
	s := v.student(studentID)
	s.addApplication(h.(*Heading), ratingPlace, priority, competitionType)
}

// SetQuit marks a student as having quit.
func (v *VarsityCalculator) SetQuit(studentID string) {
	s := v.student(studentID)
	s.mu.Lock()
	s.quit = true
	s.mu.Unlock()
}

// SetOriginalSubmitted marks a student as submitting their original.
func (v *VarsityCalculator) SetOriginalSubmitted(studentID string) {
	s := v.student(studentID)
	s.mu.Lock()
	s.originalSubmitted = true
	s.mu.Unlock()
}

// GetStudents returns a slice of all students in the varsity.
// Useful for iterating over all students in a stable manner.
func (v *VarsityCalculator) GetStudents() []*Student {
	var students []*Student
	v.students.Range(func(key, value interface{}) bool {
		students = append(students, value.(*Student))
		return true
	})
	// Sort students by ID for deterministic behavior, although not strictly required by all logic
	sort.Slice(students, func(i, j int) bool {
		return students[i].Id < students[j].Id // Use exported Id
	})
	return students
}

// CalculateAdmissions processes student applications and determines admissions for each heading
// within the varsity. The algorithm works as follows:
//
// 1. Initialization:
//   - For each heading, an `HeadingAdmissionState` is created to track admitted applications
//     for different competition types (quotas and general).
//   - A map `studentBestPlacement` is initialized to store the best (highest priority, best rating/competition type)
//     application through which a student is currently admitted.
//
// 2. Iterative Admission Process:
//   - The core of the algorithm is an iterative loop that continues as long as any student's
//     admission status changes in a pass. This ensures that displacements and subsequent
//     re-placements are handled correctly until a stable state is reached.
//   - In each iteration, the algorithm processes all non-quit students.
//   - For each student, it iterates through their applications, sorted by priority (highest priority first).
//
// 3. Application Processing for a Student:
//   - If a student is already placed via an application of a strictly higher priority than the current
//     application being considered, the algorithm skips to the student's next application (or next student
//     if no more applications for this one).
//   - If the current application is for the same heading, priority, and competition type as the student's
//     existing best placement, it's skipped as they are already optimally placed for this specific application.
//
// 4. Admission Attempt for an Application:
//   - Quota Applications (Target, Dedicated, Special):
//   - The specific quota capacity for the heading is checked.
//   - If there's space, or if the current applicant outranks the worst-ranked applicant currently
//     admitted to that quota (by ratingPlace), the applicant can be placed.
//   - If placement occurs:
//   - If the student was previously placed via a different (lower priority or worse rank) application,
//     that old placement is vacated (removed from the respective heading/list, and `studentBestPlacement` updated for the displaced student).
//   - If this placement displaces another student from the quota, that displaced student is removed
//     from the quota list and their `studentBestPlacement` entry is cleared, making them eligible for
//     reconsideration in subsequent iterations or for their lower-priority applications.
//   - The current applicant is added to the quota list (sorted by ratingPlace), and their `studentBestPlacement` is updated.
//   - `madeChangeInIteration` is set to true.
//   - The algorithm breaks from processing this student's lower-priority applications, as they have found their best possible placement for now.
//   - General Competition Applications (BVI, Regular):
//   - The effective general capacity is calculated by summing the base general capacity and any unfilled spots
//     from all quota types for that heading.
//   - If there's space in the effective general capacity, or if the current applicant outscores
//     (BVI > Regular, then by ratingPlace) the worst applicant in the general admission list, they can be placed.
//   - Placement logic (vacating old spots, displacing others, updating lists and `studentBestPlacement`)
//     is similar to quota admissions, using `outscores` for comparison.
//   - `madeChangeInIteration` is set to true, and the algorithm breaks from this student's lower-priority apps.
//
// 5. Stabilization and Result Finalization:
//   - The iterative process (step 2-4) repeats until an entire pass over all students and their applications
//     results in no changes to any student's admission status (`madeChangeInIteration` remains false).
//   - Once stable, the final list of admitted students for each heading is compiled:
//   - All applications that resulted in admission (from `quotaAdmitted` and `generalAdmitted` in each `HeadingAdmissionState`)
//     are collected for each heading.
//   - These applications are then sorted using `outscores` to determine the final ranking order within the heading.
//   - The student objects are extracted from these sorted applications.
//   - The overall results are then sorted by heading code for deterministic output.
//
// Key Principles:
//   - Student Priority: Students always aim for their highest possible priority application.
//   - Separate Quotas: Quotas have distinct capacities. Failure in one quota does not grant access to another or general for that specific application.
//   - Unfilled Quotas: Unused quota spots augment the general competition capacity for that heading.
//   - Iterative Refinement: Displacements trigger re-evaluation, ensuring a globally (within the rules) optimal placement for students.
func (v *VarsityCalculator) CalculateAdmissions() []CalculationResult {
	admissionStates := make(map[*Heading]*HeadingAdmissionState)
	v.headings.Range(func(_, value interface{}) bool {
		h := value.(*Heading)
		admissionStates[h] = NewHeadingAdmissionState(h)
		return true
	})

	studentBestPlacement := make(map[*Student]Application) // Stores the Application object
	madeChangeInIteration := true

	for madeChangeInIteration {
		madeChangeInIteration = false
		allStudents := v.GetStudents() // Get a stable list for the iteration

		for _, student := range allStudents {
			student.mu.Lock()
			quitStatus := student.quit
			// Create a snapshot of applications to iterate over (they are sorted by priority)
			currentApplicationsSnapshot := make([]Application, len(student.applications))
			copy(currentApplicationsSnapshot, student.applications)
			student.mu.Unlock()

			if quitStatus {
				continue
			}

			for _, currentApp := range currentApplicationsSnapshot { // currentApp is a copy
				targetHeading := currentApp.heading
				currentHeadingState := admissionStates[targetHeading]

				// 1. Check if student is already better placed or equally placed for this app
				if existingPlacement, hasPlacement := studentBestPlacement[student]; hasPlacement {
					if existingPlacement.priority < currentApp.priority {
						break // Already in a strictly higher priority placement
					}
					// If currentApp is for the same heading, priority, and competition type,
					// they are already optimally placed regarding this specific application.
					if existingPlacement.priority == currentApp.priority &&
						existingPlacement.heading == targetHeading &&
						existingPlacement.competitionType == currentApp.competitionType {
						continue
					}
				}

				canBePlaced := false
				var displacedApp Application // Stores the application that gets displaced
				isDisplacedAppSet := false

				// 2. Attempt to place based on competition type
				switch currentApp.competitionType {
				case CompetitionTargetQuota, CompetitionDedicatedQuota, CompetitionSpecialQuota:
					quotaType := currentApp.competitionType
					admittedAppList := currentHeadingState.quotaAdmitted[quotaType]
					capacity := 0
					switch quotaType {
					case CompetitionTargetQuota:
						capacity = targetHeading.capacities.TargetQuota
					case CompetitionDedicatedQuota:
						capacity = targetHeading.capacities.DedicatedQuota
					case CompetitionSpecialQuota:
						capacity = targetHeading.capacities.SpecialQuota
					// Add a default case to handle other competition types, though logically
					// this switch is only entered for the explicit quota types above.
					default:
						// This case should ideally not be reached if the outer switch correctly routes only quota types here.
						// However, to satisfy the compiler warning about missing iota consts.
						capacity = 0
					}

					if len(admittedAppList) < capacity {
						canBePlaced = true
					} else if capacity > 0 { // Only check for displacement if quota has capacity > 0
						worstAdmittedApp := admittedAppList[len(admittedAppList)-1] // Assumes sorted by rating
						if currentApp.ratingPlace < worstAdmittedApp.ratingPlace {
							canBePlaced = true
							displacedApp = worstAdmittedApp
							isDisplacedAppSet = true
						}
					}

					if canBePlaced {
						// If student was already placed in a different application, remove that old placement
						if existingPlacement, hasPlacement := studentBestPlacement[student]; hasPlacement {
							// Check if it's truly a different placement to avoid unnecessary removal
							if existingPlacement.student != currentApp.student || // Should be same student
								existingPlacement.heading != currentApp.heading ||
								existingPlacement.priority != currentApp.priority || // Should be different if moving
								existingPlacement.competitionType != currentApp.competitionType {
								removeStudentFromOldPlacement(existingPlacement, admissionStates)
							}
						}

						if isDisplacedAppSet {
							// Remove displacedApp from its quota list
							newList := make([]Application, 0, len(admittedAppList))
							for _, appInList := range admittedAppList {
								if !(appInList.student == displacedApp.student && appInList.priority == displacedApp.priority) {
									newList = append(newList, appInList)
								}
							}
							currentHeadingState.quotaAdmitted[quotaType] = newList
							delete(studentBestPlacement, displacedApp.student) // Displaced student loses their spot
						}

						// Add current student's application to this quota
						currentHeadingState.quotaAdmitted[quotaType] = append(currentHeadingState.quotaAdmitted[quotaType], currentApp)
						sort.Slice(currentHeadingState.quotaAdmitted[quotaType], func(i, j int) bool {
							return currentHeadingState.quotaAdmitted[quotaType][i].ratingPlace < currentHeadingState.quotaAdmitted[quotaType][j].ratingPlace
						})
						studentBestPlacement[student] = currentApp
						madeChangeInIteration = true
						break // Student placed for this app, move to next student or stop processing this student's apps
					}

				case CompetitionBVI, CompetitionRegular:
					unfilledQuotaPlaces := 0
					unfilledQuotaPlaces += maxInt(0, targetHeading.capacities.TargetQuota-len(currentHeadingState.quotaAdmitted[CompetitionTargetQuota]))
					unfilledQuotaPlaces += maxInt(0, targetHeading.capacities.DedicatedQuota-len(currentHeadingState.quotaAdmitted[CompetitionDedicatedQuota]))
					unfilledQuotaPlaces += maxInt(0, targetHeading.capacities.SpecialQuota-len(currentHeadingState.quotaAdmitted[CompetitionSpecialQuota]))
					effectiveGeneralCapacity := targetHeading.capacities.General + unfilledQuotaPlaces

					generalAdmittedAppList := currentHeadingState.generalAdmitted

					if len(generalAdmittedAppList) < effectiveGeneralCapacity {
						canBePlaced = true
					} else if effectiveGeneralCapacity > 0 { // Only check for displacement if capacity > 0
						worstAdmittedApp := generalAdmittedAppList[len(generalAdmittedAppList)-1] // Assumes sorted by outscores
						if targetHeading.outscores(currentApp, worstAdmittedApp) {
							canBePlaced = true
							displacedApp = worstAdmittedApp
							isDisplacedAppSet = true
						}
					}

					if canBePlaced {
						if existingPlacement, hasPlacement := studentBestPlacement[student]; hasPlacement {
							if existingPlacement.student != currentApp.student ||
								existingPlacement.heading != currentApp.heading ||
								existingPlacement.priority != currentApp.priority ||
								existingPlacement.competitionType != currentApp.competitionType {
								removeStudentFromOldPlacement(existingPlacement, admissionStates)
							}
						}

						if isDisplacedAppSet {
							newList := make([]Application, 0, len(generalAdmittedAppList))
							for _, appInList := range generalAdmittedAppList {
								if !(appInList.student == displacedApp.student && appInList.priority == displacedApp.priority) {
									newList = append(newList, appInList)
								}
							}
							currentHeadingState.generalAdmitted = newList
							delete(studentBestPlacement, displacedApp.student)
						}

						currentHeadingState.generalAdmitted = append(currentHeadingState.generalAdmitted, currentApp)
						sort.Slice(currentHeadingState.generalAdmitted, func(i, j int) bool {
							return targetHeading.outscores(currentHeadingState.generalAdmitted[i], currentHeadingState.generalAdmitted[j])
						})
						studentBestPlacement[student] = currentApp
						madeChangeInIteration = true
						break // Student placed
					}
				} // End switch currentApp.competitionType

				// If student was placed via currentApp, their higher-priority applications are satisfied.
				// Break from iterating this student's lower-priority applications.
				if placement, ok := studentBestPlacement[student]; ok && placement.priority == currentApp.priority && placement.heading == currentApp.heading && placement.competitionType == currentApp.competitionType {
					break
				}
			} // End loop over student's applications
		} // End loop over all students
	} // End main iteration loop

	// Prepare final results
	results := make([]CalculationResult, 0)
	v.headings.Range(func(_, value interface{}) bool {
		h := value.(*Heading)
		state := admissionStates[h]

		var admittedAppsForHeading []Application
		// Collect all applications that resulted in admission for this heading
		for _, quotaAppList := range state.quotaAdmitted {
			admittedAppsForHeading = append(admittedAppsForHeading, quotaAppList...)
		}
		admittedAppsForHeading = append(admittedAppsForHeading, state.generalAdmitted...)

		// Sort all admitted applications for this heading by the unified outscores logic
		sort.Slice(admittedAppsForHeading, func(i, j int) bool {
			return h.outscores(admittedAppsForHeading[i], admittedAppsForHeading[j])
		})

		// Extract unique students from the sorted applications for the final list
		// studentBestPlacement ensures a student is in at most one list (or one spot),
		// so all apps in admittedAppsForHeading are from different students.
		finalAdmittedStudents := make([]*Student, 0, len(admittedAppsForHeading))
		for _, app := range admittedAppsForHeading {
			finalAdmittedStudents = append(finalAdmittedStudents, app.student)
		}

		results = append(results, CalculationResult{
			Heading:  h,
			Admitted: finalAdmittedStudents, // This list of students is now sorted based on their winning applications
		})
		return true
	})

	// Sort final results by heading code for deterministic output
	sort.Slice(results, func(i, j int) bool {
		return results[i].Heading.Code() < results[j].Heading.Code()
	})

	return results
}

// SimulateOriginalsDrain drains randomly selected drainPercent% of students who DID not submit their original but
// who hasn't already quit the varsity.
func (v *VarsityCalculator) SimulateOriginalsDrain(drainPercentage int) {
	if drainPercentage == 0 {
		return
	}

	if drainPercentage < 0 || drainPercentage > 100 {
		fmt.Printf("Invalid drain percentage: %d. Must be between 0 and 100.\n", drainPercentage)
		return
	}

	var eligibleStudents []*Student
	v.students.Range(func(key, value interface{}) bool {
		student := value.(*Student)
		student.mu.Lock()
		if !student.quit && !student.originalSubmitted {
			eligibleStudents = append(eligibleStudents, student)
		}
		student.mu.Unlock()
		return true
	})

	if len(eligibleStudents) == 0 {
		return
	}

	numToDrain := (len(eligibleStudents) * drainPercentage) / 100

	// Shuffle eligible students to pick randomly
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(eligibleStudents), func(i, j int) {
		eligibleStudents[i], eligibleStudents[j] = eligibleStudents[j], eligibleStudents[i]
	})

	for i := 0; i < numToDrain && i < len(eligibleStudents); i++ {
		studentToDrain := eligibleStudents[i]
		studentToDrain.mu.Lock()
		studentToDrain.quit = true
		studentToDrain.mu.Unlock()
	}
}
