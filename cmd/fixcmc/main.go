package main

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/vultisig/airdrop-registry/config"
	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/internal/services"
	"github.com/vultisig/airdrop-registry/internal/tokens"
)

func main() {
	cachedData := cache.New(10*time.Hour, 1*time.Hour)
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	storage, err := services.NewStorage(cfg)
	if err != nil {
		panic(fmt.Errorf("failed to initialize storage: %w", err))
	}
	defer func() {
		if cerr := storage.Close(); cerr != nil {
			fmt.Println("failed to close storage db:", cerr)
		}
	}()
	cmcService, err := tokens.NewCMCService()
	if err != nil {
		panic(fmt.Errorf("failed to initialize CMC service: %w", err))
	}
	const pageSize = 1000
	var currentID uint64
	for {
		coins, err := storage.GetCoinsWithPage(currentID, pageSize)
		if err != nil {
			fmt.Printf("failed to get coins: %v\n", err)
			break
		}
		if len(coins) == 0 {
			fmt.Println("No more coins to process.")
			break
		}
		for _, coin := range coins {
			currentID = uint64(coin.ID)
			if coin.CMCId != 0 {
				continue
			}
			if coin.ContractAddress == "" {
				continue
			}
			cacheKey := fmt.Sprintf("%s_%s", coin.Chain, coin.ContractAddress)
			if cached, found := cachedData.Get(cacheKey); found {
				if cachedInt, ok := cached.(int); ok && cachedInt == -1 {
					//fmt.Printf("CMC ID does not exist for coin: %s (%s)\n", coin.Ticker, cacheKey)
					continue
				}
			}
			cmcID, err := cmcService.GetCMCID(coin.Chain, models.Coin{ContractAddress: coin.ContractAddress})
			if err != nil {
				cachedData.Set(cacheKey, -1, cache.DefaultExpiration)
				fmt.Printf("failed to get CMC ID for coin %s (%s): %v\n", coin.Ticker, cacheKey, err)
				continue
			}
			fmt.Println("Contract Address:", coin.ContractAddress, "CMC ID:", cmcID)
			if err := storage.UpdateCoinCMCIDByID(cmcID, uint64(coin.ID)); err != nil {
				fmt.Printf("failed to update coin CMC ID for %s: %v\n", coin.Ticker, err)
				continue
			}
			fmt.Printf("Updated coin %s with CMC ID %d\n", coin.Ticker, cmcID)
		}
	}
}
