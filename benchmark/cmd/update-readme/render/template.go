package render

import "text/template"

const (
	MarkerStart = "<!-- BENCHSTAT:START -->"
	MarkerEnd   = "<!-- BENCHSTAT:END -->"
)

type sectionTemplateData struct {
	MarkerStart string
	MarkerEnd   string
	UpdatedUTC  string
	GoVersion   string
	OSArch      string
	CPU         string
	Sections    []benchSection
}

var sectionTemplate = template.Must(
	template.New("benchmark-section").
		Parse(
			`{{.MarkerStart}}

## Latest Benchmark Comparison

- Updated (UTC): {{.UpdatedUTC}}
- Go: ` + "`{{.GoVersion}}`" + `
{{if .OSArch}}- OS/Arch: ` + "`{{.OSArch}}`" + `
{{end}}{{if .CPU}}- CPU: ` + "`{{.CPU}}`" + `
{{end}}
{{range .Sections}}### {{.Title}}

| Library | ops (N) | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
{{range .Rows}}| {{.Library}} | {{.Ops}} | {{.NsOp}} | {{.BOp}} | {{.Allocs}} |
{{end}}
{{end}}{{.MarkerEnd}}`,
		),
)
