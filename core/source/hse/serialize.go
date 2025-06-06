package hse

import (
	"encoding/gob"
)

func init() {
	gob.RegisterName("HseFileHeadingSource", &FileHeadingSource{})
	gob.RegisterName("HseHttpHeadingSource", &HttpHeadingSource{})
}
