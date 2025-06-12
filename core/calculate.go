package core

import (
	"fmt"
	"log" // Added import
	"log/slog"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"
)

// Capacities hold the quotas and capacities of a heading for different competition types.
// If any quota is not used fully, the remaining places are available for Regular and BVI competitors.
// Max num of students who can be admitted to a heading is the sum of all quotas and the Regular capacity.
// Max num of students who can be admitted using each quota is limited by the corresponding field in Capacities.
// If a student fails to be admitted to a heading using their quota, they fail at all;
// quotas don't fall back to each other nor to the Regular capacity.
type Capacities struct {
	Regular        int // Guaranteed number of places for Regular & BVI competitors
	TargetQuota    int
	DedicatedQuota int
	SpecialQuota   int
}

func (c Capacities) String() string {
	return fmt.Sprintf("[S%d/D%d/T%d/G%d]",
		c.SpecialQuota, c.DedicatedQuota, c.TargetQuota, c.Regular,
	)
}

// PrintRvalue returns a string representation of the Capacities struct as a Go literal.
func (c Capacities) PrintRvalue() string {
	return fmt.Sprintf(`Capacities{
			Regular:        %d,
			TargetQuota:    %d,
			DedicatedQuota: %d,
			SpecialQuota:   %d,
		}`, c.Regular, c.TargetQuota, c.DedicatedQuota, c.SpecialQuota)
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
	// The score of the student in this application, used for further uploading
	score int
}

func (a *Application) RatingPlace() int {
	return a.ratingPlace
}

func (a *Application) Priority() int {
	return a.priority
}

func (a *Application) CompetitionType() Competition {
	return a.competitionType
}

func (a *Application) Heading() *Heading {
	return a.heading
}

func (a *Application) StudentID() string {
	return a.student.id
}

func (a *Application) Score() int {
	return a.score
}

// Student represents a student in the system.
type Student struct {
	mu sync.Mutex
	// Unique identifier of the student.
	id string // Exported field
	// List of applications made by the student, sorted by priority (ascending, e.g., priority 1 first).
	applications []Application
	// If true, the student has withdrawn their application from this varsity and is ignored in calculations.
	quit bool
	// If true, the student has submitted their original documents for this varsity. Students who did can't be
	// thrown out when simulating percentage-drain of original certificates
	originalSubmitted bool
}

func (s *Student) Quit() bool {
	return s.quit
}

func (s *Student) Applications() []Application {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.applications
}

