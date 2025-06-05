package models

type Config struct {
	BackendURL     string
	ParallelChecks bool
	LokiUrl        string
	PrometheusPort string
	ListenPort     string
	OTELMetricUrl  string
}
