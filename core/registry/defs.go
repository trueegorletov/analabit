package registry

import (
	"analabit/core/registry/hse"
	"analabit/core/source"
)

var AllDefinitions = []source.VarsityDefinition{
	hse.VarsityMsk,
	hse.VarsitySpb,
	hse.VarsityNn,
	hse.VarsityPerm,
}
