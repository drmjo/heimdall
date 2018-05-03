package spec

import (
	"net/http"
	"strings"
)

type PathItem struct {
	Get     *Operation `json:"get,omitempty"`
	Post    *Operation `json:"post,omitempty"`
	Patch   *Operation `json:"patch,omitempty"`
	Put     *Operation `json:"put,omitempty"`
	Delete  *Operation `json:"delete,omitempty"`
	Servers []*Server  `json:"servers,omitempty"`
}

func (i *PathItem) GetOperation(operation string) *Operation {
	var o *Operation
	switch strings.ToUpper(operation) {
	case http.MethodGet:
		o = i.Get
	case http.MethodPost:
		o = i.Post
	case http.MethodPut:
		o = i.Put
	case http.MethodPatch:
		o = i.Patch
	case http.MethodDelete:
		o = i.Delete
	default:
		panic("unsupported operation")
	}

	return o
}
