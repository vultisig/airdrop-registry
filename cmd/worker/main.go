package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vultisig/airdrop-registry/config"
	"github.com/vultisig/airdrop-registry/internal/balance"
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
			fmt.Println("fail to close storage db: ", err)
		}
	}()
	priceResolver, err := services.NewPriceResolver()
	if err != nil {
		panic(err)
	}
	balanceResolver, err := balance.NewBalanceResolver()
	if err != nil {
		panic(err)
	}
	// get all thorchain bond providers
	if err := balanceResolver.GetTHORChainBondProviders(); err != nil {
		panic(err)
	}
	pointWorker, err := services.NewPointWorker(cfg, storage, priceResolver, balanceResolver)
	if err != nil {
		panic(err)
	}
	if err := pointWorker.Run(); err != nil {
		panic(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	pointWorker.Stop()
	log.Println("Shutting down gracefully...")
}
