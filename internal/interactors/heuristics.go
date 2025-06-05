package interactors

import (
	"net/http"
	"sync"

	"github.com/teagan42/snitchcraft/internal/interfaces"
	"github.com/teagan42/snitchcraft/internal/models"
	"github.com/teagan42/snitchcraft/plugins/heuristics"
)

func RunHeuristicChecks(
	req *http.Request,
	cfg models.Config,
) []models.HeuristicResult {
	var results []models.HeuristicResult
	if cfg.ParallelChecks {
		results = AsyncRunHeuristicChecks(req)
	} else {
		results = SyncRunHeuristicChecks(req)
	}
	return results
}

func SyncRunHeuristicChecks(
	req *http.Request,
) []models.HeuristicResult {
	var results []models.HeuristicResult

	for _, heuristic := range heuristics.RegisteredHeuristics {
		if name, ok := heuristic.Check(req); ok {
			var result = models.HeuristicResult{
				Name:  heuristic.Name(),
				Issue: name,
			}
			// Append to results
			results = append(results, result)
		}
	}

	return results
}

func AsyncRunHeuristicChecks(
	req *http.Request,
) []models.HeuristicResult {
	var wg sync.WaitGroup
	resultChan := make(chan models.HeuristicResult, 10*len(heuristics.RegisteredHeuristics))

	for _, heuristic := range heuristics.RegisteredHeuristics {
		wg.Add(1)
		go func(h interfaces.Heuristic) {
			defer wg.Done()
			if name, ok := h.Check(req); ok {
				var result = models.HeuristicResult{
					Name:  h.Name(),
					Issue: name,
				}
				// Send result to result channel
				resultChan <- result
			}
		}(heuristic)
	}

	wg.Wait()
	close(resultChan)

	var results []models.HeuristicResult
	for res := range resultChan {
		results = append(results, res)
	}

	return results
}
