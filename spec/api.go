package spec

import (
	"encoding/json"
	"log"
)

type Api struct {
	/**
	 * resolved through json Unmarshal
	 */
	OpenApi string               `json:"openapi"`
	Info    *Info                `json:"info"`
	Paths   map[string]*PathItem `json:"paths"`
}

/**
 * use to get a new Api Object can be used to
 * set defaults
 */
func NewApi(jsonData []byte) *Api {
	api := &Api{}

	err := json.Unmarshal(jsonData, api)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return api
}
