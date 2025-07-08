package oldhse

import (
	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/source/oldhse"
)

const (
	spbCode = "hse_spb"
	spbName = "ВШЭ (СПб)"
)

func spbSourcesList() []source.HeadingSource {
	return []source.HeadingSource{
		// Аналитика в экономике
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_Ave_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_Ave_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_Ave_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_Ave_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_Ave_O.xlsx",
			Capacities: core.Capacities{
				Regular:        39,
				TargetQuota:    6,
				DedicatedQuota: 9,
				SpecialQuota:   6,
			},
		},
		// Бизнес-информатика
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_BI_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_BI_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_BI_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_BI_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_BI_O.xlsx",
			Capacities: core.Capacities{
				Regular:        16,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Востоковедение
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_Vostok_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_Vostok_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_Vostok_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_Vostok_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_Vostok_O.xlsx",
			Capacities: core.Capacities{
				Regular:        21,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Дизайн
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_Design_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_Design_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_Design_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_Design_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_Design_O.xlsx",
			Capacities: core.Capacities{
				Regular:        16,
				TargetQuota:    2,
				DedicatedQuota: 0,
				SpecialQuota:   2,
			},
		},
		// История
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_ISTR_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_ISTR_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_ISTR_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_ISTR_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_ISTR_O.xlsx",
			Capacities: core.Capacities{
				Regular:        19,
				TargetQuota:    3,
				DedicatedQuota: 5,
				SpecialQuota:   3,
			},
		},
		// Компьютерные технологии, системы и сети
		&oldhse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_K_KNT_SysNet_O.xlsx",
			Capacities: core.Capacities{
				Regular:        7,
				TargetQuota:    2,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},
		// Медиакоммуникации
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_Media_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_Media_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_Media_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_Media_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_Media_O.xlsx",
			Capacities: core.Capacities{
				Regular:        15,
				TargetQuota:    2,
				DedicatedQuota: 1,
				SpecialQuota:   2,
			},
		},
		// Политология и мировая политика
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_PMP_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_PMP_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_PMP_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_PMP_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_PMP_O.xlsx",
			Capacities: core.Capacities{
				Regular:        11,
				TargetQuota:    3,
				DedicatedQuota: 8,
				SpecialQuota:   3,
			},
		},
		// Прикладная математика и  информатика
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_AMI_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_AMI_O.xlsx",
			DQListURL: "/mirror/pubs/share/947713972",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_AMI_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_AMI_O.xlsx",
			Capacities: core.Capacities{
				Regular:        21,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Прикладной анализ данных и искусственный интеллект
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_PADII_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_PADII_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_PADII_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_PADII_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_PADII_O.xlsx",
			Capacities: core.Capacities{
				Regular:        21,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Социология и социальная информатика
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_Socci_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_Socci_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_Socci_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_Socci_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_Socci_O.xlsx",
			Capacities: core.Capacities{
				Regular:        37,
				TargetQuota:    5,
				DedicatedQuota: 3,
				SpecialQuota:   5,
			},
		},
		// Управление бизнесом
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_BBA_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_BBA_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_BBA_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_BBA_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_BBA_O.xlsx",
			Capacities: core.Capacities{
				Regular:        42,
				TargetQuota:    6,
				DedicatedQuota: 6,
				SpecialQuota:   6,
			},
		},
		// Управление  и аналитика в государственном секторе
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_GMU_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_GMU_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_GMU_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_GMU_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_GMU_O.xlsx",
			Capacities: core.Capacities{
				Regular:        30,
				TargetQuota:    5,
				DedicatedQuota: 10,
				SpecialQuota:   5,
			},
		},
		// Физика
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_Physics_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_Physics_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_Physics_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_Physics_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_Physics_O.xlsx",
			Capacities: core.Capacities{
				Regular:        15,
				TargetQuota:    2,
				DedicatedQuota: 1,
				SpecialQuota:   2,
			},
		},
		// Филология
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_Philology_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_Philology_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_Philology_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_Philology_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_Philology_O.xlsx",
			Capacities: core.Capacities{
				Regular:        16,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Юриспруденция
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OM_spb_B_LAW_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_KCP_spb_B_LAW_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_SK_spb_B_LAW_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_OP_spb_B_LAW_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/spb/Bachelors/KS_BVI_spb_B_LAW_O.xlsx",
			Capacities: core.Capacities{
				Regular:        58,
				TargetQuota:    9,
				DedicatedQuota: 14,
				SpecialQuota:   9,
			},
		},

		// TODO The following SPB headings do not have list URLs determined:

		// Архитектура **
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
		// Многопрофильный конкурс "Международный бакалавриат по бизнесу и экономике"
		&oldhse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "",
			Capacities: core.Capacities{
				Regular:        39,
				TargetQuota:    6,
				DedicatedQuota: 9,
				SpecialQuota:   6,
			},
		},
		// Программирование и инжиниринг компьютерных игр
		&oldhse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "",
			Capacities: core.Capacities{
				Regular:        6,
				TargetQuota:    2,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},
		// Тексты, языки и цифровые инструменты
		&oldhse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "",
			Capacities: core.Capacities{
				Regular:        8,
				TargetQuota:    2,
				DedicatedQuota: 1,
				SpecialQuota:   2,
			},
		},
	}
}
