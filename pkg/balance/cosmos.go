package balance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func FetchThorchainBalanceOfAddress(address string) (float64, error) {
	url := fmt.Sprintf("https://thornode.ninerealms.com/cosmos/bank/v1beta1/balances/%s", address)
	return fetchSpecificCosmosBalance(url, "rune")
}

func FetchMayachainBalanceOfAddress(address string) (float64, error) {
	url := fmt.Sprintf("https://mayanode.mayachain.info/cosmos/bank/v1beta1/balances/%s", address)
	return fetchSpecificCosmosBalance(url, "cacao")
}

func fetchSpecificCosmosBalance(url, denom string) (float64, error) {
	response, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error fetching balance from %s: %v", url, err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	var data struct {
		Balances []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"balances"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, fmt.Errorf("error unmarshalling response body: %v", err)
	}

	var balance float64
	if len(data.Balances) == 0 {
		return 0, nil
	}
	for _, b := range data.Balances {
		if b.Denom == denom {
			// divide balance by 8
			balance, err = strconv.ParseFloat(b.Amount, 64)
			// fmt.Println(balance)
			if err != nil {
				return 0, fmt.Errorf("error converting balance to float: %v", err)
			}
			break
		}
	}

	return balance / 1e8, nil
}
