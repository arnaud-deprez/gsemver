package version

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/arnaud-deprez/gsemver/internal/log"
	"github.com/arnaud-deprez/gsemver/pkg/git"
)

// versionBumper type helper for the bump process
type versionBumper func(Version) Version

// BumpStrategy allows you to configure the bump strategy
type BumpStrategy struct {
	// Strategy defines the strategy to use to bump the version.
	// It can be automatic (AUTO) or manual (MAJOR, MINOR, PATCH)
	Strategy BumpStrategyType `json:"strategy"`
	// RegexMajor is the regex used to detect if a commit contains a breaking/major change
	// See RegexMinor for more details
	PatternMajor *regexp.Regexp `json:"patternMajor,omitempty"`
	// RegexMinor is the regex used to detect if a commit contains a minor change
	// If no commit match RegexMajor or RegexMinor, the change is considered as a patch
	PatternMinor *regexp.Regexp `json:"patternMinor,omitempty"`
	// BumpReleaseStrategies is a list of bump strategies for matching release branches
	BumpReleaseStrategies []BumpReleaseStrategy `json:"bumpReleaseStrategies,omitempty"`
	// BumpDefaultStrategy is a default bump strategy to apply when BumpReleaseStrategies matched.
	BumpDefaultStrategy *BumpDefaultStrategy `json:"bumpDefaultStrategy,omitempty"`
	// gitRepo is an implementation of GitRepo
	gitRepo GitRepo
}

// BumpBranchStrategy defines a method to create a versionBumper from another one
type BumpBranchStrategy interface {
	// createVersionBumperFrom create a new version bumper for the strategy based on another one (usually MAJOR, MINOR or PATCH bumpers)
	createVersionBumperFrom(bumper versionBumper) versionBumper
}

// BumpReleaseStrategy allows you to configure the bump strategy option for a matching release branch.
type BumpReleaseStrategy struct {
	BumpBranchStrategy
	// PatternReleaseBranches is the regex used to detect if the current branch is a release branch
	PatternReleaseBranches *regexp.Regexp `json:"patternReleaseBranches,omitempty"`
	// PreReleaseTemplate defines the pre-release template for the next version
	// It can be alpha, beta, or a go-template expression
	PreReleaseTemplate string `json:"preReleaseTemplate,omitempty"`
}

// createVersionBumperFrom is an implementation for BumpBranchStrategy
func (o *BumpReleaseStrategy) createVersionBumperFrom(bumper versionBumper) versionBumper {
	return func(v Version) Version {
		if o != nil && o.PreReleaseTemplate != "" {
			return v.BumpPreRelease(o.PreReleaseTemplate, false, bumper)
		}
		return bumper(v)
	}
}

// BumpDefaultStrategy allows you to configure the bump  default strategy option when no BumpReleaseStrategy matched.
type BumpDefaultStrategy struct {
	BumpBranchStrategy
	// PreReleaseTemplate defines the pre-release template for the next version
	// It can be alpha, beta, or a go-template expression
	PreReleaseTemplate string `json:"preReleaseTemplate,omitempty"`
	// PreReleaseOverwrite defines if a pre-release can be overwritten
	// If true, it will not append an index to the next version
	// If false, it will append an incremented index based on the previous same version of same class if any and 0 otherwise
	PreReleaseOverwrite bool `json:"preReleaseOverwrite"`
	// BuildMetadataTemplate defines the build metadata for the next version.
	// It can be a static value but it will usually be a go-template expression to guarantee uniqueness of each built version.
	BuildMetadataTemplate string `json:"buildMetadataTemplate,omitempty"`
}

// createVersionBumperFrom is an implementation for BumpBranchStrategy
func (o *BumpDefaultStrategy) createVersionBumperFrom(bumper versionBumper) versionBumper {
	return func(v Version) Version {
		if o != nil && o.BuildMetadataTemplate != "" { // if BuildMetadataTemplate
			return v.WithBuildMetadata(o.BuildMetadataTemplate)
		} else if o != nil && o.PreReleaseTemplate != "" {
			return v.BumpPreRelease(o.PreReleaseTemplate, o.PreReleaseOverwrite, bumper)
		}
		return bumper(v)
	}
}

/*
NewConventionalCommitBumpStrategy create a BumpStrategy following https://www.conventionalcommits.org

The strategy configuration is:

	Strategy: AUTO
	BumpReleaseStrategies: [
		{
			PatternReleaseBranches: ^(master|release/.*)$
			PreReleaseTemplate:     ""
		}
	]
	BumpDefaultStrategy: {
		PreReleaseTemplate:    ""
		PreReleaseOverwrite:   false
		BuildMetadataTemplate: ""
	}
	PatternMajor: (?m)^BREAKING CHANGE:.*$
	PatternMinor: ^feat(?:\(.+\))?:.*
*/
func NewConventionalCommitBumpStrategy(gitRepo GitRepo) *BumpStrategy {
	return &BumpStrategy{
		Strategy: AUTO,
		BumpReleaseStrategies: []BumpReleaseStrategy{
			{
				PatternReleaseBranches: regexp.MustCompile(`^(master|release/.*)$`),
				PreReleaseTemplate:     "",
			},
		},
		BumpDefaultStrategy: &BumpDefaultStrategy{
			PreReleaseTemplate:  "",
			PreReleaseOverwrite: false,
			//TODO: we should have a default value here
			BuildMetadataTemplate: "",
		},
		PatternMajor: regexp.MustCompile(`(?m)^BREAKING CHANGE:.*$`),
		PatternMinor: regexp.MustCompile(`^feat(?:\(.+\))?:.*`),
		gitRepo:      gitRepo,
	}
}

