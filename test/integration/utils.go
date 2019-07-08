package integration

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/arnaud-deprez/gsemver/internal/command"
)

const gitRepoPath = "./build/git-tmp"

func execInDir(t *testing.T, dir, cmd string) string {
	out, err := command.New(cmd).InDir(dir).Run()
	assert.NoError(t, err)
	return out
}

func execInGitRepo(t *testing.T, cmd string) string {
	return execInDir(t, gitRepoPath, cmd)
}

func createTag(t *testing.T, tag string) {
	execInGitRepo(t, fmt.Sprintf(`git tag -fa v%s -m "Release %s"`, tag, tag))
}

func commit(t *testing.T, msg string) {
	execInGitRepo(t, "git add --all")
	execInGitRepo(t, fmt.Sprintf(`git commit -am "%s"`, msg))
}

func merge(t *testing.T, from, to string) {
	execInGitRepo(t, "git checkout "+to)
	execInGitRepo(t, fmt.Sprintf(`git merge --no-ff -m "Merge from %s" %s`, from, from))
}

func mergePullRequest(t *testing.T, from, to string) {
	merge(t, from, to)
	execInGitRepo(t, "git branch -d "+from)
}

func appendToFile(t *testing.T, file, content string) {
	f, err := os.OpenFile(gitRepoPath+"/"+file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	assert.NoError(t, err)
	defer f.Close()
	_, err = f.WriteString(content + "\n")
	assert.NoError(t, err)
}
