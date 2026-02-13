package main

import "net/url"

type PageData struct {
	URL            string
	H1             string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

func extractPageData(html, pageURL string) PageData {
	parsedPageURL, _ := url.Parse(pageURL)
	data := PageData{URL: pageURL}
	data.H1, _ = getFirstHTMLTagContent("h1", html)
	data.FirstParagraph, _ = getFirstHTMLTagContent("p", html)
	data.OutgoingLinks, _ = getURLsFromHTML(html, parsedPageURL)
	data.ImageURLs, _ = getImagesFromHTML(html, parsedPageURL)
	return data
}
