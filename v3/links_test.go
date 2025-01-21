package haljson

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinks(t *testing.T) {
	assert.Equal(t, &Links{Relations: make(map[string][]*Link)}, NewLinks(), "Embeds initialized incorrectly")
}

func TestLinksMarshal(t *testing.T) {
	r := NewResource[any]()
	r.Self("/")
	r.AddLink("foo", &Link{Href: "/foo"})
	deprecation := "string url"
	title := "bar item"
	typeval := "application/hal+json"
	hreflang := "en_US"
	profile := "string uri"
	err := r.AddLink("bar:baz", &Link{Href: "/bar/{item}", Templated: true, Title: title, Deprecation: deprecation, Type: typeval, HrefLang: hreflang, Profile: profile})
	assert.NotNil(t, err, "expected an error from adding link")
	assert.Equal(t, err, ErrNoCurie, "Expected ErrNoCurie to be returned")
	r.AddCurie(&Curie{Name: "bar", Templated: true, Href: "/docs/bar"})
	err = r.AddLink("bar:baz", &Link{Href: "/bar/{item}", Templated: true, Title: title, Deprecation: deprecation, Type: typeval, HrefLang: hreflang, Profile: profile})
	assert.Nil(t, err, "expected no error from adding link")

	l := r.Links

	b, err := json.MarshalIndent(l, "", "\t")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, `{
	"self": {
		"href": "/"
	},
	"curies": [
		{
			"name": "bar",
			"href": "/docs/bar",
			"templated": true
		}
	],
	"bar:baz": [
		{
			"deprecation": "string url",
			"href": "/bar/{item}",
			"hreflang": "en_US",
			"profile": "string uri",
			"templated": true,
			"title": "bar item",
			"type": "application/hal+json"
		}
	],
	"foo": [
		{
			"href": "/foo"
		}
	]
}`, string(b), "Links marshalled incorrectly")
}

func TestLinksUnmarshal(t *testing.T) {
	marshalled := `{
	"self": {
		"href": "/"
	},
	"curies": [
		{
			"name": "bar",
			"href": "/docs/bar",
			"templated": true
		}
	],
	"bar:baz": [
		{
			"deprecation": "string url",
			"href": "/bar/{item}",
			"hreflang": "en_US",
			"profile": "string uri",
			"templated": true,
			"title": "bar item",
			"type": "application/hal+json"
		}
	],
	"foo": [
		{
			"href": "/foo"
		}
	]
}`
	var inflated Links
	err := json.Unmarshal([]byte(marshalled), &inflated)
	assert.Nil(t, err, "Expected no error unmarshalling")

	r := NewResource[any]()
	r.Self("/")
	r.AddLink("foo", &Link{Href: "/foo"})
	r.AddCurie(&Curie{Name: "bar", Templated: true, Href: "/docs/bar"})
	deprecation := "string url"
	title := "bar item"
	typeval := "application/hal+json"
	hreflang := "en_US"
	profile := "string uri"
	r.AddLink("bar:baz", &Link{Href: "/bar/{item}", Templated: true, Title: title, Deprecation: deprecation, Type: typeval, HrefLang: hreflang, Profile: profile})

	assert.Equal(t, r.Links.Curies, inflated.Curies, "Links curies unmarshalled incorrectly")
	assert.Equal(t, r.Links.Relations, inflated.Relations, "Links relations unmarshalled incorrectly")

}

func TestAddLinkBeforeCurie(t *testing.T) {
	r := NewResource[any]()
	err := r.AddLink("foo:bar", &Link{Href: "/foo"})
	assert.NotNil(t, err)

	r.AddCurie(&Curie{Href: "/docs/bar/{rel}", Name: "bar"})
	err2 := r.AddLink("foo:bar", &Link{Href: "/foo"})
	assert.NotNil(t, err2)
	assert.Equal(t, err, err2, "Same errors for link before curie")
}
