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
	Deprecation string `json:"deprecation,omitempty"`
	Href        string `json:"href,omitempty"`
	HrefLang    string `json:"hreflang,omitempty"`
	Name        string `json:"name,omitempty"`
	Profile     string `json:"profile,omitempty"`
	Templated   bool   `json:"templated,omitempty"`
	Title       string `json:"title,omitempty"`
	Type        string `json:"type,omitempty"`
}

// Links is a container of Link, mapped by relation, and contains Curies
type Links struct {
	Self   *Link   `json:"-"`
	Curies []Curie `json:"curies,omitempty"`
	// When serializing to JSON we need to handle this specially
	Relations map[string][]*Link
}

// SetTitle sets the title, chainable
func (l *Link) SetTitle(title string) *Link {
	l.Title = title
	return l
}

// SetDeprecation sets deprecation, chainable
func (l *Link) SetDeprecation(deprecation string) *Link {
	l.Deprecation = deprecation
	return l
}

// SetDeprication sets deprecation, chainable
// Deprecated: Use SetDeprecation instead (typo fix)
func (l *Link) SetDeprication(deprecation string) *Link {
	return l.SetDeprecation(deprecation)
}

// SetHref sets href, chainable
func (l *Link) SetHref(href string) *Link {
	l.Href = href
	return l
}

// SetHrefLang sets hreflang, chainable
func (l *Link) SetHrefLang(lang string) *Link {
	l.HrefLang = lang
	return l
}

// SetProfile sets profile, chainable
func (l *Link) SetProfile(profile string) *Link {
	l.Profile = profile
	return l
}

// SetTemplated sets templated, chainable
func (l *Link) SetTemplated(templated bool) *Link {
	l.Templated = templated
	return l
}

// SetType sets type, chainable
func (l *Link) SetType(linkType string) *Link {
	l.Type = linkType
	return l
}

// SetName sets name, chainable
func (l *Link) SetName(name string) *Link {
	l.Name = name
	return l
}

// AddCurie adds a curie to the links
func (l *Links) AddCurie(curie *Curie) error {
	if l.Curies == nil {
		l.Curies = []Curie{}
	}
	l.Curies = append(l.Curies, *curie)
	return nil
}

// AddLink adds a link to reltype
func (l *Links) AddLink(reltype string, link *Link) error {
	// Check if curied and that if curied, curie exists
	// Note: we check > 0 to exclude relation types starting with ":"
	if strings.Index(reltype, ":") > 0 {
		parts := strings.Split(reltype, ":")
		if l.Curies == nil {
			return ErrNoCurie
		}
		curieExists := false
		for _, curie := range l.Curies {
			if parts[0] == curie.Name {
				curieExists = true
				break
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
	var bufferData []string
	if l.Self != nil {
		jsonValue, err := json.Marshal(l.Self)
		if err != nil {
			return nil, err
		}
		bufferData = append(bufferData, fmt.Sprintf("\"%s\": %s", SELF, string(jsonValue)))
	}
	if l.Curies != nil && len(l.Curies) > 0 {
		jsonValue, err := json.Marshal(l.Curies)
		if err != nil {
			return nil, err
		}
		bufferData = append(bufferData, fmt.Sprintf("\"%s\": %s", CURIES, string(jsonValue)))
	}

	// Sort keys for data
	sortedKeys := make([]string, 0, len(l.Relations))
	for k := range l.Relations {
		sortedKeys = append(sortedKeys, k)
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
	temp := make(map[string]any)
	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}
	if _, ok := temp[CURIES]; ok {
		var mycuries []Curie
		curiesArray, ok := temp[CURIES].([]any)
		if !ok {
			return fmt.Errorf("invalid curies format: expected array")
		}
		for _, curiesItem := range curiesArray {
			curiesMap, ok := curiesItem.(map[string]any)
			if !ok {
				return fmt.Errorf("invalid curie format: expected object")
			}
			var curie Curie
			for k, v := range curiesMap {
				switch k {
				case NAME:
					if name, ok := v.(string); ok {
						curie.Name = name
					}
				case HREF:
					if href, ok := v.(string); ok {
						curie.Href = href
					}
				case TEMPLATED:
					if templated, ok := v.(bool); ok {
						curie.Templated = templated
					}
				}
			}
			mycuries = append(mycuries, curie)
		}
		l.Curies = mycuries
		delete(temp, CURIES)
	}

	var self Link
	if _, ok := temp[SELF]; ok {
		selfMap, ok := temp[SELF].(map[string]any)
		if !ok {
			return fmt.Errorf("invalid self link format: expected object")
		}
		for k, v := range selfMap {
			switch k {
			case HREF:
				if href, ok := v.(string); ok {
					self.Href = href
				}
			}
		}
		l.Self = &self
		delete(temp, SELF)
	}

	l.Relations = make(map[string][]*Link)
	for rel, v := range temp {
		linksArray, ok := v.([]any)
		if !ok {
			return fmt.Errorf("invalid links format for relation %q: expected array", rel)
		}
		var links []*Link
		for _, linkItem := range linksArray {
			properties, ok := linkItem.(map[string]any)
			if !ok {
				return fmt.Errorf("invalid link format: expected object")
			}
			var link Link
			for key, property := range properties {
				switch key {
				case HREF:
					if href, ok := property.(string); ok {
						link.Href = href
					}
				case DEPRECATION:
					if deprecation, ok := property.(string); ok {
						link.Deprecation = deprecation
					}
				case HREFLANG:
					if hreflang, ok := property.(string); ok {
						link.HrefLang = hreflang
					}
				case NAME:
					if name, ok := property.(string); ok {
						link.Name = name
					}
				case PROFILE:
					if profile, ok := property.(string); ok {
						link.Profile = profile
					}
				case TITLE:
					if title, ok := property.(string); ok {
						link.Title = title
					}
				case TYPE:
					if typeval, ok := property.(string); ok {
						link.Type = typeval
					}
				case TEMPLATED:
					if templated, ok := property.(bool); ok {
						link.Templated = templated
					}
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
	return &Links{
		Relations: make(map[string][]*Link),
	}
}
