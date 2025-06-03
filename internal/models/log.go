package models

import (
	"net/http"
	"time"
)

type LogEntry struct {
	Time      time.Time   `json:"time"`
	Method    string      `json:"method"`
	Path      string      `json:"path"`
	Headers   http.Header `json:"headers"`
	Malicious bool        `json:"malicious"`
	Reasons   []string    `json:"reasons"`
}
