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
	Short: "generate",
	Long:  `TODO`,
	Run:   generateReleaseNote,
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

	err = generate.Generate(os.Stdout, pullRequests)
	if err != nil {
		failf("failed to generate release note: %s", err)
	}
}

func failf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
