package crawler

import (
	"sync"
	"time"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/fetcher"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/parser"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/utils"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/shared"
)



type Crawler struct {
	fetcher     *fetcher.Fetcher
	parser      *parser.Parser
	logger      *utils.Logger
	rateLimiter *time.Ticker
	workerPool  chan struct{}

}

func NewCrawler(fetcher *fetcher.Fetcher, parser *parser.Parser, logger *utils.Logger, rateLimiter *time.Ticker, workerPoolSize int) *Crawler {
	return &Crawler{
		fetcher:     fetcher,
		parser:      parser,
		logger:      logger,
		rateLimiter: rateLimiter,
		workerPool:  make(chan struct{}, workerPoolSize), // Worker pool size
	}
}

// Crawl recursively visits a given URL and extracts internal links within the same domain and subdomain.
// It ensures depth constraints, avoids duplicate crawling using a mutex-protected map, and filters out
// unnecessary links such as those pointing to non-HTML files or fragments.
//
// Parameters:
// - url (string): The URL to be crawled.
// - maxDepth (int): The maximum depth of recursion allowed.
// - baseURL (string): The base URL of the domain to restrict crawling.
// - delay (time.Duration): The delay between requests to avoid overloading the server.
// - used (*shared.UsedURL): A shared structure for tracking crawled URLs and visited paths, ensuring thread safety.
// - wg (*sync.WaitGroup): A WaitGroup to synchronize goroutines and ensure all crawls complete before returning.
// - logger (*utils.Logger): A logger instance for structured and detailed logging.
//
// Behavior:
// - Normalizes the URL to maintain consistency and detect duplicates.
// - Skips URLs exceeding max depth, those with invalid formats, or non-HTML file extensions.
// - Fetches links from the URL using the fetcher package, then filters internal links via the parser package.
// - Uses concurrency with goroutines to crawl multiple links in parallel, while ensuring thread safety.
//
// Errors:
// - Logs and skips invalid URLs, fetch failures, or errors during normalization and parsing.
func (c *Crawler)  Crawl(url string, maxDepth int, baseURL string, delay time.Duration, used *shared.UsedURL, wg *sync.WaitGroup, logger *utils.Logger) {
	defer wg.Done() 

	depth, err := utils.CalculateDepthFromPath(url)
	logger.Info.Printf("Depth: %d, URL: %s\n", depth, url)
	if err != nil {
		logger.Error.Println("Error calculating depth for URL:", url, err)
		return
	}

	if depth > maxDepth {
		logger.Info.Printf("[MAX DEPTH REACHED] Depth: %d, URL: %s\n", depth, url)
		return
	}

	if utils.IsExcludedFileType(url) {
		logger.Info.Printf("[SKIPPED FILE] Filetype at Depth: %d, URL: %s\n", depth, url)
		return
	}


	canonicalURL, err := utils.NormalizeURL(url,baseURL)
	if err != nil {
		logger.Error.Printf("[MALFORMED] Skipping malformed URL: %s, Error: %v\n", url, err)
		return
	}
	
	if used.IsCrawledURL(canonicalURL) {
		logger.Info.Printf("[DUPLICATE] Depth: %d, URL: %s\n", depth, canonicalURL)
		return
	}

    c.workerPool <- struct{}{}
    defer func() { <-c.workerPool }()
    <-c.rateLimiter.C

	logger.Info.Printf("[CRAWLED] Depth: %d, URL: %s\n", depth, canonicalURL)

	links, err := c.fetcher.FetchLinks(canonicalURL, logger)
	if err != nil {
		logger.Info.Printf("[ERROR] Depth: %d, URL: %s, Error: %v\n", depth, canonicalURL, err)
		return
	}
	used.AddCrawledURL(canonicalURL)

	internalLinks := c.parser.CheckInternal(url, links, logger, canonicalURL, used)
	if len(internalLinks) == 0 {
		logger.Info.Printf("[SKIPPED] No valid internal links found for URL: %s\n", canonicalURL)
		return
	}

	for _, link := range internalLinks {
		normalizedLink, err := utils.NormalizeURL(link,baseURL)
		if err != nil {
			logger.Error.Printf("[MALFORMED] Skipping malformed URL: %s, Error: %v\n", url, err)
			return
		}

		if !used.IsCrawledURL(normalizedLink) {
			wg.Add(1)
			go c.Crawl(normalizedLink, maxDepth, baseURL, delay, used, wg, logger)
		}
	}
}
