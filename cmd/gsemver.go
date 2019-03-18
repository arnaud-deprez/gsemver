package main

import (
	"fmt"
	"os"

	"github.com/arnaud-deprez/gsemver/cmd/app"
)

// Entrypoint for jx command
func main() {
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
