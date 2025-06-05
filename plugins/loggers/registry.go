package loggers

import (
	"github.com/teagan42/snitchcraft/internal/interfaces"
	"github.com/teagan42/snitchcraft/internal/models"
)

var RegisteredLoggers []func(cfg models.Config) interfaces.Logger

func RegisterLogger(l func(cfg models.Config) interfaces.Logger) {
	RegisteredLoggers = append(RegisteredLoggers, l)
}
