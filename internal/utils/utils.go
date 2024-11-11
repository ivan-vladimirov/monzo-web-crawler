package utils

import (
	"net/url"
	"strings"
)

// NormalizeURL removes fragments and query parameters, enforces HTTPS, and removes trailing slashes for consistency.
func NormalizeURL(link string) (string, error) {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return link, err
	}

	// Ensure the URL contains a valid port if specified
	if parsedURL.Port() != "" {
		_, err := url.ParseRequestURI(parsedURL.Scheme + "://" + parsedURL.Host)
		if err != nil {
			return "", err // Return an error if the port is invalid
		}
	}

	parsedURL.Scheme = "https"

	parsedURL.Fragment = ""
	parsedURL.RawQuery = ""

	// Remove trailing slash
	parsedURL.Path = strings.TrimRight(parsedURL.Path, "/")

	return parsedURL.String(), nil
}

// CalculateDepthFromPath determines the depth based on URL path segments relative to the base URL.
func CalculateDepthFromPath(currentURL string) (int, error) {
	current, err := url.Parse(currentURL)
	if err != nil {
		return 0, err
	}

	// Split the path and filter out empty segments
	pathSegments := strings.Split(strings.Trim(current.Path, "/"), "/")
	nonEmptySegments := 0
	for _, segment := range pathSegments {
		if segment != "" {
			nonEmptySegments++
		}
	}
	return nonEmptySegments, nil
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
