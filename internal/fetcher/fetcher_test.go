package fetcher_test

import (
	"errors"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/fetcher"
	"github.com/ivan-vladimirov/monzo-web-crawler/internal/utils"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

var (
	fetcherInstance *fetcher.Fetcher
	logger          *utils.Logger
	setupOnce       sync.Once //Only initialise once
)

func setup() {
	setupOnce.Do(func() {
		fetcherInstance = fetcher.NewFetcher(10 * time.Second)
		logger = utils.NewLogger()
	})
}

// TestFetchLinks_ValidURL tests fetching links from a valid URL
func TestFetchLinks_ValidURL(t *testing.T) {
	setup()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<a href="https://example.com/page1">Page1</a>
						<a href="https://example.com/page2">Page2</a>`))
	}))
	defer ts.Close()

	links, err := fetcherInstance.FetchLinks(ts.URL, logger)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedLinks := []string{
		"https://example.com/page1",
		"https://example.com/page2",
	}

	for _, link := range expectedLinks {
		if _, exists := links[link]; !exists {
			t.Errorf("Expected link %s in fetched links, but it was not found", link)
		}
	}
}

// TestFetchLinks_404Error tests handling of a 404 status code
func TestFetchLinks_404Error(t *testing.T) {
	setup()

	ts := httptest.NewServer(http.NotFoundHandler())
	defer ts.Close()

	_, err := fetcherInstance.FetchLinks(ts.URL, logger)

	if !errors.Is(err, fetcher.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

// TestFetchLinks_InvalidURL tests handling of an invalid URL
func TestFetchLinks_InvalidURL(t *testing.T) {
	setup()

	_, err := fetcherInstance.FetchLinks("://invalid-url", logger)

	if err == nil {
		t.Errorf("Expected error for invalid URL, but got none")
	}
}

func TestFetchLinks_RetryLogic(t *testing.T) {
	setup()
	attempts := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusInternalServerError) // Simulate server error
	}))
	defer ts.Close()

	_, err := fetcherInstance.FetchLinks(ts.URL, logger)

	if err == nil {
		t.Errorf("Expected error after retries, but got none")
	}

	if attempts != fetcher.MaxRetry {
		t.Errorf("Expected %d retry attempts, but got %d", fetcher.MaxRetry, attempts)
	}
}

// TestFetchLinks_SpecialCharacterURL tests handling of URLs with special characters
func TestFetchLinks_SpecialCharacterURL(t *testing.T) {
	setup()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<a href="https://example.com/path%20with%20spaces">Special Character Link</a>`))
	}))
	defer ts.Close()

	links, err := fetcherInstance.FetchLinks(ts.URL, logger)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedLink := "https://example.com/path%20with%20spaces"
	if _, exists := links[expectedLink]; !exists {
		t.Errorf("Expected link %s in fetched links, but it was not found", expectedLink)
	}
}

// TestFetchLinks_NonHTMLResponse tests handling of non-HTML responses
func TestFetchLinks_NonHTMLResponse(t *testing.T) {
	setup()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("%PDF-1.4"))
	}))
	defer ts.Close()

	links, err := fetcherInstance.FetchLinks(ts.URL, logger)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(links) != 0 {
		t.Errorf("Expected no links from non-HTML response, but got %d", len(links))
	}
}
