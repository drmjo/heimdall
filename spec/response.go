package spec

type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content"`
}
