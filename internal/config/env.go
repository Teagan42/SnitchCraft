package config

import (
	"fmt"
	"os"

	"github.com/teagan42/snitchcraft/internal/models"
)

func Load() (models.Config, error) {
	return Validate(models.Config{
		ListenPort:     getEnv("LISTEN_PORT", ":8080"),
		BackendURL:     getEnv("BACKEND_URL", ""),
		ParallelChecks: getEnv("PARALLEL_CHECKS", "true") == "true",
		LokiUrl:        getEnv("LOKI_URL", ""),
		PrometheusPort: getEnv("PROMETHEUS_PORT", ""),
		OTELMetricUrl:  getEnv("OTEL_METRIC_URL", ""),
		LogFile:        getEnv("LOG_FILE", ""),
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
