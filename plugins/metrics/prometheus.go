package metrics

import (
	"fmt"
	"sync/atomic"
	"snitchcraft/internal/interfaces"
)

type PrometheusMetrics struct{}

var totalReq uint64
var totalMalicious uint64

func (m *PrometheusMetrics) IncRequest() {
	atomic.AddUint64(&totalReq, 1)
}

func (m *PrometheusMetrics) IncMalicious() {
	atomic.AddUint64(&totalMalicious, 1)
}

func init() {
	fmt.Println("[metrics] PrometheusMetrics initialized")
}

var _ interfaces.MetricsSink = (*PrometheusMetrics)(nil)
