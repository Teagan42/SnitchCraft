package heuristics

import (
	"net/http"
	"testing"
)

func TestAnomalousHeaderCheck_Name(t *testing.T) {
	h := AnomalousHeaderCheck{}
	if h.Name() != "Header Anomaly" {
		t.Errorf("expected Name to be 'Header Anomaly', got '%s'", h.Name())
	}
}

func TestAnomalousHeaderCheck_Check(t *testing.T) {
	tests := []struct {
		name       string
		headers    map[string]string
		wantMsg    string
		wantResult bool
	}{
		{
			name:       "Missing Accept header",
			headers:    map[string]string{"Host": "example.com"},
			wantMsg:    "Missing Accept header",
			wantResult: true,
		},
		{
			name:       "Missing Host header",
			headers:    map[string]string{"Accept": "text/html"},
			wantMsg:    "Missing Host header",
			wantResult: true,
		},
		{
			name: "Referer contains suspicious domain",
			headers: map[string]string{
				"Accept":  "text/html",
				"Host":    "example.com",
				"Referer": "http://evil.com/page",
			},
			wantMsg:    "Referer contains suspicious domain",
			wantResult: true,
		},
		{
			name: "All headers present and safe referer",
			headers: map[string]string{
				"Accept":  "text/html",
				"Host":    "example.com",
				"Referer": "http://good.com/page",
			},
			wantMsg:    "",
			wantResult: false,
		},
		{
			name: "All headers present and no referer",
			headers: map[string]string{
				"Accept": "text/html",
				"Host":   "example.com",
			},
			wantMsg:    "",
			wantResult: false,
		},
	}

	h := AnomalousHeaderCheck{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://test", nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			msg, result := h.Check(req)
			if msg != tt.wantMsg || result != tt.wantResult {
				t.Errorf("Check() = (%q, %v), want (%q, %v)", msg, result, tt.wantMsg, tt.wantResult)
			}
		})
	}
}
