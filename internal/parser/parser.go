package parser

import (
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/utils"
	"net/url"
	"strings"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) CheckInternal(base string, links map[string]bool, logger *utils.Logger, parentURL string, visitedPaths *map[string]bool) []string {
	var internalUrls []string

	baseURL, err := url.Parse(base)
	if err != nil {
		logger.Error.Println("Error parsing base URL:", err)
		return internalUrls
	}
	baseHostname := baseURL.Hostname()

	for link := range links {
		cleanedLink, err := utils.NormalizeURL(strings.TrimSpace(link), parentURL)
		if err != nil {
			logger.Error.Printf("Skipping malformed URL: %s, Error: %v\n", link, err)
			continue
		}

		parsedLink, err := url.Parse(cleanedLink)

		if err != nil {
			logger.Error.Printf("Error parsing URL: %s, Error: %v\n", cleanedLink, err)
			continue
		}

		// Check if the link is internal
		if parsedLink.Hostname() != baseHostname {
			logger.Info.Println("Ignored external URL:", cleanedLink)
			continue
		}

		// Detect recursive paths
		path := parsedLink.Path
		if (*visitedPaths)[path] {
			logger.Info.Printf("Ignoring recursive path: %s\n", cleanedLink)
			continue
		}

		// Mark the path as visited
		(*visitedPaths)[path] = true

		// Add to internal URLs
		internalUrls = append(internalUrls, cleanedLink)
		logger.Info.Println("Added internal URL:", cleanedLink)
	}

	return internalUrls
}
