package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestInfo(t *testing.T) {
	sillyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, ok := getExtraInfo(r.Context())

		assert.True(t, ok)
		assert.Equal(t, "example.com:8080", data[serverAddrKey])
		assert.Equal(t, "/test", data[requestPathKey])
		assert.Equal(t, "TestClient/1.0", data[userAgentKey])
	})

	req := httptest.NewRequest("GET", "http://example.com:8080/test", nil)
	req.Header.Set("User-Agent", "TestClient/1.0")

	rec := httptest.NewRecorder()

	server := Chain(sillyHandler, Extra())
	server.ServeHTTP(rec, req)

	res := rec.Result()

	assert.Equal(t, res.StatusCode, http.StatusOK)
}
