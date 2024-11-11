package crawler

import (
	"path/filepath"
	"strings"
	"sync"
	"time"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/fetcher"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/parser"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/utils"
)

type UsedURL struct {
	URLs map[string]bool
	Mux  sync.RWMutex
}


func Crawl(url string, maxDepth int, baseURL string, delay time.Duration, used *UsedURL, wg *sync.WaitGroup, logger *utils.Logger) {
	defer wg.Done() 

	// Calculate the depth based on the URL path relative to the base
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
	// Skip URLs with file extensions (like .pdf)
	ext := strings.ToLower(filepath.Ext(url))
	if ext == ".pdf" || ext == ".jpg" || ext == ".png" || ext == ".docx" {
		logger.Info.Printf("[SKIPPED FILE] Filetype %s at Depth: %d, URL: %s\n", ext, depth, url)
		return
	}


	canonicalURL := utils.NormalizeURL(url)

	used.Mux.Lock()
	if used.URLs[canonicalURL] {
		used.Mux.Unlock()
		logger.Info.Printf("[DUPLICATE] Depth: %d, URL: %s\n", depth, canonicalURL)
		return
	}

	used.URLs[canonicalURL] = true
	used.Mux.Unlock()

	logger.Info.Printf("[CRAWLED] Depth: %d, URL: %s\n", depth, canonicalURL)

	time.Sleep(delay)

	links, err := fetcher.FetchLinks(canonicalURL, logger)
	if err != nil {
		logger.Info.Printf("[ERROR] Depth: %d, URL: %s, Error: %v\n", depth, canonicalURL, err)
		return
	}

	// Filter internal links from the fetched links
	internalLinks := parser.CheckInternal(url, links, logger, canonicalURL)

	for _, link := range internalLinks {
		normalizedLink := utils.NormalizeURL(link)

		// Lock to check if the link has already been crawled
		used.Mux.Lock()
		_, alreadyCrawled := used.URLs[normalizedLink]
		used.Mux.Unlock()

		// Only proceed if the link hasn't been crawled
		if !alreadyCrawled {
			wg.Add(1)
			go Crawl(normalizedLink, maxDepth, baseURL, delay, used, wg, logger)
		}
	}
}
