package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Testchain verifies that middleware is executed in the correct order.
func TestChain(t *testing.T) {
	logs := []string{}

	// Define two simple middleware functions that log their execution order
	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logs = append(logs, "mw1 - before")
			next.ServeHTTP(w, r)
			logs = append(logs, "mw1 - after")
		})
	}

	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logs = append(logs, "mw2 - before")
			next.ServeHTTP(w, r)
			logs = append(logs, "mw2 - after")
		})
	}

	// Define a simple handler that writes "hello world"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logs = append(logs, "handler executed")
		w.Write([]byte("hello world"))
	})

	// Chain the middleware
	finalHandler := Chain(handler, mw1, mw2)

	// Perform a test request
	req := httptest.NewRequest("GET", "http://dummy/", nil)
	rec := httptest.NewRecorder()

	finalHandler.ServeHTTP(rec, req)

	// Expected execution order:
	expectedLogs := []string{
		"mw1 - before",
		"mw2 - before",
		"handler executed",
		"mw2 - after",
		"mw1 - after",
	}

	// Validate execution order
	if len(logs) != len(expectedLogs) {
		t.Fatalf("expected [%d] log entries, got: %d", len(expectedLogs), len(logs))
	}

	for i, logEntry := range logs {
		if logEntry != expectedLogs[i] {
			t.Errorf("expected log entry %d to be %q, got %q", i, expectedLogs[i], logEntry)
		}
	}
}
