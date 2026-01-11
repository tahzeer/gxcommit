package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version string = "dev"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gxcommit version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
