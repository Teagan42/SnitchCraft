package config

import (
	"fmt"
	"os"

	"github.com/teagan42/snitchcraft/internal/models"
)

func Load() (models.Config, error) {
	return Validate(models.Config{
		BackendURL:     getEnv("BACKEND_URL", ""),
		ParallelChecks: getEnv("PARALLEL_CHECKS", "true") == "true",
		LogForwardURL:  getEnv("LOG_FORWARD_URL", ""),
		MetricsPort:    getEnv("METRICS_PORT", "9090"),
		ListenPort:     getEnv("LISTEN_PORT", ":8080"),
		OTELExporter:   getEnv("OTEL_EXPORTER", "stdout"),
	})
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func Validate(cfg models.Config) (models.Config, error) {
	if cfg.BackendURL == "" {
		return models.Config{}, fmt.Errorf("missing required env var: BACKEND_URL")
	}

	return cfg, nil
}
