package git

import (
	"github.com/arnaud-deprez/gsemver/pkg/version"
)

func NewVersionGitRepo(dir string) version.GitRepo {
	return &gitRepoCLI{
		dir:          dir,
		commitParser: &commitParser{logFormat: logFormat},
	}
}
