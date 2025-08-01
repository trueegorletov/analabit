package msu

import (
	"encoding/gob"
)

func init() {
	gob.RegisterName("MsuHTTPHeadingSource", &HTTPHeadingSource{})
}
