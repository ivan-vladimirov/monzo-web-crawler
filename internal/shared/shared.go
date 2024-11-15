package shared

import "sync"

// UsedURL is a thread-safe structure for tracking visited CrawledURLs.
type UsedURL struct {
	CrawledURLs  map[string]bool
	VisitedPaths map[string]bool
	Mux          sync.RWMutex
}

// Check if a URL is already visited (thread-safe read).
func (u *UsedURL) IsCrawledURL(url string) bool {
	u.Mux.RLock()
	defer u.Mux.RUnlock()
	return u.CrawledURLs[url]
}

// Mark a URL as visited (thread-safe write).
func (u *UsedURL) AddCrawledURL(url string) {
	u.Mux.Lock()
	defer u.Mux.Unlock()
	u.CrawledURLs[url] = true
}

// Add a visited path.
func (u *UsedURL) AddVisitedPath(path string) {
	u.Mux.Lock()
	defer u.Mux.Unlock()
	u.VisitedPaths[path] = true
}

// Check if a URL is already visited (thread-safe read).
func (u *UsedURL) IsVisitedPath(path string) bool {
	u.Mux.RLock()
	defer u.Mux.RUnlock()
	return u.VisitedPaths[path]
}
