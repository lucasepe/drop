package fileserver

import "net/http"

type contentTypeWrapper struct {
	http.ResponseWriter
	overrideContentType string
}

func (w *contentTypeWrapper) WriteHeader(statusCode int) {
	if w.overrideContentType != "" {
		w.Header().Set("Content-Type", w.overrideContentType)
	}
	w.ResponseWriter.WriteHeader(statusCode)
}
