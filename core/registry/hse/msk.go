package hse

import (
	"analabit/core"
	"analabit/core/source"
	"analabit/core/source/oldhse"
)

const (
	mskCode = "hse_msk"
	mskName = "ВШЭ (Москва)"
)

func mskSourcesList() []source.HeadingSource {
	return []source.HeadingSource{
		// Античность
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_Ant_Ist_O.xlsx",
			TQListURL: "",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_Ant_Ist_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_Ant_Ist_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_Ant_Ist_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        6,
				TargetQuota:    2,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},
		// Востоковедение
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_Vostok_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_KCP_moscow_B_Vostok_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_Vostok_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_Vostok_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_Vostok_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        44,
				TargetQuota:    2,
				DedicatedQuota: 6,
				SpecialQuota:   6,
			},
		},
		// География  глобальных изменений и геоинформационные технологии
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_GGIGT_O.xlsx",
			TQListURL: "",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_GGIGT_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_GGIGT_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_GGIGT_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        30,
				TargetQuota:    2,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},
		// Городское планирование
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_GorPlan_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_KCP_moscow_B_GorPlan_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_GorPlan_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_GorPlan_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_GorPlan_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        28,
				TargetQuota:    7,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
		},
		// Государственное и муниципальное управление
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_GMU_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_KCP_moscow_B_GMU_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_GMU_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_GMU_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_GMU_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        44,
				TargetQuota:    15,
				DedicatedQuota: 8,
				SpecialQuota:   8,
			},
		},
		// Дизайн
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_Design_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_KCP_moscow_B_Design_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_Design_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_Design_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_Design_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        51,
				TargetQuota:    0,
				DedicatedQuota: 7,
				SpecialQuota:   7,
			},
		},
		// Журналистика
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_ZHur_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_KCP_moscow_B_ZHur_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_ZHur_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_ZHur_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_ZHur_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        42,
				TargetQuota:    6,
				DedicatedQuota: 6,
				SpecialQuota:   6,
			},
		},
		// Иностранные языки и межкультурная коммуникация
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_IYAMK_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_KCP_moscow_B_IYAMK_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_IYAMK_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_IYAMK_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_IYAMK_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        23,
				TargetQuota:    4,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},
		// Инфокоммуникационные технологии и системы связи
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_ITSS_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_KCP_moscow_B_ITSS_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_ITSS_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_ITSS_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_ITSS_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        55,
				TargetQuota:    8,
				DedicatedQuota: 8,
				SpecialQuota:   8,
			},
		},
		// Информатика и вычислительная техника
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_IVT_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_KCP_moscow_B_IVT_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_IVT_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_IVT_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_IVT_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        80,
				TargetQuota:    19,
				DedicatedQuota: 13,
				SpecialQuota:   13,
			},
		},
		// Информационная безопасность
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_IB_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_KCP_moscow_B_IB_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_IB_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_IB_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_IB_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        33,
				TargetQuota:    31,
				DedicatedQuota: 8,
				SpecialQuota:   8,
			},
		},
		// История
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_ISTR_O.xlsx",
			TQListURL: "",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_ISTR_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_ISTR_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_ISTR_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        46,
				TargetQuota:    10,
				DedicatedQuota: 7,
				SpecialQuota:   7,
			},
		},
		// История искусств
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_Isk_O.xlsx",
			TQListURL: "",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_Isk_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_Isk_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_Isk_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        24,
				TargetQuota:    0,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Клеточная и молекулярная биотехнология
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_CMB_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_KCP_moscow_B_CMB_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_CMB_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_CMB_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_CMB_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        18,
				TargetQuota:    1,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Когнитивная нейробиология
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_KogNeir_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_KCP_moscow_B_KogNeir_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_KogNeir_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_KogNeir_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_KogNeir_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        10,
				TargetQuota:    1,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},
		// Компьютерная безопасность
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_KB_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_KCP_moscow_B_KB_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_KB_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_KB_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_KB_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        17,
				TargetQuota:    18,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
		},
		// Культурология
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_Cultural_O.xlsx",
			TQListURL: "",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_Cultural_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_Cultural_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_Cultural_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        32,
				TargetQuota:    0,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},
		// Маркетинг и рыночная аналитика
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_Marketing_O.xlsx",
			TQListURL: "",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_Marketing_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_Marketing_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_Marketing_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        14,
				TargetQuota:    2,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},
		// Математика
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_Math_O.xlsx",
			TQListURL: "",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_Math_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_Math_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_Math_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        42,
				TargetQuota:    6,
				DedicatedQuota: 6,
				SpecialQuota:   6,
			},
		},
		// Медиакоммуникации
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_Media_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_KCP_moscow_B_Media_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_Media_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_Media_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_Media_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        31,
				TargetQuota:    1,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},
		// Международные отношения
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_MO_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_KCP_moscow_B_MO_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_MO_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_MO_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_MO_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        25,
				TargetQuota:    10,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
		},
		// Мировая экономика
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_WE_O.xlsx",
			TQListURL: "",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_WE_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_WE_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_WE_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        43,
				TargetQuota:    5,
				DedicatedQuota: 6,
				SpecialQuota:   6,
			},
		},
		// Политология
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_Political_O.xlsx",
			TQListURL: "",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_Political_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_Political_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_Political_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        28,
				TargetQuota:    4,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},
		// Прикладная математика
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_AM_O.xlsx",
			TQListURL: "",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_AM_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_AM_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_AM_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        52,
				TargetQuota:    12,
				DedicatedQuota: 8,
				SpecialQuota:   8,
			},
		},
		// Прикладная математика и информатика
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_AMI_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_KCP_moscow_B_AMI_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_AMI_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_AMI_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_AMI_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        126,
				TargetQuota:    18,
				DedicatedQuota: 18,
				SpecialQuota:   18,
			},
		},
		// Программная инженерия
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_SE_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_KCP_moscow_B_SE_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_SE_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_SE_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_SE_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        104,
				TargetQuota:    24,
				DedicatedQuota: 16,
				SpecialQuota:   16,
			},
		},
		// Психология
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_Psy_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_KCP_moscow_B_Psy_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_Psy_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_Psy_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_Psy_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        59,
				TargetQuota:    5,
				DedicatedQuota: 8,
				SpecialQuota:   8,
			},
		},
		// Совместная программа по экономике НИУ ВШЭ и РЭШ
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_RESH_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_KCP_moscow_B_RESH_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_RESH_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_RESH_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_RESH_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        39,
				TargetQuota:    9,
				DedicatedQuota: 6,
				SpecialQuota:   6,
			},
		},
		// Совместный бакалавриат НИУ ВШЭ и Центра педагогического мастерства
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_CPM_O.xlsx",
			TQListURL: "",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_CPM_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_CPM_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_CPM_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        16,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Социология
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_Soc_O.xlsx",
			TQListURL: "",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_Soc_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_Soc_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_Soc_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        85,
				TargetQuota:    6,
				DedicatedQuota: 12,
				SpecialQuota:   12,
			},
		},
		// Стратегия и продюсирование в коммуникациях
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_Producer_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_KCP_moscow_B_Producer_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_Producer_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_Producer_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_Producer_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        7,
				TargetQuota:    1,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},
		// Управление бизнесом
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_BBA_O.xlsx",
			TQListURL: "",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_BBA_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_BBA_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_BBA_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        21,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Физика
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_Physics_O.xlsx",
			TQListURL: "",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_Physics_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_Physics_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_Physics_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        37,
				TargetQuota:    3,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
		},
		// Филология
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_Philology_O.xlsx",
			TQListURL: "",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_Philology_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_Philology_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_Philology_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        30,
				TargetQuota:    5,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
		},
		// Философия
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_Phil_O.xlsx",
			TQListURL: "",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_Phil_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_Phil_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_Phil_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        30,
				TargetQuota:    2,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},
		// Фундаментальная и компьютерная лингвистика
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_FKL_O.xlsx",
			TQListURL: "",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_FKL_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_FKL_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_FKL_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        33,
				TargetQuota:    2,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
		},
		// Химия
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_Chem_O.xlsx",
			TQListURL: "",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_Chem_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_Chem_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_Chem_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        20,
				TargetQuota:    7,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},
		// Экономика
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_Economy_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_KCP_moscow_B_Economy_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_Economy_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_Economy_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_Economy_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        65,
				TargetQuota:    15,
				DedicatedQuota: 10,
				SpecialQuota:   10,
			},
		},
		// Экономика и анализ данных
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_EcAl_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_KCP_moscow_B_EcAl_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_EcAl_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_EcAl_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_EcAl_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        28,
				TargetQuota:    7,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
		},
		// Экономика и статистика
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_Stat_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_KCP_moscow_B_Stat_O_det.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_Stat_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_Stat_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_Stat_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        22,
				TargetQuota:    5,
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},
		// Юриспруденция
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_LAW_O.xlsx",
			TQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_KCP_moscow_B_LAW_O.xlsx",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_LAW_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_LAW_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_LAW_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        98,
				TargetQuota:    22,
				DedicatedQuota: 15,
				SpecialQuota:   15,
			},
		},
		// Юриспруденция: цифровой юрист
		&oldhse.HttpHeadingSource{
			RCListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OM_moscow_B_Dlawyer_O.xlsx",
			TQListURL: "",
			DQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_SK_moscow_B_Dlawyer_O.xlsx",
			SQListURL: "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_OP_moscow_B_Dlawyer_O.xlsx",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_B_Dlawyer_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        19,
				TargetQuota:    5,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},

		// TODO The following headings do not have capacities determined:

		// Дизайн и разработка информационных продуктов
		&oldhse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_K_DisDevInfP_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        0,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   0,
			},
		},
		// Международная программа по экономике и финансам
		&oldhse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_K_icef_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        0,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   0,
			},
		},
		// Мода
		&oldhse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_K_irmo_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        0,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   0,
			},
		},
		// Прикладной анализ данных
		&oldhse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "https://enrol.hse.ru/storage/public_report_2024/moscow/Bachelors/KS_BVI_moscow_K_Data_O.xlsx",
			HeadingCapacities: core.Capacities{
				Regular:        0,
				TargetQuota:    0,
				DedicatedQuota: 0,
				SpecialQuota:   0,
			},
		},

		// TODO The following headings do not have list URLs determined:

		// Актёр
		&oldhse.HttpHeadingSource{
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
		// Бизнес-информатика
		&oldhse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "",
			HeadingCapacities: core.Capacities{
				Regular:        91,
				TargetQuota:    13,
				DedicatedQuota: 13,
				SpecialQuota:   13,
			},
		},
		// Глобальные цифровые коммуникации (онлайн)
		&oldhse.HttpHeadingSource{
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
		// Кинопроизводство
		&oldhse.HttpHeadingSource{
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
		// Компьютерные науки и анализ данных (онлайн)
		&oldhse.HttpHeadingSource{
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
		// Международная программа «Международные отношения и глобальные исследования»/ International Program «International Relations and Global Studies
		&oldhse.HttpHeadingSource{
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
		// Международный бизнес (реализуется на английском языке)
		&oldhse.HttpHeadingSource{
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
		// Монголия и Тибет
		&oldhse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "",
			HeadingCapacities: core.Capacities{
				Regular:        7,
				TargetQuota:    1,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},
		// Программа двух дипломов НИУ ВШЭ и Университета Кёнхи «Экономика и политика в Азии»/KHU-HSE Double Degree Program Economics and Politics in Asia
		&oldhse.HttpHeadingSource{
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
		// Разработка игр и цифровых продуктов
		&oldhse.HttpHeadingSource{
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
		// Реклама и связи с общественностью
		&oldhse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "",
			HeadingCapacities: core.Capacities{
				Regular:        2,
				TargetQuota:    1,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},
		// Современное искусство
		&oldhse.HttpHeadingSource{
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
		// Технологии анализа данных в бизнесе
		&oldhse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "",
			HeadingCapacities: core.Capacities{
				Regular:        21,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Турция и тюркский мир
		&oldhse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "",
			HeadingCapacities: core.Capacities{
				Regular:        7,
				TargetQuota:    1,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},
		// Управление в креативных индустриях
		&oldhse.HttpHeadingSource{
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
		// Управление цифровым продуктом
		&oldhse.HttpHeadingSource{
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
		// Химия новых материалов
		&oldhse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "",
			HeadingCapacities: core.Capacities{
				Regular:        9,
				TargetQuota:    2,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},
		// Экономический анализ (онлайн)
		&oldhse.HttpHeadingSource{
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
		// Эфиопия и арабский мир
		&oldhse.HttpHeadingSource{
			RCListURL: "",
			TQListURL: "",
			DQListURL: "",
			SQListURL: "",
			BListURL:  "",
			HeadingCapacities: core.Capacities{
				Regular:        7,
				TargetQuota:    1,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},
		// Юриспруденция: правовое регулирование бизнеса
		&oldhse.HttpHeadingSource{
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
