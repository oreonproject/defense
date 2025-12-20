package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if !cfg.General.RealTimeProtection {
		t.Error("expected real_time_protection to be true by default")
	}
	if cfg.General.LogLevel != "info" {
		t.Errorf("expected log_level 'info', got %q", cfg.General.LogLevel)
	}
	if !cfg.Firewall.Enabled {
		t.Error("expected firewall to be enabled by default")
	}
	if cfg.Notifications.Level != "all" {
		t.Errorf("expected notification level 'all', got %q", cfg.Notifications.Level)
	}
}

func TestLoadMissing(t *testing.T) {
	cfg, err := Load("/nonexistent/path/config.toml")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if cfg == nil {
		t.Fatal("expected default config, got nil")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.toml")

	cfg := Default()
	cfg.General.LogLevel = "debug"
	cfg.Notifications.Level = "critical"
	cfg.Scanning.Exclusions = []string{"/tmp", "/var/cache"}

	if err := cfg.Save(path); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if loaded.General.LogLevel != "debug" {
		t.Errorf("expected log_level 'debug', got %q", loaded.General.LogLevel)
	}
	if loaded.Notifications.Level != "critical" {
		t.Errorf("expected notification level 'critical', got %q", loaded.Notifications.Level)
	}
	if len(loaded.Scanning.Exclusions) != 2 {
		t.Errorf("expected 2 exclusions, got %d", len(loaded.Scanning.Exclusions))
	}
}

func TestSaveCreatesDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "dir", "config.toml")

	cfg := Default()
	if err := cfg.Save(path); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("config file not created: %v", err)
	}
}
