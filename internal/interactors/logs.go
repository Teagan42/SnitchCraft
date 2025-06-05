package interactors

import (
	"fmt"

	"github.com/teagan42/snitchcraft/internal/interfaces"
	"github.com/teagan42/snitchcraft/internal/models"
	"github.com/teagan42/snitchcraft/plugins/loggers"
	"github.com/teagan42/snitchcraft/utils"
)

var logsChannel chan models.Loggable = make(chan models.Loggable, 100)

func LogWorker(cfg models.Config) {
	fmt.Println("[interactors][logs] starting LogWorker...")
	utils.Do(loggers.RegisteredLoggers, func(l func(cfg models.Config) interfaces.Logger) {
		logger := l(cfg)
		if logger != nil {
			fmt.Printf("[interactors][logs] starting logger: %s\n", logger.Name())
			go func(l interfaces.Logger) {
				if err := l.Start(logsChannel); err != nil {
					fmt.Printf("[interactors][logs] error starting logger %s: %v\n", l.Name(), err)
				}
			}(logger)
		}
	})
	fmt.Println("[interactors][logs] started LogWorker...")
}
