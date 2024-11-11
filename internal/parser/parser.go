package parser

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/ivan-vladimirov/monzo-web-crawler/internal/utils"
)



// GetLastPathSegment returns the last segment of a given URL path
func GetLastPathSegment(link string) string {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return ""
	}
	pathSegments := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	return pathSegments[len(pathSegments)-1]
}

// CheckInternal filters links, keeping only internal links within the base domain, excluding subdomains and ignoring fragment links.
func CheckInternal(base string, links map[string]bool, logger *utils.Logger, parentURL string) []string {
	var internalUrls []string

	// Parse the base URL to extract its hostname
	baseURL, err := url.Parse(base)
	if err != nil {
		logger.Error.Println("Error parsing base URL:", err)
		return internalUrls // Return empty list if base URL is invalid
	}
	baseHostname := baseURL.Hostname()

	// Determine the last segment of the parent URL
	parentLastSegment := GetLastPathSegment(parentURL)

	for link := range links {
		// Skip links that are just fragments (e.g., "#section")
		if strings.HasPrefix(link, "#") {
			logger.Info.Println("Ignoring # tag:", link)
			continue
		}

		// Normalize the link to remove fragments and queries
		cleanedLink := utils.NormalizeURL(strings.TrimSpace(link))

		// Parse the cleaned link to get details about its structure
		parsedLink, err := url.Parse(cleanedLink)
		if err != nil {
			logger.Error.Println("Error parsing link:", cleanedLink, err)
			continue
		}

		// Check if the link is absolute and within the base domain
		if parsedLink.IsAbs() {
			linkHostname := parsedLink.Hostname()
			if linkHostname == baseHostname {
				// Detect recursive paths by comparing the last segment to the parent’s last segment
				if GetLastPathSegment(parsedLink.Path) == parentLastSegment {
					logger.Info.Println("Ignoring recursive link:", cleanedLink)
					continue
				}
				internalUrls = append(internalUrls, cleanedLink)
				logger.Info.Println("Added internal URL:", cleanedLink)
			} else {
				logger.Info.Println("Ignored external or subdomain URL:", cleanedLink)
			}
		} else {
			// For relative URLs, resolve them against the base URL
			resolvedURL := fmt.Sprintf("%s%s", strings.TrimRight(base, "/"), cleanedLink)
			
			// Detect recursive paths in relative URLs
			if GetLastPathSegment(resolvedURL) == parentLastSegment {
				logger.Info.Println("Ignoring recursive relative URL:", resolvedURL)
				continue
			}

			internalUrls = append(internalUrls, resolvedURL)
			logger.Info.Println("Resolved relative URL to:", resolvedURL)
		}
	}

	return internalUrls
}
