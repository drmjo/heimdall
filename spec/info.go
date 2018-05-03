package spec

type Info struct {
	Title   string   `json:"title"`
	Version string   `json:"version"`
	License *License `json:"license"`
}
