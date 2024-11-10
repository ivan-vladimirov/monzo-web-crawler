package crawler

import (
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



func Crawl(url string, maxDepth int, delay time.Duration, used *UsedURL, wg *sync.WaitGroup, logger *utils.Logger) {
	wg.Add(1)
	defer wg.Done() // Ensure Done() is called when this function returns

	// Limit crawling depth
	if maxDepth <= 0 {
		return
	}
	// Normalize the URL to a canonical form to avoid duplicates
	canonicalURL := parser.NormalizeURL(url)

	// Lock for reading/writing to shared map
	used.Mux.Lock()
	if used.URLs[canonicalURL] {
		used.Mux.Unlock()
		
		logger.Info.Println("Already crawled:", canonicalURL)
		return
	}

	used.URLs[canonicalURL] = true
	used.Mux.Unlock()

	time.Sleep(delay)

	links, err := fetcher.FetchLinks(canonicalURL, logger)
	if err != nil {
		logger.Info.Println("Error fetching ", err)
	}	

	logger.Info.Println("Crawled:", canonicalURL)

	// Process each found link concurrently with a decreased depth
	for _, link := range links {
		go func(link string) {
			logger.Info.Println("Found↳", link)
			Crawl(link, maxDepth-1, delay, used, wg, logger)
		}(link)
	}
}
