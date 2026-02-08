// Package cli provides the command-line interface for ContextKeeper.
//
// This package implements the Cobra-based CLI for managing context and
// configuration. It supports commands for adding, listing, editing,
// completing, and removing context items.
//
// # Commands
//
// The CLI provides the following main commands:
//
//	- add:     Add a new context item
//	- list:    List context items with optional filters
//	- edit:    Edit an existing context item
//	- done:    Mark a context item as completed
//	- remove:  Remove a context item
//	- status:  Show a quick overview
//	- init:    Initialize a new ContextKeeper directory
package cli

import (
	"github.com/spf13/cobra"
)

// RootCmd is the base command for the ContextKeeper CLI.
//
// All other commands are registered as subcommands of RootCmd.
// The command name is "ck" as defined by the Use field.
var RootCmd = &cobra.Command{
	Use:   "ck",
	Short: "ContextKeeper - Manage your project context",
	Long: `ContextKeeper is a CLI tool for managing context and configuration across projects.

It helps you keep track of important notes, tasks, and information organized
by project with support for tags and completion tracking.

Use "ck [command] --help" to get more information about a specific command.`,
	Example: `  # Initialize ContextKeeper in current directory
  ck init

  # Add a new context item
  ck add "Remember to update documentation"

  # Add with project and tags
  ck add "Fix bug #123" --project "web-app" --tags "bug,urgent"

  # List all active items
  ck list

  # List items for a specific project
  ck list --project "my-project"

  # Mark an item as done
  ck done abc12345

  # Edit an item in editor
  ck edit abc12345`,
}

// Execute runs the root command and handles any errors.
// This function is called from main.go to start the CLI.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		// Error handling is managed by Cobra
		// which will print the error and exit with appropriate code
	}
}