package rzgmu

import (
	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/source/rzgmu"
)

var Varsity = source.VarsityDefinition{
	Name:           "РязГМУ",
	Code:           "rzgmu",
	HeadingSources: sourcesList(),
}

func sourcesList() []source.HeadingSource {
	return []source.HeadingSource{
		&rzgmu.HTTPHeadingSource{
			ProgramName: "Лечебное дело",
			Capacities: core.Capacities{
				Regular:        19,
				TargetQuota:    237,
				DedicatedQuota: 33,
				SpecialQuota:   33,
			},
		},
		&rzgmu.HTTPHeadingSource{
			ProgramName: "Стоматология",
			Capacities: core.Capacities{
				Regular:        3,
				TargetQuota:    15,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		&rzgmu.HTTPHeadingSource{
			ProgramName: "Педиатрия",
			Capacities: core.Capacities{
				Regular:        5,
				TargetQuota:    62,
				DedicatedQuota: 9,
				SpecialQuota:   9,
			},
		},
		&rzgmu.HTTPHeadingSource{
			ProgramName: "Медико-профилактическое дело",
			Capacities: core.Capacities{
				Regular:        9,
				TargetQuota:    23,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},
		&rzgmu.HTTPHeadingSource{
			ProgramName: "Фармация",
			Capacities: core.Capacities{
				Regular:        10,
				TargetQuota:    9,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		&rzgmu.HTTPHeadingSource{
			ProgramName: "Клиническая психология",
			Capacities: core.Capacities{
				Regular:        19,
				TargetQuota:    6,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},
	}
}
