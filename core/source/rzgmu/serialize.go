package rzgmu

import (
	"encoding/gob"
)

func init() {
	gob.RegisterName("RzgmuHttpHeadingSource", &HTTPHeadingSource{})
}
