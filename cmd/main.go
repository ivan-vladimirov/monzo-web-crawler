package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/crawler"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/fetcher"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/parser"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/shared"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/utils"
	"os"
	"sync"
	"time"
)

func main() {
	logger := utils.NewLogger()

	domain := flag.String("url", "", "Starting URL for the web crawler")
	maxDepth := flag.Int("max-depth", 3, "Maximum depth to crawl")
	outputFile := flag.String("output", "output.json", "File to save the JSON output")
	delay := flag.Duration("delay", 100*time.Millisecond, "Delay between requests (e.g., 100ms, 1s)")

	flag.Parse()

	if *domain == "" {
		logger.Error.Println("USAGE: ./monzo-web-crawler -url=http://monzo.com -max-depth=3 -delay=100ms -output=mozno.json")
		return
	}

	fetcher := fetcher.NewFetcher(10 * time.Second)
	parser := parser.NewParser()
	rateLimiter := time.NewTicker(100 * time.Millisecond)

	defer rateLimiter.Stop()

	cr := crawler.NewCrawler(fetcher, parser, logger, rateLimiter, 10)

	crawled := &shared.UsedURL{
		CrawledURLs:  make(map[string]bool),
		VisitedPaths: make(map[string]bool),
	}

	wg := &sync.WaitGroup{}

	wg.Add(1)
	cr.Crawl(*domain, *maxDepth, *domain, *delay, crawled, wg, logger)
	wg.Wait()

	crawledJSON, err := json.MarshalIndent(struct {
		URLs map[string]bool `json:"urls"`
	}{
		URLs: crawled.CrawledURLs,
	}, "", "  ")

	if err != nil {
		logger.Error.Println("Error while marshalling URLs:", err)
		return
	}

	err1 := utils.SaveJSONToFile(crawledJSON, *outputFile)
	if err1 != nil {
		logger.Error.Printf("Failed to save JSON to file: %v", err)
		os.Exit(1)
	}

	fmt.Println(string(crawledJSON))
}
