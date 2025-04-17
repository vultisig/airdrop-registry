package main

import (
	log "github.com/sirupsen/logrus"

	_ "embed"

	"github.com/vultisig/airdrop-registry/config"
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/internal/services"
	"github.com/vultisig/airdrop-registry/internal/tokens"
)

const yellow = "\033[33m"
const red = "\033[31m"
const green = "\033[32m"
const reset = "\033[0m"

func main() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:      true,
		FullTimestamp:    true,
		DisableColors:    false,
		DisableTimestamp: false,
	})
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("%s[FATAL] failed to load config: %v%s", red, err, reset)
	}

	storage, err := services.NewStorage(cfg)
	if err != nil {
		log.Fatalf("%s[FATAL] failed to initialize storage: %v%s", red, err, reset)
	}
	defer func() {
		if err := storage.Close(); err != nil {
			log.Errorf("%s[ERROR] Failed to close storage: %v%s", red, err, reset)
		}
	}()

	cmcService, err := tokens.NewCMCService()
	if err != nil {
		log.Fatalf("%s[FATAL] failed to initialize CMC service: %v%s", red, err, reset)
	}
	oneInchService, err := tokens.NewOneInchService()
	if err != nil {
		log.Errorf("%s[ERROR] Failed to create OneInch service: %v%s", red, err, reset)
	}

	discoveryServices := map[common.Chain]tokens.AutoDiscoveryService{
		common.Tron:   tokens.NewTRC20DiscoveryService(common.Tron, cmcService),
		common.Solana: tokens.NewSPLDiscoveryService(cmcService),
	}
	for _, chain := range common.EVMChains {
		err := oneInchService.LoadOneInchTokens(chain)
		if err != nil {
			log.Fatalf("%s[FATAL] failed to load oneInch service: %v%s", red, err, reset)
		}
		discoveryServices[chain] = tokens.NewERC20DiscoveryService(oneInchService, cmcService)
	}

	predefinedService := tokens.NewPredefinedTokenService()
	const batchSize = 1000
	var currentID uint64

	for {
		coins, err := storage.GetCoinsWithPage(currentID, batchSize)
		if err != nil {
			log.Errorf("%s[Error] Failed to fetch coins: %v%s", red, err, reset)
		}
		if len(coins) == 0 {
			log.Infof("%s[INFO] no more coins to process%s", green, reset)
			break
		}
		currentID = uint64(coins[len(coins)-1].ID)

		for _, coin := range coins {
			if coin.CMCId == 0 || coin.Decimals == 0 {
				log.Warnf("%s[WARN] Invalid coin data - CMCId: %d, Decimals: %d, Chain: %s, Address: %s%s",
					yellow, coin.CMCId, coin.Decimals, coin.Chain, coin.ContractAddress, reset)
				continue
			}

			coinBase := models.CoinBase{
				Chain:           coin.Chain,
				Address:         coin.Address,
				ContractAddress: coin.ContractAddress,
			}

			predefinedCoin, err := predefinedService.Search(coinBase)
			if err == nil {
				if predefinedCoin.CMCId != coin.CMCId || predefinedCoin.Decimals != coin.Decimals {
					log.Warnf("%s[WARN] Coin data mismatch - System(CMCId: %d, Decimals: %d) vs User(CMCId: %d, Decimals: %d) for contract address: %s on %s%s",
						yellow, coin.CMCId, coin.Decimals, predefinedCoin.CMCId, predefinedCoin.Decimals, coin.ContractAddress, coin.Chain, reset)
				} else {
					log.Infof("%s[INFO] Coin data matches in predefined tokens - CMCId: %d, Decimals: %d for %s on %s%s",
						green, coin.CMCId, coin.Decimals, coin.ContractAddress, coin.Chain, reset)
				}
			} else {
				discoveryService, exists := discoveryServices[coin.Chain]
				if !exists {
					log.Warnf("%s[WARN] No discovery service found for chain: %s%s", yellow, coin.Chain, reset)
					continue
				}
				coinData, err := discoveryService.Search(coinBase)
				if err != nil {
					log.Errorf("%s[ERROR] Error searching contract address %s on chain %s: %v%s", red, coin.ContractAddress, coin.Chain, err, reset)
				}
				if coinData.CMCId == coin.CMCId && coinData.Decimals == coin.Decimals {
					log.Infof("%s[INFO] Coin data matches - CMCId: %d, Decimals: %d for %s on %s%s",
						green, coin.CMCId, coin.Decimals, coin.ContractAddress, coin.Chain, reset)
				} else {
					log.Warnf("%s[WARN] Coin data mismatch - System(CMCId: %d, Decimals: %d) vs User(CMCId: %d, Decimals: %d) for contract address: %s on %s%s",
						yellow, coin.CMCId, coin.Decimals, coinData.CMCId, coinData.Decimals, coin.ContractAddress, coin.Chain, reset)
				}
			}
		}
	}
}
