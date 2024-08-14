package balance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/vultisig/airdrop-registry/internal/common"
)

func (b *BalanceResolver) closer(closer io.Closer) {
	if err := closer.Close(); err != nil {
		b.logger.Error(err)
	}
}

type UtxoResult struct {
	Data map[string]struct {
		Address struct {
			Balance    float64 `json:"balance"`
			BalanceUSD float64 `json:"balance_usd"`
		} `json:"address"`
	} `json:"data"`
}

// FetchUtxoBalanceOfAddress fetches the UTXO balance of an address and it's USD value
func (b *BalanceResolver) FetchUtxoBalanceOfAddress(address string, chain common.Chain) (float64, float64, error) {
	if address == "" {
		return 0, 0, fmt.Errorf("address cannot be empty")
	}
	var chainName string
	switch chain {
	case common.Bitcoin:
		chainName = "bitcoin"
	case common.BitcoinCash:
		chainName = "bitcoin-cash"
	case common.Dash:
		chainName = "dash"
	case common.Litecoin:
		chainName = "litecoin"
	case common.Dogecoin:
		chainName = "dogecoin"
	default:
		return 0, 0, fmt.Errorf("unsupported chain: %s", chain)
	}
	url := fmt.Sprintf("%s/blockchair/%s/dashboards/address/%s?state=latest", vultisigApiProxy, chainName, address)
	response, err := http.Get(url)
	if err != nil {
		return 0, 0, fmt.Errorf("error fetching UTXO balance of address %s: %w", address, err)
	}

	defer b.closer(response.Body)
	if response.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("error fetching UTXO balance of address %s: %s", address, response.Status)
	}
	var result UtxoResult
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return 0, 0, fmt.Errorf("error unmarshalling response: %w", err)
	}
	data, ok := result.Data[address]
	if !ok {
		return 0, 0, fmt.Errorf("address data not found in response")
	}

	return data.Address.Balance, data.Address.BalanceUSD, nil

}
