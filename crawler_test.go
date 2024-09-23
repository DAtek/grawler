package grawler

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/DAtek/grawler/cache"
	"github.com/DAtek/grawler/page_loader"

	"github.com/DAtek/gotils"
	"github.com/stretchr/testify/assert"
)

func TestStartCrawling(t *testing.T) {
	t.Run("Returns all found records and analyzes all URLs", func(t *testing.T) {
		timeout := gotils.NewTimeoutMs(1100)
		go func() { panic(<-timeout.ErrorCh) }()
		defer timeout.Cancel()

		articles := []*ExapleModel{
			{Title: "title1", Content: "content1"},
			{Title: "title2", Content: "content2"},
		}
		i := 0

		analyzer := &MockAnalyzer{
			GetModel_: func() *ExapleModel {
				if i >= 2 {
					return nil
				}
				article := articles[i]
				i++
				return article
			},
			GetUrls_: func() []string { return []string{"a", "b"} },
		}

		var createAnalyzer NewAnalyzer[ExapleModel] = func(html, u *string) (IAnalyzer[ExapleModel], error) {
			return analyzer, nil
		}

		crawledUrls := map[string]struct{}{}
		crawler := NewCrawler(
			&cache.MockCache{
				Has_: func(key string) bool {
					_, ok := crawledUrls[key]
					return ok
				},
				Set_: func(key string, val string) error {
					crawledUrls[key] = struct{}{}
					return nil
				},
				Get_: func(key string) (string, error) {
					_, ok := crawledUrls[key]
					if ok {
						return key, nil
					}

					return "", errors.New("KEY_NOT_FOUND_IN_CACHE")
				},
				Delete_: func(key string) error {
					return nil
				},
			},
			createAnalyzer,
			&page_loader.MockPageLoader{
				LoadPage_: func(url string) (string, error) {
					return "", nil
				},
			},
			gotils.NewLogger(gotils.LogLevelInfo, &bytes.Buffer{}, &bytes.Buffer{}),
			"http://demo.example",
			CrawlerConfig{
				PageAnalyzers: 2,
				PageLoaders:   3,
			},
		)

		result := []*ExapleModel{}
		for item := range crawler.Crawl("asd") {
			result = append(result, item)
			if len(result) == 2 {
				break
			}
		}

		assert.Equal(t, articles, result)

		crawler.Stop()
		crawler.WaitStopped()
	})

	t.Run("Test LoadPage loads page from cache", func(t *testing.T) {
		timeout := gotils.NewTimeoutMs(100)
		go func() { panic(<-timeout.ErrorCh) }()
		defer timeout.Cancel()
		articles := []*ExapleModel{
			{Title: "title1", Content: "content1"},
			{Title: "title2", Content: "content2"},
		}
		i := 0

		analyzer := &MockAnalyzer{
			GetModel_: func() *ExapleModel {
				if i >= 2 {
					return nil
				}
				article := articles[i]
				i++
				return article
			},
			GetUrls_: func() []string { return []string{"a", "b"} },
		}

		var createAnalyzer NewAnalyzer[ExapleModel] = func(html, u *string) (IAnalyzer[ExapleModel], error) {
			return analyzer, nil
		}

		url := "asd"
		crawledUrls := map[string]struct{}{
			url: {},
		}

		crawler_ := NewCrawler(
			&cache.MockCache{
				Has_: func(key string) bool {
					_, ok := crawledUrls[key]
					return ok
				},
			},
			createAnalyzer,
			&page_loader.MockPageLoader{
				LoadPage_: func(url string) (string, error) {
					return "", nil
				},
			},
			gotils.NewLogger(gotils.LogLevelInfo, &bytes.Buffer{}, &bytes.Buffer{}),
			"http://demo.example",
			CrawlerConfig{},
		).(*crawler[ExapleModel])

		remainingUrlCh := make(chan string, 1)
		remainingUrlCh <- url
		downloadedUrlCh := make(chan string, 1)

		assert.True(t, crawler_.LoadPage(remainingUrlCh, downloadedUrlCh, 1))

		assert.Equal(t, url, <-downloadedUrlCh)
	})

	t.Run("Test LoadPage logs error if downloading fails", func(t *testing.T) {
		timeout := gotils.NewTimeoutMs(100)
		go func() { panic(<-timeout.ErrorCh) }()
		defer timeout.Cancel()
		analyzer := &MockAnalyzer{}

		var createAnalyzer NewAnalyzer[ExapleModel] = func(html, u *string) (IAnalyzer[ExapleModel], error) {
			return analyzer, nil
		}

		err := errors.New("UNEXPECTED_ERROR")
		outBuf := &bytes.Buffer{}
		crawledUrls := map[string]struct{}{}
		crawler_ := NewCrawler(
			&cache.MockCache{
				Has_: func(key string) bool {
					_, ok := crawledUrls[key]
					return ok
				},
			},
			createAnalyzer,
			&page_loader.MockPageLoader{
				LoadPage_: func(url string) (string, error) {
					return "", err
				},
			},
			gotils.NewLogger(gotils.LogLevelInfo, outBuf, &bytes.Buffer{}),
			"http://demo.example",
			CrawlerConfig{},
		).(*crawler[ExapleModel])

		remainingUrlCh := make(chan string, 1)
		remainingUrlCh <- "asd"
		downloadedUrlCh := make(chan string, 1)

		assert.True(t, crawler_.LoadPage(remainingUrlCh, downloadedUrlCh, 1))
		assert.True(t, strings.Contains(outBuf.String(), err.Error()))
	})

	t.Run("Test LoadPage logs error if writing to cache fails", func(t *testing.T) {
		timeout := gotils.NewTimeoutMs(100)
		go func() { panic(<-timeout.ErrorCh) }()
		defer timeout.Cancel()
		analyzer := &MockAnalyzer{}

		var createAnalyzer NewAnalyzer[ExapleModel] = func(html, u *string) (IAnalyzer[ExapleModel], error) {
			return analyzer, nil
		}

		err := errors.New("UNEXPECTED_ERROR")
		outBuf := &bytes.Buffer{}
		crawledUrls := map[string]struct{}{}
		crawler_ := NewCrawler(
			&cache.MockCache{
				Has_: func(key string) bool {
					_, ok := crawledUrls[key]
					return ok
				},
				Set_: func(key, val string) error {
					return err
				},
			},
			createAnalyzer,
			&page_loader.MockPageLoader{
				LoadPage_: func(url string) (string, error) {
					return "", nil
				},
			},
			gotils.NewLogger(gotils.LogLevelInfo, outBuf, &bytes.Buffer{}),
			"http://demo.example",
			CrawlerConfig{},
		).(*crawler[ExapleModel])

		remainingUrlCh := make(chan string, 1)
		remainingUrlCh <- "asd"
		downloadedUrlCh := make(chan string, 1)

		assert.True(t, crawler_.LoadPage(remainingUrlCh, downloadedUrlCh, 1))
		assert.True(t, strings.Contains(outBuf.String(), err.Error()))
	})

	t.Run("Test AnalyzePage adds base URL to new url", func(t *testing.T) {
		timeout := gotils.NewTimeoutMs(100)
		go func() { panic(<-timeout.ErrorCh) }()
		defer timeout.Cancel()

		newUrl := "/1"
		analyzer := &MockAnalyzer{
			GetModel_: func() *ExapleModel {
				return nil
			},
			GetUrls_: func() []string {
				return []string{newUrl}
			},
		}

		var createAnalyzer NewAnalyzer[ExapleModel] = func(html, u *string) (IAnalyzer[ExapleModel], error) {
			return analyzer, nil
		}

		baseUrl := "http://demo.example"
		crawler_ := NewCrawler(
			&cache.MockCache{
				Get_: func(key string) (string, error) {
					return "", nil
				},
			},
			createAnalyzer,
			&page_loader.MockPageLoader{
				LoadPage_: func(url string) (string, error) {
					return "", nil
				},
			},
			gotils.NewLogger(gotils.LogLevelInfo, &bytes.Buffer{}, &bytes.Buffer{}),
			baseUrl,
			CrawlerConfig{},
		).(*crawler[ExapleModel])

		remainingUrlCh := make(chan string, 1)
		downloadedUrlCh := make(chan string, 1)

		downloadedUrlCh <- "asd"
		resultCh := make(chan *ExapleModel, 1)

		assert.True(t, crawler_.AnalyzePage(downloadedUrlCh, remainingUrlCh, resultCh, 1))
		remainingUrl := <-remainingUrlCh
		assert.Equal(t, baseUrl+newUrl, remainingUrl)
	})

	t.Run("Test AnalyzePage logs error if creating the analyzer fails", func(t *testing.T) {
		timeout := gotils.NewTimeoutMs(100)
		go func() { panic(<-timeout.ErrorCh) }()
		defer timeout.Cancel()

		err := errors.New("UNEXPECTED_ERROR")
		var createAnalyzer NewAnalyzer[ExapleModel] = func(html, u *string) (IAnalyzer[ExapleModel], error) {
			return nil, err
		}

		outBuf := &bytes.Buffer{}
		crawler_ := NewCrawler(
			&cache.MockCache{
				Get_: func(key string) (string, error) {
					return "", nil
				},
			},
			createAnalyzer,
			&page_loader.MockPageLoader{
				LoadPage_: func(url string) (string, error) {
					return "", nil
				},
			},
			gotils.NewLogger(gotils.LogLevelInfo, outBuf, &bytes.Buffer{}),
			"http://demo.example",
			CrawlerConfig{},
		).(*crawler[ExapleModel])

		remainingUrlCh := make(chan string, 1)
		downloadedUrlCh := make(chan string, 1)
		downloadedUrlCh <- "asd"
		resultCh := make(chan *ExapleModel, 1)

		assert.True(t, crawler_.AnalyzePage(downloadedUrlCh, remainingUrlCh, resultCh, 1))
		assert.True(t, strings.Contains(outBuf.String(), err.Error()))
	})

	t.Run("Test AnalyzePage logs error if getting item from cache fails", func(t *testing.T) {
		timeout := gotils.NewTimeoutMs(100)
		go func() { panic(<-timeout.ErrorCh) }()
		defer timeout.Cancel()
		analyzer := &MockAnalyzer{}

		var createAnalyzer NewAnalyzer[ExapleModel] = func(html, u *string) (IAnalyzer[ExapleModel], error) {
			return analyzer, nil
		}

		err := errors.New("UNEXPECTED_ERROR")
		outBuf := &bytes.Buffer{}
		crawler_ := NewCrawler(
			&cache.MockCache{
				Get_: func(key string) (string, error) {
					return "", err
				},
			},
			createAnalyzer,
			&page_loader.MockPageLoader{
				LoadPage_: func(url string) (string, error) {
					return "", nil
				},
			},
			gotils.NewLogger(gotils.LogLevelInfo, outBuf, &bytes.Buffer{}),
			"http://demo.example",
			CrawlerConfig{},
		).(*crawler[ExapleModel])

		remainingUrlCh := make(chan string, 1)
		downloadedUrlCh := make(chan string, 1)
		downloadedUrlCh <- "asd"
		resultCh := make(chan *ExapleModel, 1)

		assert.True(t, crawler_.AnalyzePage(downloadedUrlCh, remainingUrlCh, resultCh, 1))
		assert.True(t, strings.Contains(outBuf.String(), err.Error()))
	})
}

