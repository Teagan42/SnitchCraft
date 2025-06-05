//go:build !integration
// +build !integration

package interactors

import (
	"errors"
	"testing"
	"time"

	"github.com/teagan42/snitchcraft/internal/interfaces"
	"github.com/teagan42/snitchcraft/internal/models"
	"github.com/teagan42/snitchcraft/plugins/metrics"
)

type mockMetricsPlugin struct {
	name      string
	startErr  error
	started   bool
	startChan chan models.RequestResult
}

func (m *mockMetricsPlugin) Name() string {
	return m.name
}
func (m *mockMetricsPlugin) Start(resultsChannel chan models.RequestResult) error {
	m.started = true
	m.startChan = resultsChannel
	return m.startErr
}

func TestMetricsWorker_StartsRegisteredPlugins(t *testing.T) {
	origPlugins := metrics.RegisteredMetricsPlugins
	defer func() { metrics.RegisteredMetricsPlugins = origPlugins }()

	mockPlugin := &mockMetricsPlugin{name: "mock", startErr: nil}
	metrics.RegisteredMetricsPlugins = []func(models.Config) interfaces.MetricsPlugin{
		func(cfg models.Config) interfaces.MetricsPlugin { return mockPlugin },
	}

	resultsChan := make(chan models.RequestResult)
	cfg := models.Config{}

	MetricsWorker(cfg, resultsChan)
	time.Sleep(300 * time.Millisecond)

	if !mockPlugin.started {
		t.Error("plugin Start was not called")
	}
	if mockPlugin.startChan != resultsChan {
		t.Error("plugin Start did not receive correct resultsChannel")
	}
}

func TestMetricsWorker_PluginReturnsNil(t *testing.T) {
	origPlugins := metrics.RegisteredMetricsPlugins
	defer func() { metrics.RegisteredMetricsPlugins = origPlugins }()

	metrics.RegisteredMetricsPlugins = []func(models.Config) interfaces.MetricsPlugin{
		func(cfg models.Config) interfaces.MetricsPlugin { return nil },
	}

	resultsChan := make(chan models.RequestResult)
	cfg := models.Config{}

	MetricsWorker(cfg, resultsChan)
	// No panic or error expected
}

func TestMetricsWorker_PluginStartReturnsError(t *testing.T) {
	origPlugins := metrics.RegisteredMetricsPlugins
	defer func() { metrics.RegisteredMetricsPlugins = origPlugins }()

	mockPlugin := &mockMetricsPlugin{name: "mock", startErr: errors.New("fail")}
	metrics.RegisteredMetricsPlugins = []func(models.Config) interfaces.MetricsPlugin{
		func(cfg models.Config) interfaces.MetricsPlugin { return mockPlugin },
	}

	resultsChan := make(chan models.RequestResult)
	cfg := models.Config{}

	MetricsWorker(cfg, resultsChan)
	// No panic or error expected, error is printed
}

func TestMetricsWorker_LogsConfig(t *testing.T) {
	// This test just ensures that the function runs and prints config.
	cfg := models.Config{}
	resultsChan := make(chan models.RequestResult)
	MetricsWorker(cfg, resultsChan)
}
