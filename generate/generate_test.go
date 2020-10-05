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

	ExpectedBreaking   []generate.PullRequest
	ExpectedFeatures   []generate.PullRequest
	ExpectedBugFixes   []generate.PullRequest
	ExpectedRefactor   []generate.PullRequest
	ExpectedTesting    []generate.PullRequest
	ExpectedDependency []generate.PullRequest
	ExpectedInternal   []generate.PullRequest
	ExpectedNoImpact   []generate.PullRequest

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
					Title:  "refactor don't worry about it!",
					Labels: []string{"refactor"},
				},
				{
					Title:  "testing don't worry about it!",
					Labels: []string{"testing"},
				},
				{
					Title:  "dependencies don't worry about it!",
					Labels: []string{"dependencies"},
				},
				{
					Title:  "internal change don't worry about it!",
					Labels: []string{"internal"},
				},
				{
					Title:  "not a change don't worry about it!",
					Labels: []string{"release/no-impact"},
				},
			},

			ExpectedBreaking:   []generate.PullRequest{{Title: "new breaking change!"}},
			ExpectedFeatures:   []generate.PullRequest{{Title: "cool new feature!"}},
			ExpectedBugFixes:   []generate.PullRequest{{Title: "squash that bug!"}},
			ExpectedRefactor:   []generate.PullRequest{{Title: "refactor don't worry about it!"}},
			ExpectedTesting:    []generate.PullRequest{{Title: "testing don't worry about it!"}},
			ExpectedDependency: []generate.PullRequest{{Title: "dependencies don't worry about it!"}},
			ExpectedInternal:   []generate.PullRequest{{Title: "internal change don't worry about it!"}},
			ExpectedNoImpact:   []generate.PullRequest{{Title: "not a change don't worry about it!"}},
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
					Labels: []string{"enhancement", "breaking", "refactor", "bug"},
				},
			},

			ExpectedBreaking: []generate.PullRequest{{Title: "new breaking change!"}},
		},
		{
			It: "groups PRs as bugs before features and no impact",

			PRs: []github.PullRequest{
				{
					Title:  "super fun pull request",
					Labels: []string{"enhancement", "refactor", "bug"},
				},
			},

			ExpectedBugFixes: []generate.PullRequest{{Title: "super fun pull request"}},
		},
		{
			It: "groups PRs as features before no impact",

			PRs: []github.PullRequest{
				{
					Title:  "best feature ever",
					Labels: []string{"enhancement", "refactor"},
				},
			},

			ExpectedFeatures: []generate.PullRequest{{Title: "best feature ever"}},
		},
		{
			It: "groups PRs as refactors before no impact",

			PRs: []github.PullRequest{
				{
					Title:  "best refactor ever",
					Labels: []string{"release/no-impact", "refactor"},
				},
			},

			ExpectedRefactor: []generate.PullRequest{{Title: "best refactor ever"}},
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
					generate.Section{
						Title: "No Impact", Icon: "🤷",
						PRs: t.ExpectedNoImpact,
						SubSections: []generate.SubSection{
							generate.SubSection{Title: "Refactors", PRs: t.ExpectedRefactor},
							generate.SubSection{Title: "Tests", PRs: t.ExpectedTesting},
							generate.SubSection{Title: "Dependencies", PRs: t.ExpectedDependency},
							generate.SubSection{Title: "Internal Changes", PRs: t.ExpectedInternal},
						},
					},
				})
			}
		})
	}
}
