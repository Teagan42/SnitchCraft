package metrics

import (
	"testing"

	"github.com/teagan42/snitchcraft/internal/interfaces"
	"github.com/teagan42/snitchcraft/internal/models"
)

func TestRegisterMetricsPlugin_AppendsToRegisteredMetricsPlugins(t *testing.T) {
	// Save and restore original state
	orig := RegisteredMetricsPlugins
	defer func() { RegisteredMetricsPlugins = orig }()

	RegisteredMetricsPlugins = nil

	mockMetricsPlugin := func(cfg models.Config) interfaces.MetricsPlugin {
		return nil
	}

	RegisterMetricsPlugin(mockMetricsPlugin)

	if len(RegisteredMetricsPlugins) != 1 {
		t.Fatalf("expected 1 metrics plugin, got %d", len(RegisteredMetricsPlugins))
	}
	if RegisteredMetricsPlugins[0] == nil {
		t.Error("registered metrics plugin is nil")
	}
}

func TestRegisterMetricsPlugin_MultipleCalls(t *testing.T) {
	orig := RegisteredMetricsPlugins
	defer func() { RegisteredMetricsPlugins = orig }()

	RegisteredMetricsPlugins = nil

	mockMetricsPlugin1 := func(cfg models.Config) interfaces.MetricsPlugin { return nil }
	mockMetricsPlugin2 := func(cfg models.Config) interfaces.MetricsPlugin { return nil }

	RegisterMetricsPlugin(mockMetricsPlugin1)
	RegisterMetricsPlugin(mockMetricsPlugin2)

	if len(RegisteredMetricsPlugins) != 2 {
		t.Fatalf("expected 2 metrics plugins, got %d", len(RegisteredMetricsPlugins))
	}
	if RegisteredMetricsPlugins[0] == nil || RegisteredMetricsPlugins[1] == nil {
		t.Error("one or more registered metrics plugins are nil")
	}
}
