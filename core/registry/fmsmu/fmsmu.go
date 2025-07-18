package fmsmu

import (
	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/source/fmsmu"
)

const (
	varsityCode = "fmsmu"
	varsityName = "ПМГМУ"
)

var Varsity = source.VarsityDefinition{
	Code:           varsityCode,
	Name:           varsityName,
	HeadingSources: sourcesList(),
}

func sourcesList() []source.HeadingSource {
	return []source.HeadingSource{
		&fmsmu.HTTPHeadingSource{
			PrettyName:           "Интеллектуальные системы в гуманитарной сфере",
			RegularListID:        "4063",
			SpecialQuotaListID:   "5097",
			DedicatedQuotaListID: "5098",
			Capacities: core.Capacities{
				Regular:        13,
				SpecialQuota:   2,
				DedicatedQuota: 2,
			},
		},

		&fmsmu.HTTPHeadingSource{
			PrettyName:           "Лечебное дело",
			RegularListID:        "4261",
			SpecialQuotaListID:   "5291",
			DedicatedQuotaListID: "5292",
			TargetQuotaListIDs: []string{
				"5270",
				"5313",
				"5269",
				"5278",
				"5301",
				"5305",
				"5295",
				"5304",
				"5312",
				"5310",
				"5303",
				"5309",
				"5302",
				"5290",
				"5293",
				"5273",
				"5296",
				"5274",
				"5275",
				"5300",
				"5280",
				"5311",
				"5307",
				"5271",
				"5276",
				"5287",
				"5306",
				"5284",
				"5286",
				"5288",
				"5299",
				"5289",
				"5308",
				"5283",
				"5272",
				"5279",
				"5281",
				"5285",
				"5298",
			},
			Capacities: core.Capacities{
				Regular:        274,
				TargetQuota:    521,
				SpecialQuota:   110,
				DedicatedQuota: 110,
			},
		},

		&fmsmu.HTTPHeadingSource{
			PrettyName:           "Менеджмент",
			RegularListID:        "4274",
			SpecialQuotaListID:   "5135",
			DedicatedQuotaListID: "5134",
			Capacities: core.Capacities{
				Regular:        14,
				SpecialQuota:   2,
				DedicatedQuota: 2,
			},
		},

		&fmsmu.HTTPHeadingSource{
			PrettyName:           "Сестринское дело",
			RegularListID:        "4281",
			SpecialQuotaListID:   "5173",
			DedicatedQuotaListID: "5172",
			TargetQuotaListIDs: []string{
				"5174",
			},
			Capacities: core.Capacities{
				Regular:        7,
				TargetQuota:    5,
				SpecialQuota:   2,
				DedicatedQuota: 2,
			},
		},

		&fmsmu.HTTPHeadingSource{
			PrettyName:           "Социальная работа",
			RegularListID:        "4283",
			SpecialQuotaListID:   "5177",
			DedicatedQuotaListID: "5176",
			Capacities: core.Capacities{
				Regular:        13,
				SpecialQuota:   2,
				DedicatedQuota: 2,
			},
		},

		&fmsmu.HTTPHeadingSource{
			PrettyName:           "Электронные и оптико-электронные приборы и системы специального назначения",
			RegularListID:        "4290",
			SpecialQuotaListID:   "5205",
			DedicatedQuotaListID: "5204",
			Capacities: core.Capacities{
				SpecialQuota:   1,
				DedicatedQuota: 1,
			},
		},

		&fmsmu.HTTPHeadingSource{
			PrettyName:           "Биоинженерия и биоинформатика",
			RegularListID:        "4059",
			SpecialQuotaListID:   "5086",
			DedicatedQuotaListID: "5085",
			TargetQuotaListIDs: []string{
				"5084",
			},
			Capacities: core.Capacities{
				Regular:        16,
				TargetQuota:    7,
				SpecialQuota:   3,
				DedicatedQuota: 3,
			},
		},

		&fmsmu.HTTPHeadingSource{
			PrettyName:           "Информационные системы и технологии",
			RegularListID:        "4065",
			SpecialQuotaListID:   "5102",
			DedicatedQuotaListID: "5101",
			Capacities: core.Capacities{
				Regular:        15,
				SpecialQuota:   3,
				DedicatedQuota: 3,
			},
		},

		&fmsmu.HTTPHeadingSource{
			PrettyName:           "Наноматериалы",
			RegularListID:        "3802",
			SpecialQuotaListID:   "5142",
			DedicatedQuotaListID: "5141",
			Capacities: core.Capacities{
				Regular:        3,
				SpecialQuota:   1,
				DedicatedQuota: 1,
			},
		},

		&fmsmu.HTTPHeadingSource{
			PrettyName:           "Педиатрия",
			RegularListID:        "4279",
			SpecialQuotaListID:   "5150",
			DedicatedQuotaListID: "5165",
			TargetQuotaListIDs: []string{
				"5143",
				"5170",
				"5151",
				"5149",
				"5147",
				"5152",
				"5167",
				"5168",
				"5162",
				"5164",
				"5161",
				"5169",
				"5148",
				"5171",
				"5163",
				"5144",
			},
			Capacities: core.Capacities{
				Regular:        13,
				TargetQuota:    131,
				SpecialQuota:   19,
				DedicatedQuota: 19,
			},
		},

		&fmsmu.HTTPHeadingSource{
			PrettyName:           "Биотехнология",
			RegularListID:        "4060",
			SpecialQuotaListID:   "5094",
			DedicatedQuotaListID: "5093",
			TargetQuotaListIDs: []string{
				"5095",
				"5096",
			},
			Capacities: core.Capacities{
				Regular:        28,
				TargetQuota:    7,
				SpecialQuota:   5,
				DedicatedQuota: 5,
			},
		},

		&fmsmu.HTTPHeadingSource{
			PrettyName:           "Клиническая психология",
			RegularListID:        "4066",
			SpecialQuotaListID:   "5110",
			DedicatedQuotaListID: "5109",
			TargetQuotaListIDs: []string{
				"5103",
				"5105",
				"5108",
				"5104",
				"5106",
			},
			Capacities: core.Capacities{
				Regular:        32,
				TargetQuota:    7,
				SpecialQuota:   5,
				DedicatedQuota: 5,
			},
		},

		&fmsmu.HTTPHeadingSource{
			PrettyName:           "Медицинская биохимия",
			RegularListID:        "4272",
			SpecialQuotaListID:   "5130",
			DedicatedQuotaListID: "5133",
			TargetQuotaListIDs: []string{
				"5131",
				"5132",
			},
			Capacities: core.Capacities{
				Regular:        20,
				TargetQuota:    4,
				SpecialQuota:   3,
				DedicatedQuota: 3,
			},
		},

		&fmsmu.HTTPHeadingSource{
			PrettyName:           "Механика и математическое моделирование",
			RegularListID:        "4276",
			SpecialQuotaListID:   "5138",
			DedicatedQuotaListID: "5139",
			Capacities: core.Capacities{
				Regular:        14,
				SpecialQuota:   2,
				DedicatedQuota: 2,
			},
		},

		&fmsmu.HTTPHeadingSource{
			PrettyName:           "Лингвистика",
			RegularListID:        "4264",
			SpecialQuotaListID:   "5113",
			DedicatedQuotaListID: "5112",
			Capacities: core.Capacities{
				Regular:        13,
				SpecialQuota:   3,
				DedicatedQuota: 3,
			},
		},

		&fmsmu.HTTPHeadingSource{
			PrettyName:           "Материаловедение и технологии материалов",
			RegularListID:        "4267",
			SpecialQuotaListID:   "5007",
			DedicatedQuotaListID: "5006",
			Capacities: core.Capacities{
				Regular:        7,
				SpecialQuota:   1,
				DedicatedQuota: 1,
			},
		},

		&fmsmu.HTTPHeadingSource{
			PrettyName:           "Медико-профилактическое дело",
			RegularListID:        "4268",
			SpecialQuotaListID:   "5117",
			DedicatedQuotaListID: "5116",
			TargetQuotaListIDs: []string{
				"5121",
				"5115",
				"5123",
				"5120",
				"5122",
				"5126",
				"5118",
				"5124",
			},
			Capacities: core.Capacities{
				Regular:        58,
				TargetQuota:    46,
				SpecialQuota:   14,
				DedicatedQuota: 14,
			},
		},

		&fmsmu.HTTPHeadingSource{
			PrettyName:           "Медицинская биофизика",
			RegularListID:        "4271",
			SpecialQuotaListID:   "5128",
			DedicatedQuotaListID: "5129",
			Capacities: core.Capacities{
				Regular:        19,
				SpecialQuota:   3,
				DedicatedQuota: 3,
			},
		},

		&fmsmu.HTTPHeadingSource{
			PrettyName:           "Стоматология",
			RegularListID:        "4286",
			SpecialQuotaListID:   "5180",
			DedicatedQuotaListID: "5187",
			TargetQuotaListIDs: []string{
				"5194",
				"5183",
				"5191",
				"5193",
				"5190",
				"5182",
				"5179",
				"5186",
				"5181",
				"5192",
				"5188",
				"5184",
				"5178",
				"5189",
			},
			Capacities: core.Capacities{
				Regular:        34,
				TargetQuota:    78,
				SpecialQuota:   15,
				DedicatedQuota: 15,
			},
		},

		&fmsmu.HTTPHeadingSource{
			PrettyName:           "Фармация",
			RegularListID:        "4288",
			SpecialQuotaListID:   "5203",
			DedicatedQuotaListID: "5202",
			TargetQuotaListIDs: []string{
				"5196",
				"5195",
				"5201",
				"5200",
			},
			Capacities: core.Capacities{
				Regular:        129,
				TargetQuota:    23,
				SpecialQuota:   20,
				DedicatedQuota: 20,
			},
		},
	}
}
