package heuristics

import (
	"fmt"
	"net/http"
	"strings"
)

type AnomalousHeaderCheck struct{}

func (h AnomalousHeaderCheck) Name() string {
	return "anomalous_header"
}

func (h AnomalousHeaderCheck) Check(r *http.Request) (string, bool) {
	if r.Header.Get("Accept") == "" {
		return "Missing Accept header", true
	}
	// Disabled for now, as it can cause issues with legitimate requests
	// if r.Header.Get("Host") == "" {
	// 	return "Missing Host header", true
	// }
	ref := strings.ToLower(r.Header.Get("Referer"))
	if ref != "" && strings.Contains(ref, "evil.com") {
		return "Referer contains suspicious domain", true
	}
	return "", false
}

func init() {
	fmt.Println("[heuristics] registering AnomalousHeaderCheck heuristic...")
	RegisterHeuristic(AnomalousHeaderCheck{})
}
