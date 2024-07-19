package balance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func FetchSuiBalanceOfAddress(address string) (float64, error) {
	chain := "sui"
	rpcUrl := getRpcUrlForChain(chain)

	if rpcUrl == "" {
		return 0, fmt.Errorf("unsupported chain: %s", chain)
	}

	payload := fmt.Sprintf(`{"jsonrpc":"2.0","method":"suix_getBalance","params":["%s"],"id":1}`, address)
	response, err := http.Post(rpcUrl, "application/json", strings.NewReader(payload))
	if err != nil {
		return 0, fmt.Errorf("error fetching balance of address %s on %s: %v", address, chain, err)
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

	totalBalance, ok := result["totalBalance"].(string)
	if !ok {
		return 0, fmt.Errorf("unexpected response format for totalBalance")
	}

	balance, err := strconv.ParseFloat(totalBalance, 64)
	if err != nil {
		return 0, fmt.Errorf("error converting balance to float: %v", err)
	}

	balance = balance / 1e9

	return balance, nil
}
