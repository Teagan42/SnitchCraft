package interactors

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/teagan42/snitchcraft/internal/models"
)

func TestStartProxyServer_Forbidden(t *testing.T) {
	// Setup config with dummy backend URL and listen port
	cfg := models.Config{
		BackendURL: "http://localhost:9999",
		ListenPort: ":0", // random port
	}

	// Start server in goroutine
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		StartProxyServer(cfg)
	}))
	defer server.Close()

	// Make request to server
	req, _ := http.NewRequest("DEBUG", server.URL, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(body, []byte("Forbidden")) {
		t.Errorf("expected forbidden message in body, got %s", string(body))
	}
}

func TestStartProxyServer_ProxySuccess(t *testing.T) {
	// Setup a dummy backend server
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Backend", "yes")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("backend response"))
	}))
	defer backend.Close()

	cfg := models.Config{
		BackendURL: backend.URL,
		ListenPort: ":0",
	}

	// Start the proxy server
	var wg sync.WaitGroup
	wg.Add(1)
	var proxy *httptest.Server
	proxy = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		StartProxyServer(cfg)
		wg.Done()
	}))
	defer proxy.Close()

	// Make request to proxy
	req, _ := http.NewRequest("GET", proxy.URL, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(body, []byte("backend response")) {
		t.Errorf("expected backend response in body, got %s", string(body))
	}
	if resp.Header.Get("X-Backend") != "yes" {
		t.Errorf("expected X-Backend header, got %v", resp.Header)
	}
}

func TestStartProxyServer_UpstreamError(t *testing.T) {
	cfg := models.Config{
		BackendURL: "http://localhost:65534", // unused port
		ListenPort: ":0",
	}

	// Start the proxy server
	var wg sync.WaitGroup
	wg.Add(1)
	var proxy *httptest.Server
	proxy = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		StartProxyServer(cfg)
		wg.Done()
	}))
	defer proxy.Close()

	// Make request to proxy
	req, _ := http.NewRequest("GET", proxy.URL, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadGateway {
		t.Errorf("expected status 502, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(body, []byte("Bad Gateway")) {
		t.Errorf("expected bad gateway message in body, got %s", string(body))
	}
}

func TestSetupTracing_NoPanic(t *testing.T) {
	// Just ensure setupTracing does not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("setupTracing panicked: %v", r)
		}
	}()
	setupTracing()
}
