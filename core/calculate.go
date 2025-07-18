package core

import (
	"bytes"
	"container/heap"
	"container/list" // Added for O(1) queue operations
	"encoding/gob"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/trueegorletov/analabit/core/utils"
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

// Helper for Capacities to get quota capacity by type.
func (c Capacities) quotaCapacity(compType Competition) int {
	switch compType {
	case CompetitionTargetQuota:
		return c.TargetQuota
	case CompetitionDedicatedQuota:
		return c.DedicatedQuota
	case CompetitionSpecialQuota:
		return c.SpecialQuota
	default:
		// This case should ideally not be reached if compType is a valid quota type.
		// Return 0 or handle as an error, depending on desired strictness.
		slog.Warn("quotaCapacity called with non-quota competition type", "type", compType.String())
		return 0
	}
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
	return a.student.IDValue
}

func (a *Application) Score() int {
	return a.score
}

// Student represents a student in the system.
type Student struct {
	mu sync.Mutex
	// Unique identifier of the student. Exported to allow gob encoding.
	IDValue string
	// List of applications made by the student, sorted by priority (ascending, e.g., priority 1 first).
	applications []*Application
	// If true, the student has withdrawn their application from this varsity and is ignored in calculations.
	quit bool
	// If true, the student has submitted their original documents for this varsity. Students who did can't be
	// thrown out when simulating percentage-drain of original certificates
	originalSubmitted bool
}

func (s *Student) Quit() bool {
	return s.quit
}

func (s *Student) Applications() []*Application {
	return s.applications
}

func (s *Student) ID() string {
	return s.IDValue
}

// application retrieves the student's highest priority application details for a specific heading.
// It panics if no application is found for the given heading for this student.
// WARNING: This might not be suitable for all contexts if a student has multiple applications
// to the same heading with different competition types/priorities and the specific one is needed.
func (s *Student) application(heading *Heading) *Application {
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

	app := Application{
		student:         s, // Link student to application
		heading:         heading,
		ratingPlace:     ratingPlace,
		priority:        priority,
		competitionType: competitionType,
		score:           score,
	}

	defer func() {
		// Ensure applications are sorted by priority (ascending).
		sort.Slice(s.applications, func(i, j int) bool {
			return s.applications[i].priority < s.applications[j].priority
		})

		s.mu.Unlock()
	}()

	// Ensure we always have only one application per heading for each student with best competition type.

	for i, existingApp := range s.applications {
		if existingApp.heading.Code() == heading.Code() {
			// If we already have an application for this heading, check if the new one is better.
			if competitionPrecedence(competitionType) > competitionPrecedence(existingApp.competitionType) {
				// Replace the existing application with the new one.
				s.applications[i] = &app

				return
			} else if competitionPrecedence(competitionType) == competitionPrecedence(existingApp.competitionType) &&
				ratingPlace < existingApp.ratingPlace {
				// If competition types are equal, keep the one with lower ratingPlace.
				s.applications[i] = &app

				return
			} else if competitionPrecedence(competitionType) == competitionPrecedence(existingApp.competitionType) && score > existingApp.score {
				// If competition types and ratingPlace are equal, keep the one with higher score.
				s.applications[i] = &app

				return
			}
			// Otherwise, do not add this application as it is worse than the existing one.
			return
		}
	}

	s.applications = append(s.applications, &app)
}

// Heading represents a program or specialization within a varsity.
type Heading struct {
	// Unique identifier for the heading. Exported so that encoding/gob can access it.
	CodeValue string
	// The varsity associated with this heading (kept unexported to avoid deep gob encoding).
	varsity *VarsityCalculator
	// Maximum number of students that can be admitted to this heading using various quotas and general competition.
	CapacitiesValue Capacities
	// A human-readable code or identifier for the heading.
	PrettyNameValue string

	// Cached vars for serialization. They are skipped by gob as we implement custom encoding but kept for runtime access.
	varsityCodeCached       string
	varsityPrettyNameCached string
}

// Capacities returns the capacities of the heading, including quotas of different types.
func (h *Heading) Capacities() Capacities {
	return h.CapacitiesValue
}

func (h *Heading) TotalCapacity() int {
	// Total capacity is the sum of all quotas and the Regular capacity.
	return h.CapacitiesValue.Regular +
		h.CapacitiesValue.TargetQuota +
		h.CapacitiesValue.DedicatedQuota +
		h.CapacitiesValue.SpecialQuota
}

// PrettyName returns the human-readable name of the heading.
func (h *Heading) PrettyName() string {
	return h.PrettyNameValue
}

// Code returns the unique identifier of the heading.
func (h *Heading) Code() string {
	return h.CodeValue
}

func (h *Heading) FullCode() string {
	// Provide safe fallback if varsity pointer is nil (e.g., after gob decoding).
	if h.varsity == nil {
		return h.CodeValue // fall back to heading code only
	}
	return fmt.Sprintf("%s:%s", h.varsity.code, h.CodeValue)
}

func (h *Heading) VarsityCode() string {
	if h.varsity != nil {
		return h.varsity.code
	}
	if h.varsityCodeCached != "" {
		return h.varsityCodeCached
	}
	return "unknown"
}

func (h *Heading) VarsityPrettyName() string {
	if h.varsity != nil {
		return h.varsity.prettyName
	}
	if h.varsityPrettyNameCached != "" {
		return h.varsityPrettyNameCached
	}
	return "unknown"
}

// outscores determines if application app1 is better than application app2 for this heading.
// A lower ratingPlace is considered better.
// ApplicationsCache of better competition type beat applications of any lower one.
func (h *Heading) outscores(app1 *Application, app2 *Application) bool {
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
	Heading *Heading
	// List of students admitted to the heading, sorted according to admission criteria (e.g., quota, BVI, Regular, then rating).
	Admitted []*Student
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
// THIS STRUCT IS REPLACED/MODIFIED SIGNIFICANTLY FOR GALE-SHAPLEY
/*
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
*/

// removeStudentFromOldPlacement removes a specific application from its recorded placement.
// THIS FUNCTION IS NO LONGER NEEDED WITH GALE-SHAPLEY
/*
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
*/

// --- Start Heap Implementations for Gale-Shapley ---

// QuotaApplicationHeap stores applications for a specific quota type for a heading.
// It's a min-heap where the root is the *least preferred* (highest ratingPlace) accepted student.
type QuotaApplicationHeap struct {
	applications []*Application
	// heading field is not strictly needed here as comparison is only on ratingPlace
}

func (h *QuotaApplicationHeap) Len() int { return len(h.applications) }
func (h *QuotaApplicationHeap) Less(i, j int) bool {
	// Higher ratingPlace is worse. We want the root to be the student with the highest ratingPlace.
	return h.applications[i].ratingPlace > h.applications[j].ratingPlace
}
func (h *QuotaApplicationHeap) Swap(i, j int) {
	h.applications[i], h.applications[j] = h.applications[j], h.applications[i]
}
func (h *QuotaApplicationHeap) Push(x interface{}) {
	h.applications = append(h.applications, x.(*Application))
}
func (h *QuotaApplicationHeap) Pop() interface{} {
	old := h.applications
	n := len(old)
	app := old[n-1]
	h.applications = old[0 : n-1]
	return app
}

// GeneralApplicationHeap stores applications for general competition (Regular/BVI) for a heading.
// It's a min-heap where the root is the *least preferred* accepted student according to outscores.
type GeneralApplicationHeap struct {
	applications []*Application
	heading      *Heading // Essential for the outscores method
}

func (h *GeneralApplicationHeap) Len() int { return len(h.applications) }
func (h *GeneralApplicationHeap) Less(i, j int) bool {
	// We want the root to be the "worst" student.
	// If app[j] outscores app[i], then app[i] is worse than app[j].
	return h.heading.outscores(h.applications[j], h.applications[i])
}
func (h *GeneralApplicationHeap) Swap(i, j int) {
	h.applications[i], h.applications[j] = h.applications[j], h.applications[i]
}
func (h *GeneralApplicationHeap) Push(x interface{}) {
	h.applications = append(h.applications, x.(*Application))
}
func (h *GeneralApplicationHeap) Pop() interface{} {
	old := h.applications
	n := len(old)
	app := old[n-1]
	h.applications = old[0 : n-1]
	return app
}

// --- End Heap Implementations ---

// HeadingAdmissionStateGS tracks the admission status for a single heading using heaps for Gale-Shapley.
type HeadingAdmissionStateGS struct {
	heading         *Heading
	quotaAdmitted   map[Competition]heap.Interface // Stores *QuotaApplicationHeap
	generalAdmitted heap.Interface                 // Stores *GeneralApplicationHeap
}

// NewHeadingAdmissionStateGS creates a new state for a heading for Gale-Shapley.
func NewHeadingAdmissionStateGS(h *Heading) *HeadingAdmissionStateGS {
	state := &HeadingAdmissionStateGS{
		heading: h,
		quotaAdmitted: map[Competition]heap.Interface{
			CompetitionTargetQuota:    &QuotaApplicationHeap{applications: make([]*Application, 0, h.CapacitiesValue.TargetQuota)},
			CompetitionDedicatedQuota: &QuotaApplicationHeap{applications: make([]*Application, 0, h.CapacitiesValue.DedicatedQuota)},
			CompetitionSpecialQuota:   &QuotaApplicationHeap{applications: make([]*Application, 0, h.CapacitiesValue.SpecialQuota)},
		},
		generalAdmitted: &GeneralApplicationHeap{applications: make([]*Application, 0, h.CapacitiesValue.Regular), heading: h},
	}
	// Initialize the heaps
	if h.CapacitiesValue.TargetQuota > 0 {
		heap.Init(state.quotaAdmitted[CompetitionTargetQuota])
	}
	if h.CapacitiesValue.DedicatedQuota > 0 {
		heap.Init(state.quotaAdmitted[CompetitionDedicatedQuota])
	}
	if h.CapacitiesValue.SpecialQuota > 0 {
		heap.Init(state.quotaAdmitted[CompetitionSpecialQuota])
	}
	if h.CapacitiesValue.Regular > 0 {
		heap.Init(state.generalAdmitted)
	}
	return state
}

type VarsityCalculator struct {
	mu         sync.Mutex
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
		newStudent := &Student{IDValue: id, applications: make([]*Application, 0)} // Use exported ID
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

	if _, ok := v.headings.Load(code); ok {
		log.Printf("DUPLICATE HEADING CODE: Overwriting heading %s with prettyName %s and capacities %v", code, prettyName, capacities)
	}
	heading := &Heading{
		CodeValue:               code,
		varsity:                 v, // Link back to varsity
		CapacitiesValue:         capacities,
		PrettyNameValue:         prettyName,
		varsityCodeCached:       v.code,
		varsityPrettyNameCached: v.prettyName,
	}
	v.headings.Store(code, heading)
}

// AddApplication adds a student's application to a specific heading.
func (v *VarsityCalculator) AddApplication(headingCode, studentID string, ratingPlace, priority int, competitionType Competition, scoresSum int) {
	headingCode = strings.TrimSpace(headingCode)
	id, err := utils.PrepareStudentID(studentID)

	if err != nil {
		slog.Warn("Invalid too-long numeric student ID found", "studentID", studentID)
		return
	}

	h, ok := v.headings.Load(headingCode)
	if !ok {
		panic(fmt.Sprintf("heading with code %s not found", headingCode))
	}
	s := v.student(id)
	s.addApplication(h.(*Heading), ratingPlace, priority, competitionType, scoresSum)
}

// isValidPrioritySequence checks if the applications have a valid priority sequence (1, 2, 3, ..., N)
func isValidPrioritySequence(applications []*Application) bool {
	if len(applications) == 0 {
		return true
	}

	// Collect all priorities
	priorities := make([]int, len(applications))
	for i, app := range applications {
		priorities[i] = app.Priority()
	}

	// Sort priorities to check for valid sequence
	sort.Ints(priorities)

	// Check if sequence starts at 1 and has no gaps or duplicates
	for i, priority := range priorities {
		expected := i + 1
		if priority != expected {
			return false
		}
	}

	return true
}

// normalizePriorities fixes priority sequence while preserving relative order
func normalizePriorities(student *Student) {
	applications := student.Applications()
	if len(applications) == 0 {
		return
	}

	// Create a slice of application-priority pairs for sorting
	type appPriorityPair struct {
		app      *Application
		priority int
	}

	pairs := make([]appPriorityPair, len(applications))
	for i, app := range applications {
		pairs[i] = appPriorityPair{app: app, priority: app.Priority()}
	}

	// Sort by original priority to preserve relative order
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].priority < pairs[j].priority
	})

	// Assign new consecutive priorities starting from 1
	student.mu.Lock()
	for i, pair := range pairs {
		pair.app.priority = i + 1
	}
	student.mu.Unlock()
}

