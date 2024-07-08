package balance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchBitcoinBalanceOfAddress(address string) (float64, error) {
	response, err := http.Get("https://mempool.space/api/address/" + address)
	if err != nil {
		return 0, fmt.Errorf("error fetching balance of address %s: %v", address, err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, fmt.Errorf("error unmarshalling response body: %v", err)
	}

	balance := data["chain_stats"].(map[string]interface{})["funded_txo_sum"].(float64)
	spent := data["chain_stats"].(map[string]interface{})["spent_txo_sum"].(float64)
	balance = balance - spent
	return balance / 100000000, nil
}
