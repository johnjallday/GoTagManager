package cmd

import (
	"log"

	"github.com/johnjallday/GoTagManager/internal/commands"
	"github.com/spf13/cobra"
)

// InfoCmd is the Cobra command for displaying workspace information
var InfoCmd = &cobra.Command{
	Use:   "info [workspace]",
	Short: "Display detailed information about a workspace",
	Long:  `Displays the accounts, tags, and aliases defined in the ws_info.toml of a specific workspace.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := commands.InfoCommand(cfg, args)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(InfoCmd)
}
