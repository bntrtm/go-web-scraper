package main

import "testing"

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
		wantErr  bool
	}{
		{
			name:     "remove scheme",
			inputURL: "https://blog.tjbainter.com/path",
			expected: "blog.tjbainter.com/path",
			wantErr:  false,
		},
		{
			name:     "remove trailing slashes",
			inputURL: "https://blog.tjbainter.com/path/",
			expected: "blog.tjbainter.com/path",
			wantErr:  false,
		},
		{
			name:     "lowercase capital letters",
			inputURL: "https://BLOG.tjbainter.com/path/",
			expected: "blog.tjbainter.com/path",
			wantErr:  false,
		},
		{
			name:     "handle bad URL",
			inputURL: ":\\badURL",
			expected: "",
			wantErr:  true,
		},
		// add more test cases here
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if tc.wantErr != (err != nil) {
				t.Errorf("Test %v - '%s' FAIL: want error: %v, but got error: %v", i, tc.name, tc.wantErr, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
