package source

import (
	"log"
	"log/slog"
	"math"
	"sync"
	"time"

	"github.com/trueegorletov/analabit/core"
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
	var headings []*HeadingData
	var applications []*ApplicationData
	for {
		select {
		case hd, ok := <-headingDataChan:
			if !ok {
				headingDataChan = nil
			} else {
				headings = append(headings, hd)
			}
		case ad, ok := <-applicationDataChan:
			if !ok {
				applicationDataChan = nil
			} else {
				applications = append(applications, ad)
			}
		}
		if headingDataChan == nil && applicationDataChan == nil {
			break
		}
	}
	// Process headings first
	for _, hd := range headings {
		v.SaveHeadingData(hd)
		v.AddHeading(hd)
	}
	// Then process applications
	for _, ad := range applications {
		v.SaveApplicationData(ad)
		v.AddApplication(ad)
		if ad.OriginalSubmitted {
			submittedOriginalsMu.Lock()
			submittedOriginals[ad.StudentID] = true
			submittedOriginalsMu.Unlock()
		}
	}
}()

	if v.Code == "spbsu" {
	for _, hs := range v.HeadingSources {
		err := Retry(func() error { return hs.LoadTo(receiver) }, 3, func(attempt int) time.Duration {
			return time.Duration(math.Pow(2, float64(attempt-1))) * 10 * time.Second
		})
		if err != nil {
			slog.Error("Failed to load source after retries", "error", err)
		}
	}
} else {
	for _, hs := range v.HeadingSources {
		sourceWg.Add(1)
		go func(s HeadingSource) {
			defer sourceWg.Done()
			err := Retry(func() error { return s.LoadTo(receiver) }, 3, func(attempt int) time.Duration {
				return time.Duration(math.Pow(2, float64(attempt-1))) * 10 * time.Second
			})
			if err != nil {
				slog.Error("Failed to load source after retries", "error", err)
			}
		}(hs)
	}
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

	v.VarsityCalculator.NormalizeApplications()

	return submittedOriginals
}

// loadAll loads all given Varsities in-place concurrently considering information about submitted originals.
// It loads all varsities in parallel, starting one goroutine per Varsity, then collects them to a slice
// and saves students who submitted there their original to a map (student ID -> Varsity.Code).
// If there are multiple Varsities for a student, uses the first one and logs a warning about that.
// Finally it sets student.quit to true for all students who submitted their original to a *different* varsity.
func loadAll(varsities []*Varsity, loadFunc func(*Varsity) map[string]bool) ([]*Varsity, map[string]string) {
	// Removed spbsu.StartWorkers(10)
	// Removed defer spbsu.StopWorkers()

	studentOriginals := make(map[string]string) // original StudentID -> Varsity.Code
	var studentOriginalsMu sync.Mutex
	var varsityLoadWg sync.WaitGroup

	const maxConcurrentLoads = 4
	sem := make(chan struct{}, maxConcurrentLoads)

	for i := range varsities {
		varsityLoadWg.Add(1)
		sem <- struct{}{} // Acquire a slot
		go func(idx int) {
			defer func() {
				<-sem // Release the slot
				varsityLoadWg.Done()
			}()
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

	// Removed spbsu.WaitJobs()

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
	log.Printf("🔍 LOAD DEBUG: LoadFromDefinitions called with %d definitions", len(defs))
	
	var varsities []*Varsity
	for i, def := range defs {
		log.Printf("🔍 LOAD DEBUG: Processing definition %d: Code=%s, Name=%s, Sources=%d", i, def.Code, def.Name, len(def.HeadingSources))
		
		// Special debug for MIREA
		if def.Code == "mirea" {
			log.Printf("🔍 LOAD DEBUG: [MIREA] Found MIREA definition with %d heading sources", len(def.HeadingSources))
			for j, src := range def.HeadingSources {
				log.Printf("🔍 LOAD DEBUG: [MIREA] Source %d: %T", j, src)
			}
		}
		
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

func LoadWithCaches(defs []VarsityDefinition, caches []*VarsityDataCache) ([]*Varsity, bool) {
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

	cacheComplete := len(newVarsities) == 0
	return append(cached, added...), cacheComplete
}
