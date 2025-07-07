package registry

import (
	"analabit/core/registry/hse"
	"analabit/core/registry/rzgmu"
	"analabit/core/registry/spbsu"
	"analabit/core/source"
)

var AllDefinitions = []source.VarsityDefinition{
	hse.VarsityMsk,
	hse.VarsitySpb,
	hse.VarsityNn,
	hse.VarsityPerm,

	spbsu.Varsity,
	rzgmu.Varsity,
}
