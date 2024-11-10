package main

import (
	"flag"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ivan-vladimirov/monzo-web-crawler/internal/crawler"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	// Initialize loggers
	InfoLogger = log.New(os.Stdout, "INFO: ", log.LstdFlags)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	// Define and parse command-line flags
	domain := flag.String("url", "", "Starting URL for the web crawler")
	maxDepth := flag.Int("max-depth", 3, "Maximum depth to crawl")
	delay := flag.Duration("delay", 100*time.Millisecond, "Delay between requests (e.g., 100ms, 1s)")

	flag.Parse()

	// Validate URL input
	if *domain == "" {
		ErrorLogger.Println("USAGE: ./monzo-web-crawler -url=http://monzo.com -max-depth=3 -delay=100ms")
		return
	}

	// Initialize the crawled URLs tracker
	crawled := &crawler.UsedURL{URLs: make(map[string]bool)}
	// Create a WaitGroup for concurrency
	wg := &sync.WaitGroup{}

	// Start crawling with configured depth and delay
	wg.Add(1)
	go crawler.Crawl(*domain, *maxDepth, *delay, crawled, wg)

	wg.Wait()
}
