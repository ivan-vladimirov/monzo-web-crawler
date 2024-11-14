package utils

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"os"
)

var validAbsoluteURLPattern = regexp.MustCompile(`^((https?|ftp):\/\/)?(([a-zA-Z0-9.-]+\.[a-zA-Z]{2,})|(\d{1,3}(\.\d{1,3}){3})|(\[([a-fA-F0-9:]+)\]))(:[0-9]{1,5})?(\/.*)?$`)
var validRelativeURLPattern = regexp.MustCompile(`^/[^?#]*$`)
var excludedFileTypes = []string{".pdf", ".jpg", ".png", ".docx"}

// SaveJSONToFile writes a given data structure as a JSON file.
// Parameters:
// - data (interface{}): The data structure to be marshaled into JSON.
// - filename (string): The name of the file to save the JSON output.
//
// Returns:
// - error: Returns an error if marshaling or file operations fail.
func SaveJSONToFile(jsonData []byte, filename string) error {
	// Create or open the output file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the JSON data directly to the file
	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}
// RandFloat generates a random float64 value between 0.0 and 1.0.
//
// Returns:
// - float64: A random float64 value in the range [0.0, 1.0].
func RandFloat() float64 {
	return float64(time.Now().UnixNano()%1000) / 1000.0
}

// IsExcludedFileType checks if the provided URL ends with a file extension
// that is part of the excluded file types list.
// Parameters:
// - url (string): The URL to check for excluded file extensions.
//
// Returns:
// - bool: True if the URL's file extension matches any in the excluded file types list; false otherwise.
func IsExcludedFileType(url string) bool {
	ext := strings.ToLower(filepath.Ext(url))
	for _, excluded := range excludedFileTypes {
		if ext == excluded {
			return true
		}
	}
	return false
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
	if !validAbsoluteURLPattern.MatchString(link) && !validRelativeURLPattern.MatchString(link) {
		return "", errors.New("invalid URL format")
	}

	parsedURL, err := url.Parse(link)
	if err != nil {
		return "", fmt.Errorf("error parsing URL: %v", err)
	}

	if port := parsedURL.Port(); port != "" && !isValidPort(port) {
		return "", errors.New("invalid port specified in URL")
	}

	if parsedURL.Scheme == "" {
		if strings.Contains(link, ".") && !strings.HasPrefix(link, "/") {
			parsedURL, err = url.Parse("http://" + link)
			if err != nil {
				return "", fmt.Errorf("error parsing bare domain: %v", err)
			}
		} else {
			base, err := url.Parse(baseURL)
			if err != nil {
				return "", fmt.Errorf("error parsing base URL: %v", err)
			}
			parsedURL = base.ResolveReference(parsedURL)
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

	pathSegments := strings.Split(strings.Trim(current.Path, "/"), "/")
	nonEmptySegments := 0
	for _, segment := range pathSegments {
		if segment != "" {
			nonEmptySegments++
		}
	}
	return nonEmptySegments, nil
}
