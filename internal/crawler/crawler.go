package crawler

import (
	"sync"
	"time"

	"github.com/ivan-vladimirov/monzo-web-crawler/internal/fetcher"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/utils"
)

type UsedURL struct {
	URLs map[string]bool
	Mux  sync.RWMutex
}



func Crawl(url string, maxDepth int, delay time.Duration, used *UsedURL, wg *sync.WaitGroup, logger *utils.Logger) {
	defer wg.Done() // Ensure Done() is called when this function returns

	// Limit crawling depth
	if maxDepth <= 0 {
		return
	}

	// Delay between requests
	time.Sleep(delay)

	// Lock for reading/writing to shared map
	used.Mux.Lock()
	if used.URLs[url] {
		used.Mux.Unlock()
		return
	}
	used.URLs[url] = true
	used.Mux.Unlock()

	// Fetch links from the page
	links := fetcher.FetchLinks(url, logger)
	logger.Info.Println("Crawled:", url)

	// Process each found link concurrently with a decreased depth
	for _, link := range links {
		wg.Add(1) 
		go func(link string) {
			logger.Info.Println("Found↳", link)
			Crawl(link, maxDepth-1, delay, used, wg, logger)
		}(link)
	}
}
