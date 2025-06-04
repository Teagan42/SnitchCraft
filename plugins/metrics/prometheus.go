package metrics

import (
	"fmt"

	"github.com/teagan42/snitchcraft/internal/interfaces"
	"github.com/teagan42/snitchcraft/internal/models"
)

type PrometheusMetrics struct{}

func (m *PrometheusMetrics) UpdateRequestStats(stats models.RequestStats) {
	// This method is not used in this implementation
}

func (m *PrometheusMetrics) UpdateHeuristicStats(stats map[string]uint64) {
	// This method is not used in this implementation
}

func init() {
	fmt.Println("[metrics] PrometheusMetrics initialized")
}

var _ interfaces.MetricsSink = (*PrometheusMetrics)(nil)
