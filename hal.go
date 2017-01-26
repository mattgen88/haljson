package haljson

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Curie struct {
	Name      string `json:"name"`
	Href      string `json:"href"`
	Templated *bool  `json:"templated,omitempty"`
}

type Link struct {
	// If curied, special key nameing must be done
	Curried    *bool   `json:"-"`
	CurrieName *string `json:"-"`
	// When serializing to JSON we need to handle this specially
	Deprecation *string `json:"deprecation,omitempty"`
	Href        string  `json:"href,omitempty"`
	HrefLang    *string `json:"hreflang,omitempty"`
	Name        *string `json:"name,omitempty"`
	Profile     *string `json:"profile,omitempty"`
	Title       *string `json:"title,omitempty"`
	Type        *string `json:"type,omitempty"`
}

func (l *Link) MarshalJSON() ([]byte, error) {
	// if Curried = true, then prefix Name with CurrieName:
	r := reflect.ValueOf(l)

	name := ""
	if l.Curried != nil && *l.Curried {
		name = fmt.Sprint(l.CurrieName, ":", l.Name)
	}
	buffer := bytes.NewBufferString("{")
	for idx, field := range []string{"Deprecation", "Href", "HrefLang", "Name", "Profile", "Title", "Type"} {
		if !reflect.Indirect(r).FieldByName(field).IsNil() {
			if idx != 0 {
				buffer.WriteString(",")
			}
			if field == "Name" {
				buffer.WriteString(fmt.Sprintf("\"%s\":\"%s\"", strings.ToLower(field), name))
			} else {
				buffer.WriteString(fmt.Sprintf("\"%s\":\"%s\"", strings.ToLower(field), reflect.Indirect(r).FieldByName(field).String()))
			}
		}
	}
	buffer.WriteString("")
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

func (l *Link) UnmarshalJSON(b []byte) error {
	return nil
}

type Links struct {
	Curies *[]Curie `json:"curies,omitempty"`
	// When serializing to JSON we need to handle this specially
	Relations map[string][]Link
}

func (l *Links) MarshalJSON() ([]byte, error) {
	return nil, nil
}
func (l *Links) UnmarshalJSON(b []byte) error {
	return nil
}

func (l *Links) AddLink(reltype string, link Link) error {
	if _, ok := l.Relations[reltype]; !ok {
		l.Relations[reltype] = []Link{}
	}
	l.Relations[reltype] = append(l.Relations[reltype], link)
	return nil
}

func (l *Links) AddCurie(curie Curie) error {
	return nil
}

type Embeds struct {
	Relations map[string][]Resource
}

func (e *Embeds) MarshalJSON() ([]byte, error) {
	return nil, nil
}

func (e *Embeds) UnmarshalJSON(b []byte) error {
	return nil
}

func (e *Embeds) AddResource(reltype string, r Resource) error {
	if _, ok := e.Relations[reltype]; !ok {
		e.Relations[reltype] = []Resource{}
	}
	e.Relations[reltype] = append(e.Relations[reltype], r)
	return nil
}

type Resource struct {
	Links  *Links  `json:"_links,omitempty"`
	Embeds *Embeds `json:"_embedded,omitempty"`
	// When serializing to JSON we need to handle this specially
	Data map[string]interface{}
}

func (r *Resource) Self(uri string) error {
	if _, ok := r.Links.Relations["self"]; ok {
		return errors.New("Self relation already exists!")
	}
	r.Links.AddLink("self", Link{Href: uri})
	return nil
}

func NewResource() *Resource {
	r := &Resource{}
	r.Links = NewLinks()
	r.Embeds = NewEmbeds()
	return r
}

func NewLinks() *Links {
	l := &Links{}
	l.Relations = make(map[string][]Link)
	return l
}

func NewEmbeds() *Embeds {
	e := &Embeds{}
	e.Relations = make(map[string][]Resource)
	return e
}
