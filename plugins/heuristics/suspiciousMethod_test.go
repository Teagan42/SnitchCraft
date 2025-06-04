package heuristics

import "testing"

func TestSuspiciousMethodCheck_Name(t *testing.T) {
	s := SuspiciousMethodCheck{}
	want := "Suspicious HTTP Method"
	if got := s.Name(); got != want {
		t.Errorf("Name() = %q, want %q", got, want)
	}
}
