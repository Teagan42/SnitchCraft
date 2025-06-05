package models

import (
	"net/http"
	"testing"
)

func TestRequestResult_StructFields(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("failed to create http.Request: %v", err)
	}

	heuristics := []HeuristicResult{
		{}, // Add more fields if HeuristicResult has them
	}

	result := RequestResult{
		HeuristicResults: heuristics,
		Request:          req,
	}

	if len(result.HeuristicResults) != 1 {
		t.Errorf("expected 1 heuristic result, got %d", len(result.HeuristicResults))
	}
	if result.Request != req {
		t.Error("Request field does not match the original http.Request")
	}
}

func TestRequestResult_EmptyFields(t *testing.T) {
	result := RequestResult{}

	if len(result.HeuristicResults) != 0 {
		t.Errorf("expected empty HeuristicResults, got %v", result.HeuristicResults)
	}
	if result.Request != nil {
		t.Errorf("expected nil Request, got %v", result.Request)
	}
}
