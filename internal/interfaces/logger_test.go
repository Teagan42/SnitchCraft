package interfaces

import (
	"errors"
	"testing"

	"github.com/teagan42/snitchcraft/internal/models"
)

// MockLogger implements the Logger interface for testing.
type MockLogger struct {
	name         string
	logCalled    bool
	startCalled  bool
	logArg       models.Loggable
	startChannel chan models.Loggable
	logErr       error
	startErr     error
}

func (m *MockLogger) Name() string {
	return m.name
}

func (m *MockLogger) Log(log models.Loggable) error {
	m.logCalled = true
	m.logArg = log
	return m.logErr
}

func (m *MockLogger) Start(logChannel chan models.Loggable) error {
	m.startCalled = true
	m.startChannel = logChannel
	return m.startErr
}

func TestMockLogger_Name(t *testing.T) {
	logger := &MockLogger{name: "testLogger"}
	if logger.Name() != "testLogger" {
		t.Errorf("expected Name to return 'testLogger', got '%s'", logger.Name())
	}
}

func TestMockLogger_Log(t *testing.T) {
	logger := &MockLogger{}
	log := models.RequestLogEntry{}
	err := logger.Log(log)
	if !logger.logCalled {
		t.Error("expected Log to set logCalled to true")
	}
	if logger.logArg.GetMethod() != log.GetMethod() ||
		logger.logArg.GetPath() != log.GetPath() ||
		logger.logArg.GetTime() != log.GetTime() {
		t.Error("expected Log to set logArg to the provided loggable")
	}

	if err != nil {
		t.Errorf("expected Log to return nil, got %v", err)
	}

	logger.logErr = errors.New("log error")
	err = logger.Log(log)
	if err == nil || err.Error() != "log error" {
		t.Errorf("expected Log to return 'log error', got %v", err)
	}
}

func TestMockLogger_Start(t *testing.T) {
	logger := &MockLogger{}
	ch := make(chan models.Loggable)
	err := logger.Start(ch)
	if !logger.startCalled {
		t.Error("expected Start to set startCalled to true")
	}
	if logger.startChannel != ch {
		t.Error("expected Start to set startChannel to the provided channel")
	}
	if err != nil {
		t.Errorf("expected Start to return nil, got %v", err)
	}

	logger.startErr = errors.New("start error")
	err = logger.Start(ch)
	if err == nil || err.Error() != "start error" {
		t.Errorf("expected Start to return 'start error', got %v", err)
	}
}
