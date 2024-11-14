package utils

import (
	"net/url"
	"strings"
	"errors"
	"strconv"
	"regexp"
	"time"
	"fmt"
)

var validAbsoluteURLPattern = regexp.MustCompile(`^(https?|http|ftp)://([a-zA-Z0-9.-]+)(:[0-9]{1,5})?(/.*)?$`)
var validRelativeURLPattern = regexp.MustCompile(`^/[^?#]*$`)

// Helper function to generate random jitter
func RandFloat() float64 {
	return float64(time.Now().UnixNano()%1000) / 1000.0
}

// NormalizeURL processes a URL to ensure consistency by removing fragments, query parameters, and trailing slashes.
// For absolute URLs, it enforces HTTPS. For relative URLs, it simply normalizes the path.
//
// Parameters:
// - link (string): The URL or path to be normalized.
//
// Returns:
// - (string, error): A normalized URL if valid, or an error if the format is invalid.
func NormalizeURL(link string, baseURL string) (string, error) {
    // Check if the link matches either the absolute or relative URL pattern
    if !validAbsoluteURLPattern.MatchString(link) && !validRelativeURLPattern.MatchString(link) {
        return "", errors.New("invalid URL format")
    }

    // Parse the URL
    parsedURL, err := url.Parse(link)
    if err != nil {
        return "", fmt.Errorf("error parsing URL: %v", err)
    }

    // Validate port if present
    if port := parsedURL.Port(); port != "" && !isValidPort(port) {
        return "", errors.New("invalid port specified in URL")
    }

    // If the URL is relative, resolve it against the base URL
    if parsedURL.Scheme == "" {
        base, err := url.Parse(baseURL)
        if err != nil {
            return "", fmt.Errorf("error parsing base URL: %v", err)
        }
        parsedURL = base.ResolveReference(parsedURL)
    }

    // Normalize the URL
    parsedURL.Scheme = "https"     // Enforce HTTPS for absolute URLs
    parsedURL.Fragment = ""        // Remove fragments
    parsedURL.RawQuery = ""        // Remove query parameters
    parsedURL.Path = strings.TrimRight(parsedURL.Path, "/") // Remove trailing slashes

    return parsedURL.String(), nil
}

func isValidPort(port string) bool {
	portNum, err := strconv.Atoi(port)
	if err != nil || portNum < 1 || portNum > 65535 {
		return false
	}
	return true
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
