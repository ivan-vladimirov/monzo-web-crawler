package parser

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/ivan-vladimirov/monzo-web-crawler/internal/utils"
)

// Function to check if the link is valid and belongs to the same domain
func CheckInternal(base string, links map[string]bool, logger *utils.Logger) []string {
	var internalUrls []string
	for link := range links {
		if !isValidURL(link, logger) {
			continue
		}

		parsedURL, err := url.Parse(link)
		if err != nil {
			logger.Error.Println("Error parsing URL:", err)
			continue
		}

		// Remove the anchor fragment
		parsedURL.Fragment = ""
		cleanedLink := parsedURL.String()

		// Process relative and absolute URLs
		if strings.HasPrefix(cleanedLink, "/") {
			resolvedURL := fmt.Sprintf("%s%s", strings.TrimRight(base, "/"), cleanedLink)
			internalUrls = append(internalUrls, resolvedURL)
			logger.Info.Println("Resolved relative URL to:", resolvedURL)
		} else if strings.HasPrefix(cleanedLink, base) {
			internalUrls = append(internalUrls, cleanedLink)
			logger.Info.Println("Added internal URL:", cleanedLink)
		}
	}
	return internalUrls
}

// Function to validate the URL's scheme
func isValidURL(link string, logger *utils.Logger) bool {
	if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") && !strings.HasPrefix(link, "/") {
		logger.Info.Println("Ignoring non-http link:", link)
		return false
	}
	return true
}
