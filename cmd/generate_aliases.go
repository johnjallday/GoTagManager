package cmd

import (
	"log"

	"github.com/johnjallday/GoTagManager/internal/commands"
	"github.com/spf13/cobra"
)

// GenerateAliasesCmd is the Cobra command for generating shell aliases
var GenerateAliasesCmd = &cobra.Command{
	Use:   "generate-aliases",
	Short: "Generate shell alias commands for .zshrc",
	Long:  `Generates alias commands based on ws_info.toml files, which can be added to your .zshrc for quick navigation.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := commands.GenerateAliasesCommand(cfg, args)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(GenerateAliasesCmd)
}
