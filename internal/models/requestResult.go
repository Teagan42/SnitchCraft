package models

import "net/http"

type RequestResult struct {
	HeuristicResults []HeuristicResult `json:"heuristic_results"`
	Request          *http.Request     `json:"request"`
	Duration         uint64            `json:"duration"`
}
