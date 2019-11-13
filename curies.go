package haljson

// Curie respresents a curie
type Curie struct {
	Name      string `json:"name"`
	Href      string `json:"href"`
	Templated bool   `json:"templated"`
}
