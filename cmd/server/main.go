package main

import (
	"log"

	"github.com/teagan42/snitchcraft/internal/config"
	"github.com/teagan42/snitchcraft/internal/interactors"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
		return
	}

	if err := interactors.StartProxyServer(cfg); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
