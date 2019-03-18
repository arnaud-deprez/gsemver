package version

import (
	"github.com/arnaud-deprez/gsemver/pkg/git"
)

// GitRepo defines common git actions used by gsemver
//go:generate mockgen -destination mock/git_repo.go github.com/arnaud-deprez/gsemver/pkg/version GitRepo
type GitRepo interface {
	// GetCommits return the list of commits between 2 revisions.
	// If no revision is provided, it does from beginning to HEAD
	GetCommits(from string, to string) ([]git.Commit, error)
	// CountCommits is similar to GetCommits but it just return the number of commits
	CountCommits(from string, to string) (int, error)
	// GetLastRelativeTag gives the last ancestor tag from HEAD
	GetLastRelativeTag(rev string) (git.Tag, error)

	// GetSymbolicRef read symbolic refs
	GetSymbolicRef(name string, short bool) (string, error)
}
