package source

import (
	"analabit/core"
	"log"
	"sync"
)

type Varsity struct {
	Code           string
	Name           string
	HeadingSources []HeadingSource
	Calculator     *core.VarsityCalculator
}

func (v *Varsity) ResetCalculator() {
	v.Calculator = core.NewVarsityCalculator(v.Code)
}

// Load loads data from all providen HeadingSources asynchronously, starting one goroutine per source.
// It sets Calculator to the clean VarsityCalculator instance and adds received data to it.
func (v *Varsity) Load() map[string]bool {
	v.ResetCalculator()

	submittedOriginals := make(map[string]bool)
	// Mutex for submittedOriginals map, as multiple application data might be processed concurrently
	// if the processor goroutine was parallelized. For current single processor, it's not strictly needed
	// but good for future-proofing if processor logic changes.
	var submittedOriginalsMu sync.Mutex

	headingDataChan := make(chan HeadingData, 100)         // Buffered channel
	applicationDataChan := make(chan ApplicationData, 100) // Buffered channel

	var sourceWg sync.WaitGroup
	var processingWg sync.WaitGroup

	processingWg.Add(1)
	go func() {
		defer processingWg.Done()
		for {
			select {
			case hd, ok := <-headingDataChan:
				if !ok {
					headingDataChan = nil // Mark as drained
				} else {
					v.Calculator.AddHeading(hd.Code, hd.Capacities, hd.PrettyName)
				}
			case ad, ok := <-applicationDataChan:
				if !ok {
					applicationDataChan = nil // Mark as drained
				} else {
					// Assuming ad.StudentID is the ORIGINAL student ID
					if ad.CompetitionType > core.CompetitionBVI && !ad.OriginalSubmitted {
						continue // Most of the quota guys never submit their originals, so consider only those who did
					}

					v.Calculator.AddApplication(ad.HeadingCode, ad.StudentID, ad.RatingPlace, ad.Priority, ad.CompetitionType)
					if ad.OriginalSubmitted {
						submittedOriginalsMu.Lock()
						submittedOriginals[ad.StudentID] = true
						submittedOriginalsMu.Unlock()
					}
				}
			}
			if headingDataChan == nil && applicationDataChan == nil {
				return // Both channels drained and closed
			}
		}
	}()

	for _, hs := range v.HeadingSources {
		sourceWg.Add(1)
		go func(s HeadingSource) {
			defer sourceWg.Done()
			err := s.LoadTo(headingDataChan, applicationDataChan)
			if err != nil {
				// TODO: Define a more robust error handling strategy. For now, log it.
				log.Printf("Error loading data from source: %v", err)
			}
		}(hs)
	}

	sourceWg.Wait()            // Wait for all sources to finish sending data
	close(headingDataChan)     // Close data channels
	close(applicationDataChan) //
	processingWg.Wait()        // Wait for the processor goroutine to finish all writes

	return submittedOriginals
}

// LoadAll loads all given Varsities in-place concurrently considering information about submitted originals.
// It loads all varsities in parallel, starting one goroutine per Varsity, then collects them to a slice
// and saves students who submitted there their original to a map (student ID -> Varsity.Code).
// If there are multiple Varsities for a student, uses the first one and logs a warning about that.
// Finally it sets student.quit to true for all students who submitted their original to a *different* varsity.
func LoadAll(varsities []*Varsity) []*Varsity {
	studentOriginals := make(map[string]string) // original StudentID -> Varsity.Code
	var studentOriginalsMu sync.Mutex
	var varsityLoadWg sync.WaitGroup

	for i := range varsities {
		varsityLoadWg.Add(1)
		go func(idx int) {
			defer varsityLoadWg.Done()
			currentVarsity := varsities[idx] // Operate on the pointer

			// Load data for this varsity and get students who submitted originals here
			submittedInThisVarsity := currentVarsity.Load() // This returns map[originalStudentID]bool

			studentOriginalsMu.Lock()
			for studentID := range submittedInThisVarsity { // studentID is original ID
				if existingVarsityCode, found := studentOriginals[studentID]; found {
					if existingVarsityCode != currentVarsity.Code { // Log only if different varsity
						log.Printf("Warning: Student %s submitted original to multiple varsities: %s and %s. Using %s.",
							studentID, existingVarsityCode, currentVarsity.Code, existingVarsityCode)
					}
				} else {
					studentOriginals[studentID] = currentVarsity.Code
				}
			}
			studentOriginalsMu.Unlock()
		}(i)
	}
	varsityLoadWg.Wait()

	// Set student.quit status
	for studentID, varsityCode := range studentOriginals {
		for i := range varsities {
			v := varsities[i]
			if v.Code == varsityCode || v.Calculator == nil {
				continue
			}

			v.Calculator.SetQuit(studentID)
		}
	}

	return varsities
}
