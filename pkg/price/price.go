package price

import (
	"fmt"
)

func FetchPrice(chain, token string) (float64, error) {
	price, err := fetchCoinGeckoPrice(chain, token)
	if err != nil {
		return 0, fmt.Errorf("fetchCoinGeckoPrice failed: %v", err)
	}

	return price, nil
}
