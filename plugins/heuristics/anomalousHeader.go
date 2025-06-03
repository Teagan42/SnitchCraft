package heuristics

import (
    "net/http"
    "strings"
)

type AnomalousHeaderCheck struct{}

func (h AnomalousHeaderCheck) Name() string {
    return "Header Anomaly"
}

func (h AnomalousHeaderCheck) Check(r *http.Request) (string, bool) {
    if r.Header.Get("Accept") == "" {
        return "Missing Accept header", true
    }
    if r.Header.Get("Host") == "" {
        return "Missing Host header", true
    }
    ref := strings.ToLower(r.Header.Get("Referer"))
    if ref != "" && strings.Contains(ref, "evil.com") {
        return "Referer contains suspicious domain", true
    }
    return "", false
}

func init() {
    RegisterCheck(AnomalousHeaderCheck{})
}