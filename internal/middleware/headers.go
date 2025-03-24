package middleware

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/lucasepe/x/config"
)

func Headers(rules config.Config) func(http.Handler) http.Handler {
	c := &headersHandler{rules: rules}
	return c.Handler
}

type headersHandler struct {
	rules config.Config
}

func (h *headersHandler) Handler(next http.Handler) http.Handler {
	fn := func(wri http.ResponseWriter, req *http.Request) {
		uriPath := strings.TrimPrefix(req.URL.Path, "/")

		cat := ""
		for _, pattern := range h.rules.Categories() {
			match, err := filepath.Match(pattern, uriPath)
			if err != nil {
				log.Printf("error while matching pattern %q vs path %s: %s", pattern, uriPath, err)
				continue
			}

			if match {
				cat = pattern
				break
			}
		}

		// Apply common rules
		for key, value := range h.rules.All("") {
			wri.Header().Set(key, value)
		}

		// Apply specific rules
		if cat != "" {
			for key, value := range h.rules.All(cat) {
				wri.Header().Set(key, value)
			}
		}

		next.ServeHTTP(wri, req)
	}

	return http.HandlerFunc(fn)
}
