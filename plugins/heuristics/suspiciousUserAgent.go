package heuristics

import (
	"net/http"
	"strings"

	"github.com/teagan42/snitchcraft/internal/interactors"
)

type SuspiciousUserAgent struct{}

func (u SuspiciousUserAgent) Name() string {
	return "Suspicious User-Agent"
}

func (u SuspiciousUserAgent) Check(r *http.Request) (string, bool) {
	ua := strings.ToLower(r.Header.Get("User-Agent"))
	if strings.Contains(ua, "curl") || strings.Contains(ua, "python") {
		return "Detected script-like User-Agent", true
	}
	return "", false
}

func init() {
	interactors.RegisterHeuristic(SuspiciousUserAgent{})
}
