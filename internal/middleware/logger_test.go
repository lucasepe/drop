package middleware

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLoggerMiddleware(t *testing.T) {
	buf := bytes.Buffer{}

	// Create a simple handler that uses the logger.
	sillyHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("Processing a lot...")
			time.Sleep(1 * time.Second)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello, World!"))

			log.Println("Done!")
		})

	route := Chain(sillyHandler, Logger())

	// Create a test request.
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	// Serve the request.
	route.ServeHTTP(rec, req)

	// Check the log output.
	fmt.Println(buf.String())
}