// NormalizeApplications iterates through all headings and normalizes the applications for each one.
// This should be called after all applications have been loaded and before any calculation is performed.
func (v *VarsityCalculator) NormalizeApplications() {
	headingToApplications := make(map[string][]*Application)

	v.students.Range(func(key, value interface{}) bool {
		student := value.(*Student)

		// Check if priorities need normalization (duplicates, gaps, wrong start, etc.)
		if !isValidPrioritySequence(student.Applications()) {
			slog.Debug("Student has invalid priority sequence, normalizing", "studentID", student.ID(), "originalPriorities", func() []int {
				priorities := make([]int, len(student.Applications()))
				for i, app := range student.Applications() {
					priorities[i] = app.Priority()
				}
				return priorities
			}())

			// Normalize priorities while preserving relative order
			normalizePriorities(student)
		}

		for _, app := range student.Applications() {
			if _, exists := headingToApplications[app.Heading().Code()]; !exists {
				headingToApplications[app.Heading().Code()] = make([]*Application, 0)
			}

			headingToApplications[app.Heading().Code()] = append(headingToApplications[app.Heading().Code()], app)
		}

		return true
	})

	v.headings.Range(func(key, value interface{}) bool {
		heading := value.(*Heading)

		normalizer := newApplicationsNormalizer(headingToApplications[heading.Code()])
		normalizer.normalize()

		return true // continue iteration
	})
}

