package grawler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrRegistry(t *testing.T) {
	t.Run("Item is being cached", func(t *testing.T) {
		r := newStringRegistry()

		r.add("a")

		newItems := r.getNew([]string{"a", "b"})

		assert.Equal(t, []string{"b"}, newItems)
	})
}
