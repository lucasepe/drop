package fileserver

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/lucasepe/drop/internal/middleware"
	"github.com/lucasepe/drop/internal/mime"
)

var (
	//go:embed assets/index.html
	defaultAutoIndexTemplate string
)

type FileServer struct {
	autoIndexTmpl *template.Template
	fileSystem    filteredFileSystem
	middlewares   []func(http.Handler) http.Handler
}

// New returns a new handler instance that serves HTTP requests with the contents of the given directory.
func New(dir string, options ...Option) http.Handler {

	autoIndexTmpl := template.New("autoIndex").Funcs(template.FuncMap{
		"humanReadableSize": humanReadableSize,
	})

	fs := &FileServer{
		fileSystem:    filteredFileSystem{http.Dir(dir)},
		autoIndexTmpl: template.Must(autoIndexTmpl.Parse(defaultAutoIndexTemplate)),
	}

	for _, option := range options {
		option(fs)
	}

	return middleware.Chain(fs, fs.middlewares...)
}

// ServeHTTP responds to an HTTP request.
func (fs *FileServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	p := path.Clean(req.URL.Path)

	file, fileInfo, err := fs.fileSystem.OpenWithStat(p)
	if err != nil {
		fs.handleError(rw, err)
		return
	}

	defer func() { _ = file.Close() }()

	fs.serveContent(rw, req, file, fileInfo)
}

func (fs *FileServer) handleError(rw http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError

	if os.IsNotExist(err) || os.IsPermission(err) {
		statusCode = http.StatusNotFound
	}

	rw.WriteHeader(statusCode)

	errorPageName := fmt.Sprintf("%d.html", statusCode)
	if file, err := fs.fileSystem.Open(errorPageName); err == nil {
		defer func() { _ = file.Close() }()

		_, _ = io.Copy(rw, file)
	}
}

func (fs *FileServer) serveContent(rw http.ResponseWriter, req *http.Request, file http.File, fileInfo os.FileInfo) {
	if !fileInfo.IsDir() {
		mimeType := mime.TypeByExtension(filepath.Ext(fileInfo.Name()))
		//log.Printf("serving file %q (mime-type: %s)\n", fileInfo.Name(), mimeType)

		wrapped := &contentTypeWrapper{
			ResponseWriter:      rw,
			overrideContentType: mimeType,
		}
		http.ServeContent(wrapped, req, fileInfo.Name(), fileInfo.ModTime(), file)
		return
	}

	// enforce trailing slash
	if !strings.HasSuffix(req.URL.Path, "/") {
		redirectTo(rw, req, fmt.Sprint(req.URL.Path, "/"))
		return
	}

	indexFilePath := path.Clean(req.URL.Path + "/index.html")

	indexFile, indexFileInfo, err := fs.fileSystem.OpenWithStat(indexFilePath)
	if err == nil {
		defer func() { _ = indexFile.Close() }()
		http.ServeContent(rw, req, indexFileInfo.Name(), indexFileInfo.ModTime(), indexFile)

		return
	}

	if !os.IsNotExist(err) {
		fs.handleError(rw, err)
		return
	}

	entries, err := file.Readdir(-1)
	if err != nil {
		fs.handleError(rw, err)
		return
	}

	var files []info
	for _, entry := range entries {
		files = append(files, info{
			Name:  entry.Name(),
			Size:  humanReadableSize(entry.Size()),
			IsDir: entry.IsDir(),
		})
	}

	rw.Header().Add("Content-Type", "text/html")

	isSubdir := strings.Trim(req.URL.Path, "/") != ""
	tmp := templateData{
		Files:          files,
		IsSubdirectory: isSubdir,
		CurrentPath:    req.URL.Path,
	}

	_ = fs.autoIndexTmpl.Execute(rw, tmp)
}

func redirectTo(rw http.ResponseWriter, req *http.Request, path string) {
	if query := req.URL.RawQuery; query != "" {
		path += fmt.Sprint("?", query)
	}

	rw.Header().Add("Location", path)
	rw.WriteHeader(http.StatusMovedPermanently)
}
