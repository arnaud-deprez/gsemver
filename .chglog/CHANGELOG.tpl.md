{{- range .Versions }}
{{- if .Tag.Previous }}
## [{{ .Tag.Name }}]({{ $.Info.RepositoryURL }}/compare/{{ .Tag.Previous.Name }}...{{ .Tag.Name }})
{{- else }}
## [{{ .Tag.Name }}]({{ $.Info.RepositoryURL }}/releases/tag/{{ .Tag.Name }})
{{- end }} ({{ datetime "2006-01-02" .Tag.Date }})
{{ range .CommitGroups }}
### {{ .Title }}
{{ range .Commits }}
- {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }} ([{{ .Hash.Short }}]({{ $.Info.RepositoryURL }}/commits/{{ .Hash.Long }})
{{- $lenRefs := len .Refs -}}
{{- if gt $lenRefs 0 }}, closes 
{{- range $idx, $ref := .Refs }}
{{- if $idx }},{{ end }} [#{{ $ref.Ref }}]({{ $.Info.RepositoryURL }}/issues/{{ $ref.Ref }})
{{- end }}
{{- end }})
{{- end }}
{{ end -}}

{{- if .RevertCommits }}
### Reverts
{{ range .RevertCommits }}
- {{ .Revert.Header }}
{{- end }}
{{ end -}}

{{- if .NoteGroups }}
{{ range .NoteGroups -}}
### {{ .Title }}
{{ range .Notes }}
- {{ .Body }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}