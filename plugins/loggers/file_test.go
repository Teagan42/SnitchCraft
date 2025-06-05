package loggers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/teagan42/snitchcraft/internal/models"
)

func TestFileLogger_Name(t *testing.T) {
	fl := &FileLogger{}
	if got := fl.Name(); got != "FileLogger" {
		t.Errorf("Name() = %q, want %q", got, "FileLogger")
	}
}

func TestFileLogger_Log_Success(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	fl := &FileLogger{Writer: writer}

	entry := models.RequestLogEntry{
		Time:    time.Now(),
		Method:  "GET",
		Path:    "/test/path",
		Headers: http.Header{"X-Test-Header": []string{"test-value"}},
		TraceID: "test-trace-id",
	}
	err := fl.Log(entry)
	if err != nil {
		t.Fatalf("Log() error = %v, want nil", err)
	}
	writer.Flush()

	// Check output contains marshaled JSON
	var out map[string]interface{}
	line := buf.String()
	line = strings.TrimPrefix(line, "[loggers][log] ")
	line = strings.TrimSpace(line)
	if err := json.Unmarshal([]byte(line), &out); err != nil {
		t.Errorf("Output is not valid JSON: %v", err)
	}
	if out["method"] != entry.Method {
		t.Errorf("Expected method %s, got %+v", entry.Method, out["method"])
	}
	if out["path"] != entry.Path {
		t.Errorf("Expected path %s, got %+v", entry.Path, out["path"])
	}
}

func TestFileLogger_Log_WriteError(t *testing.T) {
	// Writer that always fails
	writer := &errorWriter{}
	fl := &FileLogger{Writer: bufio.NewWriter(writer)}
	entry := models.RequestLogEntry{
		Time:    time.Now(),
		Method:  "GET",
		Path:    "/test/path",
		Headers: http.Header{"X-Test-Header": []string{"test-value"}},
		TraceID: "test-trace-id",
	}
	_ = fl.Log(entry) // Should print error but not panic
}

type errorWriter struct{}

func (e *errorWriter) Write(p []byte) (n int, err error) { return 0, errors.New("write error") }
func (e *errorWriter) WriteString(s string) (n int, err error) {
	return 0, errors.New("write error")
}

func TestFileLogger_Start(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	fl := &FileLogger{Writer: writer}
	ch := make(chan models.Loggable, 1)
	ch <- models.RequestLogEntry{
		Time:    time.Now(),
		Method:  "GET",
		Path:    "/test/path",
		Headers: http.Header{"X-Test-Header": []string{"test-value"}},
		TraceID: "test-trace-id",
	}
	close(ch)
	err := fl.Start(ch)
	if err != nil {
		t.Fatalf("Start() error = %v, want nil", err)
	}
	// Wait for goroutine to finish
	writer.Flush()
}

func TestNewFileLogger_NoLogFile(t *testing.T) {
	cfg := models.Config{LogFile: ""}
	logger := NewFileLogger(cfg)
	if logger != nil {
		t.Errorf("Expected nil logger when LogFile is empty")
	}
}

func TestNewFileLogger_WithLogFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "snitchcraft_test.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	cfg := models.Config{LogFile: tmpfile.Name()}
	logger := NewFileLogger(cfg)
	if logger == nil {
		t.Fatalf("Expected logger, got nil")
	}
	if _, ok := logger.(*FileLogger); !ok {
		t.Errorf("Expected *FileLogger, got %T", logger)
	}
}
