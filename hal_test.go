package haljson

import "testing"
import "github.com/stretchr/testify/assert"
import "encoding/json"

func TestResource(t *testing.T) {
	r := NewResource()
	assert.Equal(t, &Resource{Links: NewLinks(), Embeds: NewEmbeds(), Data: make(map[string]interface{})}, r, "Resource initialized correctly")
}

func TestEmbeds(t *testing.T) {
	assert.Equal(t, &Embeds{Relations: make(map[string][]Resource)}, NewEmbeds(), "Embeds initialized correctly")
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
	"bar": "baz"
}`, string(b), "marshalled resource did not match")
}

func TestResourceSelf(t *testing.T) {
	r := NewResource()
	r.Self("/")

	assert.Equal(t, &Link{Href: "/"}, r.Links.Self, "Link not set correctly")
}

func TestLinks(t *testing.T) {
	assert.Equal(t, &Links{Relations: make(map[string][]*Link)}, NewLinks(), "Embeds initialized correctly")
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
}`, string(b), "Links marshalled correctly")
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
