package cmd

import (
	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Start interactive shell mode",
	Long:  "Start an interactive shell with readline support and command completion.",
	Run:   runInteractiveShell,
}