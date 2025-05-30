package core

import (
	"analabit/utils"
	"fmt"
	"sort"
	"sync"
)

type Competition int

// Competition represents the type of competition for a student's application.
const (
	CompetitionRegular Competition = iota
	CompetitionTargetQuota
	CompetitionDedicatedQuota
	CompetitionSpecialQuota
	CompetitionBVI
)

// Application represents a student's application to a specific heading.
type Application struct {
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
	id string
	// List of applications made by the student, sorted by priority (ascending, e.g., priority 1 first).
	applications []Application
	// If true, the student has withdrawn their application from this varsity and is ignored in calculations.
	quit bool
}

// application retrieves the student's application details for a specific heading.
// It panics if no application is found for the given heading for this student.
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
	// Maximum number of students that can be admitted to this heading.
	capacity int
	// A human-readable name or identifier for the heading.
	prettyName string
}

// Capacity returns the maximum number of students that can be admitted to this heading.
func (h *Heading) Capacity() int {
	return h.capacity
}

// PrettyName returns the human-readable name of the heading.
func (h *Heading) PrettyName() string {
	return h.prettyName
}

// Code returns the unique identifier of the heading.
func (h *Heading) Code() string {
	return h.code
}

// outscores determines if student `s1` has a better (lower) rating place
// than student `s2` for this heading.
// A lower ratingPlace is considered better.
// Applications of better competition type beat applications of any lower one.
func (h *Heading) outscores(s1 *Student, s2 *Student) bool {
	app1 := s1.application(h)
	app2 := s2.application(h)

	// First, compare competition types
	if app1.competitionType < app2.competitionType {
		return false // s2 has better competition type
	}
	if app1.competitionType > app2.competitionType {
		return true // s1 has better competition type
	}

	// If competition types are the same, compare by ratingPlace
	return app1.ratingPlace < app2.ratingPlace
}

// VarsityCalculator handles the admission calculation process for a university.
type VarsityCalculator struct {
	// Unique identifier for the varsity.
	code string
	// students maps student IDs to Student objects.
	students sync.Map // Changed to sync.Map
	// headings maps heading codes to Heading objects.
	headings sync.Map // Changed to sync.Map
}

// NewVarsityCalculator creates and initializes a new VarsityCalculator.
//
// Parameters:
//   - code: The unique identifier for the varsity.
//
// Returns: A pointer to the newly created VarsityCalculator.
func NewVarsityCalculator(code string) *VarsityCalculator {
	// sync.Map does not need make, zero value is ready to use.
	return &VarsityCalculator{
		code: code,
	}
}

// CalculationResult holds the outcome of the admission process for a single heading.
type CalculationResult struct {
	// The heading for which this calculation result is generated.
	heading *Heading
	// List of students admitted to the heading, sorted by their rating place (best first).
	admitted []*Student
}

