package main

import (
	"flag"
	"sync"
	"time"

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

	crawled := &crawler.UsedURL{URLs: make(map[string]bool)}
	wg := &sync.WaitGroup{}

	crawler.Crawl(*domain, *maxDepth, *delay, crawled, wg, logger)
}