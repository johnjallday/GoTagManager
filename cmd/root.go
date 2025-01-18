package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/johnjallday/GoTagManager/config"
	"github.com/spf13/cobra"
)

var cfg *config.Config

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "GoTagManager",
	Short: "A CLI tool to manage workspaces and their aliases.",
	Long: `GoTagManager is a command-line tool designed to help you manage your various workspaces.
It allows you to list workspaces, view aliases, and generate shell aliases for quick navigation.`,
	// Uncomment the following line if your bare application has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Define persistent flags and configuration settings.
	rootCmd.PersistentFlags().StringP("config", "c", "", "Path to the configuration file")

	// Bind flags to configuration
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
}

func initConfig() {
	configPath, err := rootCmd.PersistentFlags().GetString("config")
	if err != nil {
		log.Fatalf("Error reading config flag: %v", err)
	}

	cfg, err = config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}
}
