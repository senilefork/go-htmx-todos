package cmd

import (
	"fmt"
	"os"
)

// Execute dispatches to the named subcommand.
func Execute(name string) {
	switch name {
	case "web":
		web()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %q (available: web)\n", name)
		os.Exit(1)
	}
}
