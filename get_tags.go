package main

import (
	"strings"

	gq "github.com/PuerkitoBio/goquery"
)

// getFirstHTMLTag returns the first instance of a tag
// with the given type found within the given HTML string,
// OR an empty string, if no such tag is found.
func getFirstHTMLTag(tagType string, html string) (string, error) {
	doc, err := gq.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}
	node := doc.Find(tagType).First()
	return node.Text(), nil
}
