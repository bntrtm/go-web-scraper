package main

import (
	"encoding/csv"
	"os"
	"strings"
)

func writeCSVReport(pages map[string]PageData, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(f)
	err = writer.Write([]string{"page_url", "h1", "first_paragraph", "outgoing_link_urls", "image_urls"})
	if err != nil {
		return err
	}
	for _, page := range pages {
		err = writer.Write([]string{page.URL, page.H1, page.FirstParagraph, strings.Join(page.OutgoingLinks, ";"), strings.Join(page.ImageURLs, ";")})
		if err != nil {
			return err
		}
	}
	return nil
}
