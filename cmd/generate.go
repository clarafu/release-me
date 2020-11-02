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

	// Fetch previous 50 releases from the repository and grab the commit hash
	// associated to each release
	releaseSHAs, err := client.FetchCommitsFromReleases(githubOwner, githubRepo)
	if err != nil {
		failf("failed to fetch release commit SHAs from github: %s", err)
	}

	githubBranch, _ := cmd.Flags().GetString("github-branch")

	// Starting from the latest commit on the branch, we want to walk backwards
	// and compare each commit SHA to the list of release commit SHAs. Once we
	// find a match, this is the the point at which we want to start generating
	// the release notes for.
  //
	// We also will skip over any commit SHAs that are
	// associated to a patch release because we only want to start from a
	// major/minor release.
	startingCommitSHA, err := client.FetchLatestReleaseCommitFromBranch(githubOwner, githubRepo, githubBranch, releaseSHAs)
	if err != nil {
		failf("failed to fetch latest release commit from branch: %s", err)
	}

	                       // 5904,5905
	// release/6.5.x: 	6.5.0 -------------- 6.5.1
												// 6602,5904,5905
	// release/6.6.x:  6.5.0 ---------------6.5.1--------------- current (6.6.0)

	// patchReleases []patches := {6.5.1}

	// Fetch pull requests from patch releases that were skipped over while
	// finding the starting commit SHA. These pull requests will be used to know
	// which pull requests to ignore within the release note generation. This is
	// because we don't want to include any pull requests that have already been
	// mentioned in previous patch releases
	patchReleasesPRs, err := client.FetchPullRequestsFromPatchReleases(githubOwner, githubRepo, githubBranch, releaseSHAs)
	if err != nil {
		failf("failed to fetch pull requests from patches: %s", err)
	}


	lastCommitSHA, _ := cmd.Flags().GetString("last-commit-SHA")

	// Fetch all pull requests that are associated to a commit after the starting
	// commit SHA. If the pull request is already used for a patch release, it is
	// not included.
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
