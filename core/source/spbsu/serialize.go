package spbsu

import (
	"encoding/gob"
)

func init() {
	gob.RegisterName("SpbsuHttpHeadingSource", &HttpHeadingSource{})
}
