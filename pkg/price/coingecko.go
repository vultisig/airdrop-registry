package price

import (
	"fmt"

	"github.com/JulianToledano/goingecko"
	"github.com/vultisig/airdrop-registry/config"
)

var cgClient *goingecko.Client

func initCoinGeckoClient() {
	// if cgClient == nil {
	cgClient = goingecko.NewClient(nil, config.Cfg.CoinGecko.Key, true)
	// }
}

func fetchCoinGeckoPrice(chain, token string) (float64, error) {
	initCoinGeckoClient()

	data, err := cgClient.CoinsId(chain, true, true, true, false, false, false)
	if err != nil {
		return 0, fmt.Errorf("cgClient.CoinsId failed for chain=%s token=%s: %v", chain, token, err)
	}

	return data.MarketData.CurrentPrice.Usd, nil
}
