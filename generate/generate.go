package generate

import (
	"fmt"
	"sort"
	"strings"

	"github.com/clarafu/release-me/github"
)

// The ordering of this list is the order of precedence of the labels
var ValidLabels = []string{
	"breaking",
	"bug",
	"enhancement",
	"refactor",
	"testing",
	"dependencies",
	"internal",
	"release/no-impact",
}

type PullRequestsNotLabelled struct {
	Identifiers []string
}

func (e PullRequestsNotLabelled) Error() string {
	prIdentifiers := []string{}
	for _, identifier := range e.Identifiers {
		prIdentifiers = append(prIdentifiers, "- "+identifier)
	}

	validLabels := []string{}
	for _, label := range ValidLabels {
		validLabels = append(validLabels, "- "+label)
	}

	return fmt.Sprintf(`

The following pull request(s): 
%s
	
must be labelled with at least one of:
%s`, strings.Join(prIdentifiers, "\n"), strings.Join(validLabels, "\n"))
}

type Template interface {
	Render(sections []Section) error
}

type Generator struct {
	template Template
}

func New(template Template) Generator {
	return Generator{template}
}

func (g Generator) Generate(prs []github.PullRequest) error {
	g.sortPRsByPriority(prs)

	var unlabelledPRUrls []string
	sectionPRs := make(map[string][]PullRequest)
	for _, githubPR := range prs {
		pr := PullRequest{
			Title:       githubPR.Title,
			Author:      githubPR.Author,
			Number:      githubPR.Number,
			ReleaseNote: parseReleaseNote(githubPR.Body),
		}

		var labelled bool
		for _, label := range ValidLabels {
			if githubPR.HasLabel(label) {
				sectionPRs[label] = append(sectionPRs[label], pr)
				labelled = true
				break
			}
		}

		if !labelled {
			unlabelledPRUrls = append(unlabelledPRUrls, githubPR.Url)
		}
	}

	if len(unlabelledPRUrls) > 0 {
		return PullRequestsNotLabelled{Identifiers: unlabelledPRUrls}
	}

	sections := []Section{
		Section{Title: "Breaking", Icon: "🚨", PRs: sectionPRs["breaking"]},
		Section{Title: "Features", Icon: "✈️", PRs: sectionPRs["enhancement"]},
		Section{Title: "Bug Fixes", Icon: "🐞", PRs: sectionPRs["bug"]},
		Section{Title: "No Impact", Icon: "🤷", PRs: sectionPRs["release/no-impact"],
			SubSections: []SubSection{
				{Title: "Refactors", PRs: sectionPRs["refactor"]},
				{Title: "Tests", PRs: sectionPRs["testing"]},
				{Title: "Dependencies", PRs: sectionPRs["dependencies"]},
				{Title: "Internal Changes", PRs: sectionPRs["internal"]},
			},
		},
	}

	err := g.template.Render(sections)
	if err != nil {
		return fmt.Errorf("failed to write release notes: %w", err)
	}

	return nil
}

func (g Generator) sortPRsByPriority(prs []github.PullRequest) {
	sort.Slice(prs, func(i, j int) bool {
		switch prs[i].HasLabel("priority") != prs[j].HasLabel("priority") {
		case true:
			if prs[i].HasLabel("priority") {
				return true
			}
			return false

		case false:
			// If both prs have the same priority (both have the priority label
			// or both do not), order by the pr number
			return prs[i].Number < prs[j].Number
		default:
			panic("this should never happen!")
		}
	})
}

func Validate(labels []string) bool {
	validLabelsMap := make(map[string]bool)
	for _, validLabel := range ValidLabels {
		validLabelsMap[validLabel] = true
	}

	for _, label := range labels {
		if _, exists := validLabelsMap[label]; exists {
			return true
		}
	}

	return false
}
