package fileserver

import (
	_ "embed"
	"fmt"
	"net/http"
)

func redirectTo(rw http.ResponseWriter, req *http.Request, path string) {
	if query := req.URL.RawQuery; query != "" {
		path += fmt.Sprint("?", query)
	}

	rw.Header().Add("Location", path)
	rw.WriteHeader(http.StatusMovedPermanently)
}
