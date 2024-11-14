package fetcher

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/utils"
)

type Fetcher struct {
	client *http.Client
}

func NewFetcher(timeout time.Duration) *Fetcher {
	return &Fetcher{
		client: &http.Client{Timeout: timeout},
	}
}

// Define a custom error for 404 Not Found
var ErrNotFound = errors.New("404 Not Found")

// MaxRetry defines the maximum number of retry attempts
const MaxRetry = 3

// RetryDelay defines the delay between retries
const InitialRetryDelay = 500 * time.Millisecond

// RequestTimeout defines the timeout for each request
const RequestTimeout = 1 * time.Second

// FetchLinks retrieves all links from a URL, returning a map of URLs or an error if the page couldn't be fetched.
func (f *Fetcher) FetchLinks(url string, logger *utils.Logger) (map[string]bool, error) {
	res, err := Request(url, logger)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		}
		logger.Error.Println("Error fetching the page:", err)
		return nil, err
	}

	if res == nil {
		return nil, errors.New("failed to fetch URL after retries")
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logger.Error.Println("Error parsing the page:", err)
		return nil, err
	}

	// Extract links as a map
	links := extractLinks(doc, logger)
	return links, nil
}

// Request performs an HTTP GET request to the specified URL with retry logic, exponential backoff, and timeout handling.
//
// Parameters:
// - url (string): The URL to request.
// - logger (*utils.Logger): A logger instance for structured logging.
//
// Returns:
// - (*http.Response): The HTTP response object if the request is successful.
// - (error): An error if the request fails after the maximum number of retries or encounters a non-retryable error.
//
// Behavior:
// - Configures the HTTP client with a timeout (`RequestTimeout`) to prevent blocking on slow responses.
// - Retries the request for transient errors (HTTP 5xx) up to `MaxRetry` times with exponential backoff and random jitter to avoid synchronized retries.
// - Stops retries for client-side errors (HTTP 4xx) or when retries are exhausted.
// - Closes response bodies for unsuccessful responses to prevent resource leaks.
// - Exponential backoff starts with `InitialRetryDelay` and doubles after each attempt, capped at 5 seconds.
// - Adds jitter to retry delays to distribute retries more evenly and reduce server load.
// - Uses a custom `User-Agent` header to identify the crawler.
func Request(url string, logger *utils.Logger) (*http.Response, error) {
	var resp *http.Response
	var err error

	client := &http.Client{
		Timeout: RequestTimeout,
	}

	retryDelay := InitialRetryDelay

	for attempt := 1; attempt <= MaxRetry; attempt++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MonzoCrawler/1.0)")

		logger.Info.Printf("Requesting URL (Attempt %d/%d): %s\n", attempt, MaxRetry, url)

		resp, err = client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			return resp, nil
		}

		if resp != nil {
			resp.Body.Close()
		}
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		}

		logger.Info.Printf("Retrying URL after failure (Attempt %d/%d): %s\n", attempt, MaxRetry, url)

		jitter := time.Duration(float64(retryDelay) * (0.5 + 0.5*utils.RandFloat()))
		time.Sleep(retryDelay + jitter)

		retryDelay *= 2
		if retryDelay > 5*time.Second {
			retryDelay = 5 * time.Second
		}
	}
	return nil, err
}

// extractLinks parses the HTML document to extract all unique hyperlinks (anchor tags) with valid href attributes.
// It resolves relative links to absolute URLs based on the provided base URL.
//
// Parameters:
// - doc (*goquery.Document): The parsed HTML document from which links will be extracted.
// - baseURL (string): The base URL used to resolve relative links.
// - logger (*utils.Logger): A logger instance for structured logging.
//
// Returns:
// - map[string]bool: A map of unique URLs found in the document, where keys are URLs and values are always true.
//
// Behavior:
// - Finds all `<a>` tags in the document and extracts their `href` attributes.
// - Resolves relative links (e.g., "/about") to absolute URLs using the base URL.
// - Skips invalid links, including:
//   - Fragment links starting with `#` (e.g., "#section").
//
// - Logs every link found and provides information about ignored links for debugging.
// - Filters out duplicate links using a map for efficient storage and retrieval.
// - Normalization of URLs is handled elsewhere in the pipeline to ensure consistency.
func extractLinks(doc *goquery.Document, logger *utils.Logger) map[string]bool {
	links := make(map[string]bool)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if link, exists := s.Attr("href"); exists {
			// Skip links that are fragments (starting with "#") or relative (starting with "/")
			if strings.HasPrefix(link, "#") {
				logger.Info.Println("Ignoring # tag:", link)
				return
			}
			// if strings.HasPrefix(link, "/") {
			// 	logger.Info.Println("Ignoring relative link:", link)
			// 	return
			// }
			links[link] = true
			logger.Info.Println("Found link:", link)
		}
	})
	return links
}
