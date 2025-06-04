package interactors

import (
	"net/http"
	"testing"

	"github.com/teagan42/snitchcraft/internal/interfaces"
	"github.com/teagan42/snitchcraft/internal/models"
)

type mockHeuristic struct {
	name      string
	checkName string
	checkOK   bool
}

func (m *mockHeuristic) Name() string {
	return m.name
}

func (m *mockHeuristic) Check(req *http.Request) (string, bool) {
	return m.checkName, m.checkOK
}

func setupHeuristics(hs ...interfaces.Heuristic) {
	RegisteredHeuristics = hs
}

func TestRunHeuristicChecks_Sync(t *testing.T) {
	setupHeuristics(&mockHeuristic{name: "TestHeuristic", checkName: "Issue1", checkOK: true})
	cfg := models.Config{ParallelChecks: false}
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	results := RunHeuristicChecks(req, cfg)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Name != "TestHeuristic" || results[0].Issue != "Issue1" {
		t.Errorf("unexpected result: %+v", results[0])
	}
}

func TestRunHeuristicChecks_Parallel(t *testing.T) {
	setupHeuristics(
		&mockHeuristic{name: "H1", checkName: "I1", checkOK: true},
		&mockHeuristic{name: "H2", checkName: "I2", checkOK: true},
	)
	cfg := models.Config{ParallelChecks: true}
	req, _ := http.NewRequest("POST", "http://example.com", nil)

	results := RunHeuristicChecks(req, cfg)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	names := map[string]bool{}
	issues := map[string]bool{}
	for _, r := range results {
		names[r.Name] = true
		issues[r.Issue] = true
	}
	if !names["H1"] || !names["H2"] || !issues["I1"] || !issues["I2"] {
		t.Errorf("unexpected results: %+v", results)
	}
}

func TestRunHeuristicChecks_NoHeuristics(t *testing.T) {
	setupHeuristics()
	cfg := models.Config{ParallelChecks: false}
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	results := RunHeuristicChecks(req, cfg)

	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestRunHeuristicChecks_HeuristicReturnsFalse(t *testing.T) {
	setupHeuristics(&mockHeuristic{name: "TestHeuristic", checkName: "Issue1", checkOK: false})
	cfg := models.Config{ParallelChecks: false}
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	results := RunHeuristicChecks(req, cfg)

	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
