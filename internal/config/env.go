import (
    "os"
    "strconv"
)


func Load() (Config, error) {
    return Validate(Config{
        BackendURL:    getEnv("BACKEND_URL", "http://localhost:8081"),
        LogForwardURL: getEnv("LOG_FORWARD_URL", ""),
        MetricsPort:   getEnv("METRICS_PORT", "9090"),
        ListenPort:    getEnv("LISTEN_PORT", ":8080"),
        OTELExporter:  getEnv("OTEL_EXPORTER", "stdout"),
    })
}

func getEnv(key, def string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return def
}

func Validate(cfg Config) (Config, error) {
    if cfg.BackendURL == "" {
        return nil, fmt.Errorf("issing required env var: BACKEND_URL")
    }

	return cfg, nil
}