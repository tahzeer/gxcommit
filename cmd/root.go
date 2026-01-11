package cmd

import (
	"fmt"
	"os"

	"gxcommit/internal/git"

	"github.com/spf13/cobra"
)

var (
	jiraCode string
)

var rootCmd = &cobra.Command{
	Use:   "gxcommit",
	Short: "AI-assisted Git commit generator",
	Long:  "gxcommit analyzes git diff and generates logical commits using Groq AI.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if !git.IsGitRepo() {
			return fmt.Errorf("gxcommit must be run inside a git repository")
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&jiraCode,
		"code",
		"c",
		"",
		"jira / ticket code to prefix commit messages",
	)
}