// CalculateAdmissions processes student applications and determines admissions for each heading
// within the varsity.
//
// The process iteratively places students based on application priority and rating:
//  1. Initializes empty admission lists for all headings.
//  2. For each student, considers their applications in order of priority.
//  3. If a heading has capacity, the student is provisionally admitted.
//  4. If a heading is full, the student is admitted only if they outscore the
//     worst-rated student currently admitted, who is then displaced.
//  5. If a student is admitted, any previous admission to a lower-priority heading is revoked.
//  6. Displaced students are re-considered for their other applications.
//  7. Iteration continues until a full pass occurs with no changes to student placements.
//
// Student.applications are assumed to be pre-sorted by priority (lowest number is highest priority).
// The final list of results is sorted by heading code.
func (v *VarsityCalculator) CalculateAdmissions() []CalculationResult {
	// finalAdmissions maps each heading to its list of admitted students, kept sorted by ratingPlace (best first).
	finalAdmissions := make(map[*Heading][]*Student)
	v.headings.Range(func(key, value interface{}) bool {
		h := value.(*Heading)
		finalAdmissions[h] = make([]*Student, 0, h.capacity)
		return true // continue iteration
	})

	// studentBestPlacement tracks the current best (highest priority) confirmed placement for each student.
	// (key: Student pointer, value: Application object).
	studentBestPlacement := make(map[*Student]Application)

	madeChangeInIteration := true
	for madeChangeInIteration {
		madeChangeInIteration = false // Reset for the current pass

		v.students.Range(func(key, value interface{}) bool {
			student := value.(*Student)

			student.mu.Lock()
			quitStatus := student.quit
			// Create a snapshot of applications while holding the lock
			currentApplicationsSnapshot := make([]Application, len(student.applications))
			copy(currentApplicationsSnapshot, student.applications)
			student.mu.Unlock()

			if quitStatus { // Skip students who have withdrawn.
				return true // continue Range
			}

			// Iterate through the student's applications snapshot (pre-sorted by priority).
			for _, currentApp := range currentApplicationsSnapshot {
				targetHeading := currentApp.heading

				// If student is already placed in a heading with strictly higher priority,
				// they won't consider this currentApp or any subsequent (lower priority) ones in this pass.
				if existingPlacement, hasPlacement := studentBestPlacement[student]; hasPlacement {
					if existingPlacement.priority < currentApp.priority {
						break // Already in a better priority placement.
					}
					// If currentApp is for the same heading and priority, they are already optimally placed for this app.
					if existingPlacement.priority == currentApp.priority && existingPlacement.heading == targetHeading {
						continue
					}
				}

				admittedToTargetHeading := finalAdmissions[targetHeading]
				canBePlaced := false

				if len(admittedToTargetHeading) < targetHeading.capacity {
					canBePlaced = true // Heading has space.
				} else {
					// Heading is full. Check if this student can displace the worst-admitted student.
					worstAdmittedStudent := admittedToTargetHeading[len(admittedToTargetHeading)-1]
					if targetHeading.outscores(student, worstAdmittedStudent) {
						// Displace the worst student.
						finalAdmissions[targetHeading] = admittedToTargetHeading[:len(admittedToTargetHeading)-1]
						delete(studentBestPlacement, worstAdmittedStudent)

						madeChangeInIteration = true // A student was displaced.
						canBePlaced = true
					}
				}

				if canBePlaced {
					// If student was already placed in a *different* heading, remove them from there.
					if existingPlacement, hasPlacement := studentBestPlacement[student]; hasPlacement {
						if existingPlacement.heading != targetHeading {
							oldHeadingAdmittedList := finalAdmissions[existingPlacement.heading]
							newOldHeadingAdmittedList := make([]*Student, 0, len(oldHeadingAdmittedList))
							for _, s := range oldHeadingAdmittedList {
								if s.id != student.id {
									newOldHeadingAdmittedList = append(newOldHeadingAdmittedList, s)
								}
							}
							finalAdmissions[existingPlacement.heading] = newOldHeadingAdmittedList
							madeChangeInIteration = true // Student moved headings.
						}
					} else {
						madeChangeInIteration = true // Student got their first placement.
					}

					// Add student to targetHeading's admission list and sort.
					finalAdmissions[targetHeading] = append(finalAdmissions[targetHeading], student)
					sort.Slice(finalAdmissions[targetHeading], func(i, j int) bool {
						return targetHeading.outscores(finalAdmissions[targetHeading][i], finalAdmissions[targetHeading][j])
					})

					studentBestPlacement[student] = currentApp // Update student's best known placement.

					// Student placed in this heading; won't consider lower priority applications in this pass.
					break
				}
			} // End loop over student's applications
			return true // continue Range
		}) // End loop over all students
	} // End main iteration loop

	// Prepare and sort final results.
	results := make([]CalculationResult, 0)
	v.headings.Range(func(key, value interface{}) bool {
		h := value.(*Heading)
		admittedStudentsForHeading := finalAdmissions[h]
		// Defensive sort, though list should be maintained sorted.
		sort.Slice(admittedStudentsForHeading, func(i, j int) bool {
			return h.outscores(admittedStudentsForHeading[i], admittedStudentsForHeading[j])
		})
		results = append(results, CalculationResult{
			heading:  h,
			admitted: admittedStudentsForHeading,
		})
		return true // continue iteration
	})

	return results
}

// student returns an existing Student for the given ID, or creates and stores a new one if not found.
func (v *VarsityCalculator) student(id string) *Student {
	// Attempt to load the student using the original id as the key.
	if val, ok := v.students.Load(id); ok {
		return val.(*Student)
	}

	// Student not found, prepare a new one.
	// The student's internal ID field will be the preparedID.
	preparedID, err := utils.PrepareStudentID(id)
	if err != nil {
		panic(fmt.Sprintf("failed to prepare student ID %q: %v", id, err))
	}
	newStudent := &Student{
		id:           preparedID,
		applications: make([]Application, 0),
		quit:         false, // Defaults to not quit.
	}

	// Atomically load or store. The key is the original `id`.
	// If another goroutine stored a student for this `id` in the meantime,
	// LoadOrStore will return that existing student. Otherwise, it stores newStudent.
	actual, _ := v.students.LoadOrStore(id, newStudent)
	return actual.(*Student)
}

// AddHeading creates and adds a new heading to the varsity.
func (v *VarsityCalculator) AddHeading(headingCode string, capacity int, prettyName string) {
	newHeading := &Heading{
		code:       headingCode,
		varsity:    v,
		capacity:   capacity,
		prettyName: prettyName,
	}
	v.headings.Store(headingCode, newHeading)
}

// AddApplication adds a student's application for a specific heading with a Competition.
func (v *VarsityCalculator) AddApplication(headingCode string, studentID string, ratingPlace int, priority int, competitionType Competition) {
	val, ok := v.headings.Load(headingCode)
	if !ok {
		panic("adding application to non-existent heading: " + headingCode)
	}
	heading := val.(*Heading)

	v.student(studentID).addApplication(heading, ratingPlace, priority, competitionType)
}

// SetQuit marks a student as having withdrawn their applications.
func (v *VarsityCalculator) SetQuit(studentID string) {
	if val, ok := v.students.Load(studentID); ok {
		s := val.(*Student)
		s.mu.Lock()
		s.quit = true
		s.mu.Unlock()
	}
}

func (v *VarsityCalculator) GetStudents() []*Student {
	var studentsList []*Student
	v.students.Range(func(key, value interface{}) bool {
		studentsList = append(studentsList, value.(*Student))
		return true // continue iteration
	})
	return studentsList
}
