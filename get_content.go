package main

import (
	"net/url"
	"strings"

	gq "github.com/PuerkitoBio/goquery"
)

// getFirstHTMLTagContent returns the CONTENT of the first instance of a tag
// with the given type found within the given HTML string,
// OR an empty string, if no such tag is found.
func getFirstHTMLTagContent(tagType string, html string) (string, error) {
	doc, err := gq.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}
	node := doc.Find(tagType).First()
	return node.Text(), nil
}

// getURLsFromHTML returns complete URLs from all anchor tags
// within the given html body.
func getURLsFromHTML(html string, baseURL *url.URL) ([]string, error) {
	if baseURL == nil {
		return nil, nil
	}

	doc, err := gq.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var extractedURLs []string
	doc.Find("a[href]").Each(func(_ int, s *gq.Selection) {
		href, ok := s.Attr("href")
		if !ok {
			return
		}
		href = strings.TrimSpace(href)
		if href == "" {
			return
		}
		parsedHref, err := url.Parse(href)
		if href == "" || err != nil {
			return
		}
		extractedURLs = append(extractedURLs, baseURL.ResolveReference(parsedHref).String())
	})
	return extractedURLs, nil
}

// getImagesFromHTML complete source URLs for all image tags
// within the given html body.
func getImagesFromHTML(html string, baseURL *url.URL) ([]string, error) {
	if baseURL == nil {
		return nil, nil
	}

	doc, err := gq.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var extractedURLs []string
	doc.Find("img[src]").Each(func(_ int, s *gq.Selection) {
		src, ok := s.Attr("src")
		if !ok {
			return
		}
		src = strings.TrimSpace(src)
		if src == "" {
			return
		}
		parsedSrc, err := url.Parse(src)
		if src == "" || err != nil {
			return
		}
		extractedURLs = append(extractedURLs, baseURL.ResolveReference(parsedSrc).String())
	})
	return extractedURLs, nil
}
