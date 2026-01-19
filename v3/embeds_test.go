package haljson

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmbeds(t *testing.T) {
	assert.Equal(t, &Embeds{Relations: make(map[string][]Resource[any])}, NewEmbeds(), "Embeds initialized incorrectly")
}

func TestEmbedsMarshalJSON(t *testing.T) {
	embeds := NewEmbeds()
	r1 := NewResource[any]()
	r1.Self("/foo")
	r1.Data["name"] = "test"

	embeds.Relations["items"] = []Resource[any]{*r1}

	b, err := json.Marshal(embeds)
	assert.Nil(t, err)
	assert.Contains(t, string(b), `"items"`)
	assert.Contains(t, string(b), `"name"`)
	assert.Contains(t, string(b), `"test"`)
}

func TestEmbedsUnmarshalJSON(t *testing.T) {
	jsonData := `{
		"items": [
			{
				"_links": {
					"self": {
						"href": "/foo"
					}
				},
				"name": "test"
			}
		]
	}`

	var embeds Embeds
	err := json.Unmarshal([]byte(jsonData), &embeds)
	assert.Nil(t, err)
	assert.Len(t, embeds.Relations["items"], 1)
	assert.Equal(t, "/foo", embeds.Relations["items"][0].Links.Self.Href)
	assert.Equal(t, "test", embeds.Relations["items"][0].Data["name"])
}

func TestEmbedsUnmarshalErrors(t *testing.T) {
	// Test with invalid JSON
	var embeds Embeds
	err := json.Unmarshal([]byte(`invalid json`), &embeds)
	assert.NotNil(t, err)
}
