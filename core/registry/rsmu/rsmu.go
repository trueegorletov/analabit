package rsmu

import (
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/source/rsmu"
)

var Varsity = source.VarsityDefinition{
	Name:           "РНИМУ",
	Code:           "rsmu",
	HeadingSources: sourcesList(),
}

func sourcesList() []source.HeadingSource {
	return []source.HeadingSource{
		// Generated RSMU HTTPHeadingSource entries
		// Based on registry: sample_data/rsmu/p2_560.json

		&rsmu.HTTPHeadingSource{
			ProgramName:           "Психология безопасности",
			RegularListURL:        "https://submitted.rsmu.ru/data/p3_22095.json",
			SpecialQuotaListURL:   "https://submitted.rsmu.ru/data/p3_22111.json",
			DedicatedQuotaListURL: "https://submitted.rsmu.ru/data/p3_22155.json",
		},

		&rsmu.HTTPHeadingSource{
			ProgramName:           "Социальная работа",
			RegularListURL:        "https://submitted.rsmu.ru/data/p3_22066.json",
			SpecialQuotaListURL:   "https://submitted.rsmu.ru/data/p3_22160.json",
			DedicatedQuotaListURL: "https://submitted.rsmu.ru/data/p3_22105.json",
		},

		&rsmu.HTTPHeadingSource{
			ProgramName: "Медицинская биохимия",
			TargetQuotaListURLs: []string{
				"https://submitted.rsmu.ru/data/p3_22178.json",
				"https://submitted.rsmu.ru/data/p3_22133.json",
				"https://submitted.rsmu.ru/data/p3_22065.json",
			},
			RegularListURL:        "https://submitted.rsmu.ru/data/p3_22118.json",
			SpecialQuotaListURL:   "https://submitted.rsmu.ru/data/p3_22099.json",
			DedicatedQuotaListURL: "https://submitted.rsmu.ru/data/p3_22153.json",
		},

		&rsmu.HTTPHeadingSource{
			ProgramName: "Медицинская информатика",
			TargetQuotaListURLs: []string{
				"https://submitted.rsmu.ru/data/p3_22170.json",
			},
			RegularListURL:        "https://submitted.rsmu.ru/data/p3_22084.json",
			SpecialQuotaListURL:   "https://submitted.rsmu.ru/data/p3_22102.json",
			DedicatedQuotaListURL: "https://submitted.rsmu.ru/data/p3_22137.json",
		},

		&rsmu.HTTPHeadingSource{
			ProgramName: "Стоматология",
			TargetQuotaListURLs: []string{
				"https://submitted.rsmu.ru/data/p3_22129.json",
				"https://submitted.rsmu.ru/data/p3_22104.json",
				"https://submitted.rsmu.ru/data/p3_22076.json",
				"https://submitted.rsmu.ru/data/p3_22067.json",
				"https://submitted.rsmu.ru/data/p3_22116.json",
				"https://submitted.rsmu.ru/data/p3_22139.json",
			},
			RegularListURL:        "https://submitted.rsmu.ru/data/p3_22150.json",
			SpecialQuotaListURL:   "https://submitted.rsmu.ru/data/p3_22135.json",
			DedicatedQuotaListURL: "https://submitted.rsmu.ru/data/p3_22101.json",
		},

		&rsmu.HTTPHeadingSource{
			ProgramName:           "Фармация",
			RegularListURL:        "https://submitted.rsmu.ru/data/p3_22087.json",
			SpecialQuotaListURL:   "https://submitted.rsmu.ru/data/p3_22119.json",
			DedicatedQuotaListURL: "https://submitted.rsmu.ru/data/p3_22082.json",
		},

		&rsmu.HTTPHeadingSource{
			ProgramName:           "Биомедицина",
			RegularListURL:        "https://submitted.rsmu.ru/data/p3_22181.json",
			SpecialQuotaListURL:   "https://submitted.rsmu.ru/data/p3_22078.json",
			DedicatedQuotaListURL: "https://submitted.rsmu.ru/data/p3_22144.json",
		},

		&rsmu.HTTPHeadingSource{
			ProgramName: "Клиническая психология в здравоохранении",
			TargetQuotaListURLs: []string{
				"https://submitted.rsmu.ru/data/p3_22183.json",
				"https://submitted.rsmu.ru/data/p3_22184.json",
				"https://submitted.rsmu.ru/data/p3_22127.json",
			},
			RegularListURL:        "https://submitted.rsmu.ru/data/p3_22096.json",
			SpecialQuotaListURL:   "https://submitted.rsmu.ru/data/p3_22162.json",
			DedicatedQuotaListURL: "https://submitted.rsmu.ru/data/p3_22134.json",
		},

		&rsmu.HTTPHeadingSource{
			ProgramName: "Лечебное дело",
			TargetQuotaListURLs: []string{
				"https://submitted.rsmu.ru/data/p3_22167.json",
				"https://submitted.rsmu.ru/data/p3_22140.json",
				"https://submitted.rsmu.ru/data/p3_22148.json",
				"https://submitted.rsmu.ru/data/p3_22166.json",
				"https://submitted.rsmu.ru/data/p3_22142.json",
				"https://submitted.rsmu.ru/data/p3_22069.json",
				"https://submitted.rsmu.ru/data/p3_22165.json",
				"https://submitted.rsmu.ru/data/p3_22070.json",
				"https://submitted.rsmu.ru/data/p3_22168.json",
				"https://submitted.rsmu.ru/data/p3_22097.json",
				"https://submitted.rsmu.ru/data/p3_22079.json",
				"https://submitted.rsmu.ru/data/p3_22179.json",
				"https://submitted.rsmu.ru/data/p3_22141.json",
				"https://submitted.rsmu.ru/data/p3_22157.json",
				"https://submitted.rsmu.ru/data/p3_22071.json",
				"https://submitted.rsmu.ru/data/p3_22164.json",
				"https://submitted.rsmu.ru/data/p3_22103.json",
				"https://submitted.rsmu.ru/data/p3_22120.json",
				"https://submitted.rsmu.ru/data/p3_22174.json",
				"https://submitted.rsmu.ru/data/p3_22088.json",
				"https://submitted.rsmu.ru/data/p3_22180.json",
				"https://submitted.rsmu.ru/data/p3_22126.json",
				"https://submitted.rsmu.ru/data/p3_22114.json",
				"https://submitted.rsmu.ru/data/p3_22098.json",
				"https://submitted.rsmu.ru/data/p3_22086.json",
				"https://submitted.rsmu.ru/data/p3_22163.json",
				"https://submitted.rsmu.ru/data/p3_22115.json",
				"https://submitted.rsmu.ru/data/p3_22093.json",
				"https://submitted.rsmu.ru/data/p3_22107.json",
				"https://submitted.rsmu.ru/data/p3_22176.json",
				"https://submitted.rsmu.ru/data/p3_22073.json",
				"https://submitted.rsmu.ru/data/p3_22130.json",
				"https://submitted.rsmu.ru/data/p3_22132.json",
				"https://submitted.rsmu.ru/data/p3_22090.json",
				"https://submitted.rsmu.ru/data/p3_22171.json",
			},
			RegularListURL:        "https://submitted.rsmu.ru/data/p3_22074.json",
			SpecialQuotaListURL:   "https://submitted.rsmu.ru/data/p3_22110.json",
			DedicatedQuotaListURL: "https://submitted.rsmu.ru/data/p3_22089.json",
		},

		&rsmu.HTTPHeadingSource{
			ProgramName:           "Фундаментальная медицина",
			RegularListURL:        "https://submitted.rsmu.ru/data/p3_22075.json",
			SpecialQuotaListURL:   "https://submitted.rsmu.ru/data/p3_22112.json",
			DedicatedQuotaListURL: "https://submitted.rsmu.ru/data/p3_22064.json",
		},

		&rsmu.HTTPHeadingSource{
			ProgramName:           "Медицинская биофизика",
			RegularListURL:        "https://submitted.rsmu.ru/data/p3_22177.json",
			SpecialQuotaListURL:   "https://submitted.rsmu.ru/data/p3_22159.json",
			DedicatedQuotaListURL: "https://submitted.rsmu.ru/data/p3_22108.json",
		},

		&rsmu.HTTPHeadingSource{
			ProgramName:           "Биоинформатика",
			RegularListURL:        "https://submitted.rsmu.ru/data/p3_22106.json",
			SpecialQuotaListURL:   "https://submitted.rsmu.ru/data/p3_22081.json",
			DedicatedQuotaListURL: "https://submitted.rsmu.ru/data/p3_22125.json",
		},

		&rsmu.HTTPHeadingSource{
			ProgramName: "Педиатрия",
			TargetQuotaListURLs: []string{
				"https://submitted.rsmu.ru/data/p3_22138.json",
				"https://submitted.rsmu.ru/data/p3_22124.json",
				"https://submitted.rsmu.ru/data/p3_22147.json",
				"https://submitted.rsmu.ru/data/p3_22182.json",
				"https://submitted.rsmu.ru/data/p3_22172.json",
				"https://submitted.rsmu.ru/data/p3_22077.json",
				"https://submitted.rsmu.ru/data/p3_22072.json",
				"https://submitted.rsmu.ru/data/p3_22169.json",
				"https://submitted.rsmu.ru/data/p3_22146.json",
				"https://submitted.rsmu.ru/data/p3_22121.json",
				"https://submitted.rsmu.ru/data/p3_22158.json",
				"https://submitted.rsmu.ru/data/p3_22151.json",
				"https://submitted.rsmu.ru/data/p3_22154.json",
				"https://submitted.rsmu.ru/data/p3_22092.json",
				"https://submitted.rsmu.ru/data/p3_22185.json",
				"https://submitted.rsmu.ru/data/p3_22173.json",
				"https://submitted.rsmu.ru/data/p3_22100.json",
				"https://submitted.rsmu.ru/data/p3_22085.json",
			},
			RegularListURL:        "https://submitted.rsmu.ru/data/p3_22161.json",
			SpecialQuotaListURL:   "https://submitted.rsmu.ru/data/p3_22122.json",
			DedicatedQuotaListURL: "https://submitted.rsmu.ru/data/p3_22156.json",
		},
	}
}
