package heuristics

import (
	"net/http"
	"testing"

	"github.com/teagan42/snitchcraft/internal/interfaces"
	"github.com/teagan42/snitchcraft/internal/models"
	"github.com/teagan42/snitchcraft/plugins/loggers"
)

type mockHeuristic struct {
	name string
}

func (m *mockHeuristic) Name() string {
	return m.name
}

func (m *mockHeuristic) Check(req *http.Request) (string, bool) {
	return "", false // No actual check, just a mock
}

func TestRegisterHeuristic_AppendsToRegisteredHeuristics(t *testing.T) {
	// Save and restore original state
	orig := RegisteredHeuristics
	defer func() { RegisteredHeuristics = orig }()

	RegisteredHeuristics = nil

	h := &mockHeuristic{name: "test"}
	RegisterHeuristic(h)

	if len(RegisteredHeuristics) != 1 {
		t.Fatalf("expected 1 heuristic, got %d", len(RegisteredHeuristics))
	}
	if RegisteredHeuristics[0] != h {
		t.Error("registered heuristic does not match the one added")
	}
}

func TestRegisterHeuristic_MultipleCalls(t *testing.T) {
	orig := RegisteredHeuristics
	defer func() { RegisteredHeuristics = orig }()

	RegisteredHeuristics = nil

	h1 := &mockHeuristic{name: "h1"}
	h2 := &mockHeuristic{name: "h2"}

	RegisterHeuristic(h1)
	RegisterHeuristic(h2)

	if len(RegisteredHeuristics) != 2 {
		t.Fatalf("expected 2 heuristics, got %d", len(RegisteredHeuristics))
	}
	if RegisteredHeuristics[0] != h1 || RegisteredHeuristics[1] != h2 {
		t.Error("registered heuristics are not in the correct order")
	}
}
func TestRegisterLogger_AppendsToRegisteredLoggers(t *testing.T) {
	orig := loggers.RegisteredLoggers
	defer func() { loggers.RegisteredLoggers = orig }()

	loggers.RegisteredLoggers = nil

	mockLogger := func(cfg models.Config) interfaces.Logger {
		return nil
	}

	loggers.RegisterLogger(mockLogger)

	if len(loggers.RegisteredLoggers) != 1 {
		t.Fatalf("expected 1 logger, got %d", len(loggers.RegisteredLoggers))
	}
	if loggers.RegisteredLoggers[0] == nil {
		t.Error("registered logger is nil")
	}
}

func TestRegisterLogger_MultipleCalls(t *testing.T) {
	orig := loggers.RegisteredLoggers
	defer func() { loggers.RegisteredLoggers = orig }()

	loggers.RegisteredLoggers = nil

	mockLogger1 := func(cfg models.Config) interfaces.Logger { return nil }
	mockLogger2 := func(cfg models.Config) interfaces.Logger { return nil }

	loggers.RegisterLogger(mockLogger1)
	loggers.RegisterLogger(mockLogger2)

	if len(loggers.RegisteredLoggers) != 2 {
		t.Fatalf("expected 2 loggers, got %d", len(loggers.RegisteredLoggers))
	}
	if loggers.RegisteredLoggers[0] == nil || loggers.RegisteredLoggers[1] == nil {
		t.Error("one or more registered loggers are nil")
	}
}
