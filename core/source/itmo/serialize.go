package itmo

import (
	"encoding/gob"
)

func init() {
	gob.RegisterName("ItmoHTTPHeadingSource", &HTTPHeadingSource{})
}
