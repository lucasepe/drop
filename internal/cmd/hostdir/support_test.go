package hostdir

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_LocalFS(t *testing.T) {
	const data = `
[category1]
key1: value1
key2: value2

[category2]
key1:value3
`
	tmpFile, err := os.CreateTemp("", "config-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err = tmpFile.WriteString(data); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	fsys := os.DirFS(filepath.Dir(tmpFile.Name()))
	filename := filepath.Base(tmpFile.Name())

	cfg, err := loadConfig(fsys, filename)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.Value("category1", "key1") != "value1" {
		t.Errorf("expected value1, got %s", cfg.Value("category1", "key1"))
	}

	if cfg.Value("category1", "key2") != "value2" {
		t.Errorf("expected value2, got %s", cfg.Value("category1", "key2"))
	}

	if cfg.Value("category2", "key1") != "value3" {
		t.Errorf("expected value3, got %s", cfg.Value("category2", "key1"))
	}
}
