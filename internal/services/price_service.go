package services

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"

	"github.com/vultisig/airdrop-registry/config"
	"github.com/vultisig/airdrop-registry/internal/models"
)

const CMC_Base_URL = "https://api.vultisig.com/cmc/"

type PriceResolver struct {
	logger               *logrus.Logger
	cmcMap               *CmcMapResp
	lifiBaseAddress      string
	coingeckoBaseAddress string
	priceCache           cache.Cache
	OpenSeaAPIKey        string
}

func NewPriceResolver(cfg *config.Config) (*PriceResolver, error) {
	pr := &PriceResolver{
		logger:               logrus.WithField("module", "price_resolver").Logger,
		lifiBaseAddress:      "https://li.quest",
		coingeckoBaseAddress: "https://api.vultisig.com/coingeicko/api/v3/simple/price",
		priceCache:           *cache.New(4*time.Hour, 5*time.Hour),
		OpenSeaAPIKey:        cfg.OpenSea.APIKey,
	}
	result, err := pr.getCMCMap()
	if err != nil {
		return nil, fmt.Errorf("fail to get CMC map,err: %w", err)
	}
	pr.cmcMap = result
	return pr, nil
}

type CmcMapItemPlatform struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	Slug         string `json:"slug"`
	TokenAddress string `json:"token_address"`
}
type CmcMapItem struct {
	ID       int                 `json:"id"`
	Name     string              `json:"name"`
	Symbol   string              `json:"symbol"`
	Slug     string              `json:"slug"`
	IsActive int                 `json:"is_active"`
	Platform *CmcMapItemPlatform `json:"platform"`
}
type CmcMapResp struct {
	Data []CmcMapItem `json:"data"`
}

func (p *PriceResolver) getCMCMap() (*CmcMapResp, error) {
	url := CMC_Base_URL + "/v1/cryptocurrency/map"
	resp, err := http.Get(url)
	if err != nil {
		p.logger.Error(err)
		return nil, fmt.Errorf("fail to get map from CMC,err: %w", err)
	}
	defer p.closer(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching CMC map: %s", resp.Status)
	}
	var cmcMapResp CmcMapResp
	if err := json.NewDecoder(resp.Body).Decode(&cmcMapResp); err != nil {
		return nil, fmt.Errorf("error decoding CMC map response: %w", err)
	}
	return &cmcMapResp, nil
}
func (p *PriceResolver) closer(closer io.Closer) {
	if err := closer.Close(); err != nil {
		p.logger.Error(err)
	}
}

func (p *PriceResolver) resolveIds(coinIds []models.CoinIdentity) string {
	var ids []string
	for _, coinId := range coinIds {
		ids = append(ids, strconv.Itoa(coinId.CMCId))
		//for _, item := range p.cmcMap.Data {
		//	if strings.EqualFold(item.Symbol, coinId.Ticker) &&
		//		strings.EqualFold(item.Name, coinId.Chain.String()) {
		//		ids = append(ids, strconv.Itoa(item.ID))
		//	}
		//}
	}
	return strings.Join(ids, ",")
}
func (p *PriceResolver) GetCoinGeckoPrice(priceProviderId string, currency string) (float64, error) {
	cacheKey := fmt.Sprintf("cg_%s_%s", priceProviderId, currency)
	if cachedPrice, ok := p.priceCache.Get(cacheKey); ok {
		return cachedPrice.(float64), nil
	}
	url := fmt.Sprintf("%s?ids=%s&vs_currencies=%s", p.coingeckoBaseAddress, priceProviderId, currency)
	resp, err := http.Get(url)
	if err != nil {
		p.logger.Error(err)
		return 0, fmt.Errorf("fail to get price from CoinGecko,err: %w", err)
	}
	defer p.closer(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error fetching CoinGecko price: %s", resp.Status)
	}
	var result map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return 0, fmt.Errorf("error decoding CoinGecko price response: %w", err)
	}
	if _, ok := result[priceProviderId]; ok {
		if _, ok := result[priceProviderId][currency]; ok {
			p.priceCache.Set(cacheKey, result[priceProviderId][currency], cache.DefaultExpiration)
			return result[priceProviderId][currency], nil
		}
	}
	return 0, fmt.Errorf("price not found in response")
}

