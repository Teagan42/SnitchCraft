package interactors

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"snitchcraft/internal/interfaces"
	"snitchcraft/internal/models"
	"snitchcraft/plugins/heuristics"
	"snitchcraft/plugins/loggers"
	"snitchcraft/plugins/metrics"
	"time"
)

var logSink interfaces.LogSink = &loggers.StdoutLogger{}
var metricsSink interfaces.MetricsSink = &metrics.PrometheusMetrics{}

func StartProxyServer() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		metricsSink.IncRequest()

		// Analyze heuristics
		reasons := []string{}
		for _, h := range heuristics.RegisteredHeuristics {
			if reason, bad := h.Check(r); bad {
				reasons = append(reasons, fmt.Sprintf("%s: %s", h.Name(), reason))
			}
		}

		if len(reasons) > 0 {
			metricsSink.IncMalicious()
		}

		entry := models.LogEntry{
			Time:      time.Now().UTC(),
			Method:    r.Method,
			Path:      r.URL.Path,
			Headers:   r.Header.Clone(),
			Malicious: len(reasons) > 0,
			Reasons:   reasons,
		}
		_ = logSink.Send(entry)

		resp := map[string]any{"status": "ok", "malicious": entry.Malicious, "reasons": reasons}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	fmt.Println("Proxy running at :8080")
	return http.ListenAndServe(":8080", nil)
}
