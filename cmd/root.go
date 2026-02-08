package cli

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "contextkeeper",
	Short: "A CLI tool for managing context and configuration",
	Long:  `ContextKeeper helps you manage your project's context and configuration files.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		// Handle error
	}
}
