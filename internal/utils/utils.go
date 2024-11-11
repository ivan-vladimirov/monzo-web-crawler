package utils

import (
	"net/url"
	"strings"
)

// NormalizeURL removes fragments and query parameters, enforces HTTPS, and removes trailing slashes for consistency.
func NormalizeURL(link string) string {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return link
	}

	parsedURL.Scheme = "https"

	parsedURL.Fragment = ""
	parsedURL.RawQuery = ""

	// Remove trailing slash
	parsedURL.Path = strings.TrimRight(parsedURL.Path, "/")

	return parsedURL.String()
}
// CalculateDepthFromPath determines the depth based on URL path segments relative to the base URL.
func CalculateDepthFromPath(baseURL, currentURL string) (int, error) {
	// Parse the base URL to extract the hostname and path
	base, err := url.Parse(baseURL)
	if err != nil {
		return 0, err
	}

	// Parse the current URL
	current, err := url.Parse(currentURL)
	if err != nil {
		return 0, err
	}

	// Ensure the current URL shares the same base hostname
	if base.Hostname() != current.Hostname() {
		return 0, nil // Treat it as a different domain or subdomain
	}

	// Count the path segments relative to the base path
	basePathSegments := strings.Split(strings.Trim(base.Path, "/"), "/")
	currentPathSegments := strings.Split(strings.Trim(current.Path, "/"), "/")

	// Calculate the relative depth by subtracting the base path depth
	relativeDepth := len(currentPathSegments) - len(basePathSegments)

	// Ensure depth is non-negative
	if relativeDepth < 0 {
		relativeDepth = 0
	}

	return relativeDepth, nil
}
// Function to convert user input to a valid URL
func Domain(ui string) (string, error) {
	if ui[len(ui)-1:] == "/" {
		ui = ui[:len(ui)-1]
	}

	parse, err := url.Parse(ui)
	if err != nil {
		return "", err
	}
	parse.Scheme = "http"
	return parse.String(), nil
}
