package cmd

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/arnaud-deprez/gsemver/internal/git"
	"github.com/arnaud-deprez/gsemver/internal/log"
	"github.com/arnaud-deprez/gsemver/internal/utils"
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
gsemver bump --build "issue-1.build.1"

# To use bump auto with one or many branch strategies
gsemver bump --branch-strategy='{"branchesPattern":"^miletone-1.1$","preReleaseTemplate":"beta"}' --branch-strategy='{"branchesPattern":"^miletone-2.0$","preReleaseTemplate":"alpha"}'
`
	preReleaseTemplateDesc = `Use pre-release template version such as 'alpha' which will give a version like 'X.Y.Z-alpha.N'.
If pre-release flag is present but does not contain template value, it will give a version like 'X.Y.Z-N' where 'N' is the next pre-release increment for the version 'X.Y.Z'.
You can also use go-template expression with context https://godoc.org/github.com/arnaud-deprez/gsemver/pkg/version#Context and http://masterminds.github.io/sprig functions.
This flag is not taken into account if --build-metadata is set.`

	buildMetadataTemplateDesc = `Use build metadata template which will give something like X.Y.Z+<build>.
You can also use go-template expression with context https://godoc.org/github.com/arnaud-deprez/gsemver/pkg/version#Context and http://masterminds.github.io/sprig functions.
This flag cannot be used with --pre-release* flags and take precedence over them.`

	branchStrategyDesc = `Use branch-strategy will set a strategy for a set of branches. 
The strategy is defined in json and looks like {"branchesPattern":"^milestone-.*$", "preReleaseTemplate":"alpha"} for example.
This will use pre-release alpha version for every milestone-* branches. 
You can find all available options https://godoc.org/github.com/arnaud-deprez/gsemver/pkg/version#BumpBranchesStrategy`
)

// newBumpCommands create the bump command with its subcommands
func newBumpCommands(globalOpts *globalOptions) *cobra.Command {
	return newBumpCommandsWithRun(globalOpts, run)
}

// newBumpCommandsWithRun create the bump commands from bumpOptions with run function.
// it is used for internal usage only and for test purpose.
func newBumpCommandsWithRun(globalOpts *globalOptions, run func(o *bumpOptions) error) *cobra.Command {
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

			options.PreRelease = cmd.Flags().Changed("pre-release")

			return run(options)
		},
	}

	options.addBumpFlags(cmd)

	return cmd
}

type config struct {
	MajorPattern   string
	MinorPattern   string
	BumpStrategies []struct {
		Strategy              string
		BranchesPattern       string
		PreRelease            bool
		PreReleaseTemplate    string
		PreReleaseOverwrite   bool
		BuildMetadataTemplate string
	}
}

func (c *config) createBumpStrategy() *version.BumpStrategy {
	ret := version.BumpStrategy{BumpStrategies: []version.BumpBranchesStrategy{}}
	ret.MajorPattern = regexp.MustCompile(c.MajorPattern)
	ret.MinorPattern = regexp.MustCompile(c.MinorPattern)
	for _, it := range c.BumpStrategies {
		s := version.BumpBranchesStrategy{
			Strategy:              version.ParseBumpStrategyType(it.Strategy),
			BranchesPattern:       regexp.MustCompile(it.BranchesPattern),
			PreRelease:            it.PreRelease,
			PreReleaseTemplate:    utils.NewTemplate(it.PreReleaseTemplate),
			PreReleaseOverwrite:   it.PreReleaseOverwrite,
			BuildMetadataTemplate: utils.NewTemplate(it.BuildMetadataTemplate),
		}
		ret.BumpStrategies = append(ret.BumpStrategies, s)
	}
	return &ret
}

