package utils

import (
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

// NewTemplate creates a new Template with sprig functions, returning an error if parsing fails.
func NewTemplate(value string) (*template.Template, error) {
	if value == "" {
		return nil, nil
	}
	return template.New("").Funcs(sprig.TxtFuncMap()).Parse(value)
}

// MustNewTemplate creates a new Template with sprig functions, panicking if parsing fails.
func MustNewTemplate(value string) *template.Template {
	t, err := NewTemplate(value)
	if err != nil {
		panic(err)
	}
	return t
}

// TemplateToString converts *template.Template instance to string
func TemplateToString(t *template.Template) string {
	if t != nil {
		return t.Root.String()
	}
	return ""
}
