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
			HeadingName: "Прикладные математика и физика",
			Capacities: core.Capacities{
				Regular:        87,
				TargetQuota:    12,
				DedicatedQuota: 12,
				SpecialQuota:   12,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12783/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12783/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13041/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13041/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13197/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13197/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13245/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13245/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Прикладная математика и информатика",
			Capacities: core.Capacities{
				Regular:        40,
				TargetQuota:    5,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12784/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12784/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13042/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13042/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13124/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13124/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13106/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13106/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13198/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13198/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13512/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13512/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13510/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13510/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13511/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13511/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13246/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13246/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Экономика",
			Capacities: core.Capacities{
				Regular:        35,
				TargetQuota:    5,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12785/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12785/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13043/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13043/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13199/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13199/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13247/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13247/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Применение и эксплуатация автом. систем специального назначения",
			Capacities: core.Capacities{
				Regular:        12,
				TargetQuota:    1,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12786/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12786/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13044/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13044/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13200/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13200/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13248/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13248/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Безопасность информационных технологий в правоохранительной сфере",
			Capacities: core.Capacities{
				Regular:        36,
				TargetQuota:    4,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12787/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12787/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13045/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13045/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13201/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13201/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13249/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13249/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Системный анализ  и управление",
			Capacities: core.Capacities{
				Regular:        26,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12790/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12790/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13047/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13047/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13203/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13203/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13251/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13251/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Информационно-аналитические системы безопасности",
			Capacities: core.Capacities{
				Regular:        36,
				TargetQuota:    4,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12791/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12791/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13048/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13048/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13204/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13204/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13252/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13252/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Экономическая безопасность",
			Capacities: core.Capacities{
				Regular:        21,
				TargetQuota:    2,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12792/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12792/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13049/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13049/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13205/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13205/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13253/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13253/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Электроника и автоматика физических установок*",
			Capacities: core.Capacities{
				Regular:        42,
				TargetQuota:    5,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12793/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12793/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13050/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13050/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13206/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13206/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13254/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13254/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Международные отношения",
			Capacities: core.Capacities{
				Regular:        21,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12794/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12794/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13051/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13051/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13207/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13207/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13255/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13255/original/no",
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
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12795/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12795/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13052/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13052/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13208/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13208/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13256/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13256/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Приборостроение",
			Capacities: core.Capacities{
				Regular:        56,
				TargetQuota:    8,
				DedicatedQuota: 8,
				SpecialQuota:   8,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12771/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12771/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13029/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13029/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13523/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13523/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13524/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13524/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13185/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13185/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Техническая физика",
			Capacities: core.Capacities{
				Regular:        29,
				TargetQuota:    4,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13555/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13555/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Ядерные энергетика и теплофизика",
			Capacities: core.Capacities{
				Regular:        84,
				TargetQuota:    12,
				DedicatedQuota: 12,
				SpecialQuota:   12,
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Материаловедение и технологии материалов",
			Capacities: core.Capacities{
				Regular:        14,
				TargetQuota:    2,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12781/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12781/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13039/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13039/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13195/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13195/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Информационная безопасность",
			Capacities: core.Capacities{
				Regular:        68,
				TargetQuota:    9,
				DedicatedQuota: 9,
				SpecialQuota:   9,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12747/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12747/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13126/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13126/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13108/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13108/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12769/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12769/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12789/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12789/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13027/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13027/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13516/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13516/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13514/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13514/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13515/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13515/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13518/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13518/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13517/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13517/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13513/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13513/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13252/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13252/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13602/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13602/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13519/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13519/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13580/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13580/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13183/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13183/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13239/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13239/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Электроника и наноэлектроника",
			Capacities: core.Capacities{
				Regular:        38,
				TargetQuota:    5,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12748/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12748/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13127/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13127/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13109/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13109/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12770/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12770/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13028/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13028/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13521/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13521/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13253/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13253/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13520/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13520/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13184/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13184/original/no",
			},
			DedicatedQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13240/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13240/original/no",
			},
		},
		&mephi.HTTPHeadingSource{
			HeadingName: "Информатика и вычислительная техника",
			Capacities: core.Capacities{
				Regular:        75,
				TargetQuota:    10,
				DedicatedQuota: 10,
				SpecialQuota:   10,
			},
			RegularURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12762/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12762/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12788/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/12788/original/no",
			},
			TargetQuotaURLs: []string{
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13494/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13494/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13251/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13251/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13487/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13487/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13489/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13489/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13493/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13493/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13604/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13604/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13498/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13498/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13575/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13575/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13574/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13574/original/no",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13573/original/yes",
				"https://org.mephi.ru/pupil-rating/get-rating/entity/13573/original/no",
			},
		},
	}
}
