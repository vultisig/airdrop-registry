package balance

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type tonBalanceResult struct {
	Balance uint64 `json:"balance,string"`
}

func (b *BalanceResolver) FetchTonBalanceOfAddress(address string) (float64, error) {
	url := fmt.Sprintf("%s?address=%s&use_v2=false", b.TonBalanceBaseAddress, address)
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error fetching balance of address %s on TON: %w", address, err)
	}
	defer b.closer(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error fetching balance of address %s on TON: %s", address, resp.Status)
	}
	var result tonBalanceResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("error decoding response: %w", err)
	}
	return float64(result.Balance) * 1e-9, nil
}
