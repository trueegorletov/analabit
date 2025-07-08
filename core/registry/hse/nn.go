package hse

import (
	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/source/hse"
)

const (
	nnCode = "hse_nn"
	nnName = "ВШЭ (НН)"
)

func nnSourcesList() []source.HeadingSource {
	return []source.HeadingSource{
		// Компьютерные науки и технологии (Бизнес-информатика)
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/nn/Bachelors/BD_nn_KNT_BI_O.xlsx",
			Capacities: core.Capacities{
				Regular:        44,
				TargetQuota:    7,
				DedicatedQuota: 7,
				SpecialQuota:   7,
			},
		},
		// Компьютерные науки и технологии (Прикладная математика и информатика)
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/nn/Bachelors/BD_nn_KNT_AMI_O.xlsx",
			Capacities: core.Capacities{
				Regular:        44,
				TargetQuota:    7,
				DedicatedQuota: 7,
				SpecialQuota:   7,
			},
		},
		// Компьютерные науки и технологии (Программная инженерия)
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/nn/Bachelors/BD_nn_KNT_SE_O.xlsx",
			Capacities: core.Capacities{
				Regular:        44,
				TargetQuota:    7,
				DedicatedQuota: 7,
				SpecialQuota:   7,
			},
		},
		// Международный бакалавриат по бизнесу и экономике
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/nn/Bachelors/BD_nn_MBBE_O.xlsx",
			Capacities: core.Capacities{
				Regular:        78,
				TargetQuota:    12,
				DedicatedQuota: 18,
				SpecialQuota:   12,
			},
		},
		// Технологии искусственного и дополненного интеллекта
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/nn/Bachelors/BD_nn_Ait_O.xlsx",
			Capacities: core.Capacities{
				Regular:        6,
				TargetQuota:    1,
				DedicatedQuota: 2,
				SpecialQuota:   1,
			},
		},
		// Филология
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/nn/Bachelors/BD_nn_Philology_O.xlsx",
			Capacities: core.Capacities{
				Regular:        16,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},
		// Фундаментальная и прикладная лингвистика
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/nn/Bachelors/BD_nn_Ling_O.xlsx",
			Capacities: core.Capacities{
				Regular:        22,
				TargetQuota:    3,
				DedicatedQuota: 2,
				SpecialQuota:   3,
			},
		},
		// Фундаментальная и прикладная математика
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/nn/Bachelors/BD_nn_Math_O.xlsx",
			Capacities: core.Capacities{
				Regular:        26,
				TargetQuota:    4,
				DedicatedQuota: 6,
				SpecialQuota:   4,
			},
		},
		// Юриспруденция
		&hse.HTTPHeadingSource{
			URL: "https://enrol.hse.ru/storage/public_report_2025/nn/Bachelors/BD_nn_LAW_O.xlsx",
			Capacities: core.Capacities{
				Regular:        35,
				TargetQuota:    6,
				DedicatedQuota: 8,
				SpecialQuota:   6,
			},
		},

		// ===
		// FILTERED OUT
		// ===
		// Дизайн (KCP = 0)
		// Иностранные языки и межкультурная бизнес-коммуникация (KCP = 0)
		// Программная инженерия (очно-заочная форма обучения) (KCP = 0)
		// Цифровой маркетинг (KCP = 0)
		// Экономика и бизнес (очно-заочная форма обучения) (KCP = 0)

	}
}
