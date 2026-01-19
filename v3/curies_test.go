package haljson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurieChainableMethods(t *testing.T) {
	curie := &Curie{}

	// Test chaining
	result := curie.
		SetName("test").
		SetHref("/docs/{rel}").
		SetTemplated(true)

	assert.Equal(t, "test", curie.Name)
	assert.Equal(t, "/docs/{rel}", curie.Href)
	assert.True(t, curie.Templated)
	assert.Equal(t, curie, result, "methods should return same instance for chaining")
}

func TestCurieSetName(t *testing.T) {
	curie := &Curie{}
	curie.SetName("myname")
	assert.Equal(t, "myname", curie.Name)
}

func TestCurieSetHref(t *testing.T) {
	curie := &Curie{}
	curie.SetHref("/api/{rel}")
	assert.Equal(t, "/api/{rel}", curie.Href)
}

func TestCurieSetTemplated(t *testing.T) {
	curie := &Curie{}
	curie.SetTemplated(true)
	assert.True(t, curie.Templated)

	curie.SetTemplated(false)
	assert.False(t, curie.Templated)
}
