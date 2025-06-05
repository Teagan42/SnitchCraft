package loggers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/teagan42/snitchcraft/internal/models"
)

func TestStdoutLogger_Log_ValidEntry(t *testing.T) {
	logger := &StdoutLogger{}
	entry := models.RequestLogEntry{
		Time:    time.Now(),
		Method:  "GET",
		Path:    "/test/path",
		Headers: http.Header{"X-Test-Header": []string{"test-value"}},
		TraceID: "test-trace-id",
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := logger.Log(entry)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Restore stdout and read output
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	expectedJSON, _ := json.Marshal(entry)
	expected := fmt.Sprintf("[loggers][log] %s\n", string(expectedJSON))
	if output != expected {
		t.Errorf("expected output %q, got %q", expected, output)
	}
}
