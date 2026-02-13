package main

import (
	"reflect"
	"testing"
)

func TestExtractPageData(t *testing.T) {
	cases := []struct {
		name    string
		pageURL string
		html    string
		want    PageData
	}{
		{
			name:    "basic: h1, main paragraph, relative link and img",
			pageURL: "https://blog.tjbainter.com",
			html: `
<html>
  <body>
    <h1>Hello World</h1>
    <main><p>First paragraph inside main.</p></main>
    <a href="/about">About</a>
    <img src="/logo.png" alt="Logo">
  </body>
</html>`,
			want: PageData{
				URL:            "https://blog.tjbainter.com",
				H1:             "Hello World",
				FirstParagraph: "First paragraph inside main.",
				OutgoingLinks:  []string{"https://blog.tjbainter.com/about"},
				ImageURLs:      []string{"https://blog.tjbainter.com/logo.png"},
			},
		},
		{
			name:    "fallback paragraph when no <main>",
			pageURL: "https://blog.tjbainter.com",
			html: `
<html>
  <body>
    <h1>Title</h1>
    <p>Outside paragraph wins.</p>
    <a href="/x">x</a>
    <img src="/img.png">
  </body>
</html>`,
			want: PageData{
				URL:            "https://blog.tjbainter.com",
				H1:             "Title",
				FirstParagraph: "Outside paragraph wins.",
				OutgoingLinks:  []string{"https://blog.tjbainter.com/x"},
				ImageURLs:      []string{"https://blog.tjbainter.com/img.png"},
			},
		},
		{
			name:    "malformed HTML still parsed; absolute link and image",
			pageURL: "https://blog.tjbainter.com",
			html: `
<html body>
  <h1>Messy</h1>
  <a href="https://other.com/path">Other</a>
  <img src="https://cdn.tjbainter.com/banner.jpg">
</html body>`,
			want: PageData{
				URL:            "https://blog.tjbainter.com",
				H1:             "Messy",
				FirstParagraph: "", // no <p> present
				OutgoingLinks:  []string{"https://other.com/path"},
				ImageURLs:      []string{"https://cdn.tjbainter.com/banner.jpg"},
			},
		},
		{
			name:    "no h1 and no paragraph",
			pageURL: "https://blog.tjbainter.com",
			html: `
<html>
  <body>
    <a href="/only-link">Only link</a>
    <img src="/only.png">
  </body>
</html>`,
			want: PageData{
				URL:            "https://blog.tjbainter.com",
				H1:             "",
				FirstParagraph: "",
				OutgoingLinks:  []string{"https://blog.tjbainter.com/only-link"},
				ImageURLs:      []string{"https://blog.tjbainter.com/only.png"},
			},
		},
		{
			name:    "multiple links and images preserve order",
			pageURL: "https://blog.tjbainter.com",
			html: `
<html><body>
  <h1>t</h1>
  <main><p>p</p></main>
  <a href="/a1">a1</a>
  <a href="https://x.dev/a2">a2</a>
  <img src="/i1.png">
  <img src="https://x.dev/i2.png">
</body></html>`,
			want: PageData{
				URL:            "https://blog.tjbainter.com",
				H1:             "t",
				FirstParagraph: "p",
				OutgoingLinks: []string{
					"https://blog.tjbainter.com/a1",
					"https://x.dev/a2",
				},
				ImageURLs: []string{
					"https://blog.tjbainter.com/i1.png",
					"https://x.dev/i2.png",
				},
			},
		},
		{
			name:    "invalid base URL → empty link/image slices",
			pageURL: `:\\invalidBaseURL`,
			html: `
<html>
  <body>
    <h1>Title</h1>
    <p>Paragraph</p>
    <a href="/path">path</a>
    <img src="/logo.png">
  </body>
</html>`,
			want: PageData{
				URL:            `:\\invalidBaseURL`,
				H1:             "Title",
				FirstParagraph: "Paragraph",
				OutgoingLinks:  nil,
				ImageURLs:      nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := extractPageData(tc.html, tc.pageURL)

			if got.URL != tc.want.URL {
				t.Errorf("URL: want %q, got %q", tc.want.URL, got.URL)
			}
			if got.H1 != tc.want.H1 {
				t.Errorf("H1: want %q, got %q", tc.want.H1, got.H1)
			}
			if got.FirstParagraph != tc.want.FirstParagraph {
				t.Errorf("FirstParagraph: want %q, got %q", tc.want.FirstParagraph, got.FirstParagraph)
			}
			if !reflect.DeepEqual(got.OutgoingLinks, tc.want.OutgoingLinks) {
				t.Errorf("OutgoingLinks: want %v, got %v", tc.want.OutgoingLinks, got.OutgoingLinks)
			}
			if !reflect.DeepEqual(got.ImageURLs, tc.want.ImageURLs) {
				t.Errorf("ImageURLs: want %v, got %v", tc.want.ImageURLs, got.ImageURLs)
			}
		})
	}
}
