package version

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/arnaud-deprez/gsemver/internal/log"
)

const (
	// DefaultMajorPattern defines default regular expression to match a commit message with a major change.
	DefaultMajorPattern = `(?:^.+\!:.+|(?m)^BREAKING CHANGE:.+$)`
	// DefaultMinorPattern defines default regular expression to match a commit message with a minor change.
	DefaultMinorPattern = `^(?:feat|chore|build|ci|refactor|perf)(?:\(.+\))?:.+`
	// DefaultReleaseBranchesPattern defines default regular expression to match release branches
	DefaultReleaseBranchesPattern = `^(main|master|release/.*)$`
	// DefaultPreRelease defines default pre-release activation for non release branches
	DefaultPreRelease = false
	// DefaultPreReleaseTemplate defines default pre-release go template for non release branches
	DefaultPreReleaseTemplate = ""
	// DefaultPreReleaseOverwrite defines default pre-release overwrite activation for non release branches
	DefaultPreReleaseOverwrite = false
	// DefaultBuildMetadataTemplate defines default go template used for non release branches strategy
	DefaultBuildMetadataTemplate = `{{.Commits | len}}.{{(.Commits | first).Hash.Short}}`
)

var (
	// strategyVersionBumperMap defined the link between BumpStrategyType and versionBumper.
	// Note that AUTO BumpStrategyType versionBumper is dynamically computed and therefore cannot be part of this static links
	/* const */ strategyVersionBumperMap = map[BumpStrategyType]versionBumper{
		MAJOR: Version.BumpMajor,
		MINOR: Version.BumpMinor,
		PATCH: Version.BumpPatch,
	}
)

// versionBumper type helper for the bump process
type versionBumper func(Version) Version

// BumpStrategy allows you to configure the bump strategy
type BumpStrategy struct {
	// MajorPattern is the regex used to detect if a commit contains a breaking/major change
	// See RegexMinor for more details
	MajorPattern *regexp.Regexp `json:"majorPattern,omitempty"`
	// MinorPattern is the regex used to detect if a commit contains a minor change
	// If no commit match RegexMajor or RegexMinor, the change is considered as a patch
	MinorPattern *regexp.Regexp `json:"minorPattern,omitempty"`
	// BumpStrategies is a list of bump strategies for matching branches
	BumpStrategies []BumpBranchesStrategy `json:"bumpStrategies,omitempty"`
	// gitRepo is an implementation of GitRepo
	gitRepo GitRepo
}

/*
NewConventionalCommitBumpStrategy create a BumpStrategy following https://www.conventionalcommits.org

The strategy configuration is:

	MajorPattern: (?:^.+\!:.+|(?m)^BREAKING CHANGE:.+$)
	MinorPattern: ^(?:feat|chore|build|ci|refactor|perf)(?:\(.+\))?:.+
	BumpBranchesStrategies: [
		{
			Strategy: AUTO
			BranchesPattern: ^(main|master|release/.*)$
			PreRelease:             false
			PreReleaseTemplate:     ""
			PreReleaseOverwrite:    false
			BuildMetadataTemplate:  ""
		},
		{
			Strategy: AUTO
			BranchesPattern: .*
			PreRelease:    		    false
			PreReleaseTemplate:     ""
			PreReleaseOverwrite:    false
			BuildMetadataTemplate:  "{{.Commits | len}}.{{(.Commits | first).Hash.Short}}"
		}
	]
*/
func NewConventionalCommitBumpStrategy(gitRepo GitRepo) *BumpStrategy {
	return &BumpStrategy{
		BumpStrategies: []BumpBranchesStrategy{
			*NewDefaultBumpBranchesStrategy(DefaultReleaseBranchesPattern),
			*NewBuildBumpBranchesStrategy(".*", DefaultBuildMetadataTemplate),
		},
		MajorPattern: regexp.MustCompile(DefaultMajorPattern),
		MinorPattern: regexp.MustCompile(DefaultMinorPattern),
		gitRepo:      gitRepo,
	}
}

// GoString makes BumpStrategy satisfy the GoStringer interface.
func (o BumpStrategy) GoString() string {
	var sb strings.Builder
	sb.WriteString("version.BumpStrategy{")
	sb.WriteString(fmt.Sprintf("MajorPattern: &regexp.Regexp{expr: %q}, MinorPattern: &regexp.Regexp{expr: %q}, ", o.MajorPattern, o.MinorPattern))
	sb.WriteString(fmt.Sprintf("BumpBranchesStrategies: %#v", o.BumpStrategies))
	sb.WriteString("}")
	return sb.String()
}

// SetGitRepository configures the git repository to use for the strategy
func (o *BumpStrategy) SetGitRepository(gitRepo GitRepo) {
	o.gitRepo = gitRepo
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
		log.Debug("%v", newErrorC(err, "Unable to get last relative tag"))
	}

	// Parse the last version from the tag name
	lastVersion, err := NewVersion(extractVersionFromTag(lastTag.Name))
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

func extractVersionFromTag(tagName string) string {
	return tagName[strings.LastIndex(tagName, "/")+1:]
}

// computeAutoVersionBumper computes what bump strategy to apply
func (o *BumpStrategy) computeVersionBumper(context *Context) versionBumper {
	for _, it := range o.BumpStrategies {
		if it.BranchesPattern.MatchString(context.Branch) {
			if log.IsLevelEnabled(log.DebugLevel) {
				log.Debug("BumpStrategy: will use bump %s", strings.ToUpper(it.Strategy.String()))
			}

			// find the correct bumper
			if val, ok := strategyVersionBumperMap[it.Strategy]; ok {
				return it.createVersionBumperFrom(val, context)
			} else if it.Strategy == AUTO {
				return o.computeSemverBumperFromCommits(&it, context)
			}
		}
	}

	log.Debug("BumpStrategy: not matching strategy found in %#v. versionBumperIdentity will be used", o.BumpStrategies)
	return versionBumperIdentity
}

func (o *BumpStrategy) computeSemverBumperFromCommits(bbs *BumpBranchesStrategy, context *Context) versionBumper {
	if len(context.Commits) == 0 {
		log.Debug("BumpStrategy: will not use identity bump strategy because there is not commit")
		return versionBumperIdentity
	}

	strategy := PATCH
	bumper := Version.BumpPatch
	for _, commit := range context.Commits {
		if o.MajorPattern.MatchString(commit.Message) {
			if context.LastVersion.IsUnstable() {
				log.Trace("BumpStrategy: detects a MAJOR change at %#v however the last version is unstable so it will use bump MINOR strategy", commit)
				return bbs.createVersionBumperFrom(Version.BumpMinor, context)
			}
			log.Debug("BumpStrategy: detects a MAJOR change at %#v", commit)
			return bbs.createVersionBumperFrom(Version.BumpMajor, context)
		}
		if o.MinorPattern.MatchString(commit.Message) {
			strategy = MINOR
			log.Trace("BumpStrategy: detects a MINOR change at %#v", commit)
			bumper = Version.BumpMinor
		}
	}

	log.Debug("BumpStrategy: will use bump %s strategy", strategy)
	return bbs.createVersionBumperFrom(bumper, context)
}
