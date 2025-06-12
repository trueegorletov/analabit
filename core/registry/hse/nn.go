package hse

import (
	"analabit/core"
	"analabit/core/source"
	"analabit/core/source/hse"
)

const (
	nnCode = "hse_nn"
	nnName = "Высшая школа экономики (НН)"
)

func nnSourcesList() []source.HeadingSource {
	return []source.HeadingSource{
		// Компьютерные науки и технологии
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_OM_nn_B_KNT_BI_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_KCP_nn_B_KNT_BI_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_SK_nn_B_KNT_BI_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_OP_nn_B_KNT_BI_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_BVI_nn_B_KNT_BI_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        44,
				TargetQuota:    7,
				DedicatedQuota: 7,
				SpecialQuota:   7,
			},
		},
		// Международный бакалавриат по бизнесу и экономике
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_OM_nn_B_MBBE_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_KCP_nn_B_MBBE_M_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_SK_nn_B_MBBE_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_OP_nn_B_MBBE_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_BVI_nn_B_MBBE_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        78,
				TargetQuota:    18,
				DedicatedQuota: 12,
				SpecialQuota:   12,
			},
		},
		// Технологии искусственного и дополненного интеллекта
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_OM_nn_B_Ait_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_KCP_nn_B_Ait_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_SK_nn_B_Ait_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_OP_nn_B_Ait_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_BVI_nn_B_Ait_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        6,
				TargetQuota:    2,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},
		// Филология
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_OM_nn_B_Philology_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_KCP_nn_B_Philology_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_SK_nn_B_Philology_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_OP_nn_B_Philology_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_BVI_nn_B_Philology_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        16,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Фундаментальная и прикладная
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_OM_nn_B_Ling_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_KCP_nn_B_Ling_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_SK_nn_B_Ling_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_OP_nn_B_Ling_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_BVI_nn_B_Ling_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        22,
				TargetQuota:    2,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Фундаментальная и прикладная математика
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_OM_nn_B_Math_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_KCP_nn_B_Math_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_SK_nn_B_Math_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_OP_nn_B_Math_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_BVI_nn_B_Math_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        26,
				TargetQuota:    6,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},
		// Юриспруденция
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_OM_nn_B_LAW_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_KCP_nn_B_LAW_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_SK_nn_B_LAW_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_OP_nn_B_LAW_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_BVI_nn_B_LAW_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        35,
				TargetQuota:    8,
				DedicatedQuota: 6,
				SpecialQuota:   6,
			},
		},

		// TODO The following NN headings do not have capacities determined:

		// Дизайн
		&hse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_BVI_nn_K_Design_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        0,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   0,
			},
		},
		// Иностранные языки и межкультурная бизнес-коммуникация
		&hse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_BVI_nn_K_IYAMK_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        0,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   0,
			},
		},
		// Программная инженерия
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_OM_nn_B_KNT_SE_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_KCP_nn_B_KNT_SE_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_SK_nn_B_KNT_SE_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_OP_nn_B_KNT_SE_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_BVI_nn_B_KNT_SE_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        0,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   0,
			},
		},
		// Цифровой маркетинг
		&hse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/nn/Bachelors/KS_BVI_nn_K_DM_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        0,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   0,
			},
		},

		// TODO The following NN headings do not have list URLs determined:

		// Экономика и бизнес
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
