package git

import (
	"github.com/arnaud-deprez/gsemver/pkg/version"
)

// NewVersionGitRepo creates a version.GitRepo instance for a directory
func NewVersionGitRepo(dir string) version.GitRepo {
	return &gitRepoCLI{
		dir:          dir,
		commitParser: &commitParser{logFormat: logFormat},
	}
}
