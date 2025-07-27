package mephi

import (
	"encoding/gob"
)

func init() {
	gob.RegisterName("MephiHTTPHeadingSource", &HTTPHeadingSource{})
}
