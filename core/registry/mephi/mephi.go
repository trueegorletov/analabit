package mephi

import (
	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/source/mephi"
)

const (
	varsityCode = "mephi"
	varsityName = "МИФИ"
)

var Varsity = source.VarsityDefinition{
	Code:           varsityCode,
	Name:           varsityName,
	HeadingSources: sourcesList(),
}

func sourcesList() []source.HeadingSource {
	return []source.HeadingSource{
		&mephi.HTTPHeadingSource{
			HeadingName: "Прикладная математика и информатика",
			Capacities: core.Capacities{
				Regular:        37,
				TargetQuota:    6,
				DedicatedQuota: 6,
				SpecialQuota:   6,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12744/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13476/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13474/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13249/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13601/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13178/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13177/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13236/original/no",
			},
			SpecialQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13123/original/no",
			},
		},

		&mephi.HTTPHeadingSource{
			HeadingName: "Информатика и вычислительная техника",
			Capacities: core.Capacities{
				Regular:        35,
				TargetQuota:    6,
				DedicatedQuota: 6,
				SpecialQuota:   6,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12746/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13494/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13251/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13487/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13489/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13493/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13604/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13498/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13238/original/no",
			},
			SpecialQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13125/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Программная инженерия",
			Capacities: core.Capacities{
				Regular:        34,
				TargetQuota:    6,
				DedicatedQuota: 6,
				SpecialQuota:   6,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12746/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13501/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13502/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13490/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13500/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13238/original/no",
			},
			SpecialQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13125/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Информационная безопасность",
			Capacities: core.Capacities{
				Regular:        65,
				TargetQuota:    10,
				DedicatedQuota: 10,
				SpecialQuota:   10,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12747/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13516/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13514/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13515/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13518/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13517/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13513/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13252/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13602/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13519/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13580/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13183/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13239/original/no",
			},
			SpecialQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13126/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Электроника и наноэлектроника",
			Capacities: core.Capacities{
				Regular:        35,
				TargetQuota:    6,
				DedicatedQuota: 6,
				SpecialQuota:   6,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12748/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13521/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13253/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13520/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13240/original/no",
			},
			SpecialQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13127/original/no",
			},
		},

		&mephi.HTTPHeadingSource{
			HeadingName: "Ядерная энергетика и теплофизика",
			Capacities: core.Capacities{
				Regular:        42,
				TargetQuota:    6,
				DedicatedQuota: 6,
				SpecialQuota:   6,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12750/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13538/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13537/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13536/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13539/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13242/original/no",
			},
			SpecialQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13229/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Ядерная физика и технологии",
			Capacities: core.Capacities{
				Regular:        42,
				TargetQuota:    6,
				DedicatedQuota: 6,
				SpecialQuota:   6,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12750/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13255/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13540/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13541/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13543/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13542/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13545/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13603/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13242/original/no",
			},
			SpecialQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13229/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Мехатроника и робототехника",
			Capacities: core.Capacities{
				Regular:        35,
				TargetQuota:    5,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12751/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13256/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13548/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13243/original/no",
			},
			SpecialQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13230/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Системный анализ и управление",
			Capacities: core.Capacities{
				Regular:        23,
				TargetQuota:    4,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12754/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13259/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13565/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13246/original/no",
			},
			SpecialQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13233/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Применение и эксплуатация автоматизированных систем специального назначения",
			Capacities: core.Capacities{
				Regular:        9,
				TargetQuota:    2,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12757/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13267/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13547/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13546/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13608/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13272/original/no",
			},
			SpecialQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13262/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Информационно-аналитические системы безопасности",
			Capacities: core.Capacities{
				Regular:        15,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12758/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13544/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13605/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13610/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13273/original/no",
			},
			SpecialQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13263/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Безопасность информационных технологий в правоохранительной сфере",
			Capacities: core.Capacities{
				Regular:        15,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12758/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13268/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13273/original/no",
			},
			SpecialQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13263/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Электроника и автоматика физических установок",
			Capacities: core.Capacities{
				Regular:        39,
				TargetQuota:    6,
				DedicatedQuota: 6,
				SpecialQuota:   6,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12760/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13270/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13571/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13572/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13275/original/no",
			},
			SpecialQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13265/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Ядерная энергетика и технологии",
			Capacities: core.Capacities{
				Regular:        49,
				TargetQuota:    7,
				DedicatedQuota: 7,
				SpecialQuota:   7,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12759/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13551/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13550/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13553/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13554/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13552/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13269/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13606/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13560/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13562/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13559/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13558/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13549/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13563/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13561/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13564/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13274/original/no",
			},
			SpecialQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13264/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Экономическая безопасность",
			Capacities: core.Capacities{
				Regular:        18,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12761/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13271/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13607/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13276/original/no",
			},
			SpecialQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13266/original/no",
			},
		},

		// Merged:
		// Прикладные математика и физика
		// Физика
		// Физика и астрономия
		&mephi.HTTPHeadingSource{
			HeadingName: "Физика и астрономия",
			Capacities: core.Capacities{
				Regular:        84,
				TargetQuota:    13,
				DedicatedQuota: 13,
				SpecialQuota:   13,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12745/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12765/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13250/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13486/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13485/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13484/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13483/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13477/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13237/original/no",
			},
			SpecialQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13124/original/no",
			},
		},

		//Физико-технические науки и технологии
		//Техническая физика
		//Высокотехнологические плазменные и энергетические установки
		&mephi.HTTPHeadingSource{
			HeadingName: "Физико-технические науки и технологии",
			Capacities: core.Capacities{
				Regular:        36,
				TargetQuota:    5,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12752/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13555/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13556/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13257/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13244/original/no",
			},
			SpecialQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13231/original/no",
			},
		},

		// Приборостроение
		// Merged:
		// Приборостроение
		// Фотоника и оптоинформатика
		// Лазерная техника и лазерные технологии
		// Биотехнические системы и технологии
		&mephi.HTTPHeadingSource{
			HeadingName: "Фотоника, приборостроение, оптические и биотехнические системы и технологии",
			Capacities: core.Capacities{
				Regular:        28,
				TargetQuota:    4,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
			RegularURLs: []string{},
			TargetQuotaURLs: []string{
				// Фотоника и оптоинформатика
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13525/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13526/original/no",
				// Биотехнические системы и технологии
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13530/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13529/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13527/original/no",
				// Лазерная техника и лазерные технологии
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13254/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13533/original/no",
				// Приборостроение
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13523/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13524/original/no",
			},
			DedicatedQuotaURLs: []string{},
			SpecialQuotaURLs:   []string{},
		},

		// Экономика
		// Менеджмент
		// Бизнес-информатика
		// Международные отношения
		&mephi.HTTPHeadingSource{
			HeadingName: "Экономика и управление",
			Capacities: core.Capacities{
				Regular:        10,
				TargetQuota:    2,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12786/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13260/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13570/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13569/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13568/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13567/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13566/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13261/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13248/original/no",
			},
			SpecialQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13235/original/no",
			},
		},
	}
}
