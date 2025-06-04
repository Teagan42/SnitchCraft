![Go](https://github.com/Teagan42/SnitchCraft/actions/workflows/go_tests.yml/badge.svg)

Certainly! Based on the recent enhancements in the feat/heuristic-channel branch of the SnitchCraft project, here’s an updated README.md that reflects the new architecture and features:

⸻

# SnitchCraft

SnitchCraft is a modular, pluggable HTTP proxy designed for security analysis and observability. It intercepts HTTP traffic, applies heuristic checks to detect potential threats, and forwards the results to various logging and metrics backends.

## 🚀 Features
	•	Heuristic Analysis: Intercepts HTTP requests and applies a series of heuristic checks to identify suspicious activities.
	•	Concurrent Processing: Heuristic checks are executed in parallel using Go’s goroutines and channels, ensuring high performance.
	•	Pluggable Architecture: Supports dynamic loading of plugins for heuristics, loggers, and metrics collectors.
	•	Observability: Integrates with OpenTelemetry for tracing and Prometheus for metrics collection.
	•	Configurable: Easily configurable via environment variables or configuration files.
	•	Dockerized: Comes with Docker and Docker Compose support for easy deployment.

## 📁 Project Structure

.
├── cmd/
│   └── server/           # Entry point for the proxy server
├── internal/
│   ├── config/           # Configuration loading and validation
│   ├── interactors/      # Core business logic and orchestration
│   ├── interfaces/       # Interface definitions for plugins
│   └── models/           # Data models used across the application
├── plugins/
│   ├── heuristics/       # Heuristic check implementations
│   ├── loggers/          # Logging implementations
│   └── metrics/          # Metrics collection implementations
├── Dockerfile            # Docker image definition
├── docker-compose.yml    # Docker Compose configuration
└── README.md             # Project documentation

## ⚙️ Configuration

Configuration is handled via environment variables:
	•	BACKEND_URL: The URL of the backend server to which requests are forwarded.
	•	LISTEN_PORT: The port on which the proxy server listens (default: :8080).
	•	METRICS_PORT: The port for exposing Prometheus metrics (default: :9090).
	•	OTEL_EXPORTER: The OpenTelemetry exporter to use (stdout, jaeger, etc.).

You can also provide a .env file for local development.

## 👟 Running

### 🛡️ Running the Proxy

```shell
go mod tidy
go run cmd/server/main.go
```

### 🧪 Running Tests

To run the test suite:

```shell
go test ./...
```

This will execute all unit and integration tests across the project.

## 🐳 Docker Deployment

To build and run the application using Docker Compose:

```shell
docker-compose up --build
```

This will start the proxy server along with any configured services like Jaeger and Prometheus.

## 📈 Observability
	•	Metrics: Exposed at /metrics on the configured METRICS_PORT.
	•	Tracing: Integrated with OpenTelemetry; configure OTEL_EXPORTER to your preferred backend.

## 🔌 Plugin Development

To add a new plugin:
	1.	Define the Interface: In internal/interfaces/, define the interface your plugin will implement.
	2.	Implement the Plugin: Create your plugin in the appropriate plugins/ subdirectory.
	3.	Register the Plugin: Ensure your plugin is registered during application initialization.

This modular approach allows for easy extension and customization of SnitchCraft’s capabilities.
