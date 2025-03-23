package fileserver

import "net/http"

// Option is the functional option type.
type Option func(*FileServer)

// WithMiddlewares sets the middlewares to apply before serving the request.
func WithMiddlewares(middlewares ...func(http.Handler) http.Handler) Option {
	return func(fs *FileServer) {
		fs.middlewares = middlewares
	}
}
