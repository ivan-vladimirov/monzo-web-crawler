package parser

import (
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/utils"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/shared"
	"net/url"
	"strings"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}
// CheckInternal filters and extracts internal URLs from a given set of links.
// It determines whether a link belongs to the same domain as the base URL and avoids recursive paths.
//
// Parameters:
// - base (string): The base URL of the website for determining internal links.
// - links (map[string]bool): A map of candidate links to evaluate.
// - logger (*utils.Logger): Logger instance for structured logging of errors, warnings, and progress.
// - parentURL (string): The URL of the parent page to resolve relative links.
// - visitedPaths (*map[string]bool): A pointer to a map tracking visited paths to prevent recursion.
//
// Returns:
// - []string: A list of URLs determined to be internal and not recursively visited.
//
// Behavior:
// - Parses the base URL to extract its hostname for internal link comparison.
// - Iterates over each candidate link to:
//   1. Makes sure the base has a scheme.
//   2. Normalize the link using the parent URL.
//   3. Parse and validate the normalized link.
//   4. Check if the hostname matches the base URL (i.e., the link is internal).
//   5. Skip recursive paths that have already been visited.
//   6. Add valid internal links to the result list.
// - Logs ignored links (e.g., malformed URLs, external URLs, recursive paths).
func (p *Parser) CheckInternal(base string, links map[string]bool, logger *utils.Logger, parentURL string, used *shared.UsedURL) []string {
	var internalUrls []string

	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		base = "https://" + base
		logger.Info.Printf("Base URL missing scheme, added default scheme: %s\n", base)
	}

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

		if parsedLink.Hostname() != baseHostname {
			logger.Info.Println("Ignored external URL:", cleanedLink)
			continue
		}

		path := parsedLink.Path
		if used.IsVisitedPath(path) {
			logger.Info.Printf("Ignoring recursive path: %s\n", cleanedLink)
			continue
		}

		used.AddVisitedPath(path)

		internalUrls = append(internalUrls, cleanedLink)
		logger.Info.Println("Added internal URL:", cleanedLink)
	}

	return internalUrls
}
