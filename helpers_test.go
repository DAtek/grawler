package grawler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinPath(t *testing.T) {
	t.Run("juoins path parts without parsing and escaping them", func(t *testing.T) {
		baseUrl := "http://demo.com/"
		path1 := "/regional/"
		path2 := "oesterreich/grundstueck-kaufen/seite-2?primaryAreaFrom=1500&primaryPriceTo=200000"

		assert.Equal(
			t,
			baseUrl[:(len(baseUrl)-1)]+path1+path2,
			joinPath(baseUrl, path1, path2),
		)
	})

	t.Run("Returns single element", func(t *testing.T) {
		baseUrl := "http://demo.com/"
		assert.Equal(t, baseUrl, joinPath(baseUrl))
	})

}
