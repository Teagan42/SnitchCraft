package metrics

import (
	"fmt"

	"github.com/teagan42/snitchcraft/internal/interfaces"
	"github.com/teagan42/snitchcraft/internal/models"
)

type OtelMetrics struct{}

func (m *OtelMetrics) UpdateRequestStats(stats models.RequestStats) {
	// This method is not used in this implementation
}

func (m *OtelMetrics) UpdateHeuristicStats(stats map[string]uint64) {
	// This method is not used in this implementation
}

func init() {
	fmt.Println("[metrics] OtelMetrics initialized")
}

var _ interfaces.MetricsSink = (*OtelMetrics)(nil)
