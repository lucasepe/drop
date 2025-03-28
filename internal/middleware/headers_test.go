package middleware

import (
	"context"
	"io"
	"log"
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

func TestApplyTemplate(t *testing.T) {
	log.SetOutput(io.Discard)

	mockCtx := context.WithValue(context.TODO(),
		contextKeyExtra,
		map[string]any{
			serverAddrKey:  "localhost:8080",
			remoteAddrKey:  "192.168.1.100",
			requestPathKey: "/test",
			userAgentKey:   "ciccio/v0.5.0",
		})

	tests := []struct {
		name     string
		ctx      context.Context
		input    string
		expected string
	}{
		{
			name:     "Plain string, no template",
			ctx:      mockCtx,
			input:    "Hello, World!",
			expected: "Hello, World!",
		},
		{
			name:     "Simple template replacement",
			ctx:      mockCtx,
			input:    "Server: ${ SERVER_ADDR }",
			expected: "Server: localhost:8080",
		},
		{
			name:     "Multiple template variables",
			ctx:      mockCtx,
			input:    "Remote: ${ REMOTE_ADDR }, Path: ${REQUEST_PATH }",
			expected: "Remote: 192.168.1.100, Path: /test",
		},
		{
			name:     "Invalid template syntax",
			ctx:      mockCtx,
			input:    "Server: {{ SERVER_ADDR }}",
			expected: "Server: {{ SERVER_ADDR }}",
		},
		{
			name:     "Context missing data",
			ctx:      context.TODO(),
			input:    "Server: ${SERVER_ADDR}",
			expected: "Server: ${SERVER_ADDR}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := applyTemplate(tt.ctx, tt.input)
			if output != tt.expected {
				t.Errorf("expected: %q, got: %q", tt.expected, output)
			}
		})
	}
}
