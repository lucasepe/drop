package fileserver

import (
	"net/http"
	"os"
	"strings"
)

type filteredFile struct {
	http.File
}

func (f filteredFile) Readdir(n int) ([]os.FileInfo, error) {
	files, err := f.File.Readdir(n)
	if err != nil {
		return nil, err
	}

	allowed := []os.FileInfo{}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		allowed = append(allowed, file)
	}

	return allowed, nil
}

type filteredFileSystem struct {
	http.FileSystem
}

func (fs *filteredFileSystem) Open(name string) (http.File, error) {
	if containsDotFile(name) {
		return nil, os.ErrNotExist
	}

	file, err := fs.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}

	return filteredFile{file}, err
}

func (fs *filteredFileSystem) OpenWithStat(name string) (http.File, os.FileInfo, error) {
	file, err := fs.Open(name)
	if err != nil {
		return nil, nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}

	if (fileInfo.Mode() & os.ModeSymlink) != 0 {
		return nil, nil, os.ErrNotExist
	}

	return filteredFile{file}, fileInfo, nil
}

func containsDotFile(name string) bool {
	for file := range strings.SplitSeq(name, "/") {
		if strings.HasPrefix(file, ".") {
			return true
		}
	}

	return false
}
