package heuristics

import "github.com/teagan42/snitchcraft/internal/interfaces"

var RegisteredHeuristics []interfaces.Heuristic

func RegisterHeuristic(h interfaces.Heuristic) {
	RegisteredHeuristics = append(RegisteredHeuristics, h)
}
