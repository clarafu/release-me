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

func (s *TemplateSuite) TestRemovesEmptySubSections() {
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
			Title: "Section with empty sub sections",
			Icon:  "🐞",
			SubSections: []generate.SubSection{
				generate.SubSection{
					Title: "SubSection 1",
				},
			},
		},
	}
	buf := new(bytes.Buffer)
	generate.NewReleaseNoteTemplater(buf).Render(sections)
	s.NotContains(buf.String(), "SubSection 1")
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
		generate.Section{
			Title: "Section 3",
			Icon:  "🐞",
			SubSections: []generate.SubSection{
				generate.SubSection{
					Title: "SubSection 1",
					PRs: []generate.PullRequest{
						generate.PullRequest{
							Title:  "PR Title 3",
							Number: 2,
							Author: "Michael Kerrisk",
						},
					},
				},
			},
		},
	}
	buf := new(bytes.Buffer)
	generate.NewReleaseNoteTemplater(buf).Render(sections)
	s.Contains(buf.String(), "Section 1")
	s.Contains(buf.String(), "Section 2")
	s.Contains(buf.String(), "Section 3")
	s.Contains(buf.String(), "SubSection 1")
}
