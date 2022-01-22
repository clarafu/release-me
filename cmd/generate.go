package cmd

import (
	"fmt"
	"os"
	"regexp"

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
	generateCmd.Flags().String("release-version", "", "the version that the release note will be generated for")
	generateCmd.Flags().StringSlice("ignore-authors", nil, "comma separated list of github handles, any PRs authored by these handles will be ignored.")
	generateCmd.Flags().String("ignore-release-regex", "", "a regular expression indicating releases to ignore when determining the previous release")
	generateCmd.MarkFlagRequired("release-version")
}

func generateReleaseNote(cmd *cobra.Command, args []string) {
	githubToken, _ := cmd.Flags().GetString("github-token")
	githubV4Endpoint, _ := cmd.Flags().GetString("v4-endpoint")

	client := github.New(githubToken, githubV4Endpoint)

	githubOwner, _ := cmd.Flags().GetString("github-owner")
	githubRepo, _ := cmd.Flags().GetString("github-repo")

	ignoreReleaseRegexStr, _ := cmd.Flags().GetString("ignore-release-regex")

	// Fetch previous 50 releases from the repository and grab the commit hash
	// associated to each release
	releaseSHAs, err := client.FetchCommitsFromReleases(githubOwner, githubRepo)
	if err != nil {
		failf("failed to fetch release commit SHAs from github: %s", err)
	}

	if ignoreReleaseRegexStr != "" {
		filteredReleaseSHAs := make(map[string]string)
		ignoreReleaseRegex, err := regexp.Compile(ignoreReleaseRegexStr)
		if err != nil {
			failf("invalid regex in --ignore-release-regex: %s", err)
		}
		for oid, release := range releaseSHAs {
			if !ignoreReleaseRegex.MatchString(release) {
				filteredReleaseSHAs[oid] = release
			}
		}
		releaseSHAs = filteredReleaseSHAs
	}

	githubBranch, _ := cmd.Flags().GetString("github-branch")

	// Starting from the latest commit on the branch, we want to walk backwards
	// and compare each commit SHA to the list of release commit SHAs. Once we
	// find a match, this is the the point at which we want to start generating
	// the release notes for.
	versionToRelease, _ := cmd.Flags().GetString("release-version")
	startingCommitSHA, err := client.FetchLatestReleaseCommitFromBranch(githubOwner, githubRepo, githubBranch, versionToRelease, releaseSHAs)
	if err != nil {
		failf("failed to fetch latest release commit from branch: %s", err)
	}

	lastCommitSHA, _ := cmd.Flags().GetString("last-commit-SHA")
	ignoreAuthors, _ := cmd.Flags().GetStringSlice("ignore-authors")

	// Fetch all pull requests that are associated to a commit after the starting
	// commit SHA. If the pull request is already used for a patch release, it is
	// not included.
	pullRequests, err := client.FetchPullRequestsAfterCommit(githubOwner, githubRepo, githubBranch, startingCommitSHA, lastCommitSHA, ignoreAuthors)
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
