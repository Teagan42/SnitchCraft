package interactors

import (
	"net/http"
	"sync"

	"github.com/teagan42/snitchcraft/internal/interfaces"
	"github.com/teagan42/snitchcraft/internal/models"
	"github.com/teagan42/snitchcraft/utils"
)

var RegisteredHeuristics []interfaces.Heuristic

func RegisterHeuristic(h interfaces.Heuristic) {
	RegisteredHeuristics = append(RegisteredHeuristics, h)
}

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
	return utils.Do(
		results,
		func(result models.HeuristicResult) {
			metricsChannels.heuristicsChan <- result
		},
	)
}

func SyncRunHeuristicChecks(
	req *http.Request,
) []models.HeuristicResult {
	var results []models.HeuristicResult

	for _, heuristic := range RegisteredHeuristics {
		if name, ok := heuristic.Check(req); ok {
			var result = models.HeuristicResult{
				Name:  heuristic.Name(),
				Issue: name,
			}
			// Send result to metrics channel
			metricsChannels.heuristicsChan <- result
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
	resultChan := make(chan models.HeuristicResult, len(RegisteredHeuristics))

	for _, heuristic := range RegisteredHeuristics {
		wg.Add(1)
		go func(h interfaces.Heuristic) {
			defer wg.Done()
			if name, ok := h.Check(req); ok {
				var result = models.HeuristicResult{
					Name:  h.Name(),
					Issue: name,
				}
				// Send result to metrics channel
				metricsChannels.heuristicsChan <- result
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
