package source

import (
	"encoding/gob"
	"io"
)

type VarsityDataCache struct {
	Definition   *VarsityDefinition
	Headings     []*HeadingData
	Applications []*ApplicationData
}

func NewVarsityDataCache(definition *VarsityDefinition) *VarsityDataCache {
	return &VarsityDataCache{
		Definition:   definition,
		Headings:     make([]*HeadingData, 0),
		Applications: make([]*ApplicationData, 0),
	}
}

func (c *VarsityDataCache) SaveHeadingData(heading *HeadingData) {
	if heading == nil {
		return
	}
	c.Headings = append(c.Headings, heading)
}

func (c *VarsityDataCache) SaveApplicationData(application *ApplicationData) {
	if application == nil {
		return
	}
	c.Applications = append(c.Applications, application)
}

func (c *VarsityDataCache) Reset() {
	c.Headings = nil
	c.Applications = nil
}

func SerializeList(caches []*VarsityDataCache, w io.Writer) error {
	encoder := gob.NewEncoder(w)

	return encoder.Encode(caches)
}

func DeserializeList(r io.Reader) ([]*VarsityDataCache, error) {
	decoder := gob.NewDecoder(r)
	var caches []*VarsityDataCache

	err := decoder.Decode(&caches)
	if err != nil {
		return nil, err
	}

	return caches, nil
}
