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

func StartProxyServer(cfg models.Config) error {
	fmt.Printf("[interactors][server] starting proxy server with config: %+v\n", cfg)
	go MetricsWorker(cfg, resultChan)
	go LogWorker(cfg)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[interactors][server] received request: %s %s\n", r.Method, r.URL.Path)
		r.Header.Set("X-Trace-ID", uuid.New().String())
		logsChannel <- models.RequestLogEntry{
			Time:    time.Now(),
			Method:  r.Method,
			Path:    r.URL.Path,
			Headers: r.Header,
			TraceID: r.Header.Get("X-Trace-ID"),
		}
		var results = RunHeuristicChecks(r, cfg)

		// pass to logging via channel
		resultChan <- models.RequestResult{Request: r, HeuristicResults: results}

		if len(utils.Map(results, func(r models.HeuristicResult) bool {
			return r.Issue != ""
		})) > 0 {
			// if any heuristic found an issue, return 403
			w.WriteHeader(http.StatusForbidden)
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{"error": "Forbidden"}`))
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			logsChannel <- models.ResponseLogEntry{
				Time:             time.Now(),
				Method:           r.Method,
				Path:             r.URL.Path,
				Headers:          r.Header,
				TraceID:          r.Header.Get("X-Trace-ID"),
				Malicious:        true,
				HeuristicResults: results,
				StatusCode:       http.StatusForbidden,
			}
			return
		}

		// forward request
		req, _ := http.NewRequest(r.Method, cfg.BackendURL+r.URL.RequestURI(), r.Body)
		req.Header = r.Header.Clone()

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			logsChannel <- models.ResponseLogEntry{
				Time:             time.Now(),
				Method:           r.Method,
				Path:             r.URL.Path,
				Headers:          r.Header,
				TraceID:          r.Header.Get("X-Trace-ID"),
				Malicious:        false,
				HeuristicResults: results,
				StatusCode:       http.StatusBadGateway,
			}
			http.Error(w, "upstream error", http.StatusBadGateway)
			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write([]byte(`{"error": "Bad Gateway"}`)); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			return
		}
		defer resp.Body.Close()
		logsChannel <- models.ResponseLogEntry{
			Time:             time.Now(),
			Method:           r.Method,
			Path:             r.URL.Path,
			Headers:          r.Header,
			TraceID:          r.Header.Get("X-Trace-ID"),
			Malicious:        false,
			HeuristicResults: results,
			StatusCode:       resp.StatusCode,
		}
		maps.Copy(w.Header(), resp.Header)
		w.WriteHeader(resp.StatusCode)
		if _, err := io.Copy(w, resp.Body); err != nil {
			http.Error(w, "[interactors][server] error copying response body", http.StatusInternalServerError)
			return
		}
	})

	fmt.Printf("[interactors][server] starting proxy server on %s\n", cfg.ListenPort)
	return http.ListenAndServe(cfg.ListenPort, nil)
}
