package page_loader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newMockPageLoader() IPageLoader {
	return &MockPageLoader{}
}

func TestMockPageLoader(t *testing.T) {

	t.Run("Test LoadPage", func(t *testing.T) {
		pageLoader := newMockPageLoader().(*MockPageLoader)
		expectedContent := "content"

		pageLoader.LoadPage_ = func(url string) (string, error) {
			return expectedContent, nil
		}

		result, err := pageLoader.LoadPage("")

		assert.Nil(t, err)
		assert.Equal(t, expectedContent, result)
	})

}
