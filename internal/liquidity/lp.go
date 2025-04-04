package liquidity

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

type LiquidityPositionResolver struct {
	logger            *logrus.Logger
	thorwalletBaseURL string
	tclpResolver      *TcMayaPoolPositionResolver
	mayalpResolver    *TcMayaPoolPositionResolver
	wewelpResolver    *weweLpResolver
	tgtPrice          float64
	mu                sync.RWMutex
	cachedData        *cache.Cache
}

func NewLiquidtyPositionResolver() *LiquidityPositionResolver {
	return &LiquidityPositionResolver{
		logger:         logrus.WithField("module", "liquidity_position_resolver").Logger,
		tclpResolver:   NewMidgardLPResolver("https://midgard.ninerealms.com", []string{"ETH.USDT-0XDAC17F958D2EE523A2206206994597C13D831EC7", "AVAX.USDC-0XB97EF9EF8734C71904D8002F8B6BC66DD9C48A6E", "ETH.USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48"}),
		mayalpResolver: NewMidgardLPResolver("https://midgard.mayachain.info", []string{"ETH.USDT-0XDAC17F958D2EE523A2206206994597C13D831EC7", "ETH.USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48"}),
		wewelpResolver: NewWeWeLpResolver(),
		cachedData:     cache.New(3*time.Hour, 6*time.Hour),
	}
}

// fetch Thorchain and Maya LP position from Thorwallet api
func (l *LiquidityPositionResolver) GetLiquidityPosition(address string) (float64, error) {
	if address == "" {
		return 0, fmt.Errorf("addresse cannot be empty")
	}

	tclp, err := l.tclpResolver.GetLiquidityPosition(address)
	if err != nil {
		return 0, fmt.Errorf("error getting TC liquidity position: %w", err)
	}
	mayalp, err := l.mayalpResolver.GetLiquidityPosition(address)
	if err != nil {
		return 0, fmt.Errorf("error getting Maya liquidity position: %w", err)
	}
	totalLiquidity := tclp + mayalp
	return totalLiquidity, nil
}

type tgtLPPositionResponse struct {
	StakeAmount float64 `json:"stakedAmount,string"`
	Reward      float64 `json:"reward,string"`
}

func (l *LiquidityPositionResolver) GetTGTStakePosition(addresses string) (float64, error) {
	if addresses == "" {
		return 0, nil
	}

	url := fmt.Sprintf("%s/stake/%s", l.thorwalletBaseURL, addresses)
	resp, err := http.Get(url)
	if err != nil {
		l.logger.Errorf("error fetching stake position from %s: %e", url, err)
		return 0, fmt.Errorf("error fetching stake position from %s: %e", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		l.logger.Errorf("error fetching stake position from %s: %s", url, resp.Status)
		return 0, fmt.Errorf("error fetching stake position from %s: %s", url, resp.Status)
	}
	var positions tgtLPPositionResponse
	if err := json.NewDecoder(resp.Body).Decode(&positions); err != nil {
		return 0, fmt.Errorf("error decoding stake position response: %e", err)
	}
	return positions.StakeAmount*l.GetTGTPrice() + positions.Reward, nil
}

func (l *LiquidityPositionResolver) SetTGTPrice(price float64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.tgtPrice = price
}

func (l *LiquidityPositionResolver) GetTGTPrice() float64 {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.tgtPrice
}

func (l *LiquidityPositionResolver) GetWeWeLPPosition(address string) (float64, error) {
	return l.wewelpResolver.GetLiquidityPosition(address)
}
