package plugins

import "net/http"

type SuspiciousMethodCheck struct{}

func (s SuspiciousMethodCheck) Name() string {
    return "Suspicious HTTP Method"
}

func (s SuspiciousMethodCheck) Check(r *http.Request) (string, bool) {
    switch r.Method {
    case "TRACE", "TRACK", "DEBUG", "CONNECT":
        return "Use of uncommon or dangerous HTTP method: " + r.Method, true
    default:
        return "", false
    }
}

func init() {
    RegisterCheck(SuspiciousMethodCheck{})
}