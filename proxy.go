package main

import (
    "io"
    "net/http"
    "net/url"
    "heuristics"
)

func StartProxy(listenPort string, targetHost string) {
    target, _ := url.Parse(targetHost)

    handler := func(w http.ResponseWriter, r *http.Request) {
        req, _ := http.NewRequest(r.Method, target.String()+r.URL.RequestURI(), r.Body)
        req.Header = r.Header.Clone()

        // Run plugin checks
        results := heuristics.RunChecks(r)

        // Build reason slice
        var reasons []string
        for _, r := range results {
            reasons = append(reasons, r.Name+": "+r.Reason)
        }

        LogRequest(r, reasons)

        // Forward the request
        resp, err := http.DefaultClient.Do(req)
        if err != nil {
            http.Error(w, "Upstream error", http.StatusBadGateway)
            return
        }
        defer resp.Body.Close()

        for k, v := range resp.Header {
            w.Header()[k] = v
        }
        w.WriteHeader(resp.StatusCode)
        io.Copy(w, resp.Body)
    }

    http.ListenAndServe(listenPort, http.HandlerFunc(handler))
}