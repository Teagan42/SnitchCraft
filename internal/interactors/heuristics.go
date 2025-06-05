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
		var issue, err = heuristic.Check(req)
		if err {
			results = append(results, models.HeuristicResult{
				Name:  heuristic.Name(),
				Issue: issue,
			})
		} else {
			// If heuristic does not match, we still want to return it with an empty issue
			results = append(results, models.HeuristicResult{
				Name:  heuristic.Name(),
				Issue: "",
			})
		}
	}

	return results
}

func AsyncRunHeuristicChecks(
	req *http.Request,
) []models.HeuristicResult {
	var wg sync.WaitGroup
	resultChan := make(chan models.HeuristicResult, len(heuristics.RegisteredHeuristics))

	for _, heuristic := range heuristics.RegisteredHeuristics {
		wg.Add(1)
		go func(h interfaces.Heuristic) {
			defer wg.Done()
			var issue, err = h.Check(req)
			if err {
				resultChan <- models.HeuristicResult{
					Name:  h.Name(),
					Issue: issue,
				}
			} else {
				resultChan <- models.HeuristicResult{
					Name:  h.Name(),
					Issue: "",
				}
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
