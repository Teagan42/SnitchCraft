package interactors

import (
	"github.com/teagan42/snitchcraft/internal/interfaces"
	"github.com/teagan42/snitchcraft/internal/models"
	"github.com/teagan42/snitchcraft/plugins/loggers"
)

var logsChannel chan models.Loggable = make(chan models.Loggable, 100)

func LogWorker(cfg models.Config) {
	var logSink interfaces.LogSink
	if cfg.LogForwardURL != "" {
		logSink = &loggers.StdoutLogger{}
	} else {
		logSink = &loggers.StdoutLogger{}
	}

	for {
		select {
		case log := <-logsChannel:
			if err := logSink.Log(log); err != nil {
				// Handle error, e.g., log to stderr or retry
				continue
			}
		default:
			// No logs to process, just continue
			continue
		}
	}
}
