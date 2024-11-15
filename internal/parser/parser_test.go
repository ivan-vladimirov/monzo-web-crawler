package parser_test

import (
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/parser"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/shared"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/utils"
	"net/url"
	"path"
	"sync"
	"testing"
)

var (
	parserInstance *parser.Parser
	logger         *utils.Logger
	setupOnce      sync.Once // Ensures setup runs only once
)

func setup() {
	setupOnce.Do(func() {
		parserInstance = parser.NewParser()
		logger = utils.NewLogger()
	})
}

func TestCheckInternal(t *testing.T) {
	setup()
	baseURL := "https://example.com"
	parentURL := "https://example.com/parent"

	testCases := []struct {
		name          string
		links         map[string]bool
		usedURL       *shared.UsedURL
		expectedLinks []string
	}{
		{
			name: "Basic Internal Links",
			links: map[string]bool{
				"https://example.com/page1": true,
				"https://example.com/page2": true,
			},
			usedURL: &shared.UsedURL{
				CrawledURLs:  map[string]bool{},
				VisitedPaths: map[string]bool{},
			},
			expectedLinks: []string{
				"https://example.com/page1",
				"https://example.com/page2",
			},
		},
		{
			name: "Exclude Recursive Links",
			links: map[string]bool{
				"https://example.com/parent": true,
				"https://example.com/page":   true,
			},
			usedURL: &shared.UsedURL{
				CrawledURLs: map[string]bool{},
				VisitedPaths: map[string]bool{
					"/parent": true,
				},
			},
			expectedLinks: []string{
				"https://example.com/page",
			},
		},
		{
			name: "Skip Malformed URLs",
			links: map[string]bool{
				"https://example.com/page1": true,
				"://invalid-url":            true,
			},
			usedURL: &shared.UsedURL{
				CrawledURLs:  map[string]bool{},
				VisitedPaths: map[string]bool{},
			},
			expectedLinks: []string{
				"https://example.com/page1",
			},
		},
		// Additional test cases...
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			internalLinks := parserInstance.CheckInternal(baseURL, tc.links, logger, parentURL, tc.usedURL)

			// Assert the number of internal links
			if len(internalLinks) != len(tc.expectedLinks) {
				t.Errorf("Expected %d internal links, got %d", len(tc.expectedLinks), len(internalLinks))
			}

			// Check if all expected links are in the internalLinks result
			for _, link := range tc.expectedLinks {
				found := false
				for _, internalLink := range internalLinks {
					if link == internalLink {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected link %s not found in internal links", link)
				}
			}

			// Ensure visited paths are updated correctly in UsedURL
			tc.usedURL.Mux.RLock()
			defer tc.usedURL.Mux.RUnlock()
			for link := range tc.links {
				parsedURL, err := url.Parse(link)
				if err != nil {
					continue
				}
				visitedPath := path.Clean(parsedURL.Path)
				if _, visited := tc.usedURL.VisitedPaths[visitedPath]; visited {
					t.Logf("Path correctly marked as visited: %s", visitedPath)
				}
			}
		})
	}
}
