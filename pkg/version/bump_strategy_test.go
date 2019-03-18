package version

import (
	"testing"

	"github.com/arnaud-deprez/gsemver/pkg/git"
	mock_version "github.com/arnaud-deprez/gsemver/pkg/version/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBumpVersionStrategyWithoutTag(t *testing.T) {
	assert := assert.New(t)

	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testData := []struct {
		strategy            BumpStrategy
		preRelease          string
		preReleaseOverwrite bool
		buildMetadata       string
		expected            string
	}{
		{MAJOR, "", false, "", "1.0.0"},
		{MINOR, "", false, "", "0.1.0"},
		{PATCH, "", false, "", "0.0.1"},
		{AUTO, "", false, "", "0.1.0"},
		{MAJOR, "alpha", false, "", "1.0.0-alpha.0"},
		{MINOR, "SNAPSHOT", true, "", "0.1.0-SNAPSHOT"},
		//TODO: normally a buildMetadata version should not be bumped, whatever value you give to BumpStrategy.
		{0, "", false, "build.8", "0.0.0+build.8"},
	}

	for _, tc := range testData {
		gitRepo := mock_version.NewMockGitRepo(ctrl)
		gitRepo.EXPECT().GetLastRelativeTag("HEAD").Times(1).Return(git.Tag{}, nil)
		// no commit so it should return the same version
		gitRepo.EXPECT().GetCommits("", "HEAD").Times(1).Return([]git.Commit{
			{
				Author:    git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
				Committer: git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
				Hash:      git.NewHash("1234567890"),
				Message:   `feat: init import`,
			},
		}, nil)

		strategy := NewConventionalCommitBumpStrategyOptions(gitRepo)
		strategy.Strategy = tc.strategy
		strategy.PreRelease = tc.preRelease
		strategy.PreReleaseOverwrite = tc.preReleaseOverwrite
		strategy.BuildMetadata = tc.buildMetadata
		version, err := strategy.Bump()
		assert.Nil(err)
		assert.Equal(tc.expected, version.String())
	}
}

func TestBumpVersionStrategyNoDeltaCommit(t *testing.T) {
	assert := assert.New(t)

	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testData := []struct {
		from                string
		strategy            BumpStrategy
		preRelease          string
		preReleaseOverwrite bool
		buildMetadata       string
		expected            string
	}{
		{"v1.1.0-alpha.0", MAJOR, "", false, "", "2.0.0"},
		{"1.1.0", PATCH, "", false, "", "1.1.1"},
		{"1.2.0", MINOR, "", false, "", "1.3.0"},
		{"1.2.0", MINOR, "alpha", false, "", "1.3.0-alpha.0"},
		{"1.2.0", MAJOR, "SNAPSHOT", true, "", "2.0.0-SNAPSHOT"},
		// in AUTO the version should not change
		{"1.2.0", AUTO, "", false, "", "1.2.0"},
		{"1.2.0", AUTO, "alpha", false, "", "1.2.0"},
		{"1.2.0", AUTO, "SNAPSHOT", true, "", "1.2.0"},
		{"1.2.0", AUTO, "", false, "build.1", "1.2.0"},
	}

	for _, tc := range testData {
		gitRepo := mock_version.NewMockGitRepo(ctrl)
		gitRepo.EXPECT().GetLastRelativeTag("HEAD").Times(1).Return(git.Tag{Name: tc.from}, nil)
		// no commit so it should return the same version
		gitRepo.EXPECT().GetCommits(tc.from, "HEAD").Times(1).Return([]git.Commit{}, nil)

		strategy := NewConventionalCommitBumpStrategyOptions(gitRepo)
		strategy.Strategy = tc.strategy
		strategy.PreRelease = tc.preRelease
		strategy.PreReleaseOverwrite = tc.preReleaseOverwrite
		strategy.BuildMetadata = tc.buildMetadata
		version, err := strategy.Bump()

		assert.Nil(err)
		assert.Equal(tc.expected, version.String())
	}
}

