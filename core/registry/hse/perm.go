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
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_OM_perm_B_IYAMK_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_KCP_perm_B_IYAMK_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_SK_perm_B_IYAMK_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_OP_perm_B_IYAMK_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_BVI_perm_B_IYAMK_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        8,
				TargetQuota:    3,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},
		// Международный бакалавриат по бизнесу и экономике
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_OM_perm_B_MBBE_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_KCP_perm_B_MBBE_M_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_SK_perm_B_MBBE_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_OP_perm_B_MBBE_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_BVI_perm_B_MBBE_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        45,
				TargetQuota:    11,
				DedicatedQuota: 7,
				SpecialQuota:   7,
			},
		},
		// Разработка информационных систем для бизнеса
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_OM_perm_B_IsystemsB_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_KCP_perm_B_IsystemsB_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_SK_perm_B_IsystemsB_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_OP_perm_B_IsystemsB_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_BVI_perm_B_IsystemsB_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        30,
				TargetQuota:    5,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
		},
		// Юриспруденция
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_OM_perm_B_LAW_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_KCP_perm_B_LAW_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_SK_perm_B_LAW_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_OP_perm_B_LAW_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_BVI_perm_B_LAW_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        13,
				TargetQuota:    3,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},

		// TODO The following Perm headings do not have list URLs determined:

		// Менеджмент в креативных индустриях
		&hse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "",
			HeadingCapacities: core.Capacities{
				Regular:        14,
				TargetQuota:    2,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},
		// Программные системы и автоматизация процессов разработки (онлайн)
		&hse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "",
			HeadingCapacities: core.Capacities{
				Regular:        0,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   0,
			},
		},
		// Управление бизнесом (онлайн)
		&hse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "",
			HeadingCapacities: core.Capacities{
				Regular:        0,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   0,
			},
		},
	}
}
