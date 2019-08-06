package version

import (
	"strings"
	"text/template"

	"github.com/arnaud-deprez/gsemver/internal/log"
	"github.com/arnaud-deprez/gsemver/pkg/git"
)

// NewContext returns a new Context
func NewContext(branch string, lastVersion *Version, lastTag *git.Tag, commits []git.Commit) *Context {
	return &Context{
		Branch:      branch,
		LastVersion: lastVersion,
		LastTag:     lastTag,
		Commits:     commits,
	}
}

// Context represents the context data used to compute the next version.
// This context is also used as template data.
type Context struct {
	// Branch is the current branch name
	Branch string
	// LastVersion is a semver version representation of the last git tag
	LastVersion *Version
	// LastTag is the last git tag
	LastTag *git.Tag
	// Commits is the list of commits from the previous tag until now
	Commits []git.Commit
}

// EvalTemplate evaluates the given template against the current context
func (c *Context) EvalTemplate(template *template.Template) string {
	if template == nil || c == nil {
		return ""
	}
	var sb strings.Builder
	err := template.Execute(&sb, c)
	if err != nil {
		// Stop the program
		log.Fatal("TemplateContext: fails to evaluate template caused by %v", err)
	}

	return sb.String()
}
