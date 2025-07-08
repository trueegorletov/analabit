package main

import (
	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/source/oldhse"
	"os"
)

func main() {
	Ser()
}

// Ser prints the sample heading source data for HSE.
func Ser() {
	defs := []source.VarsityDefinition{
		{
			Name: "HSE",
			Code: "hse_msk",
			HeadingSources: []source.HeadingSource{
				&oldhse.FileHeadingSource{
					RCListPath:        "./sample_data/hse/rc.xlsx",
					TQListPath:        "./sample_data/hse/tq.xlsx",
					DQListPath:        "./sample_data/hse/dq.xlsx",
					SQListPath:        "./sample_data/hse/sq.xlsx",
					BListPath:         "./sample_data/hse/bvi.xlsx",
					HeadingCapacities: core.Capacities{Regular: 25, TargetQuota: 2, DedicatedQuota: 2, SpecialQuota: 2}, // Arbitrary capacity, as it's part of the struct
				},
			},
		},
	}

	varsities := source.LoadFromDefinitions(defs)

	var caches []*source.VarsityDataCache

	for _, v := range varsities {
		caches = append(caches, v.VarsityDataCache)
	}

	f, _ := os.Create("varsities.gob")

	err := source.SerializeList(caches, f)

	if err != nil {
		panic(err)
	}

	f.Close()

	f, _ = os.Open("varsities.gob")

	_, err = source.DeserializeList(f)

	if err != nil {
		panic(err)
	}

	f.Close()
}