// SetQuit marks a student as having quit.
func (v *VarsityCalculator) SetQuit(studentID string) {
	id, err := utils.PrepareStudentID(studentID)

	if err != nil {
		slog.Warn("Invalid too-long numeric student ID found", "studentID", studentID)
		return
	}

	s := v.student(id)
	s.mu.Lock()
	s.quit = true
	s.mu.Unlock()

	v.quitStudentsMu.Lock()
	v.quitStudents[id] = true
	v.quitStudentsMu.Unlock()
}

// SetOriginalSubmitted marks a student as submitting their original.
func (v *VarsityCalculator) SetOriginalSubmitted(studentID string) {
	id, err := utils.PrepareStudentID(studentID)

	if err != nil {
		slog.Warn("Invalid too-long numeric student ID found", "studentID", studentID)
		return
	}

	s := v.student(id)
	s.mu.Lock()
	s.originalSubmitted = true
	s.mu.Unlock()

	v.originalSubmittedStudentsMu.Lock()
	v.originalSubmittedStudents[id] = true
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
		return students[i].IDValue < students[j].IDValue // Use exported ID
	})
	return students
}

func (v *VarsityCalculator) GetStudent(studentID string) *Student {
	studentID, err := utils.PrepareStudentID(strings.TrimSpace(studentID))

	if err != nil {
		slog.Warn("Invalid too-long numeric student ID found", "studentID", studentID)
		return nil
	}

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
		return headings[i].PrettyNameValue < headings[j].PrettyNameValue
	})
	return headings
}

