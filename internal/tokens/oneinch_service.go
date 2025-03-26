package tokens

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
)

type token struct {
	Address  string `json:"address"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
	Name     string `json:"name"`
	LogoURI  string `json:"logoURI"`
}

type tokensResponse struct {
	Tokens map[string]token `json:"tokens"`
}

var chainIDs map[common.Chain]int = map[common.Chain]int{
	common.Ethereum:    1,
	common.Avalanche:   43114,
	common.Base:        8453,
	common.Blast:       81457,
	common.Arbitrum:    42161,
	common.Polygon:     137,
	common.Optimism:    10,
	common.BscChain:    56,
	common.CronosChain: 25,
}

type oneinchService struct {
	logger         *logrus.Logger
	oneinchBaseURL string
	cachedData     *cache.Cache
	coinBase       []models.CoinBase
}

func NewOneinchService() *oneinchService {
	return &oneinchService{
		logger:         logrus.WithField("module", "oneinch_service").Logger,
		oneinchBaseURL: "https://api.vultisig.com/1inch",
		cachedData:     cache.New(10*time.Hour, 10*time.Hour),
		coinBase:       []models.CoinBase{},
	}
}

func (o *oneinchService) Load1inchTokens(chain common.Chain) ([]models.Coin, error) {
	if _, ok := chainIDs[chain]; !ok {
		return nil, fmt.Errorf("chain: %s is not supported", chain)
	}
	if cachedData, found := o.cachedData.Get(chain.String()); found {
		if coins, ok := cachedData.([]models.Coin); ok {
			return coins, nil
		}
	}
	url := fmt.Sprintf("%s/swap/v6.0/%d/tokens", o.oneinchBaseURL, chainIDs[chain])
	resp, err := http.Get(url)
	if err != nil {
		o.logger.Error(err)
		return nil, fmt.Errorf("fail to get tokens from, err %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching tokens from %s: %s", url, resp.Status)
	}
	var tokensResponse tokensResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokensResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	coins := make([]models.Coin, 0, len(tokensResponse.Tokens))
	for _, token := range tokensResponse.Tokens {
		isNative := false
		if token.Address == ethereum {
			isNative = true
		}
		coins = append(coins, models.Coin{
			Ticker:          token.Symbol,
			ContractAddress: token.Address,
			Decimals:        token.Decimals,
			IsNative:        isNative,
			Logo:            token.LogoURI,
		})
	}

	o.cachedData.Set(chain.String(), coins, 10*time.Hour)
	return coins, nil
}

func (o *oneinchService) GetTokenDetailsByContract(chain, contract string) (models.CoinBase, error) {
	key := chain + contract
	if cachedData, found := o.cachedData.Get(key); found {
		if coin, ok := cachedData.(models.CoinBase); ok {
			return coin, nil
		}
	}
	url := fmt.Sprintf("%s/token-details/v1.0/details/%s/%s", o.oneinchBaseURL, chain, contract)
	resp, err := http.Get(url)
	if err != nil {
		o.logger.WithFields(logrus.Fields{
			"error":        err,
			"url":          url,
			"contractAddr": contract,
		}).Error("Failed to fetch token details")
		return models.CoinBase{}, fmt.Errorf("failed to fetch token details: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return models.CoinBase{}, nil
	}
	if resp.StatusCode != http.StatusOK {
		o.logger.WithFields(logrus.Fields{
			"statusCode":   resp.StatusCode,
			"url":          url,
			"contractAddr": contract,
		}).Error("1inch token details API request failed")
		return models.CoinBase{}, fmt.Errorf("1inch token details API request failed with status: %d", resp.StatusCode)
	}
	var tokenData tokenData
	if err := json.NewDecoder(resp.Body).Decode(&tokenData); err != nil {
		o.logger.WithFields(logrus.Fields{
			"error":        err,
			"contractAddr": contract,
		}).Error("Failed to decode token details response")
		return models.CoinBase{}, fmt.Errorf("failed to decode token details response: %w", err)
	}
	coin := models.CoinBase{
		Decimals:        tokenData.Assets.Decimals,
		Ticker:          tokenData.Assets.Symbol,
		ContractAddress: contract,
	}
	o.logger.WithFields(logrus.Fields{
		"contractAddr": coin.ContractAddress,
		"decimals":     tokenData.Assets.Decimals,
		"tokenName":    tokenData.Assets.Name,
	}).Debug("Token details retrieved successfully")
	o.cachedData.Set(key, coin, 10*time.Hour)
	return coin, nil
}
