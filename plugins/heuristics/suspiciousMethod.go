package heuristics

import (
	"fmt"
	"net/http"
)

type SuspiciousMethodCheck struct{}

func (s SuspiciousMethodCheck) Name() string {
	return "suspicious_http_method"
}

func (s SuspiciousMethodCheck) Check(r *http.Request) (string, bool) {
	switch r.Method {
	case "TRACE", "TRACK", "DEBUG", "CONNECT":
		return "Use of uncommon or dangerous HTTP method: " + r.Method, true
	default:
		return "", false
	}
}

func init() {
	fmt.Printf("[heuristics] registering SuspiciousMethodCheck heuristic...\n")
	RegisterHeuristic(SuspiciousMethodCheck{})
}
