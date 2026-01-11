package cmd

import (
	"fmt"

	"github.com/tahzeer/gxcommit/internal/app"
	"github.com/tahzeer/gxcommit/internal/executil"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:     "run",
	Aliases: []string{"r"},
	Short:   "Generate and run commits immediately",
	RunE: func(cmd *cobra.Command, args []string) error {
		script := app.GenerateScript(jiraCode)
		if script == "" {
			return nil
		}

		fmt.Println("Executing generated commits...")
		return executil.RunScript(script)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
