package main

import (
	"github.com/teagan42/snitchcraft/internal/config"
	"github.com/teagan42/snitchcraft/internal/interactors"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	if err := interactors.StartProxyServer(cfg); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
