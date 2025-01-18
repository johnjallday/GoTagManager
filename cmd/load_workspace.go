package cmd

import (
	"log"

	"github.com/johnjallday/GoTagManager/internal/commands"
	"github.com/spf13/cobra"
)

// LoadWorkspaceCmd is the Cobra command for loading a workspace
var LoadWorkspaceCmd = &cobra.Command{
	Use:   "load_workspace [workspace]",
	Short: "Load a workspace and display its information",
	Long: `Loads the specified workspace, displays the contents of its ws_info.toml,
and lists all files and directories within the workspace. If no workspace is specified,
it will list all available workspaces and prompt you to select one.`,
	Args: cobra.MaximumNArgs(1), // Allow 0 or 1 argument
	Run: func(cmd *cobra.Command, args []string) {
		var workspaceName string
		var err error

		if len(args) == 1 {
			workspaceName = args[0]
		} else {
			// No argument provided; prompt user to select a workspace
			workspaceName, err = commands.SelectWorkspaceInteractive(cfg)
			if err != nil {
				log.Fatalf("Error selecting workspace: %v", err)
			}
		}

		err = commands.LoadWorkspaceCommand(cfg, workspaceName)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(LoadWorkspaceCmd)
}
