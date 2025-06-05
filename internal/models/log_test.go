package models

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestRequestLogEntry_ImplementsLoggable(t *testing.T) {
	var _ Loggable = RequestLogEntry{}
}

func TestRequestLogEntry_Getters(t *testing.T) {
	now := time.Now()
	headers := http.Header{"X-Test": []string{"val"}}
	entry := RequestLogEntry{
		Time:    now,
		Method:  "GET",
		Path:    "/test",
		Headers: headers,
		TraceID: "abc123",
	}

	if got := entry.GetTraceID(); got != "abc123" {
		t.Errorf("GetTraceID() = %v, want %v", got, "abc123")
	}
	if got := entry.GetTime(); !got.Equal(now) {
		t.Errorf("GetTime() = %v, want %v", got, now)
	}
	if got := entry.GetMethod(); got != "GET" {
		t.Errorf("GetMethod() = %v, want %v", got, "GET")
	}
	if got := entry.GetPath(); got != "/test" {
		t.Errorf("GetPath() = %v, want %v", got, "/test")
	}
	if got := entry.GetHeaders(); !reflect.DeepEqual(got, headers) {
		t.Errorf("GetHeaders() = %v, want %v", got, headers)
	}
	if got := entry.GetType(); got != "request" {
		t.Errorf("GetType() = %v, want %v", got, "request")
	}
}

func TestResponseLogEntry_Getters(t *testing.T) {
	now := time.Now()
	headers := http.Header{"X-Test": []string{"val"}}
	heuristics := []HeuristicResult{
		{Name: "h1", Issue: "issue1"},
	}
	entry := ResponseLogEntry{
		Time:             now,
		Method:           "POST",
		Path:             "/api",
		Headers:          headers,
		TraceID:          "xyz789",
		Malicious:        true,
		HeuristicResults: heuristics,
		StatusCode:       403,
	}

	if got := entry.GetTraceID(); got != "xyz789" {
		t.Errorf("GetTraceID() = %v, want %v", got, "xyz789")
	}
	if got := entry.GetTime(); !got.Equal(now) {
		t.Errorf("GetTime() = %v, want %v", got, now)
	}
	if got := entry.GetMethod(); got != "POST" {
		t.Errorf("GetMethod() = %v, want %v", got, "POST")
	}
	if got := entry.GetPath(); got != "/api" {
		t.Errorf("GetPath() = %v, want %v", got, "/api")
	}
	if got := entry.GetHeaders(); !reflect.DeepEqual(got, headers) {
		t.Errorf("GetHeaders() = %v, want %v", got, headers)
	}
	if got := entry.IsMalicious(); got != true {
		t.Errorf("IsMalicious() = %v, want %v", got, true)
	}
	if got := entry.GetHeuristicResults(); !reflect.DeepEqual(got, heuristics) {
		t.Errorf("GetHeuristicResults() = %v, want %v", got, heuristics)
	}
	if got := entry.GetStatusCode(); got != 403 {
		t.Errorf("GetStatusCode() = %v, want %v", got, 403)
	}
	if got := entry.GetType(); got != "response" {
		t.Errorf("GetType() = %v, want %v", got, "response")
	}
}
