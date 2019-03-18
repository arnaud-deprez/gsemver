package app

import (
	"github.com/arnaud-deprez/gsemver/pkg/gsemver/cmd"
)

// Run runs the command
func Run() error {
	cmd := cmd.NewDefaultRootCommand()
	return cmd.Execute()
}
