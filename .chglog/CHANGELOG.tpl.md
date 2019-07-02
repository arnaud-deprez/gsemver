{{- $root := . -}}

{{- range .Versions }}
## [{{ .Tag.Name }}]({{ $root.Info.RepositoryURL }}/releases/tag/{{ .Tag.Name }}) ({{ datetime "2006-01-02" .Tag.Date }})
{{ range .CommitGroups }}
### {{ .Title }}
{{ range .Commits }}
- {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }} ([{{ .Hash.Short }}]({{ $root.Info.RepositoryURL }}/commits/{{ .Hash.Long }})
{{- $lenRefs := len .Refs -}}
{{- if gt $lenRefs 0 }}, closes 
{{- range $idx, $ref := .Refs }}
{{if $idx }}, {{end}}[#{{ $ref.Ref }}]({{ $root.Info.RepositoryURL }}/issues/{{ $ref.Ref }})
{{- end }}
{{- end }})
{{- end }}
{{ end -}}

{{- if .RevertCommits -}}
### Reverts
{{ range .RevertCommits }}
- {{ .Revert.Header }}
{{ end }}
{{ end -}}

{{- if .NoteGroups -}}
{{ range .NoteGroups -}}
### {{ .Title }}
{{ range .Notes }}
{{ .Body }}
{{ end }}
{{ end -}}
{{ end -}}
{{ end -}}

{{- if .Versions }}
[Unreleased]: {{ .Info.RepositoryURL }}/compare/{{ $latest := index .Versions 0 }}{{ $latest.Tag.Name }}...HEAD
{{ range .Versions -}}
{{ if .Tag.Previous -}}
[{{ .Tag.Name }}]: {{ $.Info.RepositoryURL }}/compare/{{ .Tag.Previous.Name }}...{{ .Tag.Name }}
{{ end -}}
{{ end -}}
{{ end -}}