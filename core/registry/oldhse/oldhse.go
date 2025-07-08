package oldhse

import (
	"github.com/trueegorletov/analabit/core/source"
)

var VarsitySpb = source.VarsityDefinition{
	Code:           spbCode,
	Name:           spbName,
	HeadingSources: spbSourcesList(),
}

var VarsityMsk = source.VarsityDefinition{
	Code:           mskCode,
	Name:           mskName,
	HeadingSources: mskSourcesList(),
}
var VarsityNn = source.VarsityDefinition{
	Code:           nnCode,
	Name:           nnName,
	HeadingSources: nnSourcesList(),
}
var VarsityPerm = source.VarsityDefinition{
	Code:           permCode,
	Name:           permName,
	HeadingSources: permSourcesList(),
}
