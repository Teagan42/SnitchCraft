package loggers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/teagan42/snitchcraft/internal/interfaces"
	"github.com/teagan42/snitchcraft/internal/models"
)

type LokiLogger struct {
	client  *http.Client
	lokiURL string
	labels  map[string]string
}

func (l *LokiLogger) Name() string {
	return "LokiLogger"
}

func (l *LokiLogger) Start(logChannel chan models.Loggable) error {
	go func() {
		for log := range logChannel {
			if err := l.Log(log); err != nil {
				fmt.Printf("Failed to log to Loki: %v\n", err)
			}
		}
	}()
	return nil
}

func (l *LokiLogger) Log(log models.Loggable) error {
	// Marshal the log entry
	entryBytes, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("[loggers][loki] Failed to marshal loggable: %w", err)
	}

	// Prepare Loki payload
	lokiPayload := map[string]interface{}{
		"streams": []map[string]interface{}{
			{
				"stream": l.labels,
				"values": [][]string{
					{
						fmt.Sprintf("%d000000", log.GetTime().Unix()), // nanoseconds as string
						string(entryBytes),
					},
				},
			},
		},
	}

	payloadBytes, err := json.Marshal(lokiPayload)
	if err != nil {
		return fmt.Errorf("[loggers][loki] failed to marshal Loki payload: %w", err)
	}

	req, err := http.NewRequest("POST", l.lokiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("[loggers][loki] failed to create request to Loki: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := l.client.Do(req)
	if err != nil {
		return fmt.Errorf("[loggers][loki] failed to send log to Loki: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("[loggers][loki] received non-2xx from Loki: %s", resp.Status)
	}

	return nil
}

func NewLokiLogger(cfg models.Config) interfaces.Logger {
	if cfg.LokiUrl == "" {
		fmt.Printf("[loggers] not initializing LokiLogger: missing env LOKI_URL\n")
		return nil // No Loki URL configured, skip logger

		fmt.Printf("[loggers] initializing LokiLogger with URL: %s\n", cfg.LokiUrl)
	}
	return &LokiLogger{
		client:  &http.Client{},
		lokiURL: cfg.LokiUrl,
		labels: map[string]string{
			"job":  "snitchcraft",
			"host": "localhost", // optional override with env if needed
		},
	}
}

func init() {
	RegisterLogger(NewLokiLogger)
}
