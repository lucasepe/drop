package fileserver

import (
	"errors"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestFilteredFile_Readdir(t *testing.T) {
	mockFiles := []os.FileInfo{
		mockFileInfo{name: "visible.txt"},
		mockFileInfo{name: ".hidden"},
		mockFileInfo{name: "public.log"},
		mockFileInfo{name: ".config"},
	}

	mock := &mockFile{files: mockFiles}
	filtered := filteredFile{File: mock}

	result, err := filtered.Readdir(-1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"visible.txt", "public.log"}
	if len(result) != len(expected) {
		t.Errorf("expected %d files, got %d", len(expected), len(result))
	}

	for i, file := range result {
		if file.Name() != expected[i] {
			t.Errorf("expected file %q, got %q", expected[i], file.Name())
		}
	}

	// Test case where Readdir returns an error
	errorMock := &mockFile{errorOnRead: true}
	filteredError := filteredFile{File: errorMock}

	_, err = filteredError.Readdir(-1)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestFilteredFileSystem_Open(t *testing.T) {
	mockFS := &mockFileSystem{
		files: map[string]bool{
			"visible.txt": true,
			".hidden":     true,
			"public.log":  true,
		},
	}

	fs := &filteredFileSystem{mockFS}

	_, err := fs.Open("visible.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = fs.Open(".hidden")
	if err == nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected error for hidden file, got nil")
	}

	_, err = fs.Open("nonexistent.txt")
	if err == nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected error for nonexistent file, got nil")
	}
}

func TestFilteredFileSystem_OpenWithStat(t *testing.T) {
	mockFS := &mockFileSystem{
		files: map[string]bool{
			"visible.txt": true,
			".hidden":     true,
			"symlink":     true,
		},
	}

	fs := &filteredFileSystem{mockFS}

	_, fileInfo, err := fs.OpenWithStat("visible.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fileInfo == nil {
		t.Fatalf("expected valid fileInfo, got nil")
	}

	_, _, err = fs.OpenWithStat(".hidden")
	if err == nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected error for hidden file, got nil")
	}

	_, _, err = fs.OpenWithStat("nonexistent.txt")
	if err == nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected error for nonexistent file, got nil")
	}
}

type mockFileInfo struct {
	name string
	mode os.FileMode
}

func (m mockFileInfo) Name() string       { return m.name }
func (m mockFileInfo) Size() int64        { return 0 }
func (m mockFileInfo) Mode() os.FileMode  { return m.mode }
func (m mockFileInfo) ModTime() time.Time { return time.Now() }
func (m mockFileInfo) IsDir() bool        { return false }
func (m mockFileInfo) Sys() any           { return nil }

type mockFile struct {
	files       []os.FileInfo
	errorOnRead bool
	errorOnStat bool
}

func (m *mockFile) Readdir(n int) ([]os.FileInfo, error) {
	if m.errorOnRead {
		return nil, errors.New("mock read error")
	}
	return m.files, nil
}

func (m *mockFile) Close() error                   { return nil }
func (m *mockFile) Read([]byte) (int, error)       { return 0, io.EOF }
func (m *mockFile) Seek(int64, int) (int64, error) { return 0, nil }
func (m *mockFile) Stat() (os.FileInfo, error) {
	if m.errorOnStat {
		return nil, errors.New("mock stat error")
	}
	return mockFileInfo{name: "mock", mode: 0}, nil
}

type mockFileSystem struct {
	files map[string]bool
}

func (m *mockFileSystem) Open(name string) (http.File, error) {
	if _, exists := m.files[name]; !exists {
		return nil, os.ErrNotExist
	}
	return &mockFile{}, nil
}
