package cmd

import (
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current configuration",
	Long:  "Display the current encryption/decryption settings.",
	Run: func(cmd *cobra.Command, args []string) {
		showSettings()
	},
}
