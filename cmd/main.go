package main

import (
	"flag"
	"sync"
	"time"

	"github.com/ivan-vladimirov/monzo-web-crawler/internal/crawler"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/utils"
)

func main() {
	// Initialize the logger
	logger := utils.NewLogger()

	// Define and parse command-line flags
	domain := flag.String("url", "", "Starting URL for the web crawler")
	maxDepth := flag.Int("max-depth", 3, "Maximum depth to crawl")
	delay := flag.Duration("delay", 100*time.Millisecond, "Delay between requests (e.g., 100ms, 1s)")
	flag.Parse()

	// Validate URL input
	if *domain == "" {
		logger.Error.Println("USAGE: ./monzo-web-crawler -url=http://monzo.com -max-depth=3 -delay=100ms")
		return
	}

	// Initialize the crawled URLs tracker
	crawled := &crawler.UsedURL{URLs: make(map[string]bool)}
	wg := &sync.WaitGroup{}

	// Start crawling
	crawler.Crawl(*domain, *maxDepth, *delay, crawled, wg, logger)
}