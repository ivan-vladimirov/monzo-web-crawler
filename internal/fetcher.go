package fetcher

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/parser"
)

func FetchLinks(path string) []string {
	res, err := Request(path)
	if err != nil {
		ErrorLogger.Println("Error fetching the page:", err)
		return nil
	}
	doc, _ := goquery.NewDocumentFromResponse(res)
	links := extractLinks(doc)
	return parser.CheckInternal(path, links)
}

func Request(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	return client.Do(req)
}

func extractLinks(doc *goquery.Document) map[string]bool {
	links := make(map[string]bool)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if link, exists := s.Attr("href"); exists {
			links[link] = true
		}
	})
	return links
}
x