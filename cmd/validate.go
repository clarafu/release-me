package cmd

import (
	"fmt"

	"github.com/clarafu/release-me/github"
	"github.com/spf13/cobra"
)

var ValidLabels = map[string]bool{
	"breaking":          true,
	"release/no-impact": true,
	"enhancement":       true,
	"bug":               true,
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validates ..",
	Long:  `TODO`,
	Run:   validate,
}

func init() {
	validateCmd.Flags().Int("pr-number", 0, "checks the existance of labels required for generating release notes")

	validateCmd.MarkFlagRequired("pr-number")
}

func validate(cmd *cobra.Command, args []string) {
	githubToken, _ := cmd.Flags().GetString("github-token")

	client := github.New(githubToken)

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

	var hasValidLabels bool
dance:
	for _, label := range labels {
		if _, exists := ValidLabels[label]; exists {
			hasValidLabels = true
			break dance
		}
	}

	if !hasValidLabels {
		failf("pull request #%d does not have valid labels", prNumber)
	}

	fmt.Printf("pull request #%d has valid labels", prNumber)
}

