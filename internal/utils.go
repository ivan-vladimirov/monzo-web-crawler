package utils

import (
	"net/url"
)

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
