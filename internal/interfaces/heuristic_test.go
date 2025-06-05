package interfaces

import (
	"net/http"
	"testing"
)

type mockHeuristic struct{}

func (m mockHeuristic) Name() string {
	return "mock"
}
func (m mockHeuristic) Check(r *http.Request) (string, bool) {
	return "checked", true
}
func TestHeuristicInterfaceImplementation(t *testing.T) {

	var h interface{} = mockHeuristic{}
	if _, ok := h.(interface {
		Name() string
		Check(*http.Request) (string, bool)
	}); !ok {
		t.Errorf("mockHeuristic does not implement Heuristic interface")
	}
}

func TestHeuristicCheckMethod(t *testing.T) {
	h := mockHeuristic{}
	issue, ok := h.Check(nil)
	if !ok {
		t.Errorf("expected ok to be true, got false")
	}
	if issue != "checked" {
		t.Errorf("expected issue to be 'checked', got %q", issue)
	}
}
