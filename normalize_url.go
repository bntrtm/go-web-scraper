package main

import (
	"fmt"
	"net/url"
	"strings"
)

func normalizeURL(urlString string) (string, error) {
	url, err := url.Parse(urlString)
	if err != nil {
		return "", fmt.Errorf("could not normalize URL")
	}
	norm := url.Host + url.Path
	norm = strings.ToLower(norm)
	norm = strings.TrimSuffix(norm, "/")
	return norm, nil
}
