package fetcher

import (
	"errors"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/parser"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/utils"
)

// Define a custom error for 404 Not Found
var ErrNotFound = errors.New("404 Not Found")

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
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	logger.Info.Println("Requesting URL:", url)
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Check if the response is 404 Not Found
	if resp.StatusCode == http.StatusNotFound {
		resp.Body.Close()
		return nil, ErrNotFound // Return the custom 404 error
	}

	return resp, nil
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
