package heuristics

import (
	"fmt"
	"net/http"
	"regexp"
)

var sqlInjectionRegex = regexp.MustCompile(`(?i)(union\s+select|or\s+1=1|drop\s+table)`)

type SQLInjectionCheck struct{}

func (s SQLInjectionCheck) Name() string {
	return "SQL Injection"
}

func (s SQLInjectionCheck) Check(r *http.Request) (string, bool) {
	if sqlInjectionRegex.MatchString(r.URL.RawQuery) {
		return "Query contains possible SQL injection", true
	}
	return "", false
}

func init() {
	fmt.Println("[heuristics] registering SQLInjectionCheck heuristic...")
	RegisterHeuristic(SQLInjectionCheck{})
}
