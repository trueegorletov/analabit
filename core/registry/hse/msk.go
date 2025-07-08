package hse

import (
	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/source/hse"
)

const (
	mskCode = "hse_msk"
	mskName = "ВШЭ (Москва)"
)

func mskSourcesList() []source.HeadingSource {
	return []source.HeadingSource{
		// Античность (История)
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_Ant_Ist_O.xlsx",
			Capacities: core.Capacities{
				Regular:        6,
				TargetQuota:    2,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},
		// Античность (Филология)
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_Ant_Fil_O.xlsx",
			Capacities: core.Capacities{
				Regular:        6,
				TargetQuota:    2,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},
		// Бизнес-информатика
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_BI_O.xlsx",
			Capacities: core.Capacities{
				Regular:        91,
				TargetQuota:    13,
				DedicatedQuota: 13,
				SpecialQuota:   13,
			},
		},
		// Востоковедение
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_Vostok_O.xlsx",
			Capacities: core.Capacities{
				Regular:        44,
				TargetQuota:    2,
				DedicatedQuota: 6,
				SpecialQuota:   6,
			},
		},
		// География глобальных изменений и геоинформационные технологии
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_GGIGT_O.xlsx",
			Capacities: core.Capacities{
				Regular:        30,
				TargetQuota:    2,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},
		// Городское планирование
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_GorPlan_O.xlsx",
			Capacities: core.Capacities{
				Regular:        28,
				TargetQuota:    7,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
		},
		// Государственное и муниципальное управление
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_GMU_O.xlsx",
			Capacities: core.Capacities{
				Regular:        44,
				TargetQuota:    15,
				DedicatedQuota: 8,
				SpecialQuota:   8,
			},
		},
		// Дизайн
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_Design_O.xlsx",
			Capacities: core.Capacities{
				Regular:        51,
				TargetQuota:    0,
				DedicatedQuota: 7,
				SpecialQuota:   7,
			},
		},
		// Журналистика
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_ZHur_O.xlsx",
			Capacities: core.Capacities{
				Regular:        42,
				TargetQuota:    6,
				DedicatedQuota: 6,
				SpecialQuota:   6,
			},
		},
		// Иностранные языки и межкультурная коммуникация
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_IYAMK_O.xlsx",
			Capacities: core.Capacities{
				Regular:        23,
				TargetQuota:    4,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},
		// Инфокоммуникационные технологии и системы связи
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_ITSS_O.xlsx",
			Capacities: core.Capacities{
				Regular:        55,
				TargetQuota:    8,
				DedicatedQuota: 8,
				SpecialQuota:   8,
			},
		},
		// Информатика и вычислительная техника
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_IVT_O.xlsx",
			Capacities: core.Capacities{
				Regular:        80,
				TargetQuota:    19,
				DedicatedQuota: 13,
				SpecialQuota:   13,
			},
		},
		// Информационная безопасность
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_IB_O.xlsx",
			Capacities: core.Capacities{
				Regular:        33,
				TargetQuota:    31,
				DedicatedQuota: 8,
				SpecialQuota:   8,
			},
		},
		// История
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_ISTR_O.xlsx",
			Capacities: core.Capacities{
				Regular:        46,
				TargetQuota:    10,
				DedicatedQuota: 7,
				SpecialQuota:   7,
			},
		},
		// История искусств
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_Isk_O.xlsx",
			Capacities: core.Capacities{
				Regular:        24,
				TargetQuota:    0,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Клеточная и молекулярная биотехнология
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_CMB_O.xlsx",
			Capacities: core.Capacities{
				Regular:        18,
				TargetQuota:    1,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Когнитивная нейробиология
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_KogNeir_O.xlsx",
			Capacities: core.Capacities{
				Regular:        10,
				TargetQuota:    1,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},
		// Компьютерная безопасность
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_KB_O.xlsx",
			Capacities: core.Capacities{
				Regular:        17,
				TargetQuota:    18,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
		},
		// Культурология
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_Cultural_O.xlsx",
			Capacities: core.Capacities{
				Regular:        32,
				TargetQuota:    0,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},
		// Маркетинг и рыночная аналитика
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_Marketing_O.xlsx",
			Capacities: core.Capacities{
				Regular:        14,
				TargetQuota:    2,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},
		// Математика
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_Math_O.xlsx",
			Capacities: core.Capacities{
				Regular:        42,
				TargetQuota:    6,
				DedicatedQuota: 6,
				SpecialQuota:   6,
			},
		},
		// Медиакоммуникации
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_Media_O.xlsx",
			Capacities: core.Capacities{
				Regular:        31,
				TargetQuota:    1,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},
		// Международные отношения
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_MO_O.xlsx",
			Capacities: core.Capacities{
				Regular:        25,
				TargetQuota:    10,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
		},
		// Мировая экономика
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_WE_O.xlsx",
			Capacities: core.Capacities{
				Regular:        43,
				TargetQuota:    5,
				DedicatedQuota: 6,
				SpecialQuota:   6,
			},
		},
		// Монголия и Тибет
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_Mongol_O.xlsx",
			Capacities: core.Capacities{
				Regular:        7,
				TargetQuota:    1,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},
		// Политология
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_Political_O.xlsx",
			Capacities: core.Capacities{
				Regular:        28,
				TargetQuota:    4,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},
		// Прикладная математика
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_AM_O.xlsx",
			Capacities: core.Capacities{
				Regular:        52,
				TargetQuota:    12,
				DedicatedQuota: 8,
				SpecialQuota:   8,
			},
		},
		// Прикладная математика и информатика
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_AMI_O.xlsx",
			Capacities: core.Capacities{
				Regular:        126,
				TargetQuota:    18,
				DedicatedQuota: 18,
				SpecialQuota:   18,
			},
		},
		// Программная инженерия
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_SE_O.xlsx",
			Capacities: core.Capacities{
				Regular:        104,
				TargetQuota:    24,
				DedicatedQuota: 16,
				SpecialQuota:   16,
			},
		},
		// Психология
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_Psy_O.xlsx",
			Capacities: core.Capacities{
				Regular:        59,
				TargetQuota:    5,
				DedicatedQuota: 8,
				SpecialQuota:   8,
			},
		},
		// Реклама и связи с общественностью
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_AD_O.xlsx",
			Capacities: core.Capacities{
				Regular:        2,
				TargetQuota:    1,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},
		// Реклама и связи с общественностью (Медиакоммуникации)
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_ADM_O.xlsx",
			Capacities: core.Capacities{
				Regular:        2,
				TargetQuota:    1,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},
		// Совместная программа по экономике НИУ ВШЭ и РЭШ
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_RESH_O.xlsx",
			Capacities: core.Capacities{
				Regular:        39,
				TargetQuota:    9,
				DedicatedQuota: 6,
				SpecialQuota:   6,
			},
		},
		// Совместный бакалавриат НИУ ВШЭ и Центра педагогического мастерства
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_CPM_O.xlsx",
			Capacities: core.Capacities{
				Regular:        16,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Социология
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_Soc_O.xlsx",
			Capacities: core.Capacities{
				Regular:        85,
				TargetQuota:    6,
				DedicatedQuota: 12,
				SpecialQuota:   12,
			},
		},
		// Стратегия и продюсирование в коммуникациях
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_Producer_O.xlsx",
			Capacities: core.Capacities{
				Regular:        7,
				TargetQuota:    1,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},
		// Технологии анализа данных в бизнесе
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_TADinBis_O.xlsx",
			Capacities: core.Capacities{
				Regular:        21,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Турция и тюркский мир
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_Turckey_O.xlsx",
			Capacities: core.Capacities{
				Regular:        7,
				TargetQuota:    1,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},
		// Управление бизнесом
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_BBA_O.xlsx",
			Capacities: core.Capacities{
				Regular:        21,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Физика
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_Physics_O.xlsx",
			Capacities: core.Capacities{
				Regular:        37,
				TargetQuota:    3,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
		},
		// Филология
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_Philology_O.xlsx",
			Capacities: core.Capacities{
				Regular:        30,
				TargetQuota:    5,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
		},
		// Философия
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_Phil_O.xlsx",
			Capacities: core.Capacities{
				Regular:        30,
				TargetQuota:    2,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},
		// Фундаментальная и компьютерная лингвистика
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_FKL_O.xlsx",
			Capacities: core.Capacities{
				Regular:        33,
				TargetQuota:    2,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
		},
		// Химия
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_Chem_O.xlsx",
			Capacities: core.Capacities{
				Regular:        20,
				TargetQuota:    7,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},
		// Химия новых материалов
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_ChimNewMat_O.xlsx",
			Capacities: core.Capacities{
				Regular:        9,
				TargetQuota:    2,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},
		// Экономика
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_Economy_O.xlsx",
			Capacities: core.Capacities{
				Regular:        65,
				TargetQuota:    15,
				DedicatedQuota: 10,
				SpecialQuota:   10,
			},
		},
		// Экономика и анализ данных
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_EcAl_O.xlsx",
			Capacities: core.Capacities{
				Regular:        28,
				TargetQuota:    7,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
		},
		// Экономика и статистика
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_Stat_O.xlsx",
			Capacities: core.Capacities{
				Regular:        22,
				TargetQuota:    5,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},
		// Эфиопия и арабский мир
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/moscow/Bachelors/BD_moscow_Efiop_O.xlsx",
			Capacities: core.Capacities{
				Regular:        7,
				TargetQuota:    1,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// ===
		// FILTERED OUT
		// ===
		// Глобальные цифровые коммуникации (Медиакоммуникации) (KCP = 0)
		// Глобальные цифровые коммуникации (Реклама и связи с общественностью) (KCP = 0)
		// Дизайн и разработка информационных продуктов (KCP = 0)
		// Кинопроизводство (KCP = 0)
		// Компьютерные науки и анализ данных (KCP = 0)
		// Международная программа по экономике и финансам (KCP = 0)
		// Международный бизнес (KCP = 0)
		// Мода (KCP = 0)
		// Прикладной анализ данных (KCP = 0)
		// Разработка игр и цифровых продуктов (KCP = 0)
		// Разработка игр и цифровых продуктов (Дизайн) (KCP = 0)
		// Современное искусство (KCP = 0)
		// Управление в креативных индустриях (KCP = 0)
		// Управление цифровым продуктом (KCP = 0)
		// Экономический анализ (KCP = 0)
		// Юриспруденция (KCP = 0)
		// Юриспруденция: правовое регулирование бизнеса (очно-заочная форма обучения) (KCP = 0)
		// Юриспруденция: цифровой юрист (KCP = 0)

		// ===
		// FOUND ANOMALIES
		// ===
		// Program found in lists but not in capacities: Актер
		// Program found in lists but not in capacities: Международная программа «Международные отношения и глобальные исследования» (Международные отношения)
		// Program found in lists but not in capacities: Международная программа «Международные отношения и глобальные исследования» (Публичная политика и социальные науки)
		// Program found in lists but not in capacities: Программа двух дипломов НИУ ВШЭ и Университета Кёнхи "Экономика и политика Азии" (Зарубежное регионоведение)
		// Program found in lists but not in capacities: Программа двух дипломов НИУ ВШЭ и Университета Кёнхи "Экономика и политика Азии" (Публичная политика и социальные науки)
		// Program found in capacities but not in lists: программа двух дипломов и университета кёнхи "экономика и политика в азии"/khu-hse double degree program economics and politics in asia
		// Program found in capacities but not in lists: актёр
		// Program found in capacities but not in lists: международная программа "международные отношения и глобальные исследования"/ international program "international relations and global studies "

	}
}
