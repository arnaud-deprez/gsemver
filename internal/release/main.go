package main

import (
	"fmt"
	"os"

	"github.com/arnaud-deprez/gsemver/internal/git"
	"github.com/arnaud-deprez/gsemver/internal/log"
	"github.com/arnaud-deprez/gsemver/pkg/version"
)

// Entrypoint for gsemver command
func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Error("Unable to retrieve working directory: %v", err)
		os.Exit(1)
	}

	gitRepo := git.NewVersionGitRepo(dir)
	bumper := version.NewConventionalCommitBumpStrategyOptions(gitRepo)
	version, err := bumper.Bump()

	if err != nil {
		log.Error("Cannot bump version caused by: %v", err)
		os.Exit(1)
	}
	fmt.Println(version)
}
