package haljson

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResource(t *testing.T) {
	r := NewResource()
	assert.Equal(t, &Resource{Links: NewLinks(), Embeds: NewEmbeds(), Data: make(map[string]interface{})}, r, "Resource initialized incorrectly")
}

func TestResourceMarshal(t *testing.T) {
	r := NewResource()
	r.Self("/")
	r.Data["bar"] = "baz"

	rEmbed := NewResource()
	rEmbed.Self("/foo")

	r.AddCurie(&Curie{Href: "/docs/bar/{rel}", Templated: true, Name: "bar"})

	r.AddLink("bar:foo", &Link{Href: "/bar/foo"})

	r.AddEmbed("foo", rEmbed)

	b, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, `{
	"_links": {
		"self": {
			"href": "/"
		},
		"curies": [
			{
				"name": "bar",
				"href": "/docs/bar/{rel}",
				"templated": true
			}
		],
		"bar:foo": [
			{
				"href": "/bar/foo"
			}
		]
	},
	"_embedded": {
		"foo": [
			{
				"_links": {
					"self": {
						"href": "/foo"
					}
				}
			}
		]
	},
	"bar": "baz"
}`, string(b), "marshalled resource did not match")
}

func TestResourceUnmarshal(t *testing.T) {
	t.Skip()
	marshaled := `{
	"foo": {"lel": "lawl"},
	"bar": "baz"
}`

	r := &Resource{}
	r.Data = make(map[string]interface{})
	r.Data["bar"] = "baz"
	var foo map[string]interface{}
	foo = make(map[string]interface{})
	foo["lel"] = "lawl"
	r.Data["foo"] = foo

	var inflated Resource
	err := json.Unmarshal([]byte(marshaled), &inflated)
	assert.Nil(t, err, "error in unmarshal")
	assert.Equal(t, r.Data, inflated.Data, "data was not the same")
	assert.Equal(t, r.Links, inflated.Links, "links was not the same")

	// Reflate and test equivalency still
	b, err := json.MarshalIndent(inflated, "", "\t")
	assert.Nil(t, err)
	assert.Equal(t, marshaled, string(b), "failed to marshal, unmarshal, marshal to the same thing")
}

func TestResourceSelf(t *testing.T) {
	r := NewResource()
	r.Self("/")

	assert.Equal(t, &Link{Href: "/"}, r.Links.Self, "Link not set correctly")
}
