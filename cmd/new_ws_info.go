package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/johnjallday/GoTagManager/internal/tag"
	"github.com/spf13/cobra"
)

// NewWsInfoCmd is a command that creates a new ws_info.toml for the given workspace
var NewWsInfoCmd = &cobra.Command{
	Use:   "new_ws_info [workspace]",
	Short: "Create a default ws_info.toml in the specified workspace",
	Long: `Creates a new ws_info.toml in the specified workspace directory if 
it does not already exist.`,
	Args: cobra.ExactArgs(1), // We need exactly 1 argument: workspace name
	Run: func(cmd *cobra.Command, args []string) {
		workspaceName := args[0]

		// Build the path to the workspace from the config root
		workspacePath := filepath.Join(cfg.RootDirectory, workspaceName)

		if err := tag.CreateWsInfoToml(workspacePath); err != nil {
			log.Fatalf("Error creating ws_info.toml: %v", err)
		} else {
			fmt.Printf("ws_info.toml successfully created/updated in %s.\n", workspacePath)
		}
	},
}

func init() {
	rootCmd.AddCommand(NewWsInfoCmd)
}
