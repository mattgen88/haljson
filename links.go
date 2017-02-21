package haljson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Link represents a link
type Link struct {
	Deprecation *string `json:"deprecation,omitempty"`
	Href        string  `json:"href,omitempty"`
	HrefLang    *string `json:"hreflang,omitempty"`
	Profile     *string `json:"profile,omitempty"`
	Templated   *bool   `json:"templated,omitempty"`
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

// AddCurie adds a curie to the links
func (l *Links) AddCurie(curie *Curie) error {
	if l.Curies == nil {
		l.Curies = &[]Curie{}
	}
	*l.Curies = append(*l.Curies, *curie)
	return nil
}

// AddLink adds a link to reltype
func (l *Links) AddLink(reltype string, link *Link) error {
	// Check if curied and that if curied, curie exists
	curieExists := false
	if strings.Index(reltype, ":") > 0 {
		parts := strings.Split(reltype, ":")
		var curies *[]Curie
		curies = l.Curies
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
	if _, ok := l.Relations[reltype]; !ok {
		l.Relations[reltype] = []*Link{}
	}

	l.Relations[reltype] = append(l.Relations[reltype], link)
	return nil
}

// MarshalJSON to marshal Links properly
func (l *Links) MarshalJSON() ([]byte, error) {
	// @TODO sort keys
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

	// Sort keys for data
	var sortedKeys = make([]string, len(l.Relations))
	i := 0
	for k := range l.Relations {
		sortedKeys[i] = k
		i++
	}
	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		jsonValue, err := json.Marshal(l.Relations[key])
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
	var temp = make(map[string]interface{})
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
					link.Deprecation = &deprecation
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
				case TEMPLATED:
					var templated bool
					templated = property.(bool)
					link.Templated = &templated
				}
			}
			links = append(links, &link)
		}
		l.Relations[rel] = links
	}
	return nil
}

// NewLinks creates and initializes Links
func NewLinks() *Links {
	l := &Links{}
	l.Relations = make(map[string][]*Link)
	return l
}
