package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/johnjallday/GoTagManager/config"
	"github.com/johnjallday/GoTagManager/internal/workspace"
)

// ListWorkspacesCommand lists all workspaces
func ListWorkspacesCommand(cfg *config.Config, args []string) error {
	workspaces, err := workspace.ListWorkspaces(cfg.RootDirectory)
	if err != nil {
		return fmt.Errorf("failed to list workspaces: %w", err)
	}

	if len(workspaces) == 0 {
		fmt.Println("No valid workspaces found.")
		return nil
	}

	fmt.Println("Workspaces:")
	for _, ws := range workspaces {
		fmt.Printf("- %s\n", filepath.Base(ws))
	}
	return nil
}

// ListAliasesCommand lists all aliases across workspaces
func ListAliasesCommand(cfg *config.Config, args []string) error {
	workspaces, err := workspace.ListWorkspaces(cfg.RootDirectory)
	if err != nil {
		return fmt.Errorf("failed to list workspaces: %w", err)
	}

	if len(workspaces) == 0 {
		fmt.Println("No valid workspaces found.")
		return nil
	}

	allAliases, err := workspace.ListAliases(workspaces, cfg.RootDirectory)
	if err != nil {
		return fmt.Errorf("failed to list aliases: %w", err)
	}

	if len(allAliases) == 0 {
		fmt.Println("No aliases found in any workspace.")
		return nil
	}

	fmt.Println("Aliases for each workspace:")

	// Sort aliases for consistent output.
	aliasNames := make([]string, 0, len(allAliases))
	for alias := range allAliases {
		aliasNames = append(aliasNames, alias)
	}
	sort.Strings(aliasNames)

	for _, alias := range aliasNames {
		workspaceName := allAliases[alias]
		fmt.Printf("Alias: %s => Workspace: %s\n", alias, workspaceName)
	}
	return nil
}

// GenerateAliasesCommand generates shell aliases
func GenerateAliasesCommand(cfg *config.Config, args []string) error {
	workspaces, err := workspace.ListWorkspaces(cfg.RootDirectory)
	if err != nil {
		return fmt.Errorf("failed to list workspaces: %w", err)
	}

	if len(workspaces) == 0 {
		fmt.Println("No valid workspaces found.")
		return nil
	}

	allAliases, err := workspace.ListAliases(workspaces, cfg.RootDirectory)
	if err != nil {
		return fmt.Errorf("failed to list aliases: %w", err)
	}

	if len(allAliases) == 0 {
		fmt.Println("No aliases found in any workspace.")
		return nil
	}

	// Map workspace names to their paths for quick lookup
	workspacePathMap := make(map[string]string)
	for _, workspacePath := range workspaces {
		workspaceName := filepath.Base(workspacePath)
		workspacePathMap[workspaceName] = workspacePath
	}

	fmt.Println("# Generated Aliases for GoTagManager")
	// Sort aliases for consistent output.
	aliasNames := make([]string, 0, len(allAliases))
	for alias := range allAliases {
		aliasNames = append(aliasNames, alias)
	}
	sort.Strings(aliasNames)

	for _, alias := range aliasNames {
		workspaceName := allAliases[alias]
		workspacePath, exists := workspacePathMap[workspaceName]
		if !exists {
			fmt.Printf("# Warning: Workspace '%s' not found for alias '%s'\n", workspaceName, alias)
			continue
		}
		fmt.Printf("alias %s=\"cd '%s'\"\n", alias, workspacePath)
	}
	return nil
}

// InfoCommand displays information about a specific workspace
func InfoCommand(cfg *config.Config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("workspace name is required")
	}
	workspaceName := args[0]
	wsPath := filepath.Join(cfg.RootDirectory, workspaceName)
	wsInfoPath := filepath.Join(wsPath, "ws_info.toml")

	info, err := workspace.ParseWSInfo(wsInfoPath)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", wsInfoPath, err)
	}

	fmt.Printf("Contents of %s:\n", wsInfoPath)

	// Print Accounts
	if len(info.Accounts) > 0 {
		fmt.Printf("Accounts:\n")
		for key, value := range info.Accounts {
			fmt.Printf("  %s = %s\n", key, value)
		}
	}

	// Print Tags
	if len(info.Info.Tags) > 0 {
		fmt.Printf("Tags:\n")
		for _, tag := range info.Info.Tags {
			fmt.Printf("  - %s\n", tag)
		}
	}

	// Print Aliases
	if len(info.Info.Aliases) > 0 {
		fmt.Printf("Aliases:\n")
		for _, alias := range info.Info.Aliases {
			fmt.Printf("  - %s\n", alias)
		}
	}
	return nil
}

