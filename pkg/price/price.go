package price

import (
	"fmt"
)

func FetchPrice(chain, token string) (float64, error) {
	// if token starts with 0x and length is 42, then it's a contract address
	if len(token) == 42 && token[:2] == "0x" {
		price, err := fetchCoinGeckoTokenPrice(chain, token)
		if err != nil {
			return 0, fmt.Errorf("fetchCoinGeckoTokenPrice failed: %v", err)
		}
		return price, nil
	}

	price, err := fetchCoinGeckoPrice(chain, token)
	if err != nil {
		return 0, fmt.Errorf("fetchCoinGeckoPrice failed: %v", err)
	}

	return price, nil
}
