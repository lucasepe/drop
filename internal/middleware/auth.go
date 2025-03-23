package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/lucasepe/drop/internal/crypto/sha256"
	"github.com/lucasepe/x/config"
)

func BasicAuth(users config.Config) func(http.Handler) http.Handler {
	c := &basicAuthHandler{users: users}
	return c.Handler
}

type basicAuthHandler struct {
	users config.Config
}

func (h *basicAuthHandler) Handler(next http.Handler) http.Handler {
	fn := func(wri http.ResponseWriter, req *http.Request) {
		username, password, ok := req.BasicAuth()
		if !ok {
			wri.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			wri.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !h.isAuthorized(username, password) {
			wri.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(wri, req)
	}

	return http.HandlerFunc(fn)
}

func (h *basicAuthHandler) isAuthorized(username, password string) bool {
	want := h.users.Value("", username)
	if want == "" {
		return false
	}

	idx := strings.LastIndex(want, "$")
	if idx <= 0 {
		log.Printf("invalid hash for user: %s", username)
		return false
	}

	sha256Crypt := sha256.New()
	salt := want[:idx]
	got, err := sha256Crypt.Generate([]byte(password), []byte(salt))
	if err != nil {
		log.Printf("unable to generate sha256 crypt digest for user '%s': %s", username, err)
		return false
	}

	return (got == want)
}

/*
func BasicAuth(users config.Config, realm string) func(http.Handler) http.Handler {

	sha256Crypt := sha256.New()

	return func(next http.Handler) http.Handler {
		fn := func(wri http.ResponseWriter, req *http.Request) {
			usr, pwd, ok := req.BasicAuth()
			if !ok {
				wri.Header().Set("WWW-Authenticate", fmt.Sprintf("Basic realm=%q", realm))
				wri.WriteHeader(http.StatusUnauthorized)
				return
			}

			want, exists := credentials[usr]
			if !exists {
				wri.WriteHeader(http.StatusForbidden)
				return
			}

			idx := strings.LastIndex(want, "$")
			if idx <= 0 {
				log.Printf("invalid hash for user: %s", usr)
				wri.WriteHeader(http.StatusForbidden)
				return
			}

			salt := want[:idx]
			got, err := sha256Crypt.Generate([]byte(pwd), []byte(salt))
			if err != nil {
				log.Printf("unable to generate sha256 crypt digest for user '%s': %s", usr, err)
				wri.WriteHeader(http.StatusForbidden)
				return
			}
			if got != want {
				wri.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(wri, req)
		}

		return http.HandlerFunc(fn)
	}
}
*/
