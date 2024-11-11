package fetcher

import (
	"errors"
	"net/http"
	"strings"
	"time"
	"github.com/PuerkitoBio/goquery"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/parser"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/utils"
)

// Define a custom error for 404 Not Found
var ErrNotFound = errors.New("404 Not Found")

// MaxRetry defines the maximum number of retry attempts
const MaxRetry = 3

// RetryDelay defines the delay between retries
const RetryDelay = 500 * time.Millisecond

func FetchLinks(url string, logger *utils.Logger) ([]string, error) {
	res, err := Request(url, logger)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			logger.Info.Println("404 Not Found:", url)
			return nil, ErrNotFound
		}
		logger.Error.Println("Error fetching the page:", err)
		return nil, err
	}

	doc, _ := goquery.NewDocumentFromResponse(res)
	links := extractLinks(doc, logger)
	return parser.CheckInternal(url, links, logger), nil
}

func Request(url string, logger *utils.Logger) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 1; attempt <= MaxRetry; attempt++ {
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
		logger.Info.Printf("Requesting URL (Attempt %d/%d): %s\n", attempt, MaxRetry, url)

		resp, err = client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			return resp, nil // Return immediately on success
		}

		// Close response body if not successful to avoid resource leaks
		if resp != nil {
			resp.Body.Close()
		}

		// Check for 404 specifically and exit retry loop if found
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		}

		// Log the retry and delay before the next attempt
		logger.Info.Printf("Retrying URL after failure (Attempt %d/%d): %s\n", attempt, MaxRetry, url)
		time.Sleep(RetryDelay)
	}

	// Return the last error after max retries have been exhausted
	return nil, err
}

func extractLinks(doc *goquery.Document, logger *utils.Logger) map[string]bool {
	links := make(map[string]bool)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if link, exists := s.Attr("href"); exists {
			if strings.Contains(link, "#") {
				logger.Info.Println("Ignoring # tag:", link)
				return
			}
			links[link] = true
		}
	})
	return links
}
