package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/arnaud-deprez/gsemver/internal/git"
	"github.com/arnaud-deprez/gsemver/internal/log"
	"github.com/arnaud-deprez/gsemver/pkg/version"
)

const (
	bumpDesc = `
This will compute and print the next semver compatible version of your project based on commits logs, tags and current branch.

The version will look like <X>.<Y>.<Z>[-<pre-release>][+<metadata>] where:
- X is the Major number
- Y is the Minor number
- Z is the Patch number
- pre-release is the pre-release identifiers (optional)
- metadata is the build metadata identifiers (optional)

More info on the semver spec https://semver.org/spec/v2.0.0.html.

It can work in 2 fashions, the automatic or manual.

Automatic way assumes: 
- your previous tags are semver compatible.
- you follow some conventions in your commit and ideally https://www.conventionalcommits.org
- you follow some branch convention for your releases (eg. a release should be done on master or release/* branches) 

Base on this information, it is able to compute the next version.

The manual way is less restrictive and just assumes your previous tags are semver compatible.
`
	bumpExample = `
# To bump automatically:
gsemver bump

# Or more explicitly
gsemver bump auto

# To bump manually the major number:
gsemver bump major

# To bump manually the minor number:
gsemver bump minor

# To bump manually the patch number:
gsemver bump patch

# To use a pre-release version
gsemver bump --pre-release alpha

# To use a pre-release version without indexation (maven like SNAPSHOT)
gsemver bump minor --pre-release SNAPSHOT --pre-release-overwrite true

# To use version with build metadata
gsemver bump --build-metadata "issue-1.build.1"

# To use bump auto with one or many branch strategies
gsemver bump --branch-strategy='{"branchPattern":"^miletone-1.1$","preReleaseTemplate":"beta"}' --branch-strategy='{"branchPattern":"^miletone-2.0$","preReleaseTemplate":"alpha"}'
`
)

// NewBumpCommands create the bump command with its subcommands
func newBumpCommands(globalOpts *globalOptions) *cobra.Command {
	options := &bumpOptions{
		globalOptions: globalOpts,
	}

	cmd := &cobra.Command{
		Use:       "bump [strategy]",
		Short:     "Bump to next version",
		Long:      bumpDesc,
		Example:   bumpExample,
		ValidArgs: []string{"auto", "major", "minor", "patch"},
		Args:      cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.configureLogger()

			options.Cmd = cmd
			options.Args = args
			if len(args) == 0 {
				options.Bump = "auto"
			} else {
				options.Bump = args[0]
			}

			return options.run()
		},
	}

	options.addBumpFlags(cmd)

	return cmd
}

// BumpOptions type to represent the available options for the bump commands
// It extends GlobalOptions.
type bumpOptions struct {
	*globalOptions
	// Bump is mapped to pkg/version/BumpStrategyOptions#Strategy
	Bump string
	// PreReleaseTemplate is mapped to pkg/version/BumpStrategyOptions#PreRelease
	PreReleaseTemplate string
	// PreReleaseOverwrite is mapped to pkg/version/BumpStrategyOptions#PreReleaseOverwrite
	PreReleaseOverwrite bool
	// BuildMetadataTemplate is mapped to pkg/version/BumpStrategyOptions#BuildMetadata
	BuildMetadataTemplate string
	BranchStrategies      []string
}

func (o *bumpOptions) addBumpFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.PreReleaseTemplate, "pre-release", "", "", "Use pre-release template version such as `alpha` which will give a version like `X.Y.Z-alpha.N`. This flag is not taken into account if --build-metadata is set.")
	cmd.Flags().BoolVarP(&o.PreReleaseOverwrite, "pre-release-overwrite", "", false, "Use pre-release overwrite option to remove the pre-release identifier suffix which will give a version like `X.Y.Z-SNAPSHOT` if pre-release=SNAPSHOT")
	cmd.Flags().StringVarP(&o.BuildMetadataTemplate, "build-metadata", "", "", "Use build metadata template which will give something like X.Y.Z+<build>. This flag cannot be used with --pre-release* flags and take precedence over them.")
	cmd.Flags().StringArrayVarP(&o.BranchStrategies, "branch-strategy", "", []string{}, `Use branch-strategy will set a strategy for a set of branches. 
	The strategy is defined in json and looks like {"branchPattern":"^milestone-.*$", "preReleaseTemplate":"alpha"} to use pre-release alpha version for every milestone-* branches.`)

	o.Cmd = cmd
}

func (o *bumpOptions) run() error {
	log.Debug("Run bump command with configuration: %#v", o)

	version, err := o.createBumpStrategy().Bump()
	if err != nil {
		return err
	}
	fmt.Fprintf(o.globalOptions.ioStreams.Out, "%v", version)
	return nil
}

func (o *bumpOptions) createBumpStrategy() *version.BumpStrategy {
	ret := version.NewConventionalCommitBumpStrategy(git.NewVersionGitRepo(o.CurrentDir))
	ret.Strategy = version.ParseBumpStrategyType(o.Bump)
	ret.BumpDefaultStrategy.PreReleaseTemplate = o.PreReleaseTemplate
	ret.BumpDefaultStrategy.PreReleaseOverwrite = o.PreReleaseOverwrite
	ret.BumpDefaultStrategy.BuildMetadataTemplate = o.BuildMetadataTemplate

	for _, s := range o.
		BranchStrategies {
		var b version.BumpReleaseStrategy
		json.Unmarshal([]byte(s), &b)
		ret.BumpReleaseStrategies = append(ret.BumpReleaseStrategies, b)
	}
	return ret
}
