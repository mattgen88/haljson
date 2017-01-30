package haljson
import "testing"
import "github.com/stretchr/testify/assert"

func TestResource(t *testing.T) {
  r := NewResource()
  assert.Equal(t, &Resource{Links: NewLinks(), Embeds: NewEmbeds(), Data: make(map[string]interface{})}, r, "Resource initialized correctly")
}

func TestEmbeds(t *testing.T) {
  assert.Equal(t, &Embeds{Relations: make(map[string][]Resource)}, NewEmbeds(), "Embeds initialized correctly")
}

func TestEmbedsMarshal(t *testing.T) {
  t.Skip()
}

func TestEmbedsUnmarshal(t *testing.T) {
  t.Skip()
}

func TestLinks(t *testing.T) {
  assert.Equal(t, &Links{Relations: make(map[string][]*Link)}, NewLinks(), "Embeds initialized correctly")
}

func TestLinksMarshal(t *testing.T) {
  t.Skip()
}

func TestLinksUnmarshal(t *testing.T) {
  t.Skip()
}

func TestLink(t *testing.T) {
  t.Skip()
}
