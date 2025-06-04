package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	// Unset all relevant env vars to test defaults
	if err := os.Unsetenv("BACKEND_URL"); err != nil {
		t.Fatalf("Failed to unset BACKEND_URL: %v", err)
	}
	if err := os.Unsetenv("PARALLEL_CHECKS"); err != nil {
		t.Fatalf("Failed to unset PARALLEL_CHECKS: %v", err)
	}
	if err := os.Unsetenv("LOG_FORWARD_URL"); err != nil {
		t.Fatalf("Failed to unset LOG_FORWARD_URL: %v", err)
	}
	if err := os.Unsetenv("METRICS_PORT"); err != nil {
		t.Fatalf("Failed to unset METRICS_PORT: %v", err)
	}
	if err := os.Unsetenv("LISTEN_PORT"); err != nil {
		t.Fatalf("Failed to unset LISTEN_PORT: %v", err)
	}
	if err := os.Unsetenv("OTEL_EXPORTER"); err != nil {
		t.Fatalf("Failed to unset OTEL_EXPORTER: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.BackendURL != "http://localhost:8081" {
		t.Errorf("BackendURL = %q, want %q", cfg.BackendURL, "http://localhost:8081")
	}
	if cfg.ParallelChecks != true {
		t.Errorf("ParallelChecks = %v, want true", cfg.ParallelChecks)
	}
	if cfg.LogForwardURL != "" {
		t.Errorf("LogForwardURL = %q, want empty string", cfg.LogForwardURL)
	}
	if cfg.MetricsPort != "9090" {
		t.Errorf("MetricsPort = %q, want %q", cfg.MetricsPort, "9090")
	}
	if cfg.ListenPort != ":8080" {
		t.Errorf("ListenPort = %q, want %q", cfg.ListenPort, ":8080")
	}
	if cfg.OTELExporter != "stdout" {
		t.Errorf("OTELExporter = %q, want %q", cfg.OTELExporter, "stdout")
	}
}

func TestLoad_WithEnvVars(t *testing.T) {
	if err := os.Setenv("BACKEND_URL", "http://example.com"); err != nil {
		t.Fatalf("Failed to set BACKEND_URL: %v", err)
	}
	if err := os.Setenv("PARALLEL_CHECKS", "false"); err != nil {
		t.Fatalf("Failed to set PARALLEL_CHECKS: %v", err)
	}
	if err := os.Setenv("LOG_FORWARD_URL", "http://log.example.com"); err != nil {
		t.Fatalf("Failed to set LOG_FORWARD_URL: %v", err)
	}
	if err := os.Setenv("METRICS_PORT", "1234"); err != nil {
		t.Fatalf("Failed to set METRICS_PORT: %v", err)
	}
	if err := os.Setenv("LISTEN_PORT", ":9999"); err != nil {
		t.Fatalf("Failed to set LISTEN_PORT: %v", err)
	}
	if err := os.Setenv("OTEL_EXPORTER", "otlp"); err != nil {
		t.Fatalf("Failed to set OTEL_EXPORTER: %v", err)
	}

	defer func() {
		os.Unsetenv("BACKEND_URL")
		os.Unsetenv("PARALLEL_CHECKS")
		os.Unsetenv("LOG_FORWARD_URL")
		os.Unsetenv("METRICS_PORT")
		os.Unsetenv("LISTEN_PORT")
		os.Unsetenv("OTEL_EXPORTER")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.BackendURL != "http://example.com" {
		t.Errorf("BackendURL = %q, want %q", cfg.BackendURL, "http://example.com")
	}
	if cfg.ParallelChecks != false {
		t.Errorf("ParallelChecks = %v, want false", cfg.ParallelChecks)
	}
	if cfg.LogForwardURL != "http://log.example.com" {
		t.Errorf("LogForwardURL = %q, want %q", cfg.LogForwardURL, "http://log.example.com")
	}
	if cfg.MetricsPort != "1234" {
		t.Errorf("MetricsPort = %q, want %q", cfg.MetricsPort, "1234")
	}
	if cfg.ListenPort != ":9999" {
		t.Errorf("ListenPort = %q, want %q", cfg.ListenPort, ":9999")
	}
	if cfg.OTELExporter != "otlp" {
		t.Errorf("OTELExporter = %q, want %q", cfg.OTELExporter, "otlp")
	}
}

func TestLoad_MissingBackendURL(t *testing.T) {
	os.Setenv("BACKEND_URL", "")
	defer os.Unsetenv("BACKEND_URL")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() expected error for missing BACKEND_URL, got nil")
	}
}
