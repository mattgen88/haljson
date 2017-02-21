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
	r := NewResource()
	r.Self("/")
	r.AddLink("foo", &Link{Href: "/foo"})

	l := r.Links

	b, err := json.MarshalIndent(l, "", "\t")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, `{
	"self": {
		"href": "/"
	},
	"foo": [
		{
			"href": "/foo"
		}
	]
}`, string(b), "Links marshalled incorrectly")
}

func TestAddLinkBeforeCurie(t *testing.T) {
	r := NewResource()
	err := r.AddLink("foo:bar", &Link{Href: "/foo"})
	assert.NotNil(t, err)

	r.AddCurie(&Curie{Href: "/docs/bar/{rel}", Name: "bar"})
	err2 := r.AddLink("foo:bar", &Link{Href: "/foo"})
	assert.NotNil(t, err2)
	assert.Equal(t, err, err2, "Same errors for link before curie")
}
