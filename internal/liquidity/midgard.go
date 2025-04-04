package liquidity

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

type TcMayaPoolPositionResolver struct {
	baseAddress string
	usdPools    []string
	poolCache   *cache.Cache
	logger      *logrus.Logger
}

func NewMidgardLPResolver(baseAddress string, usdPools []string) *TcMayaPoolPositionResolver {
	return &TcMayaPoolPositionResolver{
		baseAddress: baseAddress,
		usdPools:    usdPools,
		poolCache:   cache.New(3*time.Hour, 6*time.Hour),
		logger:      logrus.WithField("module", "midgard_lp_resolver").Logger,
	}
}

type pool struct {
	Asset         string  `json:"asset"`
	AssetDepth    string  `json:"assetDepth"`
	AssetPrice    float64 `json:"assetPrice,string"`
	AssetPriceUSD float64 `json:"assetPriceUSD,string"`
}
type member struct {
	AssetDeposit float64 `json:"assetDeposit,string"`
	Pool         string  `json:"pool"`
	RuneDeposit  int64   `json:"runeDeposit,string"`
	CacaoDeposit int64   `json:"cacaoDeposit,string"`
}

type members struct {
	Pools []member `json:"pools"`
}

func (tcm *TcMayaPoolPositionResolver) GetLiquidityPosition(address string) (float64, error) {
	var totalLiquidity float64
	members, err := tcm.GetMemberPosition(address)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch member from %s: %w", tcm.baseAddress, err)
	}
	if members.Pools == nil {
		return 0, nil
	}
	nativeTokenPrice, err := tcm.getNativeTokenPrice()
	if err != nil {
		return 0, fmt.Errorf("failed to get native token price: %w", err)
	}

	for _, memberPool := range members.Pools {
		assetPrice, err := tcm.getAssetPrice(memberPool.Pool)
		if err != nil {
			tcm.logger.WithError(err).WithField("asset", memberPool.Pool).Error("failed to get asset price")
			return 0, fmt.Errorf("failed to get asset price for pool %s: %w", memberPool.Pool, err)
		}

		totalLiquidity += float64(memberPool.CacaoDeposit+memberPool.RuneDeposit)*nativeTokenPrice +
			memberPool.AssetDeposit*assetPrice
	}

	return totalLiquidity, nil
}

func (tcm *TcMayaPoolPositionResolver) GetMemberPosition(address string) (members, error) {
	if address == "" {
		return members{}, nil
	}
	url := fmt.Sprintf("%s/v2/member/%s", tcm.baseAddress, address)
	resp, err := http.Get(url)
	if err != nil {
		return members{}, fmt.Errorf("error fetching member from %s: %v", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return members{}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return members{}, fmt.Errorf("error fetching member from %s: %s", url, resp.Status)
	}
	var members members
	if err := json.NewDecoder(resp.Body).Decode(&members); err != nil {
		return members, fmt.Errorf("error decoding member response: %w", err)
	}
	return members, nil
}

func (tcm *TcMayaPoolPositionResolver) refreshCache() error {
	url := fmt.Sprintf("%s/v2/pools", tcm.baseAddress)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching pools from %s: %v", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error fetching pools from %s: %s", url, resp.Status)
	}
	var pools []pool
	if err := json.NewDecoder(resp.Body).Decode(&pools); err != nil {
		return fmt.Errorf("error decoding pools response: %w", err)
	}
	for _, p := range pools {
		tcm.poolCache.Set(p.Asset, p, cache.DefaultExpiration)
	}
	return nil
}

func (tcm *TcMayaPoolPositionResolver) getAssetPrice(asset string) (float64, error) {
	if cached, ok := tcm.poolCache.Get(asset); ok {
		if pool, ok := cached.(pool); ok {
			return pool.AssetPriceUSD, nil
		}
	}
	err := tcm.refreshCache()
	if err != nil {
		return 0, fmt.Errorf("failed to refresh pool cache: %w", err)
	}
	if cached, ok := tcm.poolCache.Get(asset); ok {
		if pool, ok := cached.(pool); ok {
			return pool.AssetPriceUSD, nil
		}
	}
	return 0, nil
}

func (tcm *TcMayaPoolPositionResolver) getNativeTokenPrice() (float64, error) {
	var nativePriceSum float64
	var count int
	for _, usdpool := range tcm.usdPools {
		if cached, ok := tcm.poolCache.Get(usdpool); ok {
			if pool, ok := cached.(pool); ok {
				if pool.AssetPrice > 0 && pool.AssetPriceUSD > 0 {
					nativePriceSum += pool.AssetPriceUSD / pool.AssetPrice
					count++
				}
			}
		}
	}
	if count > 0 {
		return nativePriceSum / float64(count), nil
	}

	err := tcm.refreshCache()
	if err != nil {
		return 0, fmt.Errorf("failed to refresh pool cache: %w", err)
	}
	nativePriceSum = 0
	count = 0
	for _, usdpool := range tcm.usdPools {
		if cached, ok := tcm.poolCache.Get(usdpool); ok {
			if pool, ok := cached.(pool); ok {
				if pool.AssetPrice > 0 && pool.AssetPriceUSD > 0 {
					nativePriceSum += pool.AssetPriceUSD / pool.AssetPrice
					count++
				}
			}
		}
	}
	if count > 0 {
		return nativePriceSum / float64(count), nil
	}
	return 0, nil
}
