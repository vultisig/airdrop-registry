package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vultisig/airdrop-registry/config"
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

		}
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down gracefully...")
}
