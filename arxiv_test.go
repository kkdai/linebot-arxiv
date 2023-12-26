package main

import (
	"testing"
)

func TestNormalizeArxivURL(t *testing.T) {
	testCases := []struct {
		inputURL      string
		expectedURL   string
		expectedError bool
	}{
		{"https://arxiv.org/abs/1234.56789", "https://arxiv.org/abs/1234.56789", false},
		{"http://arxiv.org/abs/1234.56789", "https://arxiv.org/abs/1234.56789", false},
		{"https://arxiv.org/pdf/1234.56789", "https://arxiv.org/abs/1234.56789", false},
		{"https://huggingface.co/papers/1234.56789", "https://arxiv.org/abs/1234.56789", false},
		{"https://huggingface.co/papers/1234.56789v2", "https://arxiv.org/abs/1234.56789v2", false},
		{"https://example.com/abs/1234.56789", "", true},
		{"https://arxiv.org/abs/invalid-id", "", true},
	}

	for _, tc := range testCases {
		normalizedURL, err := NormalizeArxivURL(tc.inputURL)
		if tc.expectedError {
			if err == nil {
				t.Errorf("NormalizeArxivURL(%q) expected an error, but got none", tc.inputURL)
			}
		} else {
			if err != nil {
				t.Errorf("NormalizeArxivURL(%q) returned an unexpected error: %v", tc.inputURL, err)
			}
			if normalizedURL != tc.expectedURL {
				t.Errorf("NormalizeArxivURL(%q) = %q, want %q", tc.inputURL, normalizedURL, tc.expectedURL)
			}
		}
	}
}
