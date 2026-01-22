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

	return cfg, nil
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
