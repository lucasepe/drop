package fileserver

import (
	"fmt"
	"testing"
)

func TestHumanReadableSize(t *testing.T) {
	tests := []struct {
		size     int64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1023, "1023 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
		{1099511627776, "1.0 TB"},
		{1125899906842624, "1.0 PB"},
		{1152921504606846976, "1.0 EB"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d bytes", tt.size), func(t *testing.T) {
			got := humanReadableSize(tt.size)
			if got != tt.expected {
				t.Errorf("for size %d, expected %q but got %q", tt.size, tt.expected, got)
			}
		})
	}
}
