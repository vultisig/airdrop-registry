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

func fetchCoinGeckoTokenPrice(chain, contract string) (float64, error) {
	initCoinGeckoClient()

	data, err := cgClient.SimpleTokenPrice(chain, contract, "usd", false, false, false, false)
	if err != nil {
		return 0, fmt.Errorf("cgClient.SimpleTokenPrice failed for chain=%s contract=%s: %v", chain, contract, err)
	}

	if data == nil {
		return 0, fmt.Errorf("cgClient.SimpleTokenPrice failed for chain=%s contract=%s: data is nil", chain, contract)
	}

	tokenData, ok := data[contract]
	if !ok {
		return 0, fmt.Errorf("no data found for contract=%s on chain=%s", contract, chain)
	}

	price, ok := tokenData["usd"]
	if !ok {
		return 0, fmt.Errorf("no USD price found for contract=%s on chain=%s", contract, chain)
	}

	return price, nil
}
