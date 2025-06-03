package loggers

import (
	"encoding/json"
	"fmt"
	"snitchcraft/internal/interfaces"
	"snitchcraft/internal/models"
)

type StdoutLogger struct{}

func (s *StdoutLogger) Send(entry models.LogEntry) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

var _ interfaces.LogSink = (*StdoutLogger)(nil)
