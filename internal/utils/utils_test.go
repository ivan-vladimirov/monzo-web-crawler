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