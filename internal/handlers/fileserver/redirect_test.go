package fileserver

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestRedirectTo(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		query    string
		expected string
	}{
		{"No Query", "/path", "", "/path"},
		{"With Query", "/path", "foo=bar", "/path?foo=bar"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{URL: &url.URL{Path: tt.url, RawQuery: tt.query}}
			rw := httptest.NewRecorder()

			redirectTo(rw, req, tt.url)

			if location := rw.Header().Get("Location"); location != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, location)
			}

			if status := rw.Code; status != http.StatusMovedPermanently {
				t.Errorf("expected status %d, got %d", http.StatusMovedPermanently, status)
			}
		})
	}
}
