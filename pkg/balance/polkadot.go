package balance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/vultisig/airdrop-registry/pkg/utils"
)

func FetchPolkadotBalanceOfAddress(address string) (float64, error) {
	payload := fmt.Sprintf(`{"jsonrpc":"2.0","method":"api.rpc.eth_getBalance","params":["%s"],"id":1}`, address)
	response, err := http.Post("https://rpc.polkadot.io", "application/json", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return 0, fmt.Errorf("error fetching balance of address %s on Polkadot: %v", address, err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	fmt.Println(string(body))
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, fmt.Errorf("error unmarshalling response body: %v", err)
	}

	// Extract balance from the response
	balanceHex, ok := data["result"].(string)
	if !ok {
		return 0, fmt.Errorf("unexpected response format")
	}

	// Convert hex balance to float64
	balance, err := utils.HexToFloat64(balanceHex)
	if err != nil {
		return 0, fmt.Errorf("error converting balance to float: %v", err)
	}

	return balance, nil
}
