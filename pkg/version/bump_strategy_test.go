package version

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/arnaud-deprez/gsemver/pkg/git"
	mock_version "github.com/arnaud-deprez/gsemver/pkg/version/mock"
)

func TestBumpVersionStrategyWithoutTag(t *testing.T) {
	assert := assert.New(t)

	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testData := []struct {
		strategy            BumpStrategyType
		branch              string
		preRelease          string
		preReleaseOverwrite bool
		buildMetadata       string
		expected            string
	}{
		{MAJOR, "dummy", "", false, "", "1.0.0"},
		{MINOR, "dummy", "", false, "", "0.1.0"},
		{PATCH, "dummy", "", false, "", "0.0.1"},
		{AUTO, "master", "", false, "", "0.1.0"},
		{AUTO, "feature/test", "", false, "{{ .Commits | len }}.{{ (.Commits | first).Hash.Short }}", "0.0.0+1.1234567"},
		{MAJOR, "dummy", "alpha", false, "", "1.0.0-alpha.0"},
		{MINOR, "dummy", "SNAPSHOT", true, "", "0.1.0-SNAPSHOT"},
		{0, "dummy", "", false, "build.8", "0.0.0+build.8"},
	}

	for _, tc := range testData {
		gitRepo := mock_version.NewMockGitRepo(ctrl)
		gitRepo.EXPECT().FetchTags().Times(1).Return(nil)
		gitRepo.EXPECT().GetLastRelativeTag("HEAD").Times(1).Return(git.Tag{}, nil)
		// no commit so it should return the same version
		gitRepo.EXPECT().GetCommits("", "HEAD").Times(1).Return([]git.Commit{
			{
				Author:    git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
				Committer: git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
				Hash:      git.Hash("1234567890"),
				Message:   `feat: init import`,
			},
		}, nil)
		gitRepo.EXPECT().GetCurrentBranch().Times(1).Return(tc.branch, nil)

		strategy := NewConventionalCommitBumpStrategy(gitRepo)
		strategy.Strategy = tc.strategy
		strategy.BumpDefaultStrategy.PreReleaseTemplate = NewTemplate(tc.preRelease)
		strategy.BumpDefaultStrategy.PreReleaseOverwrite = tc.preReleaseOverwrite
		strategy.BumpDefaultStrategy.BuildMetadataTemplate = NewTemplate(tc.buildMetadata)
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
		strategy            BumpStrategyType
		branch              string
		preRelease          string
		preReleaseOverwrite bool
		buildMetadata       string
		expected            string
	}{
		{"v1.1.0-alpha.0", MAJOR, "dummy", "", false, "", "2.0.0"},
		{"v1.1.0", PATCH, "dummy", "", false, "", "1.1.1"},
		{"v1.2.0", MINOR, "dummy", "", false, "", "1.3.0"},
		{"v1.2.0", MINOR, "dummy", "alpha", false, "", "1.3.0-alpha.0"},
		{"1.2.0", MAJOR, "dummy", "SNAPSHOT", true, "", "2.0.0-SNAPSHOT"},
		{"v1.2.0", MAJOR, "dummy", "SNAPSHOT", true, "", "2.0.0-SNAPSHOT"},
		// in AUTO the version should not change
		{"1.2.0", AUTO, "master", "", false, "", "1.2.0"},
		{"v1.2.0", AUTO, "master", "", false, "", "1.2.0"},
		{"v1.2.0", AUTO, "feature/test", "", false, "", "1.2.0"},
		{"v1.2.0", AUTO, "master", "alpha", false, "", "1.2.0"},
		{"v1.2.0", AUTO, "feature/test", "alpha", false, "", "1.2.0"},
		{"v1.2.0", AUTO, "master", "SNAPSHOT", true, "", "1.2.0"},
		{"v1.2.0", AUTO, "feature/test", "SNAPSHOT", true, "", "1.2.0"},
		{"v1.2.0", AUTO, "master", "", false, "build.1", "1.2.0"},
		{"v1.2.0", AUTO, "feature/test", "", false, "build.1", "1.2.0"},
	}

	for _, tc := range testData {
		gitRepo := mock_version.NewMockGitRepo(ctrl)
		gitRepo.EXPECT().FetchTags().Times(1).Return(nil)
		gitRepo.EXPECT().GetLastRelativeTag("HEAD").Times(1).Return(git.Tag{Name: tc.from}, nil)
		// no commit so it should return the same version
		gitRepo.EXPECT().GetCommits(tc.from, "HEAD").Times(1).Return([]git.Commit{}, nil)
		gitRepo.EXPECT().GetCurrentBranch().Times(1).Return(tc.branch, nil)

		strategy := NewConventionalCommitBumpStrategy(gitRepo)
		strategy.Strategy = tc.strategy
		strategy.BumpDefaultStrategy.PreReleaseTemplate = NewTemplate(tc.preRelease)
		strategy.BumpDefaultStrategy.PreReleaseOverwrite = tc.preReleaseOverwrite
		strategy.BumpDefaultStrategy.BuildMetadataTemplate = NewTemplate(tc.buildMetadata)
		version, err := strategy.Bump()

		assert.Nil(err)
		assert.Equal(tc.expected, version.String())
	}
}

