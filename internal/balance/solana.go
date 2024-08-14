package balance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchSolanaBalanceOfAddress(address string) (float64, error) {
	payload := fmt.Sprintf(`{"jsonrpc":"2.0","method":"getBalance","params":["%s"],"id":1}`, address)
	response, err := http.Post("https://api.mainnet-beta.solana.com", "application/json", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return 0, fmt.Errorf("error fetching balance of address %s on Solana: %v", address, err)
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

	result, ok := data["result"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("unexpected response format")
	}

	value, ok := result["value"].(float64)
	if !ok {
		return 0, fmt.Errorf("unexpected response format")
	}

	return value / 1000000000, nil // convert lamports to SOL
}
