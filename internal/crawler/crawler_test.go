package crawler_test

import (
	"sync"
	"testing"
	"time"

	"github.com/ivan-vladimirov/monzo-web-crawler/internal/crawler"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/fetcher"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/parser"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/shared"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/utils"
)

var (
	crawlerInstance *crawler.Crawler
	logger          *utils.Logger
	setupOnce       sync.Once
)

func setup() {
	setupOnce.Do(func() {
		logger = utils.NewLogger()
		fetcherInstance := fetcher.NewFetcher(10 * time.Second)
		parserInstance := parser.NewParser()
		crawlerInstance = crawler.NewCrawler(fetcherInstance, parserInstance, logger, time.NewTicker(100*time.Millisecond), 10)
	})
}

func TestCrawl(t *testing.T) {
	setup()
	baseURL := "https://example.com"
	mockUsed := &shared.UsedURL{
		CrawledURLs:  make(map[string]bool),
		VisitedPaths: make(map[string]bool),
	}
	var wg sync.WaitGroup

	t.Run("Skip Max Depth", func(t *testing.T) {
		mockUsed.CrawledURLs = make(map[string]bool)
		wg.Add(1)
		crawlerInstance.Crawl("https://example.com/depth/4", 1, baseURL, 0, mockUsed, &wg, logger)
		wg.Wait()

		if len(mockUsed.CrawledURLs) > 0 {
			t.Errorf("Expected no CrawledURLs to be crawled, but got %d", len(mockUsed.CrawledURLs))
		}
	})

	t.Run("Skip File Types", func(t *testing.T) {
		mockUsed.CrawledURLs = make(map[string]bool)
		wg.Add(1)
		crawlerInstance.Crawl("https://example.com/file.pdf", 3, baseURL, 0, mockUsed, &wg, logger)
		wg.Wait()

		if len(mockUsed.CrawledURLs) > 0 {
			t.Errorf("Expected no CrawledURLs to be crawled, but got %d", len(mockUsed.CrawledURLs))
		}
	})

	t.Run("Avoid Duplicate CrawledURLs", func(t *testing.T) {
		mockUsed.CrawledURLs = make(map[string]bool)
		mockUsed.CrawledURLs["https://example.com/duplicate"] = true

		wg.Add(1)
		crawlerInstance.Crawl("https://example.com/duplicate", 3, baseURL, 0, mockUsed, &wg, logger)
		wg.Wait()

		if len(mockUsed.CrawledURLs) > 1 {
			t.Errorf("Expected 1 URL to remain, but got %d", len(mockUsed.CrawledURLs))
		}
	})

	t.Run("Handle FetchLinks Error", func(t *testing.T) {
		mockUsed.CrawledURLs = make(map[string]bool)
		wg.Add(1)
		crawlerInstance.Crawl("https://invalid-url", 3, baseURL, 0, mockUsed, &wg, logger)
		wg.Wait()

		if len(mockUsed.CrawledURLs) > 0 {
			t.Errorf("Expected no CrawledURLs to be crawled, but got %d", len(mockUsed.CrawledURLs))
		}
	})
}
