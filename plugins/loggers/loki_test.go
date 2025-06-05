package loggers

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/teagan42/snitchcraft/internal/interfaces"
	"github.com/teagan42/snitchcraft/internal/models"
)

// mockLoggable implements models.Loggable for testing
type mockLoggable struct {
	Time time.Time
	Msg  string
}

func (m *mockLoggable) GetTime() time.Time { return m.Time }
func (m *mockLoggable) GetTraceID() string { return "test-trace-id" }
func (m *mockLoggable) GetMethod() string  { return "GET" }
func (m *mockLoggable) GetPath() string    { return "/test/path" }
func (m *mockLoggable) GetHeaders() http.Header {
	return http.Header{"X-Test-Header": []string{"test-value"}}
}
func (m *mockLoggable) GetType() string { return "test-loggable" }

func TestLokiLogger_Name(t *testing.T) {
	logger := &LokiLogger{}
	if logger.Name() != "LokiLogger" {
		t.Errorf("expected Name to be LokiLogger, got %s", logger.Name())
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newTestClient(fn roundTripFunc) *http.Client {
	return &http.Client{Transport: fn}
}

func TestLokiLogger_Log_Success(t *testing.T) {
	called := false
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		called = true
		if req.Method != "POST" {
			t.Errorf("expected POST, got %s", req.Method)
		}
		body, _ := io.ReadAll(req.Body)
		if !bytes.Contains(body, []byte("test-message")) {
			t.Errorf("expected body to contain log message")
		}
		return &http.Response{
			StatusCode: 204,
			Body:       io.NopCloser(bytes.NewBuffer(nil)),
		}, nil
	})

	logger := &LokiLogger{
		client:  client,
		lokiURL: "http://loki",
		labels:  map[string]string{"job": "snitchcraft"},
	}

	log := &mockLoggable{Time: time.Now(), Msg: "test-message"}
	err := logger.Log(log)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !called {
		t.Error("expected HTTP client to be called")
	}
}

func TestLokiLogger_Log_RequestError(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("fail")
	})
	logger := &LokiLogger{
		client:  client,
		lokiURL: "http://loki",
		labels:  map[string]string{"job": "snitchcraft"},
	}
	log := &mockLoggable{Time: time.Now(), Msg: "test-message"}
	err := logger.Log(log)
	if err == nil || err.Error() == "" {
		t.Error("expected error from HTTP client")
	}
}

func TestLokiLogger_Log_Non2xxStatus(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 400,
			Status:     "400 Bad Request",
			Body:       io.NopCloser(bytes.NewBuffer(nil)),
		}, nil
	})
	logger := &LokiLogger{
		client:  client,
		lokiURL: "http://loki",
		labels:  map[string]string{"job": "snitchcraft"},
	}
	log := &mockLoggable{Time: time.Now(), Msg: "test-message"}
	err := logger.Log(log)
	if err == nil || err.Error() == "" {
		t.Error("expected error for non-2xx status")
	}
}

func TestLokiLogger_Start_LogsError(t *testing.T) {
	logger := &LokiLogger{
		client: newTestClient(func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("fail")
		}),
		lokiURL: "http://loki",
		labels:  map[string]string{"job": "snitchcraft"},
	}
	logChan := make(chan models.Loggable, 1)
	logChan <- &mockLoggable{Time: time.Now(), Msg: "fail"}
	close(logChan)
	_ = logger.Start(logChan)
	// No panic expected, error is printed
}

func TestNewLokiLogger_NoURL(t *testing.T) {
	cfg := models.Config{LokiUrl: ""}
	logger := NewLokiLogger(cfg)
	if logger != nil {
		t.Error("expected nil logger when no LokiUrl")
	}
}

func TestNewLokiLogger_WithURL(t *testing.T) {
	cfg := models.Config{LokiUrl: "http://loki"}
	logger := NewLokiLogger(cfg)
	if logger == nil {
		t.Error("expected logger instance")
	}
	if logger.Name() != "LokiLogger" {
		t.Errorf("expected Name LokiLogger, got %s", logger.Name())
	}
}

var _ interfaces.Logger = &LokiLogger{}
