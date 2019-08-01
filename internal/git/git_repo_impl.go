package git

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/arnaud-deprez/gsemver/internal/command"
	"github.com/arnaud-deprez/gsemver/internal/log"
	"github.com/arnaud-deprez/gsemver/pkg/git"
	"github.com/arnaud-deprez/gsemver/pkg/version"
)

const (
	gitRepoBranchEnv = "GIT_BRANCH"
)

type gitRepoCLI struct {
	version.GitRepo
	dir          string
	commitParser *commitParser
}

// FetchTags implements version.GitRepo.FetchTags
func (g *gitRepoCLI) FetchTags() error {
	_, err := gitCmd(g).
		WithArgs(
			"fetch",
			"--tags",
		).Run()
	return err
}

// GetCommits implements version.GitRepo.Getcommits
func (g *gitRepoCLI) GetCommits(from string, to string) ([]git.Commit, error) {
	rev := parseRev(from, to)
	out, err := gitCmd(g).
		WithArgs(
			"log",
			rev,
			"--no-decorate",
			"--pretty="+g.commitParser.logFormat,
		).Run()

	if err != nil {
		return nil, err
	}

	return g.commitParser.Parse(out), nil
}

// CountCommits implements version.GitRepo.CountCommits
func (g *gitRepoCLI) CountCommits(from string, to string) (int, error) {
	rev := parseRev(from, to)
	cmd := gitCmd(g).WithArgs("rev-list", "--ancestry-path", "--count", rev)
	out, err := cmd.Run()
	if err != nil {
		return -1, err
	}
	count, err := strconv.Atoi(out)
	if err != nil {
		return -1, err
	}
	return count, err
}

// GetLastRelativeTag - use git describe to retrieve the last relative tag
func (g *gitRepoCLI) GetLastRelativeTag(rev string) (git.Tag, error) {
	cmd := gitCmd(g).WithArgs("describe", "--tags", "--abbrev=0", "--match", "*[0-9]*.[0-9]*.[0-9]*", "--first-parent", rev)
	out, err := cmd.Run()
	if err != nil {
		return git.Tag{}, err
	}
	return git.Tag{Name: strings.TrimSpace(out)}, nil
}

// GetCurrentBranch - use git symbolic-ref to retrieve the current branch name
func (g *gitRepoCLI) GetCurrentBranch() (string, error) {
	branch, err := gitCmd(g).
		WithArgs("symbolic-ref", "--short", "HEAD").
		Run()

	// Then it is probably because we are in detached mode in CI server.
	if err != nil {
		// Most of the time during CI build, the build occurred in a detached HEAD state.
		// And so we can retrieve the current branch name from environment variable.
		branchFromEnv := getCurrentBranchFromEnv()
		if branchFromEnv == "" {
			return "", fmt.Errorf("Unable to retrieve branch name from `git symbolic-ref HEAD` nor %s environment variable", gitRepoBranchEnv)
		}
		return branchFromEnv, nil
	}

	return branch, nil
}

func getCurrentBranchFromEnv() string {
	// We will use CI GIT_BRANCH environment variable.
	// This need to be mapped with real environment variable from your CI server.
	// TODO: eventually add support for most CI environment variable out of the box.
	log.Trace("GitRepo: retrieve branch name from %s env variable", gitRepoBranchEnv)
	return strings.TrimSpace(os.Getenv(gitRepoBranchEnv))
}

func gitCmd(g *gitRepoCLI) *command.Command {
	return command.New("git").InDir(g.dir)
}

func parseRev(from string, to string) string {
	if to == "" {
		to = "HEAD"
	}
	if from == "" {
		return to
	}
	return fmt.Sprintf("%s..%s", from, to)
}
