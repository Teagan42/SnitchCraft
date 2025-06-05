package config

import (
	"os"
	"testing"

	"github.com/teagan42/snitchcraft/internal/models"
)

func TestGetEnvReturnsEnvVar(t *testing.T) {
	os.Setenv("TEST_ENV_VAR", "value")
	defer os.Unsetenv("TEST_ENV_VAR")

	val := getEnv("TEST_ENV_VAR", "default")
	if val != "value" {
		t.Errorf("expected 'value', got '%s'", val)
	}
}

func TestGetEnvReturnsDefault(t *testing.T) {
	os.Unsetenv("TEST_ENV_VAR")

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
	os.Setenv("BACKEND_URL", "http://localhost")
	defer os.Unsetenv("BACKEND_URL")
	os.Unsetenv("LISTEN_PORT")
	os.Unsetenv("PARALLEL_CHECKS")

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
