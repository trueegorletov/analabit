package source

import (
	"encoding/gob"
	"io"
)

type VarsityDataCache struct {
	Definition        *VarsityDefinition
	HeadingsCache     []*HeadingData
	ApplicationsCache []*ApplicationData
}

func NewVarsityDataCache(definition *VarsityDefinition) *VarsityDataCache {
	return &VarsityDataCache{
		Definition:        definition,
		HeadingsCache:     make([]*HeadingData, 0),
		ApplicationsCache: make([]*ApplicationData, 0),
	}
}

func (c *VarsityDataCache) SaveHeadingData(heading *HeadingData) {
	if heading == nil {
		return
	}
	c.HeadingsCache = append(c.HeadingsCache, heading)
}

func (c *VarsityDataCache) SaveApplicationData(application *ApplicationData) {
	if application == nil {
		return
	}
	c.ApplicationsCache = append(c.ApplicationsCache, application)
}

func (c *VarsityDataCache) Reset() {
	c.HeadingsCache = nil
	c.ApplicationsCache = nil
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
