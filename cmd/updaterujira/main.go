package main

import (
	"fmt"

	"github.com/vultisig/airdrop-registry/config"
	"github.com/vultisig/airdrop-registry/internal/balance"
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
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
	balanceResolver, err := balance.NewBalanceResolver()
	if err != nil {
		panic(err)
	}
	startId := uint64(0)
	rujiCoins := make([]models.CoinDBModel, 0)
	for {
		// fetch all coins
		coins, err := storage.GetCoinsWithPage(startId, 200)
		if err != nil {
			fmt.Println("fail to get coins: ", err)
			return
		}
		if len(coins) == 0 {
			break
		}
		for _, coin := range coins {
			startId = uint64(coin.ID)
			if coin.Chain == common.THORChain {
				if coin.Ticker == "RUJIRA" {
					rujiCoins = append(rujiCoins, coin)
				}
			}
		}
	}
	fmt.Println("Found ", len(rujiCoins), " RUJIRA coins")
	cnt := 0
	balances := make([]float64, 0)
	for _, coin := range rujiCoins {
		balance, err := balanceResolver.GetBalanceWithRetry(coin)
		if err != nil {
			fmt.Printf("fail to get balance for coin %d: %v\n", coin.ID, err)
			return
		}
		if balance > 0 {
			cnt++
			fmt.Printf("Address: %s, Balance: %f\n", coin.Address, balance)
		}
		balances = append(balances, balance)
	}
	if len(balances) != len(rujiCoins) {
		fmt.Println("Error: balances length does not match coins length")
		return
	}
	// update coin balance in db
	for i, coin := range rujiCoins {
		if balances[i] > 0 {
			coin.Balance = fmt.Sprintf("%.8f", balances[i])
			if err := storage.UpdateCoinBalance(uint64(coin.ID), balances[i]); err != nil {
				panic(fmt.Errorf("failed to update coin balance for %d: %w", coin.ID, err))
			}
		}
	}
	if err := storage.UpdateVaultBalance(); err != nil {
		panic(fmt.Errorf("failed to update vault balance: %w", err))
	}
	fmt.Printf("Total coins with balance > 0: %d\n", cnt)
	fmt.Println("All balances updated successfully")
}
