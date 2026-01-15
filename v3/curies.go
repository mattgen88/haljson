package haljson

// Curie represents a curie
type Curie struct {
	Name      string `json:"name,omitempty"`
	Href      string `json:"href,omitempty"`
	Templated bool   `json:"templated,omitempty"`
}

// SetName sets the Curie name, chainable
func (c *Curie) SetName(name string) *Curie {
	c.Name = name
	return c
}

// SetHref sets the Curie href, chainable
func (c *Curie) SetHref(href string) *Curie {
	c.Href = href
	return c
}

// SetTemplated sets the Curie templated flag, chainable
func (c *Curie) SetTemplated(templated bool) *Curie {
	c.Templated = templated
	return c
}
