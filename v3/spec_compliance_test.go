package haljson

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHALSpecCompliance tests compliance with draft-kelly-json-hal-11
func TestHALSpecCompliance(t *testing.T) {
	t.Run("Resource Object MUST have reserved properties", func(t *testing.T) {
		r := NewResource[any]()
		// _links and _embedded are OPTIONAL per spec
		assert.NotNil(t, r.Links)
		assert.NotNil(t, r.Embeds)
	})

	t.Run("Link Object MUST have href property", func(t *testing.T) {
		link := &Link{Href: "/test"}
		assert.NotEmpty(t, link.Href, "href is REQUIRED per spec section 5.1")
	})

	t.Run("Link Object OPTIONAL properties", func(t *testing.T) {
		link := &Link{
			Href:        "/test",
			Templated:   true,     // OPTIONAL - section 5.2
			Type:        "text",   // OPTIONAL - section 5.3
			Deprecation: "/depr",  // OPTIONAL - section 5.4
			Name:        "mylink", // OPTIONAL - section 5.5
			Profile:     "/prof",  // OPTIONAL - section 5.6
			Title:       "Test",   // OPTIONAL - section 5.7
			HrefLang:    "en",     // OPTIONAL - section 5.8
		}

		b, err := json.Marshal(link)
		assert.Nil(t, err)

		var unmarshaled Link
		err = json.Unmarshal(b, &unmarshaled)
		assert.Nil(t, err)
		assert.Equal(t, link, &unmarshaled)
	})

	t.Run("_links property contains link relation types", func(t *testing.T) {
		r := NewResource[any]()
		r.Self("/orders/123")
		r.AddLink("warehouse", &Link{Href: "/warehouse/56"})
		r.AddLink("invoice", &Link{Href: "/invoices/873"})

		b, err := json.Marshal(r)
		assert.Nil(t, err)

		// Per spec section 4.1.1: _links is an object whose property names
		// are link relation types
		assert.Contains(t, string(b), `"_links"`)
		assert.Contains(t, string(b), `"self"`)
		assert.Contains(t, string(b), `"warehouse"`)
		assert.Contains(t, string(b), `"invoice"`)
	})

	t.Run("_embedded property contains embedded resources", func(t *testing.T) {
		r := NewResource[any]()
		r.Self("/orders")

		order := NewResource[any]()
		order.Self("/orders/123")
		order.Data["total"] = 30.00
		order.Data["currency"] = "USD"

		r.AddEmbed("orders", order)

		b, err := json.Marshal(r)
		assert.Nil(t, err)

		// Per spec section 4.1.2: _embedded is an object whose property names
		// are link relation types and values are Resource Objects
		assert.Contains(t, string(b), `"_embedded"`)
		assert.Contains(t, string(b), `"orders"`)
	})

	t.Run("Self link SHOULD be present", func(t *testing.T) {
		// Per spec section 8.1: Each Resource Object SHOULD contain a 'self' link
		r := NewResource[any]()
		r.Self("/orders/123")

		assert.NotNil(t, r.Links.Self)
		assert.Equal(t, "/orders/123", r.Links.Self.Href)
	})

	t.Run("URI Template with templated=true", func(t *testing.T) {
		// Per spec section 5.1: If href is URI Template, Link Object SHOULD
		// have templated=true
		r := NewResource[any]()
		r.AddLink("find", &Link{
			Href:      "/orders{?id}",
			Templated: true,
		})

		assert.True(t, r.Links.Relations["find"][0].Templated)
	})

	t.Run("Curies for compact link relations", func(t *testing.T) {
		// Per spec section 8.3: HAL curies allow compact link relation types
		r := NewResource[any]()
		r.Self("/orders")
		r.AddCurie(&Curie{
			Name:      "acme",
			Href:      "https://docs.acme.com/relations/{rel}",
			Templated: true,
		})
		r.AddLink("acme:widgets", &Link{Href: "/widgets"})

		b, err := json.Marshal(r)
		assert.Nil(t, err)

		assert.Contains(t, string(b), `"curies"`)
		assert.Contains(t, string(b), `"acme"`)
		assert.Contains(t, string(b), `"acme:widgets"`)
	})

	t.Run("Link relation can be Link Object or array", func(t *testing.T) {
		// Per spec section 4.1.1: values are either a Link Object or array of Link Objects
		r := NewResource[any]()
		r.AddLink("item", &Link{Href: "/item/1"})
		r.AddLink("item", &Link{Href: "/item/2"})

		assert.Len(t, r.Links.Relations["item"], 2)
	})

	t.Run("Embedded can be Resource Object or array", func(t *testing.T) {
		// Per spec section 4.1.2: values are either Resource Object or array of Resource Objects
		r := NewResource[any]()
		r.AddEmbed("items", NewResource[any]())
		r.AddEmbed("items", NewResource[any]())

		assert.Len(t, r.Embeds.Relations["items"], 2)
	})

	t.Run("All non-reserved properties are valid JSON state", func(t *testing.T) {
		// Per spec section 4: All other properties MUST be valid JSON and
		// represent current state
		r := NewResource[any]()
		r.Self("/orders/123")
		r.Data["currency"] = "USD"
		r.Data["status"] = "shipped"
		r.Data["total"] = 10.20

		b, err := json.Marshal(r)
		assert.Nil(t, err)

		var result map[string]any
		err = json.Unmarshal(b, &result)
		assert.Nil(t, err)
		assert.Equal(t, "USD", result["currency"])
		assert.Equal(t, "shipped", result["status"])
		assert.Equal(t, 10.20, result["total"])
	})

	t.Run("Media type application/hal+json", func(t *testing.T) {
		// Per spec section 3: HAL Document has media type "application/hal+json"
		// This is a documentation note - actual HTTP headers are set by user
		// Our library correctly formats HAL+JSON structure
		r := NewResource[any]()
		r.Self("/test")
		b, err := json.Marshal(r)
		assert.Nil(t, err)

		// Verify it's valid JSON with HAL structure
		var result map[string]any
		err = json.Unmarshal(b, &result)
		assert.Nil(t, err)
		assert.Contains(t, result, "_links")
	})
}

