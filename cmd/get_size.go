package cmd

import (
	"log"

	"github.com/johnjallday/GoTagManager/internal/commands"
	"github.com/spf13/cobra"
)

// GetSizeCmd is the Cobra command for getting the size of a workspace
var GetSizeCmd = &cobra.Command{
	Use:   "get_size [workspace]",
	Short: "Calculate and display the size of a workspace",
	Long: `Calculates the total size of the specified workspace by summing the sizes of all files within it.
If no workspace is specified, it will list all available workspaces and prompt you to select one.`,
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

		err = commands.GetSizeCommand(cfg, workspaceName)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(GetSizeCmd)
}
