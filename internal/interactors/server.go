package interactors

import (
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/teagan42/snitchcraft/internal/models"
	"github.com/teagan42/snitchcraft/utils"

	"maps"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

var resultChan = make(chan models.RequestResult, 10)

func StartProxyServer(cfg models.Config) error {
	setupTracing()

	go MetricsWorker(cfg)
	go LogWorker(cfg)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("X-Trace-ID", uuid.New().String())
		logsChannel <- models.RequestLogEntry{
			Time:    time.Now(),
			Method:  r.Method,
			Path:    r.URL.Path,
			Headers: r.Header,
			TraceID: r.Header.Get("X-Trace-ID"),
		}
		var results []models.HeuristicResult = RunHeuristicChecks(r, cfg)

		// pass to logging via channel
		resultChan <- models.RequestResult{Request: r, HeuristicResults: results}

		if len(utils.Map(results, func(r models.HeuristicResult) bool {
			return r.Issue != ""
		})) > 0 {
			// if any heuristic found an issue, return 403
			w.WriteHeader(http.StatusForbidden)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error": "Forbidden"}`))
			logsChannel <- models.ResponseLogEntry{
				Time:             time.Now(),
				Method:           r.Method,
				Path:             r.URL.Path,
				Headers:          r.Header,
				TraceID:          r.Header.Get("X-Trace-ID"),
				Malicious:        true,
				HeuristicResults: results,
				StatusCode:       http.StatusForbidden,
			}
			return
		}

		// forward request
		req, _ := http.NewRequest(r.Method, cfg.BackendURL+r.URL.RequestURI(), r.Body)
		req.Header = r.Header.Clone()

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			logsChannel <- models.ResponseLogEntry{
				Time:             time.Now(),
				Method:           r.Method,
				Path:             r.URL.Path,
				Headers:          r.Header,
				TraceID:          r.Header.Get("X-Trace-ID"),
				Malicious:        false,
				HeuristicResults: results,
				StatusCode:       http.StatusBadGateway,
			}
			http.Error(w, "upstream error", http.StatusBadGateway)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error": "Bad Gateway"}`))
			return
		}
		defer resp.Body.Close()
		logsChannel <- models.ResponseLogEntry{
			Time:             time.Now(),
			Method:           r.Method,
			Path:             r.URL.Path,
			Headers:          r.Header,
			TraceID:          r.Header.Get("X-Trace-ID"),
			Malicious:        false,
			HeuristicResults: results,
			StatusCode:       resp.StatusCode,
		}
		maps.Copy(w.Header(), resp.Header)
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	return http.ListenAndServe(cfg.ListenPort, nil)
}

func setupTracing() {
	exp, _ := stdouttrace.New(stdouttrace.WithPrettyPrint())
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(resource.Empty()),
	)
	otel.SetTracerProvider(tp)
}
