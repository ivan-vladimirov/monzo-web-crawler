package main

import (
	"flag"
	"sync"
	"time"
	"fmt"
	"encoding/json"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/fetcher"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/parser"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/crawler"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/utils"
)

func main() {
	logger := utils.NewLogger()

	domain := flag.String("url", "", "Starting URL for the web crawler")
	maxDepth := flag.Int("max-depth", 3, "Maximum depth to crawl")
	delay := flag.Duration("delay", 100*time.Millisecond, "Delay between requests (e.g., 100ms, 1s)")
	flag.Parse()

	if *domain == "" {
		logger.Error.Println("USAGE: ./monzo-web-crawler -url=http://monzo.com -max-depth=3 -delay=100ms")
		return
	}
	fetcher := fetcher.NewFetcher(10 * time.Second)
	parser := parser.NewParser()
	rateLimiter := time.NewTicker(100 * time.Millisecond)
	defer rateLimiter.Stop()

	cr := crawler.NewCrawler(fetcher, parser, logger, rateLimiter, 10)

	crawled := &crawler.UsedURL{
		CrawledURLs:         make(map[string]bool),
		VisitedPaths: make(map[string]bool),
	}
	wg := &sync.WaitGroup{}

	wg.Add(1)
	cr.Crawl(*domain, *maxDepth, *domain, *delay, crawled, wg, logger)
	wg.Wait()

	prettyPrinted, err := json.MarshalIndent(struct {
		URLs map[string]bool `json:"urls"`
	}{
		URLs: crawled.CrawledURLs,
	}, "", "  ")
	if err != nil {
		logger.Error.Println("Error while marshalling URLs:", err)
		return
	}

	fmt.Println(string(prettyPrinted))
}