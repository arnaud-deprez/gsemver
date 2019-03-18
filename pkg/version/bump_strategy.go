package version

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/arnaud-deprez/gsemver/internal/log"
	"github.com/arnaud-deprez/gsemver/pkg/git"
)

// BumpStrategy represents the SemVer number that needs to be bumped
type BumpStrategy int

const (
	// PATCH means to bump the patch number
	PATCH BumpStrategy = iota
	// MINOR means to bump the minor number
	MINOR
	// MAJOR means to bump the patch number
	MAJOR
	// AUTO means to apply the automatic strategy based on commit history
	AUTO
)

func (b BumpStrategy) String() string {
	return [...]string{"PATCH", "MINOR", "MAJOR", "AUTO"}[b]
}

// ParseBumpStrategy converts string value to BumpStrategy
func ParseBumpStrategy(value string) BumpStrategy {
	switch strings.ToLower(value) {
	case "major":
		return MAJOR
	case "minor":
		return MINOR
	case "patch":
		return PATCH
	default:
		return AUTO
	}
}

// BumpStrategyOptions allows you to configure the bump strategy
type BumpStrategyOptions struct {
	Strategy             BumpStrategy
	PreRelease           string
	PreReleaseOverwrite  bool
	BuildMetadata        string
	RegexReleaseBranches *regexp.Regexp
	RegexMajor           *regexp.Regexp
	RegexMinor           *regexp.Regexp
	gitRepo              GitRepo
}

// NewConventionalCommitBumpStrategyOptions create a BumpStrategyOptions following https://www.conventionalcommits.org:
func NewConventionalCommitBumpStrategyOptions(gitRepo GitRepo) *BumpStrategyOptions {
	return &BumpStrategyOptions{
		Strategy:             AUTO,
		PreRelease:           "",
		PreReleaseOverwrite:  false,
		BuildMetadata:        "",
		RegexReleaseBranches: regexp.MustCompile(`^(master|release/.*)$`),
		RegexMajor:           regexp.MustCompile(`(?m)^BREAKING CHANGE:.*$`),
		RegexMinor:           regexp.MustCompile(`^feat(?:\(.+\))?:.*$`),
		gitRepo:              gitRepo,
	}
}

// WithPreRelease sets the pre-release name
func (o *BumpStrategyOptions) WithPreRelease(value string, override bool) *BumpStrategyOptions {
	o.PreRelease = value
	o.PreReleaseOverwrite = override
	return o
}

// Bump performs the version bumping based on the strategy
func (o *BumpStrategyOptions) Bump() (Version, error) {
	// This assumes we used annotated tags for the release. Annotated tag are created with: git tag -a -m "<message>" <tag>
	// Annotated tags adds timestamp, author and message to a tag compared to lightweight tag which does not contain any of these information.
	// Thanks to that git describe will only show the more recent annotated tag if many annotated tags are on the same commit.
	// However if you use lightweight tags there are many on the same commit, it just takes the first one.
	lastTag, err := o.gitRepo.GetLastRelativeTag("HEAD")
	if err != nil {
		// just log for debug but the program can continue
		log.Debug("Unable to get last relative tag because '%s'", err)
	}

	// Parse the last version from the tag name
	lastVersion, err := NewVersion(lastTag.Name)
	if err != nil {
		return zeroVersion, err
	}

	// Check if describe is a tag, if so return the version that matches this tag
	commits, cErr := o.gitRepo.GetCommits(lastTag.Name, "HEAD")
	if cErr != nil {
		// Oops
		return zeroVersion, err
	}

	var versionBumper versionBumper
	// If strategy is auto, we should convert it to MAJOR, MINOR or PATCH
	if o.Strategy == AUTO {
		versionBumper = o.detectVersionBumper(&lastVersion, commits)
	} else {
		switch o.Strategy {
		case MAJOR:
			versionBumper = Version.BumpMajor
			break
		case MINOR:
			versionBumper = Version.BumpMinor
			break
		case PATCH:
			versionBumper = Version.BumpPatch
			break
		default:
			return zeroVersion, Error{message: fmt.Sprintf("Unable to create versionBumper with strategy <%v>", o.Strategy)}
		}
	}

	if o.Strategy != AUTO || len(commits) > 0 {
		if o.BuildMetadata != "" { // if BuildMetadata
			return lastVersion.WithBuildMetadata(o.BuildMetadata), nil
		} else if o.PreRelease != "" { // if pre-release bump
			// will automatically suffix the pre-release with an identifier. Eg: *-alpha.0
			return lastVersion.BumpPreRelease(o.PreRelease, o.PreReleaseOverwrite, versionBumper), nil
		}
	}

	// otherwise
	return versionBumper(lastVersion), nil
}

// detectVersionBumper detects what bump strategy to apply based on commits history
func (o *BumpStrategyOptions) detectVersionBumper(v *Version, commits []git.Commit) versionBumper {
	if len(commits) == 0 {
		return versionBumperIdentity
	}

	bumper := Version.BumpPatch
	for _, commit := range commits {
		if o.RegexMajor.MatchString(commit.Message) {
			if v.IsUnstable() {
				return Version.BumpMinor
			}
			return Version.BumpMajor
		} else if o.RegexMinor.MatchString(commit.Message) {
			bumper = Version.BumpMinor
		}
	}
	return bumper
}
