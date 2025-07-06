package hse

import (
	"analabit/core"
	"analabit/core/source"
	"analabit/core/source/hse"
)

const (
	permCode = "hse_perm"
	permName = "ВШЭ (Пермь)"
)

func permSourcesList() []source.HeadingSource {
	return []source.HeadingSource{
		// Иностранные языки и межкультурная коммуникация в бизнесе
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/perm/Bachelors/BD_perm_IYAMK_O.xlsx",
			Capacities: core.Capacities{
				Regular:        8,
				TargetQuota:    2,
				DedicatedQuota: 3,
				SpecialQuota:   2,
			},
		},
		// Международный бакалавриат по бизнесу и экономике
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/perm/Bachelors/BD_perm_MBBE_O.xlsx",
			Capacities: core.Capacities{
				Regular:        45,
				TargetQuota:    7,
				DedicatedQuota: 11,
				SpecialQuota:   7,
			},
		},
		// Менеджмент в креативных индустриях
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/perm/Bachelors/BD_perm_MCI_O.xlsx",
			Capacities: core.Capacities{
				Regular:        14,
				TargetQuota:    2,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},
		// Разработка информационных систем для бизнеса (Бизнес-информатика)
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/perm/Bachelors/BD_perm_IsystemsB_O.xlsx",
			Capacities: core.Capacities{
				Regular:        30,
				TargetQuota:    5,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
		},
		// Разработка информационных систем для бизнеса (Программная инженерия)
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/perm/Bachelors/BD_perm_IsystemsP_O.xlsx",
			Capacities: core.Capacities{
				Regular:        30,
				TargetQuota:    5,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
		},
		// Юриспруденция
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/perm/Bachelors/BD_perm_LAW_O.xlsx",
			Capacities: core.Capacities{
				Regular:        13,
				TargetQuota:    2,
				DedicatedQuota: 3,
				SpecialQuota:   2,
			},
		},

		// ===
		// FILTERED OUT
		// ===
		// Дизайн (KCP = 0)
		// Программные системы и автоматизация процессов разработки (KCP = 0)
		// Управление бизнесом (KCP = 0)

		// ===
		// FOUND ANOMALIES
		// ===
		// Program found in capacities but not in lists: -
	}
}
