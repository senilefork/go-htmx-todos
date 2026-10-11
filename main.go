// main.go is deliberately tiny: it just dispatches to a subcommand in
// cmd/ (the starter does this so `go run . seed`, `go run . web`, etc.
// share one binary). We only have one command, but keep the shape.
package main

import (
	"os"

	"go-htmx-jobs/cmd"
)

func main() {
	name := "web"
	if len(os.Args) > 1 {
		name = os.Args[len(os.Args)-1]
	}
	cmd.Execute(name)
}
