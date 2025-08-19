package main

import (
	"log"

	"github.com/vultisig/airdrop-registry/config"
	"github.com/vultisig/airdrop-registry/internal/handlers"
	"github.com/vultisig/airdrop-registry/internal/services"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	storage, err := services.NewStorage(cfg)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := storage.Close(); err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}()
	api, err := handlers.NewApi(cfg, storage)
	if err != nil {
		panic(err)
	}
	if err := api.Start(); err != nil {
		panic(err)
	}
}
