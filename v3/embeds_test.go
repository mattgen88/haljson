package haljson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmbeds(t *testing.T) {
	assert.Equal(t, &Embeds{Relations: make(map[string][]Resource[any])}, NewEmbeds(), "Embeds initialized incorrectly")
}
