package balance

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
)

type tronBalanceResult struct {
	Data []struct {
		Balance               int                 `json:"balance"`
		Trc20                 []map[string]string `json:"trc20"`
		LatestConsumeFreeTime int64               `json:"latest_consume_free_time"`
		NetWindowSize         int                 `json:"net_window_size"`
		NetWindowOptimized    bool                `json:"net_window_optimized"`
	} `json:"data"`
	Success bool `json:"success"`
}

func (b *BalanceResolver) FetchTronBalanceOfAddress(address, contract string, decimal int) (float64, error) {
	url := fmt.Sprintf("%s/v1/accounts/%s", b.tronBalanceBaseAddress, address)
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error fetching balance of address %s (%s) on Tron: %w", address, contract, err)
	}
	defer b.closer(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error fetching balance of address %s (%s) on Tron: %s", address, contract, resp.Status)
	}
	var result tronBalanceResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("error decoding response: %w", err)
	}
	if !result.Success || result.Data == nil || len(result.Data) == 0 {
		return 0, fmt.Errorf("failed to get balance of address %s (%s) on Tron", address, contract)
	}
	if contract == "" {
		return float64(result.Data[0].Balance) * math.Pow10(-1*decimal), nil
	}
	for i := 0; i < len(result.Data[0].Trc20); i++ {
		for k, v := range result.Data[0].Trc20[i] {
			//Tron contract address is case sensitive
			if k == contract {
				value, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return 0, fmt.Errorf("error parsing trc20 balance: %w", err)
				}
				return float64(value) * math.Pow10(-1*decimal), nil
			}
		}
	}
	return 0, nil
}
