package gsemver

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

const (
	docsDesc = `
Generate documentation files for gsemver.
This command can generate documentation for gsemver in the following formats:
- Markdown
- Man pages
It can also generate bash autocompletions.	
`
	docsExample = `
# For Markdown documentation: 
gsemver docs markdown -dir docs/

# For Man page format:
gsemver docs man -dir docs/

# For bash completion:
gsemver docs bash -dir docs/
`
)

type docsOptions struct {
	*globalOptions
	dest          string
	docTypeString string
	topCmd        *cobra.Command
}

func newDocsCommands(globalOpts *globalOptions) *cobra.Command {
	o := &docsOptions{
		globalOptions: globalOpts,
	}

	cmd := &cobra.Command{
		Use:       "docs",
		Short:     "Generate documentation as markdown or man pages",
		Long:      docsDesc,
		Example:   docsExample,
		Hidden:    true,
		ValidArgs: []string{"markdown", "man", "bash"},
		Args:      cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.topCmd = cmd.Root()
			if len(args) == 0 {
				o.docTypeString = "mardown"
			} else {
				o.docTypeString = args[0]
			}
			return o.run()
		},
	}

	f := cmd.Flags()
	f.StringVar(&o.dest, "dir", "./", "directory to which documentation is written")

	return cmd
}

func (o *docsOptions) run() error {
	o.topCmd.DisableAutoGenTag = true
	switch o.docTypeString {
	case "markdown", "mdown", "md":
		return doc.GenMarkdownTree(o.topCmd, o.dest)
	case "man":
		manHdr := &doc.GenManHeader{Title: "gsemver", Section: "1"}
		return doc.GenManTree(o.topCmd, manHdr, o.dest)
	case "bash":
		return o.topCmd.GenBashCompletionFile(filepath.Join(o.dest, "completions.bash"))
	default:
		return errors.Errorf("unknown doc type %q. Try 'markdown' or 'man'", o.docTypeString)
	}
}
