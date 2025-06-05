package heuristics

import "testing"

func TestSuspiciousMethodCheck_Name(t *testing.T) {
	s := SuspiciousMethodCheck{}
	want := "suspicious_http_method"
	if got := s.Name(); got != want {
		t.Errorf("Name() = %q, want %q", got, want)
	}
}
