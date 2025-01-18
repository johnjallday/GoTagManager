package workspace

// WorkspaceInfo represents the structure of ws_info.toml.
type WorkspaceInfo struct {
	Accounts map[string]string `toml:"accounts"`
	Info     InfoSection       `toml:"info"`
}

// InfoSection represents the [info] table in ws_info.toml.
type InfoSection struct {
	Tags    []string `toml:"tags"`
	Aliases []string `toml:"aliases"`
}
