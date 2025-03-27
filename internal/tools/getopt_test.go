package tools_test

import (
	"testing"

	"github.com/lucasepe/drop/internal/tools"
	"github.com/lucasepe/x/getopt"
)

func TestStr(t *testing.T) {
	tests := []struct {
		name     string
		opts     []getopt.OptArg
		lookup   []string
		fallback string
		expected string
	}{
		{
			name: "Option found",
			opts: []getopt.OptArg{
				{Option: "-a", Argument: "valueA"},
				{Option: "-b", Argument: "valueB"},
			},
			lookup:   []string{"-b"},
			fallback: "default",
			expected: "valueB",
		},
		{
			name: "Option not found, fallback used",
			opts: []getopt.OptArg{
				{Option: "-a", Argument: "valueA"},
			},
			lookup:   []string{"-x"},
			fallback: "default",
			expected: "default",
		},
		{
			name:     "Empty opts list, fallback used",
			opts:     []getopt.OptArg{},
			lookup:   []string{"-a"},
			fallback: "default",
			expected: "default",
		},
		{
			name: "Multiple matches, first one returned",
			opts: []getopt.OptArg{
				{Option: "-a", Argument: "first"},
				{Option: "-b", Argument: "second"},
				{Option: "-c", Argument: "third"},
			},
			lookup:   []string{"-b", "-c"},
			fallback: "default",
			expected: "second",
		},
		{
			name: "Lookup list empty, fallback used",
			opts: []getopt.OptArg{
				{Option: "-a", Argument: "valueA"},
			},
			lookup:   []string{},
			fallback: "default",
			expected: "default",
		},
		{
			name: "Option found but empty value, fallback used",
			opts: []getopt.OptArg{
				{Option: "-a", Argument: ""},
			},
			lookup:   []string{"-a"},
			fallback: "default",
			expected: "default",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tools.Str(tc.opts, tc.lookup, tc.fallback)
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestFindOptVal(t *testing.T) {
	tests := []struct {
		name     string
		opts     []getopt.OptArg
		lookup   []string
		expected string
	}{
		{
			name: "Option found",
			opts: []getopt.OptArg{
				{Option: "-a", Argument: "valueA"},
				{Option: "-b", Argument: "valueB"},
			},
			lookup:   []string{"-b"},
			expected: "valueB",
		},
		{
			name: "Option not found",
			opts: []getopt.OptArg{
				{Option: "-a", Argument: "valueA"},
			},
			lookup:   []string{"-x"},
			expected: "",
		},
		{
			name:     "Empty opts list",
			opts:     []getopt.OptArg{},
			lookup:   []string{"-a"},
			expected: "",
		},
		{
			name: "Multiple matches, first one returned",
			opts: []getopt.OptArg{
				{Option: "-a", Argument: "first"},
				{Option: "-b", Argument: "second"},
				{Option: "-c", Argument: "third"},
			},
			lookup:   []string{"-b", "-c"},
			expected: "second",
		},
		{
			name: "Lookup list empty",
			opts: []getopt.OptArg{
				{Option: "-a", Argument: "valueA"},
			},
			lookup:   []string{},
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tools.FindOptVal(tc.opts, tc.lookup)
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}
