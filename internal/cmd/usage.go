package cmd

import (
	"fmt"
	"io"
	"strings"

	xtext "github.com/lucasepe/x/text"
)

const (
	appName = "drop"
)

func usage(wri io.Writer) {
	var (
		desc = []string{
			"Lightweight and secure HTTP server for hosting static files from a specified folder.",
		}

		donateInfo = []string{
			"If you find this tool helpful consider supporting with a donation.",
			"Every bit helps cover development time and fuels future improvements.\n",
			"Your support truly makes a difference — thank you!\n",
			"  * https://www.paypal.com/donate/?hosted_button_id=FV575PVWGXZBY\n",
		}
	)

	fmt.Fprintln(wri)
	fmt.Fprint(wri, "┌┬┐┬─┐┌─┐┌─┐\n")
	fmt.Fprint(wri, " ││├┬┘│ │├─┘\n")
	fmt.Fprint(wri, "─┴┘┴└─└─┘┴ \n")

	fmt.Fprintln(wri)
	for _, el := range desc {
		fmt.Fprintf(wri, "%s\n\n", xtext.Wrap(el, 65))
	}
	fmt.Fprintln(wri)

	fmt.Fprint(wri, "USAGE:\n\n")
	fmt.Fprintf(wri, "  %s [FLAGS] /path/to/dir \n\n", appName)

	fmt.Fprint(wri, "ACTIONS:\n\n")
	fmt.Fprint(wri, "  --help          Display this help and exit.\n")
	fmt.Fprint(wri, "  --version       Output version information and exit.\n")
	fmt.Fprintln(wri)

	fmt.Fprint(wri, "FLAGS:\n\n")
	fmt.Fprint(wri, "  -a <ip:port>    Server address (default: 127.0.0.1:8080).\n")
	fmt.Fprint(wri, "  -c <file>       TLS certificate.\n")
	fmt.Fprint(wri, "  -k <file>       TLS private key.\n")
	fmt.Fprintln(wri)

	fmt.Fprint(wri, "EXAMPLES:\n\n")
	fmt.Fprint(wri, " » Serve the current directory:\n\n")
	fmt.Fprintf(wri, "     %s\n\n", appName)
	fmt.Fprint(wri, " » Serve the '/www/public' directory:\n\n")
	fmt.Fprintf(wri, "     %s /www/public \n\n", appName)

	fmt.Fprint(wri, "SUPPORT:\n\n")
	fmt.Fprint(wri, xtext.Indent(strings.Join(donateInfo, "\n"), "  "))
	fmt.Fprint(wri, "\n\n")

	fmt.Fprintln(wri, "Copyright (c) 2025 Luca Sepe")
}
