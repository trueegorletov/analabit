package spbsu

import (
	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/source/spbsu"
)

var Varsity = source.VarsityDefinition{
	Code:           "spbsu",
	Name:           "СПбГУ",
	HeadingSources: sourcesList(),
}

func sourcesList() []source.HeadingSource {
	return []source.HeadingSource{
		// Английский язык и литература (с дополнительными квалификациями «Учитель английского языка и литературы» / «Специалист в области перевода»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Английский язык и литература (с дополнительными квалификациями «Учитель английского языка и литературы» / «Специалист в области перевода»)",
			RegularListID:        1920,
			TargetQuotaListIDs:   []int{1776},
			DedicatedQuotaListID: -1,
			SpecialQuotaListID:   1766,
			Capacities: core.Capacities{
				Regular:        3,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   1,
			},
		},

		// Испанский язык, литература и перевод (с дополнительными квалификациями «Учитель иностранного языка и литературы» / «Специалист в области перевода»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Испанский язык, литература и перевод (с дополнительными квалификациями «Учитель иностранного языка и литературы» / «Специалист в области перевода»)",
			RegularListID:        1133,
			TargetQuotaListIDs:   []int{570},
			DedicatedQuotaListID: -1,
			SpecialQuotaListID:   469,
			Capacities: core.Capacities{
				Regular:        3,
				TargetQuota:    1,
				DedicatedQuota: 0,
				SpecialQuota:   1,
			},
		},

		// Английский язык как иностранный в обучении и коммуникации/ English as a Foreign Language in Teaching and Communication (с дополнительными квалификациями «Учитель английского языка» / «Специалист в области перевода»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Английский язык как иностранный в обучении и коммуникации/ English as a Foreign Language in Teaching and Communication (с дополнительными квалификациями «Учитель английского языка» / «Специалист в области перевода»)",
			RegularListID:        338,
			TargetQuotaListIDs:   []int{2096},
			DedicatedQuotaListID: -1,
			SpecialQuotaListID:   1556,
			Capacities: core.Capacities{
				Regular:        3,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   1,
			},
		},

		// Археология (с дополнительными квалификациями «Учитель истории и обществознания» / «Хранитель музейных ценностей» / «Специалист по учету музейных предметов»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Археология (с дополнительными квалификациями «Учитель истории и обществознания» / «Хранитель музейных ценностей» / «Специалист по учету музейных предметов»)",
			RegularListID:        833,
			TargetQuotaListIDs:   []int{1391},
			DedicatedQuotaListID: 945,
			SpecialQuotaListID:   721,
			Capacities: core.Capacities{
				Regular:        6,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// Исламоведение (с дополнительной квалификацией «Специалист в области перевода»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Исламоведение (с дополнительной квалификацией «Специалист в области перевода»)",
			RegularListID:        1981,
			TargetQuotaListIDs:   []int{91},
			DedicatedQuotaListID: -1,
			SpecialQuotaListID:   1561,
			Capacities: core.Capacities{
				Regular:        3,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   1,
			},
		},

		// Международные отношения (с дополнительной квалификацией «Специалист в области перевода»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Международные отношения (с дополнительной квалификацией «Специалист в области перевода»)",
			RegularListID:        1499,
			TargetQuotaListIDs:   []int{1693},
			DedicatedQuotaListID: 417,
			SpecialQuotaListID:   150,
			Capacities: core.Capacities{
				Regular:        35,
				TargetQuota:    0,
				DedicatedQuota: 7,
				SpecialQuota:   7,
			},
		},

		// Международная журналистика (с дополнительной квалификацией «Графический дизайнер»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Международная журналистика (с дополнительной квалификацией «Графический дизайнер»)",
			RegularListID:        726,
			TargetQuotaListIDs:   []int{680},
			DedicatedQuotaListID: 1923,
			SpecialQuotaListID:   971,
			Capacities: core.Capacities{
				Regular:        14,
				TargetQuota:    0,
				DedicatedQuota: 2,
				SpecialQuota:   0,
			},
		},

		// Прикладная математика, программирование и искусственный интеллект (с дополнительной квалификацией «Системный аналитик»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Прикладная математика, программирование и искусственный интеллект (с дополнительной квалификацией «Системный аналитик»)",
			RegularListID:        753,
			TargetQuotaListIDs:   []int{681},
			DedicatedQuotaListID: 830,
			SpecialQuotaListID:   2043,
			Capacities: core.Capacities{
				Regular:        51,
				TargetQuota:    8,
				DedicatedQuota: 8,
				SpecialQuota:   8,
			},
		},

		// AI360: Математика машинного обучения
		&spbsu.HttpHeadingSource{
			PrettyName:           "AI360: Математика машинного обучения",
			RegularListID:        399,
			TargetQuotaListIDs:   []int{963},
			DedicatedQuotaListID: 1600,
			SpecialQuotaListID:   1697,
			Capacities: core.Capacities{
				Regular:        26,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   2,
			},
		},

		// Астрономия
		&spbsu.HttpHeadingSource{
			PrettyName:           "Астрономия",
			RegularListID:        1568,
			TargetQuotaListIDs:   []int{},
			DedicatedQuotaListID: 647,
			SpecialQuotaListID:   1017,
			Capacities: core.Capacities{
				Regular:        16,
				TargetQuota:    0,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},

		// Французский язык, литература и перевод (с дополнительными квалификациями «Учитель французского языка и литературы» / «Специалист в области перевода»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Французский язык, литература и перевод (с дополнительными квалификациями «Учитель французского языка и литературы» / «Специалист в области перевода»)",
			RegularListID:        877,
			TargetQuotaListIDs:   []int{1169},
			DedicatedQuotaListID: -1,
			SpecialQuotaListID:   187,
			Capacities: core.Capacities{
				Regular:        3,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   0,
			},
		},

		// Общая и прикладная фонетика
		&spbsu.HttpHeadingSource{
			PrettyName:           "Общая и прикладная фонетика",
			RegularListID:        2154,
			TargetQuotaListIDs:   []int{309},
			DedicatedQuotaListID: 227,
			SpecialQuotaListID:   1036,
			Capacities: core.Capacities{
				Regular:        4,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// Всеобщая история (с дополнительными квалификациями «Учитель истории и обществознания» / «Педагог дополнительного образования»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Всеобщая история (с дополнительными квалификациями «Учитель истории и обществознания» / «Педагог дополнительного образования»)",
			RegularListID:        1554,
			TargetQuotaListIDs:   []int{738},
			DedicatedQuotaListID: 146,
			SpecialQuotaListID:   624,
			Capacities: core.Capacities{
				Regular:        21,
				TargetQuota:    0,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},

		// Государственное и муниципальное управление
		&spbsu.HttpHeadingSource{
			PrettyName:           "Государственное и муниципальное управление",
			RegularListID:        1974,
			TargetQuotaListIDs:   []int{1040},
			DedicatedQuotaListID: 2150,
			SpecialQuotaListID:   811,
			Capacities: core.Capacities{
				Regular:        8,
				TargetQuota:    3,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},

		// Технологии программирования
		&spbsu.HttpHeadingSource{
			PrettyName:           "Технологии программирования",
			RegularListID:        1779,
			TargetQuotaListIDs:   []int{1084},
			DedicatedQuotaListID: 623,
			SpecialQuotaListID:   97,
			Capacities: core.Capacities{
				Regular:        35,
				TargetQuota:    0,
				DedicatedQuota: 6,
				SpecialQuota:   6,
			},
		},

		// Почвоведение
		&spbsu.HttpHeadingSource{
			PrettyName:           "Почвоведение",
			RegularListID:        1922,
			TargetQuotaListIDs:   []int{622},
			DedicatedQuotaListID: 1825,
			SpecialQuotaListID:   1315,
			Capacities: core.Capacities{
				Regular:        7,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// Менеджмент (с дополнительной квалификацией «Бизнес-аналитик»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Менеджмент (с дополнительной квалификацией «Бизнес-аналитик»)",
			RegularListID:        654,
			TargetQuotaListIDs:   []int{1190, 1969, 506, 2199},
			DedicatedQuotaListID: 1440,
			SpecialQuotaListID:   2243,
			Capacities: core.Capacities{
				Regular:        51,
				TargetQuota:    5,
				DedicatedQuota: 8,
				SpecialQuota:   8,
			},
		},

		// Отечественная филология (Русский язык и литература) (с дополнительной квалификацией «Учитель русского языка и литературы»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Отечественная филология (Русский язык и литература) (с дополнительной квалификацией «Учитель русского языка и литературы»)",
			RegularListID:        1112,
			TargetQuotaListIDs:   []int{167},
			DedicatedQuotaListID: 33,
			SpecialQuotaListID:   791,
			Capacities: core.Capacities{
				Regular:        47,
				TargetQuota:    4,
				DedicatedQuota: 7,
				SpecialQuota:   7,
			},
		},

		// Прикладная, компьютерная и математическая лингвистика (английский язык)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Прикладная, компьютерная и математическая лингвистика (английский язык)",
			RegularListID:        1861,
			TargetQuotaListIDs:   []int{795},
			DedicatedQuotaListID: 922,
			SpecialQuotaListID:   1207,
			Capacities: core.Capacities{
				Regular:        5,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// Религиоведение (с дополнительными квалификациями «Учитель обществознания» / «Учитель основ религиозных культур и светской этики» / «Педагог дополнительного образования» / «Экскурсовод (гид)»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Религиоведение (с дополнительными квалификациями «Учитель обществознания» / «Учитель основ религиозных культур и светской этики» / «Педагог дополнительного образования» / «Экскурсовод (гид)»)",
			RegularListID:        2012,
			TargetQuotaListIDs:   []int{882},
			DedicatedQuotaListID: 76,
			SpecialQuotaListID:   129,
			Capacities: core.Capacities{
				Regular:        7,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// История Центральной Азии (с дополнительной квалификацией «Специалист в области перевода»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "История Центральной Азии (с дополнительной квалификацией «Специалист в области перевода»)",
			RegularListID:        852,
			TargetQuotaListIDs:   []int{1450},
			DedicatedQuotaListID: 2120,
			SpecialQuotaListID:   1512,
			Capacities: core.Capacities{
				Regular:        2,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// История Кавказа (Армения) (с дополнительной квалификацией «Специалист в области перевода»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "История Кавказа (Армения) (с дополнительной квалификацией «Специалист в области перевода»)",
			RegularListID:        1157,
			TargetQuotaListIDs:   []int{1955},
			DedicatedQuotaListID: 1716,
			SpecialQuotaListID:   1668,
			Capacities: core.Capacities{
				Regular:        2,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   1,
			},
		},

		// Науки о данных (с дополнительной квалификацией «Специалист по большим данным»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Науки о данных (с дополнительной квалификацией «Специалист по большим данным»)",
			RegularListID:        407,
			TargetQuotaListIDs:   []int{1137},
			DedicatedQuotaListID: 239,
			SpecialQuotaListID:   1216,
			Capacities: core.Capacities{
				Regular:        20,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   2,
			},
		},

		// Физика (с дополнительной квалификацией «Программист»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Физика (с дополнительной квалификацией «Программист»)",
			RegularListID:        1333,
			TargetQuotaListIDs:   []int{1206},
			DedicatedQuotaListID: 418,
			SpecialQuotaListID:   594,
			Capacities: core.Capacities{
				Regular:        67,
				TargetQuota:    5,
				DedicatedQuota: 9,
				SpecialQuota:   9,
			},
		},

		// Гидрометеорология
		&spbsu.HttpHeadingSource{
			PrettyName:           "Гидрометеорология",
			RegularListID:        864,
			TargetQuotaListIDs:   []int{1318, 933},
			DedicatedQuotaListID: 571,
			SpecialQuotaListID:   942,
			Capacities: core.Capacities{
				Regular:        9,
				TargetQuota:    1,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},

		// Прикладные компьютерные технологии и искусственный интеллект
		&spbsu.HttpHeadingSource{
			PrettyName:           "Прикладные компьютерные технологии и искусственный интеллект",
			RegularListID:        1109,
			TargetQuotaListIDs:   []int{1387},
			DedicatedQuotaListID: 1217,
			SpecialQuotaListID:   282,
			Capacities: core.Capacities{
				Regular:        15,
				TargetQuota:    0,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},

		// Психология служебной деятельности
		&spbsu.HttpHeadingSource{
			PrettyName:           "Психология служебной деятельности",
			RegularListID:        983,
			TargetQuotaListIDs:   []int{1381},
			DedicatedQuotaListID: 542,
			SpecialQuotaListID:   1605,
			Capacities: core.Capacities{
				Regular:        7,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// Юриспруденция (с углубленным изучением китайского языка и права КНР) (с дополнительной квалификацией «Специалист в области перевода»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Юриспруденция (с углубленным изучением китайского языка и права КНР) (с дополнительной квалификацией «Специалист в области перевода»)",
			RegularListID:        1044,
			TargetQuotaListIDs:   []int{1293},
			DedicatedQuotaListID: 1064,
			SpecialQuotaListID:   1588,
			Capacities: core.Capacities{
				Regular:        14,
				TargetQuota:    2,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},

		// Новогреческий язык, византийская и новогреческая филология
		&spbsu.HttpHeadingSource{
			PrettyName:           "Новогреческий язык, византийская и новогреческая филология",
			RegularListID:        948,
			TargetQuotaListIDs:   []int{2046},
			DedicatedQuotaListID: 479,
			SpecialQuotaListID:   1872,
			Capacities: core.Capacities{
				Regular:        2,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   0,
			},
		},

		// Русский язык как иностранный (с дополнительной квалификацией «Учитель русского языка и литературы»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Русский язык как иностранный (с дополнительной квалификацией «Учитель русского языка и литературы»)",
			RegularListID:        27,
			TargetQuotaListIDs:   []int{1124},
			DedicatedQuotaListID: -1,
			SpecialQuotaListID:   1952,
			Capacities: core.Capacities{
				Regular:        2,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   0,
			},
		},

		// Картография и геоинформатика
		&spbsu.HttpHeadingSource{
			PrettyName:           "Картография и геоинформатика",
			RegularListID:        264,
			TargetQuotaListIDs:   []int{973},
			DedicatedQuotaListID: 1148,
			SpecialQuotaListID:   1650,
			Capacities: core.Capacities{
				Regular:        9,
				TargetQuota:    2,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},

		// Языки балто-славянского культурного пространства (литовский, польский и белорусский языки)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Языки балто-славянского культурного пространства (литовский, польский и белорусский языки)",
			RegularListID:        904,
			TargetQuotaListIDs:   []int{1641},
			DedicatedQuotaListID: 1408,
			SpecialQuotaListID:   1473,
			Capacities: core.Capacities{
				Regular:        2,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   0,
			},
		},

		// Перевод и переводоведение (английский язык и языки романского культурного ареала)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Перевод и переводоведение (английский язык и языки романского культурного ареала)",
			RegularListID:        747,
			TargetQuotaListIDs:   []int{1931},
			DedicatedQuotaListID: 283,
			SpecialQuotaListID:   2000,
			Capacities: core.Capacities{
				Regular:        14,
				TargetQuota:    0,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},

		// Культурология (с дополнительными квалификациями «Учитель обществознания» / «Педагог дополнительного образования»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Культурология (с дополнительными квалификациями «Учитель обществознания» / «Педагог дополнительного образования»)",
			RegularListID:        682,
			TargetQuotaListIDs:   []int{},
			DedicatedQuotaListID: 1313,
			SpecialQuotaListID:   2108,
			Capacities: core.Capacities{
				Regular:        16,
				TargetQuota:    0,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},

		// Музеология и охрана объектов культурного и природного наследия (с дополнительной квалификацией «Экскурсовод (гид)»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Музеология и охрана объектов культурного и природного наследия (с дополнительной квалификацией «Экскурсовод (гид)»)",
			RegularListID:        1496,
			TargetQuotaListIDs:   []int{1979},
			DedicatedQuotaListID: 1047,
			SpecialQuotaListID:   2011,
			Capacities: core.Capacities{
				Regular:        10,
				TargetQuota:    1,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},

		// Ассириология (языки, история и культура Древнего Ближнего Востока) (с дополнительной квалификацией «Специалист в области перевода»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Ассириология (языки, история и культура Древнего Ближнего Востока) (с дополнительной квалификацией «Специалист в области перевода»)",
			RegularListID:        1008,
			TargetQuotaListIDs:   []int{1356},
			DedicatedQuotaListID: 243,
			SpecialQuotaListID:   1432,
			Capacities: core.Capacities{
				Regular:        2,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// Японская филология (с дополнительной квалификацией «Специалист в области перевода»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Японская филология (с дополнительной квалификацией «Специалист в области перевода»)",
			RegularListID:        376,
			TargetQuotaListIDs:   []int{},
			DedicatedQuotaListID: -1,
			SpecialQuotaListID:   1617,
			Capacities: core.Capacities{
				Regular:        5,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   1,
			},
		},

		// Реклама и связи с общественностью (с дополнительной квалификацией «Специалист по продвижению и распространению продукции средств массовой информации»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Реклама и связи с общественностью (с дополнительной квалификацией «Специалист по продвижению и распространению продукции средств массовой информации»)",
			RegularListID:        1550,
			TargetQuotaListIDs:   []int{1429},
			DedicatedQuotaListID: 2166,
			SpecialQuotaListID:   2204,
			Capacities: core.Capacities{
				Regular:        18,
				TargetQuota:    1,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},

		// Программная инженерия
		&spbsu.HttpHeadingSource{
			PrettyName:           "Программная инженерия",
			RegularListID:        888,
			TargetQuotaListIDs:   []int{573, 697},
			DedicatedQuotaListID: 2229,
			SpecialQuotaListID:   65,
			Capacities: core.Capacities{
				Regular:        28,
				TargetQuota:    7,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
		},

		// Стоматология
		&spbsu.HttpHeadingSource{
			PrettyName:           "Стоматология",
			RegularListID:        1170,
			TargetQuotaListIDs:   []int{1621, 414},
			DedicatedQuotaListID: 1715,
			SpecialQuotaListID:   1765,
			Capacities: core.Capacities{
				Regular:        3,
				TargetQuota:    13,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},

		// Организация туристской деятельности (с изучением Индии / Арабских стран) (с дополнительной квалификацией «Руководитель гостиничных комплексов и иных средств размещения»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Организация туристской деятельности (с изучением Индии / Арабских стран) (с дополнительной квалификацией «Руководитель гостиничных комплексов и иных средств размещения»)",
			RegularListID:        1558,
			TargetQuotaListIDs:   []int{},
			DedicatedQuotaListID: -1,
			SpecialQuotaListID:   850,
			Capacities: core.Capacities{
				Regular:        4,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   1,
			},
		},

		// Скандинавские языки и литературы (шведский и норвежский языки) (с дополнительными квалификациями «Учитель иностранного языка и литературы» / «Специалист в области перевода»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Скандинавские языки и литературы (шведский и норвежский языки) (с дополнительными квалификациями «Учитель иностранного языка и литературы» / «Специалист в области перевода»)",
			RegularListID:        1949,
			TargetQuotaListIDs:   []int{1863},
			DedicatedQuotaListID: 1218,
			SpecialQuotaListID:   1456,
			Capacities: core.Capacities{
				Regular:        6,
				TargetQuota:    0,
				DedicatedQuota: 2,
				SpecialQuota:   1,
			},
		},

		// Теоретическое и экспериментальное языкознание (английский язык) (с дополнительной квалификацией «Учитель иностранного языка»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Теоретическое и экспериментальное языкознание (английский язык) (с дополнительной квалификацией «Учитель иностранного языка»)",
			RegularListID:        312,
			TargetQuotaListIDs:   []int{540},
			DedicatedQuotaListID: 1565,
			SpecialQuotaListID:   123,
			Capacities: core.Capacities{
				Regular:        2,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// Иудейская теология (с дополнительными квалификациями «Педагог дополнительного образования» / «Учитель основ религиозных культур и светской этики»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Иудейская теология (с дополнительными квалификациями «Педагог дополнительного образования» / «Учитель основ религиозных культур и светской этики»)",
			RegularListID:        737,
			TargetQuotaListIDs:   []int{1900},
			DedicatedQuotaListID: 1092,
			SpecialQuotaListID:   1373,
			Capacities: core.Capacities{
				Regular:        5,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   0,
			},
		},

		// История Ирана и Афганистана (с дополнительной квалификацией «Специалист в области перевода»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "История Ирана и Афганистана (с дополнительной квалификацией «Специалист в области перевода»)",
			RegularListID:        1497,
			TargetQuotaListIDs:   []int{2036},
			DedicatedQuotaListID: 1103,
			SpecialQuotaListID:   1366,
			Capacities: core.Capacities{
				Regular:        3,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// Большие данные и распределенная цифровая платформа
		&spbsu.HttpHeadingSource{
			PrettyName:           "Большие данные и распределенная цифровая платформа",
			RegularListID:        1152,
			TargetQuotaListIDs:   []int{2143},
			DedicatedQuotaListID: 173,
			SpecialQuotaListID:   2077,
			Capacities: core.Capacities{
				Regular:        33,
				TargetQuota:    0,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},

		// Прикладные физика и математика (с дополнительной квалификацией «Программист»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Прикладные физика и математика (с дополнительной квалификацией «Программист»)",
			RegularListID:        566,
			TargetQuotaListIDs:   []int{1774},
			DedicatedQuotaListID: 447,
			SpecialQuotaListID:   1396,
			Capacities: core.Capacities{
				Regular:        21,
				TargetQuota:    0,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},

		// Химическое материаловедение (с дополнительной квалификацией «Специалист в области перевода научно-технической литературы»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Химическое материаловедение (с дополнительной квалификацией «Специалист в области перевода научно-технической литературы»)",
			RegularListID:        1415,
			TargetQuotaListIDs:   []int{1182, 405},
			DedicatedQuotaListID: 766,
			SpecialQuotaListID:   171,
			Capacities: core.Capacities{
				Regular:        13,
				TargetQuota:    0,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},

		// Монгольско-тибетская филология (с дополнительной квалификацией «Специалист в области перевода»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Монгольско-тибетская филология (с дополнительной квалификацией «Специалист в области перевода»)",
			RegularListID:        465,
			TargetQuotaListIDs:   []int{},
			DedicatedQuotaListID: 301,
			SpecialQuotaListID:   2062,
			Capacities: core.Capacities{
				Regular:        3,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   1,
			},
		},

		// Искусственный интеллект и наука о данных
		&spbsu.HttpHeadingSource{
			PrettyName:           "Искусственный интеллект и наука о данных",
			RegularListID:        1364,
			TargetQuotaListIDs:   []int{1029, 1714, 597},
			DedicatedQuotaListID: 1812,
			SpecialQuotaListID:   2065,
			Capacities: core.Capacities{
				Regular:        40,
				TargetQuota:    4,
				DedicatedQuota: 4,
				SpecialQuota:   6,
			},
		},

		// Прикладная информатика в области искусств и гуманитарных наук (с дополнительной квалификацией «Специалист по информационным системам»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Прикладная информатика в области искусств и гуманитарных наук (с дополнительной квалификацией «Специалист по информационным системам»)",
			RegularListID:        323,
			TargetQuotaListIDs:   []int{1525},
			DedicatedQuotaListID: 961,
			SpecialQuotaListID:   1744,
			Capacities: core.Capacities{
				Regular:        4,
				TargetQuota:    0,
				DedicatedQuota: 3,
				SpecialQuota:   1,
			},
		},

		// Психология
		&spbsu.HttpHeadingSource{
			PrettyName:           "Психология",
			RegularListID:        637,
			TargetQuotaListIDs:   []int{205},
			DedicatedQuotaListID: 1270,
			SpecialQuotaListID:   800,
			Capacities: core.Capacities{
				Regular:        35,
				TargetQuota:    5,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
		},

		// Клиническая психология
		&spbsu.HttpHeadingSource{
			PrettyName:           "Клиническая психология",
			RegularListID:        1723,
			TargetQuotaListIDs:   []int{2184},
			DedicatedQuotaListID: 860,
			SpecialQuotaListID:   1337,
			Capacities: core.Capacities{
				Regular:        28,
				TargetQuota:    0,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},

		// Юриспруденция
		&spbsu.HttpHeadingSource{
			PrettyName:           "Юриспруденция",
			RegularListID:        1819,
			TargetQuotaListIDs:   []int{1255, 923, 1618, 1821},
			DedicatedQuotaListID: 2157,
			SpecialQuotaListID:   736,
			Capacities: core.Capacities{
				Regular:        66,
				TargetQuota:    7,
				DedicatedQuota: 11,
				SpecialQuota:   11,
			},
		},

		// Иностранные языки в глобальных коммуникациях (с дополнительными квалификациями «Учитель иностранного языка» / «Специалист в области перевода»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Иностранные языки в глобальных коммуникациях (с дополнительными квалификациями «Учитель иностранного языка» / «Специалист в области перевода»)",
			RegularListID:        2172,
			TargetQuotaListIDs:   []int{711},
			DedicatedQuotaListID: -1,
			SpecialQuotaListID:   1004,
			Capacities: core.Capacities{
				Regular:        3,
				TargetQuota:    1,
				DedicatedQuota: 0,
				SpecialQuota:   1,
			},
		},

		// Физическая культура и спорт (с дополнительной квалификацией «Специалист по фитнесу (фитнес-тренер)»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Физическая культура и спорт (с дополнительной квалификацией «Специалист по фитнесу (фитнес-тренер)»)",
			RegularListID:        1365,
			TargetQuotaListIDs:   []int{1829},
			DedicatedQuotaListID: 9,
			SpecialQuotaListID:   1420,
			Capacities: core.Capacities{
				Regular:        7,
				TargetQuota:    0,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},

		// Искусства и гуманитарные науки (с дополнительной квалификацией «Педагог дополнительного образования»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Искусства и гуманитарные науки (с дополнительной квалификацией «Педагог дополнительного образования»)",
			RegularListID:        1278,
			TargetQuotaListIDs:   []int{},
			DedicatedQuotaListID: 245,
			SpecialQuotaListID:   896,
			Capacities: core.Capacities{
				Regular:        35,
				TargetQuota:    0,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
		},

		// Прикладная математика, процессы управления и искусственный интеллект (с дополнительной квалификацией «Программист»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Прикладная математика, процессы управления и искусственный интеллект (с дополнительной квалификацией «Программист»)",
			RegularListID:        787,
			TargetQuotaListIDs:   []int{1604},
			DedicatedQuotaListID: 1316,
			SpecialQuotaListID:   1433,
			Capacities: core.Capacities{
				Regular:        106,
				TargetQuota:    14,
				DedicatedQuota: 19,
				SpecialQuota:   19,
			},
		},

		// Фундаментальная механика (с дополнительной квалификацией «Специалист по научно-исследовательским разработкам»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Фундаментальная механика (с дополнительной квалификацией «Специалист по научно-исследовательским разработкам»)",
			RegularListID:        1266,
			TargetQuotaListIDs:   []int{1694},
			DedicatedQuotaListID: 252,
			SpecialQuotaListID:   1220,
			Capacities: core.Capacities{
				Regular:        30,
				TargetQuota:    0,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},

		// Геология
		&spbsu.HttpHeadingSource{
			PrettyName:           "Геология",
			RegularListID:        478,
			TargetQuotaListIDs:   []int{2191},
			DedicatedQuotaListID: 1772,
			SpecialQuotaListID:   1061,
			Capacities: core.Capacities{
				Regular:        33,
				TargetQuota:    2,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
		},

		// Лечебное дело
		&spbsu.HttpHeadingSource{
			PrettyName:           "Лечебное дело",
			RegularListID:        1769,
			TargetQuotaListIDs:   []int{1458, 1629, 103},
			DedicatedQuotaListID: 751,
			SpecialQuotaListID:   216,
			Capacities: core.Capacities{
				Regular:        5,
				TargetQuota:    32,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
		},

		// Бизнес-информатика
		&spbsu.HttpHeadingSource{
			PrettyName:           "Бизнес-информатика",
			RegularListID:        1298,
			TargetQuotaListIDs:   []int{1741},
			DedicatedQuotaListID: 1395,
			SpecialQuotaListID:   789,
			Capacities: core.Capacities{
				Regular:        9,
				TargetQuota:    2,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},

		// Туризм (с дополнительной квалификацией «Руководитель гостиничного комплекса и иных средств размещения»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Туризм (с дополнительной квалификацией «Руководитель гостиничного комплекса и иных средств размещения»)",
			RegularListID:        111,
			TargetQuotaListIDs:   []int{62},
			DedicatedQuotaListID: 1046,
			SpecialQuotaListID:   1860,
			Capacities: core.Capacities{
				Regular:        2,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// Христианская теология (с дополнительной квалификацией «Учитель основ религиозных культур и светской этики»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Христианская теология (с дополнительной квалификацией «Учитель основ религиозных культур и светской этики»)",
			RegularListID:        1951,
			TargetQuotaListIDs:   []int{},
			DedicatedQuotaListID: 1519,
			SpecialQuotaListID:   1106,
			Capacities: core.Capacities{
				Regular:        5,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   1,
			},
		},

		// Китайская филология (с дополнительной квалификацией «Специалист в области перевода»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Китайская филология (с дополнительной квалификацией «Специалист в области перевода»)",
			RegularListID:        125,
			TargetQuotaListIDs:   []int{1916},
			DedicatedQuotaListID: -1,
			SpecialQuotaListID:   2195,
			Capacities: core.Capacities{
				Regular:        4,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   1,
			},
		},

		// Современное программирование (с дополнительной квалификацией «Программист»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Современное программирование (с дополнительной квалификацией «Программист»)",
			RegularListID:        1856,
			TargetQuotaListIDs:   []int{1603},
			DedicatedQuotaListID: 555,
			SpecialQuotaListID:   1080,
			Capacities: core.Capacities{
				Regular:        34,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// Химия (с дополнительными квалификациями «Специалист в области перевода научно-технической литературы» / «Учитель химии»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Химия (с дополнительными квалификациями «Специалист в области перевода научно-технической литературы» / «Учитель химии»)",
			RegularListID:        1475,
			TargetQuotaListIDs:   []int{199, 619, 1611, 1806, 373, 1740},
			DedicatedQuotaListID: 1457,
			SpecialQuotaListID:   1917,
			Capacities: core.Capacities{
				Regular:        48,
				TargetQuota:    5,
				DedicatedQuota: 8,
				SpecialQuota:   8,
			},
		},

		// Фармация
		&spbsu.HttpHeadingSource{
			PrettyName:           "Фармация",
			RegularListID:        1372,
			TargetQuotaListIDs:   []int{102},
			DedicatedQuotaListID: 2118,
			SpecialQuotaListID:   1122,
			Capacities: core.Capacities{
				Regular:        2,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// Экономика (с углубленным изучением экономики Китая и китайского языка)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Экономика (с углубленным изучением экономики Китая и китайского языка)",
			RegularListID:        1104,
			TargetQuotaListIDs:   []int{1154},
			DedicatedQuotaListID: 964,
			SpecialQuotaListID:   361,
			Capacities: core.Capacities{
				Regular:        11,
				TargetQuota:    1,
				DedicatedQuota: 1,
				SpecialQuota:   2,
			},
		},

		// Экономико-математические методы
		&spbsu.HttpHeadingSource{
			PrettyName:           "Экономико-математические методы",
			RegularListID:        1593,
			TargetQuotaListIDs:   []int{2088},
			DedicatedQuotaListID: 208,
			SpecialQuotaListID:   831,
			Capacities: core.Capacities{
				Regular:        6,
				TargetQuota:    1,
				DedicatedQuota: 2,
				SpecialQuota:   1,
			},
		},

		// Организация туристской деятельности (с углубленным изучением китайского языка) (с дополнительной квалификацией «Руководитель гостиничного комплекса и иных средств размещения»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Организация туристской деятельности (с углубленным изучением китайского языка) (с дополнительной квалификацией «Руководитель гостиничного комплекса и иных средств размещения»)",
			RegularListID:        532,
			TargetQuotaListIDs:   []int{},
			DedicatedQuotaListID: 610,
			SpecialQuotaListID:   634,
			Capacities: core.Capacities{
				Regular:        3,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// Теория и практика межкультурной коммуникации (английский язык) (с дополнительными квалификациями «Учитель английского языка» / «Специалист в области перевода»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Теория и практика межкультурной коммуникации (английский язык) (с дополнительными квалификациями «Учитель английского языка» / «Специалист в области перевода»)",
			RegularListID:        1026,
			TargetQuotaListIDs:   []int{1925},
			DedicatedQuotaListID: -1,
			SpecialQuotaListID:   656,
			Capacities: core.Capacities{
				Regular:        3,
				TargetQuota:    1,
				DedicatedQuota: 0,
				SpecialQuota:   1,
			},
		},

		// Теория и история искусств (с дополнительной квалификацией «Экскурсовод (гид)»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Теория и история искусств (с дополнительной квалификацией «Экскурсовод (гид)»)",
			RegularListID:        720,
			TargetQuotaListIDs:   []int{},
			DedicatedQuotaListID: 1334,
			SpecialQuotaListID:   2224,
			Capacities: core.Capacities{
				Regular:        8,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// Журналистика (с дополнительной квалификацией «Специалист по информационным ресурсам»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Журналистика (с дополнительной квалификацией «Специалист по информационным ресурсам»)",
			RegularListID:        1737,
			TargetQuotaListIDs:   []int{182},
			DedicatedQuotaListID: 1656,
			SpecialQuotaListID:   1704,
			Capacities: core.Capacities{
				Regular:        16,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},

		// Графический дизайн (с дополнительной квалификацией «Графический дизайнер»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Графический дизайн (с дополнительной квалификацией «Графический дизайнер»)",
			RegularListID:        2207,
			TargetQuotaListIDs:   []int{},
			DedicatedQuotaListID: 1982,
			SpecialQuotaListID:   2176,
			Capacities: core.Capacities{
				Regular:        9,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// Корейская филология (с дополнительной квалификацией «Специалист в области перевода»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Корейская филология (с дополнительной квалификацией «Специалист в области перевода»)",
			RegularListID:        1562,
			TargetQuotaListIDs:   []int{},
			DedicatedQuotaListID: -1,
			SpecialQuotaListID:   280,
			Capacities: core.Capacities{
				Regular:        4,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   1,
			},
		},

		// Кхмерская филология (с дополнительной квалификацией «Специалист в области перевода»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Кхмерская филология (с дополнительной квалификацией «Специалист в области перевода»)",
			RegularListID:        803,
			TargetQuotaListIDs:   []int{},
			DedicatedQuotaListID: -1,
			SpecialQuotaListID:   400,
			Capacities: core.Capacities{
				Regular:        4,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   1,
			},
		},

		// Языки и культура Западной Африки (языки манде) (с дополнительной квалификацией «Специалист в области перевода»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Языки и культура Западной Африки (языки манде) (с дополнительной квалификацией «Специалист в области перевода»)",
			RegularListID:        537,
			TargetQuotaListIDs:   []int{1382},
			DedicatedQuotaListID: 617,
			SpecialQuotaListID:   1734,
			Capacities: core.Capacities{
				Regular:        3,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// География (с дополнительной квалификацией «Учитель географии»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "География (с дополнительной квалификацией «Учитель географии»)",
			RegularListID:        1508,
			TargetQuotaListIDs:   []int{412},
			DedicatedQuotaListID: 684,
			SpecialQuotaListID:   770,
			Capacities: core.Capacities{
				Regular:        30,
				TargetQuota:    2,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},

		// Экономика
		&spbsu.HttpHeadingSource{
			PrettyName:           "Экономика",
			RegularListID:        464,
			TargetQuotaListIDs:   []int{772, 1546, 696},
			DedicatedQuotaListID: 426,
			SpecialQuotaListID:   1498,
			Capacities: core.Capacities{
				Regular:        44,
				TargetQuota:    11,
				DedicatedQuota: 7,
				SpecialQuota:   7,
			},
		},

		// Социологические исследования в цифровом обществе (с дополнительной квалификацией «Социолог: специалист по фундаментальным и прикладным социологическим исследованиям»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Социологические исследования в цифровом обществе (с дополнительной квалификацией «Социолог: специалист по фундаментальным и прикладным социологическим исследованиям»)",
			RegularListID:        936,
			TargetQuotaListIDs:   []int{1807},
			DedicatedQuotaListID: 191,
			SpecialQuotaListID:   268,
			Capacities: core.Capacities{
				Regular:        10,
				TargetQuota:    0,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},

		// Немецкий язык, литература и перевод (с дополнительными квалификациями «Учитель немецкого языка и литературы» / «Специалист в области перевода»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Немецкий язык, литература и перевод (с дополнительными квалификациями «Учитель немецкого языка и литературы» / «Специалист в области перевода»)",
			RegularListID:        1811,
			TargetQuotaListIDs:   []int{507},
			DedicatedQuotaListID: -1,
			SpecialQuotaListID:   1761,
			Capacities: core.Capacities{
				Regular:        3,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   1,
			},
		},

		// История России (с дополнительными квалификациями «Учитель истории и обществознания» / «Педагог дополнительного образования»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "История России (с дополнительными квалификациями «Учитель истории и обществознания» / «Педагог дополнительного образования»)",
			RegularListID:        2028,
			TargetQuotaListIDs:   []int{1612},
			DedicatedQuotaListID: 134,
			SpecialQuotaListID:   1097,
			Capacities: core.Capacities{
				Regular:        31,
				TargetQuota:    6,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},

		// Реставрация произведений изобразительного и декоративно-прикладного искусства (с дополнительной квалификацией «Специалист по техническим процессам художественной деятельности»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Реставрация произведений изобразительного и декоративно-прикладного искусства (с дополнительной квалификацией «Специалист по техническим процессам художественной деятельности»)",
			RegularListID:        678,
			TargetQuotaListIDs:   []int{},
			DedicatedQuotaListID: 107,
			SpecialQuotaListID:   1542,
			Capacities: core.Capacities{
				Regular:        4,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   0,
			},
		},

		// Художник кино и телевидения по костюму (с дополнительными квалификациями «Учитель рисования» / «Педагог дополнительного образования»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Художник кино и телевидения по костюму (с дополнительными квалификациями «Учитель рисования» / «Педагог дополнительного образования»)",
			RegularListID:        1607,
			TargetQuotaListIDs:   []int{},
			DedicatedQuotaListID: 2162,
			SpecialQuotaListID:   1359,
			Capacities: core.Capacities{
				Regular:        3,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   1,
			},
		},

		// Филология телугу (телугу, санскрит, хинди) (с дополнительной квалификацией «Специалист в области перевода»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Филология телугу (телугу, санскрит, хинди) (с дополнительной квалификацией «Специалист в области перевода»)",
			RegularListID:        262,
			TargetQuotaListIDs:   []int{},
			DedicatedQuotaListID: 1970,
			SpecialQuotaListID:   1967,
			Capacities: core.Capacities{
				Regular:        4,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// Социальная работа (с дополнительной квалификацией «Специалист по работе с семьей»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Социальная работа (с дополнительной квалификацией «Специалист по работе с семьей»)",
			RegularListID:        954,
			TargetQuotaListIDs:   []int{695},
			DedicatedQuotaListID: 1833,
			SpecialQuotaListID:   874,
			Capacities: core.Capacities{
				Regular:        6,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// Философия (с дополнительными квалификациями «Учитель обществознания» / «Педагог дополнительного образования»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Философия (с дополнительными квалификациями «Учитель обществознания» / «Педагог дополнительного образования»)",
			RegularListID:        1889,
			TargetQuotaListIDs:   []int{1183},
			DedicatedQuotaListID: 1749,
			SpecialQuotaListID:   2035,
			Capacities: core.Capacities{
				Regular:        30,
				TargetQuota:    0,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},

		// Декоративно-прикладное искусство (с дополнительной квалификацией «Педагог дополнительного образования»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Декоративно-прикладное искусство (с дополнительной квалификацией «Педагог дополнительного образования»)",
			RegularListID:        2111,
			TargetQuotaListIDs:   []int{},
			DedicatedQuotaListID: 630,
			SpecialQuotaListID:   7,
			Capacities: core.Capacities{
				Regular:        6,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   0,
			},
		},

		// История арабских стран (с дополнительной квалификацией «Специалист в области перевода»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "История арабских стран (с дополнительной квалификацией «Специалист в области перевода»)",
			RegularListID:        1421,
			TargetQuotaListIDs:   []int{292},
			DedicatedQuotaListID: -1,
			SpecialQuotaListID:   600,
			Capacities: core.Capacities{
				Regular:        3,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   1,
			},
		},

		// Османистика (языки, история и культура османской Турции) (с дополнительной квалификацией «Специалист в области перевода»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Османистика (языки, история и культура османской Турции) (с дополнительной квалификацией «Специалист в области перевода»)",
			RegularListID:        117,
			TargetQuotaListIDs:   []int{},
			DedicatedQuotaListID: 180,
			SpecialQuotaListID:   2091,
			Capacities: core.Capacities{
				Regular:        3,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// Математика (с дополнительной квалификацией «Специалист по научно-исследовательским разработкам»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Математика (с дополнительной квалификацией «Специалист по научно-исследовательским разработкам»)",
			RegularListID:        1248,
			TargetQuotaListIDs:   []int{713, 842},
			DedicatedQuotaListID: 1229,
			SpecialQuotaListID:   1705,
			Capacities: core.Capacities{
				Regular:        37,
				TargetQuota:    1,
				DedicatedQuota: 6,
				SpecialQuota:   6,
			},
		},

		// Биология (с дополнительной квалификацией «Учитель биологии»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Биология (с дополнительной квалификацией «Учитель биологии»)",
			RegularListID:        716,
			TargetQuotaListIDs:   []int{81, 1670},
			DedicatedQuotaListID: 2232,
			SpecialQuotaListID:   878,
			Capacities: core.Capacities{
				Regular:        75,
				TargetQuota:    5,
				DedicatedQuota: 10,
				SpecialQuota:   10,
			},
		},

		// Кадастр недвижимости: оценка и информационное обеспечение (с дополнительной квалификацией «Специалист в области инженерно-геодезических изысканий для градостроительной деятельности»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Кадастр недвижимости: оценка и информационное обеспечение (с дополнительной квалификацией «Специалист в области инженерно-геодезических изысканий для градостроительной деятельности»)",
			RegularListID:        1660,
			TargetQuotaListIDs:   []int{1006},
			DedicatedQuotaListID: 2196,
			SpecialQuotaListID:   2048,
			Capacities: core.Capacities{
				Regular:        7,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// Управление персоналом (с дополнительной квалификацией «Консультант в области управления персоналом»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Управление персоналом (с дополнительной квалификацией «Консультант в области управления персоналом»)",
			RegularListID:        2233,
			TargetQuotaListIDs:   []int{105},
			DedicatedQuotaListID: 980,
			SpecialQuotaListID:   1135,
			Capacities: core.Capacities{
				Regular:        6,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// Прикладная социология
		&spbsu.HttpHeadingSource{
			PrettyName:           "Прикладная социология",
			RegularListID:        367,
			TargetQuotaListIDs:   []int{603},
			DedicatedQuotaListID: 1632,
			SpecialQuotaListID:   1438,
			Capacities: core.Capacities{
				Regular:        7,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// Социология (с дополнительной квалификацией «Социолог: специалист по фундаментальным и прикладным социологическим исследованиям»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Социология (с дополнительной квалификацией «Социолог: специалист по фундаментальным и прикладным социологическим исследованиям»)",
			RegularListID:        1973,
			TargetQuotaListIDs:   []int{2241},
			DedicatedQuotaListID: 510,
			SpecialQuotaListID:   1442,
			Capacities: core.Capacities{
				Regular:        23,
				TargetQuota:    0,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},

		// Классическая филология (древнегреческий и латинский языки; античная литература) (с дополнительной квалификацией «Учитель древнегреческого и латинского языков»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Классическая филология (древнегреческий и латинский языки; античная литература) (с дополнительной квалификацией «Учитель древнегреческого и латинского языков»)",
			RegularListID:        562,
			TargetQuotaListIDs:   []int{672},
			DedicatedQuotaListID: 1799,
			SpecialQuotaListID:   906,
			Capacities: core.Capacities{
				Regular:        7,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// Математика и компьютерные науки
		&spbsu.HttpHeadingSource{
			PrettyName:           "Математика и компьютерные науки",
			RegularListID:        1633,
			TargetQuotaListIDs:   []int{374, 577, 1367},
			DedicatedQuotaListID: 1024,
			SpecialQuotaListID:   270,
			Capacities: core.Capacities{
				Regular:        2,
				TargetQuota:    6,
				DedicatedQuota: 6,
				SpecialQuota:   4,
			},
		},

		// Инженерно-ориентированная физика
		&spbsu.HttpHeadingSource{
			PrettyName:           "Инженерно-ориентированная физика",
			RegularListID:        122,
			TargetQuotaListIDs:   []int{16, 2038},
			DedicatedQuotaListID: 1130,
			SpecialQuotaListID:   1386,
			Capacities: core.Capacities{
				Regular:        9,
				TargetQuota:    0,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},

		// Нефтегазовая геология
		&spbsu.HttpHeadingSource{
			PrettyName:           "Нефтегазовая геология",
			RegularListID:        1196,
			TargetQuotaListIDs:   []int{912},
			DedicatedQuotaListID: 401,
			SpecialQuotaListID:   51,
			Capacities: core.Capacities{
				Regular:        8,
				TargetQuota:    3,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},

		// Политология
		&spbsu.HttpHeadingSource{
			PrettyName:           "Политология",
			RegularListID:        1803,
			TargetQuotaListIDs:   []int{2104},
			DedicatedQuotaListID: 229,
			SpecialQuotaListID:   700,
			Capacities: core.Capacities{
				Regular:        28,
				TargetQuota:    0,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},

		// Организация туристской деятельности (с изучением тайского языка) (с дополнительной квалификацией «Руководитель гостиничных комплексов и иных средств размещения»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Организация туристской деятельности (с изучением тайского языка) (с дополнительной квалификацией «Руководитель гостиничных комплексов и иных средств размещения»)",
			RegularListID:        1736,
			TargetQuotaListIDs:   []int{},
			DedicatedQuotaListID: -1,
			SpecialQuotaListID:   1466,
			Capacities: core.Capacities{
				Regular:        4,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   1,
			},
		},

		// Сравнительно-историческое языкознание (английский язык)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Сравнительно-историческое языкознание (английский язык)",
			RegularListID:        289,
			TargetQuotaListIDs:   []int{1953},
			DedicatedQuotaListID: 240,
			SpecialQuotaListID:   1035,
			Capacities: core.Capacities{
				Regular:        2,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// Языки славянского мира и русский язык как иностранный (с дополнительной квалификацией «Учитель русского языка и литературы»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Языки славянского мира и русский язык как иностранный (с дополнительной квалификацией «Учитель русского языка и литературы»)",
			RegularListID:        2202,
			TargetQuotaListIDs:   []int{1518},
			DedicatedQuotaListID: 1493,
			SpecialQuotaListID:   1121,
			Capacities: core.Capacities{
				Regular:        2,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   0,
			},
		},

		// Дизайн среды (с дополнительной квалификацией «Художник-оформитель»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Дизайн среды (с дополнительной квалификацией «Художник-оформитель»)",
			RegularListID:        1764,
			TargetQuotaListIDs:   []int{},
			DedicatedQuotaListID: 193,
			SpecialQuotaListID:   214,
			Capacities: core.Capacities{
				Regular:        7,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// Фундаментальная математика (с дополнительной квалификацией «Педагог дополнительного образования»)
		&spbsu.HttpHeadingSource{
			PrettyName:           "Фундаментальная математика (с дополнительной квалификацией «Педагог дополнительного образования»)",
			RegularListID:        1756,
			TargetQuotaListIDs:   []int{1146},
			DedicatedQuotaListID: 1392,
			SpecialQuotaListID:   172,
			Capacities: core.Capacities{
				Regular:        15,
				TargetQuota:    0,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},

		// Программирование и информационные технологии
		&spbsu.HttpHeadingSource{
			PrettyName:           "Программирование и информационные технологии",
			RegularListID:        585,
			TargetQuotaListIDs:   []int{1926},
			DedicatedQuotaListID: 1935,
			SpecialQuotaListID:   2004,
			Capacities: core.Capacities{
				Regular:        21,
				TargetQuota:    0,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},

		// Экология и природопользование
		&spbsu.HttpHeadingSource{
			PrettyName:           "Экология и природопользование",
			RegularListID:        1695,
			TargetQuotaListIDs:   []int{326, 1814},
			DedicatedQuotaListID: 2020,
			SpecialQuotaListID:   1837,
			Capacities: core.Capacities{
				Regular:        14,
				TargetQuota:    1,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},

		// --- Headings with missing or zero regular capacity ---
		// // Классическая филология (древнегреческий и латинский языки
		// &spbsu.HttpHeadingSource{
		// 	PrettyName: "Классическая филология (древнегреческий и латинский языки",
		// 	RegularListID: -1,
		// 	TargetQuotaListIDs: []int{672}античная литература) (с дополнительной квалификацией «Учитель древнегреческого и латинского языков»)
		// 	DedicatedQuotaListID: -1,
		// 	SpecialQuotaListID: -1,
		// 	Capacities: core.Capacities{
		// 		Regular: 0,
		// 		TargetQuota: 1,
		// 		DedicatedQuota: 0,
		// 		SpecialQuota: 0,
		// 	},
		// },

		// // Юриспруденция (с углубленным изучением налогового права)
		// &spbsu.HttpHeadingSource{
		// 	PrettyName: "Юриспруденция (с углубленным изучением налогового права)",
		// 	RegularListID: 256,
		// 	TargetQuotaListIDs: []int{1927},
		// 	DedicatedQuotaListID: -1,
		// 	SpecialQuotaListID: -1,
		// 	Capacities: core.Capacities{
		// 		Regular: 0,
		// 		TargetQuota: 5,
		// 		DedicatedQuota: 0,
		// 		SpecialQuota: 0,
		// 	},
		// },

	}
}
