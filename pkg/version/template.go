package version

import (
	"text/template"

	"github.com/Masterminds/sprig"
)

// NewTemplate create a new Template with sprig functions
func NewTemplate(value string) *template.Template {
	if value == "" {
		return nil
	}
	return template.Must(template.New("").Funcs(sprig.TxtFuncMap()).Parse(value))
}
