package store

import (
	"github.com/datarhei/core/v16/restream/app"
)

type StoreData struct {
	Version uint64 `json:"version"`

	Process  map[string]*app.Process `json:"process"`
	Metadata struct {
		System  map[string]interface{}            `json:"system"`
		Process map[string]map[string]interface{} `json:"process"`
	} `json:"metadata"`
}

func NewStoreData() StoreData {
	c := StoreData{
		Version: 4,
	}

	c.Process = make(map[string]*app.Process)
	c.Metadata.System = make(map[string]interface{})
	c.Metadata.Process = make(map[string]map[string]interface{})

	return c
}

func (c *StoreData) IsEmpty() bool {
	if len(c.Process) != 0 {
		return false
	}

	if len(c.Metadata.Process) != 0 {
		return false
	}

	if len(c.Metadata.System) != 0 {
		return false
	}

	return true
}

func (c *StoreData) sanitize() {
	if c.Process == nil {
		c.Process = make(map[string]*app.Process)
	}
}
