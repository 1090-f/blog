package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSupportsEnvWithoutConfigFile(t *testing.T) {
	t.Setenv("BLOG_CONFIG_FILE", filepath.Join(t.TempDir(), "missing.yaml"))
	t.Setenv("BLOG_SERVER_PORT", "9090")
	t.Setenv("BLOG_DATABASE_HOST", "127.0.0.1")
	t.Setenv("BLOG_DATABASE_PORT", "3306")
	t.Setenv("BLOG_DATABASE_USER", "root")
	t.Setenv("BLOG_DATABASE_PASS", "secret")
	t.Setenv("BLOG_DATABASE_NAME", "blog")
	t.Setenv("BLOG_JWT_SECRET", "token-secret")

	_, err := Load()
	if err == nil {
		t.Fatal("expected missing explicit config file to return an error")
	}
}

func TestLoadUsesDefaultsAndEnvironmentOverrides(t *testing.T) {
	t.Setenv("BLOG_CONFIG_FILE", "")
	t.Setenv("BLOG_SERVER_PORT", "9090")
	t.Setenv("BLOG_DATABASE_HOST", "127.0.0.1")
	t.Setenv("BLOG_DATABASE_PORT", "3306")
	t.Setenv("BLOG_DATABASE_USER", "root")
	t.Setenv("BLOG_DATABASE_PASS", "secret")
	t.Setenv("BLOG_DATABASE_NAME", "blog")
	t.Setenv("BLOG_JWT_SECRET", "token-secret")
	t.Setenv("BLOG_ADMIN_BOOTSTRAP_ENABLED", "true")
	t.Setenv("BLOG_ADMIN_BOOTSTRAP_USERNAME", "admin")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config failed: %v", err)
	}

	if cfg.Server.Port != "9090" {
		t.Fatalf("expected port 9090, got %q", cfg.Server.Port)
	}
	if cfg.JWT.Expire != 7200 {
		t.Fatalf("expected default jwt expire 7200, got %d", cfg.JWT.Expire)
	}
	if cfg.Upload.MaxSizeBytes != 5*1024*1024 {
		t.Fatalf("expected default upload max size, got %d", cfg.Upload.MaxSizeBytes)
	}
	if !cfg.AdminBootstrap.Enabled {
		t.Fatal("expected admin bootstrap enabled from env")
	}
	if cfg.AdminBootstrap.Username != "admin" {
		t.Fatalf("expected admin username from env, got %q", cfg.AdminBootstrap.Username)
	}
}

func TestLoadAllowsEnvToOverrideConfigFile(t *testing.T) {
	dir := t.TempDir()
	configFile := filepath.Join(dir, "config.yaml")
	content := []byte("server:\n  port: \"8080\"\njwt:\n  secret: \"file-secret\"\n  expire: 1800\n")
	if err := os.WriteFile(configFile, content, 0o600); err != nil {
		t.Fatalf("write config file failed: %v", err)
	}

	t.Setenv("BLOG_CONFIG_FILE", configFile)
	t.Setenv("BLOG_SERVER_PORT", "9090")
	t.Setenv("BLOG_JWT_SECRET", "env-secret")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config failed: %v", err)
	}

	if cfg.Server.Port != "9090" {
		t.Fatalf("expected env port override, got %q", cfg.Server.Port)
	}
	if cfg.JWT.Secret != "env-secret" {
		t.Fatalf("expected env secret override, got %q", cfg.JWT.Secret)
	}
	if cfg.JWT.Expire != 1800 {
		t.Fatalf("expected file expire preserved, got %d", cfg.JWT.Expire)
	}
}
