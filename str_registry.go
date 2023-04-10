package grawler

import (
	mapset "github.com/deckarep/golang-set/v2"
)

type stringRegistry struct {
	analyzedUrls mapset.Set[string]
}

func newStringRegistry() *stringRegistry {
	return &stringRegistry{
		analyzedUrls: mapset.NewSet[string](),
	}
}

func (c *stringRegistry) add(item string) {
	c.analyzedUrls.Add(item)
}

func (c *stringRegistry) getNew(items []string) []string {
	allUrls := mapset.NewSet(items...)
	newUrls := allUrls.Difference(c.analyzedUrls)
	result := []string{}

	newUrls.Each(func(s string) bool {
		result = append(result, s)
		return false
	})

	return result
}
