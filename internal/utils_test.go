package main

import (
	"testing"
)

/*
Testing the fetch of google.com. Also tests checkInternal()
*/
func TestDomain(t *testing.T) {
	url, err := Domain("http://monzo.com/")
	assertEqual(t, err, nil, "")
	assertEqual(t, url, "http://monzo.com", "")

	url, err = Domain("monzo.com")
	assertEqual(t, err, nil, "")
	assertEqual(t, url, "http://monzo.com", "")
}