func TestBumpVersionStrategy(t *testing.T) {
	assert := assert.New(t)

	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gitRepo := mock_version.NewMockGitRepo(ctrl)
	from := "v0.1.0"
	gitRepo.EXPECT().GetLastRelativeTag("HEAD").Times(1).Return(git.Tag{Name: from}, nil)
	gitRepo.EXPECT().GetCommits(from, "HEAD").Times(1).Return([]git.Commit{
		{
			Author:    git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
			Committer: git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
			Hash:      git.NewHash("1234567890"),
			Message:   `This is not relevant`,
		},
	}, nil)

	strategy := &BumpStrategyOptions{Strategy: MAJOR, gitRepo: gitRepo}
	version, err := strategy.Bump()

	assert.Nil(err)
	assert.Equal("1.0.0", version.String())
}

func TestBumpVersionStrategyMinor(t *testing.T) {
	assert := assert.New(t)

	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gitRepo := mock_version.NewMockGitRepo(ctrl)
	from := "v0.1.0"
	gitRepo.EXPECT().GetLastRelativeTag("HEAD").Times(1).Return(git.Tag{Name: from}, nil)
	gitRepo.EXPECT().GetCommits(from, "HEAD").Times(1).Return([]git.Commit{
		{
			Author:    git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
			Committer: git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
			Hash:      git.NewHash("1234567890"),
			Message:   `This is not relevant`,
		},
	}, nil)

	strategy := &BumpStrategyOptions{Strategy: MINOR, gitRepo: gitRepo}
	version, err := strategy.Bump()

	assert.Nil(err)
	assert.Equal("0.2.0", version.String())
}

func TestBumpVersionStrategyPatch(t *testing.T) {
	assert := assert.New(t)

	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gitRepo := mock_version.NewMockGitRepo(ctrl)
	from := "v0.1.0"
	gitRepo.EXPECT().GetLastRelativeTag("HEAD").Times(1).Return(git.Tag{Name: from}, nil)
	gitRepo.EXPECT().GetCommits(from, "HEAD").Times(1).Return([]git.Commit{
		{
			Author:    git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
			Committer: git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
			Hash:      git.NewHash("1234567890"),
			Message:   `This is not relevant`,
		},
	}, nil)

	strategy := &BumpStrategyOptions{Strategy: PATCH, gitRepo: gitRepo}
	version, err := strategy.Bump()

	assert.Nil(err)
	assert.Equal("0.1.1", version.String())
}

func TestBumpVersionStrategyAutoShouldBreakingMinor(t *testing.T) {
	assert := assert.New(t)

	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gitRepo := mock_version.NewMockGitRepo(ctrl)
	from := "v0.1.0"
	gitRepo.EXPECT().GetLastRelativeTag("HEAD").Times(1).Return(git.Tag{Name: from}, nil)
	gitRepo.EXPECT().GetCommits(from, "HEAD").Times(1).Return([]git.Commit{
		{
			Author:    git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
			Committer: git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
			Hash:      git.NewHash("1234567890"),
			Message: `feat(version): add auto bump strategies 
		
BREAKING CHANGE: replace next option by bump for more convenience
			`,
		},
	}, nil)

	strategy := NewConventionalCommitBumpStrategyOptions(gitRepo)
	version, err := strategy.Bump()

	assert.Nil(err)
	assert.Equal("0.2.0", version.String())
}

func TestBumpVersionStrategyAutoShouldBreakingMajor(t *testing.T) {
	assert := assert.New(t)

	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gitRepo := mock_version.NewMockGitRepo(ctrl)
	from := "v1.1.0"
	gitRepo.EXPECT().GetLastRelativeTag("HEAD").Times(1).Return(git.Tag{Name: from}, nil)
	gitRepo.EXPECT().GetCommits(from, "HEAD").Times(1).Return([]git.Commit{
		{
			Author:    git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
			Committer: git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
			Hash:      git.NewHash("1234567890"),
			Message: `feat(version): add auto bump strategies 
		
BREAKING CHANGE: replace next option by bump for more convenience
			`,
		},
	}, nil)

	strategy := NewConventionalCommitBumpStrategyOptions(gitRepo)
	version, err := strategy.Bump()

	assert.Nil(err)
	assert.Equal("2.0.0", version.String())
}

