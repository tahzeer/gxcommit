package cmd

import (
	"fmt"
	"os"

	"github.com/tahzeer/gxcommit/internal/app"

	"github.com/spf13/cobra"
)

var saveScriptCmd = &cobra.Command{
	Use:     "save-script",
	Aliases: []string{"ss"},
	Short:   "Generate commit script without executing",
	RunE: func(cmd *cobra.Command, args []string) error {
		script := app.GenerateScript(jiraCode)
		if script == "" {
			return nil
		}

		err := os.WriteFile("gxcommit.sh", []byte(script), 0755)
		if err != nil {
			return err
		}

		fmt.Println("Saved commit script to gxcommit.sh")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(saveScriptCmd)
}