func TestBumpVersionStrategyMajor(t *testing.T) {
	assert := assert.New(t)

	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gitRepo := mock_version.NewMockGitRepo(ctrl)
	from := "v0.1.0"
	gitRepo.EXPECT().FetchTags().Times(1).Return(nil)
	gitRepo.EXPECT().GetLastRelativeTag("HEAD").Times(1).Return(git.Tag{Name: from}, nil)
	gitRepo.EXPECT().GetCommits(from, "HEAD").Times(1).Return([]git.Commit{
		{
			Author:    git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
			Committer: git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
			Hash:      git.Hash("1234567890"),
			Message:   `This is not relevant`,
		},
	}, nil)
	gitRepo.EXPECT().GetCurrentBranch().Times(1).Return("dummy", nil)

	strategy := &BumpStrategy{Strategy: MAJOR, gitRepo: gitRepo}
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
	gitRepo.EXPECT().FetchTags().Times(1).Return(nil)
	gitRepo.EXPECT().GetLastRelativeTag("HEAD").Times(1).Return(git.Tag{Name: from}, nil)
	gitRepo.EXPECT().GetCommits(from, "HEAD").Times(1).Return([]git.Commit{
		{
			Author:    git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
			Committer: git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
			Hash:      git.Hash("1234567890"),
			Message:   `This is not relevant`,
		},
	}, nil)
	gitRepo.EXPECT().GetCurrentBranch().Times(1).Return("dummy", nil)

	strategy := &BumpStrategy{Strategy: MINOR, gitRepo: gitRepo}
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
	gitRepo.EXPECT().FetchTags().Times(1).Return(nil)
	gitRepo.EXPECT().GetLastRelativeTag("HEAD").Times(1).Return(git.Tag{Name: from}, nil)
	gitRepo.EXPECT().GetCommits(from, "HEAD").Times(1).Return([]git.Commit{
		{
			Author:    git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
			Committer: git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
			Hash:      git.Hash("1234567890"),
			Message:   `This is not relevant`,
		},
	}, nil)
	gitRepo.EXPECT().GetCurrentBranch().Times(1).Return("dummy", nil)

	strategy := &BumpStrategy{Strategy: PATCH, gitRepo: gitRepo}
	version, err := strategy.Bump()

	assert.Nil(err)
	assert.Equal("0.1.1", version.String())
}

