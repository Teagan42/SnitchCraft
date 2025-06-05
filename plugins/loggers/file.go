package loggers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/teagan42/snitchcraft/internal/interfaces"
	"github.com/teagan42/snitchcraft/internal/models"
)

type FileLogger struct {
	Writer *bufio.Writer
}

func (f *FileLogger) Name() string {
	return "FileLogger"
}

func (f *FileLogger) Start(logChannel chan models.Loggable) error {
	go func() {
		for log := range logChannel {
			if err := f.Log(log); err != nil {
				fmt.Printf("[loggers][stdout] Failed to log to stdout: %v\n", err)
			}
		}
	}()
	return nil
}

func (f *FileLogger) Log(entry models.Loggable) error {
	data, err := json.Marshal(entry)
	if err != nil {
		fmt.Printf("[loggers][stdout] Failed to marshal loggable: %v\n", err)
		return err
	}

	if _, err := f.Writer.WriteString(fmt.Sprintf("%s\n", data)); err != nil {
		fmt.Printf("[loggers][stdout] Failed to write log entry: %v\n", err)
	} else {
		f.Writer.Flush()
	}

	return nil
}

func NewFileLogger(cfg models.Config) interfaces.Logger {
	if cfg.LogFile == "" {
		fmt.Println("[loggers] LogFile not set, skipping FileLogger initialization")
		return nil
	}
	fmt.Println("[loggers] initializing FileLogger...")
	f, err := os.OpenFile("snitchcraft.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	b := bufio.NewWriter(f)
	return &FileLogger{
		Writer: b,
	}
}

func init() {
	RegisterLogger(NewStdoutLogger)
}
