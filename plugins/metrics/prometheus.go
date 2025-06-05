package metrics

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/teagan42/snitchcraft/internal/interfaces"
	"github.com/teagan42/snitchcraft/internal/models"
	"github.com/teagan42/snitchcraft/plugins/heuristics"
	"github.com/teagan42/snitchcraft/utils"
)

type PrometheusMetrics struct {
	requestCount     *prometheus.CounterVec
	heuristicMatches *prometheus.CounterVec
	requestTime      *prometheus.HistogramVec
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
			utils.Map(heuristics.RegisteredHeuristics, func(hr interfaces.Heuristic) string {
				return fmt.Sprintf("heuristic_%s", strings.ReplaceAll(hr.Name(), " ", "_"))
			}),
		),
		requestTime: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "snitchcraft_http_request_duration_nanoseconds",
				Help:    "Duration of HTTP request processing prior to proxy in nanoseconds",
				Buckets: prometheus.DefBuckets,
			},
			append(
				[]string{"method", "path"},
				utils.Map(heuristics.RegisteredHeuristics, func(hr interfaces.Heuristic) string {
					return fmt.Sprintf("heuristic_%s", strings.ReplaceAll(hr.Name(), " ", "_"))
				})...,
			),
		),
	}

	prometheus.MustRegister(pm.requestCount, pm.heuristicMatches, pm.requestTime)
	ServePrometheusMetrics(cfg.PrometheusPort)
	fmt.Println("[metrics] PrometheusMetrics initialized")
	return pm
}

func ServePrometheusMetrics(port string) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	go func() {
		addr := fmt.Sprintf(":%s", port)
		fmt.Printf("[metrics] Prometheus metrics available at %s/metrics\n", addr)
		if err := http.ListenAndServe(addr, mux); err != nil {
			fmt.Printf("[metrics] Prometheus HTTP server error: %v\n", err)
		}
	}()
}

func (m *PrometheusMetrics) Name() string {
	return "PrometheusMetrics"
}

func (m *PrometheusMetrics) Start(resultChannel chan models.RequestResult) error {
	go func() {
		for result := range resultChannel {
			if result.Request == nil {
				fmt.Println("[metrics] Received result with nil request, skipping")
				continue
			}
			m.requestCount.WithLabelValues(result.Request.Method, result.Request.URL.Path).Inc()
			var labels = map[string]string{
				"method": result.Request.Method,
				"path":   result.Request.URL.Path,
			}
			utils.Do(heuristics.RegisteredHeuristics, func(hr interfaces.Heuristic) {
				labels[fmt.Sprintf("heuristic_%s", strings.ReplaceAll(hr.Name(), " ", "_"))] = "false"
			})
			utils.Do(result.HeuristicResults, func(r models.HeuristicResult) {
				labels[fmt.Sprintf("heuristic_%s", strings.ReplaceAll(r.Name, " ", "_"))] = fmt.Sprintf("%t", r.Issue != "")
			})
			m.requestTime.With(labels).Observe(float64(result.Duration))

			labels = map[string]string{}
			utils.Do(heuristics.RegisteredHeuristics, func(h interfaces.Heuristic) {
				labels[fmt.Sprintf("heuristic_%s", strings.ReplaceAll(h.Name(), " ", "_"))] = "false"
			})
			utils.Do(result.HeuristicResults, func(hr models.HeuristicResult) {
				labels[fmt.Sprintf("heuristic_%s", strings.ReplaceAll(hr.Name, " ", "_"))] = fmt.Sprintf("%t", hr.Issue != "")
			})
			m.heuristicMatches.With(labels).Inc()
		}
	}()
	fmt.Println("[metrics] PrometheusMetrics started")
	return nil
}

func init() {
	RegisterMetricsPlugin(NewPrometheusMetrics)
	fmt.Println("[metrics] Prometheus metrics plugin registered")
}
