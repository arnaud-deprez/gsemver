package version

import (
	"github.com/arnaud-deprez/gsemver/pkg/git"
)

// GitRepo defines common git actions used by gsemver
//go:generate mockgen -destination mock/git_repo.go github.com/arnaud-deprez/gsemver/pkg/version GitRepo
type GitRepo interface {
	// FetchTags fetches the tags from remote
	FetchTags() error
	// GetCommits return the list of commits between 2 revisions.
	// If no revision is provided, it does from beginning to HEAD
	GetCommits(from string, to string) ([]git.Commit, error)
	// CountCommits counts the number of commits between 2 revisions.
	CountCommits(from string, to string) (int, error)
	// GetLastRelativeTag gives the last ancestor tag from HEAD
	GetLastRelativeTag(rev string) (git.Tag, error)
	// GetCurrentBranch gives the current branch from HEAD
	GetCurrentBranch() (string, error)
}
