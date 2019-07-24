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
	// Strategy defines the strategy to use to bump the version.
	// It can be automatic (AUTO) or manual (MAJOR, MINOR, PATCH)
	Strategy BumpStrategy
	// PreRelease defines the pre-release class (alpha, beta, etc.) for the next version
	PreRelease string
	// PreReleaseOverwrite defines if a pre-release can be overwritten
	// If true, it will not append an index to the next version
	// If false, it will append an incremented index based on the previous same version of same class if any and 0 otherwise
	PreReleaseOverwrite bool
	// BuildMetadata defines the build metadata for the next version
	BuildMetadata string
	// RegexReleaseBranches is the regex used to detect if the current branch is a release branch
	RegexReleaseBranches *regexp.Regexp
	// RegexMajor is the regex used to detect if a commit contains a breaking/major change
	// See RegexMinor for more details
	RegexMajor *regexp.Regexp
	// RegexMinor is the regex used to detect if a commit contains a minor change
	// If no commit match RegexMajor or RegexMinor, the change is considered as a patch
	RegexMinor *regexp.Regexp
	// gitRepo is an implementation of GitRepo
	gitRepo GitRepo
}

/*
NewConventionalCommitBumpStrategyOptions create a BumpStrategyOptions following https://www.conventionalcommits.org

The strategy configuration is:

	Strategy:             AUTO
	PreRelease:           ""
	PreReleaseOverwrite:  false
	BuildMetadata:        ""
	RegexReleaseBranches: ^(master|release/.*)$
	RegexMajor:           (?m)^BREAKING CHANGE:.*$
	RegexMinor:           ^feat(?:\(.+\))?:.*
*/
func NewConventionalCommitBumpStrategyOptions(gitRepo GitRepo) *BumpStrategyOptions {
	return &BumpStrategyOptions{
		Strategy:             AUTO,
		PreRelease:           "",
		PreReleaseOverwrite:  false,
		BuildMetadata:        "",
		RegexReleaseBranches: regexp.MustCompile(`^(master|release/.*)$`),
		RegexMajor:           regexp.MustCompile(`(?m)^BREAKING CHANGE:.*$`),
		RegexMinor:           regexp.MustCompile(`^feat(?:\(.+\))?:.*`),
		gitRepo:              gitRepo,
	}
}

// GoString makes BumpStrategyOptions satisfy the GoStringer interface.
func (o BumpStrategyOptions) GoString() string {
	sb := &strings.Builder{}
	sb.WriteString("version.BumpStrategyOptions{")
	sb.WriteString("Strategy: %q, PreRelease: %q, PreReleaseOverwrite: %v, BuildMetadata: %q, ")
	sb.WriteString("RegexReleaseBranches: &regexp.Regexp{expr: %q}, RegexMajor: &regexp.Regexp{expr: %q}, RegexMinor: &regexp.Regexp{expr: %q}")
	sb.WriteString("}")
	return fmt.Sprintf(sb.String(), o.Strategy, o.PreRelease, o.PreReleaseOverwrite, o.BuildMetadata,
		o.RegexReleaseBranches, o.RegexMajor, o.RegexMinor)
}

// WithPreRelease sets the pre-release name
func (o *BumpStrategyOptions) WithPreRelease(value string, override bool) *BumpStrategyOptions {
	o.PreRelease = value
	o.PreReleaseOverwrite = override
	return o
}

