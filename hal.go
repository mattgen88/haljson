package haljson

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	// NAME is a name link, curie property
	NAME = "name"
	// HREF is a href link, curie property
	HREF = "href"
	// HREFLANG is a hreflang link property
	HREFLANG = "hreflang"
	// TEMPLATED is a templated link, curie property
	TEMPLATED = "templated"
	// PROFILE is a profile link property
	PROFILE = "profile"
	// TITLE is a title link property
	TITLE = "title"
	// TYPE is a type link property
	TYPE = "type"
	// DEPRECATION is a deprecation link property
	DEPRECATION = "deprecation"
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
	var bufferData []string
	if l.Self != nil {
		jsonValue, err := json.Marshal(l.Self)
		if err != nil {
			return nil, err
		}
		bufferData = append(bufferData, fmt.Sprintf("\"self\": %s", string(jsonValue)))
	}
	if l.Curies != nil {
		jsonValue, err := json.Marshal(l.Curies)
		if err != nil {
			return nil, err
		}
		bufferData = append(bufferData, fmt.Sprintf("\"curies\": %s", string(jsonValue)))
	}
	for key, links := range l.Relations {
		jsonValue, err := json.Marshal(links)
		if err != nil {
			return nil, err
		}
		bufferData = append(bufferData, fmt.Sprintf("\"%s\": %s", key, string(jsonValue)))
	}
	joined := strings.Join(bufferData, ",")
	buffer := bytes.NewBufferString("{")
	buffer.WriteString(joined)
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

// UnmarshalJSON to unmarshal links
func (l *Links) UnmarshalJSON(b []byte) error {
	var temp map[string]interface{}
	temp = make(map[string]interface{})
	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}
	if _, ok := temp["curies"]; ok {
		var mycuries []Curie
		for _, curies := range temp["curies"].([]interface{}) {
			var curie Curie
			for k, v := range curies.(map[string]interface{}) {
				switch k {
				case NAME:
					curie.Name = v.(string)
				case HREF:
					curie.Href = v.(string)
				case TEMPLATED:
					curie.Templated = v.(bool)
				}
			}
			mycuries = append(mycuries, curie)
		}
		l.Curies = &mycuries
		delete(temp, "curies")
	}

	var self Link
	if _, ok := temp["self"]; ok {
		for k, v := range temp["self"].(map[string]interface{}) {
			switch k {
			case HREF:
				self.Href = v.(string)
			case DEPRECATION:
				var deprecation string
				deprecation = v.(string)
				self.Deprecation = &deprecation // hehe
			case HREFLANG:
				var hreflang string
				hreflang = v.(string)
				self.HrefLang = &hreflang
			case PROFILE:
				var profile string
				profile = v.(string)
				self.Profile = &profile
			case TITLE:
				var title string
				title = v.(string)
				self.Title = &title
			case TYPE:
				var typeval string
				typeval = v.(string)
				self.Type = &typeval
			}
		}
		l.Self = &self
		delete(temp, "self")
	}

	l.Relations = make(map[string][]*Link)
	for rel, v := range temp {
		var links []*Link
		for _, properties := range v.([]interface{}) {
			var link Link
			for key, property := range properties.(map[string]interface{}) {
				switch key {
				case HREF:
					link.Href = property.(string)
				case DEPRECATION:
					var deprecation string
					deprecation = property.(string)
					link.Deprecation = &deprecation // hehe
				case HREFLANG:
					var hreflang string
					hreflang = property.(string)
					link.HrefLang = &hreflang
				case PROFILE:
					var profile string
					profile = property.(string)
					link.Profile = &profile
				case TITLE:
					var title string
					title = property.(string)
					link.Title = &title
				case TYPE:
					var typeval string
					typeval = property.(string)
					link.Type = &typeval
				}
			}
			links = append(links, &link)
		}
		l.Relations[rel] = links
	}
	return nil
}

// Embeds holds embedded relations by reltype
type Embeds struct {
	Relations map[string][]Resource
}

// MarshalJSON marshals embeds
func (e *Embeds) MarshalJSON() ([]byte, error) {
	var bufferData []string
	for key, links := range e.Relations {
		jsonValue, err := json.Marshal(links)
		if err != nil {
			return nil, err
		}
		bufferData = append(bufferData, fmt.Sprintf("\"%s\": %s", key, string(jsonValue)))
	}
	buffer := bytes.NewBufferString("{")
	buffer.WriteString(strings.Join(bufferData, ","))
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshals embeds
func (e *Embeds) UnmarshalJSON(b []byte) error {
	var temp map[string]interface{}
	temp = make(map[string]interface{})
	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}
	e.Relations = make(map[string][]Resource)
	for k, v := range temp {
		var res []Resource
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		err = json.Unmarshal(b, &res)
		if err != nil {
			return err
		}
		e.Relations[k] = res
	}
	return nil
}

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
	var links *string
	var embeds *string
	if len(r.Links.Relations) > 0 || r.Links.Self != nil {
		b, err := json.Marshal(r.Links)
		if err != nil {
			return nil, err
		}
		linkString := fmt.Sprintf("\"_links\": %s", string(b))
		links = &linkString
	}

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

	// re marshal embedded and links
	embededjson, err := json.Marshal(temp["_embedded"])
	if err != nil {
		return err
	}
	embedded := Embeds{}
	links := Links{}
	err = json.Unmarshal(embededjson, &embedded)
	if err != nil {
		return err
	}
	r.Embeds = &embedded
	delete(temp, "_embedded")

	linksjson, err := json.Marshal(temp["_links"])
	if err != nil {
		return err
	}
	err = json.Unmarshal(linksjson, &links)
	if err != nil {
		return err
	}
	r.Links = &links

	delete(temp, "_links")

	r.Data = make(map[string]interface{})
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
