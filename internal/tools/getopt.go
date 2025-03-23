package tools

import (
	"strconv"
	"strings"

	"github.com/lucasepe/x/getopt"
)

func Str(opts []getopt.OptArg, lookup []string, defaultValue ...string) string {
	val := findOptVal(opts, lookup)
	if val == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return val
}

func True(opts []getopt.OptArg, lookup []string) bool {
	val := findOptVal(opts, lookup)
	val = strings.ToUpper(strings.TrimSpace(val))
	switch val {
	case "1",
		"ENABLE", "ENABLED",
		"POSITIVE",
		"T", "TRUE",
		"Y", "YES":
		return true
	}
	return false
}

func Int(opts []getopt.OptArg, lookup []string, defaultValue int) int {
	val := findOptVal(opts, lookup)
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return i
}

func Float64(opts []getopt.OptArg, lookup []string, defaultValue float64) float64 {
	val := findOptVal(opts, lookup)
	f64, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return defaultValue
	}
	return f64
}

func Has(opts []getopt.OptArg, lookup []string) bool {
	for _, opt := range opts {
		if contains(lookup, opt.Opt()) {
			return true
		}
	}

	return false
}

func findOptVal(opts []getopt.OptArg, lookup []string) (val string) {
	for _, opt := range opts {
		if contains(lookup, opt.Opt()) {
			val = opt.Argument
			break
		}
	}

	return
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
