package haljson

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Curie respresents a curie
type Curie struct {
	Name      string `json:"name"`
	Href      string `json:"href"`
	Templated bool   `json:"templated"`
}

// Link represents a link
type Link struct {
	Deprecation *string `json:"deprecation,omitempty"`
	Href        string  `json:"href,omitempty"`
	HrefLang    *string `json:"hreflang,omitempty"`
	Profile     *string `json:"profile,omitempty"`
	Title       *string `json:"title,omitempty"`
	Type        *string `json:"type,omitempty"`
}

// Links is a container of Link, mapped by relation, and contains Curies
type Links struct {
	Self   *Link    `json:"-"`
	Curies *[]Curie `json:"curies,omitempty"`
	// When serializing to JSON we need to handle this specially
	Relations map[string][]*Link
}

// MarshalJSON to marshal Links properly
func (l *Links) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	firstrun := true
	if l.Self != nil {
		firstrun = false
		jsonValue, err := json.Marshal(l.Self)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(fmt.Sprintf("\"self\": %s", string(jsonValue)))
	}
	if l.Curies != nil {
		if !firstrun {
			buffer.WriteString(",")
		} else {
			firstrun = false
		}
		jsonValue, err := json.Marshal(l.Curies)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(fmt.Sprintf("\"curies\": %s", string(jsonValue)))
	}
	for key, links := range l.Relations {
		if !firstrun {
			buffer.WriteString(",")
		} else {
			firstrun = false
		}
		jsonValue, err := json.Marshal(links)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(fmt.Sprintf("\"%s\": %s", key, string(jsonValue)))
	}
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

// UnmarshalJSON to unmarshal links
func (l *Links) UnmarshalJSON(b []byte) error {
	return nil
}

// Embeds holds embedded relations by reltype
type Embeds struct {
	Relations map[string][]Resource
}

// MarshalJSON marshals embeds
func (e *Embeds) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	firstrun := true
	for key, links := range e.Relations {
		if !firstrun {
			buffer.WriteString(",")
		} else {
			firstrun = false
		}
		jsonValue, err := json.Marshal(links)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(fmt.Sprintf("\"%s\": %s", key, string(jsonValue)))
	}
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshals embeds
func (e *Embeds) UnmarshalJSON(b []byte) error {
	return nil
}

// Resource represents a Resource with Links and Embeds with Data
type Resource struct {
	Links  *Links  `json:"_links,omitempty"`
	Embeds *Embeds `json:"_embedded,omitempty"`
	// When serializing to JSON we need to handle this specially
	Data map[string]interface{}
}

// Self is used to add a self link
func (r *Resource) Self(uri string) {
	r.Links.Self = &Link{Href: uri}
}

// AddLink adds a link to reltype
func (r *Resource) AddLink(reltype string, link *Link) error {
	if _, ok := r.Links.Relations[reltype]; !ok {
		r.Links.Relations[reltype] = []*Link{}
	}

	// Check if curied and that if curied, curie exists
	curieExists := false
	if strings.Index(reltype, ":") > 0 {
		parts := strings.Split(reltype, ":")
		var curies *[]Curie
		curies = r.Links.Curies
		if curies == nil {
			return errors.New("Must add curie before adding a curied link")
		}
		for _, curie := range *curies {
			if parts[0] == curie.Name {
				curieExists = true
			}
		}
		if !curieExists {
			return errors.New("Must add curie before adding a curied link")
		}
	}
	r.Links.Relations[reltype] = append(r.Links.Relations[reltype], link)
	return nil
}

// AddEmbed adds a Resource by reltype
func (r *Resource) AddEmbed(reltype string, embed *Resource) error {
	if _, ok := r.Embeds.Relations[reltype]; !ok {
		r.Embeds.Relations[reltype] = []Resource{}
	}
	r.Embeds.Relations[reltype] = append(r.Embeds.Relations[reltype], *embed)
	return nil
}

// AddCurie adds a curie to the links
func (r *Resource) AddCurie(curie *Curie) error {
	if r.Links.Curies == nil {
		r.Links.Curies = &[]Curie{}
	}
	*r.Links.Curies = append(*r.Links.Curies, *curie)
	return nil
}

// MarshalJSON marshals a resource properly
func (r *Resource) MarshalJSON() ([]byte, error) {
	var obj map[string]interface{}
	obj = make(map[string]interface{})
	for key, val := range r.Data {
		obj[key] = val
	}
	if len(r.Links.Relations) > 0 || r.Links.Self != nil {
		obj["_links"] = r.Links
	}

	if len(r.Embeds.Relations) > 0 {
		obj["_embedded"] = r.Embeds
	}
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// UnmarshalJSON unmarshals embeds
func (r *Resource) UnmarshalJSON(b []byte) error {
	return nil
}

// NewResource creates a Resource and initializes it
func NewResource() *Resource {
	r := &Resource{}
	r.Data = make(map[string]interface{})
	r.Links = NewLinks()
	r.Embeds = NewEmbeds()
	return r
}

// NewLinks creates and initializes Links
func NewLinks() *Links {
	l := &Links{}
	l.Relations = make(map[string][]*Link)
	return l
}

// NewEmbeds creates and initializes Embeds
func NewEmbeds() *Embeds {
	e := &Embeds{}
	e.Relations = make(map[string][]Resource)
	return e
}
