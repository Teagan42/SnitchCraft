package interactors

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/teagan42/snitchcraft/internal/models"
	"github.com/teagan42/snitchcraft/utils"

	"maps"
)

var resultChan = make(chan models.RequestResult, 10)

func LogResponse(r *http.Request, results []models.HeuristicResult, start time.Time) {
	logsChannel <- models.ResponseLogEntry{
		Time:             time.Now(),
		Method:           r.Method,
		Path:             r.URL.Path,
		Headers:          maps.Clone(r.Header),
		StatusCode:       r.Response.StatusCode,
		Duration:         uint64(time.Since(start).Nanoseconds()),
		HeuristicResults: results,
		TraceID:          r.Header.Get("X-Trace-ID"),
		Malicious:        len(utils.Map(results, func(hr models.HeuristicResult) bool { return hr.Issue != "" })) > 0,
	}
}

func GetHandler(cfg models.Config) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		var start = time.Now()
		fmt.Printf("[interactors][server] received request: %s %s\n", r.Method, r.URL.Path)
		r.Header.Set("X-Trace-ID", uuid.New().String())
		var results = RunHeuristicChecks(r, cfg)

		// pass to logging via channel
		resultChan <- models.RequestResult{
			Request:          r,
			HeuristicResults: results,
			Duration:         uint64(time.Since(start).Milliseconds()),
		}
		fmt.Printf("[interactors][server] heuristic checks completed for %s %s, results: %+v\n", r.Method, r.URL.Path, results)
		if len(utils.Filter(results, func(r models.HeuristicResult) bool {
			return len(r.Issue) > 0
		})) > 0 {
			r.Response = &http.Response{
				StatusCode: http.StatusForbidden,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(io.Reader(nil)), // no body for forbidden response
				Status:     "403 Forbidden",
				Request:    r,
			}
			LogResponse(r, results, start)
			w.WriteHeader(http.StatusForbidden)
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{"error": "Forbidden"}`))
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			resultChan <- models.RequestResult{Request: r, HeuristicResults: results}
			return
		}

		// forward request
		req, _ := http.NewRequest(r.Method, cfg.BackendURL+r.URL.RequestURI(), r.Body)
		req.Header = r.Header.Clone()

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			r.Response = &http.Response{
				StatusCode: http.StatusBadGateway,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(io.Reader(nil)), // no body for bad gateway response
				Status:     "502 Bad Gateway",
				Request:    r,
			}
			LogResponse(r, results, start)
			http.Error(w, "upstream error", http.StatusBadGateway)
			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write([]byte(`{"error": "Bad Gateway"}`)); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			return
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Printf("[interactors][server] error closing response body: %v\n", err)
			}
		}()
		r.Response = resp
		LogResponse(r, results, start)
		maps.Copy(w.Header(), resp.Header)
		w.WriteHeader(resp.StatusCode)
		if _, err := io.Copy(w, resp.Body); err != nil {
			http.Error(w, "[interactors][server] error copying response body", http.StatusInternalServerError)
			return
		}
	}

	return handler
}

func StartProxyServer(cfg models.Config) error {
	fmt.Printf("[interactors][server] starting proxy server with config: %+v\n", cfg)
	go MetricsWorker(cfg, resultChan)
	go LogWorker(cfg)

	http.HandleFunc("/", GetHandler(cfg))

	fmt.Printf("[interactors][server] starting proxy server on %s\n", cfg.ListenPort)
	return http.ListenAndServe(cfg.ListenPort, nil)
}
