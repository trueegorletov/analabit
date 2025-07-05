package source

import (
	"analabit/core"
	"log"
	"log/slog"
	"sync"
)

type VarsityDefinition struct {
	Code           string
	Name           string
	HeadingSources []HeadingSource
}

type Varsity struct {
	*VarsityDefinition
	*core.VarsityCalculator
	*VarsityDataCache
}

func (v *Varsity) Prepare() {
	v.VarsityCalculator = core.NewVarsityCalculator(v.Code, v.Name)

	if v.VarsityDataCache == nil {
		v.VarsityDataCache = NewVarsityDataCache(v.VarsityDefinition)
	}
}

func (v *Varsity) Clone() *Varsity {
	if v.VarsityDataCache == nil {
		panic("trying to copy an unloaded varsity")
	}

	clone := Varsity{
		v.VarsityDefinition, nil, v.VarsityDataCache,
	}
	clone.loadFromCache()

	vc := clone.VarsityCalculator
	v.VarsityCalculator.ForEachQuit(func(studentID string) {
		vc.SetQuit(studentID)
	})
	v.VarsityCalculator.ForEachOriginalSubmitted(func(studentID string) {
		vc.SetOriginalSubmitted(studentID)
	})

	return &clone
}

type channelledReceiver struct {
	headingDataChan     chan *HeadingData
	applicationDataChan chan *ApplicationData
}

func (cr *channelledReceiver) PutHeadingData(hd *HeadingData) {
	cr.headingDataChan <- hd
}

func (cr *channelledReceiver) PutApplicationData(ad *ApplicationData) {
	cr.applicationDataChan <- ad
}

func (v *Varsity) AddHeading(hd *HeadingData) {
	v.VarsityCalculator.AddHeading(hd.Code, hd.Capacities, hd.PrettyName)
}

func (v *Varsity) AddApplication(ad *ApplicationData) {
	if ad.CompetitionType > core.CompetitionBVI && !ad.OriginalSubmitted {
		return // Most of the quota guys never submit their originals, so consider only those who did
	}

	v.VarsityCalculator.AddApplication(ad.HeadingCode, ad.StudentID, ad.RatingPlace, ad.Priority, ad.CompetitionType, ad.ScoresSum)
}

// LoadFromDefinitions loads data from all given HeadingSources asynchronously, starting one goroutine per source.
// It sets Calculator to the clean VarsityCalculator instance and adds received data to it.
func (v *Varsity) loadFromSources() map[string]bool {
	v.Prepare()

	submittedOriginals := make(map[string]bool)
	// Mutex for submittedOriginals map, as multiple application data might be processed concurrently
	// if the processor goroutine was parallelized. For current single processor, it's not strictly needed
	// but good for future-proofing if processor logic changes.
	var submittedOriginalsMu sync.Mutex

	headingDataChan := make(chan *HeadingData, 100)         // Buffered channel
	applicationDataChan := make(chan *ApplicationData, 100) // Buffered channel

	receiver := &channelledReceiver{
		headingDataChan:     headingDataChan,
		applicationDataChan: applicationDataChan,
	}

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
					v.SaveHeadingData(hd)
					v.AddHeading(hd)
				}
			case ad, ok := <-applicationDataChan:
				if !ok {
					applicationDataChan = nil // Mark as drained
				} else {
					v.SaveApplicationData(ad)
					v.AddApplication(ad)

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
			err := s.LoadTo(receiver)
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

	// After all applications are loaded, normalize them.
	v.VarsityCalculator.NormalizeApplications()

	return submittedOriginals
}

func (v *Varsity) loadFromCache() map[string]bool {
	v.Prepare()

	submittedOriginals := make(map[string]bool)

	for _, hd := range v.HeadingsCache {
		v.AddHeading(hd)
	}

	for _, ad := range v.ApplicationsCache {
		v.AddApplication(ad)

		if ad.OriginalSubmitted {
			submittedOriginals[ad.StudentID] = true
		}
	}

	// After all applications are loaded from cache, normalize them.
	v.VarsityCalculator.NormalizeApplications()

	return submittedOriginals
}

