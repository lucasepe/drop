package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/lucasepe/drop/internal/cmd/hostdir"
	"github.com/lucasepe/x/getopt"
)

var (
	BuildKey = buildKey{}
)

type Action int

const (
	NoAction Action = iota
	HostDir
	ShowHelp
	ShowVersion
)

func Run(ctx context.Context) error {
	leftOvers, opts, err := getopt.GetOpt(
		os.Args[1:],
		"a:c:k:",
		[]string{
			"help",
			"version",
		},
	)
	if err != nil {
		return err
	}

	op := chosenAction(opts)
	if (op == ShowHelp) || (op == NoAction) {
		usage(os.Stdout)
		return nil
	}

	if op == ShowVersion {
		bld := ctx.Value(BuildKey).(string)
		fmt.Fprintf(os.Stdout, "%s - build: %s\n", appName, bld)
		return nil
	}

	switch op {
	case HostDir:
		err = hostdir.Do(leftOvers, opts)
	}

	return err
}

func chosenAction(opts []getopt.OptArg) Action {
	for _, opt := range opts {
		switch opt.Opt() {
		case "--help":
			return ShowHelp
		case "--version":
			return ShowVersion
		}
	}

	return HostDir
}

type (
	buildKey struct{}
)
