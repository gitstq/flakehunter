// Package engine stores run history and computes flake statistics.
// The engine is pure standard library and fully offline.
package engine

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds local FlakeHunter settings stored at .flakehunter/config.json.
type Config struct {
	MaxRuns              int     `json:"max_runs"`
	MinRuns              int     `json:"min_runs"`
	MaxFlakeRate         float64 `json:"max_flake_rate"`
	MaxFlaky             int     `json:"max_flaky"`
	Framework            string  `json:"framework"`
	QuarantineStableRuns int     `json:"quarantine_stable_runs"`
}

// DefaultConfig returns the built-in defaults.
func DefaultConfig() Config {
	return Config{
		MaxRuns:              100,
		MinRuns:              3,
		MaxFlakeRate:         5,
		MaxFlaky:             0,
		Framework:            "pytest",
		QuarantineStableRuns: 5,
	}
}

// LoadConfig reads config.json from dir, falling back to defaults when absent.
func LoadConfig(dir string) (Config, error) {
	cfg := DefaultConfig()
	data, err := os.ReadFile(filepath.Join(dir, "config.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	if cfg.MaxRuns <= 0 {
		cfg.MaxRuns = 100
	}
	if cfg.MinRuns < 2 {
		cfg.MinRuns = 2
	}
	if cfg.QuarantineStableRuns <= 0 {
		cfg.QuarantineStableRuns = 5
	}
	return cfg, nil
}

// SaveConfig writes the config file, creating dir when needed.
func SaveConfig(dir string, cfg Config) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "config.json"), append(data, '\n'), 0o644)
}
