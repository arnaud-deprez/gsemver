package main

import (
	"fmt"
	"os"

	"github.com/arnaud-deprez/gsemver/cmd/gsemver"
)

// Entrypoint for gsemver command
func main() {
	if err := gsemver.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
