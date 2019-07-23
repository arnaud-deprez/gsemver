package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/arnaud-deprez/gsemver/internal/log"
	"github.com/arnaud-deprez/gsemver/internal/version"
)

const (
	versionDesc = `
Show the version for gsemver.

This will print a representation the version of gsemver.
The output will look something like this:

version.BuildInfo{Version:"0.1.0", GitCommit:"acfe51b15f9a1f12d47a20f88c29e5364916ae57", GitTreeState:"clean", BuildDate:"2019-07-02T07:44:00Z", GoVersion:"go1.12.6", Compiler:"gc", Platform:"darwin/amd64"}

- Version is the semantic version of the release.
- GitCommit is the SHA for the commit that this version was built from.
- GitTreeState is "clean" if there are no local code changes when this binary was
  built, and "dirty" if the binary was built from locally modified code.
- BuildDate is the build date in ISO-8601 format at UTC.
- GoVersion is the go version with which it has been built.
- Compiler is the go compiler with which it has been built.
- Platform is the current OS platform on which it is running and for which it has been built.
`
	versionExample = `
# Print version of gsemver
$ gsemver version
`
)

// newVersionCommands create the version commands
func newVersionCommands(globalOpts *globalOptions) *cobra.Command {
	options := &versionOptions{
		globalOptions: globalOpts,
	}

	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Print the CLI version information",
		Long:    versionDesc,
		Example: versionExample,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return options.run()
		},
	}

	options.addVersionFlags(cmd)

	return cmd
}

type versionOptions struct {
	*globalOptions
	short bool
}

func (o *versionOptions) addVersionFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.short, "short", false, "print the version number")

	o.Cmd = cmd
}

func (o *versionOptions) run() error {
	log.Debug("Run version command with configuration: %#v", o)

	fmt.Fprintln(o.globalOptions.ioStreams.Out, formatVersion(o.short))
	return nil
}

func formatVersion(short bool) string {
	if short {
		return version.GetVersion()
	}
	return fmt.Sprintf("%#v", version.Get())
}
