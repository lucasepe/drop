package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lucasepe/x/config"
)

func TestHeadersMiddleware(t *testing.T) {
	const (
		sample = `
X-Greetings: Hello World!

[*.mod]
X-Type: Go Module File

[[]
Wrong: Rule
`
	)

	tests := []struct {
		name           string
		urlPath        string
		expectedHeader map[string]string
	}{
		{
			name:    "Global header applied",
			urlPath: "/any/path",
			expectedHeader: map[string]string{
				"X-Greetings": "Hello World!",
			},
		},
		{
			name:    "Pattern-specific header applied",
			urlPath: "/module.mod",
			expectedHeader: map[string]string{
				"X-Greetings": "Hello World!",
				"X-Type":      "Go Module File",
			},
		},
		{
			name:    "No pattern match, only global headers",
			urlPath: "/other.txt",
			expectedHeader: map[string]string{
				"X-Greetings": "Hello World!",
			},
		},
	}

	rules, err := config.Parse(strings.NewReader(sample))
	if err != nil {
		t.Fatal(err)
	}

	sillyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := Chain(sillyHandler, Headers(rules))

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.urlPath, nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			for key, expectedValue := range tc.expectedHeader {
				if got := rec.Header().Get(key); got != expectedValue {
					t.Errorf("header %q: got %q, expected %q", key, got, expectedValue)
				}
			}
		})
	}
}
