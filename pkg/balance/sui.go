package balance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchSuiBalanceOfAddress(address string) (float64, error) {
	apiUrl := fmt.Sprintf("https://sui.api.com/v1/balance/%s", address)
	response, err := http.Get(apiUrl)
	if err != nil {
		return 0, fmt.Errorf("error fetching balance of address %s on Sui: %v", address, err)
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

	balance := data["result"].(float64)

	return balance, nil
}
