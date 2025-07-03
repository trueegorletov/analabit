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
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_Ave_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_Ave_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_Ave_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_Ave_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_Ave_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        39,
				TargetQuota:    6,
				DedicatedQuota: 9,
				SpecialQuota:   6,
			},
		},
		// Бизнес-информатика
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_BI_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_BI_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_BI_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_BI_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_BI_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        16,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Востоковедение
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_Vostok_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_Vostok_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_Vostok_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_Vostok_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_Vostok_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        21,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Дизайн
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_Design_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_Design_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_Design_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_Design_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_Design_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        16,
				TargetQuota:    2,
				DedicatedQuota: 0,
				SpecialQuota:   2,
			},
		},
		// История
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_ISTR_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_ISTR_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_ISTR_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_ISTR_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_ISTR_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        19,
				TargetQuota:    3,
				DedicatedQuota: 5,
				SpecialQuota:   3,
			},
		},
		// Компьютерные технологии, системы и сети
		&hse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_K_KNT_SysNet_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        7,
				TargetQuota:    2,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},
		// Медиакоммуникации
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_Media_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_Media_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_Media_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_Media_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_Media_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        15,
				TargetQuota:    2,
				DedicatedQuota: 1,
				SpecialQuota:   2,
			},
		},
		// Политология и мировая политика
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_PMP_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_PMP_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_PMP_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_PMP_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_PMP_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        11,
				TargetQuota:    3,
				DedicatedQuota: 8,
				SpecialQuota:   3,
			},
		},
		// Прикладная математика и  информатика
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_AMI_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_AMI_O.xlsx",
			DQListURL: "/mirror/pubs/share/947713972",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_AMI_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_AMI_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        21,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Прикладной анализ данных и искусственный интеллект
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_PADII_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_PADII_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_PADII_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_PADII_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_PADII_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        21,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Социология и социальная информатика
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_Socci_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_Socci_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_Socci_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_Socci_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_Socci_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        37,
				TargetQuota:    5,
				DedicatedQuota: 3,
				SpecialQuota:   5,
			},
		},
		// Управление бизнесом
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_BBA_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_BBA_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_BBA_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_BBA_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_BBA_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        42,
				TargetQuota:    6,
				DedicatedQuota: 6,
				SpecialQuota:   6,
			},
		},
		// Управление  и аналитика в государственном секторе
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_GMU_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_GMU_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_GMU_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_GMU_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_GMU_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        30,
				TargetQuota:    5,
				DedicatedQuota: 10,
				SpecialQuota:   5,
			},
		},
		// Физика
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_Physics_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_Physics_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_Physics_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_Physics_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_Physics_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        15,
				TargetQuota:    2,
				DedicatedQuota: 1,
				SpecialQuota:   2,
			},
		},
		// Филология
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_Philology_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_Philology_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_Philology_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_Philology_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_Philology_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        16,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Юриспруденция
		&hse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_LAW_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_LAW_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_LAW_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_LAW_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_LAW_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        58,
				TargetQuota:    9,
				DedicatedQuota: 14,
				SpecialQuota:   9,
			},
		},

		// TODO The following SPB headings do not have list URLs determined:

		// Архитектура **
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
		// Многопрофильный конкурс "Международный бакалавриат по бизнесу и экономике"
		&hse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "",
			HeadingCapacities: core.Capacities{
				Regular:        39,
				TargetQuota:    6,
				DedicatedQuota: 9,
				SpecialQuota:   6,
			},
		},
		// Программирование и инжиниринг компьютерных игр
		&hse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "",
			HeadingCapacities: core.Capacities{
				Regular:        6,
				TargetQuota:    2,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},
		// Тексты, языки и цифровые инструменты
		&hse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "",
			HeadingCapacities: core.Capacities{
				Regular:        8,
				TargetQuota:    2,
				DedicatedQuota: 1,
				SpecialQuota:   2,
			},
		},
	}
}
