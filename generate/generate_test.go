package generate_test

import (
	"testing"

	"github.com/clarafu/release-me/generate"
	"github.com/clarafu/release-me/generate/mocks"
	"github.com/clarafu/release-me/github"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestGenerate(t *testing.T) {
	suite.Run(t, &GenerateSuite{
		Assertions: require.New(t),
	})
}

type GenerateSuite struct {
	suite.Suite
	*require.Assertions
}

type GenerateTest struct {
	It string

	PRs []github.PullRequest

	ExpectedBreaking []generate.PullRequest
	ExpectedFeatures []generate.PullRequest
	ExpectedBugFixes []generate.PullRequest
	ExpectedMisc     []generate.PullRequest

	GenerateErr error
}

func (s *GenerateSuite) TestGenerate() {
	for _, t := range []GenerateTest{
		{
			It: "groups the PRs by label",

			PRs: []github.PullRequest{
				{
					Title:  "new breaking change!",
					Labels: []string{"breaking"},
				},
				{
					Title:  "cool new feature!",
					Labels: []string{"enhancement"},
				},
				{
					Title:  "squash that bug!",
					Labels: []string{"bug"},
				},
				{
					Title:  "don't worry about it!",
					Labels: []string{"misc"},
				},
			},

			ExpectedBreaking: []generate.PullRequest{{Title: "new breaking change!"}},
			ExpectedFeatures: []generate.PullRequest{{Title: "cool new feature!"}},
			ExpectedBugFixes: []generate.PullRequest{{Title: "squash that bug!"}},
			ExpectedMisc:     []generate.PullRequest{{Title: "don't worry about it!"}},
		},
		{
			It: "sorts PRs by number",

			PRs: []github.PullRequest{
				{
					Number: 1,
					Labels: []string{"enhancement"},
				},
				{
					Number: 3,
					Labels: []string{"enhancement"},
				},
				{
					Number: 2,
					Labels: []string{"enhancement"},
				},
			},

			ExpectedFeatures: []generate.PullRequest{{Number: 1}, {Number: 2}, {Number: 3}},
		},
		{
			It: "sorts PRs with priority label first",

			PRs: []github.PullRequest{
				{
					Number: 1,
					Labels: []string{"enhancement"},
				},
				{
					Number: 3,
					Labels: []string{"enhancement", "priority"},
				},
				{
					Number: 2,
					Labels: []string{"enhancement"},
				},
			},

			ExpectedFeatures: []generate.PullRequest{{Number: 3}, {Number: 1}, {Number: 2}},
		},
		{
			It: "groups PRs as breaking first",

			PRs: []github.PullRequest{
				{
					Title:  "new breaking change!",
					Labels: []string{"enhancement", "breaking", "misc", "bug"},
				},
			},

			ExpectedBreaking: []generate.PullRequest{{Title: "new breaking change!"}},
		},
		{
			It: "groups PRs as misc before bugs and features",

			PRs: []github.PullRequest{
				{
					Title:  "super fun pull request",
					Labels: []string{"enhancement", "misc", "bug"},
				},
			},

			ExpectedMisc: []generate.PullRequest{{Title: "super fun pull request"}},
		},
		{
			It: "groups PRs as misc before features",

			PRs: []github.PullRequest{
				{
					Title:  "best feature ever",
					Labels: []string{"enhancement", "misc"},
				},
			},

			ExpectedMisc: []generate.PullRequest{{Title: "best feature ever"}},
		},
		{
			It: "fails when PR does not have appropriate label",

			PRs: []github.PullRequest{
				{
					Number: 1,
					Url:    "http://pr/1",
					Title:  "no labels",
				},
				{
					Number: 2,
					Url:    "http://pr/2",
					Title:  "correctly labelled PR",
					Labels: []string{"enhancement"},
				},
				{
					Number: 3,
					Url:    "http://pr/3",
					Title:  "other no labels",
				},
			},

			GenerateErr: generate.PullRequestsNotLabelled{
				Identifiers: []string{
					"http://pr/1",
					"http://pr/3",
				},
			},
		},
		{
			It: "parses pull request release note description from header Release Note",

			PRs: []github.PullRequest{
				{
					Title:  "Fist of the North Star",
					Labels: []string{"enhancement"},
					Body: `# Description
blah

## Release Note

omai wa mo shindeiru

## End of Pull request

its over`,
				},
			},

			ExpectedFeatures: []generate.PullRequest{
				{
					Title:       "Fist of the North Star",
					ReleaseNote: "omai wa mo shindeiru",
				},
			},
		},
		{
			It: "parses pull request release note description from header Release Notes",

			PRs: []github.PullRequest{
				{
					Title:  "Fist of the North Star",
					Labels: []string{"enhancement"},
					Body: `## Release Notes

omai wa mo shindeiru`,
				},
			},

			ExpectedFeatures: []generate.PullRequest{
				{
					Title:       "Fist of the North Star",
					ReleaseNote: "omai wa mo shindeiru",
				},
			},
		},
		{
			It: "parses pull request description from header case insensitive",

			PRs: []github.PullRequest{
				{
					Title:  "Fist of the North Star",
					Labels: []string{"enhancement"},
					Body: `## release note

omai wa mo shindeiru`,
				},
			},

			ExpectedFeatures: []generate.PullRequest{
				{
					Title:       "Fist of the North Star",
					ReleaseNote: "omai wa mo shindeiru",
				},
			},
		},
	} {
		s.Run(t.It, func() {
			fakeTemplate := new(mocks.Template)
			fakeTemplate.On("Render", mock.Anything).Return(nil)

			generator := generate.New(fakeTemplate)

			err := generator.Generate(t.PRs)
			if t.GenerateErr != nil {
				s.Equal(err.Error(), t.GenerateErr.Error())
			} else {
				s.NoError(err)

				fakeTemplate.AssertCalled(s.T(), "Render", []generate.Section{
					generate.Section{Title: "Breaking", Icon: "🚨", PRs: t.ExpectedBreaking},
					generate.Section{Title: "Features", Icon: "✈️", PRs: t.ExpectedFeatures},
					generate.Section{Title: "Bug Fixes", Icon: "🐞", PRs: t.ExpectedBugFixes},
					generate.Section{Title: "Miscellaneous", Icon: "🤷", PRs: t.ExpectedMisc},
				})
			}
		})
	}
}
