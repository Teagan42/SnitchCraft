package interfaces

import "github.com/teagan42/snitchcraft/internal/models"

type MetricsPlugin interface {
	Name() string
	Start(chan models.RequestResult) error
}
