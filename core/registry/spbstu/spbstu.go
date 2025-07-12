// Package spbstu provides registry definitions for St. Petersburg State Technical University (SPbSTU).
package spbstu

import (
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/source/spbstu"
)

var Varsity = source.VarsityDefinition{
	Code:           "spbstu",
	Name:           "СПбПУ",
	HeadingSources: sourcesList(),
}

// sourcesList returns the list of SPbSTU HeadingSource definitions.
func sourcesList() []source.HeadingSource {
	return []source.HeadingSource{
		&spbstu.HTTPHeadingSource{
			PrettyName:           "Автоматизация технологических процессов и производств",
			RegularListID:        12,
			TargetQuotaListIDs:   []int{1670, 1673, 1672, 1676, 1679, 1678, 1675, 1671, 1677, 1674},
			DedicatedQuotaListID: 14,
			SpecialQuotaListID:   13,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Атомные станции: проектирование, эксплуатация и инжиниринг",
			RegularListID:        20,
			TargetQuotaListIDs:   []int{1682, 1686, 1685, 1684, 1687, 1683},
			DedicatedQuotaListID: 22,
			SpecialQuotaListID:   21,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Атомные станции: проектирование, эксплуатация и инжиниринг (Проектирование и эксплуатация атомных станций)",
			RegularListID:        25,
			TargetQuotaListIDs:   []int{1688, 1689},
			DedicatedQuotaListID: 27,
			SpecialQuotaListID:   26,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Бизнес-информатика",
			RegularListID:        35,
			TargetQuotaListIDs:   []int{1691, 1690, 1692},
			DedicatedQuotaListID: 37,
			SpecialQuotaListID:   36,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Биотехнические системы и технологии",
			RegularListID:        54,
			TargetQuotaListIDs:   []int{57},
			DedicatedQuotaListID: 56,
			SpecialQuotaListID:   55,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Биотехнология",
			RegularListID:        87,
			TargetQuotaListIDs:   []int{1693, 1694},
			DedicatedQuotaListID: 89,
			SpecialQuotaListID:   88,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Гостиничное дело",
			RegularListID:        135,
			TargetQuotaListIDs:   []int{138},
			DedicatedQuotaListID: 137,
			SpecialQuotaListID:   136,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Государственное и муниципальное управление",
			RegularListID:        141,
			TargetQuotaListIDs:   []int{1698, 1699},
			DedicatedQuotaListID: 143,
			SpecialQuotaListID:   142,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Дизайн",
			RegularListID:        167,
			TargetQuotaListIDs:   []int{},
			DedicatedQuotaListID: 169,
			SpecialQuotaListID:   168,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Дизайн архитектурной среды",
			RegularListID:        181,
			TargetQuotaListIDs:   []int{184},
			DedicatedQuotaListID: 183,
			SpecialQuotaListID:   182,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Зарубежное регионоведение",
			RegularListID:        193,
			TargetQuotaListIDs:   []int{196},
			DedicatedQuotaListID: 195,
			SpecialQuotaListID:   194,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Издательское дело",
			RegularListID:        208,
			TargetQuotaListIDs:   []int{211},
			DedicatedQuotaListID: 210,
			SpecialQuotaListID:   209,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Инноватика",
			RegularListID:        225,
			TargetQuotaListIDs:   []int{1702, 1701, 1700},
			DedicatedQuotaListID: 227,
			SpecialQuotaListID:   226,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Интеллектуальные системы в гуманитарной сфере",
			RegularListID:        243,
			TargetQuotaListIDs:   []int{246},
			DedicatedQuotaListID: 245,
			SpecialQuotaListID:   244,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Инфокоммуникационные технологии и системы связи",
			RegularListID:        260,
			TargetQuotaListIDs:   []int{1706, 1707, 1705},
			DedicatedQuotaListID: 262,
			SpecialQuotaListID:   261,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Информатика и вычислительная техника",
			RegularListID:        287,
			TargetQuotaListIDs:   []int{1709, 1713, 1711, 1710, 1714, 1712},
			DedicatedQuotaListID: 289,
			SpecialQuotaListID:   288,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Информационная безопасность",
			RegularListID:        298,
			TargetQuotaListIDs:   []int{1728, 1722, 1724, 1723, 1725, 1727, 1729, 1721, 1726},
			DedicatedQuotaListID: 300,
			SpecialQuotaListID:   299,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Информационная безопасность автоматизированных систем",
			RegularListID:        315,
			TargetQuotaListIDs:   []int{1731, 1732, 1733, 1730},
			DedicatedQuotaListID: 317,
			SpecialQuotaListID:   316,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Информационно-аналитические системы безопасности",
			RegularListID:        320,
			TargetQuotaListIDs:   []int{1734, 1735},
			DedicatedQuotaListID: 322,
			SpecialQuotaListID:   321,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Информационные системы и технологии",
			RegularListID:        340,
			TargetQuotaListIDs:   []int{1737, 1739, 1736, 1738},
			DedicatedQuotaListID: 342,
			SpecialQuotaListID:   341,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Компьютерная безопасность",
			RegularListID:        360,
			TargetQuotaListIDs:   []int{1743, 1741, 1744, 1745, 1748, 1742, 1746, 1740, 1747},
			DedicatedQuotaListID: 362,
			SpecialQuotaListID:   361,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Конструкторско-технологическое обеспечение машиностроительных производств",
			RegularListID:        376,
			TargetQuotaListIDs:   []int{1761, 1759, 1760, 1766, 1756, 1762, 1764, 1765, 1767, 1768, 1754, 1757, 1763, 1755, 1753, 1758},
			DedicatedQuotaListID: 378,
			SpecialQuotaListID:   377,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Лингвистика",
			RegularListID:        403,
			TargetQuotaListIDs:   []int{1769, 1770},
			DedicatedQuotaListID: 405,
			SpecialQuotaListID:   404,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Математика и компьютерные науки",
			RegularListID:        423,
			TargetQuotaListIDs:   []int{1771, 1773, 1772},
			DedicatedQuotaListID: 425,
			SpecialQuotaListID:   424,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Математическое обеспечение и администрирование информационных систем",
			RegularListID:        450,
			TargetQuotaListIDs:   []int{1775, 1774},
			DedicatedQuotaListID: 452,
			SpecialQuotaListID:   451,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Материаловедение и технологии материалов",
			RegularListID:        468,
			TargetQuotaListIDs:   []int{1776},
			DedicatedQuotaListID: 470,
			SpecialQuotaListID:   469,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Машиностроение",
			RegularListID:        499,
			TargetQuotaListIDs:   []int{1790, 1795, 1791, 1785, 1796, 1784, 1794, 1797, 1786, 1788, 1792, 1781, 1787, 1793, 1782, 1789, 1783},
			DedicatedQuotaListID: 501,
			SpecialQuotaListID:   500,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Менеджмент",
			RegularListID:        530,
			TargetQuotaListIDs:   []int{1802, 1800, 1803, 1798, 1801, 1799},
			DedicatedQuotaListID: 532,
			SpecialQuotaListID:   531,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Менеджмент (Международный бизнес (международная образовательная программа))",
			RegularListID:        535,
			TargetQuotaListIDs:   []int{},
			DedicatedQuotaListID: 537,
			SpecialQuotaListID:   536,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Металлургия (Цифровые технологии в металлургии)",
			RegularListID:        580,
			TargetQuotaListIDs:   []int{1812, 1811, 1810, 1809},
			DedicatedQuotaListID: 582,
			SpecialQuotaListID:   581,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Механика и математическое моделирование",
			RegularListID:        643,
			TargetQuotaListIDs:   []int{1814, 1813, 1815},
			DedicatedQuotaListID: 645,
			SpecialQuotaListID:   644,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Мехатроника и робототехника",
			RegularListID:        648,
			TargetQuotaListIDs:   []int{1823, 1824, 1821, 1820, 1822, 1819},
			DedicatedQuotaListID: 650,
			SpecialQuotaListID:   649,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Наземные транспортно-технологические средства",
			RegularListID:        684,
			TargetQuotaListIDs:   []int{1826, 1825},
			DedicatedQuotaListID: 686,
			SpecialQuotaListID:   685,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Нанотехнологии и микросистемная техника",
			RegularListID:        693,
			TargetQuotaListIDs:   []int{1827, 1828},
			DedicatedQuotaListID: 695,
			SpecialQuotaListID:   694,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Прикладная информатика",
			RegularListID:        765,
			TargetQuotaListIDs:   []int{1831, 1832},
			DedicatedQuotaListID: 767,
			SpecialQuotaListID:   766,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Прикладная математика и информатика",
			RegularListID:        776,
			TargetQuotaListIDs:   []int{1836, 1837, 1839, 1835, 1838},
			DedicatedQuotaListID: 778,
			SpecialQuotaListID:   777,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Прикладная механика",
			RegularListID:        791,
			TargetQuotaListIDs:   []int{794},
			DedicatedQuotaListID: 793,
			SpecialQuotaListID:   792,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Прикладные математика и физика",
			RegularListID:        833,
			TargetQuotaListIDs:   []int{836},
			DedicatedQuotaListID: 835,
			SpecialQuotaListID:   834,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Программная инженерия",
			RegularListID:        847,
			TargetQuotaListIDs:   []int{1861, 1854, 1856, 1863, 1857, 1858, 1862, 1855, 1853, 1852, 1859, 1860},
			DedicatedQuotaListID: 849,
			SpecialQuotaListID:   848,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Продукты питания животного происхождения",
			RegularListID:        866,
			TargetQuotaListIDs:   []int{869},
			DedicatedQuotaListID: 868,
			SpecialQuotaListID:   867,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Психолого-педагогическое образование",
			RegularListID:        883,
			TargetQuotaListIDs:   []int{886},
			DedicatedQuotaListID: 885,
			SpecialQuotaListID:   884,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Радиотехника",
			RegularListID:        900,
			TargetQuotaListIDs:   []int{1871, 1867, 1870, 1869, 1868, 1873, 1872, 1874, 1866},
			DedicatedQuotaListID: 902,
			SpecialQuotaListID:   901,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Реклама и связи с общественностью",
			RegularListID:        920,
			TargetQuotaListIDs:   []int{923},
			DedicatedQuotaListID: 922,
			SpecialQuotaListID:   921,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Сервис",
			RegularListID:        944,
			TargetQuotaListIDs:   []int{947},
			DedicatedQuotaListID: 946,
			SpecialQuotaListID:   945,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Статистика",
			RegularListID:        975,
			TargetQuotaListIDs:   []int{978},
			DedicatedQuotaListID: 977,
			SpecialQuotaListID:   976,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Строительство",
			RegularListID:        1004,
			TargetQuotaListIDs:   []int{1877, 1880, 1881, 1882, 1879, 1878},
			DedicatedQuotaListID: 1006,
			SpecialQuotaListID:   1005,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Строительство уникальных зданий и сооружений",
			RegularListID:        1055,
			TargetQuotaListIDs:   []int{1883, 1884, 1885},
			DedicatedQuotaListID: 1057,
			SpecialQuotaListID:   1056,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Таможенное дело",
			RegularListID:        1060,
			TargetQuotaListIDs:   []int{1063},
			DedicatedQuotaListID: 1062,
			SpecialQuotaListID:   1061,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Теплоэнергетика и теплотехника",
			RegularListID:        1108,
			TargetQuotaListIDs:   []int{1892, 1890, 1886, 1888, 1891, 1889, 1887},
			DedicatedQuotaListID: 1110,
			SpecialQuotaListID:   1109,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Техническая физика",
			RegularListID:        1141,
			TargetQuotaListIDs:   []int{1893, 1894},
			DedicatedQuotaListID: 1143,
			SpecialQuotaListID:   1142,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Технология продукции и организация общественного питания",
			RegularListID:        1177,
			TargetQuotaListIDs:   []int{1895},
			DedicatedQuotaListID: 1179,
			SpecialQuotaListID:   1178,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Технология транспортных процессов",
			RegularListID:        1190,
			TargetQuotaListIDs:   []int{1193},
			DedicatedQuotaListID: 1192,
			SpecialQuotaListID:   1191,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Технология художественной обработки материалов",
			RegularListID:        1206,
			TargetQuotaListIDs:   []int{1209},
			DedicatedQuotaListID: 1208,
			SpecialQuotaListID:   1207,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Техносферная безопасность",
			RegularListID:        1217,
			TargetQuotaListIDs:   []int{1900},
			DedicatedQuotaListID: 1219,
			SpecialQuotaListID:   1218,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Товароведение",
			RegularListID:        1249,
			TargetQuotaListIDs:   []int{1252},
			DedicatedQuotaListID: 1251,
			SpecialQuotaListID:   1250,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Торговое дело",
			RegularListID:        1254,
			TargetQuotaListIDs:   []int{1257},
			DedicatedQuotaListID: 1256,
			SpecialQuotaListID:   1255,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Торговое дело (Международная торговля (международная образовательная программа))",
			RegularListID:        1259,
			TargetQuotaListIDs:   []int{},
			DedicatedQuotaListID: 1261,
			SpecialQuotaListID:   1260,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Туризм",
			RegularListID:        1294,
			TargetQuotaListIDs:   []int{1297},
			DedicatedQuotaListID: 1296,
			SpecialQuotaListID:   1295,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Управление в технических системах",
			RegularListID:        1306,
			TargetQuotaListIDs:   []int{1904, 1903, 1902, 1901},
			DedicatedQuotaListID: 1308,
			SpecialQuotaListID:   1307,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Управление качеством",
			RegularListID:        1327,
			TargetQuotaListIDs:   []int{1908},
			DedicatedQuotaListID: 1329,
			SpecialQuotaListID:   1328,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Физика",
			RegularListID:        1350,
			TargetQuotaListIDs:   []int{1910, 1911},
			DedicatedQuotaListID: 1352,
			SpecialQuotaListID:   1351,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Экономика",
			RegularListID:        1428,
			TargetQuotaListIDs:   []int{1915, 1914, 1912, 1916, 1913},
			DedicatedQuotaListID: 1430,
			SpecialQuotaListID:   1429,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Экономика (Экономика цифрового предприятия (международная образовательная программа))",
			RegularListID:        1433,
			TargetQuotaListIDs:   []int{},
			DedicatedQuotaListID: 1435,
			SpecialQuotaListID:   1434,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Экономическая безопасность",
			RegularListID:        1443,
			TargetQuotaListIDs:   []int{1919, 1921, 1922, 1918, 1920},
			DedicatedQuotaListID: 1445,
			SpecialQuotaListID:   1444,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Электроника и наноэлектроника",
			RegularListID:        1460,
			TargetQuotaListIDs:   []int{1928, 1929, 1926, 1927},
			DedicatedQuotaListID: 1462,
			SpecialQuotaListID:   1461,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Электроэнергетика и электротехника",
			RegularListID:        1503,
			TargetQuotaListIDs:   []int{1931, 1934, 1938, 1932, 1935, 1939, 1930, 1933, 1940, 1936, 1937},
			DedicatedQuotaListID: 1505,
			SpecialQuotaListID:   1504,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Энергетическое машиностроение",
			RegularListID:        1569,
			TargetQuotaListIDs:   []int{1963, 1962, 1961, 1958, 1957, 1960, 1959},
			DedicatedQuotaListID: 1571,
			SpecialQuotaListID:   1570,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Юриспруденция",
			RegularListID:        1611,
			TargetQuotaListIDs:   []int{1965, 1964},
			DedicatedQuotaListID: 1613,
			SpecialQuotaListID:   1612,
		},

		&spbstu.HTTPHeadingSource{
			PrettyName:           "Ядерная энергетика и теплофизика",
			RegularListID:        1624,
			TargetQuotaListIDs:   []int{1966, 1968},
			DedicatedQuotaListID: 1626,
			SpecialQuotaListID:   1625,
		},
	}
}