func TestBumpVersionStrategyAutoBreakingChangeOnInitialDevelopmentRelease(t *testing.T) {
	assert := assert.New(t)

	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testData := []struct {
		from     string
		branch   string
		expected string
	}{
		{"v0.1.0", "master", "0.2.0"},
		{"v0.1.0", "feature/test", "0.1.0+1.1234567"},
	}

	for _, tc := range testData {
		gitRepo := mock_version.NewMockGitRepo(ctrl)
		gitRepo.EXPECT().FetchTags().Times(1).Return(nil)
		gitRepo.EXPECT().GetLastRelativeTag("HEAD").Times(1).Return(git.Tag{Name: tc.from}, nil)
		gitRepo.EXPECT().GetCommits(tc.from, "HEAD").Times(1).Return([]git.Commit{
			{
				Author:    git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
				Committer: git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
				Hash:      git.Hash("1234567890"),
				Message: `feat(version): add auto bump strategies 
		
BREAKING CHANGE: replace next option by bump for more convenience
			`,
			},
		}, nil)
		gitRepo.EXPECT().GetCurrentBranch().Times(1).Return(tc.branch, nil)

		strategy := NewConventionalCommitBumpStrategy(gitRepo)
		version, err := strategy.Bump()

		assert.Nil(err)
		assert.Equal(tc.expected, version.String())
	}
}

func TestBumpVersionStrategyAutoBreakingChange(t *testing.T) {
	assert := assert.New(t)

	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testData := []struct {
		from     string
		branch   string
		expected string
	}{
		{"v1.1.0", "master", "2.0.0"},
		{"v1.1.0", "feature/test", "1.1.0+2.1234567"},
	}

	for _, tc := range testData {
		gitRepo := mock_version.NewMockGitRepo(ctrl)
		gitRepo.EXPECT().FetchTags().Times(1).Return(nil)
		gitRepo.EXPECT().GetLastRelativeTag("HEAD").Times(1).Return(git.Tag{Name: tc.from}, nil)
		gitRepo.EXPECT().GetCommits(tc.from, "HEAD").Times(1).Return([]git.Commit{
			{
				Author:    git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
				Committer: git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
				Hash:      git.Hash("1234567890"),
				Message: `feat(version): add auto bump strategies 
		
BREAKING CHANGE: replace next option by bump for more convenience
			`,
			},
			{
				Author:    git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
				Committer: git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
				Hash:      git.Hash("1234567890"),
				Message:   `feat(version): add pre-release option`,
			},
		}, nil)
		gitRepo.EXPECT().GetCurrentBranch().Times(1).Return(tc.branch, nil)

		strategy := NewConventionalCommitBumpStrategy(gitRepo)
		version, err := strategy.Bump()

		assert.Nil(err)
		assert.Equal(tc.expected, version.String())
	}
}

func TestBumpVersionStrategyAutoWithNewFeature(t *testing.T) {
	assert := assert.New(t)

	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testData := []struct {
		from     string
		branch   string
		expected string
	}{
		{"v1.1.0", "master", "1.2.0"},
		{"v1.1.0", "feature/test", "1.1.0+1.1234567"},
	}

	for _, tc := range testData {
		gitRepo := mock_version.NewMockGitRepo(ctrl)
		gitRepo.EXPECT().FetchTags().Times(1).Return(nil)
		gitRepo.EXPECT().GetLastRelativeTag("HEAD").Times(1).Return(git.Tag{Name: tc.from}, nil)
		gitRepo.EXPECT().GetCommits(tc.from, "HEAD").Times(1).Return([]git.Commit{
			{
				Author:    git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
				Committer: git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
				Hash:      git.Hash("1234567890"),
				Message:   `feat(version): add pre-release option`,
			},
		}, nil)
		gitRepo.EXPECT().GetCurrentBranch().Times(1).Return(tc.branch, nil)

		strategy := NewConventionalCommitBumpStrategy(gitRepo)
		version, err := strategy.Bump()

		assert.Nil(err)
		assert.Equal(tc.expected, version.String())
	}
}

