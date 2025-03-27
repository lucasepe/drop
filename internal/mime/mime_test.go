package mime_test

import (
	"testing"

	"github.com/lucasepe/drop/internal/mime"
)

func TestTypeByExtension(t *testing.T) {
	const (
		fallbackMimeType = "application/octet-stream"
	)

	tests := []struct {
		name     string
		ext      string
		expected string
	}{
		{
			name:     "Valid extension",
			ext:      ".txt",
			expected: "text/plain",
		},
		{
			name:     "Valid extension, case insensitive",
			ext:      ".JPG",
			expected: "image/jpeg",
		},
		{
			name:     "Empty extension",
			ext:      "",
			expected: fallbackMimeType,
		},
		{
			name:     "Extension with leading dot",
			ext:      ".html",
			expected: "text/html",
		},
		{
			name:     "Extension with mixed case",
			ext:      ".PnG",
			expected: "image/png",
		},
		{
			name:     "Unknown extension",
			ext:      ".aaa",
			expected: fallbackMimeType,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := mime.TypeByExtension(tc.ext)
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}
