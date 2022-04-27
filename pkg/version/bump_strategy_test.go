package version

import (
	"fmt"
	"regexp"
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
		strategy              BumpStrategyType
		branch                string
		preRelease            bool
		preReleaseTemplate    string
		preReleaseOverwrite   bool
		buildMetadataTemplate string
		expected              string
	}{
		{MAJOR, "dummy", false, "", false, "", "1.0.0"},
		{MINOR, "dummy", false, "", false, "", "0.1.0"},
		{PATCH, "dummy", false, "", false, "", "0.0.1"},
		{MAJOR, "dummy", true, "", false, "", "1.0.0-0"},
		{MAJOR, "dummy", true, "alpha", false, "", "1.0.0-alpha.0"},
		{MINOR, "dummy", true, "SNAPSHOT", true, "", "0.1.0-SNAPSHOT"},
		{0, "dummy", false, "", false, "build.8", "0.0.0+build.8"},
		{AUTO, "master", false, "", false, "", "0.1.0"},
		{AUTO, "master", true, "", false, "", "0.1.0-0"},
		{AUTO, "master", true, "alpha", false, "", "0.1.0-alpha.0"},
		{AUTO, "master", false, "", false, "build.1", "0.0.0+build.1"},
		{AUTO, "feature/test", false, "", false, "{{ .Commits | len }}.{{ (.Commits | first).Hash.Short }}", "0.0.0+1.1234567"},
	}

	for idx, tc := range testData {
		t.Run(fmt.Sprintf("Case %d", idx), func(t *testing.T) {
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
			strategy.BumpStrategies = []BumpBranchesStrategy{*NewBumpAllBranchesStrategy(tc.strategy, tc.preRelease, tc.preReleaseTemplate, tc.preReleaseOverwrite, tc.buildMetadataTemplate)}
			version, err := strategy.Bump()
			assert.Nil(err)
			assert.Equal(tc.expected, version.String())
		})
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
		preRelease          bool
		preReleaseTemplate  string
		preReleaseOverwrite bool
		buildMetadata       string
		expected            string
	}{
		{"v1.1.0-alpha.0", MAJOR, "dummy", false, "", false, "", "2.0.0"},
		{"v1.1.0", PATCH, "dummy", false, "", false, "", "1.1.1"},
		{"v1.2.0", MINOR, "dummy", false, "", false, "", "1.3.0"},
		{"v1.2.0", MINOR, "dummy", true, "alpha", false, "", "1.3.0-alpha.0"},
		{"1.2.0", MAJOR, "dummy", true, "SNAPSHOT", true, "", "2.0.0-SNAPSHOT"},
		{"v1.2.0", MAJOR, "dummy", true, "SNAPSHOT", true, "", "2.0.0-SNAPSHOT"},
		// in AUTO the version should not change
		{"1.2.0", AUTO, "master", false, "", false, "", "1.2.0"},
		{"v1.2.0", AUTO, "master", false, "", false, "", "1.2.0"},
		{"v1.2.0", AUTO, "feature/test", false, "", false, "", "1.2.0"},
		{"v1.2.0", AUTO, "master", true, "", false, "", "1.2.0"},
		{"v1.2.0", AUTO, "master", true, "alpha", false, "", "1.2.0"},
		{"v1.2.0", AUTO, "feature/test", true, "alpha", false, "", "1.2.0"},
		{"v1.2.0", AUTO, "master", true, "SNAPSHOT", true, "", "1.2.0"},
		{"v1.2.0", AUTO, "feature/test", true, "SNAPSHOT", true, "", "1.2.0"},
		{"v1.2.0", AUTO, "master", false, "", false, "build.1", "1.2.0"},
		{"v1.2.0", AUTO, "feature/test", false, "", false, "build.1", "1.2.0"},
	}

	for idx, tc := range testData {
		t.Run(fmt.Sprintf("Case %d", idx), func(t *testing.T) {
			gitRepo := mock_version.NewMockGitRepo(ctrl)
			gitRepo.EXPECT().FetchTags().Times(1).Return(nil)
			gitRepo.EXPECT().GetLastRelativeTag("HEAD").Times(1).Return(git.Tag{Name: tc.from}, nil)
			// no commit so it should return the same version
			gitRepo.EXPECT().GetCommits(tc.from, "HEAD").Times(1).Return([]git.Commit{}, nil)
			gitRepo.EXPECT().GetCurrentBranch().Times(1).Return(tc.branch, nil)

			strategy := NewConventionalCommitBumpStrategy(gitRepo)
			strategy.BumpStrategies = []BumpBranchesStrategy{*NewBumpAllBranchesStrategy(tc.strategy, tc.preRelease, tc.preReleaseTemplate, tc.preReleaseOverwrite, tc.buildMetadata)}
			version, err := strategy.Bump()

			assert.Nil(err)
			assert.Equal(tc.expected, version.String())
		})
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

	strategy := &BumpStrategy{BumpStrategies: []BumpBranchesStrategy{*NewBumpAllBranchesStrategy(MAJOR, false, "", false, "")}, gitRepo: gitRepo}
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

	strategy := &BumpStrategy{BumpStrategies: []BumpBranchesStrategy{*NewBumpAllBranchesStrategy(MINOR, false, "", false, "")}, gitRepo: gitRepo}
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

	strategy := &BumpStrategy{BumpStrategies: []BumpBranchesStrategy{*NewBumpAllBranchesStrategy(PATCH, false, "", false, "")}, gitRepo: gitRepo}
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

	for idx, tc := range testData {
		t.Run(fmt.Sprintf("Case %d", idx), func(t *testing.T) {
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
		})
	}
}

