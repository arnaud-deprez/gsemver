package utils

import (
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

// NewTemplate create a new Template with sprig functions
func NewTemplate(value string) *template.Template {
	if value == "" {
		return nil
	}
	return template.Must(template.New("").Funcs(sprig.TxtFuncMap()).Parse(value))
}

// TemplateToString	converts *template.Template instance to string
func TemplateToString(t *template.Template) string {
	if t != nil {
		return t.Root.String()
	}
	return ""
}
