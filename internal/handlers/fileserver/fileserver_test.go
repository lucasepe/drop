package fileserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lucasepe/x/http/memfs"
	"github.com/stretchr/testify/assert"
)

func TestFileServer(t *testing.T) {
	mockFS, err := memfs.New(
		map[string][]byte{
			"/1.txt":  []byte("Hello!"),
			"/2.html": []byte("<h2>Hello!</h2>"),
			"/3.md":   []byte("## Hello!"),
			"/4.jpg":  []byte("asdasd"),
		})

	if err != nil {
		t.Fatal(err)
	}

	fs := New(mockFS)
	/*
		fs := &fileServer{
			fileSystem: filteredFileSystem{mockFS},
			autoIndexTmpl: template.Must(template.New("autoIndex").
				Parse("<html><body>{{range .Files}}<div>{{.Name}}</div>{{end}}</body></html>")),
		}
	*/

	t.Run("Test AutoIndex", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		fs.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusOK)
		assert.Equal(t, rr.Header().Get("Content-Type"), "text/html")
		assert.Equal(t, rr.Body.Len(), 3085)
	})

	t.Run("Test NotFound", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/nonexistent.html", nil)
		rr := httptest.NewRecorder()

		fs.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusNotFound)
	})

	t.Run("Test Serve Content for File", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/3.md", nil)
		rr := httptest.NewRecorder()

		fs.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusOK)
		assert.Equal(t, rr.Body.String(), "## Hello!")
	})

	//t.Run("Test Error Page Serving", func(t *testing.T) {
	//	fs.handleError(httptest.NewRecorder(), os.ErrNotExist)
	//})
}
