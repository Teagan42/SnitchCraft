package interfaces

import "net/http"

type Heuristic interface {
	Name() string
	Check(*http.Request) (string, bool)
}

type HeuristicCheck interface {
	Name() string
	Check(r *http.Request) (string, bool) // returns (reason, isMalicious)
}
