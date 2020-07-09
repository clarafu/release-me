package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "releaseme",
		Short: "",
		Long:  `TODO`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().String("github-owner", "", "the login field of a github user or organization")
	rootCmd.PersistentFlags().String("github-repo", "", "the name of the github repository")
	rootCmd.PersistentFlags().String("github-token", "", "github oauth token to authenticate with")

	rootCmd.MarkFlagRequired("github-token")
	rootCmd.MarkFlagRequired("github-owner")
	rootCmd.MarkFlagRequired("github-repo")

	rootCmd.AddCommand(generateCmd)
}
