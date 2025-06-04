package models

import (
	"sync/atomic"

	"github.com/teagan42/snitchcraft/utils"
)

type RequestStatsCounter struct {
	TotalRequests  atomic.Uint64
	TotalMalicious atomic.Uint64
	ReturnCodes    utils.ConcurrentCounter
}

type RequestStats struct {
	TotalRequests  uint64            `json:"total_requests"`
	TotalMalicious uint64            `json:"total_malicious"`
	ReturnCodes    map[string]uint64 `json:"return_codes"`
}

func (rs *RequestStats) ToCounter() RequestStatsCounter {
	return RequestStatsCounter{
		TotalRequests:  atomic.Uint64{},
		TotalMalicious: atomic.Uint64{},
		ReturnCodes:    utils.ConcurrentCounter{},
	}
}

func (rs *RequestStatsCounter) ToModel() RequestStats {
	return RequestStats{
		TotalRequests:  rs.TotalRequests.Load(),
		TotalMalicious: rs.TotalMalicious.Load(),
		ReturnCodes:    rs.ReturnCodes.Snapshot(),
	}
}
