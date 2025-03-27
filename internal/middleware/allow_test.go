package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAllowedMethods(t *testing.T) {
	sillyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name           string
		method         string
		expectedCode   int
		expectedHeader string
	}{
		{
			name:           "GET method",
			method:         http.MethodGet,
			expectedCode:   http.StatusOK,
			expectedHeader: "",
		},
		{
			name:           "HEAD method",
			method:         http.MethodHead,
			expectedCode:   http.StatusOK,
			expectedHeader: "",
		},
		{
			name:           "OPTIONS method",
			method:         http.MethodOptions,
			expectedCode:   http.StatusNoContent,
			expectedHeader: "GET, HEAD, OPTIONS",
		},
		{
			name:           "POST method (method not allowed)",
			method:         http.MethodPost,
			expectedCode:   http.StatusMethodNotAllowed,
			expectedHeader: "",
		},
		{
			name:           "DELETE method (method not allowed)",
			method:         http.MethodDelete,
			expectedCode:   http.StatusMethodNotAllowed,
			expectedHeader: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "http://example.com", nil)
			rec := httptest.NewRecorder()

			middleware := AllowedMethods()
			middleware(sillyHandler).ServeHTTP(rec, req)

			if rec.Code != tc.expectedCode {
				t.Errorf("expected status code %d, got %d", tc.expectedCode, rec.Code)
			}

			if tc.expectedHeader != "" {
				allowHeader := rec.Header().Get("Allow")
				if allowHeader != tc.expectedHeader {
					t.Errorf("expected Allow header %s, got %s", tc.expectedHeader, allowHeader)
				}
			}
		})
	}
}
