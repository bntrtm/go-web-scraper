package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

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

func getHTML(rawURL string) (string, error) {
	var buf io.Reader
	client := http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest(http.MethodGet, rawURL, buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "text/html; charset=utf-8")
	req.Header.Set("User-Agent", "Crawler/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("status code %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		return "", fmt.Errorf("expected content-type text/html, got %s", ct)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	html := string(bytes)
	return html, nil
}

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		log.Fatal(err)
	}
	baseHostname := baseURL.Hostname()
	curURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		log.Fatal(err)
	}
	currentHostname := curURL.Hostname()
	if baseHostname != currentHostname {
		return
	}

	normCurrent, err := normalizeURL(rawCurrentURL)
	if err != nil {
		log.Fatal(err)
	}
	if _, ok := pages[normCurrent]; ok {
		pages[normCurrent] += 1
		return
	} else {
		pages[normCurrent] = 1
	}
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		log.Printf("error: %s\n", err)
	}
	fmt.Println("______________________________")
	fmt.Printf("SCRAPING LINKS FROM HTML: %s\n", rawCurrentURL)
	fmt.Println("______________________________")
	// fmt.Println(html)
	urls, err := getURLsFromHTML(html, curURL)
	if err != nil {
		log.Fatal(err)
	}
	for _, url := range urls {
		crawlPage(rawBaseURL, url, pages)
	}
}
