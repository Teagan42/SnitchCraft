package interfaces

import (
	"errors"
	"testing"

	"github.com/teagan42/snitchcraft/internal/models"
)

// mockMetricsPlugin is a mock implementation of the MetricsPlugin interface for testing.
type mockMetricsPlugin struct {
	nameCalled  bool
	startCalled bool
	startErr    error
}

func (m *mockMetricsPlugin) Name() string {
	m.nameCalled = true
	return "mock"
}

func (m *mockMetricsPlugin) Start(ch chan models.RequestResult) error {
	m.startCalled = true
	return m.startErr
}

func TestMockMetricsPlugin_Name(t *testing.T) {
	mock := &mockMetricsPlugin{}
	got := mock.Name()
	if got != "mock" {
		t.Errorf("Name() = %q, want %q", got, "mock")
	}
	if !mock.nameCalled {
		t.Error("Name() did not set nameCalled to true")
	}
}

func TestMockMetricsPlugin_Start_Success(t *testing.T) {
	mock := &mockMetricsPlugin{}
	ch := make(chan models.RequestResult)
	go func() {
		close(ch)
	}()
	err := mock.Start(ch)
	if err != nil {
		t.Errorf("Start() returned error: %v, want nil", err)
	}
	if !mock.startCalled {
		t.Error("Start() did not set startCalled to true")
	}
}

func TestMockMetricsPlugin_Start_Error(t *testing.T) {
	mock := &mockMetricsPlugin{startErr: errors.New("fail")}
	ch := make(chan models.RequestResult)
	go func() {
		close(ch)
	}()
	err := mock.Start(ch)
	if err == nil {
		t.Error("Start() returned nil, want error")
	}
	if err.Error() != "fail" {
		t.Errorf("Start() error = %v, want %v", err, "fail")
	}
	if !mock.startCalled {
		t.Error("Start() did not set startCalled to true")
	}
}
