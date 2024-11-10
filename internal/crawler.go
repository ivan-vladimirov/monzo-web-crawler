package crawler

import (
	"log"
	"sync"
	"time"

	"github.com/ivan-vladimirov/monzo-web-crawler/internal/fetcher"
)

type UsedURL struct {
	URLs map[string]bool
	Mux  sync.RWMutex
}

func Crawl(url string, maxDepth int, delay time.Duration, used *UsedURL, wg *sync.WaitGroup) {
	defer wg.Done()

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
	links := fetcher.FetchLinks(url)
	log.Println("Crawled:", url)

	// Process each found link concurrently with a decreased depth
	wgLinks := &sync.WaitGroup{}
	for _, link := range links {
		wgLinks.Add(1)
		go func(link string) {
			defer wgLinks.Done()
			Crawl(link, maxDepth-1, delay, used, wgLinks)
			log.Println("Found↳", link)
		}(link)
	}

	wgLinks.Wait()
}