package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(wri http.ResponseWriter, req *http.Request) {
			startTime := time.Now()

			rec := &statusRecorder{ResponseWriter: wri, statusCode: http.StatusOK}
			next.ServeHTTP(rec, req)

			elapsedTime := time.Since(startTime)
			log.Printf("%d %s %s [%s]", rec.statusCode, req.Method, req.URL.Path, elapsedTime)
		}

		return http.HandlerFunc(fn)
	}
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	return r.ResponseWriter.Write(b)
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
