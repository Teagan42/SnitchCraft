package interfaces

import "github.com/teagan42/snitchcraft/internal/models"

type LogSink interface {
	Log(log models.Loggable) error
}
