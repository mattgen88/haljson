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

func TestLinkChainableMethods(t *testing.T) {
	link := &Link{}

	// Test chaining
	result := link.
		SetHref("/api/users").
		SetTitle("User List").
		SetDeprecation("http://deprecated.example.com").
		SetHrefLang("en-US").
		SetProfile("http://profile.example.com").
		SetTemplated(true).
		SetType("application/json").
		SetName("user-list")

	assert.Equal(t, "/api/users", link.Href)
	assert.Equal(t, "User List", link.Title)
	assert.Equal(t, "http://deprecated.example.com", link.Deprecation)
	assert.Equal(t, "en-US", link.HrefLang)
	assert.Equal(t, "http://profile.example.com", link.Profile)
	assert.True(t, link.Templated)
	assert.Equal(t, "application/json", link.Type)
	assert.Equal(t, "user-list", link.Name)
	assert.Equal(t, link, result, "methods should return same instance for chaining")
}

func TestLinksUnmarshalErrors(t *testing.T) {
	// Test invalid curies format
	invalidCuries := `{"curies": "not-an-array"}`
	var links Links
	err := json.Unmarshal([]byte(invalidCuries), &links)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid curies format")

	// Test invalid curie item format
	invalidCurieItem := `{"curies": [123]}`
	var links2 Links
	err = json.Unmarshal([]byte(invalidCurieItem), &links2)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid curie format")

	// Test invalid self format
	invalidSelf := `{"self": "not-an-object"}`
	var links3 Links
	err = json.Unmarshal([]byte(invalidSelf), &links3)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid self link format")

	// Test invalid relation links format
	invalidRelation := `{"foo": "not-an-array"}`
	var links4 Links
	err = json.Unmarshal([]byte(invalidRelation), &links4)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid links format")

	// Test invalid link item format
	invalidLinkItem := `{"foo": [123]}`
	var links5 Links
	err = json.Unmarshal([]byte(invalidLinkItem), &links5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid link format")
}

func TestLinksAddCurieWithNilCuries(t *testing.T) {
	links := NewLinks()
	links.Curies = nil
	err := links.AddCurie(&Curie{Name: "test", Href: "/test"})
	assert.Nil(t, err)
	assert.NotNil(t, links.Curies)
	assert.Len(t, links.Curies, 1)
}

func TestLinksAddLinkWithColonPrefix(t *testing.T) {
	// Test that relation types starting with ":" don't trigger curie check
	links := NewLinks()
	err := links.AddLink(":special", &Link{Href: "/special"})
	assert.Nil(t, err)
	assert.Len(t, links.Relations[":special"], 1)
}

func TestSelfLinkWithAllProperties(t *testing.T) {
	// Test that self links can have all Link properties, not just href
	jsonData := `{
		"self": {
			"href": "/api/users",
			"title": "User List",
			"type": "application/hal+json",
			"deprecation": "http://deprecated.example.com",
			"hreflang": "en-US",
			"profile": "http://profile.example.com",
			"templated": true,
			"name": "user-list"
		}
	}`

	var links Links
	err := json.Unmarshal([]byte(jsonData), &links)
	assert.Nil(t, err)
	assert.NotNil(t, links.Self)
	assert.Equal(t, "/api/users", links.Self.Href)
	assert.Equal(t, "User List", links.Self.Title)
	assert.Equal(t, "application/hal+json", links.Self.Type)
	assert.Equal(t, "http://deprecated.example.com", links.Self.Deprecation)
	assert.Equal(t, "en-US", links.Self.HrefLang)
	assert.Equal(t, "http://profile.example.com", links.Self.Profile)
	assert.True(t, links.Self.Templated)
	assert.Equal(t, "user-list", links.Self.Name)
}

func TestLinksMarshalErrors(t *testing.T) {
	// Note: Link marshaling errors are hard to trigger as Link struct only contains
	// basic types (string, bool) which always marshal successfully.
	// Testing that marshal works correctly with various combinations.
	
	links := NewLinks()
	links.Self = &Link{Href: "/"}
	links.Curies = []Curie{{Name: "doc", Href: "/docs/{rel}", Templated: true}}
	links.Relations["test"] = []*Link{{Href: "/test"}}
	
	b, err := json.Marshal(links)
	assert.Nil(t, err)
	assert.Contains(t, string(b), `"self"`)
	assert.Contains(t, string(b), `"curies"`)
	assert.Contains(t, string(b), `"test"`)
}

func TestLinkNameInRelation(t *testing.T) {
	// Test that link "name" property is properly unmarshaled in relation links
	jsonData := `{
		"items": [
			{
				"href": "/api/items",
				"name": "item-link"
			}
		]
	}`

	var links Links
	err := json.Unmarshal([]byte(jsonData), &links)
	assert.Nil(t, err)
	assert.Len(t, links.Relations["items"], 1)
	assert.Equal(t, "/api/items", links.Relations["items"][0].Href)
	assert.Equal(t, "item-link", links.Relations["items"][0].Name)
}
