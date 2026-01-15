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

func TestResourceWithTypedData(t *testing.T) {
	type User struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	// Create resource with User as the data type
	r := NewResource[User]()
	r.Self("/users/123")
	
	// When T is User, Data is map[string]User
	// So we store User struct directly as the value
	r.Data["profile"] = User{Name: "John", Email: "john@example.com"}
	r.Data["settings"] = User{Name: "Admin", Email: "admin@example.com"}

	// Marshal to JSON
	b, err := json.Marshal(r)
	assert.Nil(t, err)
	assert.Contains(t, string(b), `"name":"John"`)
	assert.Contains(t, string(b), `"email":"john@example.com"`)

	// Test unmarshal with typed data
	var r2 Resource[User]
	err = json.Unmarshal(b, &r2)
	assert.Nil(t, err)
	
	// Verify the data was properly typed as User structs
	assert.Equal(t, "John", r2.Data["profile"].Name)
	assert.Equal(t, "john@example.com", r2.Data["profile"].Email)
	assert.Equal(t, "Admin", r2.Data["settings"].Name)
	assert.Equal(t, "admin@example.com", r2.Data["settings"].Email)
	assert.Equal(t, "/users/123", r2.Links.Self.Href)
}

func TestResourceUnmarshalErrors(t *testing.T) {
	// Test with invalid JSON
	var r Resource[any]
	err := json.Unmarshal([]byte(`invalid json`), &r)
	assert.NotNil(t, err)

	// Test with data that can't be converted to type
	type StrictType struct {
		ID int `json:"id"`
	}
	invalidData := `{"field": "not-a-number"}`
	var r2 Resource[StrictType]
	err = json.Unmarshal([]byte(invalidData), &r2)
	assert.NotNil(t, err)
}

func TestResourceMarshalWithNilLinks(t *testing.T) {
	r := &Resource[any]{
		Links:  nil,
		Embeds: NewEmbeds(),
		Data:   make(map[string]any),
	}
	r.Data["test"] = "value"

	b, err := json.Marshal(r)
	assert.Nil(t, err)
	assert.Contains(t, string(b), `"test":"value"`)
	assert.NotContains(t, string(b), `"_links"`)
}

func TestResourceMarshalWithEmptyLinks(t *testing.T) {
	r := NewResource[any]()
	r.Data["test"] = "value"

	b, err := json.Marshal(r)
	assert.Nil(t, err)
	assert.Contains(t, string(b), `"test":"value"`)
	// Empty links should not be marshaled
	assert.NotContains(t, string(b), `"_links"`)
}

func TestResourceMarshalWithNilEmbeds(t *testing.T) {
	r := &Resource[any]{
		Links:  NewLinks(),
		Embeds: nil,
		Data:   make(map[string]any),
	}
	r.Links.Self = &Link{Href: "/"}
	r.Data["test"] = "value"

	b, err := json.Marshal(r)
	assert.Nil(t, err)
	assert.Contains(t, string(b), `"test":"value"`)
	assert.NotContains(t, string(b), `"_embedded"`)
}

func TestResourceMarshalWithEmptyEmbeds(t *testing.T) {
	r := NewResource[any]()
	r.Self("/")
	r.Data["test"] = "value"

	b, err := json.Marshal(r)
	assert.Nil(t, err)
	assert.Contains(t, string(b), `"test":"value"`)
	// Empty embeds should not be marshaled
	assert.NotContains(t, string(b), `"_embedded"`)
}

func TestResourceUnmarshalWithCuries(t *testing.T) {
	jsonData := `{
		"_links": {
			"self": {"href": "/"},
			"curies": [
				{"name": "doc", "href": "/docs/{rel}", "templated": true}
			]
		},
		"test": "value"
	}`

	var r Resource[any]
	err := json.Unmarshal([]byte(jsonData), &r)
	assert.Nil(t, err)
	assert.Equal(t, "/", r.Links.Self.Href)
	assert.Len(t, r.Links.Curies, 1)
	assert.Equal(t, "doc", r.Links.Curies[0].Name)
	assert.Equal(t, "value", r.Data["test"])
}
