package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"handfree-work/web-restic/internal/config"
)

func TestLoadDevConfigAndModeOverride(t *testing.T) {
	configDir := t.TempDir()
	writeConfig(t, configDir, "dev.yaml", `
name: web-restic
mode: dev
server:
  port: ":3000"
database:
  path: "./data/db.sqlite"
`)
	writeConfig(t, configDir, "prod.yaml", `
mode: prod
server:
  port: ":8080"
`)

	cfg, err := config.Load(configDir, "prod")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Mode != "prod" {
		t.Fatalf("Mode = %q, want prod", cfg.Mode)
	}
	if cfg.Server.Port != ":8080" {
		t.Fatalf("Server.Port = %q, want :8080", cfg.Server.Port)
	}
	if cfg.Database.Path != "./data/db.sqlite" {
		t.Fatalf("Database.Path = %q, want ./data/db.sqlite", cfg.Database.Path)
	}
}

func TestLoadUsesDevConfigWhenModeIsDev(t *testing.T) {
	configDir := t.TempDir()
	writeConfig(t, configDir, "dev.yaml", "mode: dev\nserver:\n  port: ':3000'\n")

	cfg, err := config.Load(configDir, "dev")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Server.Port != ":3000" {
		t.Fatalf("Server.Port = %q, want :3000", cfg.Server.Port)
	}
}

func writeConfig(t *testing.T, directory, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(directory, name), []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
}
