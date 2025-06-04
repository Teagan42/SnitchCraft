package interactors

import (
	"net/http"
	"sync/atomic"

	"github.com/teagan42/snitchcraft/internal/interfaces"
	"github.com/teagan42/snitchcraft/internal/models"
	"github.com/teagan42/snitchcraft/plugins/metrics"
	"github.com/teagan42/snitchcraft/utils"
)

type MetricsChannels struct {
	heuristicsChan chan models.HeuristicResult
	requestChan    chan *http.Request
}

var metricsChannels = &MetricsChannels{
	heuristicsChan: make(chan models.HeuristicResult, 100),
	requestChan:    make(chan *http.Request, 100),
}

func MetricsWorker(cfg models.Config) {
	var heuristicCounter utils.ConcurrentCounter = utils.ConcurrentCounter{}
	var requestStats models.RequestStatsCounter = models.RequestStatsCounter{
		TotalRequests:  atomic.Uint64{},
		TotalMalicious: atomic.Uint64{},
		ReturnCodes:    utils.ConcurrentCounter{},
	}
	var metricsSink interfaces.MetricsSink
	if cfg.OTELExporter != "" {
		metricsSink = &metrics.PrometheusMetrics{}
	} else {
		metricsSink = &metrics.PrometheusMetrics{}
	}

	for {
		select {
		case req := <-metricsChannels.requestChan:
			requestStats.TotalRequests.Add(1)
			requestStats.ReturnCodes.Inc(req.Method)
			metricsSink.UpdateRequestStats(requestStats.ToModel())
		case result := <-metricsChannels.heuristicsChan:
			if result.Issue == "" {
				continue // No issue, nothing to do
			}
			heuristicCounter.Inc(result.Name)
			metricsSink.UpdateHeuristicStats(heuristicCounter.Snapshot())
		default:
			// No messages, just continue
			continue
		}
	}
}
