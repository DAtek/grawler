package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewMockCache() ICache {
	return &MockCache{}
}

func TestMockCache(t *testing.T) {
	t.Run("Test Get", func(t *testing.T) {
		c := NewMockCache().(*MockCache)

		expectedContent := "I want a big garden with lot of trees"
		c.Get_ = func(key string) (string, error) {
			return expectedContent, nil
		}

		res, err := c.Get("")

		assert.Nil(t, err)
		assert.Equal(t, expectedContent, res)
	})

	t.Run("Test Set", func(t *testing.T) {
		c := NewMockCache().(*MockCache)

		c.Set_ = func(key, val string) error {
			return nil
		}

		assert.Nil(t, c.Set("", ""))
	})

	t.Run("Test Has", func(t *testing.T) {
		c := NewMockCache().(*MockCache)

		c.Has_ = func(key string) bool {
			return true
		}

		assert.True(t, c.Has(""))
	})

	t.Run("Test Delete", func(t *testing.T) {
		c := NewMockCache().(*MockCache)

		c.Delete_ = func(key string) error {
			return nil
		}

		assert.Nil(t, c.Delete(""))
	})
}
