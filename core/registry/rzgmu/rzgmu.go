package rzgmu

import (
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/source/rzgmu"
)

var Varsity = source.VarsityDefinition{
	Name:           "РязГМУ",
	Code:           "rzgmu",
	HeadingSources: sourcesList(),
}

func sourcesList() []source.HeadingSource {
	return []source.HeadingSource{
		&rzgmu.HTTPHeadingSource{
			URL:         "https://rzgmu.ru/images/upload/priem/2025/hod/l_b.pdf",
			ProgramName: "Лечебное дело",
		},
		&rzgmu.HTTPHeadingSource{
			URL:         "https://rzgmu.ru/images/upload/priem/2025/hod/s_b.pdf",
			ProgramName: "Стоматология",
		},
		&rzgmu.HTTPHeadingSource{
			URL:         "https://rzgmu.ru/images/upload/priem/2025/hod/p_b.pdf",
			ProgramName: "Педиатрия",
		},
		&rzgmu.HTTPHeadingSource{
			URL:         "https://rzgmu.ru/images/upload/priem/2025/hod/m_b.pdf",
			ProgramName: "Медико-профилактическое дело",
		},
		&rzgmu.HTTPHeadingSource{
			URL:         "https://rzgmu.ru/images/upload/priem/2025/hod/f_b.pdf",
			ProgramName: "Фармация",
		},
		&rzgmu.HTTPHeadingSource{
			URL:         "https://rzgmu.ru/images/upload/priem/2025/hod/kp_b.pdf",
			ProgramName: "Клиническая психология",
		},
	}
}
