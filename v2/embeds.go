package haljson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

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

// NewEmbeds creates and initializes Embeds
func NewEmbeds() *Embeds {
	e := &Embeds{}
	e.Relations = make(map[string][]Resource)
	return e
}
