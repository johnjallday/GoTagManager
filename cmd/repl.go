package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/johnjallday/GoTagManager/internal/commands"
	"github.com/johnjallday/GoTagManager/internal/workspace"
	"github.com/spf13/cobra"
)

// REPLCmd represents the repl command
var REPLCmd = &cobra.Command{
	Use:   "repl",
	Short: "Start the GoTagManager REPL",
	Long:  `Starts an interactive Read-Eval-Print Loop (REPL) for GoTagManager.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to GoTagManager REPL! Type 'help' to see available commands.")

		p := prompt.New(
			executor,
			completer,
			prompt.OptionPrefix("GoTagManager> "),
			prompt.OptionTitle("GoTagManager REPL"),
			prompt.OptionHistory([]string{}),
		)
		p.Run()
	},
}

func init() {
	rootCmd.AddCommand(REPLCmd)
}

// executor processes the input and executes the corresponding command
func executor(input string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}

	// Split input into command and arguments
	args := strings.Fields(input)
	command := strings.ToLower(args[0])

	switch command {
	case "help":
		printHelp()
	case "exit", "quit":
		fmt.Println("Goodbye!")
		os.Exit(0)
	case "list":
		err := commands.ListWorkspacesCommand(cfg, args[1:])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "aliases":
		err := commands.ListAliasesCommand(cfg, args[1:])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "generate-aliases":
		err := commands.GenerateAliasesCommand(cfg, args[1:])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "info":
		err := commands.InfoCommand(cfg, args[1:])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "load_workspace":
		var workspaceName string
		var err error
		if len(args) >= 2 {
			workspaceName = args[1]
		} else {
			// Interactive selection within REPL using go-prompt
			workspaceName, err = commands.SelectWorkspaceInteractive(cfg)
			if err != nil {
				fmt.Printf("Error selecting workspace: %v\n", err)
				return
			}
		}
		err = commands.LoadWorkspaceCommand(cfg, workspaceName)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "get_size":
		var workspaceName string
		var err error
		if len(args) >= 2 {
			workspaceName = args[1]
		} else {
			// Interactive selection within REPL using go-prompt
			workspaceName, err = commands.SelectWorkspaceInteractive(cfg)
			if err != nil {
				fmt.Printf("Error selecting workspace: %v\n", err)
				return
			}
		}
		err = commands.GetSizeCommand(cfg, workspaceName)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}

// completer provides auto-completion suggestions
func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "list", Description: "List all workspaces"},
		{Text: "aliases", Description: "List all aliases"},
		{Text: "generate-aliases", Description: "Generate shell aliases"},
		{Text: "info", Description: "Display workspace information"},
		{Text: "load_workspace", Description: "Load a workspace and display its information"},
		{Text: "get_size", Description: "Calculate and display the size of a workspace"},
		{Text: "help", Description: "Show help information"},
		{Text: "exit", Description: "Exit the REPL"},
		{Text: "quit", Description: "Exit the REPL"},
	}

	// Handle auto-completion for commands that require workspace names
	if strings.HasPrefix(d.TextBeforeCursor(), "info ") ||
		strings.HasPrefix(d.TextBeforeCursor(), "load_workspace ") ||
		strings.HasPrefix(d.TextBeforeCursor(), "get_size ") {
		// Suggest workspace names
		workspaces, err := workspace.ListWorkspaces(cfg.RootDirectory)
		if err == nil {
			for _, ws := range workspaces {
				s = append(s, prompt.Suggest{Text: filepath.Base(ws), Description: "Workspace"})
			}
		}
	}

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

// printHelp displays help information in REPL
func printHelp() {
	helpText := `
Available Commands:
  list                     List all workspaces
  aliases                  List all aliases for each workspace
  generate-aliases         Generate shell alias commands for .zshrc
  info [workspace]         Display detailed information about a workspace
  load_workspace [workspace]  Load a workspace and display its information
  get_size [workspace]     Calculate and display the size of a workspace
  help                     Show help information
  exit, quit               Exit the REPL
`
	fmt.Println(helpText)
}
