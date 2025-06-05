package metrics

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/teagan42/snitchcraft/internal/models"
)

func UnregisterPrometheusMetrics(pm *PrometheusMetrics) {
	prometheus.Unregister(pm.requestCount)
	prometheus.Unregister(pm.heuristicMatches)
}

func TestNewPrometheusMetrics_NoPort(t *testing.T) {
	cfg := models.Config{PrometheusPort: ""}
	metrics := NewPrometheusMetrics(cfg)
	if metrics != nil {
		t.Error("expected nil when PrometheusPort is empty")
	}
}

func TestNewPrometheusMetrics_WithPort(t *testing.T) {
	cfg := models.Config{PrometheusPort: "9090"}
	pm := NewPrometheusMetrics(cfg).(*PrometheusMetrics)

	defer UnregisterPrometheusMetrics(pm)

	if pm == nil {
		t.Fatal("expected PrometheusMetrics instance")
	}
	if pm.Name() != "PrometheusMetrics" {
		t.Errorf("expected Name PrometheusMetrics, got %s", pm.Name())
	}
}

func TestPrometheusMetrics_Start_CountsRequestsAndHeuristics(t *testing.T) {
	cfg := models.Config{PrometheusPort: "9090"}
	pm := NewPrometheusMetrics(cfg).(*PrometheusMetrics)
	defer UnregisterPrometheusMetrics(pm)
	resultChan := make(chan models.RequestResult, 2)

	req1 := models.RequestResult{
		Request: &http.Request{
			Method: "GET",
			URL:    &url.URL{Path: "/foo"},
		},
		HeuristicResults: []models.HeuristicResult{
			{Name: "h1", Issue: "issue1"},
			{Name: "h2", Issue: ""},
		},
	}
	req2 := models.RequestResult{
		Request: &http.Request{
			Method: "POST",
			URL:    &url.URL{Path: "/bar"},
		},
		HeuristicResults: []models.HeuristicResult{
			{Name: "h1", Issue: "issue2"},
		},
	}
	resultChan <- req1
	resultChan <- req2
	close(resultChan)

	_ = pm.Start(resultChan)
	time.Sleep(100 * time.Millisecond) // Let goroutine process

	// Check requestCount
	metric := &dto.Metric{}
	if err := pm.requestCount.WithLabelValues("GET", "/foo").Write(metric); err != nil {
		t.Errorf("failed to get metric: %v", err)
	}
	if metric.GetCounter().GetValue() != 1 {
		t.Errorf("expected GET /foo count 1, got %v", metric.GetCounter().GetValue())
	}
	if err := pm.requestCount.WithLabelValues("POST", "/bar").Write(metric); err != nil {
		t.Errorf("failed to get metric: %v", err)
	}
	if metric.GetCounter().GetValue() != 1 {
		t.Errorf("expected POST /bar count 1, got %v", metric.GetCounter().GetValue())
	}

	// Check heuristicMatches
	if err := pm.heuristicMatches.WithLabelValues("h1").Write(metric); err != nil {
		t.Errorf("failed to get heuristic metric: %v", err)
	}
	if metric.GetCounter().GetValue() != 2 {
		t.Errorf("expected h1 count 2, got %v", metric.GetCounter().GetValue())
	}
	if err := pm.heuristicMatches.WithLabelValues("h2").Write(metric); err != nil {
		t.Errorf("failed to get heuristic metric: %v", err)
	}
	if metric.GetCounter().GetValue() != 0 {
		t.Errorf("expected h2 count 0, got %v", metric.GetCounter().GetValue())
	}
}

func TestPrometheusMetrics_Name(t *testing.T) {
	cfg := models.Config{PrometheusPort: "9090"}
	pm := NewPrometheusMetrics(cfg).(*PrometheusMetrics)
	defer UnregisterPrometheusMetrics(pm)
	if pm.Name() != "PrometheusMetrics" {
		t.Errorf("expected PrometheusMetrics, got %s", pm.Name())
	}
}