func newMockCrawler() ICrawler[ExapleModel] {
	return &MockCrawler[ExapleModel]{}
}

func TestMockCrawler(t *testing.T) {
	t.Run("Test Crawl", func(t *testing.T) {
		crawler := newMockCrawler().(*MockCrawler[ExapleModel])

		crawler.Crawl_ = func(startingUrl string) <-chan *ExapleModel {
			ch := make(chan *ExapleModel, 1)

			go func() {
				ch <- &ExapleModel{Title: startingUrl, Content: "content"}
			}()

			return ch
		}

		startingUrl := "http://example.com"
		resultCh := crawler.Crawl(startingUrl)
		result := <-resultCh

		assert.Equal(t, startingUrl, result.Title)
	})

	t.Run("Test Stop", func(t *testing.T) {
		crawler := newMockCrawler().(*MockCrawler[ExapleModel])

		stopped := false
		crawler.Stop_ = func() {
			stopped = true
		}

		crawler.Stop()

		assert.True(t, stopped)
	})

	t.Run("Test WaitStopped", func(t *testing.T) {
		crawler := newMockCrawler().(*MockCrawler[ExapleModel])

		stopped := false
		crawler.WaitStopped_ = func() {
			stopped = true
		}

		crawler.WaitStopped()

		assert.True(t, stopped)
	})
}
