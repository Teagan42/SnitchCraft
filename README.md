![Go](https://github.com/Teagan42/SnitchCraft/actions/workflows/go_tests.yml/badge.svg)

# SnitchCraft

A modular HTTP reverse proxy that uses plugin-based heuristics for security inspection, with pluggable logging and metrics backends.

## Features
- Reverse proxy handler
- Plugin-based heuristic analyzers
- Plugin-based logging
- Plugin-based metrics
- Interface-based logging + metrics
- OTEL-ready, Loki-compatible

## Layout
- `cmd/` — app entrypoint
- `internal/` — interfaces and shared models
- `plugins/` — actual plugin implementations

## Included Heuristics
- Suspicious User-Agent (e.g., curl, python)
- SQL Injection detection
- Anomalous headers (missing or fake headers)
- Suspicious HTTP methods (TRACE, TRACK, CONNECT, etc.)

## Usage
```bash
go mod tidy
go run cmd/server/main.go
```

## Configuration
Change the backend URL in `main.go`:
```go
StartProxy(":8080", "http://localhost:8081")
```

## Logging Output Example
```json
{
  "time": "2025-06-03T15:20:01Z",
  "method": "GET",
  "path": "/api/user?id=1 OR 1=1",
  "headers": {
    "User-Agent": ["curl/7.68.0"]
  },
  "malicious": true,
  "reasons": [
    "Suspicious User-Agent: Detected script-like User-Agent",
    "SQL Injection: Query contains possible SQL injection"
  ]
}
```

## Run Tests
```bash
go test ./... -cover -v
```