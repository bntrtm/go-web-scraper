package main

import "testing"

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
			actual, err := getFirstHTMLTag(tc.inputTagType, tc.inputHTML)
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
