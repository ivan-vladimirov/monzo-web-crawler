package utils_test

import (
	"testing"

	"github.com/ivan-vladimirov/monzo-web-crawler/internal/utils"
)

func TestCalculateDepthFromPath(t *testing.T) {
	testCases := []struct {
		url      string
		expected int
	}{
		// Basic paths
		{"https://example.com", 0},
		{"https://example.com/", 0},
		{"https://example.com/level1", 1},
		{"https://example.com/level1/level2", 2},
		{"https://example.com/level1/level2/level3", 3},

		// Root path
		{"https://example.com/", 0},
		{"https://example.com/////", 0},

		// Empty and malformed URLs
		{"", 0},
		{"https://", 0},

		// Trailing slashes
		{"https://example.com/level1/", 1},
		{"https://example.com/level1/level2/", 2},

		// Multiple slashes in paths
		{"https://example.com//level1///level2", 2},
		{"https://example.com///level1/level2///level3", 3},

		// URLs with file extensions
		{"https://example.com/level1/file.pdf", 2},
		{"https://example.com/level1/level2/file.txt", 3},

		// URLs with query parameters and fragments
		{"https://example.com/level1?query=123", 1},
		{"https://example.com/level1#fragment", 1},
		{"https://example.com/level1/level2?query=123#fragment", 2},
	}

	for _, tc := range testCases {
		t.Run(tc.url, func(t *testing.T) {
			depth, err := utils.CalculateDepthFromPath(tc.url)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if depth != tc.expected {
				t.Errorf("For URL %s, expected depth %d, but got %d", tc.url, tc.expected, depth)
			}
		})
	}
}


func TestNormalizeURL(t *testing.T) {
	testCases := []struct {
		inputURL      string
		expectedURL   string
		expectError   bool
	}{
		// Basic URLs
		{"http://example.com", "https://example.com", false},
		{"https://example.com", "https://example.com", false},
		{"http://example.com/", "https://example.com", false},
		{"https://example.com/path", "https://example.com/path", false},
		{"http://example.com/path/", "https://example.com/path", false},

		// Malformed Ports
		{"https://example.com:port/path", "", true},
		{"http://example.com:1234path", "", true},
		{"https://example.com:8080", "https://example.com:8080", false}, // Valid port
		{"https://example.com:99999999", "", true}, // Invalid port


		// Different Schemes
		{"ftp://example.com/path", "https://example.com/path", false},
		{"example.com/path", "https://example.com/path", false},
		{"http://example.com/path/to/resource", "https://example.com/path/to/resource", false},

		// Trailing Slashes
		{"https://example.com/", "https://example.com", false},
		{"https://example.com/path/", "https://example.com/path", false},
		{"http://example.com/path/to/", "https://example.com/path/to", false},

		// Fragments and Query Parameters
		{"https://example.com/path#section", "https://example.com/path", false},
		{"http://example.com/path?query=value", "https://example.com/path", false},
		{"https://example.com/path?query=value#section", "https://example.com/path", false},

		// Subdomains and Paths
		{"http://sub.example.com/path", "https://sub.example.com/path", false},
		{"https://example.com/path/to/resource", "https://example.com/path/to/resource", false},

		// Special Characters in Path
		{"http://example.com/path/with%20spaces", "https://example.com/path/with%20spaces", false},
		{"https://example.com/path/with_unicode_ñ", "https://example.com/path/with_unicode_%C3%B1", false},
		{"https://example.com/!@#$%^&*()", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.inputURL, func(t *testing.T) {
			normalizedURL, err := utils.NormalizeURL(tc.inputURL, "http://example.com")

			// Check for expected error state
			if tc.expectError && err == nil {
				t.Errorf("Expected error for input URL %s, but got none", tc.inputURL)
			} else if !tc.expectError && err != nil {
				t.Errorf("Did not expect error for input URL %s, but got %v", tc.inputURL, err)
			}

			// Check for expected normalized URL
			if normalizedURL != tc.expectedURL {
				t.Errorf("For input URL %s, expected normalized URL %s, but got %s", tc.inputURL, tc.expectedURL, normalizedURL)
			}
		})
	}
}