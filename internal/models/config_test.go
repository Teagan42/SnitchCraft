package models

import (
	"reflect"
	"testing"
)

func TestConfig_DefaultValues(t *testing.T) {
	cfg := Config{}
	if cfg.BackendURL != "" {
		t.Errorf("expected BackendURL to be empty, got %q", cfg.BackendURL)
	}
	if cfg.ParallelChecks != false {
		t.Errorf("expected ParallelChecks to be false, got %v", cfg.ParallelChecks)
	}
	if cfg.LokiUrl != "" {
		t.Errorf("expected LokiUrl to be empty, got %q", cfg.LokiUrl)
	}
	if cfg.PrometheusPort != "" {
		t.Errorf("expected PrometheusPort to be empty, got %q", cfg.PrometheusPort)
	}
	if cfg.ListenPort != "" {
		t.Errorf("expected ListenPort to be empty, got %q", cfg.ListenPort)
	}
	if cfg.OTELMetricUrl != "" {
		t.Errorf("expected OTELMetricUrl to be empty, got %q", cfg.OTELMetricUrl)
	}
}

func TestConfig_FieldAssignment(t *testing.T) {
	want := Config{
		BackendURL:     "http://backend",
		ParallelChecks: true,
		LokiUrl:        "http://loki",
		PrometheusPort: "9090",
		ListenPort:     "8080",
		OTELMetricUrl:  "http://otel",
	}
	got := Config{
		BackendURL:     "http://backend",
		ParallelChecks: true,
		LokiUrl:        "http://loki",
		PrometheusPort: "9090",
		ListenPort:     "8080",
		OTELMetricUrl:  "http://otel",
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected %+v, got %+v", want, got)
	}
}
