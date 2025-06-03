package heuristics
import (
    "net/http"
    "sync"
)

type HeuristicResult struct {
    Name   string
    Reason string
}

type HeuristicCheck interface {
    Name() string
    Check(r *http.Request) (string, bool) // returns (reason, isMalicious)
}