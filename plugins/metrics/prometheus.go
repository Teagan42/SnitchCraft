package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/teagan42/snitchcraft/internal/interfaces"
	"github.com/teagan42/snitchcraft/internal/models"
	"github.com/teagan42/snitchcraft/utils"
)

type PrometheusMetrics struct {
	requestCount     *prometheus.CounterVec
	heuristicMatches *prometheus.CounterVec
}

func NewPrometheusMetrics(cfg models.Config) interfaces.MetricsPlugin {
	if cfg.PrometheusPort == "" {
		fmt.Println("[metrics] Prometheus port not configured, skipping initialization")
		return nil
	}
	pm := &PrometheusMetrics{
		requestCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "snitchcraft_http_requests_total",
				Help: "Total number of HTTP requests received",
			},
			[]string{"method", "path"},
		),
		heuristicMatches: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "snitchcraft_heuristic_matches_total",
				Help: "Total number of heuristic matches by type",
			},
			[]string{"heuristic"},
		),
	}

	prometheus.MustRegister(pm.requestCount, pm.heuristicMatches)
	fmt.Println("[metrics] PrometheusMetrics initialized")
	return pm
}

func (m *PrometheusMetrics) Name() string {
	return "PrometheusMetrics"
}

func (m *PrometheusMetrics) Start(resultChannel chan models.RequestResult) error {
	go func() {
		for {
			select {
			case result := <-resultChannel:
				if result.Request == nil {
					fmt.Println("[metrics] Received result with nil request, skipping")
					continue
				}
				m.requestCount.WithLabelValues(result.Request.Method, result.Request.URL.Path).Inc()
				for _, heuristicResult := range utils.Filter(result.HeuristicResults, func(r models.HeuristicResult) bool {
					return r.Issue != ""
				}) {
					m.heuristicMatches.WithLabelValues(heuristicResult.Name).Inc()
				}
			default:
				// No results to process, sleep briefly to avoid busy waiting
			}
		}
	}()

	return nil
}

func init() {
	RegisterMetricsPlugin(NewPrometheusMetrics)
	fmt.Println("[metrics] Prometheus metrics plugin registered")
}
