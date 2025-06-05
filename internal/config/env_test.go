package config

import (
	"os"
	"testing"

	"github.com/teagan42/snitchcraft/internal/models"
)

func TestGetEnvReturnsEnvVar(t *testing.T) {
	if err := os.Setenv("TEST_ENV_VAR", "value"); err != nil {
		t.Fatalf("failed to set env var: %v", err)
	}
	defer func() {
		if err := os.Unsetenv("TEST_ENV_VAR"); err != nil {
			t.Fatalf("failed to unset env var: %v", err)
		}
	}()

	val := getEnv("TEST_ENV_VAR", "default")
	if val != "value" {
		t.Errorf("expected 'value', got '%s'", val)
	}
}

func TestGetEnvReturnsDefault(t *testing.T) {
	if err := os.Unsetenv("TEST_ENV_VAR"); err != nil {
		t.Fatalf("failed to unset env var: %v", err)
	}

	val := getEnv("TEST_ENV_VAR", "default")
	if val != "default" {
		t.Errorf("expected 'default', got '%s'", val)
	}
}

func TestValidateReturnsErrorIfBackendURLMissing(t *testing.T) {
	cfg := models.Config{}
	_, err := Validate(cfg)
	if err == nil {
		t.Error("expected error for missing BACKEND_URL, got nil")
	}
}

func TestValidateReturnsConfigIfBackendURLPresent(t *testing.T) {
	cfg := models.Config{BackendURL: "http://localhost"}
	got, err := Validate(cfg)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got.BackendURL != "http://localhost" {
		t.Errorf("expected BackendURL to be 'http://localhost', got '%s'", got.BackendURL)
	}
}

func TestLoadReturnsConfigWithDefaults(t *testing.T) {
	if err := os.Setenv("BACKEND_URL", "http://localhost"); err != nil {
		t.Fatalf("failed to set BACKEND_URL: %v", err)
	}
	defer func() {
		if err := os.Unsetenv("BACKEND_URL"); err != nil {
			t.Fatalf("failed to unset BACKEND_URL: %v", err)
		}
	}()
	if err := os.Unsetenv("LISTEN_PORT"); err != nil {
		t.Fatalf("failed to unset LISTEN_PORT: %v", err)
	}
	if err := os.Unsetenv("PARALLEL_CHECKS"); err != nil {
		t.Fatalf("failed to unset PARALLEL_CHECKS: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ListenPort != ":8080" {
		t.Errorf("expected ListenPort ':8080', got '%s'", cfg.ListenPort)
	}
	if cfg.BackendURL != "http://localhost" {
		t.Errorf("expected BackendURL 'http://localhost', got '%s'", cfg.BackendURL)
	}
	if !cfg.ParallelChecks {
		t.Errorf("expected ParallelChecks true, got false")
	}
}
