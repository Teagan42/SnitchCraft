package models

type Config struct {
	BackendURL     string
	ParallelChecks bool
	LogForwardURL  string
	MetricsPort    string
	ListenPort     string
	OTELExporter   string
}
