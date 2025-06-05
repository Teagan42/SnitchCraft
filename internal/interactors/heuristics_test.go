package interactors

import (
	"net/http"
	"testing"

	"github.com/teagan42/snitchcraft/internal/interfaces"
	"github.com/teagan42/snitchcraft/internal/models"
	"github.com/teagan42/snitchcraft/plugins/heuristics"
)

func TestRunHeuristicChecks_Sync(t *testing.T) {
	origHeuristics := heuristics.RegisteredHeuristics
	defer func() { heuristics.RegisteredHeuristics = origHeuristics }()

	mockHeuristic := &mockHeuristic{
		name: "TestHeuristic",
		checkFunc: func(req *http.Request) (string, bool) {
			return "TestIssue", true
		},
	}
	heuristics.RegisteredHeuristics = []interfaces.Heuristic{mockHeuristic}

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	cfg := models.Config{ParallelChecks: false}

	results := RunHeuristicChecks(req, cfg)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Name != "TestHeuristic" || results[0].Issue != "TestIssue" {
		t.Errorf("unexpected result: %+v %+v", results[0], results)
	}
}

func TestRunHeuristicChecks_Async(t *testing.T) {
	origHeuristics := heuristics.RegisteredHeuristics
	defer func() { heuristics.RegisteredHeuristics = origHeuristics }()

	mockHeuristic := &mockHeuristic{
		name: "AsyncHeuristic",
		checkFunc: func(req *http.Request) (string, bool) {
			return "AsyncIssue", true
		},
	}
	heuristics.RegisteredHeuristics = []interfaces.Heuristic{mockHeuristic}

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	cfg := models.Config{ParallelChecks: true}

	results := RunHeuristicChecks(req, cfg)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Name != "AsyncHeuristic" || results[0].Issue != "AsyncIssue" {
		t.Errorf("unexpected result: %+v", results[0])
	}
}

type mockHeuristic struct {
	name      string
	checkFunc func(req *http.Request) (string, bool)
}

func (m *mockHeuristic) Name() string {
	return m.name
}

func (m *mockHeuristic) Check(req *http.Request) (string, bool) {
	return m.checkFunc(req)
}
