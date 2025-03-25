package middleware

import "net/http"

func AllowedMethods() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(wri http.ResponseWriter, req *http.Request) {
			switch req.Method {
			case http.MethodGet, http.MethodHead:
				next.ServeHTTP(wri, req)
			case http.MethodOptions:
				wri.Header().Set("Allow", "GET, HEAD, OPTIONS")
				wri.WriteHeader(http.StatusNoContent)
			default:
				http.Error(wri, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		}

		return http.HandlerFunc(fn)
	}
}