// CalculateAdmissions performs the main admission calculation logic for the varsity.
func (v *VarsityCalculator) CalculateAdmissions() []CalculationResult {
	slog.Debug("CalculateAdmissions: starting")
	v.mu.Lock()
	slog.Debug("CalculateAdmissions: mutex locked")
	defer func() {
		v.mu.Unlock()
		slog.Debug("CalculateAdmissions: mutex unlocked")
	}()

	if v.wasted {
		panic(fmt.Errorf("attempt to use a wasted calculator for varsity %s", v.code))
	}
	defer func() { v.wasted = true }()

	// 1. Initialization
	admissionStates := make(map[*Heading]*HeadingAdmissionStateGS)
	allHeadings := v.Headings()
	for _, h := range allHeadings {
		admissionStates[h] = NewHeadingAdmissionStateGS(h)
	}

	freeStudentsQueue := list.New()         // Changed from slice to container/list
	allStudents := v.Students()             // Students() already sorts in order if needed
	studentMap := make(map[string]*Student) // For quick lookup of student objects by ID

	studentNextProposalIndex := make(map[string]int)    // student.ID() -> index
	provisionalMatches := make(map[string]*Application) // student.ID() -> provisionally accepted Application

	for _, s := range allStudents {
		studentMap[s.ID()] = s
		studentNextProposalIndex[s.ID()] = 0 // Initialize proposal index
		if !s.Quit() {
			freeStudentsQueue.PushBack(s) // Changed to PushBack
		}
	}

	slog.Debug("CalculateAdmissions: initialized", "freeStudents", freeStudentsQueue.Len())

	// 2. Iteration (Gale-Shapley main loop)
	for freeStudentsQueue.Len() > 0 { // Changed to Len()
		element := freeStudentsQueue.Front() // Dequeue student (O(1))
		student := element.Value.(*Student)
		freeStudentsQueue.Remove(element) // O(1)

		studentID := student.ID()

		slog.Debug("Processing proposals for student", "studentID", studentID, "currentProposalIndex", studentNextProposalIndex[studentID], "apps", len(student.Applications()))

		// Student makes proposals in order of their preference list
		for proposalIdx := studentNextProposalIndex[studentID]; proposalIdx < len(student.Applications()); proposalIdx++ {
			app := student.Applications()[proposalIdx]
			heading := app.Heading()
			headingState := admissionStates[heading]

			slog.Debug("Student", "studentID", studentID, "appIndex", proposalIdx+1, "priority", app.Priority(), "headingCode", heading.Code(), "competitionType", app.CompetitionType(), "ratingPlace", app.RatingPlace())

			var targetHeap heap.Interface
			var capacity int
			isQuotaHeap := false

			switch app.CompetitionType() {
			case CompetitionTargetQuota, CompetitionDedicatedQuota, CompetitionSpecialQuota:
				targetHeap = headingState.quotaAdmitted[app.CompetitionType()]
				capacity = heading.CapacitiesValue.quotaCapacity(app.CompetitionType())
				isQuotaHeap = true
			case CompetitionRegular, CompetitionBVI:
				targetHeap = headingState.generalAdmitted

				currentFilledTargetQuota := 0
				if tqHeap, ok := headingState.quotaAdmitted[CompetitionTargetQuota]; ok {
					currentFilledTargetQuota = tqHeap.Len()
				}
				currentFilledDedicatedQuota := 0
				if dqHeap, ok := headingState.quotaAdmitted[CompetitionDedicatedQuota]; ok {
					currentFilledDedicatedQuota = dqHeap.Len()
				}
				currentFilledSpecialQuota := 0
				if sqHeap, ok := headingState.quotaAdmitted[CompetitionSpecialQuota]; ok {
					currentFilledSpecialQuota = sqHeap.Len()
				}

				currentFilledQuotasOverall := currentFilledTargetQuota + currentFilledDedicatedQuota + currentFilledSpecialQuota

				// The capacity for the generalAdmitted heap is the remaining total capacity of the heading
				// after accounting for students already provisionally admitted to specific quotas.
				capacity = heading.TotalCapacity() - currentFilledQuotasOverall
				if capacity < 0 { // Ensure capacity is not negative
					capacity = 0
				}

			default:
				slog.Debug("Student", "studentID", studentID, "unknownCompetitionType", app.CompetitionType(), "headingCode", heading.Code(), "action", "skippingThisApp")
				studentNextProposalIndex[studentID] = proposalIdx + 1 // Mark this proposal as considered
				continue                                              // Try next application
			}

			if capacity == 0 {
				slog.Debug("Student", "studentID", studentID, "headingCode", heading.Code(), "competitionType", app.CompetitionType(), "reason", "has0CapacityRejected")
				studentNextProposalIndex[studentID] = proposalIdx + 1 // Mark this proposal as considered
				continue                                              // Try next application
			}

			// Proposal attempt
			acceptedThisProposal := false
			var displacedApplication *Application
			wasDisplaced := false

			if targetHeap.Len() < capacity {
				heap.Push(targetHeap, app)
				acceptedThisProposal = true
				slog.Debug("Student", "studentID", studentID, "acceptedToHeading", heading.Code(), "competitionType", app.CompetitionType(), "reason", "capacityAvailable")
			} else {
				// Heap is full, compare with the worst student currently in the heap (root)
				worstAppInHeap := targetHeap.(interface{ Peek() *Application }).Peek()

				isNewAppBetter := false
				if isQuotaHeap {
					isNewAppBetter = app.RatingPlace() < worstAppInHeap.RatingPlace()
				} else { // General heap
					isNewAppBetter = heading.outscores(app, worstAppInHeap)
				}

				if isNewAppBetter {
					displacedApplication = heap.Pop(targetHeap).(*Application)
					wasDisplaced = true
					heap.Push(targetHeap, app)
					acceptedThisProposal = true
					slog.Debug("Student", "studentID", studentID, "acceptedToHeading", heading.Code(), "competitionType", app.CompetitionType(), "outscoresWorst", worstAppInHeap.StudentID())
				} else {
					slog.Debug("Student", "studentID", studentID, "rejectedByHeading", heading.Code(), "competitionType", app.CompetitionType(), "doesNotOutscoreWorst", worstAppInHeap.StudentID())
					// Student is rejected for this specific proposal, will try their next one.
				}
			}

			studentNextProposalIndex[studentID] = proposalIdx + 1 // Mark this proposal as considered

			if acceptedThisProposal {
				// If student was previously matched, that old match is now broken.
				if oldMatch, ok := provisionalMatches[studentID]; ok {
					slog.Debug("Student", "studentID", studentID, "isNowMatchedTo", heading.Code(), "breakingOldMatchWith", oldMatch.Heading().Code())
					// The spot at oldMatch.Heading() effectively becomes more available.
					// No explicit removal from old heading's heap needed here for the student *being moved*.
				}
				provisionalMatches[studentID] = app // Update/set provisional match

				if wasDisplaced {
					displacedStudentID := displacedApplication.StudentID()
					slog.Debug("Student", "displacedStudentID", displacedStudentID, "displacedFromHeading", displacedApplication.Heading().Code(), "by", studentID)

					delete(provisionalMatches, displacedStudentID) // Displaced student loses their provisional match

					// Add displaced student back to the free queue to find a new match
					displacedStudentObj := studentMap[displacedStudentID]
					if displacedStudentObj != nil {
						// Check if already in queue to prevent duplicates if logic allows (though GS typically processes one student fully)
						// For simplicity here, we add. If a student is processed multiple times due to re-queuing,
						// their studentNextProposalIndex ensures they don't re-propose to same.
						freeStudentsQueue.PushBack(displacedStudentObj) // Changed to PushBack
						slog.Debug("DisplacedStudent", "displacedStudentID", displacedStudentID, "action", "addedBackToFreeQueue")
					} else {
						slog.Debug("ERROR", "couldNotFindStudentObjectForDisplacedID", displacedStudentID)
					}
				}
				// Student is now provisionally matched, so they stop proposing in this turn.
				goto nextStudentInQueue // Break out of the inner proposal loop for the current student
			}
			// If rejected, loop continues to the student's next application.
		} // End loop over student's applications

		// If loop finishes, student has exhausted all their preferences or found a match.
		// If they found a match, `goto nextStudentInQueue` was hit.
		// If they exhausted preferences without a match, they remain unmatched.
		// Check if student has a match for logging
		{
			hasMatch := false
			if match, ok := provisionalMatches[studentID]; ok && match.student != nil {
				hasMatch = true
			}
			slog.Debug("Student", "studentID", studentID, "finishedProposingForThisTurn", hasMatch)
		}

	nextStudentInQueue: // Label for goto
	} // End while freeStudentsQueue is not empty

	slog.Debug("CalculateAdmissions: all proposals processed")

	// 3. Collect Results
	finalAdmissionsByHeading := make(map[*Heading][]*Student)
	for _, app := range provisionalMatches { // Iterate over final matches
		h := app.Heading()
		finalAdmissionsByHeading[h] = append(finalAdmissionsByHeading[h], app.student)
	}

	var results []CalculationResult
	for _, h := range allHeadings { // Use allHeadings to ensure all headings are in results
		admittedStudents := finalAdmissionsByHeading[h]

		// Sort admitted students for this heading based on the heading's preference criteria for consistent output
		// This requires getting their applications for this heading.
		// The provisionalMatches map stores the winning application.

		// Create a temporary list of winning applications for this heading to sort
		winningAppsForHeading := make([]*Application, 0, len(admittedStudents))
		for _, student := range admittedStudents {
			if app, ok := provisionalMatches[student.ID()]; ok && app.Heading() == h {
				winningAppsForHeading = append(winningAppsForHeading, app)
			}
		}

		sort.Slice(winningAppsForHeading, func(i, j int) bool {
			return h.outscores(winningAppsForHeading[i], winningAppsForHeading[j])
		})

		// Reconstruct sorted student list
		sortedAdmittedStudents := make([]*Student, len(winningAppsForHeading))
		for i, app := range winningAppsForHeading {
			sortedAdmittedStudents[i] = app.student
		}

		results = append(results, CalculationResult{
			Heading:  h,
			Admitted: sortedAdmittedStudents,
		})
	}

	// Sort final results by heading code for deterministic output
	sort.Slice(results, func(i, j int) bool {
		return results[i].Heading.Code() < results[j].Heading.Code()
	})

	slog.Debug("CalculateAdmissions: finished")
	return results
}

