package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lucasepe/x/config"
)

func TestBasicAuth(t *testing.T) {

	const (
		sample = `
# printf "ciccio:$(openssl passwd -5 -salt 'nntldico' '123456')\n"
ciccio: $5$nntldico$ptvLNSpQq7lzE2qPinKdtAE1T/pUVvvgndAn57Wv8q3
		`
	)

	users, err := config.Parse(strings.NewReader(sample))
	if err != nil {
		t.Fatal(err)
	}

	sillyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello!")
	})

	chain := Chain(sillyHandler, BasicAuth(users))

	server := httptest.NewServer(chain)
	defer server.Close()

	tests := []struct {
		name         string
		username     string
		password     string
		wantStatus   int
		wantResponse string
	}{
		{
			name:         "Auth OK",
			username:     "ciccio",
			password:     "123456",
			wantStatus:   http.StatusOK,
			wantResponse: "Hello!",
		},
		{
			name:       "Wrong Password",
			username:   "ciccio",
			password:   "wrongpassword",
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "Wrong Username",
			username:   "wrongusername",
			password:   "123456",
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "No auth",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", server.URL, nil)

			if tt.username != "" || tt.password != "" {
				req.SetBasicAuth(tt.username, tt.password)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("got status code [%d], expected [%d]", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}
