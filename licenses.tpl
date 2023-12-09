{{ range . }}
## {{ .Name }}

- Version: {{ .Version }}
- License: [{{ .LicenseName }}]({{ .LicenseURL }})

{{ end }}