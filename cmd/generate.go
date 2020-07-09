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
	Run:   generateReleaseNote,
}

func init() {
	generateCmd.Flags().String("github-branch", "master", "the branch name of the github repository to pull the pull requests from")
}

func generateReleaseNote(cmd *cobra.Command, args []string) {
	githubToken, _ := cmd.Flags().GetString("github-token")

	client := github.New(githubToken)

	githubOwner, _ := cmd.Flags().GetString("github-owner")
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
