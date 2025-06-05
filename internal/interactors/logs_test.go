package interactors

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/teagan42/snitchcraft/internal/interfaces"
	"github.com/teagan42/snitchcraft/internal/models"
	"github.com/teagan42/snitchcraft/plugins/loggers"
)

// mockLogger implements interfaces.Logger for testing.
type mockLogger struct {
	name         string
	startCalled  bool
	startErr     error
	startChannel chan models.Loggable
}

func (m *mockLogger) Name() string { return m.name }
func (m *mockLogger) Start(ch chan models.Loggable) error {
	m.startCalled = true
	m.startChannel = ch
	return m.startErr
}
func (m *mockLogger) Log(log models.Loggable) error {
	if m.startErr != nil {
		return m.startErr
	}
	// Simulate logging by just returning nil
	return nil
}

func TestLogWorker_StartsRegisteredLoggers(t *testing.T) {
	var mu sync.Mutex
	var started []string

	// Save original RegisteredLoggers and restore after test
	origRegistered := loggers.RegisteredLoggers
	defer func() { loggers.RegisteredLoggers = origRegistered }()

	mock1 := &mockLogger{name: "mock1"}
	mock2 := &mockLogger{name: "mock2"}

	loggers.RegisteredLoggers = []func(models.Config) interfaces.Logger{
		func(cfg models.Config) interfaces.Logger {
			mu.Lock()
			defer mu.Unlock()
			started = append(started, "mock1")
			return mock1
		},
		func(cfg models.Config) interfaces.Logger {
			mu.Lock()
			defer mu.Unlock()
			started = append(started, "mock2")
			return mock2
		},
	}

	cfg := models.Config{}
	LogWorker(cfg)

	time.Sleep(200 * time.Millisecond) // Allow goroutines to start

	if !mock1.startCalled {
		t.Errorf("mock1 logger Start was not called")
	}
	if !mock2.startCalled {
		t.Errorf("mock2 logger Start was not called")
	}
	if mock1.startChannel != logsChannel {
		t.Errorf("mock1 logger did not receive the correct channel")
	}
	if mock2.startChannel != logsChannel {
		t.Errorf("mock2 logger did not receive the correct channel")
	}
}

func TestLogWorker_LoggerStartError(t *testing.T) {
	origRegistered := loggers.RegisteredLoggers
	defer func() { loggers.RegisteredLoggers = origRegistered }()

	mockErr := errors.New("start error")
	mock := &mockLogger{name: "mockErr", startErr: mockErr}

	loggers.RegisteredLoggers = []func(models.Config) interfaces.Logger{
		func(cfg models.Config) interfaces.Logger { return mock },
	}

	cfg := models.Config{}
	LogWorker(cfg)

	time.Sleep(200 * time.Millisecond) // Allow goroutine to attempt start

	if !mock.startCalled {
		t.Errorf("mock logger Start was not called")
	}
}

func TestLogWorker_NilLogger(t *testing.T) {
	origRegistered := loggers.RegisteredLoggers
	defer func() { loggers.RegisteredLoggers = origRegistered }()

	loggers.RegisteredLoggers = []func(models.Config) interfaces.Logger{
		func(cfg models.Config) interfaces.Logger { return nil },
	}

	cfg := models.Config{}
	LogWorker(cfg)
	// Should not panic or call Start
}