func (p *PriceResolver) GetLiFiPrice(chain, contractAddress string) (float64, error) {
	url := fmt.Sprintf("%s/v1/token?chain=%s&token=%s", p.lifiBaseAddress, chain, contractAddress)
	resp, err := http.Get(url)
	if err != nil {
		p.logger.Error(err)
		return 0, fmt.Errorf("fail to get price from LiQuest,err: %w", err)
	}
	defer p.closer(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error fetching LiQuest price: %s", resp.Status)
	}
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		p.logger.Errorf("Error parsing JSON: %s", err)
		return 0, fmt.Errorf("error decoding LiQuest price response: %w", err)
	}
	if _, ok := result["priceUSD"]; !ok {
		return 0, fmt.Errorf("priceUSD not found in response")
	}
	if _, ok := result["priceUSD"].(string); !ok {
		return 0, fmt.Errorf("priceUSD is not string")
	}
	//convert "0.45" to float64
	strPrice := result["priceUSD"].(string)
	price, err := strconv.ParseFloat(strPrice, 32)
	if err != nil {
		return 0, fmt.Errorf("error parsing price: %w", err)
	}
	return price, nil
}
func (p *PriceResolver) GetMidgardCacaoPrices() (float64, error) {
	if cachedPrice, ok := p.priceCache.Get("midgard_cacao"); ok {
		return cachedPrice.(float64), nil
	}
	// fetch from https://midgard.mayachain.info/v2/debug/usd
	resp, err := http.Get("https://midgard.mayachain.info/v2/debug/usd")
	if err != nil {
		p.logger.Error(err)
		return 0, fmt.Errorf("fail to get price from Midgard,err: %w", err)
	}
	defer p.closer(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error fetching Midgard price: %s", resp.Status)
	}
	/*
			sample response:
			ETH.USDT-0XDAC17F958D2EE523A2206206994597C13D831EC7 - originalDepth: 11302361418928360 runeDepth: 1130236 assetDepth: 70335388460771 cacaoPriceUsd: 0.62
		ETH.USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48 - originalDepth: 7207733475181527 runeDepth: 720773 assetDepth: 44619085750427 cacaoPriceUsd: 0.62
		cacaoPriceUSD: 0.6223070193364945
	*/
	//parse the response
	str, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response: %w", err)
	}
	lines := strings.Split(string(str), "\n")
	for _, line := range lines {
		if strings.Contains(line, "cacaoPriceUSD") {
			priceStr := strings.Split(line, ":")[1]
			priceStr = strings.TrimSpace(priceStr)
			price, err := strconv.ParseFloat(priceStr, 64)
			if err != nil {
				return 0, fmt.Errorf("error parsing price: %w", err)
			}
			p.priceCache.Set("midgard_cacao", price, 4*time.Hour)
			return price, nil
		}
	}
	return 0, fmt.Errorf("price not found in response")
}
func (p *PriceResolver) GetAllTokenPrices(coinIds []models.CoinIdentity) (map[int]float64, error) {
	strIds := p.resolveIds(coinIds)
	url := CMC_Base_URL + "/v2/cryptocurrency/quotes/latest?id=" + strIds
	resp, err := http.Get(url)
	if err != nil {
		p.logger.Error(err)
		return nil, fmt.Errorf("fail to get prices from CMC,err: %w", err)
	}
	defer p.closer(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching CMC prices: %s", resp.Status)
	}
	type CmcQuoteResp struct {
		Data map[string]struct {
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Symbol string `json:"symbol"`
			Slug   string `json:"slug"`
			Quote  struct {
				USD struct {
					Price float64 `json:"price"`
				} `json:"USD"`
			} `json:"quote"`
		} `json:"data"`
	}
	var cmcQuoteResp CmcQuoteResp
	if err := json.NewDecoder(resp.Body).Decode(&cmcQuoteResp); err != nil {
		return nil, fmt.Errorf("error decoding CMC quote response: %w", err)
	}
	priceMap := make(map[int]float64)
	for _, item := range cmcQuoteResp.Data {
		priceMap[item.ID] = item.Quote.USD.Price
	}
	return priceMap, nil
}

type OpenSeaBestCollectionResponse struct {
	Listings []struct {
		Price struct {
			Current struct {
				Currency string `json:"currency"`
				Decimals int    `json:"decimals"`
				Value    string `json:"value"`
			} `json:"current"`
		} `json:"price"`
	} `json:"listings"`
}

func (p *PriceResolver) GetOpenSeaCollectionMinPrice(collectionSlug string) (float64, error) {
	key := fmt.Sprintf("opensea_%s", collectionSlug)
	//check cache first
	if cached, ok := p.priceCache.Get(key); ok {
		if price, ok := cached.(float64); ok {
			return price, nil
		}
	}

	url := fmt.Sprintf("https://api.opensea.io/api/v2/listings/collection/%s/best", collectionSlug)
	// add x-api-key header
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get collection from OpenSea,err: %w", err)
	}
	req.Header.Add("x-api-key", p.OpenSeaAPIKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to get collection from OpenSea,err: %w", err)
	}
	defer p.closer(resp.Body)
	var openseaResp OpenSeaBestCollectionResponse
	if err := json.NewDecoder(resp.Body).Decode(&openseaResp); err != nil {
		return 0, fmt.Errorf("failed to decode response from OpenSea,err: %w", err)
	}
	if len(openseaResp.Listings) == 0 {
		return 0, fmt.Errorf("no listing found in response")
	}
	if !strings.EqualFold(openseaResp.Listings[0].Price.Current.Currency, "ETH") {
		return 0, fmt.Errorf("currency not ETH")
	}
	rawValue := openseaResp.Listings[0].Price.Current.Value
	valueBigInt, ok := new(big.Int).SetString(rawValue, 10)
	var pricePerEth float64
	if !ok {
		// Handle error, maybe set pricePerEth to 0 or return an error
		pricePerEth = 0.0
	} else {
		// Convert to float64 for division
		valueFloat := new(big.Float).SetInt(valueBigInt)
		ethDiv := big.NewFloat(1e18)
		pricePerEthFloat, _ := new(big.Float).Quo(valueFloat, ethDiv).Float64()
		// pricePerEthFloat is your value in ETH
		pricePerEth = pricePerEthFloat
		// Use pricePerEth as needed
	}

	//pricePerEth := float64(openseaResp.Listings[0].Price.Current.Value) / 1e18
	priceMap, err := p.GetAllTokenPrices([]models.CoinIdentity{
		models.CoinIdentity{
			CMCId: 1027,
		},
	})
	if err != nil {
		return 0, fmt.Errorf("fail to resolve ids,err: %w", err)
	}
	if price, ok := priceMap[1027]; ok {
		//add to cache
		p.priceCache.Add(key, price*pricePerEth, 3*time.Hour)
		return price * pricePerEth, nil
	}
	return 0, fmt.Errorf("ETH price not found in response")
}
