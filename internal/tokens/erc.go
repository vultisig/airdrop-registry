package tokens

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
)

const ethereum string = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"

type (
	tokenData struct {
		Assets assets `json:"assets"`
	}
	assets struct {
		Name     string `json:"name"`
		Symbol   string `json:"symbol"`
		Decimals int    `json:"decimals"`
		Status   string `json:"status"`
		ID       string `json:"id"`
	}
)

type ercDiscoveryService struct {
	logger         *logrus.Logger
	baseAddress    string
	cmcService     *CMCService
	oneInchService *oneInchService
}

func NewERC20DiscoveryService(oneInchService *oneInchService, cmcService *CMCService) AutoDiscoveryService {
	return &ercDiscoveryService{
		logger:         logrus.WithField("module", "oneInch_evm_base_service").Logger,
		baseAddress:    "https://api.vultisig.com/1inch",
		cmcService:     cmcService,
		oneInchService: oneInchService,
	}
}

func (e *ercDiscoveryService) Discover(address string, chain common.Chain) ([]models.CoinBase, error) {
	// Validate inputs
	if address == "" {
		return nil, fmt.Errorf("empty address provided")
	}

	chainID, ok := chainIDs[chain]
	if !ok {
		return nil, fmt.Errorf("unsupported chain: %v", chain)
	}
	url := fmt.Sprintf("%s/balance/v1.2/%d/balances/%s", e.baseAddress, chainID, address)
	resp, err := http.Get(url)
	if err != nil {
		e.logger.WithFields(logrus.Fields{
			"error": err,
			"url":   url,
			"chain": chain,
		}).Error("Failed to fetch account balances")
		return nil, fmt.Errorf("failed to fetch balances: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		e.logger.WithFields(logrus.Fields{
			"statusCode": resp.StatusCode,
			"url":        url,
			"chain":      chain,
		}).Error("API request failed")
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	accounts := make(map[string]string)
	if err := json.NewDecoder(resp.Body).Decode(&accounts); err != nil {
		e.logger.WithError(err).Error("Failed to decode response")
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	coins := make([]models.CoinBase, 0)
	for contract, balance := range accounts {
		if balance == "0" {
			continue
		}
		if contract == ethereum {
			continue
		}
		coins = append(coins, models.CoinBase{
			Address:         address,
			Balance:         balance,
			Chain:           chain,
			ContractAddress: contract,
		})
	}
	// Check if coinBase is nil
	if len(coins) == 0 {
		e.logger.Debug("No tokens found with non-zero balance")
		return coins, nil
	}

	for i, coin := range coins {
		tokenDetails, err := e.oneInchService.GetTokenDetailsByContract(chain.String(), coin.ContractAddress)
		if err != nil {
			e.logger.WithError(err).Error("Failed to fetch token details")
			return nil, fmt.Errorf("failed to fetch token details: %w", err)
		}

		cmcId, err := e.cmcService.GetCMCIDByContract(chain.String(), coin.ContractAddress)
		if err != nil {
			e.logger.WithError(err).Error("Failed to fetch cmc id")
			return nil, fmt.Errorf("failed to fetch cmc id: %w", err)
		}

		coins[i].Decimals = tokenDetails.Decimals
		coins[i].Ticker = tokenDetails.Ticker
		coins[i].ContractAddress = tokenDetails.ContractAddress
		coins[i].CMCId = cmcId
	}
	return coins, nil
}

func (e *ercDiscoveryService) Search(coin models.CoinBase) (models.CoinBase, error) {
	if err != nil {
		e.logger.WithError(err).Error("Failed to fetch cmc id")
		return models.CoinBase{}, fmt.Errorf("failed to fetch cmc id: %w", err)
	}
	oneInchCoin, err := e.oneInchService.GetTokenDetailsByContract(coin.Chain.String(), coin.ContractAddress)
	if err != nil {
		e.logger.WithError(err).Error("Failed to fetch token details")
		return models.CoinBase{}, fmt.Errorf("failed to fetch token details: %w", err)
	}

	coin.Ticker = oneInchCoin.Ticker
	coin.Decimals = oneInchCoin.Decimals
	coin.PriceProviderID = oneInchCoin.PriceProviderID
	coin.Logo = oneInchCoin.Logo
	coin.CMCId = cmcId
	return coin, nil
}
