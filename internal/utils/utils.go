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
func CalculateDepthFromPath(currentURL string) (int, error) {
	current, err := url.Parse(currentURL)
	if err != nil {
		return 0, err
	}

	currentPathSegments := len(strings.Split(strings.Trim(current.Path, "/"), "/"))

	return currentPathSegments, nil
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
