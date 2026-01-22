package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
)

// Config holds all gastop configuration.
type Config struct {
	RefreshInterval    time.Duration `toml:"refresh_interval"`
	StuckThresholdMins int           `toml:"stuck_threshold_minutes"`
	LogLines           int           `toml:"log_lines"`
	ShowLogs           bool          `toml:"show_logs"`

	Paths   PathsConfig   `toml:"paths"`
	Filters FiltersConfig `toml:"filters"`
}

// PathsConfig holds path settings.
type PathsConfig struct {
	GTBinary string `toml:"gt_binary"`
	BDBinary string `toml:"bd_binary"`
	TownRoot string `toml:"town_root"`
}

// FiltersConfig holds default filter settings.
type FiltersConfig struct {
	Status     []string `toml:"status"`
	ShowClosed bool     `toml:"show_closed"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		RefreshInterval:    1 * time.Second,
		StuckThresholdMins: 30,
		LogLines:           10,
		ShowLogs:           true,
		Paths: PathsConfig{
			GTBinary: "gt",
			BDBinary: "bd",
			TownRoot: "", // Auto-detect
		},
		Filters: FiltersConfig{
			Status:     []string{"open", "in_progress"},
			ShowClosed: false,
		},
	}
}

// Load reads configuration from standard locations.
func Load() (*Config, error) {
	cfg := DefaultConfig()

	// Try config file locations in order
	paths := []string{
		filepath.Join(os.Getenv("HOME"), ".config", "gastop", "config.toml"),
	}

	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		paths = append([]string{filepath.Join(xdg, "gastop", "config.toml")}, paths...)
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			if _, err := toml.DecodeFile(p, cfg); err != nil {
				return nil, err
			}
			break
		}
	}

	// Auto-detect town root if not set
	if cfg.Paths.TownRoot == "" {
		cfg.Paths.TownRoot = detectTownRoot()
	}

	return cfg, nil
}

// detectTownRoot tries to find a Gas Town workspace by looking for markers.
func detectTownRoot() string {
	// Check GT_TOWN_ROOT environment variable first
	if root := os.Getenv("GT_TOWN_ROOT"); root != "" {
		return root
	}

	// Common locations to check
	home := os.Getenv("HOME")
	candidates := []string{
		filepath.Join(home, "gt", "gastop"),
		filepath.Join(home, "gt"),
		filepath.Join(home, "gastop"),
	}

	// Check current directory and parents
	cwd, err := os.Getwd()
	if err == nil {
		dir := cwd
		for i := 0; i < 5; i++ { // Check up to 5 levels up
			if isTownRoot(dir) {
				return dir
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	// Check common locations
	for _, c := range candidates {
		if isTownRoot(c) {
			return c
		}
	}

	return ""
}

// isTownRoot checks if a directory looks like a Gas Town workspace.
func isTownRoot(dir string) bool {
	// Look for .beads directory or config.json with rig markers
	beadsDir := filepath.Join(dir, ".beads")
	if info, err := os.Stat(beadsDir); err == nil && info.IsDir() {
		return true
	}

	// Check for mayor directory (towns have this)
	mayorDir := filepath.Join(dir, "mayor")
	if info, err := os.Stat(mayorDir); err == nil && info.IsDir() {
		return true
	}

	return false
}

// Save writes configuration to the default location.
func (c *Config) Save() error {
	dir := filepath.Join(os.Getenv("HOME"), ".config", "gastop")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path := filepath.Join(dir, "config.toml")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return toml.NewEncoder(f).Encode(c)
}