// GoString makes BumpStrategy satisfy the GoStringer interface.
func (o BumpStrategy) GoString() string {
	var sb strings.Builder
	sb.WriteString("version.BumpStrategy{")
	sb.WriteString(fmt.Sprintf("Strategy: %q, PatternMajor: &regexp.Regexp{expr: %q}, PatternMinor: &regexp.Regexp{expr: %q}", o.Strategy, o.PatternMajor, o.PatternMajor))
	// TODO: implements the rest of GoString()
	sb.WriteString("}")
	return sb.String()
}

// AddBumpReleaseStrategy add a bump release strategy for a branch pattern
func (o *BumpStrategy) AddBumpReleaseStrategy(pattern string, preReleaseTemplate string) *BumpStrategy {
	s := BumpReleaseStrategy{
		PatternReleaseBranches: regexp.MustCompile(pattern),
		PreReleaseTemplate:     preReleaseTemplate,
	}
	o.BumpReleaseStrategies = append(o.BumpReleaseStrategies, s)
	return o
}

// WithBumpDevelopmentStrategy set a bump development strategy
func (o *BumpStrategy) WithBumpDevelopmentStrategy(preReleaseTemplate string, preReleaseOverwrite bool, buildMetadataTemplate string) *BumpStrategy {
	s := &BumpDefaultStrategy{
		PreReleaseTemplate:    preReleaseTemplate,
		PreReleaseOverwrite:   preReleaseOverwrite,
		BuildMetadataTemplate: buildMetadataTemplate,
	}
	o.BumpDefaultStrategy = s
	return o
}

// Bump performs the version bumping based on the strategy
func (o *BumpStrategy) Bump() (Version, error) {
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

	// Bump the version
	return versionBumper(lastVersion), nil
}

// computeAutoVersionBumper computes what bump strategy to apply
func (o *BumpStrategy) computeVersionBumper(currentBranch string, v *Version, commits []git.Commit) versionBumper {
	if log.IsLevelEnabled(log.DebugLevel) && o.Strategy != AUTO {
		log.Debug("BumpStrategy: will use bump %s", strings.ToUpper(o.Strategy.String()))
	}

	switch o.Strategy {
	case MAJOR:
		return o.BumpDefaultStrategy.createVersionBumperFrom(Version.BumpMajor)
	case MINOR:
		return o.BumpDefaultStrategy.createVersionBumperFrom(Version.BumpMinor)
	case PATCH:
		return o.BumpDefaultStrategy.createVersionBumperFrom(Version.BumpPatch)
	case AUTO:
		return o.computeAutoVersionBumper(currentBranch, v, commits)
	default:
		log.Debug("BumpStrategy: will not use any bump strategy because the strategy is unknown")
		return versionBumperIdentity
	}
}

// computeAutoVersionBumper detects what bump strategy to apply based on commits history in auto mode
func (o *BumpStrategy) computeAutoVersionBumper(currentBranch string, v *Version, commits []git.Commit) versionBumper {
	if len(commits) == 0 {
		log.Debug("BumpStrategy: will not use identity bump strategy because there is not commit")
		return versionBumperIdentity
	}

	semverBumper := o.computeSemverBumperFromCommits(currentBranch, v, commits)

	for _, it := range o.BumpReleaseStrategies {
		if it.PatternReleaseBranches.MatchString(currentBranch) {
			return it.createVersionBumperFrom(semverBumper)
		}
	}
	return o.BumpDefaultStrategy.createVersionBumperFrom(semverBumper)
}

func (o *BumpStrategy) computeSemverBumperFromCommits(currentBranch string, v *Version, commits []git.Commit) versionBumper {
	strategy := PATCH
	bumper := Version.BumpPatch
	for _, commit := range commits {
		if o.PatternMajor.MatchString(commit.Message) {
			if v.IsUnstable() {
				log.Trace("BumpStrategy: detects a MAJOR change at %#v however the last version is unstable so it will use bump MINOR strategy", commit)
				return Version.BumpMinor
			}
			log.Trace("BumpStrategy: detects a MAJOR change at %#v", commit)
			return Version.BumpMajor
		} else if o.PatternMinor.MatchString(commit.Message) {
			strategy = MINOR
			log.Trace("BumpStrategy: detects a MINOR change at %#v", commit)
			bumper = Version.BumpMinor
		}
	}
	log.Debug("BumpStrategy: will use bump %s strategy", strategy)
	return bumper
}
