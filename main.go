package main

import (
	"fmt"
	"os"

	app "github.com/arnaud-deprez/gsemver/cmd"
)

// Entrypoint for gsemver command
func main() {
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
