package haljson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
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
	return r.Links.AddLink(reltype, link)
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
	return r.Links.AddCurie(curie)
}

// MarshalJSON marshals a resource properly
func (r *Resource) MarshalJSON() ([]byte, error) {

	// Marshal links
	var links *string
	if r.Links != nil && (len(r.Links.Relations) > 0 || r.Links.Self != nil) {
		b, err := json.Marshal(r.Links)
		if err != nil {
			return nil, err
		}
		linkString := fmt.Sprintf("\"%s\": %s", LINKS, string(b))
		links = &linkString
	}

	// Marshal Embeds
	var embeds *string
	if r.Embeds != nil && len(r.Embeds.Relations) > 0 {
		b, err := json.Marshal(r.Embeds)
		if err != nil {
			return nil, err
		}
		embedString := fmt.Sprintf("\"%s\": %s", EMBEDDED, string(b))
		embeds = &embedString
	}

	// Sort keys for data
	var sortedKeys = make([]string, len(r.Data))
	i := 0
	for k := range r.Data {
		sortedKeys[i] = k
		i++
	}
	sort.Strings(sortedKeys)

	// Marshal the data
	var dataBuffer []string
	for _, key := range sortedKeys {
		b, err := json.Marshal(r.Data[key])
		if err != nil {
			return nil, err
		}
		dataBuffer = append(dataBuffer, fmt.Sprintf("\"%s\": %s", key, string(b)))
	}

	// Produce JSON
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

	embedded := NewEmbeds()

	if temp[EMBEDDED] != nil {
		// re marshal embedded and links
		embededjson, err := json.Marshal(temp[EMBEDDED])
		if err != nil {
			return err
		}
		err = json.Unmarshal(embededjson, &embedded)
		if err != nil {
			return err
		}
	}
	r.Embeds = embedded
	delete(temp, EMBEDDED)

	links := NewLinks()
	if temp[LINKS] != nil {
		linksjson, err := json.Marshal(temp[LINKS])
		if err != nil {
			return err
		}
		err = json.Unmarshal(linksjson, &links)
		if err != nil {
			return err
		}
	}

	r.Links = links
	delete(temp, LINKS)

	if temp[CURIES] != nil {
		var curies *[]Curie
		curiesjson, err := json.Marshal(temp[CURIES])
		if err != nil {
			return err
		}
		err = json.Unmarshal(curiesjson, &curies)
		if err != nil {
			return err
		}
		r.Links.Curies = curies
	}
	delete(temp, CURIES)

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
