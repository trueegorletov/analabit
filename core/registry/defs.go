package registry

import (
	"github.com/trueegorletov/analabit/core/registry/hse"
	"github.com/trueegorletov/analabit/core/registry/rzgmu"
	"github.com/trueegorletov/analabit/core/registry/spbsu"
	"github.com/trueegorletov/analabit/core/source"
)

var AllDefinitions = []source.VarsityDefinition{
	hse.VarsityMsk,
	hse.VarsitySpb,
	hse.VarsityNn,
	hse.VarsityPerm,

	spbsu.Varsity,
	rzgmu.Varsity,
}
