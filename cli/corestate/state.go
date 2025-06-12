package corestate

import (
	"analabit/core"
	"analabit/core/drainer"
	"analabit/core/source"
	"sync"
	"sync/atomic"
)

var (
	LoadedVarsities []*source.Varsity
	PrimaryResults  map[string][]core.CalculationResult        // Key: Varsity Code
	DrainedResults  map[string]map[int][]drainer.DrainedResult // Key: Varsity Code, Key: Drain Percent

	// Background simulation tracking
	TotalSimulations     int32
	CompletedSimulations atomic.Int32
	SimulationsDone      bool // Flag to indicate all simulations are complete

	// Mutexes
	ResultsMutex sync.RWMutex

	// Initialization control
	Initialized  bool       // Exported, was 'initialized'
	InitMutex    sync.Mutex // Was already exported
	InitError    error
	CrawlingDone bool
)

func InitializeState() {
	PrimaryResults = make(map[string][]core.CalculationResult)
	DrainedResults = make(map[string]map[int][]drainer.DrainedResult)
	SimulationsDone = false
	CrawlingDone = false
	Initialized = false // Use exported field
	// Reset atomic counter if it's being reused across initializations (e.g. in tests)
	CompletedSimulations.Store(0)
}

// IsInitialized() and SetInitialized() are removed as PersistentPreRunE will manage the Initialized flag directly under InitMutex lock.
