package liquidity

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/sirupsen/logrus"
)

type LiquidityPositionResolver struct {
	logger            *logrus.Logger
	thorwalletBaseURL string
	thornodeBaseURL   string
	tcyPrice          float64
	midgardBaseURL    string
	mu                sync.RWMutex
}

func NewLiquidtyPositionResolver() *LiquidityPositionResolver {
	return &LiquidityPositionResolver{
		logger:            logrus.WithField("module", "liquidity_position_resolver").Logger,
		thorwalletBaseURL: "https://api-v2-prod.thorwallet.org",
		thornodeBaseURL:   "https://thornode.ninerealms.com",
		midgardBaseURL:    "https://midgard.ninerealms.com/v2",
	}
}

type poolPositionResponse struct {
	RuneOrCacaoAddedUsd float64 `json:"runeOrCacaoAddedUsd,string"`
	AssetAddedUsd       float64 `json:"assetAddedUsd,string"`
}

// fetch Thorchain and Maya LP position from Thorwallet api
func (l *LiquidityPositionResolver) GetLiquidityPosition(address string) (float64, error) {
	if address == "" {
		return 0, fmt.Errorf("address cannot be empty")
	}
	url := fmt.Sprintf("%s/pools/positions?addresses=%s", l.thorwalletBaseURL, address)
	resp, err := http.Get(url)
	if err != nil {
		l.logger.Errorf("error fetching liquidity position from %s: %e", url, err)
		return 0, fmt.Errorf("error fetching liquidity position from %s: %e", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		l.logger.Errorf("error fetching liquidity position from %s: %s", url, resp.Status)
		return 0, fmt.Errorf("error fetching liquidity position from %s: %s", url, resp.Status)
	}
	var positions map[string][]poolPositionResponse
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading liquidity position response: %e", err)
	}
	l.logger.Infof("response(%s) : %s", url, string(buf))
	if err := json.Unmarshal(buf, &positions); err != nil {
		return 0, fmt.Errorf("error decoding liquidity position response: %e", err)
	}
	var totalLiquidity float64
	if positions == nil {
		l.logger.Errorf("no liquidity position found for address %s", address)
		return 0, nil
	}
	for _, v := range positions {
		for _, p := range v {
			totalLiquidity += p.RuneOrCacaoAddedUsd + p.AssetAddedUsd
		}
	}
	return totalLiquidity, nil
}

// fetch Thorchain TCY LP position from thornode api
func (l *LiquidityPositionResolver) GetTCYStakePosition(address string) (float64, error) {
	if address == "" {
		return 0, nil
	}
	url := fmt.Sprintf("%s/thorchain/tcy_staker/%s", l.thornodeBaseURL, address)
	resp, err := http.Get(url)
	if err != nil {
		l.logger.Errorf("error fetching liquidity position from %s: %e", url, err)
		return 0, fmt.Errorf("error fetching liquidity position from %s: %e", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			l.logger.Warnf("bad request for address %s, possibly not a TCY staker", address)
			return 0, nil
		}
		l.logger.Errorf("error fetching liquidity position from %s: %s", url, resp.Status)
		return 0, fmt.Errorf("error fetching liquidity position from %s: %s", url, resp.Status)
	}
	var lp tcyLPPositionResponse
	if err := json.NewDecoder(resp.Body).Decode(&lp); err != nil {
		l.logger.WithError(err).Error("Failed to decode response")
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return float64(lp.Amount) * l.GetTCYPrice(), nil
}

type tcyLPPositionResponse struct {
	Address string `json:"address"`
	Amount  int64  `json:"amount,string"`
}

func (l *LiquidityPositionResolver) SetTCYPrice(price float64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.tcyPrice = price

}
func (l *LiquidityPositionResolver) GetTCYPrice() float64 {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.tcyPrice
}
