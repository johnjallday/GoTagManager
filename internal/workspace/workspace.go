package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// ParseWSInfo parses the ws_info.toml file into a WorkspaceInfo struct.
func ParseWSInfo(filePath string) (*WorkspaceInfo, error) {
	var info WorkspaceInfo

	// Open the TOML file.
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the TOML data into the WorkspaceInfo struct.
	if _, err := toml.NewDecoder(file).Decode(&info); err != nil {
		return nil, err
	}

	return &info, nil
}

// ListWorkspaces returns a list of workspace directories containing ws_info.toml.
func ListWorkspaces(root string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var workspaces []string
	for _, entry := range entries {
		// Skip hidden directories and ensure the entry is a directory.
		if entry.Name()[0] == '.' || !entry.IsDir() {
			continue
		}

		wsInfoPath := filepath.Join(root, entry.Name(), "ws_info.toml")
		if _, err := os.Stat(wsInfoPath); err == nil {
			workspaces = append(workspaces, filepath.Join(root, entry.Name()))
		}
	}
	return workspaces, nil
}

// ListAliases collects all aliases from each workspace's ws_info.toml.
// It returns a map where the key is the alias name and the value is the workspace name.
func ListAliases(workspaces []string, rootDirectory string) (map[string]string, error) {
	aliases := make(map[string]string)

	for _, workspacePath := range workspaces {
		wsInfoPath := filepath.Join(workspacePath, "ws_info.toml")
		info, err := ParseWSInfo(wsInfoPath)
		if err != nil {
			fmt.Printf("Failed to parse %s: %v\n", wsInfoPath, err)
			continue
		}

		for _, alias := range info.Info.Aliases {
			workspaceName := filepath.Base(workspacePath)
			if _, exists := aliases[alias]; exists {
				fmt.Printf("Warning: Duplicate alias '%s' found in workspace '%s'. Overwriting previous definition.\n", alias, workspaceName)
			}
			aliases[alias] = workspaceName
		}
	}

	return aliases, nil
}

// ListFilesAndDirectories lists all files and directories in the given workspace path.
// It returns two slices: one for directories and one for files.
func ListFilesAndDirectories(workspacePath string) ([]string, []string, error) {
	entries, err := os.ReadDir(workspacePath)
	if err != nil {
		return nil, nil, err
	}

	var dirs []string
	var files []string
	for _, entry := range entries {
		name := entry.Name()

		// Skip the ws_info.toml file and all hidden files/directories
		if name == "ws_info.toml" || strings.HasPrefix(name, ".") {
			continue
		}

		if entry.IsDir() {
			dirs = append(dirs, name)
		} else {
			files = append(files, name)
		}
	}

	return dirs, files, nil
}

// GetWorkspaceSize calculates the total size of all files within the workspace
func GetWorkspaceSize(workspacePath string) (int64, error) {
	var totalSize int64 = 0

	err := filepath.Walk(workspacePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Skip files or directories that cause errors
			fmt.Printf("Warning: Skipping '%s' due to error: %v\n", path, err)
			return nil
		}

		// Skip the ws_info.toml file itself
		if filepath.Base(path) == "ws_info.toml" {
			return nil
		}

		if !info.IsDir() {
			totalSize += info.Size()
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return totalSize, nil
}
