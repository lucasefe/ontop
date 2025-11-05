package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config represents the application configuration
type Config struct {
	UI UIConfig `toml:"ui"`
}

// UIConfig holds user interface preferences
type UIConfig struct {
	ViewMode string `toml:"view_mode"` // "column" or "row"
}

// GetConfigPath returns the absolute path to the config file.
// Follows XDG Base Directory specification: ~/.config/ontop/ontop.toml
func GetConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "ontop", "ontop.toml")
}

// Load reads the config file from the standard location and returns
// a Config struct. If the file doesn't exist or can't be parsed,
// returns a Config with default values.
//
// Returns:
//   - Config with loaded values (or defaults on error)
//   - error (nil if successful or file not found, non-nil if parse failed)
func Load() (Config, error) {
	// Default config
	cfg := Config{
		UI: UIConfig{
			ViewMode: "column",
		},
	}

	path := GetConfigPath()
	if path == "" {
		return cfg, fmt.Errorf("unable to determine home directory")
	}

	// If file doesn't exist, return defaults (not an error)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}

	// Try to parse the file
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		// Return defaults on parse error, but return the error so caller can log
		return Config{
			UI: UIConfig{
				ViewMode: "column",
			},
		}, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate and sanitize view_mode
	if cfg.UI.ViewMode != "column" && cfg.UI.ViewMode != "row" {
		cfg.UI.ViewMode = "column"
	}

	return cfg, nil
}

// Save writes the config to the standard location. Creates the directory
// if it doesn't exist. Uses atomic write (temp file + rename).
//
// Parameters:
//   - cfg: The Config struct to persist
//
// Returns:
//   - error (nil if successful, non-nil if write failed)
func Save(cfg Config) error {
	path := GetConfigPath()
	if path == "" {
		return fmt.Errorf("unable to determine home directory")
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create temp file in same directory
	tmpFile, err := os.CreateTemp(dir, ".ontop.toml.tmp.*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	// Ensure temp file is cleaned up on error
	defer func() {
		if tmpFile != nil {
			tmpFile.Close()
			os.Remove(tmpPath)
		}
	}()

	// Encode config to TOML and write to temp file
	encoder := toml.NewEncoder(tmpFile)
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	// Close temp file before rename
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}
	tmpFile = nil // Mark as closed so defer doesn't close again

	// Atomic rename
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	// Set permissions
	if err := os.Chmod(path, 0644); err != nil {
		// Non-fatal, just log
		return fmt.Errorf("config saved but failed to set permissions: %w", err)
	}

	return nil
}
