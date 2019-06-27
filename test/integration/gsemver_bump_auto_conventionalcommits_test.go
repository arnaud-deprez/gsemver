package integration

import (
	"os"
	"reflect"
	"regexp"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/arnaud-deprez/gsemver/internal/git"
	"github.com/arnaud-deprez/gsemver/pkg/version"
)

const (
	README = "README.md"
)

var (
	bumper *version.BumpStrategyOptions
)

func beforeAll(t *testing.T) {
	assert.NoError(t, os.RemoveAll(GIT_REPO_PATH))
	os.MkdirAll(GIT_REPO_PATH, 0755)
	execInGitRepo(t, "git init")
	execInGitRepo(t, "git status")
	gitRepo := git.NewVersionGitRepo(GIT_REPO_PATH)
	bumper = version.NewConventionalCommitBumpStrategyOptions(gitRepo)
}

func afterAll(t *testing.T) {
	out := execInGitRepo(t, "git log --oneline --decorate --graph --all")
	appendToFile(t, "git.log", out)
}

func beforeEach(t *testing.T) {}

func afterEach(t *testing.T) {
	v, err := bumper.Bump()
	assert.NoError(t, err)
	if v.String() != "0.0.0" {
		execInGitRepo(t, "git checkout master")
	}
}

func TestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	beforeAll(t)

	tests := []func(t *testing.T){
		testFirstVersionWithoutCommit,
		testWithFirstFeatureCommit,
		testWithFirstFixCommit,
		testCreateFeaturePullRequest,
		testCreateFeature2PullRequest,
		testCreateFixPullRequestOnMaster,
		testCreateReleaseBranchWithFix,
		testMergeDirectlyReleaseBranchShouldHaveSameVersionOnMaster,
		testCreateFeature3PullRequest,
		testMergeReleaseBranch,
		testCreateFixPullRequestInReleaseBranch,
		testMerge2ReleaseBranch,
	}

	for _, tf := range tests {
		t.Run(runtime.FuncForPC(reflect.ValueOf(tf).Pointer()).Name(), func(t *testing.T) {
			beforeEach(t)
			tf(t)
			afterEach(t)
		})
	}

	afterAll(t)
}

func testFirstVersionWithoutCommit(t *testing.T) {
	assert := assert.New(t)
	v, err := bumper.Bump()
	assert.NoError(err)
	assert.Equal("0.0.0", v.String())
}

func testWithFirstFeatureCommit(t *testing.T) {
	assert := assert.New(t)
	appendToFile(t, README, "First feature")
	commit(t, "feat: add README.md")

	v, err := bumper.Bump()
	assert.NoError(err)
	assert.Equal("0.1.0", v.String())

	createTag(t, v.String())
	v2, err := bumper.Bump()
	assert.NoError(err)
	assert.Equal(v, v2)
}

func testWithFirstFixCommit(t *testing.T) {
	assert := assert.New(t)
	appendToFile(t, README, "First fix")
	commit(t, "fix(doc): fix documentation")

	v, err := bumper.Bump()
	assert.NoError(err)
	assert.Equal("0.1.1", v.String())

	createTag(t, v.String())
	v2, err := bumper.Bump()
	assert.NoError(err)
	assert.Equal(v, v2)
}

func testCreateFeaturePullRequest(t *testing.T) {
	assert := assert.New(t)
	branch := "feature/awesome-1"
	execInGitRepo(t, "git checkout -b "+branch)
	appendToFile(t, README, "Awesome feature with breaking change")
	commit(t, `feat: my awesome change
	
BREAKING CHANGE: this is a breaking change but should not bump major as it is a development release`)

	v, err := bumper.Bump()
	assert.NoError(err)
	assert.Regexp(regexp.MustCompile(`0.1.1\+1\..*`), v.String())

	mergePullRequest(t, branch, "master")

	v, err = bumper.Bump()
	assert.NoError(err)
	assert.Equal("0.2.0", v.String())
	createTag(t, v.String())

	time.Sleep(1 * time.Second)
	createTag(t, "1.0.0")

	v, err = bumper.Bump()
	assert.NoError(err)
	assert.Equal("1.0.0", v.String())
}

func testCreateFeature2PullRequest(t *testing.T) {
	assert := assert.New(t)
	branch := "feature/awesome-2"
	execInGitRepo(t, "git checkout -b "+branch)
	appendToFile(t, README, "Awesome 2nd feature")
	commit(t, `feat: my awesome 2nd change`)

	v, err := bumper.Bump()
	assert.NoError(err)
	assert.Regexp(regexp.MustCompile(`1.0.0\+1\..*`), v.String())

	mergePullRequest(t, branch, "master")
	v, err = bumper.Bump()
	assert.NoError(err)
	assert.Equal("1.1.0", v.String())
	createTag(t, v.String())
}

