package spbstu

import (
	"encoding/gob"
)

func init() {
	gob.Register(&HTTPHeadingSource{})
}
