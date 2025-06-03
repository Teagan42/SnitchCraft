package interactors

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"snitchcraft/internal/config/env"
	"snitchcraft/internal/interfaces"
	"snitchcraft/internal/models"
	"snitchcraft/plugins/heuristics"
	"snitchcraft/plugins/loggers"
	"snitchcraft/plugins/metrics"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

type ResultMsg struct {
	Request *http.Request
	Reasons []string
}

var logSink interfaces.LogSink = &loggers.StdoutLogger{}
var metricsSink interfaces.MetricsSink = &metrics.PrometheusMetrics{}
var resultChan = make(chan ResultMsg, 10)

func StartProxyServer(cfg env.Config) error {
	setupTracing()

	go resultWorker()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var reasons []string
		for _, h := range heuristics.RegisteredHeuristics {
			if reason, bad := h.Check(r); bad {
				reasons = append(reasons, fmt.Sprintf("%s: %s", h.Name(), reason))
			}
		}

		metricsSink.IncRequest()
		if len(reasons) > 0 {
			metricsSink.IncMalicious()
		}

		// pass to logging via channel
		resultChan <- ResultMsg{Request: r, Reasons: reasons}

		// forward request
		req, _ := http.NewRequest(r.Method, cfg.BackendURL+r.URL.RequestURI(), r.Body)
		req.Header = r.Header.Clone()

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "upstream error", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		for k, v := range resp.Header {
			w.Header()[k] = v
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	return http.ListenAndServe(cfg.ListenPort, nil)
}

func resultWorker() {
	for msg := range resultChan {
		e := models.LogEntry{
			Time:      time.Now().UTC(),
			Method:    msg.Request.Method,
			Path:      msg.Request.URL.RequestURI(),
			Headers:   msg.Request.Header.Clone(),
			Malicious: len(msg.Reasons) > 0,
			Reasons:   msg.Reasons,
		}
		_ = logSink.Send(e)
	}
}

func setupTracing() {
	exp, _ := stdouttrace.New(stdouttrace.WithPrettyPrint())
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(resource.Empty()),
	)
	otel.SetTracerProvider(tp)
}