func testCreateFixPullRequestOnMaster(t *testing.T) {
	assert := assert.New(t)
	branch := "bug/fix-1"
	execInGitRepo(t, "git checkout -b "+branch)
	appendToFile(t, README, "Bug fix on master")
	commit(t, `fix: my bug fix on master`)

	v, err := bumper.Bump()
	assert.NoError(err)
	assert.Regexp(regexp.MustCompile(`1.1.0\+1\..*`), v.String())

	mergePullRequest(t, branch, "master")
	v, err = bumper.Bump()
	assert.NoError(err)
	assert.Equal("1.1.1", v.String())
	createTag(t, v.String())
}

func testCreateReleaseBranchWithFix(t *testing.T) {
	assert := assert.New(t)
	releaseBranch := "release/1.1.x"
	execInGitRepo(t, "git checkout -b "+releaseBranch)

	v, err := bumper.Bump()
	assert.NoError(err)
	assert.Equal("1.1.1", v.String())

	branch := "fix/fix-2"
	execInGitRepo(t, "git checkout -b "+branch)
	appendToFile(t, "README-1.1.x.md", "Bug fix 2 on "+releaseBranch)
	commit(t, `fix: my bug fix 2 on `+releaseBranch)

	v, err = bumper.Bump()
	assert.NoError(err)
	assert.Regexp(regexp.MustCompile(`1.1.1\+1\..*`), v.String())

	mergePullRequest(t, branch, releaseBranch)

	v, err = bumper.Bump()
	assert.NoError(err)
	assert.Equal("1.1.2", v.String())
	createTag(t, v.String())
}

func testMergeDirectlyReleaseBranchShouldHaveSameVersionOnMaster(t *testing.T) {
	assert := assert.New(t)

	// to merge into release branch, we should first perform the merge in a working branch
	branch := "feature/merge-direct-release-1.1.x"
	execInGitRepo(t, "git checkout -b "+branch)
	merge(t, "release/1.1.x", branch)
	mergePullRequest(t, branch, "master")

	v, err := bumper.Bump()
	assert.NoError(err)
	// We should have the same version as semantically nothing is different from release and master branch
	assert.Equal(`1.1.2`, v.String())

	// revert this change
	execInGitRepo(t, "git reset --hard v1.1.1")
}

func testCreateFeature3PullRequest(t *testing.T) {
	assert := assert.New(t)
	branch := "feature/awesome-3"
	execInGitRepo(t, "git checkout -b "+branch)
	appendToFile(t, README, "Awesome 3rd feature")
	commit(t, `feat: my awesome 3rd change`)

	v, err := bumper.Bump()
	assert.NoError(err)
	assert.Regexp(regexp.MustCompile(`1.1.1\+1\..*`), v.String())

	mergePullRequest(t, branch, "master")
	v, err = bumper.Bump()
	assert.NoError(err)
	assert.Equal("1.2.0", v.String())
	createTag(t, v.String())
}

func testMergeReleaseBranch(t *testing.T) {
	assert := assert.New(t)

	// to merge into release branch, we should first perform the merge in a working branch
	branch := "feature/merge-release-1.1.x"
	execInGitRepo(t, "git checkout -b "+branch)
	merge(t, "release/1.1.x", branch)
	mergePullRequest(t, branch, "master")

	v, err := bumper.Bump()
	assert.NoError(err)
	assert.Equal(`1.2.1`, v.String())
	createTag(t, v.String())
}

func testCreateFixPullRequestInReleaseBranch(t *testing.T) {
	assert := assert.New(t)
	releaseBranch := "release/1.1.x"
	execInGitRepo(t, "git checkout "+releaseBranch)

	branch := "fix/fix-3"
	execInGitRepo(t, "git checkout -b "+branch)
	appendToFile(t, "README-1.1.x.md", "Bug fix 3 on "+releaseBranch)
	commit(t, `fix: my bug fix 3 on `+releaseBranch)

	v, err := bumper.Bump()
	assert.NoError(err)
	assert.Regexp(regexp.MustCompile(`1.1.2\+1\..*`), v.String())

	mergePullRequest(t, branch, releaseBranch)

	v, err = bumper.Bump()
	assert.NoError(err)
	assert.Equal("1.1.3", v.String())
	createTag(t, v.String())
}

func testMerge2ReleaseBranch(t *testing.T) {
	assert := assert.New(t)
	// to merge into release branch, we should first perform the merge in a working branch
	branch := "feature/merge2-release-1.1.x"
	execInGitRepo(t, "git checkout -b "+branch)
	merge(t, "release/1.1.x", branch)
	mergePullRequest(t, branch, "master")

	v, err := bumper.Bump()
	assert.NoError(err)
	assert.Equal(`1.2.2`, v.String())
	createTag(t, v.String())
}