// Peek is helper method for heaps (assuming heap is not empty)
func (h *QuotaApplicationHeap) Peek() *Application {
	return h.applications[0]
}
func (h *GeneralApplicationHeap) Peek() *Application {
	return h.applications[0]
}

// processStudentApplications, tryAdmitQuota, tryAdmitGeneral, canAdmit, admissionInfo
// are no longer needed with the Gale-Shapley implementation.
/*
// processStudentApplications processes a single student's applications.
func (v *VarsityCalculator) processStudentApplications(
// ... (omitted old implementation)
}

// Returns true if admitted, the displaced student (if any), and the displaced application (if any).
func (v *VarsityCalculator) tryAdmitQuota(
// ... (omitted old implementation)
}

// tryAdmitGeneral attempts to admit a student to general competition (Regular or BVI).
// Returns true if admitted, the displaced student (if any), and the displaced application (if any).
func (v *VarsityCalculator) tryAdmitGeneral(
// ... (omitted old implementation)
}

// canAdmit checks if a student can be admitted to a heading based on their current best placement.
// A student can be admitted if the new application's heading has higher priority than their current best,
// or if they have no current best placement.
func (v *VarsityCalculator) canAdmit(
// ... (omitted old implementation)
}

// This struct was assumed by the logging I added. It needs to be defined.
// It stores information about a student's admission to a particular heading.
type admissionInfo struct {
	app          Application            // The specific application that led to admission
	headingState *HeadingAdmissionState // The state of the heading they were admitted to
}
*/

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

	drainableCount := 0

	var drainableStudents []*Student
	v.students.Range(func(key, value interface{}) bool {
		student := value.(*Student)
		student.mu.Lock()

		if !student.quit && !student.originalSubmitted {
			drainableCount += 1
			drainableStudents = append(drainableStudents, student)
		}

		student.mu.Unlock()
		return true
	})

	if len(drainableStudents) == 0 {
		return
	}

	numToDrain := (len(drainableStudents) * drainPercent) / 100

	// Shuffle eligible students to pick randomly
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(drainableStudents), func(i, j int) {
		drainableStudents[i], drainableStudents[j] = drainableStudents[j], drainableStudents[i]
	})

	for i := 0; i < numToDrain && i < len(drainableStudents); i++ {
		studentToDrain := drainableStudents[i]
		studentToDrain.mu.Lock()
		studentToDrain.quit = true
		studentToDrain.mu.Unlock()
	}

	v.drainedPercent = drainPercent
}

