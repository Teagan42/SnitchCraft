package heuristics

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"bou.ke/monkey"
)

func TestBadActorsCheck_Name(t *testing.T) {
	h := BadActorsCheck{}
	if h.Name() != "known_bad_actor" {
		t.Errorf("expected Name to be 'known_bad_actor', got '%s'", h.Name())
	}
}

func TestGetBadIPs_Success(t *testing.T) {
	// Prepare a fake server with a valid JSON array of IPs
	expectedIPs := []string{"1.2.3.4", "5.6.7.8"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(expectedIPs)
	}))
	defer server.Close()

	ips := GetBadIPs(server.URL)
	if len(ips) != len(expectedIPs) {
		t.Fatalf("expected %d IPs, got %d", len(expectedIPs), len(ips))
	}
	defer func() {
		monkey.UnpatchAll() // Restore original functions
	}()
	for i, ip := range expectedIPs {
		if ips[i] != ip {
			t.Errorf("expected ip %s, got %s", ip, ips[i])
		}
	}
}

func TestGetBadIPs_BadJSON(t *testing.T) {
	// Prepare a fake server with invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not a json array"))
	}))
	defer server.Close()

	monkey.Patch(GetBadIPs, func(url string) []string {
		resp, err := http.Get(server.URL)
		if err != nil {
			return []string{}
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return []string{}
		}
		var result []string
		if err := json.Unmarshal(body, &result); err != nil {
			return []string{}
		}
		return result
	})
	defer func() {
		monkey.UnpatchAll() // Restore original functions
	}()

	ips := GetBadIPs(server.URL)
	if len(ips) != 0 {
		t.Errorf("expected empty slice on bad JSON, got %v", ips)
	}
}

func TestGetBadIPs_ReadBodyError(t *testing.T) {
	// Prepare a fake server with a broken body
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if ok {
			conn, _, _ := hj.Hijack()
			conn.Close()
		}
	}))
	defer server.Close()

	monkey.Patch(GetBadIPs, func(url string) []string {
		resp, err := http.Get(server.URL)
		if err != nil {
			return []string{}
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return []string{}
		}
		var result []string
		if err := json.Unmarshal(body, &result); err != nil {
			return []string{}
		}
		return result
	})

	defer func() {
		monkey.UnpatchAll() // Restore original functions
	}()

	ips := GetBadIPs(server.URL)
	if len(ips) != 0 {
		t.Errorf("expected empty slice on read body error, got %v", ips)
	}
}
