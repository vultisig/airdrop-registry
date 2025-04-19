package tokens

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	lru "github.com/hashicorp/golang-lru/v2"
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
	common.Ethereum:  1,
	common.Avalanche: 43114,
	common.Base:      8453,
	//common.Blast:       81457,
	common.Arbitrum: 42161,
	common.Polygon:  137,
	common.Optimism: 10,
	common.BscChain: 56,
	//common.CronosChain: 25,
}

type oneInchService struct {
	logger         *logrus.Logger
	oneInchBaseURL string
	cachedData     *lru.Cache[string, models.CoinBase]
	coinBase       []models.CoinBase
}

func NewOneInchService() (*oneInchService, error) {

	cache, err := lru.New[string, models.CoinBase](20000)
	if err != nil {
		return nil,fmt.Errorf("failed to create LRU cache")
	}
	return &oneInchService{
		logger:         logrus.WithField("module", "oneinch_service").Logger,
		oneInchBaseURL: "https://api.vultisig.com/1inch",
		cachedData:     cache,
		coinBase:       []models.CoinBase{},
	}, nil
}
func (o *oneInchService) IsChainSupported(chain common.Chain) bool {
	if _, ok := chainIDs[chain]; ok {
		return true
	}
	return false
}

func (o *oneInchService) LoadOneInchTokens(chain common.Chain) error {
	if _, ok := chainIDs[chain]; !ok {
		return fmt.Errorf("chain: %s is not supported", chain)
	}

	keys := o.cachedData.Keys()
	chainPrefix := chain.String() + "_"
	for _, key := range keys {
		if strings.HasPrefix(key, chainPrefix) {
			return nil 
		}
	}
	url := fmt.Sprintf("%s/swap/v6.0/%d/tokens", o.oneInchBaseURL, chainIDs[chain])
	resp, err := http.Get(url)
	if err != nil {
		o.logger.Error(err)
		return fmt.Errorf("fail to get tokens from, err %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error fetching tokens from %s: %s", url, resp.Status)
	}
	var tokensResponse tokensResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokensResponse); err != nil {
		return fmt.Errorf("error unmarshalling response: %w", err)
	}
	for _, token := range tokensResponse.Tokens {
		isNative := false
		if token.Address == ethereum {
			isNative = true
		}
		cacheKey := o.getCacheKey(chain.String(), token.Address)
		o.cachedData.Add(cacheKey, models.CoinBase{
			Ticker:          token.Symbol,
			ContractAddress: token.Address,
			Decimals:        token.Decimals,
			IsNative:        isNative,
			Logo:            token.LogoURI,
		})
	}
	return nil
}

func (o *oneInchService) GetTokenDetailsByContract(chain common.Chain, contract string) (models.CoinBase, error) {
	chainID, ok := chainIDs[chain]
	if !ok {
		return models.CoinBase{}, fmt.Errorf("chain: %s is not supported", chain)
	}
	cacheKey := o.getCacheKey(chain.String(), contract)
	if cachedData, found := o.cachedData.Get(cacheKey); found {
		return cachedData, nil
	}
	url := fmt.Sprintf("%s/token-details/v1.0/details/%d/%s", o.oneInchBaseURL, chainID, contract)
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
	o.cachedData.Add(cacheKey, coin)
	return coin, nil
}

func (o *oneInchService) getCacheKey(chain, contract string) string {
	return fmt.Sprintf("%s_%s", chain, contract)
}
