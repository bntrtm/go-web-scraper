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

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() { <-cfg.concurrencyControl }()

	cfg.mu.Lock()
	shouldReturn := false
	if cfg.maxPages != 0 {
		shouldReturn = len(cfg.pages) >= cfg.maxPages
	}
	cfg.mu.Unlock()
	if shouldReturn {
		return
	}

	curURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		log.Fatal(err)
	}
	currentHostname := curURL.Hostname()
	if cfg.baseURL.Hostname() != currentHostname {
		return
	}

	normCurrent, err := normalizeURL(rawCurrentURL)
	if err != nil {
		log.Fatal(err)
	}
	if isFirst := cfg.addPageVisit(normCurrent); !isFirst {
		return
	}
	htmlBody, err := getHTML(rawCurrentURL)
	if err != nil {
		log.Printf("error: %s\n", err)
	}

	pageData := extractPageData(htmlBody, rawCurrentURL)
	cfg.mu.Lock()
	cfg.pages[normCurrent] = pageData
	cfg.mu.Unlock()
	for _, url := range pageData.OutgoingLinks {
		cfg.wg.Go(func() {
			cfg.crawlPage(url)
		})
	}
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	if _, ok := cfg.pages[normalizedURL]; ok {
		return false
	} else {
		cfg.mu.Lock()
		cfg.pages[normalizedURL] = PageData{}
		cfg.mu.Unlock()
		return true
	}
}
