package resolver

import "strings"

// CompetitionIDs represents the competition list IDs for different quota types
type CompetitionIDs struct {
	RegularBVI     string `json:"regular_bvi"`     // Regular and BVI (льготники) quota
	DedicatedQuota string `json:"dedicated_quota"` // Dedicated quota (целевая квота)
	SpecialQuota   string `json:"special_quota"`   // Special quota (особая квота)
	TargetQuota    string `json:"target_quota"`    // Target quota (квота для поступающих по договорам о целевом обучении)
}

// getMSUCompetitionIDs returns the mapping of MSU program names to their competition list IDs
// for different quota types. These IDs were extracted from the MSU admissions data.
func getMSUCompetitionIDs() map[string]CompetitionIDs {
	return map[string]CompetitionIDs{
		"Астрономия": {
			RegularBVI:     "193244",
			DedicatedQuota: "193246",
			SpecialQuota:   "193245",
			TargetQuota:    "",
		},
		"Фундаментальная и прикладная химия": {
			RegularBVI:     "193254",
			DedicatedQuota: "193257",
			SpecialQuota:   "193255",
			TargetQuota:    "193256",
		},
		"Фундаментальная физико-химическая инженерия": {
			RegularBVI:     "193460",
			DedicatedQuota: "193463",
			SpecialQuota:   "193461",
			TargetQuota:    "193462",
		},
		"Физико-химическая биология. Общая биология": {
			RegularBVI:     "193259",
			DedicatedQuota: "193262",
			SpecialQuota:   "193260",
			TargetQuota:    "193261",
		},
		"Биоинженерия и биотехнология, биофизика": {
			RegularBVI:     "193264",
			DedicatedQuota: "193267",
			SpecialQuota:   "193265",
			TargetQuota:    "",
		},
		"Социология": {
			RegularBVI:     "193423",
			DedicatedQuota: "193426",
			SpecialQuota:   "193424",
			TargetQuota:    "",
		},
		"Социология, общий профиль": {
			RegularBVI:     "193485",
			DedicatedQuota: "193488",
			SpecialQuota:   "193486",
			TargetQuota:    "",
		},
		"Государственный и муниципальный аудит": {
			RegularBVI:     "193470",
			DedicatedQuota: "193473",
			SpecialQuota:   "193471",
			TargetQuota:    "193472",
		},
		"Экономика": {
			RegularBVI:     "193378",
			DedicatedQuota: "193381",
			SpecialQuota:   "193379",
			TargetQuota:    "193380",
		},
		"Экономика, профиль экономическая теория и аналитика": {
			RegularBVI:     "",
			DedicatedQuota: "",
			SpecialQuota:   "",
			TargetQuota:    "",
		},
		"Биотехнология": {
			RegularBVI:     "193562",
			DedicatedQuota: "193565",
			SpecialQuota:   "193563",
			TargetQuota:    "",
		},
		"Регионоведение России": {
			RegularBVI:     "193445",
			DedicatedQuota: "193447",
			SpecialQuota:   "193446",
			TargetQuota:    "",
		},
		"Многопрофильный конкурс География, Картография и геоинформатика , Гидрометеорология, Экология и природопользование": {
			RegularBVI:     "388983",
			DedicatedQuota: "388985",
			SpecialQuota:   "388984",
			TargetQuota:    "388986",
		},
		"Образовательная программа \"Развитие стран Азии и Африки\"": {
			RegularBVI:     "389021",
			DedicatedQuota: "389024",
			SpecialQuota:   "389022",
			TargetQuota:    "390502",
		},
		"Юриспруденция (финансово-правовой профиль, профиль Расследование экономических преступлений)": {
			RegularBVI:     "193475",
			DedicatedQuota: "193478",
			SpecialQuota:   "193476",
			TargetQuota:    "193477",
		},
		"Юриспруденция (государственно-правовой, гражданско-правовой, уголовно-правовой профили)": {
			RegularBVI:     "193388",
			DedicatedQuota: "193391",
			SpecialQuota:   "193389",
			TargetQuota:    "193390",
		},
		"Международно-правовая (международно-правовой профиль)": {
			RegularBVI:     "",
			DedicatedQuota: "",
			SpecialQuota:   "",
			TargetQuota:    "",
		},
		"Менеджмент": {
			RegularBVI:     "193382",
			DedicatedQuota: "193385",
			SpecialQuota:   "193383",
			TargetQuota:    "193646",
		},
		"Менеджмент, общий профиль": {
			RegularBVI:     "",
			DedicatedQuota: "",
			SpecialQuota:   "",
			TargetQuota:    "",
		},
		"Менеджмент в культуре": {
			RegularBVI:     "",
			DedicatedQuota: "",
			SpecialQuota:   "",
			TargetQuota:    "",
		},
		"Менеджмент в спорте": {
			RegularBVI:     "",
			DedicatedQuota: "",
			SpecialQuota:   "",
			TargetQuota:    "",
		},
		"Музейный и галерейный менеджмент": {
			RegularBVI:     "",
			DedicatedQuota: "",
			SpecialQuota:   "",
			TargetQuota:    "",
		},
		"Менеджмент и маркетинг в сфере культуры и культурный туризм": {
			RegularBVI:     "",
			DedicatedQuota: "",
			SpecialQuota:   "",
			TargetQuota:    "",
		},
		"Менеджмент, профиль управление эффективностью и развитием бизнеса": {
			RegularBVI:     "",
			DedicatedQuota: "",
			SpecialQuota:   "",
			TargetQuota:    "",
		},
		"Управление бизнесом и предпринимательство": {
			RegularBVI:     "",
			DedicatedQuota: "",
			SpecialQuota:   "",
			TargetQuota:    "",
		},
		"Государственное и муниципальное управление": {
			RegularBVI:     "193523",
			DedicatedQuota: "193526",
			SpecialQuota:   "193524",
			TargetQuota:    "193525",
		},
		"Биоинженерия и биоинформатика": {
			RegularBVI:     "193455",
			DedicatedQuota: "193458",
			SpecialQuota:   "193456",
			TargetQuota:    "193652",
		},
		"Философия": {
			RegularBVI:     "193359",
			DedicatedQuota: "193362",
			SpecialQuota:   "193360",
			TargetQuota:    "193361",
		},
		"Экология и природопользование": {
			RegularBVI:     "193269",
			DedicatedQuota: "193272",
			SpecialQuota:   "193270",
			TargetQuota:    "",
		},
		"Экология и природопользование, экологический менеджмент и экобезопасность. Радиоэкология, управление земельными ресурсами и биологический контроль окружающей среды": {
			RegularBVI:     "193278",
			DedicatedQuota: "193281",
			SpecialQuota:   "193279",
			TargetQuota:    "",
		},
		"Биология почв. Химия почв. Земельные ресурсы и функционирование почв. Физика, мелиорация и эрозия почв. Агрохимия и агроэкология": {
			RegularBVI:     "193274",
			DedicatedQuota: "193277",
			SpecialQuota:   "193275",
			TargetQuota:    "193276",
		},
		"Геология": {
			RegularBVI:     "193284",
			DedicatedQuota: "193287",
			SpecialQuota:   "193285",
			TargetQuota:    "193286",
		},
		"Туризм": {
			RegularBVI:     "193289",
			DedicatedQuota: "193292",
			SpecialQuota:   "193290",
			TargetQuota:    "",
		},
		"Фармация": {
			RegularBVI:     "193304",
			DedicatedQuota: "193307",
			SpecialQuota:   "193305",
			TargetQuota:    "",
		},
		"Математика": {
			RegularBVI:     "193310",
			DedicatedQuota: "193313",
			SpecialQuota:   "193311",
			TargetQuota:    "193312",
		},
		"Механика": {
			RegularBVI:     "193314",
			DedicatedQuota: "193317",
			SpecialQuota:   "193315",
			TargetQuota:    "391190",
		},
		"Фундаментальная математика и математическая физика": {
			RegularBVI:     "193318",
			DedicatedQuota: "193321",
			SpecialQuota:   "193319",
			TargetQuota:    "",
		},
		"Космические исследования и космонавтика": {
			RegularBVI:     "193567",
			DedicatedQuota: "193570",
			SpecialQuota:   "193568",
			TargetQuota:    "193569",
		},
		"История": {
			RegularBVI:     "193325",
			DedicatedQuota: "193328",
			SpecialQuota:   "193326",
			TargetQuota:    "",
		},
		"История международных отношений": {
			RegularBVI:     "193329",
			DedicatedQuota: "193332",
			SpecialQuota:   "193330",
			TargetQuota:    "",
		},
		"СЛАВЯНСКАЯ ГРЕЧЕСКАЯ И КЛАССИЧЕСКАЯ ФИЛОЛОГИЯ": {
			RegularBVI:     "193347",
			DedicatedQuota: "193350",
			SpecialQuota:   "193348",
			TargetQuota:    "",
		},
		"Русский язык и литература": {
			RegularBVI:     "193339",
			DedicatedQuota: "193342",
			SpecialQuota:   "193340",
			TargetQuota:    "193341",
		},
		"Зарубежная филология": {
			RegularBVI:     "193343",
			DedicatedQuota: "193346",
			SpecialQuota:   "193344",
			TargetQuota:    "",
		},
		"Фундаментальная и прикладная лингвистика": {
			RegularBVI:     "193351",
			DedicatedQuota: "193354",
			SpecialQuota:   "193352",
			TargetQuota:    "193353",
		},
		"Религиоведение": {
			RegularBVI:     "193363",
			DedicatedQuota: "193366",
			SpecialQuota:   "193364",
			TargetQuota:    "",
		},
		"Реклама и связи с общественностью": {
			RegularBVI:     "193367",
			DedicatedQuota: "193370",
			SpecialQuota:   "193368",
			TargetQuota:    "",
		},
		"Журналистика": {
			RegularBVI:     "193394",
			DedicatedQuota: "193397",
			SpecialQuota:   "193395",
			TargetQuota:    "",
		},
		"Международная журналистика": {
			RegularBVI:     "193398",
			DedicatedQuota: "193401",
			SpecialQuota:   "193399",
			TargetQuota:    "",
		},
		"Медиакоммуникации": {
			RegularBVI:     "193402",
			DedicatedQuota: "193405",
			SpecialQuota:   "193403",
			TargetQuota:    "",
		},
		"Психология служебной деятельности": {
			RegularBVI:     "193413",
			DedicatedQuota: "193415",
			SpecialQuota:   "193414",
			TargetQuota:    "",
		},
		"Педагогика и психология девиантного поведения": {
			RegularBVI:     "193416",
			DedicatedQuota: "193419",
			SpecialQuota:   "193417",
			TargetQuota:    "193418",
		},
		"Прикладные математика и физика": {
			RegularBVI:     "193465",
			DedicatedQuota: "193468",
			SpecialQuota:   "193466",
			TargetQuota:    "391193",
		},
		"Общая политология": {
			RegularBVI:     "193496",
			DedicatedQuota: "193499",
			SpecialQuota:   "193497",
			TargetQuota:    "",
		},
		"Политический менеджмент и связи с общественностью": {
			RegularBVI:     "193501",
			DedicatedQuota: "193503",
			SpecialQuota:   "193502",
			TargetQuota:    "",
		},
		"Политология": {
			RegularBVI:     "193528",
			DedicatedQuota: "193531",
			SpecialQuota:   "193529",
			TargetQuota:    "",
		},
		"Продюсерство": {
			RegularBVI:     "193510",
			DedicatedQuota: "193512",
			SpecialQuota:   "193511",
			TargetQuota:    "",
		},
		"Математические и компьютерные методы решения задач естествознания": {
			RegularBVI:     "193230",
			DedicatedQuota: "193233",
			SpecialQuota:   "193231",
			TargetQuota:    "193232",
		},
		"Управление персоналом, общий профиль": {
			RegularBVI:     "193491",
			DedicatedQuota: "193494",
			SpecialQuota:   "193492",
			TargetQuota:    "",
		},
		"Управление персоналом, профиль стратегическое управление человеческими ресурсами": {
			RegularBVI:     "193518",
			DedicatedQuota: "193521",
			SpecialQuota:   "193519",
			TargetQuota:    "193520",
		},
		"Фундаментальная и прикладная физика": {
			RegularBVI:     "193240",
			DedicatedQuota: "193242",
			SpecialQuota:   "193241",
			TargetQuota:    "193243",
		},
		"Физика частиц и экстремальных состояний материи": {
			RegularBVI:     "193249",
			DedicatedQuota: "193251",
			SpecialQuota:   "193250",
			TargetQuota:    "",
		},
		"История искусств": {
			RegularBVI:     "193333",
			DedicatedQuota: "193335",
			SpecialQuota:   "193334",
			TargetQuota:    "",
		},
		"Культурология, прагматика и менеджмент культуры": {
			RegularBVI:     "193371",
			DedicatedQuota: "193373",
			SpecialQuota:   "193372",
			TargetQuota:    "",
		},
		"Культурология": {
			RegularBVI:     "193448",
			DedicatedQuota: "193450",
			SpecialQuota:   "193449",
			TargetQuota:    "",
		},
		"Клиническая психология": {
			RegularBVI:     "193409",
			DedicatedQuota: "193412",
			SpecialQuota:   "193410",
			TargetQuota:    "193411",
		},
		"Экспертная деятельность в управлении социально-политическими проектами": {
			RegularBVI:     "193427",
			DedicatedQuota: "193430",
			SpecialQuota:   "193428",
			TargetQuota:    "",
		},
		"Специальный перевод. Лингвистическое обеспечение межгосударственных отношений": {
			RegularBVI:     "193433",
			DedicatedQuota: "193435",
			SpecialQuota:   "193434",
			TargetQuota:    "",
		},
		"Перевод и переводоведение": {
			RegularBVI:     "193556",
			DedicatedQuota: "193558",
			SpecialQuota:   "193557",
			TargetQuota:    "",
		},
		"Группа программ Европейские исследования, Американские исследования": {
			RegularBVI:     "193441",
			DedicatedQuota: "193444",
			SpecialQuota:   "193442",
			TargetQuota:    "",
		},
		"Телевидение": {
			RegularBVI:     "193480",
			DedicatedQuota: "193483",
			SpecialQuota:   "193481",
			TargetQuota:    "",
		},
		"Конфликтология": {
			RegularBVI:     "193505",
			DedicatedQuota: "193508",
			SpecialQuota:   "193506",
			TargetQuota:    "",
		},
		"Международные отношения": {
			RegularBVI:     "193535",
			DedicatedQuota: "193538",
			SpecialQuota:   "193536",
			TargetQuota:    "",
		},
		"Глобальные политические процессы и дипломатия": {
			RegularBVI:     "193541",
			DedicatedQuota: "193544",
			SpecialQuota:   "193542",
			TargetQuota:    "",
		},
		"Глобальная экономика и управление. Глобальная энергетика и международный бизнес": {
			RegularBVI:     "193546",
			DedicatedQuota: "193549",
			SpecialQuota:   "193547",
			TargetQuota:    "390499",
		},
		"Международное гуманитарное сотрудничество": {
			RegularBVI:     "193551",
			DedicatedQuota: "193554",
			SpecialQuota:   "193552",
			TargetQuota:    "",
		},
		"Фундаментальная информатика и информационные технологии": {
			RegularBVI:     "193234",
			DedicatedQuota: "193237",
			SpecialQuota:   "193235",
			TargetQuota:    "193236",
		},
		"Химия, физика и механика материалов": {
			RegularBVI:     "193294",
			DedicatedQuota: "193297",
			SpecialQuota:   "193295",
			TargetQuota:    "",
		},
		"Лечебное дело": {
			RegularBVI:     "193299",
			DedicatedQuota: "193303",
			SpecialQuota:   "193300",
			TargetQuota:    "193301",
		},
		"Лингвистика": {
			RegularBVI:     "",
			DedicatedQuota: "",
			SpecialQuota:   "",
			TargetQuota:    "",
		},
		"Группа программ Межкультурная коммуникация, Теория и методика преподавания иностранных языков": {
			RegularBVI:     "193437",
			DedicatedQuota: "193440",
			SpecialQuota:   "193438",
			TargetQuota:    "193439",
		},
	}
}

// GetCompetitionIDsForProgram returns competition IDs for a given program name (case-insensitive)
func GetCompetitionIDsForProgram(programName string) (CompetitionIDs, bool) {
	competitionMap := getMSUCompetitionIDs()

	// Try exact match first
	if ids, exists := competitionMap[programName]; exists {
		return ids, true
	}

	// Try case-insensitive match
	upperProgramName := strings.ToUpper(programName)
	for name, ids := range competitionMap {
		if strings.ToUpper(name) == upperProgramName {
			return ids, true
		}
	}

	return CompetitionIDs{}, false
}
