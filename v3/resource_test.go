package haljson

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResource(t *testing.T) {
	r := NewResource[any]()
	assert.Equal(t, &Resource[any]{Links: NewLinks(), Embeds: NewEmbeds(), Data: make(map[string]any)}, r, "Resource initialized incorrectly")
}

func TestResourceMarshal(t *testing.T) {
	r := NewResource[any]()
	r.Self("/")
	r.Data["bar"] = "baz"

	rEmbed := NewResource[any]()
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
	marshaled := `{
	"_links": {
		"self": {
			"href": "/"
		}
	},
	"_embedded": {
		"foo": [
			{
				"_links": {
					"self": {
						"href": "/foo"
					}
				},
				"test": "foo"
			}
		]
	},
	"bar": "baz",
	"foo": {
		"lel": "lawl"
	}
}`

	r := NewResource[any]()
	r.Self("/")
	r.Data["bar"] = "baz"
	var foo = make(map[string]any)
	foo["lel"] = "lawl"
	r.Data["foo"] = foo
	rEmbed := NewResource[any]()
	rEmbed.Self("/foo")
	rEmbed.Data["test"] = "foo"
	r.AddEmbed("foo", rEmbed)

	var inflated Resource[any]
	err := json.Unmarshal([]byte(marshaled), &inflated)
	assert.Nil(t, err, "error in unmarshal")
	assert.Equal(t, r.Data, inflated.Data, "data was not the same")
	assert.Equal(t, r.Links, inflated.Links, "links was not the same")
	assert.Equal(t, r.Embeds, inflated.Embeds, "embeds was not the same")

	// Reflate and test equivalency still
	b, err := json.MarshalIndent(r, "", "\t")
	assert.Nil(t, err)
	assert.Equal(t, marshaled, string(b), "failed to marshal, unmarshal, marshal to the same thing")
}

func TestResourceSelf(t *testing.T) {
	r := NewResource[any]()
	r.Self("/")

	assert.Equal(t, &Link{Href: "/"}, r.Links.Self, "Link not set correctly")
}
