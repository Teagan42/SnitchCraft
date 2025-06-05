package metrics

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/teagan42/snitchcraft/internal/interfaces"
	"github.com/teagan42/snitchcraft/internal/models"
	"github.com/teagan42/snitchcraft/plugins/heuristics"
)

func UnregisterPrometheusMetrics(pm *PrometheusMetrics) {
	prometheus.Unregister(pm.requestCount)
	prometheus.Unregister(pm.heuristicMatches)
	prometheus.Unregister(pm.requestTime)
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	prometheus.DefaultGatherer = prometheus.DefaultRegisterer.(prometheus.Gatherer)
}

func TestNewPrometheusMetrics_NoPort(t *testing.T) {
	cfg := models.Config{PrometheusPort: ""}
	metrics := NewPrometheusMetrics(cfg)
	if metrics != nil {
		UnregisterPrometheusMetrics(metrics.(*PrometheusMetrics))
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

type testHeuristic struct {
	NameField string
	CheckFunc func(req *http.Request) (string, bool)
}

func (th *testHeuristic) Name() string {
	return th.NameField
}
func (th *testHeuristic) Check(req *http.Request) (string, bool) {
	return th.CheckFunc(req)
}

var h1 = &testHeuristic{
	NameField: "h1",
	CheckFunc: func(req *http.Request) (string, bool) {
		return "issue1", true
	},
}
var h2 = &testHeuristic{
	NameField: "h2",
	CheckFunc: func(req *http.Request) (string, bool) {
		return "", false
	},
}

func TestPrometheusMetrics_Start_CountsRequestsAndHeuristics(t *testing.T) {
	originalHeuristics := heuristics.RegisteredHeuristics
	heuristics.RegisteredHeuristics = []interfaces.Heuristic{
		h1,
		h2,
	}
	cfg := models.Config{PrometheusPort: "9090"}
	pm := NewPrometheusMetrics(cfg).(*PrometheusMetrics)
	defer func() {
		UnregisterPrometheusMetrics(pm)
		heuristics.RegisteredHeuristics = originalHeuristics
	}()
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
			{Name: "h2", Issue: "issue3"},
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
	if err := pm.heuristicMatches.With(prometheus.Labels{
		"heuristic_h1": "true",
		"heuristic_h2": "false",
	}).Write(metric); err != nil {
		t.Errorf("failed to get heuristic metric: %v", err)
	}
	if metric.GetCounter().GetValue() != 1 {
		t.Errorf("expected h1 count 2, got %v", metric.GetCounter().GetValue())
	}
	if err := pm.heuristicMatches.With(prometheus.Labels{
		"heuristic_h1": "false",
		"heuristic_h2": "false",
	}).Write(metric); err != nil {
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
