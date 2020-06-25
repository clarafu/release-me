package generate

import (
	"io"
	"sort"
	"fmt"

	"github.com/clarafu/release-me/github"
)

type PullRequestNotLabelled struct {
	Number int
}

func (e PullRequestNotLabelled) Error() string {
	return fmt.Sprintf(`Pull request #%d must be labelled with at least one of:
	
- breaking
- enhancement
- bug
- release/no-impact`, e.Number)
}

func Generate(w io.Writer, prs []github.PullRequest) error {
	sort.Slice(prs, func(i, j int) bool {
		return prs[i].Number < prs[j].Number
	})

	var breakingPRs, noImpactPRs, featurePRs, bugFixPRs []PullRequest
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

		return PullRequestNotLabelled{Number: pr.Number}
	}

	sections := []Section{
		Section{Title: "Breaking", Icon: "🚨", PRs: breakingPRs},
		Section{Title: "Features", Icon: "✈️", PRs: featurePRs},
		Section{Title: "Bug Fixes", Icon: "🐞", PRs: bugFixPRs},
		Section{Title: "No Impact", Icon: "🤷", PRs: noImpactPRs},
	}

	err := writeReleaseNotes(w, sections)
	if err != nil {
		return fmt.Errorf("failed to write release notes: %w", err)
	}

	return nil
}