package cmd

import (
	"fmt"
	"strconv"

	"github.com/clarafu/release-me/generate"
	"github.com/clarafu/release-me/github"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validates the pull request has correct labels.",
	Long: `Ensures that the pull request given has at least one of the labels
	required to properly generate a release note using the "generate"
	command.`,
	Run: validate,
}

func init() {
	validateCmd.Flags().Int("pr-number", 0, "checks the existance of labels required for generating release notes")

	validateCmd.MarkFlagRequired("pr-number")
}

func validate(cmd *cobra.Command, args []string) {
	githubToken, _ := cmd.Flags().GetString("github-token")
	githubV4Endpoint, _ := cmd.Flags().GetString("v4-endpoint")

	client := github.New(githubToken, githubV4Endpoint)

	githubOwner, _ := cmd.Flags().GetString("github-owner")
	githubRepo, _ := cmd.Flags().GetString("github-repo")

	prNumber, err := cmd.Flags().GetInt("pr-number")
	if err != nil {
		failf("failed to get pr number: %s", err)
	}

	labels, err := client.FetchLabelsForPullRequest(githubOwner, githubRepo, prNumber)
	if err != nil {
		failf("failed fetch labels for pull request: %s", err)
	}

	hasValidLabels := generate.Validate(labels)
	if !hasValidLabels {
		failf("invalid pull request %s", generate.PullRequestsNotLabelled{Identifiers: []string{strconv.Itoa(prNumber)}})
	}

	fmt.Printf("pull request #%d has valid labels\n", prNumber)
}
