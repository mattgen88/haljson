package haljson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// Resource represents a Resource with Links and Embeds with Data
type Resource struct {
	Links  *Links  `json:"_links,omitempty"`
	Embeds *Embeds `json:"_embedded,omitempty"`
	// When serializing to JSON we need to handle this specially
	Data map[string]interface{} `json:"-"`
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
			return ErrNoCurie
		}
		for _, curie := range *curies {
			if parts[0] == curie.Name {
				curieExists = true
			}
		}
		if !curieExists {
			return ErrNoCurie
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
	var links *string
	if len(r.Links.Relations) > 0 || r.Links.Self != nil {
		b, err := json.Marshal(r.Links)
		if err != nil {
			return nil, err
		}
		linkString := fmt.Sprintf("\"_links\": %s", string(b))
		links = &linkString
	}

	var embeds *string
	if len(r.Embeds.Relations) > 0 {
		b, err := json.Marshal(r.Embeds)
		if err != nil {
			return nil, err
		}
		embedString := fmt.Sprintf("\"_embedded\": %s", string(b))
		embeds = &embedString
	}

	var dataBuffer []string
	for key, val := range r.Data {
		b, err := json.Marshal(val)
		if err != nil {
			return nil, err
		}
		dataBuffer = append(dataBuffer, fmt.Sprintf("\"%s\": %s", key, string(b)))
	}

	var joined []string

	buffer := bytes.NewBufferString("{")

	if links != nil {
		joined = append(joined, *links)
	}

	if embeds != nil {
		joined = append(joined, *embeds)
	}

	if len(dataBuffer) > 0 {
		joined = append(joined, strings.Join(dataBuffer, ","))
	}

	buffer.WriteString(strings.Join(joined, ","))
	buffer.WriteString("}")

	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshals embeds
func (r *Resource) UnmarshalJSON(b []byte) error {
	var temp map[string]interface{}
	temp = make(map[string]interface{})
	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}

	if temp["_embedded"] != nil {
		// re marshal embedded and links
		embededjson, err := json.Marshal(temp["_embedded"])
		if err != nil {
			return err
		}
		embedded := NewEmbeds()
		err = json.Unmarshal(embededjson, &embedded)
		if err != nil {
			return err
		}
		r.Embeds = embedded
	}
	delete(temp, "_embedded")

	if temp["_links"] != nil {
		linksjson, err := json.Marshal(temp["_links"])
		if err != nil {
			return err
		}
		links := NewLinks()
		err = json.Unmarshal(linksjson, &links)
		if err != nil {
			return err
		}
		r.Links = links

	}
	delete(temp, "_links")

	// Whatever is left over shove into Data
	r.Data = temp

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
