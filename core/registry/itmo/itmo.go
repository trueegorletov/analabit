package itmo

import (
	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/source/itmo"
)

const (
	varsityCode = "itmo"
	varsityName = "ИТМО"
)

var Varsity = source.VarsityDefinition{
	Code:           varsityCode,
	Name:           varsityName,
	HeadingSources: sourcesList(),
}

func sourcesList() []source.HeadingSource {
	return []source.HeadingSource{
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2190",
			PrettyName: "Прикладная математика и информатика",
			Capacities: core.Capacities{Regular: 170}, // fallback total КЦП
		},
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2191",
			PrettyName: "Математическое обеспечение и администрирование информационных систем",
			Capacities: core.Capacities{Regular: 30}, // fallback total КЦП
		},
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2192",
			PrettyName: "Физика",
			Capacities: core.Capacities{Regular: 25}, // fallback total КЦП
		},
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2193",
			PrettyName: "Химия",
			Capacities: core.Capacities{Regular: 25}, // fallback total КЦП
		},
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2194",
			PrettyName: "Экология и природопользование",
			Capacities: core.Capacities{Regular: 25}, // fallback total КЦП
		},
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2195",
			PrettyName: "Информатика и вычислительная техника",
			Capacities: core.Capacities{Regular: 30}, // fallback total КЦП
		},
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2196",
			PrettyName: "Информационные системы и технологии",
			Capacities: core.Capacities{Regular: 151}, // fallback total КЦП
		},
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2197",
			PrettyName: "Прикладная информатика",
			Capacities: core.Capacities{Regular: 28}, // fallback total КЦП
		},
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2198",
			PrettyName: "Программная инженерия",
			Capacities: core.Capacities{Regular: 164}, // fallback total КЦП
		},
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2199",
			PrettyName: "Информационная безопасность",
			Capacities: core.Capacities{Regular: 93}, // fallback total КЦП
		},
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2200",
			PrettyName: "Инфокоммуникационные технологии и системы связи",
			Capacities: core.Capacities{Regular: 81}, // fallback total КЦП
		},
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2201",
			PrettyName: "Конструирование и технология электронных средств",
			Capacities: core.Capacities{Regular: 14}, // fallback total КЦП
		},
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2202",
			PrettyName: "Приборостроение",
			Capacities: core.Capacities{Regular: 16}, // fallback total КЦП
		},
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2203",
			PrettyName: "Фотоника и оптоинформатика и Лазерная техника и лазерные технологии",
			Capacities: core.Capacities{Regular: 86}, // fallback total КЦП
		},
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2204",
			PrettyName: "Биотехнические системы и технологии",
			Capacities: core.Capacities{Regular: 28}, // fallback total КЦП
		},
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2205",
			PrettyName: "Электроэнергетика и электротехника",
			Capacities: core.Capacities{Regular: 16}, // fallback total КЦП
		},
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2206",
			PrettyName: "Автоматизация технологических процессов и производств и Мехатроника и робототехника",
			Capacities: core.Capacities{Regular: 80}, // fallback total КЦП
		},
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2207",
			PrettyName: "Техническая физика",
			Capacities: core.Capacities{Regular: 80}, // fallback total КЦП
		},
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2208",
			PrettyName: "Химическая технология",
			Capacities: core.Capacities{Regular: 45}, // fallback total КЦП
		},
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2209",
			PrettyName: "Энерго- и ресурсосберегающие процессы в химической технологии, нефтехимии и биотехнологии",
			Capacities: core.Capacities{Regular: 20}, // fallback total КЦП
		},
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2210",
			PrettyName: "Биотехнология",
			Capacities: core.Capacities{Regular: 65}, // fallback total КЦП
		},
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2211",
			PrettyName: "Системы управления движением и навигация",
			Capacities: core.Capacities{Regular: 17}, // fallback total КЦП
		},
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2212",
			PrettyName: "Управление в технических системах",
			Capacities: core.Capacities{Regular: 15}, // fallback total КЦП
		},
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2213",
			PrettyName: "Инноватика",
			Capacities: core.Capacities{Regular: 65}, // fallback total КЦП
		},
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2298",
			PrettyName: "Бизнес-информатика",
			Capacities: core.Capacities{Regular: 59}, // fallback total КЦП
		},
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2214",
			PrettyName: "Интеллектуальные системы в гуманитарной сфере",
			Capacities: core.Capacities{Regular: 42}, // fallback total КЦП
		},
		&itmo.HTTPHeadingSource{
			URL:        "https://abit.itmo.ru/rating/bachelor/budget/2215",
			PrettyName: "Дизайн",
			Capacities: core.Capacities{Regular: 10}, // fallback total КЦП
		},
	}
}
