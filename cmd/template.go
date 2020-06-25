package cmd

import (
	"io"
	"text/template"

	"github.com/clarafu/release-me/github"
)

type Section struct {
	Title string
	Icon  string
	PRs   []github.PullRequest
}

const rawTemplate = `
{{range $section := .Sections}}
{{if $section.PRs }}
<br />

# {{$section.Icon}} {{$section.Title}}

{{ range $pr := $section.PRs }}
* {{$pr.Title}} (#{{$pr.Number}}) @{{$pr.Author}} <sub><sup><a name="{{$pr.Number}}" href="#{{$pr.Number}}">:link:</a></sup></sub>
{{end}}
{{end}}
{{end}}
`

var releaseNotesTemplate = template.Must(template.New("release_notes").Parse(rawTemplate))

func WriteReleaseNotes(w io.Writer, sections []Section) error {
	return releaseNotesTemplate.Execute(w, struct{
		Sections []Section
	}{
		Sections: sections,
	})
}
