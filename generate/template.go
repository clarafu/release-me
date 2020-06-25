package generate

import (
	"io"
	"strings"
	"text/template"
)

type PullRequest struct {
	Title       string
	Number      int
	Author      string
	ReleaseNote string
}

type Section struct {
	Title string
	Icon  string
	PRs   []PullRequest
}

func indent(spaces int, v string) string {
	pad := strings.Repeat(" ", spaces)
	return pad + strings.Replace(v, "\n", "\n"+pad, -1)
}

const rawTemplate = `
{{range $section := .Sections}}
{{if $section.PRs }}
<br />

# {{$section.Icon}} {{$section.Title}}

{{ range $pr := $section.PRs }}
* **{{$pr.Title}} (#{{$pr.Number}})** @{{$pr.Author}} <sub><sup><a name="{{$pr.Number}}" href="#{{$pr.Number}}">:link:</a></sup></sub>  
{{if $pr.ReleaseNote}}{{ $pr.ReleaseNote | indent 2 }}{{end}}
{{end}}
{{end}}
{{end}}
`

var funcMap = template.FuncMap{
	"indent": indent,
}

var releaseNotesTemplate = template.Must(template.New("release_notes").Funcs(funcMap).Parse(rawTemplate))

func writeReleaseNotes(w io.Writer, sections []Section) error {
	return releaseNotesTemplate.Execute(w, struct {
		Sections []Section
	}{
		Sections: sections,
	})
}
