package balance

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func (b *BalanceResolver) FetchThorchainBalanceOfAddress(address string) (float64, error) {
	if address == "" {
		return 0, fmt.Errorf("address cannot be empty")
	}
	url := fmt.Sprintf("https://thornode.ninerealms.com/cosmos/bank/v1beta1/balances/%s", address)
	return b.fetchSpecificCosmosBalance(url, "rune", 8)
}

func (b *BalanceResolver) FetchMayachainBalanceOfAddress(address string) (float64, error) {
	url := fmt.Sprintf("https://mayanode.mayachain.info/cosmos/bank/v1beta1/balances/%s", address)
	return b.fetchSpecificCosmosBalance(url, "cacao", 10)
}
func (b *BalanceResolver) FetchCosmosBalanceOfAddress(address string) (float64, error) {
	url := fmt.Sprintf("https://cosmos-rest.publicnode.com/cosmos/bank/v1beta1/balances/%s", address)
	return b.fetchSpecificCosmosBalance(url, "uatom", 6)
}

func (b *BalanceResolver) FetchKujiraBalanceOfAddress(address string) (float64, error) {
	url := fmt.Sprintf("https://kujira-rest.publicnode.com/cosmos/bank/v1beta1/balances/%s", address)
	return b.fetchSpecificCosmosBalance(url, "ukuji", 6)
}

func (b *BalanceResolver) FetchDydxBalanceOfAddress(address string) (float64, error) {
	url := fmt.Sprintf("https://dydx-rest.publicnode.com/cosmos/bank/v1beta1/balances/%s", address)
	return b.fetchSpecificCosmosBalance(url, "adydx", 18)
}

type CosmosData struct {
	Balances []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"balances"`
}

func (b *BalanceResolver) fetchSpecificCosmosBalance(url, denom string, decimals int) (float64, error) {
	if url == "" {
		return 0, fmt.Errorf("url cannot be empty")
	}
	if denom == "" {
		return 0, fmt.Errorf("denom cannot be empty")
	}
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error fetching balance from %s: %w", url, err)
	}
	defer b.closer(resp.Body)
	if resp.StatusCode == http.StatusTooManyRequests {
		// rate limited, need to backoff and then retry
		return 0, ErrRateLimited
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error fetching balance from %s: %s", url, resp.Status)
	}
	var result CosmosData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("error unmarshalling response: %w", err)
	}

	var balance float64
	if len(result.Balances) == 0 {
		return 0, nil
	}
	for _, b := range result.Balances {
		if strings.EqualFold(b.Denom, denom) {
			balance, err = strconv.ParseFloat(b.Amount, 64)
			if err != nil {
				return 0, fmt.Errorf("error converting balance to float: %v", err)
			}
			break
		}
	}

	return balance / math.Pow10(decimals), nil
}
