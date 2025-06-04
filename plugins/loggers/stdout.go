package loggers

import (
	"encoding/json"
	"fmt"

	"github.com/teagan42/snitchcraft/internal/interfaces"
	"github.com/teagan42/snitchcraft/internal/models"
)

type StdoutLogger struct{}

func (s *StdoutLogger) Log(entry models.Loggable) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

var _ interfaces.LogSink = (*StdoutLogger)(nil)
