package heuristics

import (
	"net/http"
	"testing"
)

func TestSuspiciousUserAgent_Check(t *testing.T) {
	tests := []struct {
		name       string
		userAgent  string
		wantMsg    string
		wantResult bool
	}{
		{
			name:       "Contains curl",
			userAgent:  "curl/7.68.0",
			wantMsg:    "Detected script-like User-Agent",
			wantResult: true,
		},
		{
			name:       "Contains python",
			userAgent:  "python-requests/2.25.1",
			wantMsg:    "Detected script-like User-Agent",
			wantResult: true,
		},
		{
			name:       "Mixed case CURL",
			userAgent:  "CuRL/8.0.1",
			wantMsg:    "Detected script-like User-Agent",
			wantResult: true,
		},
		{
			name:       "Mixed case PYTHON",
			userAgent:  "PyThOn-urllib/3.9",
			wantMsg:    "Detected script-like User-Agent",
			wantResult: true,
		},
		{
			name:       "Normal browser",
			userAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			wantMsg:    "",
			wantResult: false,
		},
		{
			name:       "Empty User-Agent",
			userAgent:  "",
			wantMsg:    "",
			wantResult: false,
		},
	}

	check := SuspiciousUserAgent{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			if tt.userAgent != "" {
				req.Header.Set("User-Agent", tt.userAgent)
			}
			gotMsg, gotResult := check.Check(req)
			if gotMsg != tt.wantMsg || gotResult != tt.wantResult {
				t.Errorf("Check() = (%q, %v), want (%q, %v)", gotMsg, gotResult, tt.wantMsg, tt.wantResult)
			}
		})
	}
}
