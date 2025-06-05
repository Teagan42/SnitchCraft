package models

import (
	"testing"
)

func TestHeuristicResultFields(t *testing.T) {
	name := "TestHeuristic"
	issue := "TestIssue"
	result := HeuristicResult{
		Name:  name,
		Issue: issue,
	}

	if result.Name != name {
		t.Errorf("expected Name %q, got %q", name, result.Name)
	}
	if result.Issue != issue {
		t.Errorf("expected Issue %q, got %q", issue, result.Issue)
	}
}

func TestHeuristicResult_EmptyFields(t *testing.T) {
	result := HeuristicResult{}

	if result.Name != "" {
		t.Errorf("expected Name to be empty, got %q", result.Name)
	}
	if result.Issue != "" {
		t.Errorf("expected Issue to be empty, got %q", result.Issue)
	}
}
