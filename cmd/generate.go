package cmd

import (
	"fmt"
	"os"

	"github.com/clarafu/release-me/generate"
	"github.com/clarafu/release-me/github"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates a release note using pull requests",
	Long: `A release note is generated through fetching all the pull requests
	merged after the latest tag (release) of the repository. The release note
	is outputted to stdout.`,
	Run: generateReleaseNote,
}

func init() {
	generateCmd.Flags().String("github-branch", "master", "the branch name of the github repository to pull the pull requests from")
	generateCmd.Flags().String("last-commit-SHA", "", "will generate a release note using all prs merged up to this commit SHA. If empty, will generate release note until latest commit.")
}

func generateReleaseNote(cmd *cobra.Command, args []string) {
	githubToken, _ := cmd.Flags().GetString("github-token")

	client := github.New(githubToken)

	githubOwner, _ := cmd.Flags().GetString("github-owner")
	githubRepo, _ := cmd.Flags().GetString("github-repo")

	releaseSHAs, err := client.FetchCommitsFromReleases(githubOwner, githubRepo)
	if err != nil {
		failf("failed to fetch release commit SHAs from github: %s", err)
	}

	githubBranch, _ := cmd.Flags().GetString("github-branch")

	startingCommitSHA, err := client.FetchLatestReleaseCommitFromBranch(githubOwner, githubRepo, githubBranch, releaseSHAs)
	if err != nil {
		failf("failed to fetch latest release commit from branch: %s", err)
	}

	lastCommitSHA, _ := cmd.Flags().GetString("last-commit-SHA")

	pullRequests, err := client.FetchPullRequestsAfterCommit(githubOwner, githubRepo, githubBranch, startingCommitSHA, lastCommitSHA)
	if err != nil {
		failf("failed to fetch pull requests: %s", err)
	}

	g := generate.New(generate.NewReleaseNoteTemplater(os.Stdout))

	err = g.Generate(pullRequests)
	if err != nil {
		failf("failed to generate release note: %s", err)
	}
}

func failf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
