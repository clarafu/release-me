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
	rootCmd.PersistentFlags().String("github-token", "", "github oauth token to authenticate with")
	rootCmd.MarkFlagRequired("github-token")

	rootCmd.AddCommand(generateCmd)
}
