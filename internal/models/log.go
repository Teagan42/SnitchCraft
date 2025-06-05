package models

import (
	"net/http"
	"time"

	"github.com/teagan42/snitchcraft/utils"
)

type Loggable interface {
	GetTraceID() string
	GetTime() time.Time
	GetMethod() string
	GetPath() string
	GetHeaders() http.Header
	GetType() string
}

type RequestLogEntry struct {
	Time    time.Time   `json:"time"`
	Method  string      `json:"method"`
	Path    string      `json:"path"`
	Headers http.Header `json:"headers"`
	TraceID string      `json:"trace_id"`
}

func (r RequestLogEntry) GetTraceID() string      { return r.TraceID }
func (r RequestLogEntry) GetTime() time.Time      { return r.Time }
func (r RequestLogEntry) GetMethod() string       { return r.Method }
func (r RequestLogEntry) GetPath() string         { return r.Path }
func (r RequestLogEntry) GetHeaders() http.Header { return r.Headers }
func (r RequestLogEntry) GetType() string         { return "request" }

type ResponseLogEntry struct {
	Time             time.Time         `json:"time"`
	Method           string            `json:"method"`
	Path             string            `json:"path"`
	Headers          http.Header       `json:"headers"`
	TraceID          string            `json:"trace_id"`
	Malicious        bool              `json:"malicious"`
	HeuristicResults []HeuristicResult `json:"heuristic_results"`
	StatusCode       int               `json:"status_code"`
}

func (r ResponseLogEntry) GetTraceID() string      { return r.TraceID }
func (r ResponseLogEntry) GetTime() time.Time      { return r.Time }
func (r ResponseLogEntry) GetMethod() string       { return r.Method }
func (r ResponseLogEntry) GetPath() string         { return r.Path }
func (r ResponseLogEntry) GetHeaders() http.Header { return r.Headers }
func (r ResponseLogEntry) IsMalicious() bool {
	return utils.Filter(r.HeuristicResults, func(hr HeuristicResult) bool {
		return hr.Issue != ""
	}) != nil
}
func (r ResponseLogEntry) GetHeuristicResults() []HeuristicResult { return r.HeuristicResults }
func (r ResponseLogEntry) GetStatusCode() int                     { return r.StatusCode }
func (r ResponseLogEntry) GetType() string                        { return "response" }
