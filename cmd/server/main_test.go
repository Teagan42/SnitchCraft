package main

import (
	"errors"
	"log"
	"testing"

	"bou.ke/monkey"

	"github.com/teagan42/snitchcraft/internal/config"
	"github.com/teagan42/snitchcraft/internal/interactors"
	"github.com/teagan42/snitchcraft/internal/models"
)

// mockConfigLoader replaces config.Load for testing.
func mockConfigLoader(cfg models.Config, err error) func() (models.Config, error) {
	return func() (models.Config, error) {
		return cfg, err
	}
}

// mockStartProxyServer replaces interactors.StartProxyServer for testing.
func mockStartProxyServer(err error) func(models.Config) error {
	return func(models.Config) error {
		return err
	}
}

func TestMain_ConfigLoadError(t *testing.T) {
	monkey.Patch(config.Load, mockConfigLoader(models.Config{}, errors.New("load error")))

	calledFatal := false
	logFatalf := func(format string, v ...interface{}) {
		calledFatal = true
	}
	monkey.Patch(log.Fatalf, logFatalf)
	defer func() {
		monkey.UnpatchAll()
	}()

	main()

	if !calledFatal {
		t.Error("expected log.Fatalf to be called on config load error")
	}
}

func TestMain_StartProxyServerError(t *testing.T) {
	monkey.Patch(config.Load, mockConfigLoader(models.Config{}, nil))
	monkey.Patch(interactors.StartProxyServer, mockStartProxyServer(errors.New("server error")))

	calledFatal := false
	logFatalf := func(format string, v ...interface{}) {
		calledFatal = true
	}
	monkey.Patch(log.Fatalf, logFatalf)
	defer func() {
		monkey.UnpatchAll()
	}()

	main()

	if !calledFatal {
		t.Error("expected log.Fatalf to be called on server error")
	}
}

func TestMain_Success(t *testing.T) {
	monkey.Patch(config.Load, mockConfigLoader(models.Config{}, nil))
	monkey.Patch(interactors.StartProxyServer, mockStartProxyServer(nil))

	calledFatal := false
	logFatalf := func(format string, v ...interface{}) {
		calledFatal = true
	}
	monkey.Patch(log.Fatalf, logFatalf)
	defer func() {
		monkey.UnpatchAll()
	}()

	main()

	if calledFatal {
		t.Error("did not expect log.Fatalf to be called on success")
	}
}

// --- Test helpers for log.Fatalf interception ---

var logFatalfFunc = log.Fatalf

func init() {
	logFatalfFunc = log.Fatalf
}
