package interactors

import (
	"fmt"

	"github.com/teagan42/snitchcraft/internal/interfaces"
	"github.com/teagan42/snitchcraft/internal/models"
	"github.com/teagan42/snitchcraft/plugins/metrics"
	"github.com/teagan42/snitchcraft/utils"
)

func MetricsWorker(cfg models.Config, resultsChannel chan models.RequestResult) {
	fmt.Printf("[interactors][metrics] starting MetricsWorker with config: %+v\n", cfg)

	utils.Do(metrics.RegisteredMetricsPlugins, func(p func(cfg models.Config) interfaces.MetricsPlugin) {
		plugin := p(cfg)
		if plugin != nil {
			fmt.Printf("[interactors][metrics] starting metrics plugin: %s\n", plugin.Name())
			if err := plugin.Start(resultsChannel); err != nil {
				fmt.Printf("[interactors[metrics] failed to start metrics plugin %s: %v\n", plugin.Name(), err)
			}
		}
	})
	fmt.Println("[interactors][metrics] started MetricsWorker...")
}
