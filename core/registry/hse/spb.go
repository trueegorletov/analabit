package hse

import (
	"analabit/core"
	"analabit/core/source"
	"analabit/core/source/hse"
)

const (
	spbCode = "hse_spb"
	spbName = "ВШЭ (СПб)"
)

func spbSourcesList() []source.HeadingSource {
	return []source.HeadingSource{
		// Аналитика в экономике
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/spb/Bachelors/BD_spb_Ave_O.xlsx",
			Capacities: core.Capacities{
				Regular:        39,
				TargetQuota:    6,
				DedicatedQuota: 9,
				SpecialQuota:   6,
			},
		},
		// Бизнес-информатика
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/spb/Bachelors/BD_spb_BI_O.xlsx",
			Capacities: core.Capacities{
				Regular:        16,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Востоковедение
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/spb/Bachelors/BD_spb_Vostok_O.xlsx",
			Capacities: core.Capacities{
				Regular:        21,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Дизайн
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/spb/Bachelors/BD_spb_Design_O.xlsx",
			Capacities: core.Capacities{
				Regular:        16,
				TargetQuota:    2,
				DedicatedQuota: 0,
				SpecialQuota:   2,
			},
		},
		// История
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/spb/Bachelors/BD_spb_ISTR_O.xlsx",
			Capacities: core.Capacities{
				Regular:        19,
				TargetQuota:    3,
				DedicatedQuota: 5,
				SpecialQuota:   3,
			},
		},
		// Компьютерные технологии, системы и сети
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/spb/Bachelors/BD_spb_KNT_SysNet_O.xlsx",
			Capacities: core.Capacities{
				Regular:        7,
				TargetQuota:    2,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},
		// Медиакоммуникации
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/spb/Bachelors/BD_spb_Media_O.xlsx",
			Capacities: core.Capacities{
				Regular:        15,
				TargetQuota:    2,
				DedicatedQuota: 1,
				SpecialQuota:   2,
			},
		},
		// Политология и мировая политика (Международные отношения)
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/spb/Bachelors/BD_spb_PMP_O.xlsx",
			Capacities: core.Capacities{
				Regular:        11,
				TargetQuota:    3,
				DedicatedQuota: 8,
				SpecialQuota:   3,
			},
		},
		// Политология и мировая политика (Политология)
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/spb/Bachelors/BD_spb_P_WP_O.xlsx",
			Capacities: core.Capacities{
				Regular:        11,
				TargetQuota:    3,
				DedicatedQuota: 8,
				SpecialQuota:   3,
			},
		},
		// Прикладная математика и информатика
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/spb/Bachelors/BD_spb_AMI_O.xlsx",
			Capacities: core.Capacities{
				Regular:        21,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Прикладной анализ данных и искусственный интеллект
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/spb/Bachelors/BD_spb_PADII_O.xlsx",
			Capacities: core.Capacities{
				Regular:        21,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Программирование и инжиниринг компьютерных игр
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/spb/Bachelors/BD_spb_P_ICG_O.xlsx",
			Capacities: core.Capacities{
				Regular:        6,
				TargetQuota:    2,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},
		// Социология и социальная информатика
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/spb/Bachelors/BD_spb_Socci_O.xlsx",
			Capacities: core.Capacities{
				Regular:        37,
				TargetQuota:    5,
				DedicatedQuota: 3,
				SpecialQuota:   5,
			},
		},
		// Тексты, языки и цифровые инструменты
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/spb/Bachelors/BD_spb_TLDI_O.xlsx",
			Capacities: core.Capacities{
				Regular:        8,
				TargetQuota:    2,
				DedicatedQuota: 1,
				SpecialQuota:   2,
			},
		},
		// Управление бизнесом
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/spb/Bachelors/BD_spb_BBA_O.xlsx",
			Capacities: core.Capacities{
				Regular:        42,
				TargetQuota:    6,
				DedicatedQuota: 6,
				SpecialQuota:   6,
			},
		},
		// Управление и аналитика в государственном секторе
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/spb/Bachelors/BD_spb_GMU_O.xlsx",
			Capacities: core.Capacities{
				Regular:        30,
				TargetQuota:    5,
				DedicatedQuota: 10,
				SpecialQuota:   5,
			},
		},
		// Физика
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/spb/Bachelors/BD_spb_Physics_O.xlsx",
			Capacities: core.Capacities{
				Regular:        15,
				TargetQuota:    2,
				DedicatedQuota: 1,
				SpecialQuota:   2,
			},
		},
		// Филология
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/spb/Bachelors/BD_spb_Philology_O.xlsx",
			Capacities: core.Capacities{
				Regular:        16,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Юриспруденция
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/spb/Bachelors/BD_spb_LAW_O.xlsx",
			Capacities: core.Capacities{
				Regular:        58,
				TargetQuota:    9,
				DedicatedQuota: 14,
				SpecialQuota:   9,
			},
		},

		// ===
		// FILTERED OUT
		// ===
		// Архитектура (KCP = 0)

		// ===
		// FOUND ANOMALIES
		// ===
		// Program found in lists but not in capacities: Бакалаврская программа «Международный бакалавриат по бизнесу и экономике»
		// Program found in capacities but not in lists: многопрофильный конкурс "международный бакалавриат по бизнесу и экономике"
	}
}
