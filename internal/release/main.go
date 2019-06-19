package main

import (
	"fmt"
	"os"

	"github.com/arnaud-deprez/gsemver/internal/git"
	"github.com/arnaud-deprez/gsemver/internal/log"
	"github.com/arnaud-deprez/gsemver/pkg/version"
)

// Use to compte the next version of gsemver itself
func main() {
	gitRepo := git.NewDefaultVersionGitRepo()
	bumper := version.NewConventionalCommitBumpStrategyOptions(gitRepo)
	version, err := bumper.Bump()

	if err != nil {
		log.Error("Cannot bump version caused by: %v", err)
		os.Exit(1)
	}
	fmt.Println(version)
}
