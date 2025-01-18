package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config holds the configuration settings.
type Config struct {
	RootDirectory string `mapstructure:"root_directory"`
}

// LoadConfig initializes Viper, reads the config file, and environment variables.
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()

	if configPath != "" {
		// Use the specified config file
		v.SetConfigFile(configPath)
	} else {
		// Look for config.toml in the current directory
		v.SetConfigName("config") // name of config file (without extension)
		v.SetConfigType("toml")   // REQUIRED if the config file does not have the extension in the name
		v.AddConfigPath(".")      // look for config in the current directory
	}

	// Set default values
	v.SetDefault("root_directory", "/Users/jj/Workspace/")

	// Bind environment variables
	v.AutomaticEnv() // read in environment variables that match

	// Read the config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; proceed with defaults and environment variables
			fmt.Println("Config file not found; using default values and environment variables.")
		} else {
			// Config file was found but another error was produced
			return nil, fmt.Errorf("fatal error config file: %w", err)
		}
	} else {
		fmt.Println("Using config file:", v.ConfigFileUsed())
	}

	// Bind specific environment variables to config fields
	err := v.BindEnv("root_directory", "WORKSPACE_ROOT")
	if err != nil {
		return nil, fmt.Errorf("error binding environment variable: %w", err)
	}

	// Populate the Config struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	// Validate the root directory
	if _, err := os.Stat(cfg.RootDirectory); os.IsNotExist(err) {
		return nil, fmt.Errorf("root directory does not exist: %s", cfg.RootDirectory)
	}
	fmt.Printf("Loaded root_directory: %s\n", cfg.RootDirectory)

	return &cfg, nil
}
