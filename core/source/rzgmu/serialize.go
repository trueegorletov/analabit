package rzgmu

import (
	"encoding/gob"
)

func init() {
	gob.RegisterName("RzgmuHttpHeadingSource", &HTTPHeadingSource{})
	gob.RegisterName("RzgmuFileHeadingSource", &FileHeadingSource{})
}