func TestBumpVersionStrategyAutoShouldMinor(t *testing.T) {
	assert := assert.New(t)

	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gitRepo := mock_version.NewMockGitRepo(ctrl)
	from := "v0.1.0"
	gitRepo.EXPECT().GetLastRelativeTag("HEAD").Times(1).Return(git.Tag{Name: from}, nil)
	gitRepo.EXPECT().GetCommits(from, "HEAD").Times(1).Return([]git.Commit{
		{
			Author:    git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
			Committer: git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
			Hash:      git.NewHash("1234567890"),
			Message:   `feat(version): add pre-release option`,
		},
	}, nil)

	strategy := NewConventionalCommitBumpStrategyOptions(gitRepo)
	version, err := strategy.Bump()

	assert.Nil(err)
	assert.Equal("0.2.0", version.String())
}

func TestBumpVersionStrategyAutoShouldPatch(t *testing.T) {
	assert := assert.New(t)

	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gitRepo := mock_version.NewMockGitRepo(ctrl)
	from := "v0.1.0"
	gitRepo.EXPECT().GetLastRelativeTag("HEAD").Times(1).Return(git.Tag{Name: from}, nil)
	gitRepo.EXPECT().GetCommits(from, "HEAD").Times(1).Return([]git.Commit{
		{
			Author:    git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
			Committer: git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
			Hash:      git.NewHash("1234567890"),
			Message:   `doc: add FAQ`,
		},
	}, nil)

	strategy := NewConventionalCommitBumpStrategyOptions(gitRepo)
	version, err := strategy.Bump()

	assert.Nil(err)
	assert.Equal("0.1.1", version.String())
}

func TestBumpVersionStrategyAutoShouldPreRelease(t *testing.T) {
	assert := assert.New(t)

	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gitRepo := mock_version.NewMockGitRepo(ctrl)
	from := "v1.0.0"
	gitRepo.EXPECT().GetLastRelativeTag("HEAD").Times(1).Return(git.Tag{Name: from}, nil)
	gitRepo.EXPECT().GetCommits(from, "HEAD").Times(1).Return([]git.Commit{
		{
			Author:    git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
			Committer: git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
			Hash:      git.NewHash("1234567890"),
			Message:   `feat(version): add pre-release option`,
		},
	}, nil)

	strategy := NewConventionalCommitBumpStrategyOptions(gitRepo).WithPreRelease("alpha", false)
	version, err := strategy.Bump()

	assert.Nil(err)
	assert.Equal("1.1.0-alpha.0", version.String())
}

func TestBumpVersionStrategyAutoShouldPreReleaseMavenLike(t *testing.T) {
	assert := assert.New(t)

	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gitRepo := mock_version.NewMockGitRepo(ctrl)
	from := "v1.0.0"
	gitRepo.EXPECT().GetLastRelativeTag("HEAD").Times(1).Return(git.Tag{Name: from}, nil)
	gitRepo.EXPECT().GetCommits(from, "HEAD").Times(1).Return([]git.Commit{
		{
			Author:    git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
			Committer: git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
			Hash:      git.NewHash("1234567890"),
			Message:   `feat(version): add pre-release option`,
		},
	}, nil)

	strategy := NewConventionalCommitBumpStrategyOptions(gitRepo).WithPreRelease("SNAPSHOT", true)
	version, err := strategy.Bump()

	assert.Nil(err)
	assert.Equal("1.1.0-SNAPSHOT", version.String())
}
