package main

import (
    "encoding/json"
    "log"
    "net/http"
    "time"
)

type LogEntry struct {
    Time      time.Time       `json:"time"`
    Method    string          `json:"method"`
    Path      string          `json:"path"`
    Headers   http.Header     `json:"headers"`
    Malicious bool            `json:"malicious"`
    Reasons   []string        `json:"reasons"`
}

func LogRequest(r *http.Request, reasons []string) {
    entry := LogEntry{
        Time:      time.Now().UTC(),
        Method:    r.Method,
        Path:      r.URL.RequestURI(),
        Headers:   r.Header.Clone(),
        Malicious: len(reasons) > 0,
        Reasons:   reasons,
    }

    jsonData, _ := json.Marshal(entry)
    log.Println(string(jsonData))
}