package main

import (
	"log"

	"table_collab/cmd/server/config"
	"table_collab/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	srv := server.New(cfg)

	log.Printf("🚀 Starting TableCollab on %s", cfg.Server.Address)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
