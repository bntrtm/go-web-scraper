package main

import (
	"net/url"
	"reflect"
	"strings"
	"testing"
)

const testHTML = `<html>
  <body>
    <h1>This is some HTML</h1>
    <main>
      <p>It has this first paragraph.</p>
      <p>It also has this second paragraph.</p>
    </main>
  </body>
</html>`

func TestGetHTMLTags(t *testing.T) {
	tests := []struct {
		name         string
		inputHTML    string
		inputTagType string
		expected     string
	}{
		{
			name:      "empty string with no tag",
			inputHTML: "",
			expected:  "",
		},
		{
			name:         "extract first <h1> content",
			inputHTML:    testHTML,
			inputTagType: "h1",
			expected:     "This is some HTML",
		},
		{
			name:         "extract first <p>",
			inputHTML:    testHTML,
			inputTagType: "p",
			expected:     "It has this first paragraph.",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getFirstHTMLTagContent(tc.inputTagType, tc.inputHTML)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected content: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetURLsFromHTML(t *testing.T) {
	cases := []struct {
		name          string
		inputURL      string
		inputBody     string
		expected      []string
		errorContains string
	}{
		{
			name:     "absolute URL",
			inputURL: "https://blog.tjbainter.com",
			inputBody: `
<html>
	<body>
		<a href="https://blog.tjbainter.com">
			<span>TJBainter.com</span>
		</a>
	</body>
</html>
`,
			expected: []string{"https://blog.tjbainter.com"},
		},
		{
			name:     "relative URL",
			inputURL: "https://blog.tjbainter.com",
			inputBody: `
<html>
	<body>
		<a href="/path/one">
			<span>TJBainter.com</span>
		</a>
	</body>
</html>
`,
			expected: []string{"https://blog.tjbainter.com/path/one"},
		},
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.tjbainter.com",
			inputBody: `
<html>
	<body>
		<a href="/path/one">
			<span>TJBainter.com</span>
		</a>
		<a href="https://other.com/path/one">
			<span>TJBainter.com</span>
		</a>
	</body>
</html>
`,
			expected: []string{"https://blog.tjbainter.com/path/one", "https://other.com/path/one"},
		},
		{
			name:     "no href",
			inputURL: "https://blog.tjbainter.com",
			inputBody: `
<html>
	<body>
		<a>
			<span>TJBainter.com</span>
		</a>
	</body>
</html>
`,
			expected: nil,
		},
		{
			name:     "bad HTML",
			inputURL: "https://blog.tjbainter.com",
			inputBody: `
<html body>
	<a href="path/one">
		<span>TJBainter.com</span>
	</a>
</html body>
`,
			expected: []string{"https://blog.tjbainter.com/path/one"},
		},
		{
			name:     "invalid href URL",
			inputURL: "https://blog.tjbainter.com",
			inputBody: `
<html>
	<body>
		<a href=":\\invalidURL">
			<span>TJBainter.com</span>
		</a>
	</body>
</html>
`,
			expected: nil,
		},
	}

	for i, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, err := url.Parse(tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: couldn't parse input URL: %v", i, tc.name, err)
				return
			}

			actual, err := getURLsFromHTML(tc.inputBody, baseURL)

			if err != nil && !strings.Contains(err.Error(), tc.errorContains) {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			} else if err != nil && tc.errorContains == "" {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			} else if err == nil && tc.errorContains != "" {
				t.Errorf("Test %v - '%s' FAIL: expected error containing '%v', got none.", i, tc.name, tc.errorContains)
				return
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - '%s' FAIL: expected URLs %v, got URLs %v", i, tc.name, tc.expected, actual)
				return
			}
		})
	}
}

func TestGetImagesFromHTMLRelative(t *testing.T) {
	inputURL := "https://blog.tjbainter.com"
	inputBody := `<html><body><img src="/logo.png" alt="Logo"></body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getImagesFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://blog.tjbainter.com/logo.png"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetImagesFromHTMLMultiple(t *testing.T) {
	inputURL := "https://blog.tjbainter.com"
	inputBody := `<html><body>
		<img src="/logo.png" alt="Logo">
		<img src="https://cdn.tjbainter.com/banner.jpg">
	</body></html>`

	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	actual, err := getImagesFromHTML(inputBody, parsedURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{
		"https://blog.tjbainter.com/logo.png",
		"https://cdn.tjbainter.com/banner.jpg",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}
