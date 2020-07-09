package generate

import (
	"fmt"
	"sort"
	"strings"

	"github.com/clarafu/release-me/github"
)

type PullRequestsNotLabelled struct {
	Urls []string
}

func (e PullRequestsNotLabelled) Error() string {
	prUrls := []string{}
	for _, url := range e.Urls {
		prUrls = append(prUrls, "- "+url)
	}
	return fmt.Sprintf(`

The following pull requests: 
%s
	
must be labelled with at least one of:
- breaking
- enhancement
- bug
- release/no-impact`, strings.Join(prUrls, "\n"))
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

	var breakingPRs, noImpactPRs, featurePRs, bugFixPRs []PullRequest
	var unlabelledPRUrls []string
	for _, githubPR := range prs {
		pr := PullRequest{
			Title:       githubPR.Title,
			Author:      githubPR.Author,
			Number:      githubPR.Number,
			ReleaseNote: parseReleaseNote(githubPR.Body),
		}

		if githubPR.HasLabel("breaking") {
			breakingPRs = append(breakingPRs, pr)
			continue
		}

		if githubPR.HasLabel("release/no-impact") {
			noImpactPRs = append(noImpactPRs, pr)
			continue
		}

		if githubPR.HasLabel("enhancement") {
			featurePRs = append(featurePRs, pr)
			continue
		}

		if githubPR.HasLabel("bug") {
			bugFixPRs = append(bugFixPRs, pr)
			continue
		}

		unlabelledPRUrls = append(unlabelledPRUrls, githubPR.Url)
	}

	if len(unlabelledPRUrls) > 0 {
		return PullRequestsNotLabelled{Urls: unlabelledPRUrls}
	}

	sections := []Section{
		Section{Title: "Breaking", Icon: "🚨", PRs: breakingPRs},
		Section{Title: "Features", Icon: "✈️", PRs: featurePRs},
		Section{Title: "Bug Fixes", Icon: "🐞", PRs: bugFixPRs},
		Section{Title: "No Impact", Icon: "🤷", PRs: noImpactPRs},
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
