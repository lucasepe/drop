package middleware

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/lucasepe/x/text"
)

func Debug() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(wri http.ResponseWriter, req *http.Request) {

			dump, err := httputil.DumpRequest(req, true)
			if err == nil {
				info := bytes.ReplaceAll(dump, []byte("\r\n"), []byte("\n"))
				log.Print(fmt.Sprintf("[HTTP DEBUG] Request:\n%s\n",
					text.IndentBytes(info, []byte("  >>> "))))
			}

			// Wrappa ResponseWriter per intercettare la response
			rec := &responseRecorder{ResponseWriter: wri, body: &bytes.Buffer{}, statusCode: http.StatusOK}
			next.ServeHTTP(rec, req)

			log.Printf("[HTTP DEBUG] Response:\n"+
				"📌 Status: %d\n"+
				"📋 Headers: %v\n"+
				"📝 Body: %s\n",
				rec.statusCode, rec.Header(),
				strings.ReplaceAll(rec.body.String(), "\r\n", "\n"))
		}

		return http.HandlerFunc(fn)
	}
}

type responseRecorder struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)                  // Scrive nel buffer
	return r.ResponseWriter.Write(b) // Scrive nella response originale
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
