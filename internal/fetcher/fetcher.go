package fetcher

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/parser"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/utils"
)

func FetchLinks(url string, logger *utils.Logger) []string {
	res, err := Request(url, logger)
	if err != nil {
		logger.Error.Println("Error fetching the page:", err)
		return nil
	}
	doc, _ := goquery.NewDocumentFromResponse(res)
	links := extractLinks(doc, logger)
	return parser.CheckInternal(url, links, logger)
}

func Request(url string, logger *utils.Logger) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	logger.Info.Println("Requesting URL:", url)
	return client.Do(req)
}

func extractLinks(doc *goquery.Document, logger *utils.Logger) map[string]bool {
	links := make(map[string]bool)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if link, exists := s.Attr("href"); exists {
			link = strings.TrimSpace(link)

			// if strings.Contains(link, "#") {
			// 	logger.Info.Println("Anchor links not supported:", link)
			// 	return
			// }

			// Add only if the link is non-empty after trimming
			if link != "" {
				links[link] = true
			}
		}
	})
	return links
}