// Bump performs the version bumping based on the strategy
func (o *BumpStrategyOptions) Bump() (Version, error) {
	log.Debug("BumpStrategy: bump with configuration: %#v", o)

	// Make sure we have the tags
	err := o.gitRepo.FetchTags()
	if err != nil {
		return zeroVersion, newErrorC(err, "Cannot fetch tags")
	}

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

	currentBranch, err := o.gitRepo.GetCurrentBranch()
	if err != nil {
		return zeroVersion, newErrorC(err, "Cannot get current branch name")
	}

	// Check if describe is a tag, if so return the version that matches this tag
	commits, cErr := o.gitRepo.GetCommits(lastTag.Name, "HEAD")
	if cErr != nil {
		// Oops
		return zeroVersion, err
	}

	log.Debug("BumpStrategy: look for appropriate version bumper with %#v, lastVersion=%v, branch=%v", lastTag, lastVersion, currentBranch)
	versionBumper := o.computeVersionBumper(currentBranch, &lastVersion, commits)

	if o.Strategy != AUTO || len(commits) > 0 {
		if o.BuildMetadata != "" { // if BuildMetadata
			return lastVersion.WithBuildMetadata(o.BuildMetadata), nil
		} else if o.PreRelease != "" && // if PreRelease
			(o.Strategy != AUTO || o.RegexReleaseBranches.MatchString(currentBranch) || !lastVersion.HasSamePreReleaseIdentifiers(o.PreRelease)) {
			// if AUTO
			//   if branch = master/release
			//     bump
			//   else if branch != master/release && !HasSamePreReleaseIdentifiers
			//     bump
			// else bump
			// will automatically suffix the pre-release with an identifier. Eg: *-alpha.0
			return lastVersion.BumpPreRelease(o.PreRelease, o.PreReleaseOverwrite, versionBumper), nil
		}
	}

	// otherwise
	return versionBumper(lastVersion), nil
}

// computeAutoVersionBumper computes what bump strategy to apply
func (o *BumpStrategyOptions) computeVersionBumper(currentBranch string, v *Version, commits []git.Commit) versionBumper {
	if log.IsLevelEnabled(log.DebugLevel) && o.Strategy != AUTO {
		log.Debug("BumpStrategy: will use bump %s", strings.ToUpper(o.Strategy.String()))
	}

	switch o.Strategy {
	case MAJOR:
		return Version.BumpMajor
	case MINOR:
		return Version.BumpMinor
	case PATCH:
		return Version.BumpPatch
	case AUTO:
		return o.computeAutoVersionBumper(currentBranch, v, commits)
	default:
		log.Debug("BumpStrategy: will not use any bump strategy because the strategy is unknown")
		return versionBumperIdentity
	}
}

// computeAutoVersionBumper detects what bump strategy to apply based on commits history in auto mode
func (o *BumpStrategyOptions) computeAutoVersionBumper(currentBranch string, v *Version, commits []git.Commit) versionBumper {
	if len(commits) == 0 {
		log.Debug("BumpStrategy: will not use any bump strategy because there is not commit")
		return versionBumperIdentity
	}

	if !o.RegexReleaseBranches.MatchString(currentBranch) {
		log.Debug("BumpStrategy: will use build metadata strategy because the branch:%v is not a release branch matching the regex %v", currentBranch, o.RegexReleaseBranches)
		return func(v Version) Version {
			lastCommitShortHash := commits[0].Hash.Short()
			count := len(commits)
			return v.WithBuildMetadata(fmt.Sprintf("%d.%s", count, lastCommitShortHash.String()))
		}
	}

	strategy := PATCH
	bumper := Version.BumpPatch
	for _, commit := range commits {
		if o.RegexMajor.MatchString(commit.Message) {
			if v.IsUnstable() {
				log.Trace("BumpStrategy: detects a MAJOR change at %#v however the last version is unstable so it will use bump MINOR strategy", commit)
				return Version.BumpMinor
			}
			log.Trace("BumpStrategy: detects a MAJOR change at %#v", commit)
			return Version.BumpMajor
		} else if o.RegexMinor.MatchString(commit.Message) {
			strategy = MINOR
			log.Trace("BumpStrategy: detects a MINOR change at %#v", commit)
			bumper = Version.BumpMinor
		}
	}
	log.Debug("BumpStrategy: will use bump %s strategy", strategy)
	return bumper
}
