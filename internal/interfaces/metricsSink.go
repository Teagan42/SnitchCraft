package interfaces

import "github.com/teagan42/snitchcraft/internal/models"

type MetricsSink interface {
	UpdateHeuristicStats(stats map[string]uint64)
	UpdateRequestStats(stats models.RequestStats)
}
