package cmd

import (
	"log"

	"github.com/johnjallday/GoTagManager/internal/commands"
	"github.com/spf13/cobra"
)

// AliasesCmd is the Cobra command for listing aliases
var AliasesCmd = &cobra.Command{
	Use:   "aliases",
	Short: "List all aliases for each workspace",
	Long:  `Displays all aliases defined in the ws_info.toml files across all workspaces.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := commands.ListAliasesCommand(cfg, args)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(AliasesCmd)
}
