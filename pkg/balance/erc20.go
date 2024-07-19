package balance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type TokenInfo struct {
	Address  string   `json:"address"`
	Symbol   string   `json:"symbol"`
	Decimals int      `json:"decimals"`
	Name     string   `json:"name"`
	LogoURI  string   `json:"logoURI"`
	Eip2612  bool     `json:"eip2612"`
	Tags     []string `json:"tags"`
}

func GetTokenInfo(addresses []string, chain string) (map[string]TokenInfo, error) {
	addressesParam := strings.Join(addresses, ",")
	apiURL := fmt.Sprintf("https://api.vultisig.com/1inch/token/v1.2/%s/custom?addresses=%s", chain, addressesParam)

	response, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching token info: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var result map[string]TokenInfo
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response body: %v", err)
	}

	return result, nil
}

func FetchTokensWithBalance(address, chain string) (map[string]string, error) {
	apiURL := fmt.Sprintf("https://api.vultisig.com/1inch/balance/v1.2/%s/balances/%s", chain, address)

	response, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching token balances: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var balances map[string]string
	err = json.Unmarshal(body, &balances)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response body: %v", err)
	}

	nonZeroBalances := make(map[string]string)
	for address, balance := range balances {
		if balance != "0" {
			nonZeroBalances[address] = balance
		}
	}

	return nonZeroBalances, nil
}