func (s *Student) ID() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.id // Exported field
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
func (s *Student) addApplication(heading *Heading, ratingPlace int, priority int, competitionType Competition, score int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	app := Application{
		student:         s, // Link student to application
		heading:         heading,
		ratingPlace:     ratingPlace,
		priority:        priority,
		competitionType: competitionType,
		score:           score,
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
	// A human-readable code or identifier for the heading.
	prettyName string
}

// Capacities returns the capacities of the heading, including quotas of different types.
func (h *Heading) Capacities() Capacities {
	return h.capacities
}

func (h *Heading) TotalCapacity() int {
	// Total capacity is the sum of all quotas and the Regular capacity.
	return h.capacities.Regular +
		h.capacities.TargetQuota +
		h.capacities.DedicatedQuota +
		h.capacities.SpecialQuota
}

// PrettyName returns the human-readable code of the heading.
func (h *Heading) PrettyName() string {
	return h.prettyName
}

// Code returns the unique identifier of the heading.
func (h *Heading) Code() string {
	return h.code
}

func (h *Heading) FullCode() string {
	return fmt.Sprintf("%s:%s", h.varsity.code, h.code)
}

func (h *Heading) VarsityCode() string {
	return h.varsity.code
}

func (h *Heading) VarsityPrettyName() string {
	return h.varsity.prettyName
}

// outscores determines if application app1 is better than application app2 for this heading.
// A lower ratingPlace is considered better.
// ApplicationsCache of better competition type beat applications of any lower one.
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

func (r *CalculationResult) PassingScore() (int, error) {
	if len(r.Admitted) == 0 {
		return 0, fmt.Errorf("no students admitted for heading %s", r.Heading.FullCode())
	}

	// return score of the last admitted student with heading code == r's heading code
	lastAdmitted := r.Admitted[len(r.Admitted)-1]

	for _, app := range lastAdmitted.applications {
		if app.heading.FullCode() == r.Heading.FullCode() {
			return app.score, nil
		}
	}

	return 0, fmt.Errorf("no needle application found for last admitted student %s in heading %s", lastAdmitted.ID(), r.Heading.FullCode())
}

func (r *CalculationResult) LastAdmittedRatingPlace() (int, error) {
	if len(r.Admitted) == 0 {
		return 0, fmt.Errorf("no students admitted for heading %s", r.Heading.FullCode())
	}

	// return ratingPlace of the last admitted student with heading code == r's heading code
	lastAdmitted := r.Admitted[len(r.Admitted)-1]

	for _, app := range lastAdmitted.applications {
		if app.heading.FullCode() == r.Heading.FullCode() {
			return app.ratingPlace, nil
		}
	}

	return 0, fmt.Errorf("no needle application found for last admitted student %s in heading %s", lastAdmitted.ID(), r.Heading.FullCode())
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
	mu         sync.Mutex // Added mutex for the calculator instance
	code       string
	prettyName string
	students   sync.Map // Stores *Student, keyed by student ID (string)
	headings   sync.Map // Stores *Heading, keyed by heading code (string)

	quitStudentsMu sync.RWMutex
	quitStudents   map[string]bool

	originalSubmittedStudentsMu sync.RWMutex
	originalSubmittedStudents   map[string]bool

	drainedPercent int
	wasted         bool
}

// NewVarsityCalculator creates a new varsity calculator.
func NewVarsityCalculator(code, prettyName string) *VarsityCalculator {
	return &VarsityCalculator{
		code:                      strings.TrimSpace(code),
		prettyName:                strings.TrimSpace(prettyName),
		students:                  sync.Map{},
		headings:                  sync.Map{},
		quitStudents:              make(map[string]bool),
		originalSubmittedStudents: make(map[string]bool),
	}
}

func (v *VarsityCalculator) checkNotWasted() {
	if v.wasted {
		panic("varsity calculator instance was already used for calculations, cannot be reused")
	}
}

// student returns an existing Student for the given ID, or creates and stores a new one if not found.
func (v *VarsityCalculator) student(id string) *Student {
	s, ok := v.students.Load(id)
	if !ok {
		newStudent := &Student{id: id, applications: make([]Application, 0)} // Use exported ID
		s, _ = v.students.LoadOrStore(id, newStudent)
	}
	return s.(*Student)
}

func (v *VarsityCalculator) GetHeading(code string) *Heading {
	code = strings.TrimSpace(code)
	h, ok := v.headings.Load(code)
	if !ok {
		return nil // Heading not found
	}
	return h.(*Heading) // Return the heading if found
}

// AddHeading adds a new heading to the varsity.
func (v *VarsityCalculator) AddHeading(code string, capacities Capacities, prettyName string) {
	code = strings.TrimSpace(code)
	prettyName = strings.TrimSpace(prettyName)

	heading := &Heading{
		code:       code,
		varsity:    v, // Link back to varsity
		capacities: capacities,
		prettyName: prettyName,
	}
	v.headings.Store(code, heading)
}

// AddApplication adds a student's application to a specific heading.
func (v *VarsityCalculator) AddApplication(headingCode string, studentID string, ratingPlace int, priority int, competitionType Competition, score int) {
	headingCode = strings.TrimSpace(headingCode)

	h, ok := v.headings.Load(headingCode)
	if !ok {
		panic(fmt.Sprintf("heading with code %s not found", headingCode))
	}
	s := v.student(studentID)
	s.addApplication(h.(*Heading), ratingPlace, priority, competitionType, score)
}

// SetQuit marks a student as having quit.
func (v *VarsityCalculator) SetQuit(studentID string) {
	s := v.student(studentID)
	s.mu.Lock()
	s.quit = true
	s.mu.Unlock()

	v.quitStudentsMu.Lock()
	v.quitStudents[studentID] = true
	v.quitStudentsMu.Unlock()
}

// SetOriginalSubmitted marks a student as submitting their original.
func (v *VarsityCalculator) SetOriginalSubmitted(studentID string) {
	s := v.student(studentID)
	s.mu.Lock()
	s.originalSubmitted = true
	s.mu.Unlock()

	v.originalSubmittedStudentsMu.Lock()
	v.originalSubmittedStudents[studentID] = true
	v.originalSubmittedStudentsMu.Unlock()
}

// Students returns a slice of all students in the varsity.
// Useful for iterating over all students in a stable manner.
func (v *VarsityCalculator) Students() []*Student {
	var students []*Student
	v.students.Range(func(key, value interface{}) bool {
		students = append(students, value.(*Student))
		return true
	})
	// Sort students by ID for deterministic behavior, although not strictly required by all logic
	sort.Slice(students, func(i, j int) bool {
		return students[i].id < students[j].id // Use exported ID
	})
	return students
}

func (v *VarsityCalculator) GetStudent(studentID string) *Student {
	studentID = strings.TrimSpace(studentID)
	s, ok := v.students.Load(studentID)
	if !ok {
		return nil // Student not found
	}
	return s.(*Student) // Return the student if found
}

func (v *VarsityCalculator) ForEachQuit(f func(studentID string)) {
	v.quitStudentsMu.RLock()

	for id, val := range v.quitStudents {
		if !val {
			slog.Warn("Found non-quit student in quitStudents map", "studentID", id)
			continue
		}

		f(id)
	}

	v.quitStudentsMu.RUnlock()
}

func (v *VarsityCalculator) ForEachOriginalSubmitted(f func(studentID string)) {
	v.originalSubmittedStudentsMu.RLock()

	for id, val := range v.originalSubmittedStudents {
		if !val {
			slog.Warn("Found non-original-submitted student in originalSubmittedStudents map", "studentID", id)
			continue
		}

		f(id)
	}

	v.originalSubmittedStudentsMu.RUnlock()
}

func (v *VarsityCalculator) DrainedPercent() int {
	return v.drainedPercent
}

func (v *VarsityCalculator) Headings() []*Heading {
	var headings []*Heading
	v.headings.Range(func(_, value interface{}) bool {
		headings = append(headings, value.(*Heading))
		return true
	})
	// Sort headings by code for deterministic behavior
	sort.Slice(headings, func(i, j int) bool {
		return headings[i].prettyName < headings[j].prettyName
	})
	return headings
}

// CalculateAdmissions performs the main admission calculation logic for the varsity.
// It processes all students and their applications, considering quotas and competition types,
// to determine which students are admitted to which headings.
// The current implementation involves iterating through students, then their applications,
// attempting admission which may involve sorting lists of admitted candidates for each quota/general competition
// and potentially re-evaluating other applications if a student is displaced.
// Complexity: Potentially high, roughly S * A * (C_log_C + D), where S is students, A is applications per student,
// C is heading capacity (due to sorting admitted lists), and D is cost of handling displacements (which can be recursive).
// For N total applications, and H headings, sorting final admitted lists per heading adds H * C_log_C.
func (v *VarsityCalculator) CalculateAdmissions() []CalculationResult {
	log.Println("CalculateAdmissions: Starting...")
	v.mu.Lock() // Use the new mutex
	log.Println("CalculateAdmissions: VarsityCalculator mutex locked.")
	defer func() {
		v.mu.Unlock() // Use the new mutex
		log.Println("CalculateAdmissions: VarsityCalculator mutex unlocked.")
	}()

	if v.wasted {
		panic(fmt.Errorf("attempt to use a wasted calculator for varsity %s", v.code))
	}
	defer func() { v.wasted = true }()

	// Initialize admission states for each heading
	admissionStates := make(map[*Heading]*HeadingAdmissionState) // Corrected type to HeadingAdmissionState
	log.Println("CalculateAdmissions: Initializing admission states for headings...")
	headingsCount := 0
	for _, h := range v.Headings() { // Corrected: Call Headings() method
		admissionStates[h] = NewHeadingAdmissionState(h) // Corrected: Use NewHeadingAdmissionState(h)
		headingsCount++
	}
	log.Printf("CalculateAdmissions: Initialized admission states for %d headings.", headingsCount)

	// studentBestPlacement tracks the best (highest priority) heading a student is currently admitted to.
	// Key: student ID, Value: *admissionInfo
	studentBestPlacement := make(map[string]*admissionInfo) // Assuming admissionInfo is defined elsewhere or should be defined.

	// Process all students
	// Convert sync.Map to slice for deterministic processing and length calculation
	var studentsToProcess []*Student
	v.students.Range(func(key, value interface{}) bool {
		studentsToProcess = append(studentsToProcess, value.(*Student))
		return true
	})
	// Optionally sort studentsToProcess by ID if deterministic order is critical here
	// sort.Slice(studentsToProcess, func(i, j int) bool { return studentsToProcess[i].ID() < studentsToProcess[j].ID() })

	log.Printf("CalculateAdmissions: Starting to process %d students...", len(studentsToProcess))
	studentProcessingTimes := make(map[string]time.Duration)

	for i, student := range studentsToProcess { // Iterate over the slice
		studentProcessingStart := time.Now()
		log.Printf("CalculateAdmissions: Processing student %d/%d ID: %s (Quit: %t, OriginalSubmitted: %t, Apps: %d)",
			i+1, len(studentsToProcess), student.ID(), student.Quit(), student.originalSubmitted, len(student.Applications()))

		if student.Quit() {
			log.Printf("CalculateAdmissions: Student %s quit, skipping.", student.ID())
			continue
		}
		v.processStudentApplications(student, admissionStates, studentBestPlacement)
		studentProcessingTimes[student.ID()] = time.Since(studentProcessingStart)
		log.Printf("CalculateAdmissions: Finished processing student %s in %v.", student.ID(), studentProcessingTimes[student.ID()])
	}
	log.Println("CalculateAdmissions: Finished processing all students.")

	// Collect results
	log.Println("CalculateAdmissions: Collecting results...")
	var results []CalculationResult
	for _, h := range v.Headings() { // Corrected: Iterate over the slice returned by v.Headings()
		log.Printf("CalculateAdmissions: Collecting results for heading %s (%s)", h.PrettyName(), h.Code())
		state := admissionStates[h]
		if state == nil { // Add a nil check for safety, though it shouldn't happen if initialized correctly
			log.Printf("CalculateAdmissions: WARNING - nil admissionState for heading %s (%s). Skipping.", h.PrettyName(), h.Code())
			continue
		}

		var admittedAppsForHeading []Application
		// Collect all applications that resulted in admission for this heading
		for competitionType, quotaAppList := range state.quotaAdmitted {
			log.Printf("CalculateAdmissions: Heading %s, Quota %s, Admitted count: %d", h.Code(), competitionType, len(quotaAppList))
			admittedAppsForHeading = append(admittedAppsForHeading, quotaAppList...)
		}
		log.Printf("CalculateAdmissions: Heading %s, General Admitted count: %d", h.Code(), len(state.generalAdmitted))
		admittedAppsForHeading = append(admittedAppsForHeading, state.generalAdmitted...)
		log.Printf("CalculateAdmissions: Heading %s, Total admitted apps before final sort: %d", h.Code(), len(admittedAppsForHeading))

		// Sort all admitted applications for this heading by the unified outscores logic
		log.Printf("CalculateAdmissions: Sorting %d admitted applications for heading %s...", len(admittedAppsForHeading), h.Code())
		sortStart := time.Now()
		sort.Slice(admittedAppsForHeading, func(i, j int) bool {
			return h.outscores(admittedAppsForHeading[i], admittedAppsForHeading[j])
		})
		log.Printf("CalculateAdmissions: Sorted admitted applications for heading %s in %v.", h.Code(), time.Since(sortStart))

		// Extract unique students from the sorted applications for the final list
		finalAdmittedStudents := make([]*Student, 0, len(admittedAppsForHeading))
		for _, app := range admittedAppsForHeading {
			finalAdmittedStudents = append(finalAdmittedStudents, app.student)
		}
		log.Printf("CalculateAdmissions: Heading %s, Final admitted student count: %d", h.Code(), len(finalAdmittedStudents))

		results = append(results, CalculationResult{
			Heading:  h,
			Admitted: finalAdmittedStudents,
		})
	}

	// Sort final results by heading code for deterministic output
	sort.Slice(results, func(i, j int) bool {
		return results[i].Heading.Code() < results[j].Heading.Code()
	})

	return results
}

// processStudentApplications processes a single student's applications.
func (v *VarsityCalculator) processStudentApplications(
	student *Student,
	admissionStates map[*Heading]*HeadingAdmissionState,
	studentBestPlacement map[string]*admissionInfo,
) {
	log.Printf("processStudentApplications: Starting for student %s.", student.ID())
	// Iterate through the student's applications, sorted by priority
	apps := student.Applications() // Applications are already sorted by priority
	log.Printf("processStudentApplications: Student %s has %d applications.", student.ID(), len(apps))

	for i, app := range apps {
		log.Printf("processStudentApplications: Student %s, processing app %d/%d to heading %s (Priority: %d, Type: %s, Rating: %d)",
			student.ID(), i+1, len(apps), app.Heading().Code(), app.Priority(), app.CompetitionType(), app.RatingPlace())

		// Check if student can be admitted to this heading based on current best placement
		if !v.canAdmit(student, app, studentBestPlacement) {
			log.Printf("processStudentApplications: Student %s cannot be admitted to heading %s due to better current placement or same priority.", student.ID(), app.Heading().Code())
			continue // Student is already admitted to a higher or same priority heading
		}

		state := admissionStates[app.Heading()]
		if state == nil { // Add a nil check for safety
			log.Printf("processStudentApplications: WARNING - nil admissionState for heading %s (%s) for student %s. Skipping app.", app.Heading().PrettyName(), app.Heading().Code(), student.ID())
			continue
		}

		var admitted bool
		var displacedStudent *Student
		var displacedApplication Application // To store the actual application that was displaced

		switch app.CompetitionType() {
		case CompetitionTargetQuota, CompetitionDedicatedQuota, CompetitionSpecialQuota:
			log.Printf("processStudentApplications: Student %s, attempting quota admission to %s for type %s.", student.ID(), app.Heading().Code(), app.CompetitionType())
			admitted, displacedStudent, displacedApplication = v.tryAdmitQuota(app, state, studentBestPlacement)
		case CompetitionRegular, CompetitionBVI: // BVI is treated like regular but with higher "score"
			log.Printf("processStudentApplications: Student %s, attempting general admission to %s for type %s.", student.ID(), app.Heading().Code(), app.CompetitionType())
			admitted, displacedStudent, displacedApplication = v.tryAdmitGeneral(app, state, studentBestPlacement)
		default:
			log.Printf("processStudentApplications: Student %s, unknown competition type %s for app to %s.", student.ID(), app.CompetitionType(), app.Heading().Code())
			continue
		}

		if admitted {
			log.Printf("processStudentApplications: Student %s ADMITTED to heading %s via app (Priority: %d, Type: %s). Updating best placement.",
				student.ID(), app.Heading().Code(), app.Priority(), app.CompetitionType())
			// Update student's best placement
			// Before updating, if this student was already in studentBestPlacement for a *different, lower-priority* heading,
			// they need to be removed from that old heading's admission list.
			if oldPlacement, ok := studentBestPlacement[student.ID()]; ok && oldPlacement.app.Heading() != app.Heading() {
				log.Printf("processStudentApplications: Student %s was previously admitted to %s (Priority: %d), now admitted to %s (Priority: %d). Removing from old heading.",
					student.ID(), oldPlacement.app.Heading().Code(), oldPlacement.app.Priority(), app.Heading().Code(), app.Priority())
				removeStudentFromOldPlacement(oldPlacement.app, admissionStates)
			}
			studentBestPlacement[student.ID()] = &admissionInfo{app: app, headingState: state}

			if displacedStudent != nil {
				log.Printf("processStudentApplications: Student %s was displaced from heading %s by student %s. Finding new placement for %s.",
					displacedStudent.ID(), displacedApplication.Heading().Code(), student.ID(), displacedStudent.ID())
				// If a student was displaced, try to find a new best placement for them
				// The student is already removed from the specific heading's list by tryAdmitQuota/tryAdmitGeneral.
				// We need to remove them from studentBestPlacement map as their previous best is gone.
				delete(studentBestPlacement, displacedStudent.ID())
				v.processStudentApplications(displacedStudent, admissionStates, studentBestPlacement) // Recursive call
			}
		} else {
			log.Printf("processStudentApplications: Student %s could not be admitted to heading %s (Quota: %s, General: %s).",
				student.ID(), app.Heading().Code(), app.competitionType, app.competitionType)
		}
	}
}

// Returns true if admitted, the displaced student (if any), and the displaced application (if any).
func (v *VarsityCalculator) tryAdmitQuota(
	app Application,
	state *HeadingAdmissionState,
	studentBestPlacement map[string]*admissionInfo,
) (bool, *Student, Application) {
	log.Printf("tryAdmitQuota: Student %s, App to %s (Rating: %d, Type: %s)",
		app.student.ID(), app.heading.Code(), app.ratingPlace, app.competitionType)

	quotaCapacity := 0
	switch app.competitionType {
	case CompetitionTargetQuota:
		quotaCapacity = state.heading.Capacities().TargetQuota
	case CompetitionDedicatedQuota:
		quotaCapacity = state.heading.Capacities().DedicatedQuota
	case CompetitionSpecialQuota:
		quotaCapacity = state.heading.Capacities().SpecialQuota
	default:
		log.Printf("tryAdmitQuota: Unknown quota type %s for student %s", app.competitionType, app.student.ID())
		return false, nil, Application{}
	}

	log.Printf("tryAdmitQuota: Heading %s, Quota Type %s, Capacity: %d, Currently Admitted in this quota: %d",
		app.heading.Code(), app.competitionType, quotaCapacity, len(state.quotaAdmitted[app.competitionType]))

	// **** THIS IS THE CRUCIAL FIX for tryAdmitQuota ****
	if quotaCapacity == 0 {
		log.Printf("tryAdmitQuota: Quota capacity for type %s is 0 for heading %s. Cannot admit student %s.", app.competitionType, app.heading.Code(), app.student.ID())
		return false, nil, Application{}
	}
	// **** END OF FIX ****

	currentAdmittedInQuota := state.quotaAdmitted[app.competitionType]

	if len(currentAdmittedInQuota) < quotaCapacity {
		// Quota not full, admit student
		log.Printf("tryAdmitQuota: Admitting student %s to %s (Type: %s) - quota not full.", app.student.ID(), app.heading.Code(), app.competitionType)
		state.quotaAdmitted[app.competitionType] = append(currentAdmittedInQuota, app)
		// Sort by ratingPlace (lower is better)
		sort.Slice(state.quotaAdmitted[app.competitionType], func(i, j int) bool {
			return state.quotaAdmitted[app.competitionType][i].ratingPlace < state.quotaAdmitted[app.competitionType][j].ratingPlace
		})
		return true, nil, Application{}
	}

	// Quota is full, check if this student outscores the worst student in quota
	worstAdmittedApp := currentAdmittedInQuota[len(currentAdmittedInQuota)-1] // Assumes sorted by ratingPlace, worst is last
	if app.ratingPlace < worstAdmittedApp.ratingPlace {
		// This student is better than the worst admitted student
		log.Printf("tryAdmitQuota: Student %s (Rating: %d) outscores worst student %s (Rating: %d) in quota %s for %s. Admitting and displacing.",
			app.student.ID(), app.ratingPlace, worstAdmittedApp.student.ID(), worstAdmittedApp.ratingPlace, app.competitionType, app.heading.Code())

		displacedStudent := worstAdmittedApp.student
		// Remove worst student
		state.quotaAdmitted[app.competitionType] = currentAdmittedInQuota[:len(currentAdmittedInQuota)-1]
		// Add new student
		state.quotaAdmitted[app.competitionType] = append(state.quotaAdmitted[app.competitionType], app)
		// Re-sort
		sort.Slice(state.quotaAdmitted[app.competitionType], func(i, j int) bool {
			return state.quotaAdmitted[app.competitionType][i].ratingPlace < state.quotaAdmitted[app.competitionType][j].ratingPlace
		})
		// Remove displaced student from their overall best placement
		if currentBest, ok := studentBestPlacement[displacedStudent.ID()]; ok && currentBest.app.student.ID() == displacedStudent.ID() && currentBest.app.Heading() == worstAdmittedApp.Heading() {
			log.Printf("tryAdmitQuota: Removing displaced student %s from their best placement map for heading %s.", displacedStudent.ID(), worstAdmittedApp.Heading().Code())
			delete(studentBestPlacement, displacedStudent.ID())
		}
		return true, displacedStudent, worstAdmittedApp // Return the displaced application
	}
	log.Printf("tryAdmitQuota: Student %s (Rating: %d) does not outscore worst student %s (Rating: %d) in quota %s for %s. Not admitting.",
		app.student.ID(), app.ratingPlace, worstAdmittedApp.student.ID(), worstAdmittedApp.ratingPlace, app.competitionType, app.heading.Code())
	return false, nil, Application{} // Return empty Application if no displacement
}

// tryAdmitGeneral attempts to admit a student to general competition (Regular or BVI).
// Returns true if admitted, the displaced student (if any), and the displaced application (if any).
func (v *VarsityCalculator) tryAdmitGeneral(
	app Application,
	state *HeadingAdmissionState,
	studentBestPlacement map[string]*admissionInfo,
) (bool, *Student, Application) { // Added displacedApplication to return
	log.Printf("tryAdmitGeneral: Student %s, App to %s (Rating: %d, Type: %s), General Capacity: %d, Currently Admitted: %d",
		app.student.ID(), app.heading.Code(), app.ratingPlace, app.competitionType,
		state.heading.Capacities().Regular, // BVI uses Regular capacity
		len(state.generalAdmitted))

	generalCapacity := state.heading.Capacities().Regular

	// **** THIS IS THE CRUCIAL FIX ****
	if generalCapacity == 0 {
		log.Printf("tryAdmitGeneral: General capacity is 0 for heading %s. Cannot admit student %s.", app.heading.Code(), app.student.ID())
		return false, nil, Application{}
	}
	// **** END OF FIX ****

	if len(state.generalAdmitted) < generalCapacity {
		// General capacity not full, admit student
		log.Printf("tryAdmitGeneral: Admitting student %s to %s (Type: %s) - general capacity not full.", app.student.ID(), app.heading.Code(), app.competitionType)
		state.generalAdmitted = append(state.generalAdmitted, app)
		// Sort by BVI > Regular, then ratingPlace
		sort.Slice(state.generalAdmitted, func(i, j int) bool {
			return state.heading.outscores(state.generalAdmitted[i], state.generalAdmitted[j])
		})
		return true, nil, Application{}
	}

	// General capacity is full, check if this student outscores the worst student
	// The list is sorted by outscores, so the "worst" is the last one.
	worstAdmittedApp := state.generalAdmitted[len(state.generalAdmitted)-1]
	if state.heading.outscores(app, worstAdmittedApp) {
		// This student is better
		log.Printf("tryAdmitGeneral: Student %s (App: %v) outscores worst student %s (App: %v) in general for %s. Admitting and displacing.",
			app.student.ID(), app, worstAdmittedApp.student.ID(), worstAdmittedApp, app.heading.Code())
		displacedStudent := worstAdmittedApp.student
		// Remove worst student
		state.generalAdmitted = state.generalAdmitted[:len(state.generalAdmitted)-1]
		// Add new student
		state.generalAdmitted = append(state.generalAdmitted, app)
		// Re-sort
		sort.Slice(state.generalAdmitted, func(i, j int) bool {
			return state.heading.outscores(state.generalAdmitted[i], state.generalAdmitted[j])
		})
		// Remove displaced student from their overall best placement
		if currentBest, ok := studentBestPlacement[displacedStudent.ID()]; ok && currentBest.app.student.ID() == displacedStudent.ID() && currentBest.app.Heading() == worstAdmittedApp.Heading() {
			log.Printf("tryAdmitGeneral: Removing displaced student %s from their best placement map for heading %s.", displacedStudent.ID(), worstAdmittedApp.Heading().Code())
			delete(studentBestPlacement, displacedStudent.ID())
		}
		return true, displacedStudent, worstAdmittedApp // Return the displaced application
	}
	log.Printf("tryAdmitGeneral: Student %s (App: %v) does not outscore worst student %s (App: %v) in general for %s. Not admitting.",
		app.student.ID(), app, worstAdmittedApp.student.ID(), worstAdmittedApp, app.heading.Code())
	return false, nil, Application{} // Return empty Application if no displacement
}

// canAdmit checks if a student can be admitted to a heading based on their current best placement.
// A student can be admitted if the new application's heading has higher priority than their current best,
// or if they have no current best placement.
func (v *VarsityCalculator) canAdmit(
	student *Student,
	newApp Application,
	studentBestPlacement map[string]*admissionInfo, // Assuming admissionInfo is defined
) bool {
	currentBest, ok := studentBestPlacement[student.ID()]
	if !ok {
		log.Printf("canAdmit: Student %s has no current best placement. Can attempt admission to %s (Priority: %d).", student.ID(), newApp.heading.Code(), newApp.priority)
		return true // No current placement, can always try
	}
	// Student has a current placement. Can only admit to newApp if newApp has strictly higher priority.
	// Lower priority number means higher priority.
	can := newApp.priority < currentBest.app.priority
	if can {
		log.Printf("canAdmit: Student %s current best is %s (Priority: %d). New app to %s (Priority: %d) is HIGHER priority. Can attempt.",
			student.ID(), currentBest.app.heading.Code(), currentBest.app.priority, newApp.heading.Code(), newApp.priority)
	} else {
		log.Printf("canAdmit: Student %s current best is %s (Priority: %d). New app to %s (Priority: %d) is LOWER or SAME priority. Cannot attempt.",
			student.ID(), currentBest.app.heading.Code(), currentBest.app.priority, newApp.heading.Code(), newApp.priority)
	}
	return can
}

// This struct was assumed by the logging I added. It needs to be defined.
// It stores information about a student's admission to a particular heading.
type admissionInfo struct {
	app          Application            // The specific application that led to admission
	headingState *HeadingAdmissionState // The state of the heading they were admitted to
}

// Helper for Capacities to get quota capacity by type, assuming it might be useful.
// This is not part of the original code but added for clarity in tryAdmitQuota.
func (c Capacities) quotaCapacity(compType Competition) int {
	switch compType {
	case CompetitionTargetQuota:
		return c.TargetQuota
	case CompetitionDedicatedQuota:
		return c.DedicatedQuota
	case CompetitionSpecialQuota:
		return c.SpecialQuota
	}
	return 0
}

// SimulateOriginalsDrain drains randomly selected drainPercent% of students who DID not submit their original but
// who hasn't already quit the varsity.
func (v *VarsityCalculator) SimulateOriginalsDrain(drainPercent int) {
	v.checkNotWasted()

	if drainPercent == 0 {
		return
	}

	if drainPercent < 0 || drainPercent > 100 {
		slog.Warn("Invalid value, must be in [0, 100]", "drainPercent", drainPercent)
		return
	}

	if v.drainedPercent != 0 {
		slog.Warn("Already drained students, cannot drain again", "drainedPercent", v.drainedPercent)
		return
	}

	nonQuitCount := 0

	var eligibleStudents []*Student
	v.students.Range(func(key, value interface{}) bool {
		student := value.(*Student)
		student.mu.Lock()
		if !student.quit {
			nonQuitCount += 1
		}

		if !student.quit && !student.originalSubmitted {
			eligibleStudents = append(eligibleStudents, student)
		}
		student.mu.Unlock()
		return true
	})

	if len(eligibleStudents) == 0 {
		return
	}

	numToDrain := (len(eligibleStudents) * drainPercent) / 100

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

	v.drainedPercent = drainPercent
}
