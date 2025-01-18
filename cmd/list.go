package cmd

import (
	"log"

	"github.com/johnjallday/GoTagManager/internal/commands"
	"github.com/spf13/cobra"
)

// ListCmd is the Cobra command for listing workspaces
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all workspaces",
	Long:  `Lists all directories in the specified root that contain a ws_info.toml file.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := commands.ListWorkspacesCommand(cfg, args)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(ListCmd)
}
