package main

import (
	"analabit/core"
	"analabit/core/source"
	"analabit/core/source/hse"
	"fmt"
	"log"
	"sync"
)

func main() {
	PrintHseFileHeadingSourceData()
}

// PrintHseFileHeadingSourceData prints the sample heading source data for HSE.
func PrintHseFileHeadingSourceData() {
	s := hse.FileHeadingSource{
		RCListPath:        "./sample_data/main/rc.xlsx",
		TQListPath:        "./sample_data/main/tq.xlsx",
		DQListPath:        "./sample_data/main/dq.xlsx",
		SQListPath:        "./sample_data/main/sq.xlsx",
		BListPath:         "./sample_data/main/bvi.xlsx",
		HeadingCapacities: core.Capacities{25, 2, 2, 2}, // Arbitrary capacity, as it's part of the struct
	}

	headingsChan := make(chan source.HeadingData, 3)          // Buffer size can be adjusted
	applicationsChan := make(chan source.ApplicationData, 15) // Buffer size can be adjusted

	var wg sync.WaitGroup
	var loadErr error

	// Goroutine to load data
	wg.Add(1)
	go func() {
		defer wg.Done()
		loadErr = s.LoadTo(nil)
	}()

	// Goroutines to process data from channels
	var headings []source.HeadingData
	var applications []source.ApplicationData

	wg.Add(1)
	go func() {
		defer wg.Done()
		for h := range headingsChan {
			headings = append(headings, h)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for a := range applicationsChan {
			applications = append(applications, a)
		}
	}()

	// Wait for all goroutines to complete
	wg.Wait()

	if loadErr != nil {
		log.Fatalf("Error loading data: %v", loadErr)
	}

	fmt.Printf("--- Headings --- (%d loaded)", len(headings))
	if len(headings) == 0 {
		fmt.Println("No heading data loaded.")
	}
	for _, h := range headings {
		fmt.Printf("Code: %s, Name: %s, Capacities: %d\n", h.Code, h.PrettyName, h.Capacities)
	}

	fmt.Printf("--- Applications --- (%d loaded)", len(applications))
	if len(applications) == 0 {
		fmt.Println("No application data loaded.")
	}
	for _, a := range applications {
		fmt.Printf("HeadingCode: %s, StudentID: %s, Scores: %d, Place: %d, Priority: %d, Competition: %s, Original: %t\n",
			a.HeadingCode, a.StudentID, a.ScoresSum, a.RatingPlace, a.Priority, a.CompetitionType, a.OriginalSubmitted)
	}
}
