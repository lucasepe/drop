package middleware

import (
	"net/http"
	"strings"
)

// Content Security Policy (CSP) HTTP headers Middleware
// that follows the OWASP recommendations.
// https://owasp.org/www-project-secure-headers/index.htm
func CSP(addr string) func(http.Handler) http.Handler {
	cspHeaderRows := []string{
		"default-src 'self'; img-src 'self';",
		"style-src 'unsafe-inline' http://{ADDR} https://{ADDR} https://cdnjs.cloudflare.com;",
		"script-src 'self' http://{ADDR} https://{ADDR} 'wasm-unsafe-eval';",
		"connect-src 'self' http://{ADDR} https://{ADDR};",
		"font-src 'self' http://{ADDR} https://{ADDR} https://cdnjs.cloudflare.com;",
	}

	return func(next http.Handler) http.Handler {
		cspHeaderVal := strings.ReplaceAll(strings.Join(cspHeaderRows, " "), "{ADDR}", addr)

		fn := func(wri http.ResponseWriter, req *http.Request) {
			wri.Header().Set("Content-Security-Policy", cspHeaderVal)
			wri.Header().Set("X-Content-Type-Options", "nosniff")
			wri.Header().Set("X-Frame-Options", "DENY")
			wri.Header().Set("X-XSS-Protection", "1; mode=block")
			wri.Header().Set("Referrer-Policy", "no-referrer")

			next.ServeHTTP(wri, req)
		}

		return http.HandlerFunc(fn)
	}
}
