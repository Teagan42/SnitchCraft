package metrics

import (
	"github.com/teagan42/snitchcraft/internal/interfaces"
	"github.com/teagan42/snitchcraft/internal/models"
)

var RegisteredMetricsPlugins []func(models.Config) interfaces.MetricsPlugin

func RegisterMetricsPlugin(m func(models.Config) interfaces.MetricsPlugin) {
	RegisteredMetricsPlugins = append(RegisteredMetricsPlugins, m)
}
