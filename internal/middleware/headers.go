package middleware

import (
	"net/http"
)

func Headers(headers map[string]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		fn := func(wri http.ResponseWriter, req *http.Request) {
			for key, value := range headers {
				wri.Header().Set(key, value)
			}

			next.ServeHTTP(wri, req)
		}

		return http.HandlerFunc(fn)
	}
}
