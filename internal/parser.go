package parser

import (
	"fmt"
	"strings"
)

// Function to check if the link is valid and belongs to the same domain
func CheckInternal(base string, links map[string]bool) []string {
	var internalUrls []string
	for link := range links {
		if !isValidURL(link) {
			continue
		}
		if strings.HasPrefix(link, base) {
			internalUrls = append(internalUrls, link)
		} else if strings.HasPrefix(link, "/") {
			resolvedURL := fmt.Sprintf("%s%s", base, link)
			internalUrls = append(internalUrls, resolvedURL)
		}
	}
	return internalUrls
}

// Function to validate the URL's scheme
func isValidURL(link string) bool {
	if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") && !strings.HasPrefix(link, "/") {
		InfoLogger.Println("Ignoring non-http link:", link)
		return false
	}
	return true
}