func TestBumpVersionStrategyAutoWithPatch(t *testing.T) {
	assert := assert.New(t)

	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testData := []struct {
		from     string
		branch   string
		expected string
	}{
		{"v1.1.0", "master", "1.1.1"},
		{"v1.1.0", "feature/test", "1.1.0+1.1234567"},
		{"v1.1.0", "release/1.1.x", "1.1.1"},
	}

	for _, tc := range testData {
		gitRepo := mock_version.NewMockGitRepo(ctrl)
		gitRepo.EXPECT().FetchTags().Times(1).Return(nil)
		gitRepo.EXPECT().GetLastRelativeTag("HEAD").Times(1).Return(git.Tag{Name: tc.from}, nil)
		gitRepo.EXPECT().GetCommits(tc.from, "HEAD").Times(1).Return([]git.Commit{
			{
				Author:    git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
				Committer: git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
				Hash:      git.Hash("1234567890"),
				Message:   `fix: typo error`,
			},
		}, nil)
		gitRepo.EXPECT().GetCurrentBranch().Times(1).Return(tc.branch, nil)

		strategy := NewConventionalCommitBumpStrategy(gitRepo)
		version, err := strategy.Bump()

		assert.Nil(err)
		assert.Equal(tc.expected, version.String())
	}
}

func TestBumpVersionStrategyAutoWithPreReleaseStrategyAndNewFeature(t *testing.T) {
	assert := assert.New(t)

	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testData := []struct {
		from       string
		branch     string
		preRelease string
		expected   string
	}{
		{"v1.1.0", "master", "alpha", "1.2.0"},
		{"v1.1.0", "milestone-1.2", "alpha", "1.2.0-alpha.0"},
		{"v1.2.0-alpha.0", "milestone-1.2", "alpha", "1.2.0-alpha.1"},
		{"v1.1.0", "feature/test", "alpha", "1.1.0+1.1234567"},
		{"v1.1.0-alpha.0", "feature/test", "alpha", "1.1.0-alpha.0+1.1234567"},
	}

	for _, tc := range testData {
		gitRepo := mock_version.NewMockGitRepo(ctrl)
		gitRepo.EXPECT().FetchTags().Times(1).Return(nil)
		gitRepo.EXPECT().GetLastRelativeTag("HEAD").Times(1).Return(git.Tag{Name: tc.from}, nil)
		gitRepo.EXPECT().GetCommits(tc.from, "HEAD").Times(1).Return([]git.Commit{
			{
				Author:    git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
				Committer: git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
				Hash:      git.Hash("1234567890"),
				Message:   `feat(version): add pre-release option`,
			},
		}, nil)
		gitRepo.EXPECT().GetCurrentBranch().Times(1).Return(tc.branch, nil)

		strategy := NewConventionalCommitBumpStrategy(gitRepo).AddBumpBranchesStrategy("milestone-1.2", tc.preRelease, false, "")
		version, err := strategy.Bump()

		assert.Nil(err)
		assert.Equal(tc.expected, version.String())
	}
}

func TestBumpVersionStrategyAutoWithPreReleaseMavenLike(t *testing.T) {
	assert := assert.New(t)

	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gitRepo := mock_version.NewMockGitRepo(ctrl)
	from := "v1.0.0"
	gitRepo.EXPECT().FetchTags().Times(1).Return(nil)
	gitRepo.EXPECT().GetLastRelativeTag("HEAD").Times(1).Return(git.Tag{Name: from}, nil)
	gitRepo.EXPECT().GetCommits(from, "HEAD").Times(1).Return([]git.Commit{
		{
			Author:    git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
			Committer: git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
			Hash:      git.Hash("1234567890"),
			Message:   `feat(version): add pre-release option`,
		},
	}, nil)
	gitRepo.EXPECT().GetCurrentBranch().Times(1).Return("feature/xyz", nil)

	strategy := NewConventionalCommitBumpStrategy(gitRepo)
	strategy.BumpDefaultStrategy.PreReleaseTemplate = NewTemplate("SNAPSHOT")
	strategy.BumpDefaultStrategy.PreReleaseOverwrite = true
	strategy.BumpDefaultStrategy.BuildMetadataTemplate = nil
	version, err := strategy.Bump()

	assert.Nil(err)
	assert.Equal("1.1.0-SNAPSHOT", version.String())
}
