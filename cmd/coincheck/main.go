package main

import (
	"log"

	_ "embed"

	"github.com/vultisig/airdrop-registry/config"
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/internal/services"
	"github.com/vultisig/airdrop-registry/internal/tokens"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	storage, err := services.NewStorage(cfg)
	if err != nil {
		log.Fatalf("failed to initialize storage: %v", err)
	}

	cmcService, err := tokens.NewCMCService()
	if err != nil {
		storage.Close()
		log.Fatalf("failed to initialize CMC service: %v", err)
	}
	oneInchService := tokens.NewOneInchService()

	// Initialize discovery services map
	discoveryServices := map[common.Chain]tokens.AutoDiscoveryService{
		// Native blockchain-specific implementations
		common.Tron:   tokens.NewTRC20DiscoveryService(common.Tron, cmcService),
		common.Solana: tokens.NewSPLDiscoveryService(cmcService),
	}

	// Add EVM-compatible chains
	for _, chain := range common.EVMChains {
		discoveryServices[chain] = tokens.NewERC20DiscoveryService(oneInchService, cmcService)
	}

	defer func() {
		if err := storage.Close(); err != nil {
			log.Printf("Failed to close storage: %v", err)
		}
	}()

	const batchSize = 1000
	var currentID uint64
	predefinedService := tokens.NewPredefinedTokenService()
	for {
		coins, err := storage.GetCoinsWithPage(currentID, batchSize)
		if err != nil {
			log.Fatalf("Failed to fetch coins: %v", err)
		}
		if len(coins) == 0 {
			log.Println("no more coins to process")
			break
		}
		currentID = uint64(coins[len(coins)-1].ID)
		for _, coin := range coins {

			if coin.CMCId == 0 || coin.Decimals == 0 {
				log.Printf("Invalid coin data - CMCId: %d, Decimals: %d, Chain: %s, Address: %s \n\n",
					coin.CMCId, coin.Decimals, coin.Chain, coin.ContractAddress)

				continue
			}

			// check in predefined tokens
			predefinedCoin, err := predefinedService.Search(models.CoinBase{
				Chain:           coin.Chain,
				Address:         coin.Address,
				ContractAddress: coin.ContractAddress,
			},
			)
			if err == nil {
				if predefinedCoin.CMCId != coin.CMCId || predefinedCoin.Decimals != coin.Decimals {
					log.Printf("mismatch found in predefined tokens - System(CMCId: %d, Decimals: %d) vs User(CMCId: %d, Decimals: %d) for %s on %s \n\n",
						coin.CMCId, coin.Decimals, predefinedCoin.CMCId, predefinedCoin.Decimals, coin.ContractAddress, coin.Chain)
				} else {
					log.Printf("Coin data matches in predefined tokens - CMCId: %d, Decimals: %d for %s on %s \n\n",
						coin.CMCId, coin.Decimals, coin.ContractAddress, coin.Chain)
				}
				continue
			} else {
				log.Printf("Coin not found in predefined tokens - CMCId: %d, Decimals: %d for %s on %s \n\n",
					coin.CMCId, coin.Decimals, coin.ContractAddress, coin.Chain)
			}

			discoveryService, exists := discoveryServices[coin.Chain]
			if !exists {
				// set colour to yellow

				log.Printf("No discovery service found for chain: %s\n\n", coin.Chain)
				continue
			}
			// if coin.ContractAddress == "" {
			// 	coin.ContractAddress = coin.Address
			// }
			coinData, err := discoveryService.Search(models.CoinBase{
				Chain:           coin.Chain,
				Address:         coin.Address,
				ContractAddress: coin.ContractAddress,
			})
			if err != nil {
				log.Printf("Error searching coin %s on chain %s: %v \n\n", coin.ContractAddress, coin.Chain, err)
				continue
			}
			if coinData.CMCId != coin.CMCId || coinData.Decimals != coin.Decimals {
				log.Printf("mismatch found - System(CMCId: %d, Decimals: %d) vs User(CMCId: %d, Decimals: %d) for %s on %s \n\n",
					coin.CMCId, coin.Decimals, coinData.CMCId, coinData.Decimals, coin.ContractAddress, coin.Chain)
			}

			if coinData.CMCId == coin.CMCId && coinData.Decimals == coin.Decimals {
				log.Printf("Coin data matches - CMCId: %d, Decimals: %d for %s on %s \n\n",
					coin.CMCId, coin.Decimals, coin.ContractAddress, coin.Chain)
			} else {
				log.Printf("Coin data mismatch - System(CMCId: %d, Decimals: %d) vs User(CMCId: %d, Decimals: %d) for %s on %s \n\n",
					coin.CMCId, coin.Decimals, coinData.CMCId, coinData.Decimals, coin.ContractAddress, coin.Chain)
			}
		}
	}
}
