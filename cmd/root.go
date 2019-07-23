package cmd

import (
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	log "github.com/arnaud-deprez/gsemver/internal/log"
)

const (
	optionVerbose  = "verbose"
	optionLogLevel = "log-level"
)

// ioStreams provides the standard names for iostreams.  This is useful for embedding and for unit testing.
// Inconsistent and different names make it hard to read and review code
type ioStreams struct {
	// In think, os.Stdin
	In io.Reader
	// Out think, os.Stdout
	Out io.Writer
	// ErrOut think, os.Stderr
	ErrOut io.Writer
}

// newIOStreams creates a IOStreams
func newIOStreams(in io.Reader, out, err io.Writer) *ioStreams {
	return &ioStreams{
		In:     in,
		Out:    out,
		ErrOut: err,
	}
}

// globalOptions provides the global options of the CLI
type globalOptions struct {
	// Cmd is the current *cobra.Command
	Cmd *cobra.Command
	// Args contains all the non options args for the command
	Args []string
	// CurrentDir is the directory from where the command has been executed.
	CurrentDir string
	// Verbose enables verbose output
	Verbose bool
	// LogLevel sets the log level (panic, fatal, error, warning, info, debug)
	LogLevel string
	// ioStreams contains the input, output and error stream
	ioStreams *ioStreams
}

var globalUsage = `Simple CLI to manage semver compliant version from your git tags
`

// newDefaultRootCommand creates the `gsemver` command with default arguments
func newDefaultRootCommand() *cobra.Command {
	return newRootCommand(os.Stdin, os.Stdout, os.Stderr)
}

// newRootCommand creates the `gsemver` command with args
func newRootCommand(in io.Reader, out, errout io.Writer) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "gsemver",
		Short: "CLI to manage semver compliant version from your git tags",
		Long:  globalUsage,
		Run:   runHelp,
	}
	// commonOpts holds the global flags that will be shared/inherited by all sub-commands created bellow
	globalOpts := &globalOptions{ioStreams: newIOStreams(in, out, errout)}
	globalOpts.addGlobalFlags(cmds)

	cmds.AddCommand(
		newBumpCommands(globalOpts),
		newVersionCommands(globalOpts),
		// Hidden documentation generator command: 'helm docs'
		newDocsCommands(globalOpts),
	)
	return cmds
}

// addGlobalFlags adds the common flags to the given command
func (o *globalOptions) addGlobalFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&o.Verbose, optionVerbose, "v", false, "Enables verbose output by setting log level to debug. This is a shortland to --log-level debug.")
	cmd.PersistentFlags().StringVarP(&o.LogLevel, optionLogLevel, "", "info", "Sets the logging level (fatal, error, warning, info, debug, trace)")

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable to retrieve working directory: %v", err)
	}

	o.CurrentDir = dir
	o.Cmd = cmd
}

func (o *globalOptions) configureLogger() {
	if o.Verbose && strings.ToLower(o.LogLevel) != "trace" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevelS(o.LogLevel)
	}
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
