package main

import (
	"log"
	"snitchcraft/internal/config/env"
	"snitchcraft/internal/interactors"
)

func main() {
	cfg, err := env.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	if err := interactors.StartProxyServer(cfg); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
