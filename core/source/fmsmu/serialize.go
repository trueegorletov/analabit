package fmsmu

import (
	"encoding/gob"
)

func init() {
	gob.Register(&HTTPHeadingSource{})
}