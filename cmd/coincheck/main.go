package main

import (
	_ "embed"

	"github.com/sirupsen/logrus"

	"github.com/vultisig/airdrop-registry/config"
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/internal/services"
	"github.com/vultisig/airdrop-registry/internal/tokens"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:      true,
		FullTimestamp:    true,
		DisableColors:    false,
		DisableTimestamp: false,
		DisableLevelTruncation: true,
	})

	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to load config")
	}

	storage, err := services.NewStorage(cfg)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to initialize storage")
	}
	defer func() {
		if err := storage.Close(); err != nil {
			logrus.WithError(err).Errorf("Failed to close storage")
		}
	}()

	cmcService, err := tokens.NewCMCService()
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to initialize CMC service")
	}
	oneInchService, err := tokens.NewOneInchService()
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to initialize OneInch service")
	}

	discoveryServices := map[common.Chain]tokens.AutoDiscoveryService{
		common.Tron:   tokens.NewTRC20DiscoveryService(common.Tron, cmcService),
		common.Solana: tokens.NewSPLDiscoveryService(cmcService),
	}
	for _, chain := range common.EVMChains {
		//vultisig return not found for Blast and Cronos chain
		//Zksync is no supported
		if chain == common.Blast || chain == common.CronosChain ||chain == common.Zksync {
			continue
		}
		err := oneInchService.LoadOneInchTokens(chain)
		if err != nil {
			logrus.WithError(err).Fatalf("Failed to load oneInch service")
		}
		discoveryServices[chain] = tokens.NewERC20DiscoveryService(oneInchService, cmcService)
	}

	predefinedService := tokens.NewPredefinedTokenService()
	const batchSize = 1000
	var currentID uint64

	for {
		coins, err := storage.GetCoinsWithPage(currentID, batchSize)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to fetch coins")
		}
		if len(coins) == 0 {
			logrus.Infof("No more coins to process")
			break
		}
		currentID = uint64(coins[len(coins)-1].ID)

		for _, coin := range coins {
			if coin.CMCId == 0 || coin.Decimals == 0 {
				logrus.Warnf("Invalid coin data - CMCId: %d, Decimals: %d, Chain: %s, Address: %s",
					coin.CMCId, coin.Decimals, coin.Chain, coin.ContractAddress)
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
					logrus.Warnf("Coin data mismatch - System(CMCId: %d, Decimals: %d) vs User(CMCId: %d, Decimals: %d) for contract address: %s on %s",
						coin.CMCId, coin.Decimals, predefinedCoin.CMCId, predefinedCoin.Decimals, coin.ContractAddress, coin.Chain)
				} else {
					logrus.Infof("Coin data matches in predefined tokens - CMCId: %d, Decimals: %d for %s on %s",
						coin.CMCId, coin.Decimals, coin.ContractAddress, coin.Chain)
				}
			} else {
				discoveryService, exists := discoveryServices[coin.Chain]
				if !exists {
					logrus.Warnf("No discovery service found for chain: %s", coin.Chain)
					continue
				}
				coinData, err := discoveryService.Search(coinBase)
				if err != nil {
					logrus.Errorf("Error searching contract address %s on chain %s: %v", coin.ContractAddress, coin.Chain, err)
				}
				if coinData.CMCId == coin.CMCId && coinData.Decimals == coin.Decimals {
					logrus.Infof("Coin data matches - CMCId: %d, Decimals: %d for %s on %s",
						coin.CMCId, coin.Decimals, coin.ContractAddress, coin.Chain)
				} else {
					logrus.Warnf("Coin data mismatch - System(CMCId: %d, Decimals: %d) vs User(CMCId: %d, Decimals: %d) for contract address: %s on %s",
						coin.CMCId, coin.Decimals, coinData.CMCId, coinData.Decimals, coin.ContractAddress, coin.Chain)
				}
			}
		}
	}
}