// BumpOptions type to represent the available options for the bump commands
// It extends GlobalOptions.
type bumpOptions struct {
	*globalOptions
	viperConfig config
	// Bump is mapped to pkg/version/BumpStrategyOptions#Strategy
	Bump string
	// PreRelease is mapped to pkg/version/BumpStrategyOptions#PreRelease
	// It is set to true only if explicitly set by the user
	PreRelease bool
	// PreReleaseTemplate is mapped to pkg/version/BumpStrategyOptions#PreReleaseTemplate
	PreReleaseTemplate string
	// PreReleaseOverwrite is mapped to pkg/version/BumpStrategyOptions#PreReleaseOverwrite
	PreReleaseOverwrite bool
	// BuildMetadataTemplate is mapped to pkg/version/BumpStrategyOptions#BuildMetadataTemplate
	BuildMetadataTemplate string
	// BranchStrategies is mapped to pkg/version/BumpStrategyOptions#BranchStrategies
	BranchStrategies []string
}

func (o *bumpOptions) addBumpFlags(cmd *cobra.Command) {
	cmd.Flags().String("major-pattern", "", "Use major-pattern option to define your regular expression to match a breaking change commit message")
	cmd.Flags().String("minor-pattern", "", "Use major-pattern option to define your regular expression to match a minor change commit message")
	cmd.Flags().StringVar(&o.PreReleaseTemplate, "pre-release", "", preReleaseTemplateDesc)
	cmd.Flags().BoolVar(&o.PreReleaseOverwrite, "pre-release-overwrite", false, "Use pre-release overwrite option to remove the pre-release identifier suffix which will give a version like `X.Y.Z-SNAPSHOT` if pre-release=SNAPSHOT")
	cmd.Flags().StringVar(&o.BuildMetadataTemplate, "build-metadata", "", buildMetadataTemplateDesc)
	cmd.Flags().StringArrayVar(&o.BranchStrategies, "branch-strategy", []string{}, branchStrategyDesc)

	viper.BindPFlag("majorPattern", cmd.Flags().Lookup("major-pattern"))
	viper.BindPFlag("minorPattern", cmd.Flags().Lookup("minor-pattern"))

	viper.SetDefault("majorPattern", version.DefaultMajorPattern)
	viper.SetDefault("minorPattern", version.DefaultMinorPattern)
	viper.SetDefault("bumpStrategies", []interface{}{
		map[string]interface{}{
			"strategy":        "AUTO",
			"branchesPattern": version.DefaultReleaseBranchesPattern,
		},
		map[string]interface{}{
			"strategy":              "AUTO",
			"branchesPattern":       ".*",
			"buildMetadataTemplate": version.DefaultBuildMetadataTemplate,
		},
	})

	o.Cmd = cmd
}

func (o *bumpOptions) createBumpStrategy() *version.BumpStrategy {
	// ret := *version.NewBumpStrategy(git.NewVersionGitRepo(o.CurrentDir))
	viper.Unmarshal(&o.viperConfig)
	ret := o.viperConfig.createBumpStrategy()
	ret.SetGitRepository(git.NewVersionGitRepo(o.CurrentDir))

	for id, s := range o.BranchStrategies {
		if id == 0 {
			// reset branch strategy
			ret.BumpStrategies = []version.BumpBranchesStrategy{}
		}
		var b version.BumpBranchesStrategy
		json.Unmarshal([]byte(s), &b)
		ret.BumpStrategies = append(ret.BumpStrategies, b)
	}

	if o.Bump != "AUTO" || o.PreRelease || o.PreReleaseTemplate != "" || o.PreReleaseOverwrite || o.BuildMetadataTemplate != "" {
		// configure default BumpBranchesStrategy
		defaultStrategy := *version.NewBumpAllBranchesStrategy(version.ParseBumpStrategyType(o.Bump), o.PreRelease, o.PreReleaseTemplate, o.PreReleaseOverwrite, o.BuildMetadataTemplate)
		var bumpStrategies []version.BumpBranchesStrategy
		ret.BumpStrategies = append(bumpStrategies, defaultStrategy)
		//TODO: ret.BumpStrategies = append(bumpStrategies, defaultStrategy, ret.BumpStrategies...)
	}

	return ret
}

func run(o *bumpOptions) error {
	log.Debug("Run bump command with configuration: %#v", o)

	version, err := o.createBumpStrategy().Bump()
	if err != nil {
		return err
	}
	fmt.Fprintf(o.globalOptions.ioStreams.Out, "%v", version)
	return nil
}
