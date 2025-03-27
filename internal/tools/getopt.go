package tools

import (
	"slices"

	"github.com/lucasepe/x/getopt"
)

func Str(opts []getopt.OptArg, lookup []string, fallback string) string {
	val := FindOptVal(opts, lookup)
	if val == "" {
		return fallback
	}
	return val
}

func FindOptVal(opts []getopt.OptArg, lookup []string) (val string) {
	for _, opt := range opts {
		if slices.Contains(lookup, opt.Opt()) {
			val = opt.Argument
			break
		}
	}

	return
}
