package interfaces

import "snitchcraft/internal/models"

type LogSink interface {
	Send(log models.LogEntry) error
}
