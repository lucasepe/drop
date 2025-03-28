package middleware

import (
	"context"
	"net/http"
)

func Extra() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(wri http.ResponseWriter, req *http.Request) {

			ctx := context.WithValue(req.Context(),
				contextKeyExtra,
				map[string]any{
					serverAddrKey:  req.Host,
					serverNameKey:  req.URL.Hostname(),
					remoteAddrKey:  req.RemoteAddr,
					requestUriKey:  req.RequestURI,
					requestPathKey: req.URL.Path,
					userAgentKey:   req.UserAgent(),
					refererKey:     req.Referer(),
					originKey:      req.Header.Get("Origin"),
				})

			next.ServeHTTP(wri, req.WithContext(ctx))
		}

		return http.HandlerFunc(fn)

	}
}

func getExtraInfo(ctx context.Context) (map[string]any, bool) {
	info, ok := ctx.Value(contextKeyExtra).(map[string]any)
	if !ok {
		return map[string]any{}, false
	}

	return info, true
}

const (
	serverAddrKey  = "SERVER_ADDR"
	serverNameKey  = "SERVER_NAME"
	remoteAddrKey  = "REMOTE_ADDR"
	requestUriKey  = "REQUEST_URI"
	requestPathKey = "REQUEST_PATH"
	userAgentKey   = "USER_AGENT"
	refererKey     = "REFERER"
	originKey      = "ORIGIN"
)

type contextKey string

func (c contextKey) String() string {
	return "drop." + string(c)
}

var (
	contextKeyExtra = contextKey("extra")
)