func TestBumpVersionStrategyAutoBreakingChangeOnInitialDevelopmentReleaseShortForm(t *testing.T) {
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

	for idx, tc := range testData {
		t.Run(fmt.Sprintf("Case %d", idx), func(t *testing.T) {
			gitRepo := mock_version.NewMockGitRepo(ctrl)
			gitRepo.EXPECT().FetchTags().Times(1).Return(nil)
			gitRepo.EXPECT().GetLastRelativeTag("HEAD").Times(1).Return(git.Tag{Name: tc.from}, nil)
			gitRepo.EXPECT().GetCommits(tc.from, "HEAD").Times(1).Return([]git.Commit{
				{
					Author:    git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
					Committer: git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
					Hash:      git.Hash("1234567890"),
					Message:   `feat(version)!: add auto bump strategies`,
				},
			}, nil)
			gitRepo.EXPECT().GetCurrentBranch().Times(1).Return(tc.branch, nil)

			strategy := NewConventionalCommitBumpStrategy(gitRepo)
			version, err := strategy.Bump()

			assert.Nil(err)
			assert.Equal(tc.expected, version.String())
		})
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

	for idx, tc := range testData {
		t.Run(fmt.Sprintf("Case %d", idx), func(t *testing.T) {
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
		})
	}
}

func TestBumpVersionStrategyAutoBreakingChangeShortForm(t *testing.T) {
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

	for idx, tc := range testData {
		t.Run(fmt.Sprintf("Case %d", idx), func(t *testing.T) {
			gitRepo := mock_version.NewMockGitRepo(ctrl)
			gitRepo.EXPECT().FetchTags().Times(1).Return(nil)
			gitRepo.EXPECT().GetLastRelativeTag("HEAD").Times(1).Return(git.Tag{Name: tc.from}, nil)
			gitRepo.EXPECT().GetCommits(tc.from, "HEAD").Times(1).Return([]git.Commit{
				{
					Author:    git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
					Committer: git.Signature{Name: "Arnaud Deprez", Email: "xxx@example.com"},
					Hash:      git.Hash("1234567890"),
					Message:   `fix(version)!: add auto bump strategies`,
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
		})
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

	for idx, tc := range testData {
		t.Run(fmt.Sprintf("Case %d", idx), func(t *testing.T) {
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
		})
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

	for idx, tc := range testData {
		t.Run(fmt.Sprintf("Case %d", idx), func(t *testing.T) {
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
		})
	}
}

func TestBumpVersionStrategyAutoWithPreReleaseStrategyAndNewFeature(t *testing.T) {
	assert := assert.New(t)

	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testData := []struct {
		from               string
		branch             string
		preRelease         bool
		preReleaseTemplate string
		expected           string
	}{
		{"v1.1.0", "master", true, "alpha", "1.2.0"},
		{"v1.1.0", "milestone-1.2", true, "alpha", "1.2.0-alpha.0"},
		{"v1.2.0-alpha.0", "milestone-1.2", true, "alpha", "1.2.0-alpha.1"},
		{"v1.1.0", "feature/test", true, "alpha", "1.1.0+1.1234567"},
		{"v1.1.0-alpha.0", "feature/test", true, "alpha", "1.1.0-alpha.0+1.1234567"},
	}

	for idx, tc := range testData {
		t.Run(fmt.Sprintf("Case %d", idx), func(t *testing.T) {
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

			strategy := &BumpStrategy{
				gitRepo: gitRepo,
				BumpStrategies: []BumpBranchesStrategy{
					*NewDefaultBumpBranchesStrategy(DefaultReleaseBranchesPattern),
					*NewBumpBranchesStrategy(AUTO, "milestone-1.2", tc.preRelease, tc.preReleaseTemplate, false, ""),
					*NewBumpAllBranchesStrategy(AUTO, DefaultPreRelease, DefaultPreReleaseTemplate, DefaultPreReleaseOverwrite, DefaultBuildMetadataTemplate),
				},
				MajorPattern: regexp.MustCompile(DefaultMajorPattern),
				MinorPattern: regexp.MustCompile(DefaultMinorPattern),
			}
			version, err := strategy.Bump()

			assert.Nil(err)
			assert.Equal(tc.expected, version.String())
		})
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
	strategy.BumpStrategies = []BumpBranchesStrategy{*NewBumpAllBranchesStrategy(AUTO, true, "SNAPSHOT", true, "")}
	version, err := strategy.Bump()

	assert.Nil(err)
	assert.Equal("1.1.0-SNAPSHOT", version.String())
}

func TestSetGitRepository(t *testing.T) {
	assert := assert.New(t)
	s := &BumpStrategy{}
	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gitRepo := mock_version.NewMockGitRepo(ctrl)
	s.SetGitRepository(gitRepo)

	assert.Equal(gitRepo, s.gitRepo)
}
func ExampleBumpStrategy_GoString() {
	gitRepo := mock_version.NewMockGitRepo(nil)
	s := NewConventionalCommitBumpStrategy(gitRepo)
	fmt.Printf("%#v\n", s)
	// Output: version.BumpStrategy{MajorPattern: &regexp.Regexp{expr: "(?:^.+\\!:.*$|(?m)^BREAKING CHANGE:.*$)"}, MinorPattern: &regexp.Regexp{expr: "^(?:feat|chore|build|ci|refactor|perf)(?:\\(.+\\))?:.*$"}, BumpBranchesStrategies: []version.BumpBranchesStrategy{version.BumpBranchesStrategy{Strategy: AUTO, BranchesPattern: &regexp.Regexp{expr: "^(master|release/.*)$"}, PreRelease: false, PreReleaseTemplate: &template.Template{text: ""}, PreReleaseOverwrite: false, BuildMetadataTemplate: &template.Template{text: ""}}, version.BumpBranchesStrategy{Strategy: AUTO, BranchesPattern: &regexp.Regexp{expr: ".*"}, PreRelease: false, PreReleaseTemplate: &template.Template{text: ""}, PreReleaseOverwrite: false, BuildMetadataTemplate: &template.Template{text: "{{.Commits | len}}.{{(.Commits | first).Hash.Short}}"}}}}
}
