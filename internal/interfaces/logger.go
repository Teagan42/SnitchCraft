package interfaces

import "github.com/teagan42/snitchcraft/internal/models"

type Logger interface {
	Name() string
	Log(log models.Loggable) error
	Start(logChannel chan models.Loggable) error
}
