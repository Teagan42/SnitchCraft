package heuristics

import (
	"net/http/httptest"
	"testing"
)

func TestSuspiciousUserAgent(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", "python-requests")

	check := SuspiciousUserAgent{}
	reason, isMalicious := check.Check(req)
	if !isMalicious {
		t.Fatal("expected request to be flagged as malicious")
	}
	if reason == "" {
		t.Fatal("expected non-empty reason")
	}
}

func TestSQLInjectionCheck(t *testing.T) {
	req := httptest.NewRequest("GET", "/?q=1+OR+1%3D1", nil)

	check := SQLInjectionCheck{}
	reason, isMalicious := check.Check(req)
	if !isMalicious {
		t.Fatal("expected SQLi to be detected")
	}
	if reason == "" {
		t.Fatal("expected non-empty reason")
	}
}

func TestAnomalousHeaderCheck(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Referer", "http://evil.com")

	check := AnomalousHeaderCheck{}
	reason, isMalicious := check.Check(req)
	if !isMalicious {
		t.Fatal("expected header anomaly to be detected")
	}
	if reason == "" {
		t.Fatal("expected non-empty reason")
	}
}

func TestSuspiciousMethodCheck(t *testing.T) {
	req := httptest.NewRequest("TRACE", "/", nil)

	check := SuspiciousMethodCheck{}
	reason, isMalicious := check.Check(req)
	if !isMalicious {
		t.Fatal("expected suspicious method to be detected")
	}
	if reason == "" {
		t.Fatal("expected non-empty reason")
	}
}
func TestAsyncRunHeuristicChecks_SingleHeuristic(t *testing.T) {
	// Mock heuristic
	type mockHeuristic struct{}