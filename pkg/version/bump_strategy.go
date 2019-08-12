package version

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/arnaud-deprez/gsemver/internal/log"
)

// versionBumper type helper for the bump process
type versionBumper func(Version) Version

// BumpStrategy allows you to configure the bump strategy
type BumpStrategy struct {
	// Strategy defines the strategy to use to bump the version.
	// It can be automatic (AUTO) or manual (MAJOR, MINOR, PATCH)
	Strategy BumpStrategyType `json:"strategy"`
	// MajorPattern is the regex used to detect if a commit contains a breaking/major change
	// See RegexMinor for more details
	MajorPattern *regexp.Regexp `json:"majorPattern,omitempty"`
	// MinorPattern is the regex used to detect if a commit contains a minor change
	// If no commit match RegexMajor or RegexMinor, the change is considered as a patch
	MinorPattern *regexp.Regexp `json:"minorPattern,omitempty"`
	// BumpBranchesStrategies is a list of bump strategies for matching branches
	BumpBranchesStrategies []BumpBranchesStrategy `json:"bumpBranchesStrategies,omitempty"`
	// BumpDefaultStrategy is a default bump strategy to apply when BumpReleaseStrategies matched.
	BumpDefaultStrategy *BumpBranchesStrategy `json:"bumpDefaultStrategy,omitempty"`
	// gitRepo is an implementation of GitRepo
	gitRepo GitRepo
}

/*
NewConventionalCommitBumpStrategy create a BumpStrategy following https://www.conventionalcommits.org

The strategy configuration is:

	Strategy: AUTO
	BumpBranchesStrategies: [
		{
			PatternReleaseBranches: ^(master|release/.*)$
			PreRelease:             false
			PreReleaseTemplate:     ""
			PreReleaseOverwrite:    false
			BuildMetadataTemplate:  ""
		}
	]
	BumpDefaultStrategy: {
		PatternReleaseBranches: .*
		PreRelease:    		    false
		PreReleaseTemplate:     ""
		PreReleaseOverwrite:    false
		BuildMetadataTemplate:  ""
	}
	MajorPattern: (?m)^BREAKING CHANGE:.*$
	MinorPattern: ^feat(?:\(.+\))?:.*
*/
func NewConventionalCommitBumpStrategy(gitRepo GitRepo) *BumpStrategy {
	return &BumpStrategy{
		Strategy: AUTO,
		BumpBranchesStrategies: []BumpBranchesStrategy{
			*NewDefaultBumpBranchesStrategy(`^(master|release/.*)$`),
		},
		BumpDefaultStrategy: NewFallbackBumpBranchesStrategy(false, "", false, "{{ .Commits | len }}.{{ (.Commits | first).Hash.Short }}"),
		MajorPattern:        regexp.MustCompile(`(?m)^BREAKING CHANGE:.*$`),
		MinorPattern:        regexp.MustCompile(`^feat(?:\(.+\))?:.*`),
		gitRepo:             gitRepo,
	}
}

// GoString makes BumpStrategy satisfy the GoStringer interface.
func (o BumpStrategy) GoString() string {
	var sb strings.Builder
	sb.WriteString("version.BumpStrategy{")
	sb.WriteString(fmt.Sprintf("Strategy: %q, MajorPattern: &regexp.Regexp{expr: %q}, MinorPattern: &regexp.Regexp{expr: %q}, ", o.Strategy, o.MajorPattern, o.MinorPattern))
	sb.WriteString(fmt.Sprintf("BumpBranchesStrategies: %#v", o.BumpBranchesStrategies))
	sb.WriteString(fmt.Sprintf("BumpDefaultStrategy: %#v", o.BumpDefaultStrategy))
	sb.WriteString("}")
	return sb.String()
}

// AddBumpBranchesStrategy add a bump strategy for a matching set of branches
func (o *BumpStrategy) AddBumpBranchesStrategy(s *BumpBranchesStrategy) *BumpStrategy {
	o.BumpBranchesStrategies = append(o.BumpBranchesStrategies, *s)
	return o
}

// WithBumpDevelopmentStrategy set a bump development strategy
func (o *BumpStrategy) WithBumpDevelopmentStrategy(s *BumpBranchesStrategy) *BumpStrategy {
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

	context := NewContext(currentBranch, &lastVersion, &lastTag, commits)

	log.Debug("BumpStrategy: look for appropriate version bumper with %#v, lastVersion=%v, branch=%v", lastTag, lastVersion, currentBranch)
	versionBumper := o.computeVersionBumper(context)

	// Bump the version
	return versionBumper(lastVersion), nil
}

// computeAutoVersionBumper computes what bump strategy to apply
func (o *BumpStrategy) computeVersionBumper(context *Context) versionBumper {
	if log.IsLevelEnabled(log.DebugLevel) && o.Strategy != AUTO {
		log.Debug("BumpStrategy: will use bump %s", strings.ToUpper(o.Strategy.String()))
	}

	switch o.Strategy {
	case MAJOR:
		return o.BumpDefaultStrategy.createVersionBumperFrom(Version.BumpMajor, context)
	case MINOR:
		return o.BumpDefaultStrategy.createVersionBumperFrom(Version.BumpMinor, context)
	case PATCH:
		return o.BumpDefaultStrategy.createVersionBumperFrom(Version.BumpPatch, context)
	case AUTO:
		return o.computeAutoVersionBumper(context)
	default:
		log.Debug("BumpStrategy: will not use any bump strategy because the strategy is unknown")
		return versionBumperIdentity
	}
}

// computeAutoVersionBumper detects what bump strategy to apply based on commits history in auto mode
func (o *BumpStrategy) computeAutoVersionBumper(context *Context) versionBumper {
	if len(context.Commits) == 0 {
		log.Debug("BumpStrategy: will not use identity bump strategy because there is not commit")
		return versionBumperIdentity
	}

	semverBumper := o.computeSemverBumperFromCommits(context)

	for _, it := range o.BumpBranchesStrategies {
		if it.BranchesPattern.MatchString(context.Branch) {
			return it.createVersionBumperFrom(semverBumper, context)
		}
	}
	return o.BumpDefaultStrategy.createVersionBumperFrom(semverBumper, context)
}

func (o *BumpStrategy) computeSemverBumperFromCommits(context *Context) versionBumper {
	strategy := PATCH
	bumper := Version.BumpPatch
	for _, commit := range context.Commits {
		if o.MajorPattern.MatchString(commit.Message) {
			if context.LastVersion.IsUnstable() {
				log.Trace("BumpStrategy: detects a MAJOR change at %#v however the last version is unstable so it will use bump MINOR strategy", commit)
				return Version.BumpMinor
			}
			log.Trace("BumpStrategy: detects a MAJOR change at %#v", commit)
			return Version.BumpMajor
		} else if o.MinorPattern.MatchString(commit.Message) {
			strategy = MINOR
			log.Trace("BumpStrategy: detects a MINOR change at %#v", commit)
			bumper = Version.BumpMinor
		}
	}
	log.Debug("BumpStrategy: will use bump %s strategy", strategy)
	return bumper
}
