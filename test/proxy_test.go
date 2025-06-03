package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"snitchcraft/heuristics"
)

func TestProxyMaliciousDetection(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))
	defer backend.Close()

	req := httptest.NewRequest("GET", "/?q=1+OR+1=1", nil)
	req.Header.Set("User-Agent", "curl/7.81.0")
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		target, _ := http.NewRequest(r.Method, backend.URL+r.URL.RequestURI(), r.Body)
		target.Header = r.Header.Clone()
		results := plugins.RunChecks(r)
		if len(results) == 0 {
			t.Errorf("expected malicious check to detect something, got none")
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "proxied")
	})

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "proxied") {
		t.Errorf("expected body to contain 'proxied', got '%s'", rec.Body.String())
	}
}