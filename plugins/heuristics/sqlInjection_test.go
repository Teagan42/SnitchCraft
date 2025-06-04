package heuristics

import (
	"net/http"
	"testing"
)

func TestSQLInjectionCheck_Name(t *testing.T) {
	s := SQLInjectionCheck{}
	if s.Name() != "SQL Injection" {
		t.Errorf("expected Name to be 'SQL Injection', got '%s'", s.Name())
	}
}

func TestSQLInjectionCheck_Check(t *testing.T) {
	tests := []struct {
		name       string
		rawQuery   string
		wantMsg    string
		wantResult bool
	}{
		{
			name:       "No SQL injection",
			rawQuery:   "id=5&name=foo",
			wantMsg:    "",
			wantResult: false,
		},
		{
			name:       "Union select injection",
			rawQuery:   "id=1 UNION SELECT username, password FROM users",
			wantMsg:    "Query contains possible SQL injection",
			wantResult: true,
		},
		{
			name:       "OR 1=1 injection",
			rawQuery:   "user=admin' OR 1=1 --",
			wantMsg:    "Query contains possible SQL injection",
			wantResult: true,
		},
		{
			name:       "Drop table injection",
			rawQuery:   "q=DROP TABLE users",
			wantMsg:    "Query contains possible SQL injection",
			wantResult: true,
		},
		{
			name:       "Case insensitive match",
			rawQuery:   "q=UnIoN SeLeCt * from users",
			wantMsg:    "Query contains possible SQL injection",
			wantResult: true,
		},
	}

	s := SQLInjectionCheck{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://test?"+tt.rawQuery, nil)
			msg, result := s.Check(req)
			if msg != tt.wantMsg || result != tt.wantResult {
				t.Errorf("Check() = (%q, %v), want (%q, %v)", msg, result, tt.wantMsg, tt.wantResult)
			}
		})
	}
}
