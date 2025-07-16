package rsmu

import (
	"encoding/gob"
)

func init() {
	gob.RegisterName("RsmuHTTPHeadingSource", &HTTPHeadingSource{})
}
