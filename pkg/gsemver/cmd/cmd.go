package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"

	log "github.com/arnaud-deprez/gsemver/internal/log"
)

const (
	optionVerbose  = "verbose"
	optionLogLevel = "log-level"
)

// IOStreams provides the standard names for iostreams.  This is useful for embedding and for unit testing.
// Inconsistent and different names make it hard to read and review code
type IOStreams struct {
	// In think, os.Stdin
	In io.Reader
	// Out think, os.Stdout
	Out io.Writer
	// ErrOut think, os.Stderr
	ErrOut io.Writer
}

// NewIOStreams creates a IOStreams
func NewIOStreams(in io.Reader, out, err io.Writer) *IOStreams {
	return &IOStreams{
		In:     in,
		Out:    out,
		ErrOut: err,
	}
}

// GlobalOptions provides the global options of the CLI
type GlobalOptions struct {
	Cmd        *cobra.Command
	Args       []string
	CurrentDir string
	Verbose    bool
	LogLevel   string
	ioStreams  *IOStreams
}

// addGlobalFlags adds the common flags to the given command
func (options *GlobalOptions) addGlobalFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&options.Verbose, optionVerbose, "", false, "Enables verbose output")
	cmd.PersistentFlags().StringVarP(&options.LogLevel, optionLogLevel, "", "INFO", "Sets the logging level (panic, fatal, error, warning, info, debug)")

	dir, err := os.Getwd()
	if err != nil {
		log.Error("Unable to retrieve working directory: %v", err)
		os.Exit(1)
	}

	options.CurrentDir = dir
	options.Cmd = cmd
}

var globalUsage = `Simple CLI to manage semver compliant version from your git tags
`

// NewDefaultRootCommand creates the `gsemver` command with default arguments
func NewDefaultRootCommand() *cobra.Command {
	return NewRootCommand(os.Stdin, os.Stdout, os.Stderr)
}

// NewRootCommand creates the `gsemver` command with args
func NewRootCommand(in io.Reader, out, errout io.Writer) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "gsemver",
		Short: "CLI to manage semver compliant version from your git tags",
		Long:  globalUsage,
		Run:   runHelp,
	}
	// commonOpts holds the global flags that will be shared/inherited by all sub-commands created bellow
	globalOpts := &GlobalOptions{ioStreams: NewIOStreams(in, out, errout)}
	globalOpts.addGlobalFlags(cmds)

	cmds.AddCommand(
		NewBumpCommands(globalOpts),
		// Hidden documentation generator command: 'helm docs'
		newDocsCmd(globalOpts),
	)
	return cmds
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
