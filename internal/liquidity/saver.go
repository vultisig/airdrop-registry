package liquidity

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

type SaverPositionResolver struct {
	logger            *logrus.Logger
	thorwalletBaseURL string
	poolCache         *cache.Cache
}

func NewSaverPositionResolver() *SaverPositionResolver {
	//since we fetch vault positions once a day, we can cache the pool data for 4 hours
	poolCache := cache.New(4*time.Hour, 4*time.Hour)
	return &SaverPositionResolver{
		logger:            logrus.WithField("module", "saver_position_resolver").Logger,
		thorwalletBaseURL: "https://api-v2-prod.thorwallet.org",
		poolCache:         poolCache,
	}
}

func (l *SaverPositionResolver) GetSaverPosition(address string) (float64, error) {
	positions, err := l.fetchSaverPosition(address)
	if err != nil {
		return 0, err
	}
	var totalLiquidity float64
	for _, v := range positions.SaverPosition {
		pool, err := l.getpool(v.Pool)
		if err != nil {
			return 0, err
		}
		totalLiquidity += v.AssetRedeem * pool.AssetPriceUsd
	}
	return totalLiquidity * 1e-8, nil
}

type saverResponse struct {
	SaverPosition []struct {
		AssetRedeem float64 `json:"assetRedeem,string"`
		Pool        string  `json:"pool"`
	} `json:"pools"`
}

func (l *SaverPositionResolver) fetchSaverPosition(address string) (saverResponse, error) {
	if address == "" {
		return saverResponse{}, fmt.Errorf("address cannot be empty")
	}
	url := fmt.Sprintf("%s/saver/positions?addresses=%s", l.thorwalletBaseURL, address)
	resp, err := http.Get(url)
	if err != nil {
		l.logger.Errorf("error fetching saver position from %s: %e", url, err)
		return saverResponse{}, fmt.Errorf("error fetching saver position from %s: %e", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		l.logger.Errorf("error fetching saver position from %s: %s", url, resp.Status)
		return saverResponse{}, fmt.Errorf("error fetching saver position from %s: %s", url, resp.Status)
	}
	var positions saverResponse
	if err := json.NewDecoder(resp.Body).Decode(&positions); err != nil {
		return saverResponse{}, fmt.Errorf("error decoding saver position response: %e", err)
	}
	return positions, nil
}

type poolResp struct {
	Pool                      string  `json:"pool"`
	AssetPriceUsd             float64 `json:"assetPriceUsd,string"`
	RuneOrCacaoLiquidityInUsd float64 `json:"runeOrCacaoLiquidityInUsd,string"`
}

func (l *SaverPositionResolver) getpool(pool string) (poolResp, error) {
	resp, found := l.poolCache.Get(pool)
	if found {
		return resp.(poolResp), nil
	}
	pools, err := l.fetchPools()
	if err != nil {
		return poolResp{}, err
	}
	for _, p := range pools {
		l.poolCache.Set(p.Pool, p, cache.DefaultExpiration)
	}
	return l.getpool(pool)
}
func (l *SaverPositionResolver) fetchPools() ([]poolResp, error) {
	url := fmt.Sprintf("%s/pools", l.thorwalletBaseURL)
	resp, err := http.Get(url)
	if err != nil {
		l.logger.Errorf("error fetching pools from %s: %e", url, err)
		return nil, fmt.Errorf("error fetching pools from %s: %e", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		l.logger.Errorf("error fetching pools from %s: %s", url, resp.Status)
		return nil, fmt.Errorf("error fetching pools from %s: %s", url, resp.Status)
	}
	var pools []poolResp
	if err := json.NewDecoder(resp.Body).Decode(&pools); err != nil {
		return nil, fmt.Errorf("error decoding pools response: %e", err)
	}
	return pools, nil
}
