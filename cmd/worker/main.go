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
	"github.com/vultisig/airdrop-registry/internal/volume"
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
	referralResolver := services.NewReferralResolverService(cfg.ReferralBot.BaseAddress, cfg.ReferralBot.APIKey)
	priceResolver, err := services.NewPriceResolver(cfg)
	if err != nil {
		panic(err)
	}
	balanceResolver, err := balance.NewBalanceResolver()
	if err != nil {
		panic(err)
	}
	volumeTracker, err := volume.NewVolumeResolver(cfg)
	if err != nil {
		panic(err)
	}
	pointWorker, err := services.NewPointWorker(cfg, storage, priceResolver, balanceResolver, volumeTracker, referralResolver)
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
