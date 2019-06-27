package git

import (
	"os"

	"github.com/arnaud-deprez/gsemver/internal/log"
	"github.com/arnaud-deprez/gsemver/pkg/version"
)

// NewVersionGitRepo creates a version.GitRepo instance for a directory
func NewVersionGitRepo(dir string) version.GitRepo {
	return &gitRepoCLI{
		dir:          dir,
		commitParser: &commitParser{logFormat: logFormat},
	}
}

// NewDefaultVersionGitRepo creates a version.GitRepo instance with current working dir
func NewDefaultVersionGitRepo() version.GitRepo {
	dir, err := os.Getwd()
	if err != nil {
		log.Error("Unable to retrieve working directory: %v", err)
		os.Exit(1)
	}

	return NewVersionGitRepo(dir)
}
