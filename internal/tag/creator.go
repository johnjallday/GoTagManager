package tag

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/johnjallday/GoTagManager/internal/workspace"
)

// CreateWsInfoToml checks if ws_info.toml exists in the workspacePath.
// If it does not, it creates a minimal ws_info.toml with defaults.
// If it does exist, this function can either skip or handle updates.
func CreateWsInfoToml(workspacePath string) error {
	wsInfoPath := filepath.Join(workspacePath, "ws_info.toml")

	// Check if the workspace directory actually exists
	if _, err := os.Stat(workspacePath); os.IsNotExist(err) {
		return fmt.Errorf("workspace path does not exist: %s", workspacePath)
	}

	// If ws_info.toml already exists, decide how to handle:
	if _, err := os.Stat(wsInfoPath); err == nil {
		// For example, we can simply return here:
		fmt.Printf("ws_info.toml already exists at %s; skipping creation.\n", wsInfoPath)
		return nil

		// Or if you want to overwrite, you'd remove the file or proceed to write.
	}

	// Create a minimal default structure
	info := workspace.WorkspaceInfo{
		Accounts: map[string]string{
			// Example:
			"default_account": "abc123",
		},
		Info: workspace.InfoSection{
			Tags:    []string{"example-tag"},
			Aliases: []string{"example-alias"},
		},
	}

	// Write it to ws_info.toml
	file, err := os.Create(wsInfoPath)
	if err != nil {
		return fmt.Errorf("unable to create ws_info.toml: %w", err)
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(info); err != nil {
		return fmt.Errorf("failed to encode default ws_info.toml: %w", err)
	}

	fmt.Printf("Created default ws_info.toml at %s\n", wsInfoPath)
	return nil
}
