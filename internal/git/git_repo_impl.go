package git

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/arnaud-deprez/gsemver/internal/command"
	"github.com/arnaud-deprez/gsemver/pkg/git"
	"github.com/arnaud-deprez/gsemver/pkg/version"
)

type gitRepoCLI struct {
	version.GitRepo
	dir          string
	commitParser *commitParser
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
	rev = parseRev("", rev)
	cmd := gitCmd(g).WithArgs("describe", "--abbrev=0", "--match", "v[0-9]*.[0-9]*.[0-9]*", "--first-parent", rev)
	out, err := cmd.Run()
	if err != nil {
		return git.Tag{}, err
	}
	return git.Tag{Name: strings.TrimSpace(out)}, nil
}

// git-symbolic-ref - Read symbolic refs
func (g *gitRepoCLI) GetSymbolicRef(name string, short bool) (string, error) {
	cmd := gitCmd(g).WithArgs("symbolic-ref")

	if short {
		cmd.WithArg("--short")
	}
	cmd.WithArg(parseRev("", name))

	return cmd.Run()
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
