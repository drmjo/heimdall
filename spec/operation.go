package spec

type Operation struct {
	Description string               `json:"description,omitempty"`
	Responses   map[string]*Response `json:"responses"`
	Parameters  []*Parameter         `json:"parameters,omitempty"`
}

type Parameter struct {
	Name            string                 `json:"name"`
	In              string                 `json:"in"`
	Description     string                 `json:"description,omitempty"`
	Required        bool                   `json:"required,omitempty"`
	Deprecated      bool                   `json:"deprecated,omitempty"`
	AllowEmptyValue bool                   `json:"allowEmptyValue,omitempty"`
	Schema          map[string]interface{} `json:"schema,omitempty"`
}