// loadAll loads all given Varsities in-place concurrently considering information about submitted originals.
// It loads all varsities in parallel, starting one goroutine per Varsity, then collects them to a slice
// and saves students who submitted there their original to a map (student ID -> Varsity.Code).
// If there are multiple Varsities for a student, uses the first one and logs a warning about that.
// Finally it sets student.quit to true for all students who submitted their original to a *different* varsity.
func loadAll(varsities []*Varsity, loadFunc func(*Varsity) map[string]bool) ([]*Varsity, map[string]string) {
	studentOriginals := make(map[string]string) // original StudentID -> Varsity.Code
	var studentOriginalsMu sync.Mutex
	var varsityLoadWg sync.WaitGroup

	for i := range varsities {
		varsityLoadWg.Add(1)
		go func(idx int) {
			defer varsityLoadWg.Done()
			currentVarsity := varsities[idx] // Operate on the pointer

			// loadFromSources data for this varsity and get students who submitted originals here
			submittedInThisVarsity := loadFunc(currentVarsity) // This returns map[originalStudentID]bool

			studentOriginalsMu.Lock()
			for studentID := range submittedInThisVarsity { // studentID is original ID
				if existingVarsityCode, found := studentOriginals[studentID]; found {
					if existingVarsityCode != currentVarsity.Code { // Log only if different varsity
						slog.Debug("Student submitted original to multiple varsities", "studentID", studentID, "varsityFirst", existingVarsityCode, "varsitySecond", currentVarsity.Code, "chosen", existingVarsityCode)
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
			if v.Code == varsityCode || v.VarsityCalculator == nil {
				if v.VarsityCalculator != nil {
					v.SetOriginalSubmitted(studentID)
				}

				continue
			}

			v.SetQuit(studentID)
		}
	}

	return varsities, studentOriginals
}

func LoadFromDefinitions(defs []VarsityDefinition) []*Varsity {
	var varsities []*Varsity
	for _, def := range defs {
		v := &Varsity{
			VarsityDefinition: &def,
			VarsityCalculator: nil,
		}
		varsities = append(varsities, v)
	}

	vs, _ := loadAll(varsities, func(v *Varsity) map[string]bool {
		return v.loadFromSources()
	})

	return vs
}

func LoadWithCaches(defs []VarsityDefinition, caches []*VarsityDataCache) []*Varsity {
	codeToCache := make(map[string]*VarsityDataCache)
	for _, cache := range caches {
		codeToCache[cache.Definition.Code] = cache
	}

	var newVarsities []*Varsity
	var cacheVarsities []*Varsity

	for _, def := range defs {
		v := &Varsity{
			VarsityDefinition: &def,
			VarsityCalculator: nil,
		}

		if cache, found := codeToCache[def.Code]; found {
			v.VarsityDataCache = cache
			cacheVarsities = append(cacheVarsities, v)
		} else {
			newVarsities = append(newVarsities, v)
		}
	}

	cached, cachedOrigs := loadAll(cacheVarsities, func(v *Varsity) map[string]bool {
		return v.loadFromCache()
	})

	added, addedOrigs := loadAll(newVarsities, func(v *Varsity) map[string]bool {
		return v.loadFromSources()
	})

	removedFromCached := make(map[string]bool)

	for studentID, varsityCode := range addedOrigs {
		for _, v := range cached {
			if v.Code == varsityCode || v.VarsityCalculator == nil {
				if v.VarsityCalculator != nil {
					v.SetOriginalSubmitted(studentID)
				}

				continue
			}
			v.SetQuit(studentID)
			removedFromCached[studentID] = true
		}
	}

	for studentID, varsityCode := range cachedOrigs {
		if _, found := removedFromCached[studentID]; found {
			slog.Debug("Student submitted original to multiple varsities", "studentID", studentID)
			continue
		}

		for _, v := range added {
			if v.Code == varsityCode || v.VarsityCalculator == nil {
				if v.VarsityCalculator != nil {
					v.SetOriginalSubmitted(studentID)
				}

				continue
			}
			v.SetQuit(studentID)
		}
	}

	return append(cached, added...)
}
