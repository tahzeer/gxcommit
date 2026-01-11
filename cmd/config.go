package cmd

import (
	"fmt"
	"strings"

	"gxcommit/internal/config"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage gxcommit configuration",
}

var configSetCmd = &cobra.Command{
	Use:   "set KEY=VALUE",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		parts := strings.SplitN(args[0], "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("invalid format, expected KEY=VALUE")
		}

		key := parts[0]
		value := parts[1]

		section := "groq"
		if err := config.Set(section, key, value); err != nil {
			return err
		}

		fmt.Printf("✓ %s saved to ~/.gxconfig\n", key)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
}
