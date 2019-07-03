package gsemver

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/arnaud-deprez/gsemver/internal/git"
	"github.com/arnaud-deprez/gsemver/pkg/version"
)

// NewBumpCommands create the bump command with its subcommands
func newBumpCommands(globalOpts *globalOptions) *cobra.Command {
	options := &bumpOptions{
		globalOptions: globalOpts,
	}

	cmd := &cobra.Command{
		Use:       "bump [strategy]",
		Short:     "Bump to next version",
		Long:      "",
		ValidArgs: []string{"auto", "major", "minor", "patch"},
		Args:      cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
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
	// PreRelease is mapped to pkg/version/BumpStrategyOptions#PreRelease
	PreRelease string
	// PreReleaseOverwrite is mapped to pkg/version/BumpStrategyOptions#PreReleaseOverwrite
	PreReleaseOverwrite bool
	// BuildMetadata is mapped to pkg/version/BumpStrategyOptions#BuildMetadata
	BuildMetadata string
}

func (o *bumpOptions) addBumpFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&o.PreRelease, "pre-release", "", "", "Use pre-release version such as alpha which will give a version like X.Y.Z-alpha.N")
	cmd.PersistentFlags().BoolVarP(&o.PreReleaseOverwrite, "pre-release-overwrite", "", false, "Use pre-release overwrite option to remove the pre-release identifier suffix which will give a version like X.Y.Z-SNAPSHOT if pre-release=SNAPSHOT")
	cmd.PersistentFlags().StringVarP(&o.BuildMetadata, "build", "", "", "Use build metadata which will give something like X.Y.Z+<build>")

	o.Cmd = cmd
}

func (o *bumpOptions) createBumpStrategy() *version.BumpStrategyOptions {
	ret := version.NewConventionalCommitBumpStrategyOptions(git.NewVersionGitRepo(o.CurrentDir))
	ret.Strategy = version.ParseBumpStrategy(o.Bump)
	ret.PreRelease = o.PreRelease
	ret.PreReleaseOverwrite = o.PreReleaseOverwrite
	ret.BuildMetadata = o.BuildMetadata
	return ret
}

func (o *bumpOptions) run() error {
	version, err := o.createBumpStrategy().Bump()
	if err != nil {
		return err
	}
	fmt.Fprintf(o.globalOptions.ioStreams.Out, "%v", version)
	return nil
}