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

	"github.com/lucasepe/drop/internal/mime"
)

var (
	//go:embed assets/index.html
	defaultAutoIndexTemplate string
)

// New returns a new handler instance that serves HTTP requests with the contents of the given directory.
func New(in http.FileSystem) http.Handler {
	autoIndexTmpl := template.New("autoIndex").Funcs(template.FuncMap{
		"humanReadableSize": humanReadableSize,
	})

	return &fileServer{
		fileSystem:    filteredFileSystem{in},
		autoIndexTmpl: template.Must(autoIndexTmpl.Parse(defaultAutoIndexTemplate)),
	}
}

// ServeHTTP responds to an HTTP request.
func (fs *fileServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	p := path.Clean(req.URL.Path)

	file, fileInfo, err := fs.fileSystem.OpenWithStat(p)
	if err != nil {
		fs.handleError(rw, err)
		return
	}

	defer func() { _ = file.Close() }()

	fs.serveContent(rw, req, file, fileInfo)
}

type templateData struct {
	Files          []info
	IsSubdirectory bool
	CurrentPath    string
}

type info struct {
	Name  string
	Size  string
	IsDir bool
}

type fileServer struct {
	autoIndexTmpl *template.Template
	fileSystem    filteredFileSystem
}

func (fs *fileServer) handleError(rw http.ResponseWriter, err error) {
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

func (fs *fileServer) serveContent(rw http.ResponseWriter, req *http.Request, file http.File, fileInfo os.FileInfo) {
	if !fileInfo.IsDir() {
		mimeType := mime.TypeByExtension(filepath.Ext(fileInfo.Name()))

		if req.Method == http.MethodHead {
			rw.Header().Set("Content-Type", mimeType)
			rw.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
			rw.WriteHeader(http.StatusOK)
			return
		}

		wrapped := &contentTypeWrapper{
			ResponseWriter:      rw,
			overrideContentType: mimeType,
		}

		http.ServeContent(wrapped, req, fileInfo.Name(), fileInfo.ModTime(), file)
		return
	}

	if req.Method == http.MethodHead {
		rw.WriteHeader(http.StatusNoContent)
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
