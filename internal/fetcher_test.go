package main

import (
	"fmt"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

/*
Standard assert Equals function
*/
func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

/*
Testing the fetch of google.com. Also tests checkInternal()
*/
func TestFetch(t *testing.T) {

	links := fetchLinks("http://google.com/", "http://google.com")
	assertEqual(t, len(links), 7, "")

	links = fetchLinks("http://monzo.com/", "http://monzo.com")
	assertEqual(t, len(links), 42, "")
}

/*
Testing the extraction of all links.
*/
func TestExtract(t *testing.T) {
	res, err := Request("http://google.com")
	if err != nil {
		t.Fatal(err)
	}
	doc, _ := goquery.NewDocumentFromResponse(res)
	links := extractLinks(doc)
	assertEqual(t, len(links), 29, "")

	res, err = Request("http://monzo.com")
	if err != nil {
		t.Fatal(err)
	}
	doc, _ = goquery.NewDocumentFromResponse(res)
	links = extractLinks(doc)
	assertEqual(t, len(links), 53, "")
}
