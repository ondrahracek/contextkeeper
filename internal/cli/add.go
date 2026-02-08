// Package cli provides the command-line interface for ContextKeeper.
//
// This package implements the Cobra-based CLI for managing project context
// and configuration. See the root.go file for the main command structure.
package cli

import (
	"time"

	"github.com/ondrahracek/contextkeeper/internal/config"
	"github.com/ondrahracek/contextkeeper/internal/models"
	"github.com/ondrahracek/contextkeeper/internal/storage"
	"github.com/ondrahracek/contextkeeper/internal/utils"
	"github.com/spf13/cobra"
)

// addCmd adds a new context item to storage.
//
// The command supports three ways to provide content:
//   - As a command-line argument
//   - Via stdin (when no arguments provided)
//   - Through an interactive editor (--editor flag)
var addCmd = &cobra.Command{
	Use:     "add [content]",
	Short:   "Add a new context item",
	Long:    "Add a new context item to storage. Content can be provided as an argument, via stdin, or using --editor.",
	Example: `  # Add content directly
  ck add "Remember to update documentation"

  # Add with project and tags
  ck add "Fix bug #123" --project "web-app" --tags "bug,urgent"

  # Open editor for multi-line content
  ck add --editor

  # Add from stdin
  echo "Quick note" | ck add`,
	Args:  cobra.MaximumNArgs(1),
	RunE: addCommand,
}

// Command flags for the add command.
var (
	// projectFlag specifies the project name for the new item
	projectFlag string
	// tagStr contains comma-separated tags for the new item
	tagStr string
	// useEditor opens the system editor for content input
	useEditor bool
)

// addCommand is the execution function for the add command.
// It creates a new context item and saves it to storage.
func addCommand(cmd *cobra.Command, args []string) error {
	// Load configuration to get storage path
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	var content string

	// Determine content source: argument, editor, or stdin
	switch {
	case len(args) > 0:
		content = args[0]
	case useEditor:
		content, err = utils.OpenEditor("")
		if err != nil {
			return err
		}
	default:
		// No content provided, show help
		return cmd.Help()
	}

	// Skip empty content
	if content == "" {
		return nil
	}

	// Parse and validate tags
	tags := utils.ParseTags(tagStr)
	if err := utils.ValidateTags(tags); err != nil {
		return err
	}

	// Get project from flag or config default
	project := projectFlag
	if project == "" {
		project = config.GetDefaultProject()
	}

	// Create the new item using models.ContextItem
	now := time.Now()
	item := models.ContextItem{
		ID:        utils.GenerateUUID(),
		Content:   content,
		Project:   project,
		Tags:      tags,
		CreatedAt: now,
	}

	// Initialize storage and add the item
	stor := storage.NewStorage(cfg.StoragePath)
	if err := stor.Load(); err != nil {
		return err
	}

	if err := stor.Add(item); err != nil {
		return err
	}

	if err := stor.Save(); err != nil {
		return err
	}

	cmd.Println("Added context item")
	return nil
}

// init registers the add command with the root command.
func init() {
	// Register command flags
	addCmd.Flags().StringVarP(&projectFlag, "project", "p", "", "Project name for the context item")
	addCmd.Flags().StringVarP(&tagStr, "tags", "t", "", "Tags for the context item (comma or space separated)")
	addCmd.Flags().BoolVarP(&useEditor, "editor", "e", false, "Open editor to enter content")

	// Add command to root
	RootCmd.AddCommand(addCmd)
}
