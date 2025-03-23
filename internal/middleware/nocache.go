package middleware

import "net/http"

// NoCache disable cache
func NoCache() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(wri http.ResponseWriter, req *http.Request) {

			wri.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
			wri.Header().Set("Pragma", "no-cache")
			wri.Header().Set("Expires", "0")

			next.ServeHTTP(wri, req)
		}

		return http.HandlerFunc(fn)
	}
}
