package fileserver

import (
	"fmt"
)

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

func humanReadableSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
