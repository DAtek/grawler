package grawler

import (
	"sync"
	"time"

	"github.com/DAtek/grawler/cache"
	"github.com/DAtek/grawler/page_loader"

	"github.com/DAtek/gotils"
)

type ICrawler[T any] interface {
	Crawl(startingUrl string) chan *T
	Stop()
	WaitStopped()
}

type MockCrawler[T any] struct {
	Crawl_       func(startingUrl string) chan *T
	Stop_        func()
	WaitStopped_ func()
}

func (c MockCrawler[T]) Crawl(startingUrl string) chan *T {
	return c.Crawl_(startingUrl)
}

func (c MockCrawler[T]) Stop() {
	c.Stop_()
}

func (c MockCrawler[T]) WaitStopped() {
	c.WaitStopped_()
}

type CrawlerConfig struct {
	PageLoaders         int
	PageAnalyzers       int
	RemainingUrlChSize  int
	DownloadedUrlChSize int
	ResultChSize        int
}

func (c *CrawlerConfig) validate() {
	c.PageLoaders = maxInt(c.PageLoaders, 1)
	c.PageAnalyzers = maxInt(c.PageAnalyzers, 1)
	c.RemainingUrlChSize = maxInt(c.RemainingUrlChSize, 10)
	c.DownloadedUrlChSize = maxInt(c.DownloadedUrlChSize, 10)
	c.ResultChSize = maxInt(c.ResultChSize, 10)
}

func NewCrawler[T any](
	cache cache.ICache,
	createAnalyzer NewAnalyzer[T],
	pageLoader page_loader.IPageLoader,
	logger *gotils.Logger,
	baseUrl string,
	config CrawlerConfig,
) ICrawler[T] {
	config.validate()

	return &crawler[T]{
		cache:          cache,
		createAnalyzer: createAnalyzer,
		pageLoader:     pageLoader,
		logger:         logger,
		urlRegistry:    newStringRegistry(),
		baseUrl:        baseUrl,
		stopCh:         make(chan struct{}, 1),
		wg:             &sync.WaitGroup{},
		config:         &config,
	}
}

type crawler[T any] struct {
	cache          cache.ICache
	createAnalyzer NewAnalyzer[T]
	pageLoader     page_loader.IPageLoader
	urlRegistry    *stringRegistry
	baseUrl        string
	logger         *gotils.Logger
	wg             *sync.WaitGroup
	config         *CrawlerConfig
	stopCh         chan struct{}
}

func (c crawler[T]) Crawl(startingUrl string) chan *T {
	resultCh := make(chan *T, c.config.ResultChSize)
	remainingUrlCh := make(chan string, c.config.RemainingUrlChSize)
	downloadedUrlCh := make(chan string, c.config.DownloadedUrlChSize)
	c.wg.Add(c.totalWorkers())
	c.urlRegistry.add(startingUrl)
	remainingUrlCh <- startingUrl

	pageLoader := func(i int) {
		defer func() {
			c.logger.Debug("Stopping page loader %d", i)
			c.wg.Done()
		}()

		for c.LoadPage(remainingUrlCh, downloadedUrlCh, i) {
		}
	}

	for i := 0; i < c.config.PageLoaders; i++ {
		go pageLoader(i)
	}

	pageAnalyzer := func(i int) {
		defer func() {
			c.logger.Debug("Stopping page analyzer %d", i)
			c.wg.Done()
		}()

		for c.AnalyzePage(downloadedUrlCh, remainingUrlCh, resultCh, i) {
		}
	}

	for i := 0; i < c.config.PageAnalyzers; i++ {
		go pageAnalyzer(i)
	}

	go func() {
		defer func() {
			c.logger.Debug("Stopping status logger")
			c.wg.Done()
		}()

		for {
			select {
			case <-c.stopCh:
				return
			default:
				c.logger.Debug(
					"Remaining URL ch: %d | Downloaded URL ch: %d | Model ch: %d",
					len(remainingUrlCh),
					len(downloadedUrlCh),
					len(resultCh),
				)
				time.Sleep(1 * time.Second)
			}
		}
	}()

	return resultCh
}

func (c crawler[T]) Stop() {
	for i := 0; i < c.totalWorkers(); i++ {
		c.stopCh <- struct{}{}
	}
}

func (c crawler[T]) WaitStopped() {
	c.wg.Wait()
}

func (c crawler[T]) totalWorkers() int {
	return c.config.PageLoaders + c.config.PageAnalyzers + 1
}

func (c crawler[T]) LoadPage(remainingUrlCh chan string, downloadedUrlChan chan string, i int) bool {
	select {
	case <-c.stopCh:
		return false
	case newUrl := <-remainingUrlCh:
		if c.cache.Has(newUrl) {
			c.logger.Debug("loadPage(%d) | Found in cache %s", i, newUrl)
			downloadedUrlChan <- newUrl
			return true
		}

		c.logger.Info("loadPage(%d) | Downloading from %s", i, newUrl)
		page, err := c.pageLoader.LoadPage(newUrl)
		if err != nil {
			c.logger.Error("loadPage(%d) | Error loading from '%s' Error: %s", i, newUrl, err)
			return true
		}

		if err := c.cache.Set(newUrl, page); err != nil {
			c.logger.Error("loadPage(%d) | Error saving to cache. '%s' Error: %s", i, newUrl, err)
			return true
		}
		downloadedUrlChan <- newUrl
		return true
	}

}

func (c crawler[T]) AnalyzePage(downloadedUrlCh, remainingUrlCh chan string, resultCh chan *T, i int) bool {
	select {
	case <-c.stopCh:
		return false
	case newUrl := <-downloadedUrlCh:
		page, err := c.cache.Get(newUrl)

		if err != nil {
			c.logger.Error("analyzePage(%d) | Error loading from '%s' Error: %s", i, newUrl, err)
			return true
		}

		c.logger.Debug("analyzePage(%d) | Analyzing page %s", i, newUrl)
		analyzer, err := c.createAnalyzer(&page, &newUrl)
		if err != nil {
			c.logger.Error("analyzePage(%d) | Failed to create the analyzer. URL: %s Error: %s", i, newUrl, err)
			return true
		}

		if model := analyzer.GetModel(); model != nil {
			c.logger.Info("analyzePage(%d) | Collected model for %s", i, newUrl)
			resultCh <- model
		}

		for _, newUrl := range c.urlRegistry.getNew(analyzer.GetUrls()) {
			c.urlRegistry.add(newUrl)
			if string(newUrl[0]) == "/" {
				newUrl = joinPath(c.baseUrl, newUrl)
			}
			c.logger.Debug("Adding URL: %s", newUrl)
			remainingUrlCh <- newUrl
		}
		return true
	}
}

func maxInt(x, y int) int {
	if x >= y {
		return x
	}

	return y
}
