package haljson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Resource represents a Resource with Links and Embeds with Data
type Resource[T any] struct {
	Links  *Links  `json:"_links,omitempty"`
	Embeds *Embeds `json:"_embedded,omitempty"`
	// When serializing to JSON we need to handle this specially
	Data map[string]T `json:"-"`
}

// Self is used to add a self link
func (r *Resource[T]) Self(uri string) {
	r.Links.Self = &Link{Href: uri}
}

// AddLink adds a link to reltype
func (r *Resource[T]) AddLink(reltype string, link *Link) error {
	return r.Links.AddLink(reltype, link)
}

// AddEmbed adds a Resource by reltype
func (r *Resource[T]) AddEmbed(reltype string, embed *Resource[any]) error {
	if _, ok := r.Embeds.Relations[reltype]; !ok {
		r.Embeds.Relations[reltype] = []Resource[any]{}
	}
	r.Embeds.Relations[reltype] = append(r.Embeds.Relations[reltype], *embed)
	return nil
}

// AddCurie adds a curie to the links
func (r *Resource[T]) AddCurie(curie *Curie) error {
	return r.Links.AddCurie(curie)
}

// MarshalJSON marshals a resource properly
func (r *Resource[T]) MarshalJSON() ([]byte, error) {
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
	sortedKeys := make([]string, 0, len(r.Data))
	for k := range r.Data {
		sortedKeys = append(sortedKeys, k)
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
func (r *Resource[T]) UnmarshalJSON(b []byte) error {
	temp := make(map[string]any)
	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}

	embedded := NewEmbeds()

	if _, ok := temp[EMBEDDED]; ok {
		// re marshal embedded and links
		embeddedjson, err := json.Marshal(temp[EMBEDDED])
		if err != nil {
			return err
		}
		err = json.Unmarshal(embeddedjson, &embedded)
		if err != nil {
			return err
		}
	}
	r.Embeds = embedded
	delete(temp, EMBEDDED)

	links := NewLinks()
	if _, ok := temp[LINKS]; ok {
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

	if _, ok := temp[CURIES]; ok {
		var curies []Curie
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

	// Whatever is left over shove into Data, converting from any to T
	r.Data = make(map[string]T)
	for k, v := range temp {
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		var typedValue T
		err = json.Unmarshal(data, &typedValue)
		if err != nil {
			return err
		}
		r.Data[k] = typedValue
	}
	return nil
}

// NewResource creates a Resource and initializes it
func NewResource[T any]() *Resource[T] {
	return &Resource[T]{
		Data:   make(map[string]T),
		Links:  NewLinks(),
		Embeds: NewEmbeds(),
	}
}
