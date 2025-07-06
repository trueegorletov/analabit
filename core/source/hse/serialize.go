package hse

import (
	"encoding/gob"
)

func init() {
	gob.RegisterName("HseHTTPHeadingSource", &HTTPHeadingSource{})
	gob.RegisterName("HseFileHeadingSource", &FileHeadingSource{})
}
