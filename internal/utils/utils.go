package utils

import (
	"net/url"
	"strings"
	"errors"
	"strconv"
	"regexp"
)

var validURLPattern = regexp.MustCompile(`^(https?|ftp)://([a-zA-Z0-9.-]+)(:[0-9]{1,5})?(/.*)?$`)


// NormalizeURL removes fragments and query parameters, enforces HTTPS, and removes trailing slashes for consistency.
func NormalizeURL(link string) (string, error) {
	if !validURLPattern.MatchString(link) {
		return "", errors.New("invalid URL format")
	}

	parsedURL, err := url.Parse(link)
	if err != nil {
		return link, err
	}

	if port := parsedURL.Port(); port != "" {
		if !isValidPort(port) {
			return "", errors.New("invalid port specified in URL")
		}
	}

	parsedURL.Scheme = "https"

	parsedURL.Fragment = ""
	parsedURL.RawQuery = ""

	parsedURL.Path = strings.TrimRight(parsedURL.Path, "/")

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