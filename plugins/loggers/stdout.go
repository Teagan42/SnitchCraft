package loggers

import (
	"encoding/json"
	"fmt"

	"github.com/teagan42/snitchcraft/internal/interfaces"
	"github.com/teagan42/snitchcraft/internal/models"
)

type StdoutLogger struct {
}

func (s *StdoutLogger) Name() string {
	return "StdoutLogger"
}

func (s *StdoutLogger) Start(logChannel chan models.Loggable) error {
	go func() {
		for log := range logChannel {
			if err := s.Log(log); err != nil {
				fmt.Printf("[loggers][stdout] Failed to log to stdout: %v\n", err)
			}
		}
	}()
	return nil
}

func (s *StdoutLogger) Log(entry models.Loggable) error {
	data, err := json.Marshal(entry)
	if err != nil {
		fmt.Printf("[loggers][stdout] Failed to marshal loggable: %v\n", err)
		return err
	}
	fmt.Printf("[loggers][log] %s\n", string(data))
	return nil
}

func NewStdoutLogger(cfg models.Config) interfaces.Logger {
	fmt.Println("[loggers] initializing StdoutLogger...")
	return &StdoutLogger{}
}

func init() {
	RegisterLogger(NewStdoutLogger)
}
