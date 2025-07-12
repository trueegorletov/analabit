package mirea

import (
	"encoding/gob"
)

func init() {
	gob.Register(&HTTPHeadingSource{})
}