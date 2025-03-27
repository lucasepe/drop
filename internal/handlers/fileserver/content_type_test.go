package fileserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContentTypeWrapper(t *testing.T) {
	tests := []struct {
		overrideContentType string
		statusCode          int
		expectedType        string
	}{
		{"application/json", http.StatusOK, "application/json"},
		{"text/plain", http.StatusNotFound, "text/plain"},
		{"", http.StatusInternalServerError, ""},
	}

	for _, tt := range tests {
		t.Run(tt.overrideContentType, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			wrapper := &contentTypeWrapper{
				ResponseWriter:      recorder,
				overrideContentType: tt.overrideContentType,
			}

			wrapper.WriteHeader(tt.statusCode)

			if got := recorder.Header().Get("Content-Type"); got != tt.expectedType {
				t.Errorf("expected Content-Type %q, got %q", tt.expectedType, got)
			}

			if recorder.Code != tt.statusCode {
				t.Errorf("expected status code %d, got %d", tt.statusCode, recorder.Code)
			}
		})
	}
}