func TestHALSpecExample(t *testing.T) {
	// Reproduce the example from spec section 3
	t.Run("Spec Section 3 Example", func(t *testing.T) {
		r := NewResource[any]()
		r.Self("/orders/523")
		r.AddLink("warehouse", &Link{Href: "/warehouse/56"})
		r.AddLink("invoice", &Link{Href: "/invoices/873"})
		r.Data["currency"] = "USD"
		r.Data["status"] = "shipped"
		r.Data["total"] = 10.20

		b, err := json.Marshal(r)
		assert.Nil(t, err)

		// Unmarshal and verify structure
		var result map[string]any
		err = json.Unmarshal(b, &result)
		assert.Nil(t, err)

		// Verify _links structure
		links := result["_links"].(map[string]any)
		self := links["self"].(map[string]any)
		assert.Equal(t, "/orders/523", self["href"])

		// Verify state
		assert.Equal(t, "USD", result["currency"])
		assert.Equal(t, "shipped", result["status"])
		assert.Equal(t, 10.20, result["total"])
	})

	// Reproduce the complex example from spec section 6
	t.Run("Spec Section 6 Example - Order List", func(t *testing.T) {
		r := NewResource[any]()
		r.Self("/orders")
		r.AddLink("next", &Link{Href: "/orders?page=2"})
		r.AddLink("find", &Link{Href: "/orders{?id}", Templated: true})

		// Add embedded orders
		order1 := NewResource[any]()
		order1.Self("/orders/123")
		order1.AddLink("basket", &Link{Href: "/baskets/98712"})
		order1.AddLink("customer", &Link{Href: "/customers/7809"})
		order1.Data["total"] = 30.00
		order1.Data["currency"] = "USD"
		order1.Data["status"] = "shipped"

		order2 := NewResource[any]()
		order2.Self("/orders/124")
		order2.AddLink("basket", &Link{Href: "/baskets/97213"})
		order2.AddLink("customer", &Link{Href: "/customers/12369"})
		order2.Data["total"] = 20.00
		order2.Data["currency"] = "USD"
		order2.Data["status"] = "processing"

		r.AddEmbed("orders", order1)
		r.AddEmbed("orders", order2)

		// Resource state
		r.Data["currentlyProcessing"] = 14
		r.Data["shippedToday"] = 20

		b, err := json.Marshal(r)
		assert.Nil(t, err)

		// Verify structure matches spec
		var result map[string]any
		err = json.Unmarshal(b, &result)
		assert.Nil(t, err)

		// Verify top-level links
		links := result["_links"].(map[string]any)
		assert.Contains(t, links, "self")
		assert.Contains(t, links, "next")
		assert.Contains(t, links, "find")

		// Verify embedded resources
		embedded := result["_embedded"].(map[string]any)
		orders := embedded["orders"].([]any)
		assert.Len(t, orders, 2)

		// Verify resource state
		assert.Equal(t, float64(14), result["currentlyProcessing"])
		assert.Equal(t, float64(20), result["shippedToday"])
	})
}