func (s *Student) OriginalSubmitted() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.originalSubmitted
}

// --- Custom gob encoding/decoding to preserve varsity metadata ---

func (h *Heading) GobEncode() ([]byte, error) {
	type alias struct {
		CodeValue         string
		CapacitiesValue   Capacities
		PrettyNameValue   string
		VarsityCode       string
		VarsityPrettyName string
	}

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(alias{
		CodeValue:         h.CodeValue,
		CapacitiesValue:   h.CapacitiesValue,
		PrettyNameValue:   h.PrettyNameValue,
		VarsityCode:       h.VarsityCode(),
		VarsityPrettyName: h.VarsityPrettyName(),
	}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (h *Heading) GobDecode(data []byte) error {
	type alias struct {
		CodeValue         string
		CapacitiesValue   Capacities
		PrettyNameValue   string
		VarsityCode       string
		VarsityPrettyName string
	}

	var aux alias
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&aux); err != nil {
		return err
	}

	h.CodeValue = aux.CodeValue
	h.CapacitiesValue = aux.CapacitiesValue
	h.PrettyNameValue = aux.PrettyNameValue
	// varsity pointer is nil after decoding; store cached data for getters.
	h.varsity = nil
	h.varsityCodeCached = aux.VarsityCode
	h.varsityPrettyNameCached = aux.VarsityPrettyName
	return nil
}
