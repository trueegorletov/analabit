package registry

import (
	"github.com/trueegorletov/analabit/core/registry/fmsmu"
	"github.com/trueegorletov/analabit/core/registry/hse"
	"github.com/trueegorletov/analabit/core/registry/itmo"
	"github.com/trueegorletov/analabit/core/registry/mephi"
	"github.com/trueegorletov/analabit/core/registry/mipt"
	"github.com/trueegorletov/analabit/core/registry/mirea"
	"github.com/trueegorletov/analabit/core/registry/msu"
	"github.com/trueegorletov/analabit/core/registry/rsmu"
	"github.com/trueegorletov/analabit/core/registry/rzgmu"
	"github.com/trueegorletov/analabit/core/registry/spbstu"
	"github.com/trueegorletov/analabit/core/registry/spbsu"
	"github.com/trueegorletov/analabit/core/source"
)

var AllDefinitions = []source.VarsityDefinition{
	mirea.Varsity,
	hse.VarsityMsk,
	hse.VarsitySpb,
	hse.VarsityNn,
	hse.VarsityPerm,

	spbsu.Varsity,
	spbstu.Varsity,
	rsmu.Varsity,
	rzgmu.Varsity,
	itmo.Varsity,
	mipt.Varsity,
	fmsmu.Varsity,
	mephi.Varsity,
	msu.Varsity,
}