// SelectWorkspaceInteractive lists all workspaces with numbers and prompts the user to select one.
// It returns the selected workspace name.
func SelectWorkspaceInteractive(cfg *config.Config) (string, error) {
	workspaces, err := workspace.ListWorkspaces(cfg.RootDirectory)
	if err != nil {
		return "", fmt.Errorf("failed to list workspaces: %w", err)
	}

	if len(workspaces) == 0 {
		return "", fmt.Errorf("no workspaces found in root directory '%s'", cfg.RootDirectory)
	}

	fmt.Println("Available Workspaces:")
	for i, ws := range workspaces {
		fmt.Printf("%d) %s\n", i+1, filepath.Base(ws))
	}

	// Use go-prompt's Input to get user input
	promptText := fmt.Sprintf("Enter the number of the workspace to load (1-%d): ", len(workspaces))
	input := prompt.Input(promptText, func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	})

	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("no input provided")
	}

	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(workspaces) {
		return "", fmt.Errorf("invalid input")
	}

	selectedWorkspace := filepath.Base(workspaces[choice-1])
	return selectedWorkspace, nil
}

// LoadWorkspaceCommand loads a workspace, displays its ws_info.toml, and lists all files and directories.
func LoadWorkspaceCommand(cfg *config.Config, workspaceName string) error {
	workspacePath := filepath.Join(cfg.RootDirectory, workspaceName)
	wsInfoPath := filepath.Join(workspacePath, "ws_info.toml")

	// Check if the workspace exists
	if _, err := os.Stat(workspacePath); os.IsNotExist(err) {
		return fmt.Errorf("workspace '%s' does not exist in root directory '%s'", workspaceName, cfg.RootDirectory)
	}

	// Parse ws_info.toml
	info, err := workspace.ParseWSInfo(wsInfoPath)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", wsInfoPath, err)
	}

	// Display ws_info.toml contents
	fmt.Printf("\nContents of %s:\n", wsInfoPath)

	// Print Accounts
	if len(info.Accounts) > 0 {
		fmt.Printf("Accounts:\n")
		for key, value := range info.Accounts {
			fmt.Printf("  %s = %s\n", key, value)
		}
	} else {
		fmt.Println("No Accounts defined.")
	}

	// Print Tags
	if len(info.Info.Tags) > 0 {
		fmt.Printf("Tags:\n")
		for _, tag := range info.Info.Tags {
			fmt.Printf("  - %s\n", tag)
		}
	} else {
		fmt.Println("No Tags defined.")
	}

	// Print Aliases
	if len(info.Info.Aliases) > 0 {
		fmt.Printf("Aliases:\n")
		for _, alias := range info.Info.Aliases {
			fmt.Printf("  - %s\n", alias)
		}
	} else {
		fmt.Println("No Aliases defined.")
	}

	// List all files and directories in the workspace
	fmt.Printf("\nContents of workspace '%s':\n", workspaceName)
	dirs, files, err := workspace.ListFilesAndDirectories(workspacePath)
	if err != nil {
		return fmt.Errorf("failed to list files and directories: %w", err)
	}

	if len(dirs) > 0 {
		fmt.Println("Directories:")
		for _, dir := range dirs {
			fmt.Printf("  - %s\n", dir)
		}
	} else {
		fmt.Println("No subdirectories found.")
	}

	if len(files) > 0 {
		fmt.Println("Files:")
		for _, file := range files {
			fmt.Printf("  - %s\n", file)
		}
	} else {
		fmt.Println("No files found.")
	}

	return nil
}

// GetSizeCommand calculates and displays the size of a specified workspace
func GetSizeCommand(cfg *config.Config, workspaceName string) error {
	workspacePath := filepath.Join(cfg.RootDirectory, workspaceName)
	wsInfoPath := filepath.Join(workspacePath, "ws_info.toml")

	// Check if the workspace exists
	if _, err := workspace.ParseWSInfo(wsInfoPath); err != nil {
		return fmt.Errorf("workspace '%s' does not exist or has an invalid ws_info.toml", workspaceName)
	}

	// Calculate the size
	sizeBytes, err := workspace.GetWorkspaceSize(workspacePath)
	if err != nil {
		return fmt.Errorf("failed to calculate size for workspace '%s': %w", workspaceName, err)
	}

	// Format the size into a human-readable format
	humanReadableSize := formatBytes(sizeBytes)

	fmt.Printf("Total size of workspace '%s': %s\n", workspaceName, humanReadableSize)
	return nil
}

// formatBytes converts bytes to a human-readable string
func formatBytes(bytes int64) string {
	const (
		KB = 1 << (10 * (iota + 1))
		MB
		GB
		TB
		PB
	)

	switch {
	case bytes >= PB:
		return fmt.Sprintf("%.2f PB", float64(bytes)/PB)
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/TB)
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
