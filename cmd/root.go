package cli

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "ck",
	Short: "ContextKeeper - Manage your project context",
	Long:  `A CLI tool for managing context and configuration across projects.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		// Handle error
	}
}
