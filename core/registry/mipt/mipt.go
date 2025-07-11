package mipt

import (
	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/source/mipt"
)

const (
	varsityCode = "mipt"
	varsityName = "МФТИ"
)

var Varsity = source.VarsityDefinition{
	Code:           varsityCode,
	Name:           varsityName,
	HeadingSources: sourcesList(),
}

func sourcesList() []source.HeadingSource {
	return []source.HeadingSource{

		// EXACT MATCH: 'Геокосмические науки и технологии'
		&mipt.HTTPHeadingSource{
			PrettyName:        "Геокосмические науки и технологии",
			RegularBVIListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvR2Vva29zbWljaGVza2llIG5hdWtpIGkgdGVraG5vbG9naWlfQnl1ZHpoZXRfTmEgb2JzaGNoaWtoIG9zbm92YW5peWFraC5odG1s",
			TargetQuotaListURLs: []string{
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvR2Vva29zbWljaGVza2llIG5hdWtpIGkgdGVraG5vbG9naWkgQU8gwqtHTlRzIMKrVHNlbnRyIEtlbGR5c2hhwrsgQU8gwqtSS1PCuyBBTyDCq1RzTklJbWFzaMK7IEFPIMKrTklJIFRQwrsgUEFPIMKrUktLIMKrRW5lcmdpeWHCu19Uc2VsZXZvZS5odG1s",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvR2Vva29zbWljaGVza2llIG5hdWtpIGkgdGVraG5vbG9naWkgQU8gwqtHTlRzIMKrVHNlbnRyIEtlbGR5c2hhwrsgQU8gwqtSS1PCuyBBTyDCq1RzTklJbWFzaMK7IEFPIMKrTklJIFRQwrsgUEFPIMKrUktLIMKrRW5lcmdpeWHCu19Uc2VsZXZvZS5odG1s",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvR2Vva29zbWljaGVza2llIG5hdWtpIGkgdGVraG5vbG9naWkgQU8gwqtHTlRzIMKrVHNlbnRyIEtlbGR5c2hhwrsgQU8gwqtSS1PCuyBBTyDCq1RzTklJbWFzaMK7IEFPIMKrTklJIFRQwrsgUEFPIMKrUktLIMKrRW5lcmdpeWHCu19Uc2VsZXZvZS5odG1s",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvR2Vva29zbWljaGVza2llIG5hdWtpIGkgdGVraG5vbG9naWkgQU8gwqtHTlRzIMKrVHNlbnRyIEtlbGR5c2hhwrsgQU8gwqtSS1PCuyBBTyDCq1RzTklJbWFzaMK7IEFPIMKrTklJIFRQwrsgUEFPIMKrUktLIMKrRW5lcmdpeWHCu19Uc2VsZXZvZS5odG1s",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvR2Vva29zbWljaGVza2llIG5hdWtpIGkgdGVraG5vbG9naWkgQU8gwqtHTlRzIMKrVHNlbnRyIEtlbGR5c2hhwrsgQU8gwqtSS1PCuyBBTyDCq1RzTklJbWFzaMK7IEFPIMKrTklJIFRQwrsgUEFPIMKrUktLIMKrRW5lcmdpeWHCu19Uc2VsZXZvZS5odG1s",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvR2Vva29zbWljaGVza2llIG5hdWtpIGkgdGVraG5vbG9naWkgQWt0c2lvbmVybm9lIG9ic2hjaGVzdHZvIMKrS29ycG9yYXRzaXlhIGtvc21pY2hlc2tpa2ggc2lzdGVtIHNwZXRzaWFsbm9nbyBuYXpuYWNoZW5peWEgwqtLb21ldGHCu19Uc2VsZXZvZS5odG1z",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvR2Vva29zbWljaGVza2llIG5hdWtpIGkgdGVraG5vbG9naWkgQU8gwqtOUE8gRW5lcmdvbWFzaMK7X1RzZWxldm9lLmh0bWw=",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvR2Vva29zbWljaGVza2llIG5hdWtpIGkgdGVraG5vbG9naWkgQU8gwqtPcmdhbml6YXRzaXlhIMKrQUdBVMK7X1RzZWxldm9lLmh0bWw=",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvR2Vva29zbWljaGVza2llIG5hdWtpIGkgdGVraG5vbG9naWkgSW55ZSBvcmdhbml6YXRzaWlfVHNlbGV2b2UuaHRtbA==",
			},
			DedicatedQuotaListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvR2Vva29zbWljaGVza2llIG5hdWtpIGkgdGVraG5vbG9naWlfQnl1ZHpoZXRfT3RkZWxuYXlhIGt2b3RhLmh0bWw=",
			SpecialQuotaListURL:   "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvR2Vva29zbWljaGVza2llIG5hdWtpIGkgdGVraG5vbG9naWlfQnl1ZHpoZXRfSW1leXVzaGNoaWUgb3NvYm9lIHByYXZvLmh0bWw=",
			Capacities: core.Capacities{
				Regular:        33,
				TargetQuota:    13, // 58-(33+6+6)
				DedicatedQuota: 6,
				SpecialQuota:   6,
			},
		},

		// EXACT MATCH: 'Биотехнология'
		&mipt.HTTPHeadingSource{
			PrettyName:        "Биотехнология",
			RegularBVIListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvQmlvdGVraG5vbG9naXlhX0J5dWR6aGV0X05hIG9ic2hjaGlraCBvc25vdmFuaXlha2guaHRtbA==",
			TargetQuotaListURLs: []string{
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvQmlvdGVraG5vbG9naXlhIElueWUgb3JnYW5pemF0c2lpX1RzZWxldm9lLmh0bWw=",
			},
			DedicatedQuotaListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvQmlvdGVraG5vbG9naXlhX0J5dWR6aGV0X090ZGVsbmF5YSBrdm90YS5odG1z",
			SpecialQuotaListURL:   "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvQmlvdGVraG5vbG9naXlhX0J5dWR6aGV0X0ltZXl1c2hjaGllIG9zb2JvZSBwcmF2by5odG1z",
			Capacities: core.Capacities{
				Regular:        32,
				TargetQuota:    8, // 50-(32+5+5)
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
		},

		// EXACT MATCH: 'Радиотехника и компьютерные технологии'
		&mipt.HTTPHeadingSource{
			PrettyName:        "Радиотехника и компьютерные технологии",
			RegularBVIListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUmFkaW90ZWtobmlrYSBpIGtvbXB5dXRlcm55ZSB0ZWtobm9sb2dpaV9CeXVkemhldF9OYSBvYnNoY2hpa2ggb3Nub3Zhbml5YWtoLmh0bWw=",
			TargetQuotaListURLs: []string{
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUmFkaW90ZWtobmlrYSBpIGtvbXB5dXRlcm55ZSB0ZWtobm9sb2dpaSBQdWJsaWNobm9lIGFrdHNpb25lcm5vZSBvYnNoY2hlc3R2byDCq1JhZGlvZml6aWthwrtfVHNlbGV2b2UuaHRtbA==",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUmFkaW90ZWtobmlrYSBpIGtvbXB5dXRlcm55ZSB0ZWtobm9sb2dpaSBGZWRlcmFsbm9lIGF2dG9ub21ub2UgdWNocmV6aGRlbmllIMKrR29zdWRhcnN0dmVubnl5IG5hdWNobm8taXNzbGVkb3ZhdGVsc2tpeSBpbnN0aXR1dCBhdmlhdHNpb25ueWtoIHNpc3RlbcK7X1RzZWxldm9lLmh0bWw=",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUmFkaW90ZWtobmlrYSBpIGtvbXB5dXRlcm55ZSB0ZWtobm9sb2dpaSBQdWJsaWNobm9lIGFrdHNpb25lcm5vZSBvYnNoY2hlc3R2byDCq05hdWNobm8tcHJvaXp2b2RzdHZlbm5vZSBvYmVkaW5lbmllIMKrQWxtYXrCuyBpbWVuaSBha2FkZW1pa2EgQS5BLlJhc3BsZXRpbmHCu19Uc2VsZXZvZS5odG1z",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUmFkaW90ZWtobmlrYSBpIGtvbXB5dXRlcm55ZSB0ZWtobm9sb2dpaSBBS1RzSU9ORVJOT0UgT0JTaGNoRVNUVk8gwqtNVHNTVMK7X1RzZWxldm9lLmh0bWw=",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUmFkaW90ZWtobmlrYSBpIGtvbXB5dXRlcm55ZSB0ZWtobm9sb2dpaSBJbnllIG9yZ2FuaXphdHNpaV9Uc2VsZXZvZS5odG1z",
			},
			DedicatedQuotaListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUmFkaW90ZWtobmlrYSBpIGtvbXB5dXRlcm55ZSB0ZWtobm9sb2dpaV9CeXVkemhldF9PdGRlbG5heWEga3ZvdGEuaHRtbA==",
			SpecialQuotaListURL:   "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUmFkaW90ZWtobmlrYSBpIGtvbXB5dXRlcm55ZSB0ZWtobm9sb2dpaV9CeXVkemhldF9JbWV5dXNoY2hpZSBvc29ib2UgcHJhdm8uaHRtbA==",
			Capacities: core.Capacities{
				Regular:        68,
				TargetQuota:    11,
				DedicatedQuota: 9,
				SpecialQuota:   9,
			},
		},

		// EXACT MATCH: 'Общая и прикладная физика'
		&mipt.HTTPHeadingSource{
			PrettyName:        "Общая и прикладная физика",
			RegularBVIListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvT2JzaGNoYXlhIGkgcHJpa2xhZG5heWEgZml6aWthX0J5dWR6aGV0X05hIG9ic2hjaGlraCBvc25vdmFuaXlha2guaHRtbA==",
			TargetQuotaListURLs: []string{
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvT2JzaGNoYXlhIGkgcHJpa2xhZG5heWEgZml6aWthIElueWUgb3JnYW5pemF0c2lpX1RzZWxldm9lLmh0bWw=",
			},
			DedicatedQuotaListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvT2JzaGNoYXlhIGkgcHJpa2xhZG5heWEgZml6aWthX0J5dWR6aGV0X090ZGVsbmF5YSBrdm90YS5odG1z",
			SpecialQuotaListURL:   "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvT2JzaGNoYXlhIGkgcHJpa2xhZG5heWEgZml6aWthX0J5dWR6aGV0X0ltZXl1c2hjaGllIG9zb2JvZSBwcmF2by5odG1z",
			Capacities: core.Capacities{
				Regular:        137,
				TargetQuota:    2, // 179-(137+20+20)
				DedicatedQuota: 20,
				SpecialQuota:   20,
			},
		},

		// EXACT MATCH: 'Биофизика и биоинформатика'
		&mipt.HTTPHeadingSource{
			PrettyName:        "Биофизика и биоинформатика",
			RegularBVIListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvQmlvZml6aWthIGkgYmlvaW5mb3JtYXRpa2FfQnl1ZHpoZXRfTmEgb2JzaGNoaWtoIG9zbm92YW5peWFraC5odG1z",
			TargetQuotaListURLs: []string{
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvQmlvZml6aWthIGkgYmlvaW5mb3JtYXRpa2EgSW55ZSBvcmdhbml6YXRzaWlfVHNlbGV2b2UuaHRtbA==",
			},
			DedicatedQuotaListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvQmlvZml6aWthIGkgYmlvaW5mb3JtYXRpa2FfQnl1ZHpoZXRfT3RkZWxuYXlhIGt2b3RhLmh0bWw=",
			SpecialQuotaListURL:   "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvQmlvZml6aWthIGkgYmlvaW5mb3JtYXRpa2FfQnl1ZHpoZXRfSW1leXVzaGNoaWUgb3NvYm9lIHByYXZvLmh0bWw=",
			Capacities: core.Capacities{
				Regular:        35,
				TargetQuota:    5, // 50-(35+5+5)
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
		},

		// EXACT MATCH: 'Математическое моделирование и теория управления'
		&mipt.HTTPHeadingSource{
			PrettyName:        "Математическое моделирование и теория управления",
			RegularBVIListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvTWF0ZW1hdGljaGVza29lIG1vZGVsaXJvdmFuaWUgaSB0ZW9yaXlhIHVwcmF2bGVuaXlhX0J5dWR6aGV0X05hIG9ic2hjaGlraCBvc25vdmFuaXlha2guaHRtbA==",
			TargetQuotaListURLs: []string{
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvTWF0ZW1hdGljaGVza29lIG1vZGVsaXJvdmFuaWUgaSB0ZW9yaXlhIHVwcmF2bGVuaXlhIE9CU2hjaEVTVFZPIFMgT0dSQU5JQ2hFTk5PWSBPVFZFVFNUVkVOTk9TVFl1IMKrTk0tVEVLaMK7X1RzZWxldm9lLmh0bWw=",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvTWF0ZW1hdGljaGVza29lIG1vZGVsaXJvdmFuaWUgaSB0ZW9yaXlhIHVwcmF2bGVuaXlhIElueWUgb3JnYW5pemF0c2lpX1RzZWxldm9lLmh0bWw=",
			},
			DedicatedQuotaListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvTWF0ZW1hdGljaGVza29lIG1vZGVsaXJvdmFuaWUgaSB0ZW9yaXlhIHVwcmF2bGVuaXlhX0J5dWR6aGV0X090ZGVsbmF5YSBrdm90YS5odG1z",
			SpecialQuotaListURL:   "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvTWF0ZW1hdGljaGVza29lIG1vZGVsaXJvdmFuaWUgaSB0ZW9yaXlhIHVwcmF2bGVuaXlhX0J5dWR6aGV0X0ltZXl1c2hjaGllIG9zb2JvZSBwcmF2by5odG1z",
			Capacities: core.Capacities{
				Regular:        15,
				TargetQuota:    4,
				DedicatedQuota: 5,
				SpecialQuota:   5,
			},
		},

		// EXACT MATCH: 'Техническая физика'
		&mipt.HTTPHeadingSource{
			PrettyName:        "Техническая физика",
			RegularBVIListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvVGVraG5pY2hlc2theWEgZml6aWthX0J5dWR6aGV0X05hIG9ic2hjaGlraCBvc25vdmFuaXlha2guaHRtbA==",
			TargetQuotaListURLs: []string{
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvVGVraG5pY2hlc2theWEgZml6aWthIEFPIMKrTVogwqtBcnNlbmFswrtfVHNlbGV2b2UuaHRtbA==",
			},
			DedicatedQuotaListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvVGVraG5pY2hlc2theWEgZml6aWthX0J5dWR6aGV0X090ZGVsbmF5YSBrdm90YS5odG1z",
			SpecialQuotaListURL:   "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvVGVraG5pY2hlc2theWEgZml6aWthX0J5dWR6aGV0X0ltZXl1c2hjaGllIG9zb2JvZSBwcmF2by5odG1z",
			Capacities: core.Capacities{
				Regular:        10,
				TargetQuota:    1, // 15-(10+2+2)
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},

		// EXACT MATCH: 'Системное программирование и прикладная математика'
		&mipt.HTTPHeadingSource{
			PrettyName:        "Системное программирование и прикладная математика",
			RegularBVIListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvU2lzdGVtbm9lIHByb2dyYW1taXJvdmFuaWUgaSBwcmlrbGFkbmF5YSBtYXRlbWF0aWthX0J5dWR6aGV0X05hIG9ic2hjaGlraCBvc25vdmFuaXlha2guaHRtbA==",
			TargetQuotaListURLs: []string{
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvU2lzdGVtbm9lIHByb2dyYW1taXJvdmFuaWUgaSBwcmlrbGFkbmF5YSBtYXRlbWF0aWthIElueWUgb3JnYW5pemF0c2lpX1RzZWxldm9lLmh0bWw=",
			},
			DedicatedQuotaListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvU2lzdGVtbm9lIHByb2dyYW1taXJvdmFuaWUgaSBwcmlrbGFkbmF5YSBtYXRlbWF0aWthX0J5dWR6aGV0X090ZGVsbmF5YSBrdm90YS5odG1s",
			SpecialQuotaListURL:   "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvU2lzdGVtbm9lIHByb2dyYW1taXJvdmFuaWUgaSBwcmlrbGFkbmF5YSBtYXRlbWF0aWthX0J5dWR6aGV0X0ltZXl1c2hjaGllIG9zb2JvZSBwcmF2by5odG1s",
			Capacities: core.Capacities{
				Regular:        24,
				TargetQuota:    3,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},

		// EXACT MATCH: 'Авиационные технологии и автономные транспортные системы'
		&mipt.HTTPHeadingSource{
			PrettyName:        "Авиационные технологии и автономные транспортные системы",
			RegularBVIListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvQXZpYXRzaW9ubnllIHRla2hub2xvZ2lpIGkgYXZ0b25vbW55ZSB0cmFuc3BvcnRueWUgc2lzdGVteV9CeXVkemhldF9OYSBvYnNoY2hpa2ggb3Nub3Zhbml5YWtoLmh0bWw=",
			TargetQuotaListURLs: []string{
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvQXZpYXRzaW9ubnllIHRla2hub2xvZ2lpIGkgYXZ0b25vbW55ZSB0cmFuc3BvcnRueWUgc2lzdGVteSBGZWRlcmFsbm9lIGF2dG9ub21ub2UgdWNocmV6aGRlbmllIMKrVHNlbnRyYWxueXkgaW5zdGl0dXQgYXZpYXRzaW9ubm9nbyBtb3Rvcm9zdHJvZW5peWEgaW1lbmkgUC5JX1RzZWxldm9lLmh0bWw=",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvQXZpYXRzaW9ubnllIHRla2hub2xvZ2lpIGkgYXZ0b25vbW55ZSB0cmFuc3BvcnRueWUgc2lzdGVteSBQdWJsaWNobm9lIGFrdHNpb25lcm5vZSBvYnNoY2hlc3R2byDCq05hdWNobm8tcHJvaXp2b2RzdHZlbm5vZSBvYmVkaW5lbmllIMKrQWxtYXrCuyBpbWVuaSBha2FkZW1pa2EgQV9Uc2VsZXZvZS5odG1z",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvQXZpYXRzaW9uLiB0ZWtobm9sb2dpaSBpIGF2dG9ub20uIHRyYW5zcG9ydG55ZSBzaXN0ZW15IEZlZGVyYWxub2UgYXZ0b25vbW5vZSB1Y2hyZXpoZGVuaWUgwqtUc2VudHJhbG55eSBhZXJvZ2lkcm9kaW5hbWljaGVza2l5IGluc3RpdHV0IG1vbGVrdWx5YXJub3kgZWxla3Ryb25pa2nCu19Uc2VsZXZvZS5odG1z",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvQXZpYXRzaW9uLiB0ZWtobm9sb2dpaSBpIGF2dG9ub20uIHRyYW5zcG9ydG55ZSBzaXN0ZW15IEFrdHNpb25lcm5vZSBvYnNoY2hlc3R2byDCq0luZXJ0c2lhbG55ZSB0ZWtobm9sb2dpaSDCq1Rla2hub2tvbXBsZWtzYcK7X1RzZWxldm9lLmh0bWw=",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvQXZpYXRzaW9uLiB0ZWtobm9sb2dpaSBpIGF2dG9ub20uIHRyYW5zcG9ydG55ZSBzaXN0ZW15IFB1YmxpY2hub2UgYWt0c2lvbmVybm9lIG9ic2hjaGVzdHZvIMKrQXZpYXRzaW9ubnl5IGtvbXBsZWtzIGltLiBTLlYuIElseXVzaGluYcK7X1RzZWxldm9lLmh0bWw=",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvQXZpYXRzaW9ubnllIHRla2hub2xvZ2lpIGkgYXZ0b25vbW55ZSB0cmFuc3BvcnRueWUgc2lzdGVteSBBa3RzaW9uZXJub2Ugb2JzaGNoZXN0dm8gwqtVcmFsc2tpeSB6YXZvZCBncmF6aGRhbnNrb3kgYXZpYXRzaWnCu19Uc2VsZXZvZS5odG1z",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvQXZpYXRzaW9uLiB0ZWtobm9sb2dpaSBpIGF2dG9ub20uIHRyYW5zcG9ydG55ZSBzaXN0ZW15IEFrdHNpb25lcm5vZSBvYnNoY2hlc3R2byDCq0tvbnRzZXJuIHZvemR1c2huby1rb3NtaWNoZXNrb3kgb2Jvcm9ueSDCq0FsbWF6IC0gQW50ZXnCu19Uc2VsZXZvZS5odG1z",
			},
			DedicatedQuotaListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvQXZpYXRzaW9ubnllIHRla2hub2xvZ2lpIGkgYXZ0b25vbW55ZSB0cmFuc3BvcnRueWUgc2lzdGVteV9CeXVkemhldF9PdGRlbG5heWEga3ZvdGEuaHRtbA==",
			SpecialQuotaListURL:   "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvQXZpYXRzaW9ubnllIHRla2hub2xvZ2lpIGkgYXZ0b25vbW55ZSB0cmFuc3BvcnRueWUgc2lzdGVteV9CeXVkemhldF9JbWV5dXNoY2hpZSBvc29ib2UgcHJhdm8uaHRtbA==",
			Capacities: core.Capacities{
				Regular:        14,
				TargetQuota:    12, // 28-(14+1+1)
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		// EXACT MATCH: 'Природоподобные, плазменные и ядерные технологии'
		&mipt.HTTPHeadingSource{
			PrettyName:        "Природоподобные, плазменные и ядерные технологии",
			RegularBVIListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJpcm9kb3BvZG9ibnllIHBsYXptZW5ueWUgaSB5YWRlcm55ZSB0ZWtobm9sb2dpaV9CeXVkemhldF9OYSBvYnNoY2hpa2ggb3Nub3Zhbml5YWtoLmh0bWw=",
			TargetQuotaListURLs: []string{
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvU2lzdGVtbm9lIHByb2dyYW1taXJvdmFuaWUgaSBwcmlrbGFkbmF5YSBtYXRlbWF0aWthIElueWUgb3JnYW5pemF0c2lpX1RzZWxldm9lLmh0bWw=",
			},
			DedicatedQuotaListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvU2lzdGVtbm9lIHByb2dyYW1taXJvdmFuaWUgaSBwcmlrbGFkbmF5YSBtYXRlbWF0aWthX0J5dWR6aGV0X090ZGVsbmF5YSBrdm90YS5odG1z",
			SpecialQuotaListURL:   "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvU2lzdGVtbm9lIHByb2dyYW1taXJvdmFuaWUgaSBwcmlrbGFkbmF5YSBtYXRlbWF0aWthX0J5dWR6aGV0X0ltZXl1c2hjaGllIG9zb2JvZSBwcmF2by5odG1z",
			Capacities: core.Capacities{
				Regular:        24,
				TargetQuota:    5,
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},

		// EXACT MATCH: 'Электроника и наноэлектроника'
		&mipt.HTTPHeadingSource{
			PrettyName:        "Электроника и наноэлектроника",
			RegularBVIListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvRWxla3Ryb25pa2EgaSBuYW5vZWxla3Ryb25pa2FfQnl1ZHpoZXRfTmEgb2JzaGNoaWtoIG9zbm92YW5peWFraC5odG1s",
			TargetQuotaListURLs: []string{
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvRWxla3Ryb25pa2EgaSBuYW5vZWxla3Ryb25pa2EgQWt0c2lvbmVybm9lIG9ic2hjaGVzdHZvIMKrTWlrcm9uwrtfVHNlbGV2b2UuaHRtbA==",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvRWxla3Ryb25pa2EgaSBuYW5vZWxla3Ryb25pa2EgQWt0c2lvbmVybm9lIG9ic2hjaGVzdHZvIMKrTmF1Y2huby1pc3NsZWRvdmF0ZWxza2l5IGluc3RpdHV0IG1vbGVrdWx5YXJub3kgZWxla3Ryb25pa2nCu19Uc2VsZXZvZS5odG1s",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvRWxla3Ryb25pa2EgaSBuYW5vZWxla3Ryb25pa2EgQWt0c2lvbmVybm9lIG9ic2hjaGVzdHZvIMKrTmF1Y2huby1wcm9penZvZHN0dmVubm9lIHByZWRwcml5YXRpZSDCq1Rvcml5wrtfVHNlbGV2b2UuaHRtbA==",
			},
			DedicatedQuotaListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvRWxla3Ryb25pa2EgaSBuYW5vZWxla3Ryb25pa2FfQnl1ZHpoZXRfT3RkZWxuYXlhIGt2b3RhLmh0bWw=",
			SpecialQuotaListURL:   "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvRWxla3Ryb25pa2EgaSBuYW5vZWxla3Ryb25pa2FfQnl1ZHpoZXRfSW1leXVzaGNoaWUgb3NvYm9lIHByYXZvLmh0bWw=",
			Capacities: core.Capacities{
				Regular:        22,
				TargetQuota:    5, // 35 - (22+4+4) = 5 per screenshot
				DedicatedQuota: 4,
				SpecialQuota:   4,
			},
		},

		// EXACT MATCH: 'Физика перспективных технологий'
		&mipt.HTTPHeadingSource{
			PrettyName:        "Физика перспективных технологий",
			RegularBVIListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvRml6aWthIHBlcnNwZWt0aXZueWtoIHRla2hub2xvZ2l5X0J5dWR6aGV0X05hIG9ic2hjaGlraCBvc25vdmFuaXlha2guaHRtbA==",
			TargetQuotaListURLs: []string{
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvRml6aWthIHBlcnNwZWt0aXZueWtoIHRla2hub2xvZ2l5IEFrdHNpb25lcm5vZSBvYnNoY2hlc3R2byDCq05QTyDCq09yaW9uwrtfVHNlbGV2b2UuaHRtbA==",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvRml6aWthIHBlcnNwZWt0aXZueWtoIHRla2hub2xvZ2l5IEFrdHNpb25lcm5vZSBvYnNoY2hlc3R2byDCq05hdWNobm8taXNzbGVkb3ZhdGVsc2tpeSBpbnN0aXR1dCDCq1BvbHl1c8K7IGltLiBNLkYuU3RlbG1ha2hhwrtfVHNlbGV2b2UuaHRtbA==",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvRml6aWthIHBlcnNwZWt0aXZueWtoIHRla2hub2xvZ2l5IElueWUgb3JnYW5pemF0c2lpX1RzZWxldm9lLmh0bWw=",
			},
			DedicatedQuotaListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvRml6aWthIHBlcnNwZWt0aXZueWtoIHRla2hub2xvZ2l5X0J5dWR6aGV0X090ZGVsbmF5YSBrdm90YS5odG1z",
			SpecialQuotaListURL:   "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvRml6aWthIHBlcnNwZWt0aXZueWtoIHRla2hub2xvZ2l5X0J5dWR6aGV0X0ltZXl1c2hjaGllIG9zb2JvZSBwcmF2by5odG1z",
			Capacities: core.Capacities{
				Regular:        58,
				TargetQuota:    9, // 85 - (58+9+9) = 9 per screenshot
				DedicatedQuota: 9,
				SpecialQuota:   9,
			},
		},

		// EXACT MATCH: 'Компьютерные технологии и вычислительная техника'
		&mipt.HTTPHeadingSource{
			PrettyName:        "Компьютерные технологии и вычислительная техника",
			RegularBVIListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvS29tcHl1dGVybnllIHRla2hub2xvZ2lpIGkgdnljaGlzbGl0ZWxuYXlhIHRla2huaWthX0J5dWR6aGV0X05hIG9ic2hjaGlraCBvc25vdmFuaXlha2guaHRtbA==",
			TargetQuotaListURLs: []string{
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvS29tcHl1dGVybnllIHRla2hub2xvZ2lpIGkgdnljaGlzbGl0ZWxuYXlhIHRla2huaWthIEFrdHNpb25lcm5vZSBvYnNoY2hlc3R2byBOYXVjaG5vLXByb2l6dm9kc3R2ZW5ueXkgdHNlbnRyIMKrRWxla3Ryb25ueWUgdnljaGlzbGl0ZWxuby1pbmZvcm1hdHNpb25ueWUgc2lzdGVtecK7X1RzZWxldm9lLmh0bWw=",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvS29tcHl1dGVybnllIHRla2hub2xvZ2lpIGkgdnljaGlzbGl0ZWxuYXlhIHRla2huaWthIEFLVHNJT05FUk5PRSBPQlNoY2hFU1RWTyDCq01Uc1NUwrtfVHNlbGV2b2UuaHRtbA==",
			},
			DedicatedQuotaListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvS29tcHl1dGVybnllIHRla2hub2xvZ2lpIGkgdnljaGlzbGl0ZWxuYXlhIHRla2huaWthX0J5dWR6aGV0X090ZGVsbmF5YSBrdm90YS5odG1z",
			SpecialQuotaListURL:   "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvS29tcHl1dGVybnllIHRla2hub2xvZ2lpIGkgdnljaGlzbGl0ZWxuYXlhIHRla2huaWthX0J5dWR6aGV0X0ltZXl1c2hjaGllIG9zb2JvZSBwcmF2by5odG1z",
			Capacities: core.Capacities{
				Regular:        10,
				TargetQuota:    2, // 16 - (10+2+2) = 2 per screenshot
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},

		// EXACT MATCH: 'Программная инженерия и компьютерные технологии'
		&mipt.HTTPHeadingSource{
			PrettyName:        "Программная инженерия и компьютерные технологии",
			RegularBVIListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJvZ3JhbW1uYXlhIGluemhlbmVyaXlhIGkga29tcHl1dGVybnllIHRla2hub2xvZ2lpX0J5dWR6aGV0X05hIG9ic2hjaGlraCBvc25vdmFuaXlha2guaHRtbA==",
			TargetQuotaListURLs: []string{
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJvZ3JhbW1uYXlhIGluemhlbmVyaXlhIGkga29tcHl1dGVybnllIHRla2hub2xvZ2lpIEZlZGVyYWxub2UgYXZ0b25vbW5vZSB1Y2hyZXpoZGVuaWUgwqtHb3N1ZGFyc3R2ZW5ueXkgbmF1Y2huby1pc3NsZWRvdmF0ZWxza2l5IGluc3RpdHV0IGF2aWF0c2lvbm55a2ggc2lzdGVfVHNlbGV2b2UuaHRtbA==",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJvZ3JhbW1uYXlhIGluemhlbmVyaXlhIGkga29tcHl1dGVybnllIHRla2hub2xvZ2lpIEFrdHNpb25lcm5vZSBvYnNoY2hlc3R2byDCq0l6aGV2c2tpeSBlbGVrdHJvbWVraGFuaWNoZXNraXkgemF2b2QgwqtLdXBvbMK7X1RzZWxldm9lLmh0bWw=",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJvZ3JhbW1uYXlhIGluemhlbmVyaXlhIGkga29tcHl1dGVybnllIHRla2hub2xvZ2lpIFB1YmxpY2hub2UgYWt0c2lvbmVybm9lIG9ic2hjaGVzdHZvIMKrTmF1Y2huby1wcm9penZvZHN0dmVubm9lIG9iZWRpbmVuaWUgwqtBbG1hesK7IGltZW5pIGFrYWRlbWlrYSBBLkEuUmFzcGxlX1RzZWxldm9lLmh0bWw=",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJvZ3JhbS4gaW56aGVuZXJpeWEgaSBrb21wLiB0ZWtobm9sb2dpaSBGZWRlcmFsbm9lIGF2dG9ub21ub2UgdWNocmV6aGRlbmllIMKrVHNlbnRyYWxueXkgYWVyb2dpZHJvZGluYW1pY2hlc2tpeSBpbnN0aXR1dCBpbWVuaSBwcm9mZXNzb3JhIE4uRS4gWmh1a292c2tvZ2/Cu19Uc2VsZXZvZS5odG1s",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJvZ3JhbW1uYXlhIGluemhlbmVyaXlhIGkga29tcHl1dGVybnllIHRla2hub2xvZ2lpIEFrdHNpb25lcm5vZSBvYnNoY2hlc3R2byDCq0luZXJ0c2lhbG55ZSB0ZWtobm9sb2dpaSDCq1Rla2hub2tvbXBsZWtzYcK7X1RzZWxldm9lLmh0bWw=",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJvZ3JhbW1uYXlhIGluemhlbmVyaXlhIGkga29tcHl1dGVybnllIHRla2hub2xvZ2lpIEFrdHNpb25lcm5vZSBvYnNoY2hlc3R2byDCq0tvbnRzZXJuIHZvemR1c2huby1rb3NtaWNoZXNrb3kgb2Jvcm9ueSDCq0FsbWF6IC0gQW50ZXnCu19Uc2VsZXZvZS5odG1s",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJvZ3JhbW1uYXlhIGluemhlbmVyaXlhIGkga29tcHl1dGVybnllIHRla2hub2xvZ2lpIEFrdHNpb25lcm5vZSBvYnNoY2hlc3R2byDCq01vc2tvdnNraXkgbmF1Y2huby1pc3NsZWRvdmF0ZWxza2l5IGluc3RpdHV0IMKrQWdhdMK7X1RzZWxldm9lLmh0bWw=",
			},
			DedicatedQuotaListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJvZ3JhbW1uYXlhIGluemhlbmVyaXlhIGkga29tcHl1dGVybnllIHRla2hub2xvZ2lpX0J5dWR6aGV0X090ZGVsbmF5YSBrdm90YS5odG1s",
			SpecialQuotaListURL:   "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJvZ3JhbW1uYXlhIGluemhlbmVyaXlhIGkga29tcHl1dGVybnllIHRla2hub2xvZ2lpX0J5dWR6aGV0X0ltZXl1c2hjaGllIG9zb2JvZSBwcmF2by5odG1s",
			Capacities: core.Capacities{
				Regular:        15,
				TargetQuota:    15, // 95-(68+9+9)
				DedicatedQuota: 3,
				SpecialQuota:   3,
			},
		},

		// EXACT MATCH: 'Естественные и компьютерные науки'
		&mipt.HTTPHeadingSource{
			PrettyName:        "Естественные и компьютерные науки",
			RegularBVIListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvRXN0ZXN0dmVubnllIGkga29tcHl1dGVybnllIG5hdWtpX0J5dWR6aGV0X05hIG9ic2hjaGlraCBvc25vdmFuaXlha2guaHRtbA==",
			TargetQuotaListURLs: []string{
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvRXN0ZXN0dmVubnllIGkga29tcHl1dGVybnllIG5hdWtpIEFPIMKrTklJIFRQwrtfVHNlbGV2b2UuaHRtbA==",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvRXN0ZXN0dmVubnllIGkga29tcHl1dGVybnllIG5hdWtpIElueWUgb3JnYW5pemF0c2lpX1RzZWxldm9lLmh0bWw=",
			},
			DedicatedQuotaListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvRXN0ZXN0dmVubnllIGkga29tcHl1dGVybnllIG5hdWtpX0J5dWR6aGV0X090ZGVsbmF5YSBrdm90YS5odG1s",
			SpecialQuotaListURL:   "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvRXN0ZXN0dmVubnllIGkga29tcHl1dGVybnllIG5hdWtpX0J5dWR6aGV0X0ltZXl1c2hjaGllIG9zb2JvZSBwcmF2by5odG1s",
			Capacities: core.Capacities{
				Regular:        58,
				TargetQuota:    4,
				DedicatedQuota: 6,
				SpecialQuota:   6,
			},
		},

		// EXACT MATCH: 'Программная инженерия'
		&mipt.HTTPHeadingSource{
			PrettyName:        "Программная инженерия",
			RegularBVIListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJvZ3JhbW1uYXlhIGluemhlbmVyaXlhX0J5dWR6aGV0X05hIG9ic2hjaGlraCBvc25vdmFuaXlha2guaHRtbA==",
			TargetQuotaListURLs: []string{
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJvZ3JhbW1uYXlhIGluemhlbmVyaXlhIEFrdHNpb25lcm5vZSBvYnNoY2hlc3R2byBOYXVjaG5vLXByb2l6dm9kc3R2ZW5ueXkgdHNlbnRyIMKrRWxla3Ryb25ueWUgdnljaGlzbGl0ZWxuby1pbmZvcm1hdHNpb25ueWUgc2lzdGVtecK7X1RzZWxldm9lLmh0bWw=",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJvZ3JhbW1uYXlhIGluemhlbmVyaXlhIE9CU2hjaEVTVFZPIFMgT0dSQU5JQ2hFTk5PWSBPVFZFVFNUVkVOTk9TVFl1IMKrTk0tVEVLaMK7X1RzZWxldm9lLmh0bWw=",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJvZ3JhbW1uYXlhIGluemhlbmVyaXlhIEFrdHNpb25lcm5vZSBvYnNoY2hlc3R2byDCq0tvbnRzZXJuIEtpemx5YXJza2l5IGVsZWt0cm9tZWtoYW5pY2hlc2tpeSB6YXZvZMK7X1RzZWxldm9lLmh0bWw=",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJvZ3JhbW1uYXlhIGluemhlbmVyaXlhIE9ic2hjaGVzdHZvIHMgb2dyYW5pY2hlbm5veSBvdHZldHN0dmVubm9zdHl1IMKrVHNlbnRyIGJlem9wYXNub3N0aSBpbmZvcm1hdHNpacK7X1RzZWxldm9lLmh0bWw=",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJvZ3JhbW1uYXlhIGluemhlbmVyaXlhIElueWUgb3JnYW5pemF0c2lpX1RzZWxldm9lLmh0bWw=",
			},
			DedicatedQuotaListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJvZ3JhbW1uYXlhIGluemhlbmVyaXlhX0J5dWR6aGV0X090ZGVsbmF5YSBrdm90YS5odG1s",
			SpecialQuotaListURL:   "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJvZ3JhbW1uYXlhIGluemhlbmVyaXlhX0J5dWR6aGV0X0ltZXl1c2hjaGllIG9zb2JvZSBwcmF2by5odG1s",
			Capacities: core.Capacities{
				Regular:        40,
				TargetQuota:    9, // 63-(40+7+7)
				DedicatedQuota: 7,
				SpecialQuota:   7,
			},
		},

		// EXACT MATCH: 'Прикладная математика и информатика'
		&mipt.HTTPHeadingSource{
			PrettyName:        "Прикладная математика и информатика",
			RegularBVIListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJpa2xhZG5heWEgbWF0ZW1hdGlrYSBpIGluZm9ybWF0aWthX0J5dWR6aGV0X05hIG9ic2hjaGlraCBvc25vdmFuaXlha2guaHRtbA==",
			TargetQuotaListURLs: []string{
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJpa2xhZG5heWEgbWF0ZW1hdGlrYSBpIGluZm9ybWF0aWthIEFrdHNpb25lcm5vZSBvYnNoY2hlc3R2byDCq1JhbWVuc2tvZSBwcmlib3Jvc3Ryb2l0ZWxub2Uga29uc3RydWt0b3Jza29lIGJ5dXJvwrtfVHNlbGV2b2UuaHRtbA==",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJpa2xhZG5heWEgbWF0ZW1hdGlrYSBpIGluZm9ybWF0aWthIEFrdHNpb25lcm5vZSBvYnNoY2hlc3R2byDCq05hdWNobm8taXNzbGVkb3ZhdGVsc2tpeSBpbnN0aXR1dCBhdnRvbWF0aXppcm92YW5ueWtoIHNpc3RlbSBpIGtvbXBsZWtzb3Ygc3Z5YXppIMKrTmVwdHVuwrtfVHNlbGV2b2UuaHRtbA==",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJpa2xhZG5heWEgbWF0ZW1hdGlrYSBpIGluZm9ybWF0aWthIEFrdHNpb25lcm5vZSBvYnNoY2hlc3R2byDCq0tvbnN0cnVrdG9yc2tvZSBieXVybyBwcmlib3Jvc3Ryb2VuaXlhIGltLiBha2FkZW1pa2EgQS4gRy4gU2hpcHVub3ZhwrtfVHNlbGV2b2UuaHRtbA==",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJpa2xhZG5heWEgbWF0ZW1hdGlrYSBpIGluZm9ybWF0aWthIEZlZGVyYWxub2UgYXZ0b25vbW5vZSB1Y2hyZXpoZGVuaWUgwqtHb3N1ZGFyc3R2ZW5ueXkgbmF1Y2huby1pc3NsZWRvdmF0ZWxza2l5IGluc3RpdHV0IGF2aWF0c2lvbm55a2ggc2lzdGVtwrtfVHNlbGV2b2UuaHRtbA==",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJpa2xhZG5heWEgbWF0ZW1hdGlrYSBpIGluZm9ybWF0aWthIEFrdHNpb25lcm5vZSBvYnNoY2hlc3R2byBOYXVjaG5vLXByb2l6dm9kc3R2ZW5ueXkgdHNlbnRyIMKrRWxla3Ryb25ueWUgdnljaGlzbGl0ZWxuby1pbmZvcm1hdHNpb25ueWUgc2lzdGVtecK7X1RzZWxldm9lLmh0bWw=",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJpa2xhZG5heWEgbWF0ZW1hdGlrYSBpIGluZm9ybWF0aWthIEFPIMKrTlBPIEVuZXJnb21hc2jCu19Uc2VsZXZvZS5odG1s",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJpa2xhZG5heWEgbWF0ZW1hdGlrYSBpIGluZm9ybWF0aWthIEFPIMKrTklJIFRQwrtfVHNlbGV2b2UuaHRtbA==",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJpa2xhZG5heWEgbWF0ZW1hdGlrYSBpIGluZm9ybWF0aWthIEFPIMKrTlBQIMKrR2VvZml6aWthLUtvc21vc8K7X1RzZWxldm9lLmh0bWw=",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJpa2xhZG5heWEgbWF0ZW1hdGlrYSBpIGluZm9ybWF0aWthIElueWUgb3JnYW5pemF0c2lpX1RzZWxldm9lLmh0bWw=",
			},
			DedicatedQuotaListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJpa2xhZG5heWEgbWF0ZW1hdGlrYSBpIGluZm9ybWF0aWthX0J5dWR6aGV0X090ZGVsbmF5YSBrdm90YS5odG1s",
			SpecialQuotaListURL:   "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJpa2xhZG5heWEgbWF0ZW1hdGlrYSBpIGluZm9ybWF0aWthX0J5dWR6aGV0X0ltZXl1c2hjaGllIG9zb2JvZSBwcmF2by5odG1s",
			Capacities: core.Capacities{
				Regular:        118,
				TargetQuota:    18, // 170-(118+17+17)
				DedicatedQuota: 17,
				SpecialQuota:   17,
			},
		},

		// EXACT MATCH: 'Проектирование и разработка комплексных бизнес-приложений'
		&mipt.HTTPHeadingSource{
			PrettyName:            "Проектирование и разработка комплексных бизнес-приложений",
			RegularBVIListURL:     "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJvZWt0aXJvdmFuaWUgaSByYXpyYWJvdGthIGtvbXBsZWtzbnlraCBiaXpuZXMtcHJpbG96aGVuaXlfQnl1ZHpoZXRfTmEgb2JzaGNoaWtoIG9zbm92YW5peWFraC5odG1s",
			DedicatedQuotaListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJvZWt0aXJvdmFuaWUgaSByYXpyYWJvdGthIGtvbXBsZWtzbnlraCBiaXpuZXMtcHJpbG96aGVuaXlfQnl1ZHpoZXRfT3RkZWxuYXlhIGt2b3RhLmh0bWw=",
			SpecialQuotaListURL:   "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvUHJvZWt0aXJvdmFuaWUgaSByYXpyYWJvdGthIGtvbXBsZWtzbnlraCBiaXpuZXMtcHJpbG96aGVuaXlfQnl1ZHpoZXRfSW1leXVzaGNoaWUgb3NvYm9lIHByYXZvLmh0bWw=",
			Capacities: core.Capacities{
				Regular:        8,
				TargetQuota:    0,
				DedicatedQuota: 1,
				SpecialQuota:   1,
			},
		},

		&mipt.HTTPHeadingSource{
			PrettyName:        "Фундаментальная математика",
			RegularBVIListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvRnVuZGFtZW50YWxuYXlhIG1hdGVtYXRpa2FfQnl1ZHpoZXRfTmEgb2JzaGNoaWtoIG9zbm92YW5peWFraC5odG1s",
			TargetQuotaListURLs: []string{
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvRnVuZGFtZW50YWxuYXlhIG1hdGVtYXRpa2EgRmVkZXJhbG5vZSBnb3N1ZGFyc3R2ZW5ub2UgdW5pdGFybm9lIHByZWRwcml5YXRpZSDCq1RzZW50cmFsbnl5IG5hdWNobm8taXNzbGVkb3ZhdGVsc2tpeSBpbnN0aXR1dCBraGltaWkgaSBtZWtoYW5pa2nCu19Uc2VsZXZvZS5odG1s",
				"https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvRnVuZGFtZW50YWxuYXlhIG1hdGVtYXRpa2EgSW55ZSBvcmdhbml6YXRzaWlfVHNlbGV2b2UuaHRtbA==",
			},
			DedicatedQuotaListURL: "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvRnVuZGFtZW50YWxuYXlhIG1hdGVtYXRpa2FfQnl1ZHpoZXRfT3RkZWxuYXlhIGt2b3RhLmh0bWw=",
			SpecialQuotaListURL:   "https://priem.mipt.ru/applications_v2/YmFjaGVsb3IvRnVuZGFtZW50YWxuYXlhIG1hdGVtYXRpa2FfQnl1ZHpoZXRfSW1leXVzaGNoaWUgb3NvYm9lIHByYXZvLmh0bWw=",
			Capacities: core.Capacities{
				Regular:        10,
				TargetQuota:    2,
				DedicatedQuota: 2,
				SpecialQuota:   2,
			},
		},
	}
}
