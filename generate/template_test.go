package generate_test

import (
	"bytes"
	"testing"

	"github.com/clarafu/release-me/generate"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestTemplate(t *testing.T) {
	suite.Run(t, &TemplateSuite{
		Assertions: require.New(t),
	})
}

type TemplateSuite struct {
	suite.Suite
	*require.Assertions
}

type TemplateTest struct {
	It string

	Sections       generate.Section
	ExpectedOutput string
}

func (s *TemplateSuite) TestRemovesEmptySections() {
	sections := []generate.Section{
		generate.Section{
			Title: "Section with PRs",
			Icon:  "🚨",
			PRs: []generate.PullRequest{
				generate.PullRequest{
					Title:  "PR Title",
					Number: 0,
					Author: "J.K. Rowling",
				},
			},
		},
		generate.Section{
			Title: "Section with no PRs",
			Icon:  "🐞",
		},
	}
	buf := new(bytes.Buffer)
	generate.NewReleaseNoteTemplater(buf).Render(sections)
	s.NotContains(buf.String(), "no PRs")
}

func (s *TemplateSuite) TestMultipleSections() {
	sections := []generate.Section{
		generate.Section{
			Title: "Section 1",
			Icon:  "🚨",
			PRs: []generate.PullRequest{
				generate.PullRequest{
					Title:  "PR Title",
					Number: 0,
					Author: "J.K. Rowling",
				},
			},
		},
		generate.Section{
			Title: "Section 2",
			Icon:  "🐞",
			PRs: []generate.PullRequest{
				generate.PullRequest{
					Title:  "PR Title 2",
					Number: 1,
					Author: "Steven King",
				},
			},
		},
	}
	buf := new(bytes.Buffer)
	generate.NewReleaseNoteTemplater(buf).Render(sections)
	s.Contains(buf.String(), "Section 1")
	s.Contains(buf.String(), "Section 2")
}
