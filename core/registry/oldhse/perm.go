package oldhse

import (
	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/source/oldhse"
)

const (
	permCode = "hse_perm"
	permName = "ВШЭ (Пермь)"
)

func permSourcesList() []source.HeadingSource {
	return []source.HeadingSource{
		// Иностранные языки и межкультурная коммуникация в бизнесе
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_OM_perm_B_IYAMK_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_KCP_perm_B_IYAMK_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_SK_perm_B_IYAMK_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_OP_perm_B_IYAMK_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_BVI_perm_B_IYAMK_O.xlsx",
			Capacities: core.Capacities{
				Regular:        8,
				TargetQuota:    3,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},
		// Международный бакалавриат по бизнесу и экономике
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_OM_perm_B_MBBE_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_KCP_perm_B_MBBE_M_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_SK_perm_B_MBBE_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_OP_perm_B_MBBE_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_BVI_perm_B_MBBE_O.xlsx",
			Capacities: core.Capacities{
				Regular:        45,
				TargetQuota:    11,
				DedicatedQuota: 7,
				SpecialQuota:   7,
			},
		},
		// Разработка информационных систем для бизнеса
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_OM_perm_B_IsystemsB_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_KCP_perm_B_IsystemsB_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_SK_perm_B_IsystemsB_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_OP_perm_B_IsystemsB_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_BVI_perm_B_IsystemsB_O.xlsx",
			Capacities: core.Capacities{
				Regular:        30,
				TargetQuota:    5,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
		},
		// Юриспруденция
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_OM_perm_B_LAW_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_KCP_perm_B_LAW_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_SK_perm_B_LAW_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_OP_perm_B_LAW_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/perm/Bachelors/KS_BVI_perm_B_LAW_O.xlsx",
			Capacities: core.Capacities{
				Regular:        13,
				TargetQuota:    3,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},

		// TODO The following Perm headings do not have list URLs determined:

		// Менеджмент в креативных индустриях
		&oldhse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "",
			Capacities: core.Capacities{
				Regular:        14,
				TargetQuota:    2,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},
		// Программные системы и автоматизация процессов разработки (онлайн)
		&oldhse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "",
			Capacities: core.Capacities{
				Regular:        0,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   0,
			},
		},
		// Управление бизнесом (онлайн)
		&oldhse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "",
			Capacities: core.Capacities{
				Regular:        0,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   0,
			},
		},
	}
}
