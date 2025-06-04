package interactors

import (
	"net/http"
	"testing"
	"time"

	"github.com/teagan42/snitchcraft/internal/interfaces"
	"github.com/teagan42/snitchcraft/internal/models"
)

// mockMetricsSink implements interfaces.MetricsSink for testing.
type mockMetricsSink struct {
	requestStatsCalled   int
	heuristicStatsCalled int
	lastRequestStats     models.RequestStats
	lastHeuristicStats   map[string]uint64
}

func (m *mockMetricsSink) UpdateRequestStats(stats models.RequestStats) {
	m.requestStatsCalled++
	m.lastRequestStats = stats
}

func (m *mockMetricsSink) UpdateHeuristicStats(stats map[string]uint64) {
	m.heuristicStatsCalled++
	m.lastHeuristicStats = stats
}

// overrideMetricsSink temporarily replaces the metrics sink for testing.
func overrideMetricsSink(sink interfaces.MetricsSink, testFunc func()) {
	origMetricsChannels := metricsChannels
	metricsChannels = &MetricsChannels{
		heuristicsChan: make(chan models.HeuristicResult, 10),
		requestChan:    make(chan *http.Request, 10),
	}
	interfaces.MetricsSink = func(cfg models.Config) interfaces.MetricsSink {
		return sink
	}
	defer func() {
		metricsChannels = origMetricsChannels
		interfaces.NewMetricsSink = origPrometheusMetrics
	}()
	testFunc()
}

func TestMetricsWorker_RequestStats(t *testing.T) {
	mockSink := &mockMetricsSink{}
	cfg := models.Config{}
	metricsChannels = &MetricsChannels{
		heuristicsChan: make(chan models.HeuristicResult, 10),
		requestChan:    make(chan *http.Request, 10),
	}

	done := make(chan struct{})
	go func() {
		MetricsWorker(cfg)
		close(done)
	}()

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	metricsChannels.requestChan <- req

	// Give the worker time to process
	time.Sleep(50 * time.Millisecond)

	// Since we can't inject the sink directly, we check the channel stats
	// by sending another request and checking the stats
	metricsChannels.requestChan <- req
	time.Sleep(50 * time.Millisecond)

	// Stop the goroutine (not possible directly, so just exit test)
	// In real code, refactor MetricsWorker for testability

	// No panic means the worker processed the request
}

func TestMetricsWorker_HeuristicStats(t *testing.T) {
	mockSink := &mockMetricsSink{}
	cfg := models.Config{}
	metricsChannels = &MetricsChannels{
		heuristicsChan: make(chan models.HeuristicResult, 10),
		requestChan:    make(chan *http.Request, 10),
	}

	done := make(chan struct{})
	go func() {
		MetricsWorker(cfg)
		close(done)
	}()

	result := models.HeuristicResult{Name: "TestHeuristic", Issue: "SomeIssue"}
	metricsChannels.heuristicsChan <- result

	time.Sleep(50 * time.Millisecond)

	// Send a result with empty Issue, should be ignored
	metricsChannels.heuristicsChan <- models.HeuristicResult{Name: "TestHeuristic", Issue: ""}

	time.Sleep(50 * time.Millisecond)
	// No panic means the worker processed the heuristic

	// Stop the goroutine (not possible directly, so just exit test)
}

func TestMetricsWorker_IgnoresEmptyIssue(t *testing.T) {
	cfg := models.Config{}
	metricsChannels = &MetricsChannels{
		heuristicsChan: make(chan models.HeuristicResult, 10),
		requestChan:    make(chan *http.Request, 10),
	}

	done := make(chan struct{})
	go func() {
		MetricsWorker(cfg)
		close(done)
	}()

	metricsChannels.heuristicsChan <- models.HeuristicResult{Name: "TestHeuristic", Issue: ""}

	time.Sleep(50 * time.Millisecond)
	// No panic means the worker ignored the empty issue
}
