// Package cli provides the command-line interface for ContextKeeper.
//
// This package implements the Cobra-based CLI for managing project context
// and configuration. See the root.go file for the main command structure.
package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// initCmd initializes a new ContextKeeper directory.
//
// This command creates the .contextkeeper directory structure with
// the necessary storage file (items.json).
var initCmd = &cobra.Command{
	Use:     "init",
	Short:   "Initialize ContextKeeper",
	Long:    "Create the .contextkeeper directory structure in the current directory.",
	Example: `  # Initialize in current directory
  ck init`,
	Args: cobra.NoArgs,
	RunE: initCommand,
}

// initCommand is the execution function for the init command.
// It creates the required directory structure and items.json file.
func initCommand(cmd *cobra.Command, args []string) error {
	// Define the context directory
	contextDir := ".contextkeeper"

	// Create the .contextkeeper directory
	if err := os.MkdirAll(contextDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %q: %w", contextDir, err)
	}

	// Create items.json file for storing context items
	itemsFile := filepath.Join(contextDir, "items.json")
	if _, err := os.Stat(itemsFile); err == nil {
		cmd.Printf("ContextKeeper is already initialized in %s\n", contextDir)
		return nil
	}

	if err := os.WriteFile(itemsFile, []byte("[]"), 0644); err != nil {
		return fmt.Errorf("failed to create items file %q: %w", itemsFile, err)
	}

	cmd.Printf("Initialized ContextKeeper in: %s\n", contextDir)
	cmd.Println("Run 'ck add --help' to get started.")

	return nil
}

// init registers the init command with the root command.
func init() {
	// Add command to root
	RootCmd.AddCommand(initCmd)
}