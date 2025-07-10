package mipt

import (
	"encoding/gob"
)

func init() {
	gob.RegisterName("MiptFileHeadingSource", &FileHeadingSource{})
	gob.RegisterName("MiptHTTPHeadingSource", &HTTPHeadingSource{})
}
