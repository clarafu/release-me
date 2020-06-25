package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/clarafu/release-me/github"
	"github.com/spf13/cobra"
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

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate",
	Long:  `TODO`,
	Run:  generateReleaseNote,
}

func init() {
	generateCmd.Flags().String("github-owner", "", "the login field of a github user or organization")
	generateCmd.Flags().String("github-repo", "", "the name of the github repository")
	generateCmd.Flags().String("github-branch", "master", "the branch name of the github repository to pull the pull requests from")

	generateCmd.MarkFlagRequired("github-owner")
	generateCmd.MarkFlagRequired("github-repo")
}

func generateReleaseNote(cmd *cobra.Command, args []string) {
	githubToken, _ := rootCmd.Flags().GetString("github-token")

	client := github.New(githubToken)

	githubOwner, _  := cmd.Flags().GetString("github-owner")
	githubRepo, _ := cmd.Flags().GetString("github-repo")

	commitSHA, err := client.FetchLatestReleaseCommitSHA(githubOwner, githubRepo)
	if err != nil {
		failf("failed to fetch latest release commit SHA from github: %s", err)
	}

	githubBranch, _ := cmd.Flags().GetString("github-branch")

	pullRequests, err := client.FetchPullRequestsAfterCommit(githubOwner, githubRepo, githubBranch, commitSHA)
	if err != nil {
		 failf("failed to fetch pull requests: %s", err)
	}

	sort.Slice(pullRequests, func(i, j int) bool {
		return pullRequests[i].Number < pullRequests[j].Number
	})

	var breakingPRs []github.PullRequest
	var noImpactPRs []github.PullRequest
	var featurePRs []github.PullRequest
	var bugFixPRs []github.PullRequest
	for _, pr := range pullRequests {
		if pr.HasLabel("breaking") {
			breakingPRs = append(breakingPRs, pr)
			continue
		}

		if pr.HasLabel("release/no-impact") {
			noImpactPRs = append(noImpactPRs, pr)
			continue
		}

		if pr.HasLabel("enhancement") {
			featurePRs = append(featurePRs, pr)
			continue
		}

		if pr.HasLabel("bug") {
			bugFixPRs = append(bugFixPRs, pr)
			continue
		}

		failf(PullRequestNotLabelled{Number: pr.Number}.Error())
	}

	sections := []Section{
		Section{Title: "Breaking", Icon: "🚨", PRs: breakingPRs},
		Section{Title: "Features", Icon: "✈️", PRs: featurePRs},
		Section{Title: "Bug Fixes", Icon: "🐞", PRs: bugFixPRs},
		Section{Title: "No Impact", Icon: "🤷", PRs: noImpactPRs},
	}

	err = WriteReleaseNotes(os.Stdout, sections)
	if err != nil {
		failf("failed to write release notes: %s", err)
	}
}

func failf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}